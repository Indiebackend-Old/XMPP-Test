package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
)

type OnIqFunc func(xmpp.Sender, stanza.Packet)

type TestXMPPClient struct {
	c    *xmpp.Client
	recv chan stanza.Packet
}

func (t *TestXMPPClient) SendIq(iqAttrs stanza.Attrs, payload stanza.IQPayload) stanza.IQ {
	iq, err := stanza.NewIQ(iqAttrs)
	if err != nil {
		log.Fatal(err)
	}
	iq.Payload = payload
	c, _ := t.c.SendIQ(context.TODO(), iq)
	resp := <-c
	return resp
}

func (t *TestXMPPClient) Receive() stanza.Packet {
	packet := <-t.recv
	return packet
}

func (t *TestXMPPClient) Client() *xmpp.Client {
	return t.c
}

func CreateClient(jid string, withLogger bool, fn OnIqFunc) *TestXMPPClient {
	testClient := &TestXMPPClient{
		recv: make(chan stanza.Packet),
	}

	logger := os.Stdout
	if !withLogger {
		logger = nil
	}

	config := &xmpp.Config{
		Jid:        jid,
		Credential: xmpp.Password("Test"),
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: "localhost:3010",
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		StreamLogger: logger,
	}

	router := xmpp.NewRouter()
	router.HandleFunc("iq", fn)

	client, err := xmpp.NewClient(config, router, func(err error) { fmt.Println(err.Error()) })
	if err != nil {
		log.Fatalf("%+v", err)
	}
	testClient.c = client

	return testClient
}
