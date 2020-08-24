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

type FileStatResponse struct {
	Blocks []BlockInfo
}

type BlockInfo struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Size      string `json:"size"`
}

func (h *Handler) FileStatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		return
	}

	user := r.FormValue("user")
	pod := r.FormValue("pod")
	podFile := r.FormValue("file")
	if user == "" {
		jsonhttp.BadRequest(w, "stat: \"user\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "stat: \"pod\" argument missing")
		return
	}
	if podFile == "" {
		jsonhttp.BadRequest(w, "upload: \"file\" argument missing")
		return
	}

	// get file stat
	stat, err := h.dfsAPI.FileStat(user, pod, podFile)
	if err != nil {
		fmt.Println("file stat: %w", err)
		jsonhttp.InternalServerError(w, err)
	}

	jsonhttp.OK(w, stat)
}
