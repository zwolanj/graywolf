//go:build android

package app

import (
	"context"

	"github.com/chrissnell/graywolf/pkg/messages"
)

// startMessageNotifier forwards each inbound message to the Kotlin side as a
// system notification. It runs in Go because the WebView is not alive while
// the app is backgrounded. Shutdown is awaited via messagesReloadWG.
func (a *App) startMessageNotifier(ctx context.Context) {
	if a.platformClient == nil || a.msgSvc == nil {
		return
	}
	events, unsub := a.msgSvc.EventHub().Subscribe()
	a.messagesReloadWG.Add(1)
	go func() {
		defer a.messagesReloadWG.Done()
		defer unsub()
		for {
			select {
			case <-ctx.Done():
				return
			case ev, ok := <-events:
				if !ok {
					return
				}
				if ev.Type != messages.EventMessageReceived {
					continue
				}
				a.notifyInbound(ctx, ev.MessageID)
			}
		}
	}()
}

func (a *App) notifyInbound(ctx context.Context, id uint64) {
	prefs, err := a.store.GetMessagePreferences(ctx)
	if err == nil && prefs != nil && prefs.NotificationsMuted {
		return
	}
	row, err := a.msgStore.GetByID(ctx, id)
	if err != nil {
		a.logger.Debug("message notify: lookup failed", "id", id, "err", err)
		return
	}
	n, ok := inboundNotification(row)
	if !ok {
		return
	}
	if err := a.platformClient.NotifyMessage(n); err != nil {
		a.logger.Debug("message notify: send failed", "id", id, "err", err)
	}
}
