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
	"github.com/jmozah/intOS-dfs/pkg/cookie"
	"github.com/jmozah/intOS-dfs/pkg/pod"
	"io/ioutil"
	"net/http"
	"resenje.org/jsonhttp"
)

func (h *Handler) SaveUserNameHandler(w http.ResponseWriter, r *http.Request) {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("name: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "name: \"cookie-id\" parameter missing in cookie")
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jsonhttp.BadRequest(w, "name: " + err.Error())
		return
	}
	name := &pod.Name{}
	err = json.Unmarshal(data, &name)
	if err != nil {
		jsonhttp.BadRequest(w, "name: " + err.Error())
		return
	}

	err = h.dfsAPI.SaveNameFile(pod.NameFile, sessionId, data)
	if err != nil {
		jsonhttp.BadRequest(w, "name: " + err.Error())
		return
	}

	jsonhttp.OK(w, nil)
}

func (h *Handler) GetUserNameHandler(w http.ResponseWriter, r *http.Request)  {
	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("name: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "name: \"cookie-id\" parameter missing in cookie")
		return
	}

	name, err := h.dfsAPI.GetNameFile(pod.NameFile, sessionId)
	if err != nil {
		jsonhttp.InternalServerError(w, "name: " + err.Error())
		return
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, name)
}
