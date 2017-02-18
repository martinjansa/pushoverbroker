package main

// PushoverConnector represents the connector to the Pushover API
type PushoverConnector interface {
	PostPushNotificationMessage(message PushNotification) error
}
