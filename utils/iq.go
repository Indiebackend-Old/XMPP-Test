package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gosrc.io/xmpp/stanza"
)

func AssertNoErrorsIq(t *testing.T, iq stanza.IQ) {
	// Directly testing equality for (nil, iq.Error) doesn't work for some reasons
	assert.Equal(t, true, iq.Error == nil)
}

func AssertRosterPush(t *testing.T, packet stanza.Packet, expected []stanza.RosterItem) {
	push, okPush := packet.(*stanza.IQ)
	assert.Equal(t, true, okPush)
	AssertNoErrorsIq(t, *push)

	payload, okPayload := push.Payload.(*stanza.RosterItems)
	assert.Equal(t, true, okPayload)
	assert.Equal(t, expected, payload.Items)
}
