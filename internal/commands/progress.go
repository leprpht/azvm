package commands

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

func startSpinner(ctx context.Context) func() {
	done := make(chan struct{})
	finished := make(chan struct{})

	var once sync.Once

	go func() {
		defer close(finished)

		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		frame := 0

		fmt.Fprintf(os.Stderr, "\r%s", frames[frame])

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				frame = (frame + 1) % len(frames)
				fmt.Fprintf(os.Stderr, "\r%s", frames[frame])

			case <-done:
				return

			case <-ctx.Done():
				return
			}
		}
	}()

	return func() {
		once.Do(func() {
			close(done)
			<-finished
			fmt.Fprint(os.Stderr, "\r\x1b[2K")
		})
	}
}
