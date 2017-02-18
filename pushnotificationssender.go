package main

// PushNotificationsSender represents the connector to the Pushover API
type PushNotificationsSender interface {
	PostPushNotificationMessage(message PushNotification) error
}
