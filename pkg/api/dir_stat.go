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

func (h *Handler) DirectoryStatHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	dir := r.FormValue("dir")
	if user == "" {
		jsonhttp.BadRequest(w, "stat: \"user\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "stat: \"pod\" argument missing")
		return
	}
	if dir == "" {
		jsonhttp.BadRequest(w, "stat: \"dir\" argument missing")
		return
	}

	// stat directory
	ds, err := h.dfsAPI.DirectoryStat(user, pod, dir)
	if err != nil {
		fmt.Println("stat dir: %w", err)
		jsonhttp.InternalServerError(w, err)
	}

	jsonhttp.OK(w, ds)
}
