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

func (h *Handler) UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")
	if user == "" {
		jsonhttp.BadRequest(w, "delete: \"user\" argument missing")
		return
	}
	if password == "" {
		jsonhttp.BadRequest(w, "delete: \"password\" argument missing")
		return
	}

	// delete user
	err := h.dfsAPI.DeleteUser(user, password)
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
