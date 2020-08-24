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

	"github.com/jmozah/intOS-dfs/pkg/dfs"
	p "github.com/jmozah/intOS-dfs/pkg/pod"
)

func (h *Handler) DirectoryRmdirHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		return
	}

	user := r.FormValue("user")
	pod := r.FormValue("pod")
	dir := r.FormValue("dir")
	if user == "" {
		jsonhttp.BadRequest(w, "rmdir: \"user\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "rmdir: \"pod\" argument missing")
		return
	}
	if dir == "" {
		jsonhttp.BadRequest(w, "rmdir: \"dir\" argument missing")
		return
	}

	// remove directory
	err := h.dfsAPI.RmDir(user, pod, dir)
	if err != nil {
		if err == dfs.ErrInvalidUserName || err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrPodNotOpened {
			fmt.Println("rmdir: ", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("rmdir: ", err)
		jsonhttp.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
