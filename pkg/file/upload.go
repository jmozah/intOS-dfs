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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	m "github.com/jmozah/intOS-dfs/pkg/meta"
)

var (
	NoOfParallelWorkers = runtime.NumCPU() * 4
)

func (f *File) Upload(fd multipart.File, fileName string, fileSize int64, blockSize uint32, filePath string) ([]byte, error) {
	reader := bufio.NewReader(fd)
	now := time.Now().Unix()
	meta := m.FileMetaData{
		Version:          m.FileMetaVersion,
		Path:             filepath.Dir(filePath),
		Name:             fileName,
		FileSize:         uint64(fileSize),
		BlockSize:        blockSize,
		CreationTime:     now,
		AccessTime:       now,
		ModificationTime: now,
	}

	fileINode := FileINode{}

	var totalLength uint64
	i := 0
	errC := make(chan error)
	doneC := make(chan bool)
	worker := make(chan bool, NoOfParallelWorkers)
	var wg sync.WaitGroup
	refMap := make(map[int]*FileBlock)
	refMapMu := sync.RWMutex{}
	for {
		data := make([]byte, blockSize)
		r, err := reader.Read(data)
		totalLength += uint64(r)
		if err != nil {
			if err == io.EOF {
				if totalLength < uint64(fileSize) {
					return nil, fmt.Errorf("uplaod: invalid file length of file data received")
				}
				break
			} else {
				return nil, fmt.Errorf("uplaod: %w", err)
			}
		}

		// determine the content type from the first 512 bytes of the file
		if i == 0 {
			cBytes := bytes.NewReader(data[:512])
			cReader := bufio.NewReader(cBytes)
			meta.ContentType = f.GetContentType(cReader)
		}

		wg.Add(1)
		worker <- true
		go func(counter, size int) {
			defer func() {
				<-worker
				wg.Done()
				fmt.Println("uploaded chunk: ", counter, size)
			}()
			addr, err := f.client.UploadBlob(data[:size])
			if err != nil {
				errC <- err
				return
			}

			fileBlock := &FileBlock{
				Name:    fmt.Sprintf("block-%05d", counter),
				Size:    uint32(size),
				Address: addr,
			}

			refMapMu.Lock()
			defer refMapMu.Unlock()
			refMap[counter] = fileBlock
		}(i, r)

		i++
	}

	go func() {
		wg.Wait()
		close(doneC)
	}()

	select {
	case <-doneC:
		break
	case err := <-errC:
		close(errC)
		return nil, fmt.Errorf("uplaod: %w", err)
	}

	// copy the block references to the fileInode
	for i := 0; i < len(refMap); i++ {
		fileINode.FileBlocks = append(fileINode.FileBlocks, refMap[i])
	}

	fileInodeData, err := json.Marshal(fileINode)
	if err != nil {
		return nil, fmt.Errorf("uplaod: %v", fileName)
	}

	addr, err := f.client.UploadBlob(fileInodeData)
	if err != nil {
		return nil, fmt.Errorf("uplaod: %w", err)
	}

	meta.InodeAddress = addr
	fileMetaBytes, err := json.Marshal(meta)
	if err != nil {
		return nil, fmt.Errorf("uplaod: %v", fileName)
	}
	metaAddr, err := f.client.UploadBlob(fileMetaBytes)
	if err != nil {
		return nil, fmt.Errorf("uplaod: %w", err)
	}

	f.AddToFileMap(filePath, &meta)
	return metaAddr, nil
}
