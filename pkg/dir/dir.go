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

package datapod

import (
	"strings"
	"sync"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	f "github.com/jmozah/intOS-dfs/pkg/file"
	m "github.com/jmozah/intOS-dfs/pkg/meta"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

const (
	DirectoryNameLength = 25
)

type Directory struct {
	podName string
	client  blockstore.Client
	fd      *feed.API
	acc     *account.AccountInfo
	file    *f.File
	dirMap  map[string]*DirInode // path to dirInode cache
	dirMu   *sync.RWMutex
}

type DirInode struct {
	Meta   *m.DirectoryMetaData
	Hashes [][]byte
}

func NewDirectory(podName string, client blockstore.Client, fd *feed.API, acc *account.AccountInfo, file *f.File) *Directory {
	return &Directory{
		podName: podName,
		client:  client,
		fd:      fd,
		acc:     acc,
		file:    file,
		dirMap:  make(map[string]*DirInode),
		dirMu:   &sync.RWMutex{},
	}
}

func (d *Directory) getFeed() *feed.API {
	return d.fd
}

func (d *Directory) getAccount() *account.AccountInfo {
	return d.acc
}

func (d *Directory) AddToDirectoryMap(path string, dirInode *DirInode) {
	d.dirMu.Lock()
	defer d.dirMu.Unlock()
	if !strings.HasPrefix(path, "/") {
		path = utils.PathSeperator + path
	}
	d.dirMap[path] = dirInode
}

func (d *Directory) RemoveFromDirectoryMap(path string) {
	d.dirMu.Lock()
	defer d.dirMu.Unlock()
	if !strings.HasPrefix(path, "/") {
		path = utils.PathSeperator + path
	}
	delete(d.dirMap, path)
}

func (d *Directory) GetDirFromDirectoryMap(path string) *DirInode {
	d.dirMu.Lock()
	defer d.dirMu.Unlock()
	if !strings.HasPrefix(path, "/") {
		path = utils.PathSeperator + path
	}
	for k := range d.dirMap {
		if k == path {
			return d.dirMap[path]
		}
	}
	return nil
}

func (d *Directory) GetPrefixPodFromPathMap(prefix string) *DirInode {
	d.dirMu.Lock()
	defer d.dirMu.Unlock()
	for k := range d.dirMap {
		if strings.HasPrefix(k, prefix) {
			delete(d.dirMap, k)
		}
	}
	return nil
}
