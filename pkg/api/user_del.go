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

func (h *Handler) UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	if password == "" {
		jsonhttp.BadRequest(w, "delete: \"password\" argument missing")
		return
	}

	// get values from cookie
	userName, sessionId, err := cookie.GetUserNameAndSessionId(r)
	if err != nil {
		fmt.Println("delete: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if userName == "" {
		jsonhttp.BadRequest(w, "delete: \"user\" parameter missing in cookie")
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "delete: \"cookie-id\" parameter missing in cookie")
		return
	}

	// delete user
	err = h.dfsAPI.DeleteUser(userName, password, sessionId, w)
	if err != nil {
		if err == u.ErrInvalidUserName ||
			err == u.ErrInvalidPassword ||
			err == u.ErrUserNotLoggedIn {
			fmt.Println("delete: ", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("delete: ", err)
		jsonhttp.InternalServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
