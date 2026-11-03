//go:build !android

package app

import "context"

// startMessageNotifier is a no-op off Android; desktop alerts run in the SPA.
func (a *App) startMessageNotifier(ctx context.Context) {}
