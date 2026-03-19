package gossip

import (
	"context"
	"sync"
)

// Channel ...
type Channel struct {
	*ChannelConfig

	subscribers []Transporter
	mu          sync.RWMutex
}

func NewChannel(opts ...ChannelOpt) *Channel {
	cc := DefaultChannelConfig(opts...)
	ch := new(Channel)
	ch.ChannelConfig = cc
	ch.subscribers = make([]Transporter, 0)
	return ch
}

func (ch *Channel) AddSubscriber(tp Transporter) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	ch.subscribers = append(ch.subscribers, tp)
}

func (ch *Channel) PublishCtx(ctx context.Context, messages ...TransportMessage) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	var wg sync.WaitGroup
	for _, tp := range ch.subscribers {
		if tp == nil {
			continue
		}

		wg.Add(1)
		go func(t Transporter) {
			defer wg.Done()
			done := make(chan error, 1)

			go func() {
				defer close(done)
				for i := range messages {
					if err := t.Write(messages[i]); err != nil {
						done <- err
						return
					}
				}
			}()

			select {
			case <-ctx.Done():
				return
			case err := <-done:
				if err != nil && ch.ErrHandler != nil {
					ch.ErrHandler(err, t)
				}
			}
		}(tp)
	}
	wg.Wait()
}

func (ch *Channel) Publish(p ...TransportMessage) {
	ch.PublishCtx(context.Background(), p...)
}
