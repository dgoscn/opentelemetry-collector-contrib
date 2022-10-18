// Copyright 2022, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package purefareceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver"

import (
	"time"

	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`

	// Configure the service URL for the purefa service.
	Endpoint string `mapstructure:"endpoint"`
	Arrays   []*Array
	Settings *Settings

	// Metrics allows customizing scraped metrics representation.
	//Metrics metadata.MetricsSettings `mapstructure:"metrics"`

	//MetricsPath string `mapstructure:"metrics_path"`
	// Params the parameters to the metrics endpoint.
	//Params url.Values `mapstructure:"params,omitempty"`
	// Labels static labels
	//Labels map[string]string `mapstructure:"labels,omitempty"`
}

// This struct represents the Pure Storage Flash Array
type Array struct {
	Address string `mapstructure:"Address"`
	Token   string `mapstructure:"Token"`
}

type Settings struct {

	//CollectionInterval time.Duration `mapstructure:"collection_interval"`

	MetricsArrayCollectionInterval  time.Duration `mapstructure:"MetricsArrayCollectionInterval"`
	MetricsHostCollectionInterval   time.Duration `mapstructure:"MetricsHostCollectionInterval"`
	MetricsVolumeCollectionInterval time.Duration `mapstructure:"MetricsVolumeCollectionInterval"`
	MetricsPodCollectionInterval    time.Duration `mapstructure:"MetricsPodCollectionInterval"`
}
