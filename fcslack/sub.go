package fcslack

import (
	"fmt"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

type subscription struct {
	Channel string
	Module  string
	Path    string
	Closer  node.NotifyCloser
	Counter uint32
}

// subscriptionKey is compound key of module and path
func subscriptionKey(module string, path string) string {
	return fmt.Sprintf("%s:%s", module, path)
}

func (n *subscription) id() string {
	return subscriptionKey(n.Module, n.Path)
}

// subscribe subscribes to the event stream and delegates message to sink
func (n *subscription) subscribe(b *node.Browser, sink Sink, errs ErrListener) error {
	s := b.Root().Find(n.Path)
	if s.LastErr != nil {
		return s.LastErr
	}
	var err error
	n.Closer, err = s.Notifications(n.stream(sink, errs))
	return err
}

func (s *subscription) stream(sink Sink, errs ErrListener) node.NotifyStream {
	return func(notify node.Notification) {
		txt, err := nodeutil.WriteJSON(notify.Event)
		if err != nil {
			errs(err, s)
			return
		}
		s.Counter++
		err = sink(msg{
			Channel: s.Channel,
			Text:    txt,
		})
		if err != nil {
			errs(err, s)
			return
		}
	}
}

func (s *subscription) Close() {
	if s.Closer != nil {
		// don't think we care if unsubscribe doesn't work
		s.Closer()
		s.Closer = nil
	}
}
