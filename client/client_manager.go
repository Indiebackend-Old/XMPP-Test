package client

import (
	"log"

	"indiebackend.com/messaging-test/utils"
)

type ClientManager struct {
	clients []*TestXMPPClient
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

func (c *ClientManager) Register(fn OnIqFunc) (*TestXMPPClient, string) {
	return c.register(false, fn)
}
func (c *ClientManager) RegisterWithLogger(fn OnIqFunc) (*TestXMPPClient, string) {
	return c.register(true, fn)
}

func (c *ClientManager) register(withLogger bool, fn OnIqFunc) (*TestXMPPClient, string) {
	jid := utils.RandomID() + "@indiebackend.com"
	client := CreateClient(jid, withLogger, fn)

	c.clients = append(c.clients, client)

	return client, jid
}

func (c *ClientManager) Connect() {
	for _, client := range c.clients {
		err := client.c.Connect()
		if err != nil {
			log.Fatal(err)
		}
	}
}
