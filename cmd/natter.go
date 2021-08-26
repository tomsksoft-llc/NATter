package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"NATter/config"
	"NATter/driver"
	"NATter/driver/http"
	"NATter/driver/msgbroker/kafka"
	"NATter/driver/msgbroker/nats"
	"NATter/log"
	"NATter/service"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const brokerDriverName = "broker"

type NATter struct {
	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc

	conns  map[string]driver.Conn
	router *service.Router
}

func NewNATter() *NATter {
	natter := &NATter{
		wg:    sync.WaitGroup{},
		conns: map[string]driver.Conn{},
	}

	natter.ctx, natter.cancel = context.WithCancel(context.Background())

	return natter
}

func (natter *NATter) Run() error {
	env := pflag.String("env", "devel.toml", "environment name")
	pflag.Parse()

	path := natter.specifyRelativeEnvPath(*env)

	if err := config.Load(path); err != nil {
		return err
	}

	if err := natter.connect(); err != nil {
		return err
	}

	defer natter.close()

	if err := natter.start(); err != nil {
		return err
	}

	natter.waitShutdown()

	return nil
}

func (natter *NATter) specifyRelativeEnvPath(path string) string {
	if !filepath.IsAbs(path) {
		return fmt.Sprintf(".%ccfg%c%s", os.PathSeparator, os.PathSeparator, path)
	}

	return path
}

func (natter *NATter) connect() error {
	if err := natter.setupLogger(); err != nil {
		return err
	}

	return natter.setupConns()
}

func (natter *NATter) close() {
	for _, conn := range natter.conns {
		if err := conn.Close(); err != nil {
			log.Error(err)
		}
	}
}

func (natter *NATter) setupLogger() error {
	err := log.SetConfig(&log.Config{
		Level: config.String("LOG.LOGGER_LEVEL"),
		Stream: log.ConfigStream{
			Enable: config.Bool("LOG.STREAMLOG_ENABLE"),
		},
		File: log.ConfigFile{
			Enable: config.Bool("LOG.FILELOG_ENABLE"),
			Path:   config.String("LOG.FILELOG_PATH"),
		},
		Syslog: log.ConfigSyslog{
			Name:   "NATter",
			Enable: config.Bool("LOG.SYSLOG_ENABLE"),
			Host:   config.String("LOG.SYSLOG_HOST"),
		},
	})

	return err
}

func (natter *NATter) setupConns() error {
	conn, err := natter.bootBroker()

	if err != nil {
		return err
	}

	natter.conns[brokerDriverName] = conn

	natter.conns[http.DriverName] = http.NewConn(&http.ConnConfig{
		Port: config.String("HTTP.PORT"),
	})

	return nil
}

func (natter *NATter) bootBroker() (conn driver.Conn, err error) {
	broker := config.String("MESSAGE_BROKER.BROKER")

	switch broker {
	case nats.DriverName:
		conn, err = nats.NewConn(&nats.ConnConfig{
			Servers: config.StringSlice("MESSAGE_BROKER.NATS_SERVERS"),
			Token:   config.String("MESSAGE_BROKER.NATS_TOKEN"),
			Group:   config.String("MESSAGE_BROKER.SERVICE_GROUP"),
			Name:    config.String("MESSAGE_BROKER.SERVICE_NAME"),
		})
	case kafka.DriverName:
		conn, err = kafka.NewConn(&kafka.ConnConfig{
			Version: config.String("MESSAGE_BROKER.KAFKA_VERSION"),
			Servers: config.StringSlice("MESSAGE_BROKER.KAFKA_SERVERS"),
			Group:   config.String("MESSAGE_BROKER.SERVICE_GROUP"),
		})
	default:
		err = errors.Errorf("unknown message broker: %s", broker)
	}

	return conn, err
}

func (natter *NATter) start() (err error) {
	log.Info("Starting NATter")

	if err := natter.bootRouter(); err != nil {
		return err
	}

	natter.listenOS()

	natter.router.Run(natter.ctx)

	return nil
}

func (natter *NATter) bootRouter() (err error) {
	routes, err := config.Routes()

	if err != nil {
		return err
	}

	natter.router, err = service.NewRouter(&service.RouterConfig{
		Routes: routes,
	}, natter.conns)

	return err
}

func (natter *NATter) waitShutdown() {
	defer log.Info("Ending NATter")

	<-natter.ctx.Done()
	natter.wg.Wait()
}

func (natter *NATter) listenOS() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigchan
		natter.cancel()
	}()
}
