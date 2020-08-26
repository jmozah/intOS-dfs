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

func (h *Handler) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")
	if user == "" {
		jsonhttp.BadRequest(w, "login: \"user\" argument missing")
		return
	}
	if password == "" {
		jsonhttp.BadRequest(w, "login: \"password\" argument missing")
		return
	}

	// login user
	err := h.dfsAPI.LoginUser(user, password, w, "")
	if err != nil {
		if err == u.ErrUserAlreadyLoggedIn ||
			err == u.ErrInvalidUserName ||
			err == u.ErrInvalidPassword {
			fmt.Println("login: ", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("login: ", err)
		jsonhttp.InternalServerError(w, err)
		return
	}

	jsonhttp.OK(w, nil)
}
