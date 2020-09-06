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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"resenje.org/jsonhttp"

	"github.com/jmozah/intOS-dfs/pkg/cookie"
	"github.com/jmozah/intOS-dfs/pkg/user"
)

type receiveFileResponse struct {
	FileName  string `json:"file_name"`
	Reference string `json:"reference"`
}

func (h *Handler) FileShareHandler(w http.ResponseWriter, r *http.Request) {
	podFile := r.FormValue("file")
	if podFile == "" {
		jsonhttp.BadRequest(w, "share: \"file\" argument missing")
		return
	}
	destinationRef := r.FormValue("to")
	if destinationRef == "" {
		jsonhttp.BadRequest(w, "share: \"to\" argument missing")
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("share: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "share: \"cookie-id\" parameter missing in cookie")
		return
	}

	w.Header().Set("Content-Type", " application/json")
	outEntry, err := h.dfsAPI.ShareFile(podFile, destinationRef, sessionId)
	if err != nil {
		fmt.Println("share: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "share: " + err.Error()})
		return
	}

	jsonhttp.OK(w, outEntry)
}

func (h *Handler) FileReceiveHandler(w http.ResponseWriter, r *http.Request) {
	// get the outbox entry
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jsonhttp.BadRequest(w, "receive: no data in body")
		return
	}
	inboxEntry := user.InboxEntry{}
	err = json.Unmarshal(data, &inboxEntry)
	if err != nil {
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "share: " + err.Error()})
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("share: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "share: \"cookie-id\" parameter missing in cookie")
		return
	}

	w.Header().Set("Content-Type", " application/json")
	err = h.dfsAPI.ReceiveFile(sessionId, inboxEntry)
	if err != nil {
		fmt.Println("share: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "share: " + err.Error()})
		return
	}

	jsonhttp.OK(w, &receiveFileResponse{
		FileName:  inboxEntry.FilePath,
		Reference: inboxEntry.FileMetaHash,
	})
}
