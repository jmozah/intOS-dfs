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
	"github.com/jmozah/intOS-dfs/pkg/pod"
)

type PodListResponse struct {
	Pods []string `json:"name"`
}

func (h *Handler) PodListHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("delete: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "ls pod: \"cookie-id\" parameter missing in cookie")
		return
	}

	// fetch pods and list them
	pods, err := h.dfsAPI.ListPods(sessionId)
	if err != nil {
		w.Header().Set("Content-Type", " application/json")
		if err == dfs.ErrUserNotLoggedIn ||
			err == pod.ErrPodNotOpened {
			fmt.Println("ls pod: ", err)
			jsonhttp.BadRequest(w, &ErrorMessage{Err: "ls pod: " + err.Error()})
			return
		}
		fmt.Println("ls pod: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "ls pod: " + err.Error()})
		return
	}

	if pods == nil {
		pods = make([]string, 0)
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, &PodListResponse{
		Pods: pods,
	})
}
