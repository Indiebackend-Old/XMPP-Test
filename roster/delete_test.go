package roster

import (
	"encoding/xml"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"indiebackend.com/messaging-test/client"
	"indiebackend.com/messaging-test/utils"
)

func TestRoster_Delete_One(t *testing.T) {
	var wg sync.WaitGroup
	var pushPacket stanza.Packet
	wg.Add(1)

	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {
		pushPacket = p
		wg.Done()
	})

	manager.Connect()

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: "test@indiebackend.com", Name: "test-name", Groups: []string{"test-group1", "test"}, Subscription: "remove"},
		},
	})

	wg.Wait()
	_, okResult := iq.Payload.(*stanza.RosterItems)

	// Check result IQ
	utils.AssertNoErrorsIq(t, iq)
	assert.Equal(t, false, okResult)

	// Check push IQ
	assert.NotEqual(t, nil, pushPacket)
	utils.AssertRosterPush(t, pushPacket, []stanza.RosterItem{{
		XMLName:      xml.Name{Space: "jabber:iq:roster", Local: "item"},
		Jid:          "test@indiebackend.com",
		Ask:          "",
		Name:         "",
		Subscription: "remove",
		Groups:       []string(nil),
	}})
}

func TestRoster_Remove_Self(t *testing.T) {
	manager := client.NewClientManager()
	romeo, selfJid := manager.Register(func(s xmpp.Sender, p stanza.Packet) {
	})

	manager.Connect()

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: selfJid, Subscription: "remove"},
		},
	})
	assert.Equal(t, nil, iq.Payload)
	assert.NotEqual(t, nil, iq.Error)
	assert.Equal(t, "bad-request", iq.Error.Reason)
	assert.Equal(t, iq.Error.Type, stanza.ErrorTypeModify)
}
