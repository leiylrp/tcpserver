package pool

import (
	"time"
)

type Option func(option *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

// Options contains all options which will be applied when instantiating an pool.
type Options struct {
	ExpiryDuration 		time.Duration

	PreAlloc			bool

	MaxBlockingTasks	int

	Nonblocking			bool

	PanicHandler		func(interface{})

	Logger Logger
}

func WithOptions(options Options) Option {
	return func(ops *Options) {
		*ops = options
	}
}

func WithExpiryDuration(expiryDuration time.Duration) Option {
	return func(ops *Options) {
		ops.ExpiryDuration = expiryDuration
	}
}

func WithPreAlloc(preAlloc bool) Option {
	return func(ops *Options) {
		ops.PreAlloc = preAlloc
	}
}

// WithMaxBlockingTasks sets up the maximum number of goroutines that are blocked when it reaches the capacity of pool.
func WithMaxBlockingTasks(maxBlockingTasks int) Option {
	return func(opts *Options) {
		opts.MaxBlockingTasks = maxBlockingTasks
	}
}

// WithNonblocking indicates that pool will return nil when there is no available workers.
func WithNonblocking(nonblocking bool) Option {
	return func(opts *Options) {
		opts.Nonblocking = nonblocking
	}
}

// WithPanicHandler sets up panic handler.
func WithPanicHandler(panicHandler func(interface{})) Option {
	return func(opts *Options) {
		opts.PanicHandler = panicHandler
	}
}

// WithLogger sets up a customized logger.
func WithLogger(logger Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}
