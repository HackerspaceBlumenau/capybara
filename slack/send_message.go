package slack

import (
	"errors"
	"log"

	"github.com/nlopes/slack"
)

var (
	rtm *slack.RTM

	ErrEmptyMessage = errors.New("empty message")
	ErrEmptyChannel = errors.New("empty channel")
)

func SendMessage(msg, channel string) error {
	if msg == "" {
		return ErrEmptyMessage
	}

	if channel == "" {
		return ErrEmptyChannel
	}

	log.Println("Sending new message...")
	rtm.SendMessage(rtm.NewOutgoingMessage(msg, channel))
	return nil
}

func Setup(slackToken string) {
	if slackToken == "" {
		panic("slack token cann't be empty")
	}

	api := slack.New(slackToken)
	rtm = api.NewRTM()
	go rtm.ManageConnection()
}
