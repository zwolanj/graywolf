package app

import (
	"github.com/chrissnell/graywolf/pkg/configstore"
	"github.com/chrissnell/graywolf/pkg/platformproto"
)

// inboundNotification builds the Android notification payload for a stored
// message, or reports false when the row should not raise one.
func inboundNotification(m *configstore.Message) (*platformproto.MessageNotification, bool) {
	if m == nil || m.Direction != "in" || !m.Unread || m.Text == "" {
		return nil, false
	}
	if m.IsAck || m.IsRej || m.IsBulletin || m.IsNWS {
		return nil, false
	}
	return &platformproto.MessageNotification{
		ThreadKind: m.ThreadKind,
		ThreadKey:  m.ThreadKey,
		FromCall:   m.FromCall,
		Text:       m.Text,
		MessageId:  m.ID,
	}, true
}
