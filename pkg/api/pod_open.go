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

type PodOpenResponse struct {
	Reference string `json:"reference"`
}

func (h *Handler) PodOpenHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	pod := r.FormValue("pod")
	if password == "" {
		jsonhttp.BadRequest(w, "open pod: \"password\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "open pod: \"pod\" argument missing")
		return
	}

	// get values from cookie
	userName, sessionId, existingPodName, err := cookie.GetUserNameSessionIdAndPodName(r)
	if err != nil {
		fmt.Println("delete: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if userName == "" {
		jsonhttp.BadRequest(w, "open pod: \"user\" parameter missing in cookie")
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "open pod: \"cookie-id\" parameter missing in cookie")
		return
	}

	// restart the cookie expiry
	err = cookie.ResetSessionExpiry(r, w)
	if err != nil {
		jsonhttp.BadRequest(w, err)
		return
	}

	// If a pod is already open, close the pod
	if existingPodName != "" {
		err = h.dfsAPI.ClosePod(userName, existingPodName, sessionId, w, r)
		if err != nil {
			w.Header().Set("Content-Type", " application/json")
			fmt.Println("open pod: could not close already open pod: ", err)
			jsonhttp.BadRequest(w, &ErrorMessage{err: "open pod: " + err.Error()})
			return
		}
	}

	// open pod
	_, err = h.dfsAPI.OpenPod(userName, pod, password, sessionId, w, r)
	if err != nil {
		w.Header().Set("Content-Type", " application/json")
		if err == dfs.ErrInvalidUserName || err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrInvalidPodName {
			fmt.Println("open pod:", err)
			jsonhttp.BadRequest(w, &ErrorMessage{err: "open pod: " + err.Error()})
			return
		}
		fmt.Println("open pod: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{err: "open pod: " + err.Error()})
		return
	}

	jsonhttp.OK(w, nil)
}
