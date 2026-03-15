package gossip

import (
	"sync"
)

func New() *Manager {
	return NewManager()
}

type Manager struct {
	channels map[string]*Channel

	mu sync.RWMutex
}

func NewManager() *Manager {
	m := new(Manager)
	m.channels = make(map[string]*Channel)
	return m
}

func (m *Manager) Subscribe(topic string, tp Transporter, opts ...ChannelOpt) *Channel {
	m.mu.RLock()
	ch, ok := m.channels[topic]
	m.mu.RUnlock()

	if !ok {
		if opts == nil {
			opts = make([]ChannelOpt, 0)
		}
		opts = append(opts, WithChannelTopic(topic))
		opts = append(opts, WithChannelTransport(tp))

		ch = NewChannel(opts...)
		m.mu.Lock()
		m.channels[topic] = ch
		m.mu.Unlock()
	}

	ch.AddSubscriber(tp)
	return ch
}
