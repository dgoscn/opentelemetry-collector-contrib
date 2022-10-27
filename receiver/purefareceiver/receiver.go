// Copyright 2022 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package purefareceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver/internal/metadata"
)

type receiver struct {
	cfg    *Config
	client *internal.PurefaClient

	lastRun time.Time
	mb      *metadata.MetricsBuilder
}

func newPureFlashArrayReceiver(settings component.ReceiverCreateSettings, cfg *Config) *receiver {
	return &receiver{
		cfg:    cfg,
		client: client,
	}
}

func newPureFlashArrayScraper(recv *receiver) (scraperhelper.Scraper, error) {
	return scraperhelper.NewScraper(typeStr, recv.scrape, scraperhelper.WithShutdown(recv.shutdown))
}

func (s *receiver) shutdown(context.Context) error {
	return s.client.Shutdown()
}
