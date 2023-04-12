package fcslack

import (
	"container/list"

	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/nodeutil"
)

// Service orchestrates sending notifications defined in YANG to slack
type Service struct {
	notifications map[string]*subscription
	errListeners  *list.List
	slack         *slack
	dev           device.Device
	sink          Sink
}

// NewService is constructor. Device can be local or remote device
func NewService(dev device.Device) *Service {
	s := &Service{
		dev:           dev,
		errListeners:  list.New(),
		slack:         newSlack(),
		notifications: make(map[string]*subscription),
	}
	s.sink = s.slack.Send
	return s
}

// ErrListener is for catching errors that happen in the background trying to send messages
// to slack
type ErrListener func(err error, r *subscription)

// OnError is for errors that happen in the background trying to send messages
// to slack
func (s *Service) OnError(l ErrListener) nodeutil.Subscription {
	return nodeutil.NewSubscription(s.errListeners, s.errListeners.PushBack(l))
}

func (s *Service) handleErr(err error, sub *subscription) {
	fc.Debug.Printf("%s:%s - %s", sub.Module, sub.Path, err)
	for p := s.errListeners.Front(); p != nil; p = p.Next() {
		p.Value.(ErrListener)(err, sub)
	}
}

func (s *Service) Close() {
	for _, r := range s.notifications {
		r.Close()
	}
}

func (s *Service) add(n *subscription) error {
	var err error
	b, err := s.dev.Browser(n.Module)
	if err != nil {
		return err
	}
	if err = n.subscribe(b, s.sink, s.handleErr); err != nil {
		return err
	}
	s.notifications[n.id()] = n
	return nil
}

func (s *Service) remove(key string) {
	if n, found := s.notifications[key]; found {
		n.Close()
		delete(s.notifications, key)
	}
}
