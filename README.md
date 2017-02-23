# Pushover Broker

## Goal

The goal of this project is to provide a reliable service of delivering the push notifications via Pushover (see http://pushover.net).

The Pushover Broker is designed to run locally on your device and offer your client code application the same interface as the original Pushover API, but in addition to that make sure that 100% of the push notification messages are always delivered.

## Usage

### Installation

Get & install the package using the:
    go get github.com/martinjansa/pushoverbroker

Provide the server certificates into $GOPATH/src/github.com/martinjansa/pushoverbroker/private/server.cert.pem & server.key.pem. Optionally you can use the github.com/martinjansa/pushoverbroker/utils/generateservercert.sh to generate the self-signed certificates (unsecure for production).

### Pushing messages

The Pushover Broker provides the same API as the original Pushover API at https://localhost:8499/1/messages.json. See the Pushover API documentation at https://pushover.net/api to study the usage and parameters.

The message request and response parameters are transparently forwarded to the Pushover API with the exceptions listed bellow:
 - timestamp: if specified by the client it is transparently passed to the Pushover API, if not specified and the message sending needs to be retried the timestamp of the original acceptance is passed to the Pushover API instead of empty parameter.
 - the response status code is 202 (Accepted) in case the delivery of the message to the Pushover API fails due to temporary reasons (no internet, internal server error, timeouts, etc.)
 - the receipient request for the priority messages might be locally generated and therefore not compatible and recognized with the original Pushover API (do not mix!)
 - the values in the pushover message limits might not represent the up to date information if the broker is offline and interprets the values based on the last successful response and the queue leght

Note: The following functions have not been implemented yet:
 - the returning of the response in the JSON format on /1/messages.json
 - acceptance of the messages on /1/messages.xml interface
 - all other APIs (getting of the delivery status, cancelling the priority message, etc.)

### Cancelling and getting status of the priority messages

Note: not implemented yet.

## Techology

The service is implemented in Go language, using the RESTful API via the HTTPS server. Internally the service is structured into following components:
 - server.go             - the RESTful API server that handles the clients requests and responses
 - processor.go          - message processor, internal logic of delivering messages to the external Pushover API, keeping the messages queue, providing the status information, etc.
 - pushoverconnector.go  - connector to the Pushover API, responsible for communication to the external system
 - messagerepository.go  - responsible for the persistence of the messages queue and mapping of the priority messages recipients tokens, limits, etc.

## Method

The service is implemented using the Test Driven Development for the unit test and Behavior Driven Development for the acceptance tests

## Comment

This is my first project in Go, my first RESTful service and first project implemented from the very beginning using the TDD and BDD approach, so I appreciate any feedback or hints on the code or the tests.

## Resources
- http://pushover.net
