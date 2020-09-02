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

	"github.com/jmozah/intOS-dfs/pkg/cookie"
	"github.com/jmozah/intOS-dfs/pkg/dfs"
	"github.com/jmozah/intOS-dfs/pkg/dir"
	p "github.com/jmozah/intOS-dfs/pkg/pod"
	"resenje.org/jsonhttp"
)

type ListFileResponse struct {
	Entries []dir.DirOrFileEntry `json:"entries"`
}

type DirOrFileEntry struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	Size             string `json:"size,omitempty"`
	CreationTime     string `json:"creation_time"`
	ModificationTime string `json:"modification_time"`
	AccessTime       string `json:"access_time"`
}

func (h *Handler) DirectoryLsHandler(w http.ResponseWriter, r *http.Request) {
	directory := r.FormValue("dir")
	if directory == "" {
		jsonhttp.BadRequest(w, "ls dir: \"dir\" argument missing")
		return
	}

	// get values from cookie
	sessionId, err := cookie.GetSessionIdFromCookie(r)
	if err != nil {
		fmt.Println("ls dir: ", err)
		jsonhttp.BadRequest(w, ErrInvalidCookie)
		return
	}
	if sessionId == "" {
		jsonhttp.BadRequest(w, "ls dir: \"cookie-id\" parameter missing in cookie")
		return
	}

	// list directory
	entries, err := h.dfsAPI.ListDir(directory, sessionId)
	if err != nil {
		w.Header().Set("Content-Type", " application/json")
		if err == dfs.ErrPodNotOpen || err == dfs.ErrUserNotLoggedIn ||
			err == p.ErrPodNotOpened {
			fmt.Println("ls dir: ", err)
			jsonhttp.BadRequest(w, &ErrorMessage{Err: "ls dir: " + err.Error()})
			return
		}
		fmt.Println("ls dir: ", err)
		jsonhttp.InternalServerError(w, &ErrorMessage{Err: "ls dir: " + err.Error()})
		return
	}

	if entries == nil {
		entries = make([]dir.DirOrFileEntry, 0)
	}

	w.Header().Set("Content-Type", " application/json")
	jsonhttp.OK(w, &ListFileResponse{
		Entries: entries,
	})
}
