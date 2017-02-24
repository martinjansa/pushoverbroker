package main

// PushNotificationsSender represents the connector to the Pushover API
type PushNotificationsSender interface {
	// PostPushNotificationMessage handles the message and returns error if ocurred (or nil) and response code (or 0 on POST error)
	PostPushNotificationMessage(response *PushNotificationHandlingResponse, message PushNotification) error
}
