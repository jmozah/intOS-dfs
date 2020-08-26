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
	"fmt"
	"net/http"

	"resenje.org/jsonhttp"

	"github.com/jmozah/intOS-dfs/pkg/cookie"
)

type FileStatResponse struct {
	Blocks []BlockInfo
}

type BlockInfo struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Size      string `json:"size"`
}

func (h *Handler) FileStatHandler(w http.ResponseWriter, r *http.Request) {
	podFile := r.FormValue("file")
	if podFile == "" {
		jsonhttp.BadRequest(w, "file stat: \"file\" argument missing")
		return
	}

	// get values from cookie
	userName, sessionId, podName, err := cookie.GetUserNameSessionIdAndPodName(r)
	if err != nil {
		fmt.Println("file stat: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if userName == "" {
		jsonhttp.BadRequest(w, "file stat: \"user\" parameter missing in cookie")
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "file stat: \"cookie-id\" parameter missing in cookie")
		return
	}
	if podName == "" {
		jsonhttp.BadRequest(w, "file stat: \"pod\" parameter missing in cookie")
		return
	}
	w.Header().Set("Content-Type", " application/json")

	// restart the cookie expiry
	err = cookie.ResetSessionExpiry(r, w)
	if err != nil {
		jsonhttp.BadRequest(w, &ErrorMessage{err: "file stat: " + err.Error()})
		return
	}

	// get file stat
	stat, err := h.dfsAPI.FileStat(userName, podName, podFile, sessionId)
	if err != nil {
		fmt.Println("file stat: %w", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{err: "file stat: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, stat)
}
