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
)

func (h *Handler) PodDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	userName, sessionId, podName, err := cookie.GetUserNameSessionIdAndPodName(r)
	if err != nil {
		fmt.Println("delete pod: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if userName == "" {
		jsonhttp.BadRequest(w, "delete pod: \"user\" parameter missing in cookie")
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "delete pod: \"cookie-id\" parameter missing in cookie")
		return
	}
	if podName == "" {
		jsonhttp.BadRequest(w, "delete pod: \"pod\" parameter missing in cookie")
		return
	}

	// delete pod
	err = h.dfsAPI.DeletePod(userName, podName, sessionId, w, r)
	if err != nil {
		if err == dfs.ErrInvalidUserName || err == dfs.ErrUserNotLoggedIn {
			fmt.Println("delete pod:", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("delete pod:", err)
		jsonhttp.InternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
