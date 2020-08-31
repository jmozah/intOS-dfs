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

package user

import (
	"sync"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	f "github.com/jmozah/intOS-dfs/pkg/file"
	"github.com/jmozah/intOS-dfs/pkg/pod"
)

type Info struct {
	name      string
	podName   string
	sessionId string
	feedApi   *feed.API
	account   *account.Account
	file      *f.File
	dir       *d.Directory
	pods      *pod.Pod
}

type Users struct {
	dataDir string
	client  blockstore.Client
	userMap map[string]*Info
	userMu  *sync.RWMutex
}

func NewUsers(dataDir string, client blockstore.Client) *Users {
	return &Users{
		dataDir: dataDir,
		client:  client,
		userMap: make(map[string]*Info),
		userMu:  &sync.RWMutex{},
	}
}

func (u *Users) addUserToMap(info *Info) {
	u.userMu.Lock()
	defer u.userMu.Unlock()
	u.userMap[info.sessionId] = info
}

func (u *Users) removeUserFromMap(sessionId string) {
	u.userMu.Lock()
	defer u.userMu.Unlock()
	delete(u.userMap, sessionId)
}

func (u *Users) getUserFromMap(sessionId string) *Info {
	u.userMu.Lock()
	defer u.userMu.Unlock()
	return u.userMap[sessionId]
}

func (u *Users) isUserPresentInMap(sessionId string) bool {
	u.userMu.Lock()
	defer u.userMu.Unlock()
	if _, ok := u.userMap[sessionId]; ok {
		return true
	}
	return false
}
