package app

import (
	"testing"

	"github.com/chrissnell/graywolf/pkg/configstore"
)

func TestInboundNotification(t *testing.T) {
	base := func() *configstore.Message {
		return &configstore.Message{
			ID: 7, Direction: "in", Unread: true, Text: "hi",
			FromCall: "W1ABC-9", ThreadKind: "dm", ThreadKey: "W1ABC-9",
		}
	}

	n, ok := inboundNotification(base())
	if !ok {
		t.Fatal("plain inbound DM should notify")
	}
	if n.ThreadKind != "dm" || n.ThreadKey != "W1ABC-9" || n.FromCall != "W1ABC-9" || n.Text != "hi" || n.MessageId != 7 {
		t.Errorf("unexpected payload: %+v", n)
	}

	skips := map[string]func(*configstore.Message){
		"outbound": func(m *configstore.Message) { m.Direction = "out" },
		"read":     func(m *configstore.Message) { m.Unread = false },
		"empty":    func(m *configstore.Message) { m.Text = "" },
		"ack":      func(m *configstore.Message) { m.IsAck = true },
		"rej":      func(m *configstore.Message) { m.IsRej = true },
		"bulletin": func(m *configstore.Message) { m.IsBulletin = true },
		"nws":      func(m *configstore.Message) { m.IsNWS = true },
	}
	for name, mutate := range skips {
		m := base()
		mutate(m)
		if _, ok := inboundNotification(m); ok {
			t.Errorf("%s row must not notify", name)
		}
	}
	if _, ok := inboundNotification(nil); ok {
		t.Error("nil row must not notify")
	}
}
