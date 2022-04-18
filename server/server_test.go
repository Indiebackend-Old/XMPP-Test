package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"indiebackend.com/messaging-test/client"
)

func TestVersion(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {
	})

	manager.Connect()

	res := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeGet}, &stanza.Version{})

	version, ok := res.Payload.(*stanza.Version)

	assert.Equal(t, true, ok)
	assert.Equal(t, "indiebackend-messaging", version.Name)
}
