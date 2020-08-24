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

	"github.com/jmozah/intOS-dfs/pkg/dfs"
	p "github.com/jmozah/intOS-dfs/pkg/pod"
	"resenje.org/jsonhttp"
)

type PodCreateResponse struct {
	Reference string `json:"reference"`
}

func (h *Handler) PodCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		return
	}

	user := r.FormValue("user")
	password := r.FormValue("password")
	pod := r.FormValue("pod")
	if user == "" {
		jsonhttp.BadRequest(w, "create pod: \"user\" argument missing")
		return
	}
	if password == "" {
		jsonhttp.BadRequest(w, "create pod: \"password\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "create pod: \"pod\" argument missing")
		return
	}

	// create pod
	_, err := h.dfsAPI.CreatePod(user, pod, password)
	if err != nil {
		if err == dfs.ErrInvalidUserName || err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrInvalidPodName ||
			err == p.ErrTooLongPodName ||
			err == p.ErrPodAlreadyExists ||
			err == p.ErrMaxPodsReached {
			fmt.Println("create pod: ", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("create pod: ", err)
		jsonhttp.InternalServerError(w, err)
		return
	}

	jsonhttp.Created(w, nil)
}
