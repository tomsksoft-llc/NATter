# NATter

![cover](cover.png)

NATter is a Go based proxy service for connecting a website with message brokers. It allows you to read messages from a queue and pass them to your website as a webhook and backwards.

**Table of Contents**

* [Running NATter](#running-natter)
* [Quick start](#quick-start)
* [Running tests](#running-tests)
* [Configuration](#configuration)
  * [LOG section](#log-section)
  * [MESSAGE_BROKER section](#message-broker-section)
  * [HTTP section](#http-section)
  * [ROUTES section](#routes-section)
* [Route mode](#route-mode)
* [Batching](#batching)
* [Custom URIs' specialties](#custom-uris-specialties)
* [Messaging drivers](#messaging-drivers)
  * [Implementation](#implementation)
  * [Injection](#injection)
* [License](#license)

## Running NATter
Run ```{natter_folder}/natter --env={config_path}``` using a command line.

The ```{config_path}``` can be either absolute or relative. If it is relative, it must point at a config file regarding the ```cfg``` subdirectory of a directory where a binary file is.

## Quick start
We are going to read messages from the ```user.create``` topic and write messages into the ```user.login``` in case when user logins into the website. Here is a minimal config:
```
[LOG]
LOGGER_LEVEL='debug'

[HTTP]
PORT='3000'

[MESSAGE_BROKER]
BROKER='nats'
NATS_SERVERS=['nats://127.0.0.1:4222']

[[ROUTES]]
MODE='broker-http-oneway'
TOPIC='user.create'
ENDPOINT='http://127.0.0.1/user.php'

[[ROUTES]]
MODE='http-broker-oneway'
TOPIC='user.login'
URI='/user/login'
```

 Note that the ```MESSAGE_BROKER``` section from the config above is set according to NATS so you might set it according to Kafka if you are going to work with it:
```
[MESSAGE_BROKER]
BROKER='kafka'
KAFKA_VERSION='2.8.0.0'
KAFKA_SERVERS=['127.0.0.1:9092']
```

The next thing you need to do is to host the message broker (```nats://127.0.0.1:4222``` if NATS or ```127.0.0.1:9092``` if Kafka) with the appropriate settings. Host your PHP web server (```127.0.0.1/user.php```) with the following script:
```
<?php

$payload = json_decode(file_get_contents("php://input"), true));

var_dump($payload);
```

and then run NATter.

After that publish the following data to your broker on the ```user.create``` topic:
```
{
  "id": 1,
  "name": "John Wick",
  "date":"2021-27-09 12:34:56"
}
```

Your web server will output the following:
```
// Output:
// array(3) {
//   ["id"]=>
//   int(1)
//   ["name"]=>
//   string(4) "John Wick"
//   ["datetime"]=>
//   string(19) "2021-27-09 12:34:56"
// }
```

Then to write a message on the ```user.login``` topic to your broker we need to run the following PHP script:
```
<?php

use GuzzleHttp;

$client = new GuzzleHttp\Client();
$client->post('http://127.0.0.1:3000/user/login', [
    'json' => [
        'id'       => 1,
        'datetime' => '2021-27-09 23:46:07'
    ]
]);
```

A subscriber will recieve the following data:
```
{
  "id": 1,
  "date":"2021-27-09 23:46:07"
}
```

Remember NATter is only proxy so you can choose any message encoding format for communication between system's pieces.

## Running tests
NATter provides integration testing for a config parser and a Kafka driver. It is configured by a testing file that must have the ```cfg/test.toml``` path. It has options to enable/disable testing and those ones that apply to a process directly.

## Configuration
A TOML config file includes the following sections of options:
 * **LOG** describes logging options.
 * **MESSAGE_BROKER** describes message broker options.
 * **HTTP** describes HTTP options.
 * **ROUTES** describes routing options.

There is the config file example settings in ```cfg/example.toml```.

### LOG section
The section includes the following options:
 * **LOGGER_LEVEL** defines of what level (or above) logs should be written. Possible values: ```trace```, ```debug```, ```info```, ```error```, ```fatal```, ```panic```.
 * **STREAMLOG_ENABLE** enables/disables logging to the standard output stream. Possible values: ```true```, ```false```. Default: ```false```.
 * **FILELOG_ENABLE** enables/disables logging to a file. Possible values: ```true```, ```false```. Default: ```false```.
 * **FILELOG_PATH** is a filepath where logs should be written. Could be absolute or relative.
 * **SYSLOG_ENABLE** enables/disables system logging. Possible values: ```true```, ```false```. Default: ```false```.
 * **SYSLOG_HOST** is a system log service host.

### MESSAGE_BROKER section
The section includes the following options:
 * **BROKER** is one of message brokers that should take part in messaging. Possible values: ```nats```, ```kafka```.
 * **NATS_SERVERS** is an array of addresses of a NATS cluster. It is parsed only if the ```BROKER``` is set to the ```nats```.
 * **NATS_TOKEN** is a token for access to the NATS cluster. It is parsed only if the ```BROKER``` is set to the ```nats```.
 * **KAFKA_VERSION** is used for internal library initialization, should reflect to your Kafka version. It has a string format of four numbers delimited by a dot, e.g. ```'1.2.3.4'```.
 * **KAFKA_SERVERS** is an array of addresses of a Kafka cluster. It is parsed only if the ```BROKER``` is set to the ```kafka```.
 * **SERVICE_GROUP** is a NATter service's queue group name within the message broker.
 * **SERVICE_NAME** is a NATter service's client name within the message broker.

### HTTP section
The section includes the following options:
 * **HOST** defines a host that the NATter service binds to listen and serve API and routing requests. Default: ```'0.0.0.0'```.
 * **PORT** defines a port that the NATter service binds to listen and serve API and routing requests.

The service supports IPv6.

### ROUTES section
The section represents an array of structures that describe a route that is wanted to be registered in a NATter instance. The route structure consists of the following fields:
 * required ones:
   * **MODE** is a rule by which the route transfers messages. Possible values: ```broker-http-oneway```, ```broker-http-twoway```, ```http-broker-oneway```, ```http-broker-twoway```.
   * **TOPIC** is the message broker topic from/to which the message is routed to/from a web.
 * optional ones that are set only under specific conditions depending on the route mode and some preferences:
   * **ENDPOINT** is an address to which the message is routed from the message broker ```TOPIC```.
   * **URI** is an HTTP path from which the message is routed to the message broker ```TOPIC```.
   * **ASYNC** enables/disables a route's asynchrony if the mode is the ```http-broker-twoway```. Possible values: ```true```, ```false```.  Default: ```false```.
 * **ROUTES.BATCHING** subsection options:
   * **TIMEOUT** is a frequency (in seconds) of a batch release. Default: ```0```.
   * **CAPACITY** is a number of messages in the batch to release it before the timeout expires. Default: ```0```.

## Route mode
The route mode is a rule that defines a logic of how data is proxied. Whether an HTTP or a message broker (Broker) client should be an initiator of interaction depends on the route mode. It has the 'source-recipient-direction' format where the 'source' is a connection where data is received from, the 'recipient' is a connection where data is sent to and the 'direction' (possible values: ```oneway```, ```twoway```) defines whether a response should be sent back to the initiator.

There is a few currently available modes:
 * **broker-http-oneway** means that a request from Broker should be proxied to HTTP without receiving the response (Broker -> HTTP).
 * **broker-http-twoway** means that the request from Broker should be proxied to HTTP which should deliver the response that should be proxied back to the same topic to Broker (Broker -> HTTP -> Broker).
 * **http-broker-oneway** means that the request from HTTP should be proxied to Broker without receiving the response (HTTP -> Broker).
 * **http-broker-twoway** means that the request from HTTP should be proxied to Broker which should deliver the response that should be proxied back (HTTP -> Broker -> HTTP). Depending on the ```ASYNC``` flag value the response should be proxied either synchronously or asynchronously.

Note that Broker can be either NATS or Kafka within one NATter session.

## Batching
It is possible to set the batching for the route. Batching options include the following parameters:
 * **timeout** is time after which a batch is released. The default value is 0 that means there is no timeout.
 * **capacity** is a maximum number of the messages to release the batch. The default is 0 that means the capacity is not taken into account.

_Note_: It is required to set at least one of these parameters to the non default value since otherwise the batch is accumulated but never being sent. It is also recommended to set the timeout to the non default value since otherwise there is a risk of the batch being pending for a long time.

## Custom URIs' specialties
URIs of the ```/i/*``` type reserved for the HTTP API so this type can not be used for declaration of custom URIs.

## API
NATter provides [REST HTTP API](https://github.com/tomsksoft-llc/NATter/blob/dev/driver/http/doc/api.yaml) to monitor routes registered in a specific session. This document can be opened in Swagger Editor: https://editor.swagger.io/.

## Messaging drivers
### Implementation
NATter suggests an opportunity to complement it by writing the new drivers for the messaging by other protocols. For these purposes NATter has the ```driver``` package at the top level of the folder structure that includes some interfaces and subpackages with their implementations (e.g. http, msgbroker). There are the following interfaces provided by NATter:
```
type Conn interface {
	Serve(context.Context) error
	Close() error
	Receiver(*entity.Route) Receiver
	Sender(*entity.Route) Sender
}

type Receiver interface {
	Listen(Sender) error
	ListenRequest(Sender) error
}

type Sender interface {
	Send(payload []byte) error
	Request(payload []byte) ([]byte, error)
}
```

The ```Conn``` interface is the base of the protocol messaging and provides methods to control a connection (```Serve()``` and ```Close()```) and building corresponding contact entities that provide the messaging directly (```Receiver()``` and ```Sender()```) by rules described in an instance of the ```*entity.Route``` located in the ```entity/route.go```:
```
type Route struct {
	Mode     RouteMode
	Async    bool
	Topic    string
	Endpoint string
	URI      string
}
```

```Receiver``` is a kind of endpoint that receives request messages and routes it to a specific ```Sender``` passed to its ```Listen*()``` methods. ```Listen()``` assumes receiving and routing requests without responding whereas ```ListenRequest()``` responds to external senders.

```Sender``` sends request messages by executing the ```Send()``` method and does the same and then responds by executing the ```Request()``` one.

### Injection
To bring the new drivers to life you need to modify the ```(*NATter) setupConns() error``` method located in ```cmd/natter.go``` by adding a new pair of key and value to a ```NATter```'s' ```conns map[string]driver.Conn``` map field. The key is the string name of the driver that is the part of the route mode and the value is the corresponding driver connection.

## License
NATter is released under the MIT license. See [LICENSE](LICENSE).
