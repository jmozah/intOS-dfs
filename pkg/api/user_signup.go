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

type UserSignupResponse struct {
	Reference string `json:"reference"`
	Mnemonic  string `json:"mnemonic"`
}

func (h *Handler) UserSignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		return
	}

	user := r.FormValue("user")
	password := r.FormValue("password")
	if user == "" {
		jsonhttp.BadRequest(w, "signup: \"user\" argument missing")
		return
	}
	if password == "" {
		jsonhttp.BadRequest(w, "signup: \"password\" argument missing")
		return
	}

	// create user
	reference, mnemonic, err := h.dfsAPI.CreateUser(user, password)
	if err != nil {
		if err == u.ErrUserAlreadyPresent {
			fmt.Println("signup: ", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("signup: ", err)
		jsonhttp.InternalServerError(w, err)
		return
	}

	// send the response
	jsonhttp.Created(w, &UserSignupResponse{
		Reference: reference,
		Mnemonic:  mnemonic,
	})
}
