package nats

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	nc *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsClient{nc: nc}, nil
}

func (c *NatsClient) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.nc.Publish(subject, payload)
}

func (c *NatsClient) Subscribe(subject string, callback func([]byte)) (*nats.Subscription, error) {
	return c.nc.Subscribe(subject, func(msg *nats.Msg) {
		callback(msg.Data)
	})
}

func (c *NatsClient) Close() {
	c.nc.Close()
}

func (c *NatsClient) GetConnection() *nats.Conn {
	return c.nc
}
