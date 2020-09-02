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
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strconv"

	"github.com/jmozah/intOS-dfs/pkg/cookie"
	"resenje.org/jsonhttp"
)

const (
	avatarHeight = 128
	avatarWidth  = 128
)

func (h *Handler) SaveUserAvatarHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("avatar: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "avatar: \"cookie-id\" parameter missing in cookie")
		return
	}

	w.Header().Set("Content-Type", " application/json")
	//  get the avatar parameter from the multi part
	err = r.ParseMultipartForm(defaultMaxMemory)
	if err != nil {
		fmt.Println("avatar: ", err)
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "avatar: " + err.Error()})
		return
	}
	avatars := r.MultipartForm.File["avatar"]
	if len(avatars) == 0 {
		fmt.Println("avatar: ", err)
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "avatar: parameter \"avatar\" missing"})
		return
	}
	if len(avatars) > 1 {
		fmt.Println("avatar: multiple avatars not allowed")
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "avatar: multiple avatars not allowed"})
		return
	}

	// Read the avatar file data
	fd, err := avatars[0].Open()
	if err != nil {
		fmt.Println("avatar: ", err)
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "avatar: " + err.Error()})
		return
	}
	data := make([]byte, avatars[0].Size)
	_, err = fd.Read(data)
	if err != nil {
		fmt.Println("avatar: ", err)
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "avatar: " + err.Error()})
		return
	}
	err = fd.Close()
	if err != nil {
		fmt.Println("avatar: ", err)
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "avatar: " + err.Error()})
		return
	}

	// check the avatar size
	reader := bytes.NewReader(data)
	im, _, err := image.DecodeConfig(reader)
	if err != nil {
		fmt.Println("avatar: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "avatar: " + err.Error()})
		return
	}
	if im.Height > avatarHeight || im.Width > avatarWidth {
		fmt.Println("avatar: invalid avatar size")
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "avatar: size should be less than 128x128"})
		return
	}

	// save avatar with .avatar extension
	err = h.dfsAPI.SaveAvatar(sessionId, data)
	if err != nil {
		fmt.Println("avatar: ", err)
		jsonhttp.BadRequest(w, &ErrorMessage{Err: "avatar: " + err.Error()})
		return
	}

	jsonhttp.OK(w, nil)
}

func (h *Handler) GetUserAvatarHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("avatar: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "avatar: \"cookie-id\" parameter missing in cookie")
		return
	}

	data, err := h.dfsAPI.GetAvatar(sessionId)
	if err != nil {
		fmt.Println("avatar: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "avatar: " + err.Error()})
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))

	reader := bytes.NewReader(data)
	_, err = io.Copy(w, reader)
	if err != nil {
		fmt.Println("avatar:", err)
		w.Header().Set("Content-Type", " application/json")
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "stat dir: " + err.Error()})
	}
}
