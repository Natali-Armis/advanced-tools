package slack_client

import (
	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
)

type SlackClient struct {
	token  string
	client *slack.Client
}

func GetSlackClient(authToken string) *SlackClient {
	return &SlackClient{
		token:  authToken,
		client: slack.New(authToken),
	}
}

func (s *SlackClient) SendMessage(channelID string, message string) error {
	_, _, err := s.client.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		log.Error().Msgf("failed to send message to channel %v: %v", channelID, err)
		return err
	}
	log.Info().Msgf("message sent to channel %v", channelID)
	return nil
}
