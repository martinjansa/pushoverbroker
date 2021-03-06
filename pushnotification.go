package main

import (
	"errors"
	"fmt"
)

// PushNotification represents a message with json request that is passed to the REST API
type PushNotification struct {
	Token   string `json:"token" schema:"token"`
	User    string `json:"user"  schema:"user"`
	Message string `json:"message" schema:"message"`
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

// check the validity of the PushNotification message
func (m *PushNotification) Validate() error {
	if m.Token == "" {
		return errors.New("push notification token value cannot be empty")
	}
	if m.User == "" {
		return errors.New("push notification user value cannot be empty")
	}
	if m.Message == "" {
		return errors.New("push notification message value cannot be empty")
	}
	return nil
}

// DumpToString converts the PushNotification to string
func (m *PushNotification) DumpToString() string {
	return fmt.Sprintf("token=\"%s\", user=\"%s\", message=\"%s\"", m.GetToken(), m.GetUser(), m.GetMessage())
}
