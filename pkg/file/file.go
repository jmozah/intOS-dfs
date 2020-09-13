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

package file

import (
	"sync"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	"github.com/jmozah/intOS-dfs/pkg/logging"
	m "github.com/jmozah/intOS-dfs/pkg/meta"
)

type File struct {
	podName string
	client  blockstore.Client
	fd      *feed.API
	acc     *account.AccountInfo
	fileMap map[string]*m.FileMetaData
	fileMu  *sync.RWMutex
	logger  logging.Logger
}

type FileINode struct {
	FileBlocks []*FileBlock
}

type FileBlock struct {
	Name           string
	Size           uint32
	CompressedSize uint32
	Address        []byte
}

func NewFile(podName string, client blockstore.Client, fd *feed.API, acc *account.AccountInfo, logger logging.Logger) *File {
	return &File{
		podName: podName,
		client:  client,
		fd:      fd,
		acc:     acc,
		fileMap: make(map[string]*m.FileMetaData),
		fileMu:  &sync.RWMutex{},
		logger:  logger,
	}
}

func (f *File) getClient() blockstore.Client {
	return f.client
}

func (f *File) AddToFileMap(filePath string, meta *m.FileMetaData) {
	f.fileMu.Lock()
	defer f.fileMu.Unlock()
	f.fileMap[filePath] = meta
}

func (f *File) RemoveFromFileMap(filePath string) {
	f.fileMu.Lock()
	defer f.fileMu.Unlock()
	delete(f.fileMap, filePath)
}

func (f *File) GetFromFileMap(filePath string) *m.FileMetaData {
	f.fileMu.Lock()
	defer f.fileMu.Unlock()
	if meta, ok := f.fileMap[filePath]; ok {
		return meta
	}
	return nil
}

func (f *File) IsFileAlreadyPResent(fileWithPath string) bool {
	f.fileMu.Lock()
	defer f.fileMu.Unlock()
	if _, ok := f.fileMap[fileWithPath]; ok {
		return true
	}
	return false
}
