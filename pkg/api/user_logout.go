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

func (h *Handler) UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		h.logger.Errorf("logout: invalid cookie: %v", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		h.logger.Errorf("logout: \"cookie-id\" parameter missing in cookie")
		jsonhttp.BadRequest(w, "logout: \"cookie-id\" parameter missing in cookie")
		return
	}

	// logout user
	err = h.dfsAPI.LogoutUser(sessionId, w)
	if err != nil {
		if err == u.ErrUserNotLoggedIn || err == u.ErrInvalidUserName {
			h.logger.Errorf("logout: %v", err)
			jsonhttp.BadRequest(w, "logout: "+err.Error())
			return
		}
		h.logger.Errorf("logout: %v", err)
		jsonhttp.InternalServerError(w, "logout: "+err.Error())
		return
	}
	jsonhttp.OK(w, "used logged out successfully")
}
