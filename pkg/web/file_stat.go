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

type FileStatResponse struct {
	Blocks []BlockInfo
}

type BlockInfo struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Size      string `json:"size"`
}

func (h *Handler) FileStatHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	podFile := r.FormValue("file")
	if user == "" {
		jsonhttp.BadRequest(w, "argument missing: user ")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "argument missing: pod")
		return
	}
	if podFile == "" {
		jsonhttp.BadRequest(w, "argument missing: file")
		return
	}

	// TODO: get file stat
	var blocks []BlockInfo
	block1 := BlockInfo{
		Name:      "block-00000",
		Reference: mockAddress1,
		Size:      "100",
	}
	blocks = append(blocks, block1)
	block2 := BlockInfo{
		Name:      "block-00001",
		Reference: mockAddress2,
		Size:      "100",
	}
	blocks = append(blocks, block2)
	block3 := BlockInfo{
		Name:      "block-00002",
		Reference: mockAddress3,
		Size:      "77",
	}
	blocks = append(blocks, block3)

	jsonhttp.OK(w, &FileStatResponse{
		Blocks: blocks,
	})
}
