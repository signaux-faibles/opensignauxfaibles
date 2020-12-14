package base

import (
	"context"
	"time"
)

// Cron exécute function régulièrement, avec l'interval fourni.
// Il retourne une fonction stop().
func Cron(interval time.Duration, function func()) (stop context.CancelFunc) {
	ctx, stop := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for range time.Tick(interval) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			function()
		}
	}(ctx)
	return stop
}
