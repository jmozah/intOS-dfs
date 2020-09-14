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
)

func (h *Handler) SaveUserNameHandler(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	middleName := r.FormValue("middle_name")
	surname := r.FormValue("surname")

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		h.logger.Errorf("save name: invalid cookie: %v", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		h.logger.Errorf("save name: \"cookie-id\" parameter missing in cookie")
		jsonhttp.BadRequest(w, "save name: \"cookie-id\" parameter missing in cookie")
		return
	}

	err = h.dfsAPI.SaveName(firstName, lastName, middleName, surname, sessionId)
	if err != nil {
		h.logger.Errorf("save name: %v", err)
		jsonhttp.InternalServerError(w, "save name: "+err.Error())
		return
	}
	jsonhttp.OK(w, nil)
}

func (h *Handler) GetUserNameHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		h.logger.Errorf("get name: invalid cookie: %v", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		h.logger.Errorf("get name: \"cookie-id\" parameter missing in cookie")
		jsonhttp.BadRequest(w, "get name: \"cookie-id\" parameter missing in cookie")
		return
	}

	name, err := h.dfsAPI.GetName(sessionId)
	if err != nil {
		h.logger.Errorf("get name: %v", err)
		jsonhttp.InternalServerError(w, "get name: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, name)
}
