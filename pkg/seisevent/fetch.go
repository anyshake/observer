package seisevent

import (
	"context"
	"fmt"

	"golang.org/x/sync/singleflight"
)

func fetchEvents(ctx context.Context, group *singleflight.Group, key string, fetch func(context.Context) ([]Event, error)) ([]Event, error) {
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		result := group.DoChan(key, func() (val any, err error) {
			defer func() {
				if r := recover(); r != nil {
					val, err = nil, fmt.Errorf("seismic event source %s panicked: %v", key, r)
				}
			}()
			events, err := fetch(ctx)
			if err != nil {
				return ctx.Err(), err
			}
			return events, nil
		})
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case res := <-result:
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if res.Err != nil {
				if leaderErr, _ := res.Val.(error); leaderErr != nil {
					continue
				}
				return nil, res.Err
			}
			events, _ := res.Val.([]Event)
			return events, nil
		}
	}
}
