package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/tink3rlabs/magic/observability"
)

func setupObservability(ctx context.Context, serviceName string) (*observability.Observer, error) {
	cfg := observability.DefaultConfig()
	cfg.ServiceName = serviceName
	cfg.MetricsMode = observability.MetricsModePrometheus
	cfg.EnableTracing = viper.GetBool("observability.tracing.enabled")
	cfg.TracesOTLPEndpoint = viper.GetString("observability.tracing.otlp_endpoint")
	cfg.TracesOTLPInsecure = viper.GetBool("observability.tracing.otlp_insecure")

	obs, err := observability.Init(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("observability.Init(%s): %w", serviceName, err)
	}
	return obs, nil
}
