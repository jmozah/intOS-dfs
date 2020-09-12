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

	"github.com/jmozah/intOS-dfs/pkg/cookie"
	u "github.com/jmozah/intOS-dfs/pkg/user"
)

func (h *Handler) UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	if password == "" {
		h.logger.Errorf("delete: \"password\" argument missing")
		jsonhttp.BadRequest(w, "delete: \"password\" argument missing")
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		h.logger.Errorf("delete: invalid cookie: %v", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		h.logger.Errorf("delete: \"cookie-id\" parameter missing in cookie")
		jsonhttp.BadRequest(w, "delete: \"cookie-id\" parameter missing in cookie")
		return
	}

	// delete user
	err = h.dfsAPI.DeleteUser(password, sessionId, w)
	if err != nil {
		if err == u.ErrInvalidUserName ||
			err == u.ErrInvalidPassword ||
			err == u.ErrUserNotLoggedIn {
			h.logger.Errorf("delete: %v", err)
			jsonhttp.BadRequest(w, "delete: "+err.Error())
			return
		}
		h.logger.Errorf("delete: %v", err)
		jsonhttp.InternalServerError(w, "delete: "+err.Error())
		return
	}
	jsonhttp.OK(w, "user deleted successfully")
}
