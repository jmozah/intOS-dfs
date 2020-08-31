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

func (h *Handler) PodCloseHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("close pod: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "close pod: \"cookie-id\" parameter missing in cookie")
		return
	}

	// close pod
	err = h.dfsAPI.ClosePod(sessionId)
	if err != nil {
		w.Header().Set("Content-Type", " application/json")
		if err == dfs.ErrPodNotOpen || err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrPodNotOpened {
			fmt.Println("close pod:", err)
			jsonhttp.BadRequest(w, &ErrorMessage{Err: "close pod: " + err.Error()})
			return
		}
		fmt.Println("close pod: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "close pod: " + err.Error()})
		return
	}

	jsonhttp.OK(w, nil)
}
