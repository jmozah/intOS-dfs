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

type ListFileResponse struct {
	Files       []string `json:"files"`
	Directories []string `json:"directories"`
}

func (h *Handler) DirectoryLsHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	currentDir := r.FormValue("curr_dir")
	if user == "" {
		jsonhttp.BadRequest(w, "ls dir: \"user\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "ls dir: \"pod\" argument missing")
		return
	}
	if currentDir == "" {
		jsonhttp.BadRequest(w, "ls dir: \"curr_dir\" argument missing")
		return
	}

	// list directory
	fl, dl, err := h.dfsAPI.ListDir(user, pod, currentDir)
	if err != nil {
		fmt.Println("ls dir: %w", err)
		jsonhttp.InternalServerError(w, err)
	}

	jsonhttp.OK(w, &ListFileResponse{
		Files:       fl,
		Directories: dl,
	})
}
