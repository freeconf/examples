package fcslack

import (
	vendor "github.com/slack-go/slack"
)

type slackOptions struct {
	ApiToken  string
	UserToken string
	Debug     bool
	Emulate   bool
}

type msg struct {
	Channel string
	Text    string
}

type Sink func(msg) error

type slack struct {
	opts slackOptions
	api  *vendor.Client
}

func newSlack() *slack {
	return &slack{}
}

func (s *slack) options() slackOptions {
	return s.opts
}

func (s *slack) apply(opts slackOptions) error {
	s.api = vendor.New(opts.ApiToken)
	s.opts = opts
	return nil
}

func (s *slack) Send(m msg) error {
	_, _, err := s.api.PostMessage(m.Channel, vendor.MsgOptionText(m.Text, true))
	return err
}
