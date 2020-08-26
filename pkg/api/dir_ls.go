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

	"github.com/jmozah/intOS-dfs/pkg/cookie"
	"github.com/jmozah/intOS-dfs/pkg/dfs"
	p "github.com/jmozah/intOS-dfs/pkg/pod"
)

type ListFileResponse struct {
	Files       []string `json:"files,omitempty"`
	Directories []string `json:"directories,omitempty"`
}

func (h *Handler) DirectoryLsHandler(w http.ResponseWriter, r *http.Request) {
	dir := r.FormValue("dir")
	if dir == "" {
		jsonhttp.BadRequest(w, "ls dir: \"dir\" argument missing")
		return
	}

	// get values from cookie
	userName, sessionId, podName, err := cookie.GetUserNameSessionIdAndPodName(r)
	if err != nil {
		fmt.Println("ls dir: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if userName == "" {
		jsonhttp.BadRequest(w, "ls dir: \"user\" parameter missing in cookie")
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "ls dir: \"cookie-id\" parameter missing in cookie")
		return
	}
	if podName == "" {
		jsonhttp.BadRequest(w, "ls dir: \"pod\" parameter missing in cookie")
		return
	}

	// list directory
	fl, dl, err := h.dfsAPI.ListDir(userName, podName, dir, sessionId)
	if err != nil {
		if err == dfs.ErrInvalidUserName || err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrPodNotOpened {
			fmt.Println("ls dir: ", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("ls dir: ", err)
		jsonhttp.InternalServerError(w, err)
		return
	}

	jsonhttp.OK(w, &ListFileResponse{
		Files:       fl,
		Directories: dl,
	})
}
