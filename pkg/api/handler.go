/*
Copyright © 2020 intOS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"github.com/jmozah/intOS-dfs/pkg/dfs"
	"github.com/jmozah/intOS-dfs/pkg/logging"
	"github.com/jmozah/intOS-dfs/pkg/web"
)

type Handler struct {
	dfsAPI      *dfs.DfsAPI
	WebHandlers *web.Web
	logger      logging.Logger
}

func NewHandler(dataDir, beeHost, beePort string, logger logging.Logger) (*Handler, error) {
	api, err := dfs.NewDfsAPI(dataDir, beeHost, beePort, logger)
	if err != nil {
		return nil, dfs.ErrBeeClient
	}
	return &Handler{
		dfsAPI:      api,
		WebHandlers: web.NewWeb(logger),
		logger:      logger,
	}, nil
}
