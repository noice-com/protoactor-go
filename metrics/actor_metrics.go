// Copyright (C) 2017 - 2022 Asynkron.se <http://www.asynkron.se>

package metrics

import (
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const LibName string = "protoactor"

type ActorMetrics struct {
	// Mutual Exclusion Primitive to use with ActorMailboxLength
	mu *sync.Mutex

	// MetricsID
	ID string

	// Actors
	ActorFailureCount            metric.Int64Counter
	ActorMailboxLength           metric.Int64ObservableGauge
	ActorMessageReceiveHistogram metric.Float64Histogram
	ActorRestartedCount          metric.Int64Counter
	ActorSpawnCount              metric.Int64Counter
	ActorStoppedCount            metric.Int64Counter

	// Deadletters
	DeadLetterCount metric.Int64Counter

	// Futures
	FuturesStartedCount   metric.Int64Counter
	FuturesCompletedCount metric.Int64Counter
	FuturesTimedOutCount  metric.Int64Counter

	// Threadpool
	ThreadPoolLatency metric.Int64Histogram
}

// NewActorMetrics creates a new ActorMetrics value and returns a pointer to it
func NewActorMetrics() *ActorMetrics {
	instruments := newInstruments()
	return instruments
}

// newInstruments will create instruments using a meter from
// the given provider p
func newInstruments() *ActorMetrics {
	meter := otel.Meter(LibName)
	instruments := ActorMetrics{mu: &sync.Mutex{}}

	var err error

	if instruments.ActorFailureCount, err = meter.Int64Counter(
		"protoactor_actor_failure_count",
		metric.WithDescription("Number of actor failures"),
		metric.WithUnit("1"),
	); err != nil {
		err = fmt.Errorf("failed to create ActorFailureCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ActorMessageReceiveHistogram, err = meter.Float64Histogram(
		"protoactor_actor_message_receive_duration_seconds",
		metric.WithDescription("Actor's messages received duration in seconds"),
	); err != nil {
		err = fmt.Errorf("failed to create ActorMessageReceiveHistogram instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ActorRestartedCount, err = meter.Int64Counter(
		"protoactor_actor_restarted_count",
		metric.WithDescription("Number of actors restarts"),
		metric.WithUnit("1"),
	); err != nil {
		err = fmt.Errorf("failed to create ActorRestartedCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ActorStoppedCount, err = meter.Int64Counter(
		"protoactor_actor_stopped_count",
		metric.WithDescription("Number of actors stopped"),
		metric.WithUnit("1"),
	); err != nil {
		err = fmt.Errorf("failed to create ActorStoppedCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ActorSpawnCount, err = meter.Int64Counter(
		"protoactor_actor_spawn_count",
		metric.WithDescription("Number of actors spawn"),
		metric.WithUnit("1"),
	); err != nil {
		err = fmt.Errorf("failed to create ActorSpawnCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.DeadLetterCount, err = meter.Int64Counter(
		"protoactor_deadletter_count",
		metric.WithDescription("Number of deadletters"),
		metric.WithUnit("1"),
	); err != nil {
		err = fmt.Errorf("failed to create DeadLetterCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.FuturesCompletedCount, err = meter.Int64Counter(
		"protoactor_futures_completed_count",
		metric.WithDescription("Number of futures completed"),
		metric.WithUnit("1"),
	); err != nil {
		err = fmt.Errorf("failed to create FuturesCompletedCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.FuturesStartedCount, err = meter.Int64Counter(
		"protoactor_futures_started_count",
		metric.WithDescription("Number of futures started"),
		metric.WithUnit("1"),
	); err != nil {
		err = fmt.Errorf("failed to create FuturesStartedCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.FuturesTimedOutCount, err = meter.Int64Counter(
		"protoactor_futures_timed_out_count",
		metric.WithDescription("Number of futures timed out"),
		metric.WithUnit("1"),
	); err != nil {
		err = fmt.Errorf("failed to create FuturesTimedOutCount instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	if instruments.ThreadPoolLatency, err = meter.Int64Histogram(
		"protoactor_thread_pool_latency_duration_seconds",
		metric.WithDescription("History of latency in second"),
		metric.WithUnit("ms"),
	); err != nil {
		err = fmt.Errorf("failed to create ThreadPoolLatency instrument, %w", err)
		plog.Error(err.Error(), log.Error(err))
	}

	return &instruments
}

// SetActorMailboxLengthGauge makes sure access to ActorMailboxLength is sequenced
func (am *ActorMetrics) SetActorMailboxLengthGauge(gauge metric.Int64ObservableGauge) {
	// lock our mutex
	am.mu.Lock()
	defer am.mu.Unlock()

	am.ActorMailboxLength = gauge
}
