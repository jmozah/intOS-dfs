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

func (h *Handler) DirectoryRmdirHandler(w http.ResponseWriter, r *http.Request) {
	dir := r.FormValue("dir")
	if dir == "" {
		jsonhttp.BadRequest(w, "rmdir: \"dir\" argument missing")
		return
	}

	// get values from cookie
	userName, sessionId, podName, err := cookie.GetUserNameSessionIdAndPodName(r)
	if err != nil {
		fmt.Println("rmdir: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if userName == "" {
		jsonhttp.BadRequest(w, "rmdir: \"user\" parameter missing in cookie")
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "rmdir: \"cookie-id\" parameter missing in cookie")
		return
	}
	if podName == "" {
		jsonhttp.BadRequest(w, "rmdir: \"pod\" parameter missing in cookie")
		return
	}

	// restart the cookie expiry
	err = cookie.ResetSessionExpiry(r, w)
	if err != nil {
		w.Header().Set("Content-Type", " application/json")
		jsonhttp.BadRequest(w, &ErrorMessage{err: "rmdir: " + err.Error()})
		return
	}

	// remove directory
	err = h.dfsAPI.RmDir(userName, podName, dir, sessionId)
	if err != nil {
		w.Header().Set("Content-Type", " application/json")
		if err == dfs.ErrInvalidUserName || err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrPodNotOpened {
			fmt.Println("rmdir: ", err)
			jsonhttp.BadRequest(w, &ErrorMessage{err: "rmdir: " + err.Error()})
			return
		}
		fmt.Println("rmdir: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{err: "rmdir: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
