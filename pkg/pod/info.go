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

package pod

import (
	"sync"

	"github.com/jmozah/intOS-dfs/pkg/account"
	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	f "github.com/jmozah/intOS-dfs/pkg/file"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

type Info struct {
	podName         string
	dir             *d.Directory
	file            *f.File
	accountInfo     *account.AccountInfo
	feed            *feed.API
	currentPodInode *d.DirInode
	curPodMu        sync.RWMutex
	currentDirInode *d.DirInode
	curDirMu        sync.RWMutex
}

func (i *Info) getDirectory() *d.Directory {
	return i.dir
}

func (i *Info) getFile() *f.File {
	return i.file
}

func (i *Info) getAccountInfo() *account.AccountInfo {
	return i.accountInfo
}

func (i *Info) getFeed() *feed.API {
	return i.feed
}

func (i *Info) GetCurrentPodInode() *d.DirInode {
	return i.currentPodInode
}
func (i *Info) GetCurrentDirInode() *d.DirInode {
	return i.currentDirInode
}

func (i *Info) SetCurrentPodInode(podInode *d.DirInode) {
	i.currentPodInode = podInode
}
func (i *Info) SetCurrentDirInode(podInode *d.DirInode) {
	i.currentDirInode = podInode
}

func (p *Info) IsCurrentDirRoot() bool {
	if p.currentDirInode.Meta.Path == utils.PathSeperator {
		return true
	} else {
		return false
	}
}

func (i *Info) GetCurrentPodPathOnly() string {
	return i.currentPodInode.Meta.Path
}

func (i *Info) GetCurrentPodNameOnly() string {
	return i.currentPodInode.Meta.Name
}

func (i *Info) GetCurrentPodPathAndName() string {
	return i.currentPodInode.Meta.Path + i.currentPodInode.Meta.Name
}

func (i *Info) GetCurrentDirPathOnly() string {
	return i.currentDirInode.Meta.Path
}

func (i *Info) GetCurrentDirNameOnly() string {
	return i.currentDirInode.Meta.Name
}

func (i *Info) GetCurrentDirPathAndName() string {
	return i.currentDirInode.Meta.Path + utils.PathSeperator + i.currentDirInode.Meta.Name
}
