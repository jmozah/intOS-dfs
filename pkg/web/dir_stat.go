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

package web

import (
	"net/http"

	"resenje.org/jsonhttp"
)

type DirectoryStatResponse struct {
	Version          string `json:"version"`
	PodName          string `json:"podName"`
	ParentPath       string `json:"path"`
	Name             string `json:"directoryName"`
	CreationTime     string `json:"cTime"`
	AccessTime       string `json:"aTime"`
	ModificationTime string `json:"mTime"`
}

func (h *Handler) DirectoryStatHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	dir := r.FormValue("dir")
	if user == "" {
		jsonhttp.BadRequest(w, "argument missing: user ")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "argument missing: pod")
		return
	}
	if dir == "" {
		jsonhttp.BadRequest(w, "argument missing: dir")
		return
	}

	// TODO: stat directory

	jsonhttp.OK(w, &DirectoryStatResponse{
		Version:          "1",
		PodName:          pod,
		ParentPath:       dir,
		Name:             "dir1",
		CreationTime:     "2006-01-02 15:04:05.999999999 +05:30 UTC",
		AccessTime:       "2006-01-02 15:04:05.999999999 +05:30 UTC",
		ModificationTime: "2006-01-02 15:04:05.999999999 +05:30 UTC",
	})
}
