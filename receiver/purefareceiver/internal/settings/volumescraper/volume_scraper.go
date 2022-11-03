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

package volumescraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver/internal/settings/volumescraper"

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/net"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/featuregate"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/scrapererror"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/processor/filterset"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver/internal/settings/volumescraper/internal/metadata"
)

const (
	volumeMetricsLen      = 4
	connectionsMetricsLen = 1
)

// scraper for Volume Metrics
type scraper struct {
	settings  component.ReceiverCreateSettings
	config    *Config
	mb        *metadata.MetricsBuilder
	startTime pcommon.Timestamp
	includeFS filterset.FilterSet
	excludeFS filterset.FilterSet

	// for mocking
	bootTime                             func() (uint64, error)
	ioCounters                           func(bool) ([]net.IOCountersStat, error)
	connections                          func(string) ([]net.ConnectionStat, error)
	conntrack                            func() ([]net.FilterStat, error)
	emitMetricsWithDirectionAttribute    bool
	emitMetricsWithoutDirectionAttribute bool
}

// newVolumeScraper creates a set of Volume related metrics
func newVolumeScraper(_ context.Context, settings component.ReceiverCreateSettings, cfg *Config) (*scraper, error) {
	scraper := &scraper{
		settings:                             settings,
		config:                               cfg,
		bootTime:                             host.BootTime,
		ioCounters:                           net.IOCounters,
		connections:                          net.Connections,
		conntrack:                            net.FilterCounters,
		emitMetricsWithDirectionAttribute:    featuregate.GetRegistry().IsEnabled(internal.EmitMetricsWithDirectionAttributeFeatureGateID),
		emitMetricsWithoutDirectionAttribute: featuregate.GetRegistry().IsEnabled(internal.EmitMetricsWithoutDirectionAttributeFeatureGateID),
	}

	var err error

	if len(cfg.Include.Interfaces) > 0 {
		scraper.includeFS, err = filterset.CreateFilterSet(cfg.Include.Interfaces, &cfg.Include.Config)
		if err != nil {
			return nil, fmt.Errorf("error creating volume interface include filters: %w", err)
		}
	}

	if len(cfg.Exclude.Interfaces) > 0 {
		scraper.excludeFS, err = filterset.CreateFilterSet(cfg.Exclude.Interfaces, &cfg.Exclude.Config)
		if err != nil {
			return nil, fmt.Errorf("error creating volume interface exclude filters: %w", err)
		}
	}

	return scraper, nil
}

func (s *scraper) start(context.Context, component.Host) error {
	bootTime, err := s.bootTime()
	if err != nil {
		return err
	}

	s.startTime = pcommon.Timestamp(bootTime * 1e9)
	s.mb = metadata.NewMetricsBuilder(s.config.Metrics, s.settings.BuildInfo, metadata.WithStartTime(pcommon.Timestamp(bootTime*1e9)))
	return nil
}

func (s *scraper) scrape(_ context.Context) (pmetric.Metrics, error) {
	var errors scrapererror.ScrapeErrors

	err := s.recordVolumeCounterMetrics()
	if err != nil {
		errors.AddPartial(volumeMetricsLen, err)
	}

	err = s.recordVolumeConnectionsMetrics()
	if err != nil {
		errors.AddPartial(connectionsMetricsLen, err)
	}

	err = s.recordVolumeConntrackMetrics()
	if err != nil {
		errors.AddPartial(connectionsMetricsLen, err)
	}

	return s.mb.Emit(), errors.Combine()
}

func (s *scraper) recordVolumeCounterMetrics() error {
	now := pcommon.NewTimestampFromTime(time.Now())

	// get total stats only
	ioCounters, err := s.ioCounters( /*perVolumeInterfaceController=*/ true)
	if err != nil {
		return fmt.Errorf("failed to read volume IO stats: %w", err)
	}

	// filter volume interfaces by name
	ioCounters = s.filterByInterface(ioCounters)

	if len(ioCounters) > 0 {
		s.recordVolumePacketsMetric(now, ioCounters)
		s.recordVolumeDroppedPacketsMetric(now, ioCounters)
		s.recordVolumeErrorPacketsMetric(now, ioCounters)
		s.recordVolumeIOMetric(now, ioCounters)
	}

	return nil
}

func (s *scraper) recordVolumePacketsMetric(now pcommon.Timestamp, ioCountersSlice []net.IOCountersStat) {
	for _, ioCounters := range ioCountersSlice {
		if s.emitMetricsWithoutDirectionAttribute {
			s.mb.RecordSystemVolumePacketsTransmitDataPoint(now, int64(ioCounters.PacketsSent), ioCounters.Name)
			s.mb.RecordSystemVolumePacketsReceiveDataPoint(now, int64(ioCounters.PacketsRecv), ioCounters.Name)
		}
		if s.emitMetricsWithDirectionAttribute {
			s.mb.RecordSystemVolumePacketsDataPoint(now, int64(ioCounters.PacketsSent), ioCounters.Name, metadata.AttributeDirectionTransmit)
			s.mb.RecordSystemVolumePacketsDataPoint(now, int64(ioCounters.PacketsRecv), ioCounters.Name, metadata.AttributeDirectionReceive)
		}
	}
}

