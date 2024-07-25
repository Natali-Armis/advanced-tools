package slack_client

import (
	"regexp"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
)

type SlackClient struct {
	token     string
	BotUserID string
	client    *slack.Client
	rtm       *slack.RTM
}

func GetSlackClient(authToken string) *SlackClient {
	slackClient := &SlackClient{
		token:  authToken,
		client: slack.New(authToken),
	}
	slackClient.rtm = slackClient.client.NewRTM()
	authTest, err := slackClient.client.AuthTest()
	if err != nil {
		log.Fatal().Msgf("client: failed to authenticate bot user: %v", err)
	}
	slackClient.BotUserID = authTest.UserID
	return slackClient
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

func (s *SlackClient) GetLastMessage(channelID string) (message string, userId string, err error) {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     1,
	}
	history, err := s.client.GetConversationHistory(params)
	if err != nil {
		log.Error().Msgf("failed to get conversation history from channel %v: %v", channelID, err)
		return "", "", err
	}
	if len(history.Messages) == 0 {
		log.Info().Msgf("no messages found in channel %v", channelID)
		return "", "", nil
	}
	lastMessage := history.Messages[0].Msg.Text
	user := history.Messages[0].User

	return lastMessage, user, nil
}

func (s *SlackClient) FetchHistoryUpToDate(channelID string, oldestDate time.Time) ([]slack.Message, error) {
	var allMessages []slack.Message
	oldestTimestamp := oldestDate.Unix()
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Oldest:    strconv.FormatInt(oldestTimestamp, 10),
		Limit:     1000,
	}

	for {
		history, err := s.client.GetConversationHistory(params)
		if err != nil {
			log.Error().Msgf("failed to get conversation history from channel %v: %v", channelID, err)
			return nil, err
		}
		allMessages = append(allMessages, history.Messages...)
		if !history.HasMore {
			break
		}
		params.Latest = history.Messages[len(history.Messages)-1].Timestamp
	}

	log.Info().Msgf("fetched %d messages from channel %v", len(allMessages), channelID)
	return allMessages, nil
}

func (s *SlackClient) GetLastMessageMatchPattern(channelID string, pattern string) (message string, userId string, err error) {
	message, userId, err = s.GetLastMessage(channelID)
	if err != nil {
		return "", "", err
	}
	re := regexp.MustCompile(pattern)
	if re.MatchString(message) {
		return message, userId, nil
	}
	return "", userId, nil
}

func (s *SlackClient) SendPrivateMessage(userID string, message string) error {
	channel, _, _, err := s.client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{userID},
	})
	if err != nil {
		log.Error().Msgf("failed to open IM channel with user %v: %v", userID, err)
		return err
	}

	_, _, err = s.client.PostMessage(
		channel.ID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		log.Error().Msgf("failed to send private message to user %v: %v", userID, err)
		return err
	}
	log.Info().Msgf("private message sent to user %v", userID)
	return nil
}
