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
)

func (h *Handler) FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	podFile := r.FormValue("file")
	if user == "" {
		jsonhttp.BadRequest(w, "download: \"user\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "download: \"pod\" argument missing")
		return
	}
	if podFile == "" {
		jsonhttp.BadRequest(w, "download: \"file\" argument missing")
		return
	}

	// download file from bee
	reader, reference, size, err := h.dfsAPI.DownloadFile(user, pod, podFile)
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