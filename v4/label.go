package promgrpc

import (
	"context"

	"google.golang.org/grpc/stats"
)

type rpcTag struct {
	isFailFast      string
	service         string
	method          string
	clientUserAgent string
}

// HandleRPCLabelFunc type represents a function signature that can be passed into a stats handler and used instead of default one.
// That way caller gets the ability to modify the way labels are assembled.
type HandleRPCLabelFunc func(ctx context.Context, stat stats.RPCStats) []string
