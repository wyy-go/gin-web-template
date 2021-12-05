package zap

import (
	"github.com/wyy-go/go-web-template/pkg/logger"
)

type Options struct {
	logger.Options
}

type callerSkipKey struct{}

func WithCallerSkip(i int) logger.Option {
	return logger.SetOption(callerSkipKey{}, i)
}

type namespaceKey struct{}

func WithNamespace(namespace string) logger.Option {
	return logger.SetOption(namespaceKey{}, namespace)
}

type productionKey struct{}

func WithProduction(b bool) logger.Option {
	return logger.SetOption(productionKey{}, b)
}
