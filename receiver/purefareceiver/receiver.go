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
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/purefareceiver/internal/metadata"
)

var (
	errEmptyEndpoint = errors.New("empty endpoint")
)

type purefaReceiver struct {
	cfg    *Config
	client *internal.PurefaClient

	array
	mb *metadata.MetricsBuilder
}

// New creates the Pure FlashArray receiver with the given configuration.
func newPureFlashArrayReceiver(settings component.ReceiverCreateSettings, cfg *Config) *receiver {
	return &purefaReceiver{
		cfg:    cfg,
		client: client,
	}
}

func newPureFlashArrayScraper(recv *purefaReceiver) (scraperhelper.Scraper, error) {
	return scraperhelper.NewScraper(typeStr, recv.scrape, scraperhelper.WithShutdown(recv.shutdown))
}

// Start begins collecting metrics from endpoint.
func (tailtracerRcvr *purefaReceiver) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (s *purefaReceiver) shutdown(context.Context) error {
	//s.cancel()
	return s.client.Shutdown()
}
