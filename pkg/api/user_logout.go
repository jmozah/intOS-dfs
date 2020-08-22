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

	u "github.com/jmozah/intOS-dfs/pkg/user"
	"resenje.org/jsonhttp"
)

func (h *Handler) UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	if user == "" {
		jsonhttp.BadRequest(w, "logout: \"user\" argument missing")
		return
	}

	// logout user
	err := h.dfsAPI.LogoutUser(user)
	if err != nil {
		if err == u.ErrUserNotLoggedIn || err == u.ErrInvalidUserName {
			fmt.Println("logout: ", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("logout: ", err)
		jsonhttp.InternalServerError(w, err)
		return
	}

	jsonhttp.OK(w, nil)
}
