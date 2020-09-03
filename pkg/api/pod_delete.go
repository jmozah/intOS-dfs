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
	podName := r.FormValue("pod")
	if podName == "" {
		jsonhttp.BadRequest(w, "delete pod: \"pod\" parameter missing in cookie")
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("delete pod: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "delete pod: \"cookie-id\" parameter missing in cookie")
		return
	}

	// dont allow deletion of default pod name
	if podName == dfs.DefaultPodName {
		fmt.Println("delete pod: cannot delete default pod ", dfs.DefaultPodName)
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "delete pod: cannot delete default pod " + dfs.DefaultPodName})
		return
	}

	// delete pod
	err = h.dfsAPI.DeletePod(podName, sessionId)
	if err != nil {
		w.Header().Set("Content-Type", " application/json")
		if err == dfs.ErrUserNotLoggedIn {
			fmt.Println("delete pod:", err)
			jsonhttp.BadRequest(w, &ErrorMessage{Err: "delete pod: " + err.Error()})
			return
		}
		fmt.Println("delete pod:", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "delete pod: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