func (s *scraper) recordVolumeDroppedPacketsMetric(now pcommon.Timestamp, ioCountersSlice []net.IOCountersStat) {
	for _, ioCounters := range ioCountersSlice {
		if s.emitMetricsWithoutDirectionAttribute {
			s.mb.RecordSystemVolumeDroppedTransmitDataPoint(now, int64(ioCounters.Dropout), ioCounters.Name)
			s.mb.RecordSystemVolumeDroppedReceiveDataPoint(now, int64(ioCounters.Dropin), ioCounters.Name)
		}
		if s.emitMetricsWithDirectionAttribute {
			s.mb.RecordSystemVolumeDroppedDataPoint(now, int64(ioCounters.Dropout), ioCounters.Name, metadata.AttributeDirectionTransmit)
			s.mb.RecordSystemVolumeDroppedDataPoint(now, int64(ioCounters.Dropin), ioCounters.Name, metadata.AttributeDirectionReceive)
		}
	}
}

func (s *scraper) recordVolumeErrorPacketsMetric(now pcommon.Timestamp, ioCountersSlice []net.IOCountersStat) {
	for _, ioCounters := range ioCountersSlice {
		if s.emitMetricsWithoutDirectionAttribute {
			s.mb.RecordSystemVolumeErrorsTransmitDataPoint(now, int64(ioCounters.Errout), ioCounters.Name)
			s.mb.RecordSystemVolumeErrorsReceiveDataPoint(now, int64(ioCounters.Errin), ioCounters.Name)
		}
		if s.emitMetricsWithDirectionAttribute {
			s.mb.RecordSystemVolumeErrorsDataPoint(now, int64(ioCounters.Errout), ioCounters.Name, metadata.AttributeDirectionTransmit)
			s.mb.RecordSystemVolumeErrorsDataPoint(now, int64(ioCounters.Errin), ioCounters.Name, metadata.AttributeDirectionReceive)
		}
	}
}

func (s *scraper) recordVolumeIOMetric(now pcommon.Timestamp, ioCountersSlice []net.IOCountersStat) {
	for _, ioCounters := range ioCountersSlice {
		if s.emitMetricsWithoutDirectionAttribute {
			s.mb.RecordSystemVolumeIoTransmitDataPoint(now, int64(ioCounters.BytesSent), ioCounters.Name)
			s.mb.RecordSystemVolumeIoReceiveDataPoint(now, int64(ioCounters.BytesRecv), ioCounters.Name)
		}
		if s.emitMetricsWithDirectionAttribute {
			s.mb.RecordSystemVolumeIoDataPoint(now, int64(ioCounters.BytesSent), ioCounters.Name, metadata.AttributeDirectionTransmit)
			s.mb.RecordSystemVolumeIoDataPoint(now, int64(ioCounters.BytesRecv), ioCounters.Name, metadata.AttributeDirectionReceive)
		}
	}
}

func (s *scraper) recordVolumeConnectionsMetrics() error {
	now := pcommon.NewTimestampFromTime(time.Now())

	connections, err := s.connections("tcp")
	if err != nil {
		return fmt.Errorf("failed to read TCP connections: %w", err)
	}

	tcpConnectionStatusCounts := getTCPConnectionStatusCounts(connections)

	s.recordVolumeConnectionsMetric(now, tcpConnectionStatusCounts)
	return nil
}

func getTCPConnectionStatusCounts(connections []net.ConnectionStat) map[string]int64 {
	tcpStatuses := make(map[string]int64, len(allTCPStates))
	for _, state := range allTCPStates {
		tcpStatuses[state] = 0
	}

	for _, connection := range connections {
		tcpStatuses[connection.Status]++
	}
	return tcpStatuses
}

func (s *scraper) recordVolumeConnectionsMetric(now pcommon.Timestamp, connectionStateCounts map[string]int64) {
	for connectionState, count := range connectionStateCounts {
		s.mb.RecordSystemVolumeConnectionsDataPoint(now, count, metadata.AttributeProtocolTcp, connectionState)
	}
}

func (s *scraper) filterByInterface(ioCounters []net.IOCountersStat) []net.IOCountersStat {
	if s.includeFS == nil && s.excludeFS == nil {
		return ioCounters
	}

	filteredIOCounters := make([]net.IOCountersStat, 0, len(ioCounters))
	for _, io := range ioCounters {
		if s.includeInterface(io.Name) {
			filteredIOCounters = append(filteredIOCounters, io)
		}
	}
	return filteredIOCounters
}

func (s *scraper) includeInterface(interfaceName string) bool {
	return (s.includeFS == nil || s.includeFS.Matches(interfaceName)) &&
		(s.excludeFS == nil || !s.excludeFS.Matches(interfaceName))
}
