package promgrpc

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/stats"
)

// NewServerResponsesTotalCounterVec allocates a new Prometheus CounterVec for the server and given set of options.
func NewServerResponsesTotalCounterVec(opts ...CollectorOption) *prometheus.CounterVec {
	return newResponsesTotalCounterVec("server", "responses_sent_total", "TODO", opts...)
}

// ServerResponsesTotalStatsHandler is responsible for counting number of incoming (server side) or outgoing (client side) requests.
type ServerResponsesTotalStatsHandler struct {
	baseStatsHandler
	vec *prometheus.CounterVec
}

// NewServerResponsesTotalStatsHandler ...
func NewServerResponsesTotalStatsHandler(vec *prometheus.CounterVec, opts ...StatsHandlerOption) *ServerResponsesTotalStatsHandler {
	h := &ServerResponsesTotalStatsHandler{
		baseStatsHandler: baseStatsHandler{
			collector: vec,
			options: statsHandlerOptions{
				handleRPCLabelFn: responsesTotalLabels,
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
func (h *ServerResponsesTotalStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	if _, ok := stat.(*stats.End); ok {
		switch {
		case !stat.IsClient():
			h.vec.WithLabelValues(h.options.handleRPCLabelFn(ctx, stat)...).Inc()
		}
	}
}

func responsesTotalLabels(ctx context.Context, stat stats.RPCStats) []string {
	tag := ctx.Value(tagRPCKey).(rpcTag)
	return []string{
		tag.clientUserAgent,
		status.Code(stat.(*stats.End).Error).String(),
		tag.isFailFast,
		tag.method,
		tag.service,
	}
}
