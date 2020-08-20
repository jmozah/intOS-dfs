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
	"io/ioutil"
	"net/http"

	"resenje.org/jsonhttp"
)

func (h *Handler) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	podDir := r.FormValue("pod_dir")
	blockSize := r.FormValue("block_size")
	if user == "" {
		jsonhttp.BadRequest(w, "argument missing: user ")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "argument missing: pod")
		return
	}
	if podDir == "" {
		jsonhttp.BadRequest(w, "argument missing: pod_dir")
		return
	}
	if blockSize == "" {
		jsonhttp.BadRequest(w, "argument missing: block_size")
		return
	}
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jsonhttp.BadRequest(w, "missing body")
		return
	}

	// TODO: copy file to bee

	jsonhttp.OK(w, nil)
}
