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
	userName, sessionId, podName, err := cookie.GetUserNameSessionIdAndPodName(r)
	if err != nil {
		fmt.Println("download: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if userName == "" {
		jsonhttp.BadRequest(w, "download: \"user\" parameter missing in cookie")
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "download: \"cookie-id\" parameter missing in cookie")
		return
	}
	if podName == "" {
		jsonhttp.BadRequest(w, "download: \"pod\" parameter missing in cookie")
		return
	}

	// restart the cookie expiry
	err = cookie.ResetSessionExpiry(r, w)
	if err != nil {
		jsonhttp.BadRequest(w, err)
		return
	}

	// download file from bee
	reader, reference, size, err := h.dfsAPI.DownloadFile(userName, podName, podFile, sessionId)
	if err != nil {
		fmt.Println("download: ", err)
		jsonhttp.InternalServerError(w, err)
		return
	}

	w.Header().Set("ETag", fmt.Sprintf("%q", reference))
	w.Header().Set("Content-Length", size)
	_, err = io.Copy(w, reader)
	if err != nil {
		fmt.Println("download:", err)
		jsonhttp.InternalServerError(w, err)
	}
}
