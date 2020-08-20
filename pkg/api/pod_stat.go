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

type PodStatResponse struct {
	Version          string `json:"version"`
	PodName          string `json:"name"`
	PodPath          string `json:"path"`
	CreationTime     string `json:"cTime"`
	AccessTime       string `json:"aTime"`
	ModificationTime string `json:"mTime"`
}

func (h *Handler) PodStatHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pod := r.FormValue("pod")
	if user == "" {
		jsonhttp.BadRequest(w, "stat pod: \"user\" argument missing")
		return
	}
	if pod == "" {
		jsonhttp.BadRequest(w, "stat pod: \"pod\" argument missing")
		return
	}

	// fetch pod stat
	stat, err := h.dfsAPI.PodStat(user, pod)
	if err != nil {
		fmt.Println("stat pod: %w", err)
		jsonhttp.InternalServerError(w, err)
	}

	jsonhttp.OK(w, &PodStatResponse{
		Version:          stat.Version,
		PodName:          stat.PodName,
		PodPath:          stat.PodPath,
		CreationTime:     stat.CreationTime,
		AccessTime:       stat.AccessTime,
		ModificationTime: stat.ModificationTime,
	})
}
