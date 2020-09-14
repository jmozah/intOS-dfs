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
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

type receiveFileResponse struct {
	FileName  string `json:"file_name"`
	Reference string `json:"reference"`
}

type SharingReference struct {
	Reference string `json:"sharing_reference"`
}

func (h *Handler) FileShareHandler(w http.ResponseWriter, r *http.Request) {
	podFile := r.FormValue("file")
	if podFile == "" {
		h.logger.Errorf("share: \"file\" argument missing")
		jsonhttp.BadRequest(w, "share: \"file\" argument missing")
		return
	}
	destinationRef := r.FormValue("to")
	if destinationRef == "" {
		h.logger.Errorf("share: \"to\" argument missing")
		jsonhttp.BadRequest(w, "share: \"to\" argument missing")
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		h.logger.Errorf("share: invalid cookie: %v", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		h.logger.Errorf("share: \"cookie-id\" parameter missing in cookie")
		jsonhttp.BadRequest(w, "share: \"cookie-id\" parameter missing in cookie")
		return
	}

	sharingRef, err := h.dfsAPI.ShareFile(podFile, destinationRef, sessionId)
	if err != nil {
		h.logger.Errorf("share: %v", err)
		jsonhttp.InternalServerError(w, "share: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, &SharingReference{
		Reference: sharingRef,
	})
}

func (h *Handler) FileReceiveHandler(w http.ResponseWriter, r *http.Request) {
	sharingRefString := r.FormValue("ref")
	if sharingRefString == "" {
		h.logger.Errorf("receive: \"ref\" argument missing")
		jsonhttp.BadRequest(w, "receive: \"ref\" argument missing")
		return
	}

	dir := r.FormValue("dir")
	if dir == "" {
		h.logger.Errorf("receive: \"dir\" argument missing")
		jsonhttp.BadRequest(w, "receive: \"dir\" argument missing")
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		h.logger.Errorf("receive: invalid cookie: %v", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		h.logger.Errorf("receive: \"cookie-id\" parameter missing in cookie")
		jsonhttp.BadRequest(w, "share: \"cookie-id\" parameter missing in cookie")
		return
	}

	sharingRef, err := utils.ParseSharingReference(sharingRefString)
	if err != nil {
		h.logger.Errorf("receive: invalid reference: ", err)
		jsonhttp.BadRequest(w, "receive: invalid reference:"+err.Error())
		return
	}

	filePath, fileRef, err := h.dfsAPI.ReceiveFile(sessionId, sharingRef, dir)
	if err != nil {
		h.logger.Errorf("receive: %v", err)
		jsonhttp.InternalServerError(w, "receive: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, &receiveFileResponse{
		FileName:  filePath,
		Reference: fileRef,
	})
}
