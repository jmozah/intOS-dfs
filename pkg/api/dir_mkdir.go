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

func (h *Handler) DirectoryMkdirHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	dirToCreate := r.FormValue("dir")
	if user == "" {
		jsonhttp.BadRequest(w, "mkdir: \"user\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "mkdir: \"pod\" argument missing")
		return
	}
	if dirToCreate == "" {
		jsonhttp.BadRequest(w, "mkdir: \"dir\" argument missing")
		return
	}

	// make directory
	err := h.dfsAPI.Mkdir(user, pod, dirToCreate)
	if err != nil {
		if err == dfs.ErrInvalidUserName || err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrInvalidDirectory ||
			err == p.ErrTooLongDirectoryName ||
			err == p.ErrPodNotOpened {
			fmt.Println("mkdir: ", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("mkdir: ", err)
		jsonhttp.InternalServerError(w, err)
		return
	}

	jsonhttp.Created(w, nil)
}
