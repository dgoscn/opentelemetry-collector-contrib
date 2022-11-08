// Copyright 2022 The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package arrayscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver/internal/settings/arrayscraper"

import (
	"context"
	"strconv"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/featuregate"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver/internal/settings/arrayscraper/internal/metadata"
)

const (
	emitMetricsWithDimensionAttributeFeatureGateID    = "receiver.purefa.emitMetricsWithDimensionAttribute"
	emitMetricsWithoutDimensionAttributeFeatureGateID = "receiver.purefa.emitMetricsWithoutDimensionAttribute"
)

var (
	emitMetricsWithDimensionAttributeFeatureGate = featuregate.Gate{
		ID:      emitMetricsWithDimensionAttributeFeatureGateID,
		Enabled: true,
		Description: "Some Pure FlashArray metrics reported are transitioning from being reported with a dimension " +
			"attribute to being reported with the dimension included in the metric name to adhere to the " +
			"OpenTelemetry specification. This feature gate controls emitting the old metrics with the dimension " +
			"attribute.",
	}

	emitMetricsWithoutDimensionAttributeFeatureGate = featuregate.Gate{
		ID:      emitMetricsWithoutDimensionAttributeFeatureGateID,
		Enabled: false,
		Description: "Some Pure FlashArray metrics reported are transitioning from being reported with a dimension " +
			"attribute to being reported with the dimension included in the metric name to adhere to the " +
			"OpenTelemetry specification. This feature gate controls emitting the new metrics without the dimension " +
			"attribute.",
	}
)

func init() {
	featuregate.GetRegistry().MustRegister(emitMetricsWithDimensionAttributeFeatureGate)
	featuregate.GetRegistry().MustRegister(emitMetricsWithoutDimensionAttributeFeatureGate)
}

type purefaScraper struct {
	logger                               *zap.Logger
	config                               *Config
	mb                                   *metadata.MetricsBuilder
	newClient                            newPurefaClientFunc
	emitMetricsWithDimensionAttribute    bool
	emitMetricsWithoutDimensionAttribute bool
}

func newPurefaScraper(
	settings component.ReceiverCreateSettings,
	config *Config,
) purefaScraper {
	return purefaScraper{
		logger:                               settings.Logger,
		config:                               config,
		newClient:                            newPurefaClient,
		mb:                                   metadata.NewMetricsBuilder(config.Metrics, settings.BuildInfo),
		emitMetricsWithDimensionAttribute:    featuregate.GetRegistry().IsEnabled(emitMetricsWithDimensionAttributeFeatureGateID),
		emitMetricsWithoutDimensionAttribute: featuregate.GetRegistry().IsEnabled(emitMetricsWithoutDimensionAttributeFeatureGateID),
	}
}

func (r *purefaScraper) scrape(_ context.Context) (pmetric.Metrics, error) {
	// Init client in scrape method in case there are transient errors in the
	// constructor.
	statsClient, err := r.newClient(r.config.Endpoint, r.config.Timeout)
	if err != nil {
		r.logger.Error("Failed to establish client", zap.Error(err))
		return pmetric.Metrics{}, err
	}

	allServerStats, err := statsClient.Stats()
	if err != nil {
		r.logger.Error("Failed to fetch purefa stats", zap.Error(err))
		return pmetric.Metrics{}, err
	}

	now := pcommon.NewTimestampFromTime(time.Now())

	for _, stats := range allServerStats {
		for k, v := range stats.Stats {
			switch k {
			case "array_performance_average_bytes":
				if parsedV, ok := r.parseInt(k, v); ok {
					r.mb.RecordPurefaArrayPerformanceAverageBytesDataPoint(now, parsedV)
				}
			case "array_performance_bandwidth_bytes":
				if parsedV, ok := r.parseInt(k, v); ok {
					r.mb.RecordPurefaArrayPerformanceBandwidthByteshDataPoint(now, parsedV)
				}
			case "array_performance_latency_usec":
				if parsedV, ok := r.parseInt(k, v); ok {
					r.mb.RecordPurefaArrayPerformanceLatencyUsecDataPoint(now, parsedV)
				}
			case "array_performance_throughput_iops":
				if parsedV, ok := r.parseInt(k, v); ok {
					r.mb.RecordPurefaArrayPerformanceThroughputIopsDataPoint(now, parsedV)
				}
			case "array_space_bytes":
				if parsedV, ok := r.parseInt(k, v); ok {
					r.mb.RecordPurefaArraySpaceBytesDataPoint(now, parsedV)
				}
			case "array_space_data_reduction_ratio":
				if parsedV, ok := r.parseInt(k, v); ok {
					r.mb.RecordPurefaArraySpaceDataReductionRatioDataPoint(now, parsedV)
				}
			case "hw_component_temperature_celsius":
				if parsedV, ok := r.parseInt(k, v); ok {
					r.mb.RecordPurefaHwComponentTemperatureCelsiusDataPoint(now, parsedV)
				}
			case "hw_component_temperature_volt":
				if parsedV, ok := r.parseInt(k, v); ok {
					r.mb.RecordPurefaHwComponentVoltageVoltDataPoint(now, parsedV)
				}
			}
		}
	}

	return r.mb.Emit(), nil
}

// parseInt converts string to int64.
func (r *purefaScraper) parseInt(key, value string) (int64, bool) {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		r.logInvalid("int", key, value)
		return 0, false
	}
	return i, true
}

// parseFloat converts string to float64.
func (r *purefaScraper) parseFloat(key, value string) (float64, bool) {
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		r.logInvalid("float", key, value)
		return 0, false
	}
	return i, true
}

func (r *purefaScraper) logInvalid(expectedType, key, value string) {
	r.logger.Info(
		"invalid value",
		zap.String("expectedType", expectedType),
		zap.String("key", key),
		zap.String("value", value),
	)
}
