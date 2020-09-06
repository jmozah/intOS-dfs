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

	"resenje.org/jsonhttp"

	"github.com/jmozah/intOS-dfs/pkg/cookie"
)

func (h *Handler) GetUserSharingInboxHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("get share inbox: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "get share inbox: \"cookie-id\" parameter missing in cookie")
		return
	}

	sharingInbox, err := h.dfsAPI.GetUserSharingInbox(sessionId)
	if err != nil {
		fmt.Println("get share inbox: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "get share inbox: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, sharingInbox)
}

func (h *Handler) GetUserSharingOutboxHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("get share outbox: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "get share outbox: \"cookie-id\" parameter missing in cookie")
		return
	}

	sharingOutbox, err := h.dfsAPI.GetUserSharingOutbox(sessionId)
	if err != nil {
		fmt.Println("get share outbox: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "get share outbox: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, sharingOutbox)
}
