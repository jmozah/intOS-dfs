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
)

func (h *Handler) FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	podFile := r.FormValue("pod_file")
	if user == "" {
		jsonhttp.BadRequest(w, "delete: \"user\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "delete: \"pod\" argument missing")
		return
	}
	if podFile == "" {
		jsonhttp.BadRequest(w, "delete: \"path_in_pod\" argument missing")
		return
	}

	// delete file
	err := h.dfsAPI.DeleteFile(user, pod, podFile)
	if err != nil {
		fmt.Println("delete: %w", err)
		jsonhttp.InternalServerError(w, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
