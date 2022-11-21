package app

import (
	"github.com/bytebot-chat/gateway-discord/model"
	"github.com/go-redis/redis/v8"
)

// Message is the struct that represents a message received from Discord.
// It's effectively a just wrapper around the Bytebot/Discord message struct so I can add methods to it.
type Message struct {
	*model.Message `json:"message"`
}

// handleIncomingMessage handles an incoming message from Bytebot/Discord.
func unmarshalIncomingMessage(m *redis.Message) (*Message, error) {
	// Create a new MessageSend struct
	var message model.Message

	// Unmarshal the message bytes into the struct
	err := message.UnmarshalJSON([]byte(m.Payload))
	if err != nil {
		return nil, err
	}

	// Return the message
	return &Message{&message}, nil

}

// handleOutgoingMessage handles an outgoing message to Bytebot/Discord.
func (a *App) handleOutgoingMessage(m *model.MessageSend) error {
	// Marshal the message into bytes
	bytes, err := m.MarshalJSON()
	if err != nil {
		return err
	}

	// Publish the message to the outbound topic
	res := a.redis.Publish(a.context, a.Config.OutboundTopic, bytes)
	if err != nil {
		return err
	}

	// Log the result
	a.logger.Debug().
		// We want to log enough context that we can correlate this message with other services in the logs
		Str("topic", a.Config.OutboundTopic).
		Str("dest_gateway", m.Metadata.Dest).
		Str("id", m.Metadata.ID.String()).
		Str("reply_to_user", m.PreviousMessage.Author.Username+"#"+m.PreviousMessage.Author.Discriminator).
		Str("reply_to_message_id", m.PreviousMessage.ID).
		Msg("published message")

	// Return the result
	return res.Err()
}
