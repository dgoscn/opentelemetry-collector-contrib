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

package directoriescraper

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/scrapertest"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/scrapertest/golden"
)

func TestScraper(t *testing.T) {
	f := NewFactory()
	cfg := f.CreateDefaultConfig().(*Config)
	scraper := newPurefaScraper(componenttest.NewNopReceiverCreateSettings(), cfg)
	scraper.newClient = func(endpoint string, timeout time.Duration) (client, error) {
		return &fakeClient{}, nil
	}

	actualMetrics, err := scraper.scrape(context.Background())
	require.NoError(t, err)

	expectedFile := filepath.Join("testdata", "expected_metrics", "test_scraper", "expected.json")
	expectedMetrics, err := golden.ReadMetrics(expectedFile)
	require.NoError(t, err)

	require.NoError(t, scrapertest.CompareMetrics(expectedMetrics, actualMetrics))
}

func TestScraperWithoutDirectionAttribute(t *testing.T) {
	f := NewFactory()
	cfg := f.CreateDefaultConfig().(*Config)
	scraper := newPurefaScraper(componenttest.NewNopReceiverCreateSettings(), cfg)
	scraper.emitMetricsWithDirectionAttribute = false
	scraper.emitMetricsWithoutDirectionAttribute = true
	scraper.newClient = func(endpoint string, timeout time.Duration) (client, error) {
		return &fakeClient{}, nil
	}

	actualMetrics, err := scraper.scrape(context.Background())
	require.NoError(t, err)

	expectedFile := filepath.Join("testdata", "expected_metrics", "test_scraper", "expected_without_direction.json")
	expectedMetrics, err := golden.ReadMetrics(expectedFile)
	require.NoError(t, err)

	require.NoError(t, scrapertest.CompareMetrics(expectedMetrics, actualMetrics))
}
