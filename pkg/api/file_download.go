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
	"io"
	"net/http"

	"github.com/jmozah/intOS-dfs/pkg/dfs"

	"resenje.org/jsonhttp"

	"github.com/jmozah/intOS-dfs/pkg/cookie"
)

func (h *Handler) FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	podFile := r.FormValue("file")
	if podFile == "" {
		jsonhttp.BadRequest(w, "download: \"file\" argument missing")
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("download: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "download: \"cookie-id\" parameter missing in cookie")
		return
	}

	// download file from bee
	reader, reference, size, err := h.dfsAPI.DownloadFile(podFile, sessionId)
	if err != nil {
		if err == dfs.ErrPodNotOpen {
			fmt.Println("download:", err)
			jsonhttp.BadRequest(w, &ErrorMessage{Err: "download: " + err.Error()})
			return
		}
		w.Header().Set("Content-Type", " application/json")
		fmt.Println("download: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "download: " + err.Error()})
		return
	}

	w.Header().Set("ETag", fmt.Sprintf("%q", reference))
	w.Header().Set("Content-Length", size)

	_, err = io.Copy(w, reader)
	if err != nil {
		fmt.Println("download:", err)
		w.Header().Set("Content-Type", " application/json")
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "stat dir: " + err.Error()})
	}
}
