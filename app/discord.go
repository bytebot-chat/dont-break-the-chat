package app

import (
	"github.com/bytebot-chat/gateway-discord/model"
	"github.com/go-redis/redis/v8"
)

// Message is the struct that represents a message received from Discord.
// It's effectively a just wrapper around the Bytebot/Discord message struct so I can add methods to it.
type Message struct {
	*model.Message
}

// handleIncomingMessage handles an incoming message from Bytebot/Discord.
func handleIncomingMessage(m *redis.Message) (*Message, error) {
	bbmsg := &model.Message{}
	err := bbmsg.UnmarshalJSON([]byte(m.Payload))
	if err != nil {
		return nil, err
	}

	return &Message{bbmsg}, nil
}

// handleOutgoingMessage handles an outgoing message to Bytebot/Discord.
func handleOutgoingMessage(m *Message) error {
	return nil
}
