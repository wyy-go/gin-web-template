package logger

import (
	"context"
	"io"
)

type Option func(*Options)

type Options struct {
	// The logging level the logger should log at. default is `InfoLevel`
	Level Level
	// fields to always be logged
	Fields map[string]interface{}
	// It's common to set this to a file, or leave it default which is `os.Stderr`
	Writer []io.Writer
	// Caller skip frame count for file:line info
	CallerSkipCount int
	// Alternative options
	Context context.Context
}

// WithFields set default fields for the logger
func WithFields(fields map[string]interface{}) Option {
	return func(opts *Options) {
		opts.Fields = fields
	}
}

// WithLevel set default level for the logger
func WithLevel(level Level) Option {
	return func(opts *Options) {
		opts.Level = level
	}
}

// WithWriter set default output writer for the logger
func WithWriter(w []io.Writer) Option {
	return func(opts *Options) {
		opts.Writer = w
	}
}

// WithCallerSkipCount set frame count to skip
func WithCallerSkipCount(c int) Option {
	return func(opts *Options) {
		opts.CallerSkipCount = c
	}
}

func SetOption(k, v interface{}) Option {
	return func(opts *Options) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, k, v)
	}
}
