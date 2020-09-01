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
	"github.com/jmozah/intOS-dfs/pkg/cookie"
	"net/http"
	"resenje.org/jsonhttp"
	"time"
)

func (h *Handler) LoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, loginTimeout, err := cookie.GetSessionIdAndLoginTimeFromCookie(r)
		if err != nil {
			err1 := h.dfsAPI.LogoutUser(sessionId, w)
			if err1 == nil {
				jsonhttp.BadRequest(w, ErrorMessage{Err: "Logged out: invalid cookie"})
			}
			return
		}

		// if the expiry time is over, logout the user
		loginTime, err := time.Parse(time.RFC3339, loginTimeout)
		if err != nil {
			err1 := h.dfsAPI.LogoutUser(sessionId, w)
			if err1 == nil {
				jsonhttp.BadRequest(w, ErrorMessage{Err: "Logged out: invalid login timeout"})
			}
			return
		}
		if loginTime.Before(time.Now()) {
			err = h.dfsAPI.LogoutUser(sessionId, w)
			if err == nil {
				jsonhttp.BadRequest(w, ErrorMessage{Err: "Logging out as cookie login timeout expired"})
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}