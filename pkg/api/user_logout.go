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

	"github.com/jmozah/intOS-dfs/pkg/cookie"
	u "github.com/jmozah/intOS-dfs/pkg/user"
	"resenje.org/jsonhttp"
)

func (h *Handler) UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	userName, sessionId, err := cookie.GetUserNameAndSessionId(r)
	if err != nil {
		fmt.Println("logout: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if userName == "" {
		jsonhttp.BadRequest(w, "logout: \"user\" parameter missing in cookie")
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "logout: \"cookie-id\" parameter missing in cookie")
		return
	}

	// logout user
	err = h.dfsAPI.LogoutUser(userName, sessionId, w)
	if err != nil {
		w.Header().Set("Content-Type", " application/json")
		if err == u.ErrUserNotLoggedIn || err == u.ErrInvalidUserName {
			fmt.Println("logout: ", err)
			jsonhttp.BadRequest(w, &ErrorMessage{Err: "logout: " + err.Error()})
			return
		}
		fmt.Println("logout: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "logout: " + err.Error()})
		return
	}

	jsonhttp.OK(w, nil)
}
