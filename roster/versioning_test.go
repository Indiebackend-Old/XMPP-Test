package roster

import (
	"encoding/xml"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"indiebackend.com/messaging-test/client"
)

type VersionedRoster struct {
	stanza.Roster
	Ver string `xml:"ver,attr"`
}

func TestRoster_Versioning_HasFeatureAdvertised(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	manager.Connect()

	hasVersioningAdvertised := false
	for _, n := range romeo.Client().Session.Features.Any {
		if n.Local == "ver" && n.Space == "urn:xmpp:features:rosterver" {
			hasVersioningAdvertised = true
		}
	}

	assert.Equal(t, true, hasVersioningAdvertised, "roster versioning not advertised")
}

func TestRoster_Versioning_Bootstrap(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	manager.Connect()

	// Intialize the roster
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: "test@indiebackend.com", Name: "test-name", Groups: []string{"test-group1", "test"}},
		},
	})
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: "test2@indiebackend.com", Name: "test-name", Groups: []string{"test-group1", "test"}},
		},
	})

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeGet}, &VersionedRoster{
		Roster: stanza.Roster{
			XMLName: xml.Name{
				Space: stanza.NSRoster,
				Local: "query",
			},
		},
		Ver: "",
	})

	roster, ok := iq.Payload.(*stanza.RosterItems)

	assert.Equal(t, true, ok, "invalid roster format")
	assert.Len(t, roster.Items, 2)
}

func TestRoster_Versioning_SameVersionAsServer(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	manager.Connect()

	// Intialize the roster
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: "test@indiebackend.com", Name: "test-name", Groups: []string{"test-group1", "test"}},
		},
	})
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: "test2@indiebackend.com", Name: "test-name", Groups: []string{"test-group1", "test"}},
		},
	})

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeGet}, &VersionedRoster{
		Roster: stanza.Roster{
			XMLName: xml.Name{
				Space: stanza.NSRoster,
				Local: "query",
			},
		},
		Ver: "2",
	})

	log.Println(iq.Payload)
	assert.Equal(t, nil, iq.Payload, "invalid roster format")
}

func TestRoster_Versioning_DifferentVersionFromServer(t *testing.T) {
	manager := client.NewClientManager()
	romeo, _ := manager.Register(func(s xmpp.Sender, p stanza.Packet) {})

	manager.Connect()

	// Intialize the roster
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: "test@indiebackend.com", Name: "test-name", Groups: []string{"test-group1", "test"}},
		},
	})
	romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeSet}, &stanza.RosterItems{
		Items: []stanza.RosterItem{
			{Jid: "test2@indiebackend.com", Name: "test-name", Groups: []string{"test-group1", "test"}},
		},
	})

	iq := romeo.SendIq(stanza.Attrs{Type: stanza.IQTypeGet}, &VersionedRoster{
		Roster: stanza.Roster{
			XMLName: xml.Name{
				Space: stanza.NSRoster,
				Local: "query",
			},
		},
		Ver: "notTheSameRosterVersion",
	})

	roster, ok := iq.Payload.(*stanza.RosterItems)

	assert.Equal(t, true, ok, "invalid roster format")
	assert.Len(t, roster.Items, 2)
}
