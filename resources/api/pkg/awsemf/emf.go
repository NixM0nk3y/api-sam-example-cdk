package awsemf

import (
	"api/pkg/log"
	"context"

	"github.com/prozz/aws-embedded-metrics-golang/emf"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey string

func (c contextKey) String() string {
	return "context key " + string(c)
}

//
var (
	contextKeyConfig = contextKey("config")
)

// ReadEnvConfig is
func NewEMF(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyConfig, emf.New())
}

// GetEMF is
func GetEMF(ctx context.Context) *emf.Logger {

	logger := log.Logger(ctx)

	mlogger, ok := ctx.Value(contextKeyConfig).(*emf.Logger)

	if !ok {
		logger.Panic("unable to retrieve metrics logger")
	}

	return mlogger
}
