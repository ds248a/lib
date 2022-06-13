package common

import (
	"context"
	"time"
)

// Sleep приостанавливает процесс на указанный интервал времени.
func Sleep(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
