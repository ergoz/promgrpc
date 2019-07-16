// Package promgrpc is an instrumentation package that allows capturing metrics of your gRPC based services, both the server and the client side.
// The main goal of version 4 was to make it modular without sacrificing the simplicity of use.
//
// It is still possible to integrate the package in just a few lines.
// However, if necessary, metrics can be added, removed or modified freely.
//
// Design
//
// The package does not introduce any new concepts to an already complicated environment.
// Instead, it focuses on providing implementations of interfaces exported by gRPC and Prometheus libraries.
//
// It causes no side effects nor has global state.
// Instead, it comes with handy one-liners to reduce integration overhead.
//
// The package achieved high modularity by using Inversion of Control.
// We can define three layers of abstraction, where each is configurable or if necessary replaceable.
//
// Collectors serve one purpose, storing metrics.
// These are types well known from Prometheus ecosystem, like counters, gauges, histograms or summaries.
// This package comes with a set of predefined functions that create a specific instances for each use case. For example:
//
//  func NewRequestsTotalCounterVec(Subsystem, ...CollectorOption) *prometheus.CounterVec
//
// Level higher consist of stats handlers. This layer is responsible for metrics collection.
// It is aware of a collector and knows how to use it to record event occurrences.
// Each implementation satisfies stats.Handler and prometheus.Collector interface and knows how to monitor a single dimension, e.g. a total number of received/sent requests:
//
//  func NewRequestsStatsHandler(Subsystem, *prometheus.GaugeVec, ...StatsHandlerOption) *RequestsStatsHandler
//
// Above all, there is a coordinator.
// StatsHandler combines multiple stats handlers into a single instance.
//
// Metrics
//
// The package comes with eighteen predefined metrics — nine for server and nine for client side:
//
//  grpc_client_connections
//  grpc_client_message_received_size_histogram_bytes
//  grpc_client_message_sent_size_histogram_bytes
//  grpc_client_messages_received_total
//  grpc_client_messages_sent_total
//  grpc_client_request_duration_histogram_seconds
//  grpc_client_requests_in_flight
//  grpc_client_requests_sent_total
//  grpc_client_responses_received_total
//  grpc_server_connections
//  grpc_server_message_received_size_histogram_bytes
//  grpc_server_message_sent_size_histogram_bytes
//  grpc_server_messages_received_total
//  grpc_server_messages_sent_total
//  grpc_server_request_duration_histogram_seconds
//  grpc_server_requests_in_flight
//  grpc_server_requests_received_total
//  grpc_server_responses_sent_total
//
// Configuration
//
// The package does not require any configuration whatsoever but makes it possible.
// It is beneficial for different reasons.
//
// Having all metrics enabled could not be desirable.
// Some, like histograms, can create significant overhead on the producer side.
// If performance is critical, it advisable to reduce the set of metrics.
// To do that, implement a custom version of coordinator constructor, ClientStatsHandler and ServerStatsHandler.
//
// Another good reason to change default settings is backward compatibility.
// Migration of Grafana dashboards is not an easy nor quick task.
// If the discrepancy is small and, e.g. the only necessary adjustment is changing the namespace, it is achievable by passing CollectorWithNamespace to a collector constructor.
// It is the same very known pattern from the gRPC package, with some enhancements.
// What makes it different is that both StatsHandlerOption and CollectorOption have a shareable variant, called ShareableCollectorOption and ShareableStatsHandlerOption respectively.
// Thanks to that, it is possible to pass options related to stats handlers and collectors to coordinator constructors.
// Constructors take care of moving options to the correct receivers.
//
// Mixing both strategies described above will give even greater freedom.
// However, if that is even not enough, it is possible to reimplement an entire stack for a given metric or metrics.
package promgrpc