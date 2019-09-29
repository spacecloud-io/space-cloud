package nats

import (
	"strings"

	nts "github.com/nats-io/nats.go"

	"github.com/spaceuptech/space-cloud/model"
)

// nats holds the nats driver session
type nats struct {
	conn *nts.EncodedConn
}

// Connect connects to the nats server
func Connect(connection string) (*nats, error) {
	nc, err := nts.Connect(connection)
	c, _ := nts.NewEncodedConn(nc, "json")
	if err != nil {
		return nil, err
	}
	return &nats{c}, nil
}

// Publish publishes a model.PubsubMsg to a particular subject
func (n *nats) Publish(subject string, msg *model.PubsubMsg) error {
	subject = strings.Replace(subject, "/", ".", -1)
	subject = strings.Trim(subject, ".")
	return n.conn.Publish(subject, msg)
}

// QueueSubscribe subscribes to a particular subject, using a queue
func (n *nats) QueueSubscribe(subject, queue string, ch chan *model.PubsubMsg) (model.PubsubUnsubscribe , error) {
	subject = strings.Replace(subject, "/", ".", -1)
	subject = strings.Trim(subject, ".")
	natsWildcard := ".>"
	subs1, err1 := n.conn.BindRecvQueueChan(subject + natsWildcard, queue, ch)
	if err1 != nil {
		return nil, err1
	}
	subs2, err2 := n.conn.BindRecvQueueChan(subject, queue, ch)
	if err2 != nil {
		subs1.Unsubscribe()
		return nil, err2
	}
	return func()(error) {
		err := subs1.Unsubscribe()
		if err != nil {
			return err
		}
		err = subs2.Unsubscribe()
		if err != nil {
			return err
		}
		close(ch)
		return nil
	}, nil
}
