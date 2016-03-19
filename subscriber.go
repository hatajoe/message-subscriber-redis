package subscriber

import "github.com/garyburd/redigo/redis"

type Subscriber struct {
	pubsubConn *redis.PubSubConn
	channel    string
	buffer     chan []byte
}

func NewSubscriber(conn redis.Conn, channel string, buffer chan []byte) (*Subscriber, error) {
	psc := redis.PubSubConn{conn}
	if err := psc.Subscribe(channel); err != nil {
		return nil, err
	}
	return &Subscriber{
		pubsubConn: &psc,
		channel:    channel,
		buffer:     buffer,
	}, nil
}

func (m *Subscriber) Subscribe() error {
	switch v := m.pubsubConn.Receive().(type) {
	case redis.Message:
		m.buffer <- v.Data
	case error:
		return v
	default:
	}
	return nil
}

func (m *Subscriber) UnSubscribe() error {
	if err := m.pubsubConn.Unsubscribe(m.channel); err != nil {
		return err
	}
	if err := m.pubsubConn.Close(); err != nil {
		return err
	}
	return nil
}

func (m *Subscriber) Abort() error {
	return nil
}

func (m *Subscriber) End() error {
	return nil
}
