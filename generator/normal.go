package generator

import (
	"context"
	"math/rand"

	"github.com/sirupsen/logrus"
)

// Background thread to generate gaussian values
func StartNormalGenerator(ctx context.Context, bufferSize int) (chan float64, func()) {
	logrus.Debug("Starting gaussian thread...")
	c := make(chan float64, bufferSize)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c <- rand.NormFloat64()
			}
		}
	}()
	return c, func() {
		cancel()
		defer close(c)
		for range c {
		}
	}
}
