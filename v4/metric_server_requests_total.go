package promgrpc

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

func NewServerRequestsTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newRequestsTotalCounterVec("server", "requests_received_total", "TODO", opts...)
}

type ServerRequestsTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewServerRequestsTotalStatsHandler ...
// The GaugeVec must have zero, one, two, three or four non-const non-curried labels.
// For those, the only allowed labelsFn names are "fail_fast", "handler", "service".
func NewServerRequestsTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerRequestsTotalStatsHandler {
	h := &ServerRequestsTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: requestsTotalLabels,
			},
		},
		vec: vec,
	}
	for _, opt := range opts {
		opt.apply(&h.options)
	}
	return h
}

// HandleRPC implements stats Handler interface.
func (h *ServerRequestsTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if beg, ok := stat.(*stats.Begin); ok {
		switch {
		case !beg.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}

func requestsTotalLabels(ctx context.Context, _ stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
