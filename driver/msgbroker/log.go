package msgbroker

import (
	"NATter/log"
)

func LogErrorHandle(err error, topic string) {
	ent := log.WithFields(log.Fields{
		"topic": topic,
	})

	ent.WithError(err).Error("unable handle message")
}

func LogDebugSubscribed(topic string) {
	log.WithFields(log.Fields{
		"topic": topic,
	}).Debug("subscribed")
}

func LogDebugUnsubscribed(topic string) {
	log.WithFields(log.Fields{
		"topic": topic,
	}).Debug("unsubscribed")
}

func LogDebugPublished(topic string, payload interface{}) {
	ent := log.WithFields(log.Fields{
		"topic": topic,
	})

	if payload != nil {
		ent.WithFields(log.Fields{
			"payload": log.FormatStruct(payload),
		})
	}

	ent.Debug("published")
}

func LogDebugReceived(topic string, payload interface{}) {
	ent := log.WithFields(log.Fields{
		"topic": topic,
	})

	if payload != nil {
		ent.WithFields(log.Fields{
			"payload": log.FormatStruct(payload),
		})
	}

	ent.Debug("received")
}

func LogDebugRequested(topic string, req, resp interface{}) {
	ent := log.WithFields(log.Fields{
		"topic": topic,
	})

	if req != nil {
		ent.WithFields(log.Fields{
			"request": log.FormatStruct(req),
		})
	}

	if resp != nil {
		ent.WithFields(log.Fields{
			"response": log.FormatStruct(resp),
		})
	}

	ent.Debug("requested")
}

func LogDebugResponded(topic string, resp interface{}) {
	ent := log.WithFields(log.Fields{
		"topic": topic,
	})

	if resp != nil {
		ent.WithFields(log.Fields{
			"payload": log.FormatStruct(resp),
		})
	}

	ent.Debug("responded")
}
