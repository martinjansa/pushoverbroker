package pushoverbroker

import (
	"fmt"
)

// PushNotification represents a message with json request that is passed to the REST API
type PushNotification struct {
	Token   string
	User    string
	Message string
}

// GetToken returns the API token from the push notification.
func (m *PushNotification) GetToken() string {
	return m.Token
}

// GetUser returns the user identification from the push notification
func (m *PushNotification) GetUser() string {
	return m.User
}

// GetMessage returns the message content from the push notification
func (m *PushNotification) GetMessage() string {
	return m.Message
}

// DumpToString converts the PushNotification to string
func (m *PushNotification) DumpToString() string {
	return fmt.Sprintf("token=\"%s\", user=\"%s\", message=\"%s\"", m.GetToken(), m.GetUser(), m.GetMessage())
}
