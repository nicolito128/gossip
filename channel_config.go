package gossip

type (
	ChannelOpt              func(*ChannelConfig)
	ChannelErrorHandlerFunc func(error, Transporter)
)

type ChannelConfig struct {
	Topic      string
	Tp         Transporter
	ErrHandler ChannelErrorHandlerFunc
}

func DefaultChannelConfig(opts ...ChannelOpt) *ChannelConfig {
	cc := new(ChannelConfig)
	cc.Topic = "default"

	for _, opt := range opts {
		opt(cc)
	}
	return cc
}

func WithChannelTopic(topicName string) ChannelOpt {
	return func(cc *ChannelConfig) {
		cc.Topic = topicName
	}
}

func WithChannelTransport(layer Transporter) ChannelOpt {
	return func(cc *ChannelConfig) {
		cc.Tp = layer
	}
}

func WithChannelErrorHandler(handler ChannelErrorHandlerFunc) ChannelOpt {
	return func(cc *ChannelConfig) {
		cc.ErrHandler = handler
	}
}
