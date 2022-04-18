package roster

import (
	"encoding/xml"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"indiebackend.com/messaging-test/client"
	"indiebackend.com/messaging-test/utils"
)

func TestRoster_Set_Empty(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	manager.Connect()

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{})
	assert.Equal(t, nil, iq.Payload)
	assert.NotEqual(t, nil, iq.Error)
	assert.Equal(t, "bad-request", iq.Error.Reason)
	assert.Equal(t, iq.Error.Type, stanza.ErrorTypeModify)
}

func TestRoster_Set_One(t *testing.T) {
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
			{Jid: "test@indiebackend.com", Name: "test-name", Groups: []string{"test-group1", "test"}},
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
		Name:         "test-name",
		Subscription: "",
		Groups:       []string{"test-group1", "test"},
	}})
}

func TestRoster_Set_IgnoreSubscription(t *testing.T) {
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
			{Jid: "test@indiebackend.com", Name: "test-name", Groups: []string{"test-group1", "test"}, Subscription: "subscriptionToIgnore"},
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
		Name:         "test-name",
		Subscription: "",
		Groups:       []string{"test-group1", "test"},
	}})
}

func TestRoster_Set_Multiple(t *testing.T) {
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
			{Jid: "test@indiebackend.com"},
			{Jid: "test2@indiebackend.com"},
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
		Subscription: "",
		Groups:       []string(nil),
	}})
}

func TestRoster_Set_Self(t *testing.T) {
	manager := client.NewClientManager()
	romeo, selfJid := manager.Register(func(s xmpp.Sender, p stanza.Packet) {
		log.Println("Got in self", p)
	})

	manager.Connect()

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: selfJid},
		},
	})
	assert.Equal(t, nil, iq.Payload)
	assert.NotEqual(t, nil, iq.Error)
	assert.Equal(t, "bad-request", iq.Error.Reason)
	assert.Equal(t, iq.Error.Type, stanza.ErrorTypeModify)
}

func TestRoster_Set_Update(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(2)

	manager := client.NewClientManager()

	var packets []stanza.Packet

	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {
		packets = append(packets, p)
		wg.Done()
	})

	manager.Connect()

	targetJid := "test@indiebackend.com"

	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: targetJid, Name: "Before"},
		},
	})
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: targetJid, Name: "After"},
		},
	})

	wg.Wait()

	utils.AssertRosterPush(t, packets[0], []stanza.RosterItem{{
		XMLName:      xml.Name{Space: "jabber:iq:roster", Local: "item"},
		Jid:          targetJid,
		Ask:          "",
		Subscription: "",
		Name:         "Before",
		Groups:       []string(nil),
	}})
	utils.AssertRosterPush(t, packets[1], []stanza.RosterItem{{
		XMLName:      xml.Name{Space: "jabber:iq:roster", Local: "item"},
		Jid:          targetJid,
		Ask:          "",
		Subscription: "",
		Name:         "After",
		Groups:       []string(nil),
	}})
}
