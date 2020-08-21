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
	"io/ioutil"
	"net/http"

	"resenje.org/jsonhttp"
)

type uploadFiletResponse struct {
	Reference string `json:"reference"`
}

func (h *Handler) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	podDir := r.FormValue("pod_dir")
	blockSize := r.FormValue("block_size")
	fileName := r.FormValue("file_name")
	if user == "" {
		jsonhttp.BadRequest(w, "upload: \"user\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "upload: \"pod\" argument missing")
		return
	}
	if podDir == "" {
		jsonhttp.BadRequest(w, "upload: \"pod_dir\" argument missing")
		return
	}
	if blockSize == "" {
		jsonhttp.BadRequest(w, "upload: \"block_size\" argument missing")
		return
	}
	if fileName == "" {
		jsonhttp.BadRequest(w, "upload: \"file_namee\" argument missing")
		return
	}
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jsonhttp.BadRequest(w, "missing body")
		return
	}
	fileSize := r.ContentLength

	// upload file to bee
	reference, err := h.dfsAPI.UploadFile(user, pod, fileName, fileSize, r.Body, podDir, blockSize)
	if err != nil {
		fmt.Println("upload: %w", err)
		jsonhttp.InternalServerError(w, err)
	}

	w.Header().Set("ETag", fmt.Sprintf("%q", reference))
	jsonhttp.OK(w, &uploadFiletResponse{
		Reference: reference,
	})
}
