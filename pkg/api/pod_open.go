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
	"net/http"

	"resenje.org/jsonhttp"

	"github.com/jmozah/intOS-dfs/pkg/cookie"
	"github.com/jmozah/intOS-dfs/pkg/dfs"
	p "github.com/jmozah/intOS-dfs/pkg/pod"
)

type PodOpenResponse struct {
	Reference string `json:"reference"`
}

func (h *Handler) PodOpenHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	pod := r.FormValue("pod")
	if password == "" {
		h.logger.Errorf("open pod: \"password\" argument missing")
		jsonhttp.BadRequest(w, "open pod: \"password\" argument missing")
		return
	}
	if pod == "" {
		h.logger.Errorf("open pod: \"pod\" argument missing")
		jsonhttp.BadRequest(w, "open pod: \"pod\" argument missing")
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		h.logger.Errorf("open pod: invalid cookie: %v", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		h.logger.Errorf("open pod: \"cookie-id\" parameter missing in cookie")
		jsonhttp.BadRequest(w, "open pod: \"cookie-id\" parameter missing in cookie")
		return
	}

	// open pod
	_, err = h.dfsAPI.OpenPod(pod, password, sessionId)
	if err != nil {
		if err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrInvalidPodName {
			h.logger.Errorf("open pod: %v", err)
			jsonhttp.BadRequest(w, "open pod: "+err.Error())
			return
		}
		h.logger.Errorf("open pod: %v", err)
		jsonhttp.InternalServerError(w, "open pod: "+err.Error())
		return
	}

	jsonhttp.OK(w, "pod opened successfully")
}
