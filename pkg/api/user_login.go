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
	"net/http"

	"resenje.org/jsonhttp"

	u "github.com/jmozah/intOS-dfs/pkg/user"
)

func (h *Handler) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")
	if user == "" {
		h.logger.Errorf("login: \"user\" argument missing")
		jsonhttp.BadRequest(w, "login: \"user\" argument missing")
		return
	}
	if password == "" {
		h.logger.Errorf("login: \"password\" argument missing")
		jsonhttp.BadRequest(w, "login: \"password\" argument missing")
		return
	}

	// login user
	err := h.dfsAPI.LoginUser(user, password, w, "")
	if err != nil {
		if err == u.ErrUserAlreadyLoggedIn ||
			err == u.ErrInvalidUserName ||
			err == u.ErrInvalidPassword {
			h.logger.Errorf("login: %v", err)
			jsonhttp.BadRequest(w, "login: "+err.Error())
			return
		}
		h.logger.Errorf("login: %v", err)
		jsonhttp.InternalServerError(w, "login: "+err.Error())
		return
	}
	jsonhttp.OK(w, "user logged-in successfully")
}
