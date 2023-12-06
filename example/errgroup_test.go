package util

import (
	"context"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestErrGroup(*testing.T) {
	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(3)
	g.Go(func() error { return nil })
	g.Wait()
}
