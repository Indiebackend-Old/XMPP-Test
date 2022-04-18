package roster

import (
	"encoding/xml"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"indiebackend.com/messaging-test/client"
	"indiebackend.com/messaging-test/utils"
)

func TestRoster_Get_Empty(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	manager.Connect()

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeGet}, &stanza.RosterItems{})
	roster, ok := iq.Payload.(*stanza.RosterItems)
	utils.AssertNoErrorsIq(t, iq)
	assert.Equal(t, true, ok)
	assert.Equal(t, 0, len(roster.Items))
}

func TestRoster_Get_1User(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})
	_, jid := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	manager.Connect()
	item := stanza.RosterItem{
		XMLName: xml.Name{Space: "jabber:iq:roster", Local: "item"},
		Jid:     jid,
		Name:    "test",
	}
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			item,
		},
	})

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeGet}, &stanza.RosterItems{})
	roster, ok := iq.Payload.(*stanza.RosterItems)
	utils.AssertNoErrorsIq(t, iq)
	assert.Equal(t, true, ok, "invalid roster type")
	assert.Equal(t, 1, len(roster.Items), "invalid roster length")
	assert.Equal(t, item, roster.Items[0])
}

func TestRoster_Get_MultipleUser(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	manager.Connect()

	var items []stanza.RosterItem
	itemsLength := rand.Intn(5)
	for i := 0; i < itemsLength; i++ {
		name := utils.RandomID()
		items = append(items, stanza.RosterItem{
			XMLName: xml.Name{Space: "jabber:iq:roster", Local: "item"},
			Jid:     name + "@indiebackend.com",
			Name:    name,
		})
		romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
			Items: []stanza.RosterItem{
				items[i],
			},
		})
	}

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeGet}, &stanza.RosterItems{})
	roster, ok := iq.Payload.(*stanza.RosterItems)
	utils.AssertNoErrorsIq(t, iq)
	assert.Equal(t, true, ok, "invalid roster type")
	assert.Equal(t, itemsLength, len(roster.Items), "invalid roster length")
	assert.ElementsMatch(t, items, roster.Items)
}

func TestRoster_Get_1Group(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	manager.Connect()
	item := stanza.RosterItem{
		XMLName: xml.Name{Space: "jabber:iq:roster", Local: "item"},
		Jid:     "item1@indiebackend.com",
		Name:    "test",
		Groups:  []string{"test-group"},
	}
	item.Jid = "item2@indiebackend.com"
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			item,
		},
	})

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeGet}, &stanza.RosterItems{})
	roster, ok := iq.Payload.(*stanza.RosterItems)
	utils.AssertNoErrorsIq(t, iq)
	assert.Equal(t, true, ok, "invalid roster type")
	assert.Equal(t, 1, len(roster.Items), "invalid roster length")
	assert.Equal(t, "test-group", roster.Items[0].Groups[0])
}

func TestRoster_Get_MultipleGroups(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	var groups []string
	groupsLength := rand.Intn(10)
	for i := 0; i < groupsLength; i++ {
		groups = append(groups, utils.RandomID())
	}

	manager.Connect()
	item := stanza.RosterItem{
		XMLName: xml.Name{Space: "jabber:iq:roster", Local: "item"},
		Jid:     "item1@indiebackend.com",
		Name:    "test",
		Groups:  groups,
	}
	item.Jid = "item2@indiebackend.com"
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			item,
		},
	})

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeGet}, &stanza.RosterItems{})
	roster, ok := iq.Payload.(*stanza.RosterItems)
	utils.AssertNoErrorsIq(t, iq)
	assert.Equal(t, true, ok, "invalid roster type")
	assert.Equal(t, 1, len(roster.Items), "invalid roster length")
	assert.ElementsMatch(t, groups, roster.Items[0].Groups)
}
