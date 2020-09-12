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

	"github.com/jmozah/intOS-dfs/pkg/dfs"

	"resenje.org/jsonhttp"

	"github.com/jmozah/intOS-dfs/pkg/cookie"
)

func (h *Handler) FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	podFile := r.FormValue("file")
	if podFile == "" {
		h.logger.Errorf("file delete: \"file\" argument missing")
		jsonhttp.BadRequest(w, "file delete: \"file\" argument missing")
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		h.logger.Errorf("file delete: invalid cookie: %v", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		h.logger.Errorf("file delete: \"cookie-id\" parameter missing in cookie")
		jsonhttp.BadRequest(w, "file delete: \"cookie-id\" parameter missing in cookie")
		return
	}

	// delete file
	err = h.dfsAPI.DeleteFile(podFile, sessionId)
	if err != nil {
		if err == dfs.ErrPodNotOpen {
			h.logger.Errorf("file delete: %v", err)
			jsonhttp.BadRequest(w, "file delete: "+err.Error())
			return
		}
		h.logger.Errorf("file delete: %v", err)
		jsonhttp.InternalServerError(w, "file delete: "+err.Error())
		return
	}

	jsonhttp.OK(w, "file deleted successfully")
}
