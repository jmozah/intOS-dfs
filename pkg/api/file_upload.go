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
)

type uploadFiletResponse struct {
	References []Reference
}

type Reference struct {
	FileName  string `json:"file_name"`
	Reference string `json:"reference,omitempty"`
	Error     string `json:"error,omitempty"`
}

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

func (h *Handler) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	podDir := r.FormValue("pod_dir")
	blockSize := r.FormValue("block_size")
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

	//  get the files parameter from the multi part
	err := r.ParseMultipartForm(defaultMaxMemory)
	if err != nil {
		fmt.Println("upload: ", err)
		jsonhttp.BadRequest(w, err)
		return
	}
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		fmt.Println("upload: ", err)
		jsonhttp.BadRequest(w, "parameter \"files\" missing")
		return
	}

	// upload files one by one
	var references []Reference
	for _, file := range files {
		fd, err := file.Open()
		defer func() {
			err := fd.Close()
			if err != nil {
				fmt.Println("upload: error closing file: ", err)
			}
		}()
		if err != nil {
			fmt.Println("upload: ", err)
			references = append(references, Reference{FileName: file.Filename, Error: err.Error()})
			continue
		}

		//upload file to bee
		reference, err := h.dfsAPI.UploadFile(user, pod, file.Filename, file.Size, fd, podDir, blockSize)
		if err != nil {
			fmt.Println("upload: ", err)
			references = append(references, Reference{FileName: file.Filename, Error: err.Error()})
			continue
		}
		references = append(references, Reference{FileName: file.Filename, Reference: reference})
	}

	jsonhttp.OK(w, &uploadFiletResponse{
		References: references,
	})
}
