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
	"github.com/jmozah/intOS-dfs/pkg/dfs"
)

type uploadFileResponse struct {
	References []Reference
}

type Reference struct {
	FileName  string `json:"file_name"`
	Reference string `json:"reference,omitempty"`
	Error     string `json:"error,omitempty"`
}

const (
	defaultMaxMemory  = 32 << 20 // 32 MB
	compressionHeader = "intOS-dfs-Compression"
)

func (h *Handler) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	podDir := r.FormValue("pod_dir")
	blockSize := r.FormValue("block_size")
	compression := r.Header.Get(compressionHeader)
	if podDir == "" {
		h.logger.Errorf("upload: \"pod_dir\" argument missing")
		jsonhttp.BadRequest(w, "upload: \"pod_dir\" argument missing")
		return
	}
	if blockSize == "" {
		h.logger.Errorf("upload: \"block_size\" argument missing")
		jsonhttp.BadRequest(w, "upload: \"block_size\" argument missing")
		return
	}

	if compression != "" {
		if compression != "snappy" && compression != "gzip" {
			h.logger.Errorf("upload: invalid value for \"compression\" header")
			jsonhttp.BadRequest(w, "upload: invalid value for \"compression\" header")
			return
		}
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		h.logger.Errorf("upload: invalid cookie: %v", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		h.logger.Errorf("upload: \"cookie-id\" parameter missing in cookie")
		jsonhttp.BadRequest(w, "upload: \"cookie-id\" parameter missing in cookie")
		return
	}

	//  get the files parameter from the multi part
	err = r.ParseMultipartForm(defaultMaxMemory)
	if err != nil {
		h.logger.Errorf("upload: %v", err)
		jsonhttp.BadRequest(w, "upload: "+err.Error())
		return
	}
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		h.logger.Errorf("upload: parameter \"files\" missing")
		jsonhttp.BadRequest(w, "upload: parameter \"files\" missing")
		return
	}

	// upload files one by one
	var references []Reference
	for _, file := range files {
		fd, err := file.Open()
		defer func() {
			err := fd.Close()
			if err != nil {
				h.logger.Errorf("upload: error closing file: %v", err)
			}
		}()
		if err != nil {
			h.logger.Errorf("upload: %v", err)
			references = append(references, Reference{FileName: file.Filename, Error: err.Error()})
			continue
		}

		//upload file to bee
		reference, err := h.dfsAPI.UploadFile(file.Filename, sessionId, file.Size, fd, podDir, blockSize, compression)
		if err != nil {
			if err == dfs.ErrPodNotOpen {
				h.logger.Errorf("upload: %v", err)
				jsonhttp.BadRequest(w, "upload: "+err.Error())
				return
			}
			h.logger.Errorf("upload: %v", err)
			references = append(references, Reference{FileName: file.Filename, Error: err.Error()})
			continue
		}
		references = append(references, Reference{FileName: file.Filename, Reference: reference})
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, &uploadFileResponse{
		References: references,
	})
}
