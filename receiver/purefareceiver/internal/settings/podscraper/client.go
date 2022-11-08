// Copyright 2022 The OpenTelemetry Authors
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
)

type client interface {
	Connect(ctx context.Context) error
	collectDataFromQuery(ctx context.Context, query *monitoringQuery) ([]map[string]string, error)
	Close() error
}

type newPurefaClientFunc func(endpoint string, timeout time.Duration) (client, error)

func newPurefaClient(endpoint string, timeout time.Duration) (client, error) {
	newClient, err := purefa.New(endpoint)
	if err != nil {
		return nil, err
	}

	newClient.Timeout = timeout
	return newClient, nil
}
