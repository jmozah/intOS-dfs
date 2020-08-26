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

	"github.com/jmozah/intOS-dfs/pkg/cookie"
	"github.com/jmozah/intOS-dfs/pkg/dfs"
	p "github.com/jmozah/intOS-dfs/pkg/pod"
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
	pod := r.FormValue("pod")
	if pod == "" {
		jsonhttp.BadRequest(w, "stat podd: \"pod\" argument missing")
		return
	}
	// get values from cookie
	userName, sessionId, _, err := cookie.GetUserNameSessionIdAndPodName(r)
	if err != nil {
		fmt.Println("stat pod: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if userName == "" {
		jsonhttp.BadRequest(w, "stat pod: \"user\" parameter missing in cookie")
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "stat pod: \"cookie-id\" parameter missing in cookie")
		return
	}

	// restart the cookie expiry
	err = cookie.ResetSessionExpiry(r, w)
	if err != nil {
		jsonhttp.BadRequest(w, err)
		return
	}

	// fetch pod stat
	stat, err := h.dfsAPI.PodStat(userName, pod, sessionId)
	if err != nil {
		if err == dfs.ErrInvalidUserName || err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrInvalidPodName {
			fmt.Println("stat pod: ", err)
			jsonhttp.BadRequest(w, err)
			return
		}
		fmt.Println("stat pod: ", err)
		jsonhttp.InternalServerError(w, err)
		return
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
