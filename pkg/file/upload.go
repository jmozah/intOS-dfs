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
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/golang/snappy"
	m "github.com/jmozah/intOS-dfs/pkg/meta"
	"github.com/klauspost/pgzip"
)

var (
	NoOfParallelWorkers = runtime.NumCPU() * 4
)

func (f *File) Upload(fd io.Reader, fileName string, fileSize int64, blockSize uint32, filePath, compression string) ([]byte, error) {
	reader := bufio.NewReader(fd)
	now := time.Now().Unix()
	meta := m.FileMetaData{
		Version:          m.FileMetaVersion,
		Path:             filepath.Dir(filePath),
		Name:             fileName,
		FileSize:         uint64(fileSize),
		BlockSize:        blockSize,
		Compression:      compression,
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
					return nil, fmt.Errorf("invalid file length of file data received")
				}
				break
			} else {
				return nil, err
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
			}()
			// compress the data
			uploadData := data[:size]
			if compression != "" {
				uploadData, err = compress(data[:size], compression, blockSize)
				if err != nil {
					errC <- err
				}
			}

			addr, err := f.client.UploadBlob(uploadData, true)
			if err != nil {
				errC <- err
				return
			}

			fileBlock := &FileBlock{
				Name:           fmt.Sprintf("block-%05d", counter),
				Size:           uint32(size),
				CompressedSize: uint32(len(uploadData)),
				Address:        addr,
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
		return nil, err
	}

	// copy the block references to the fileInode
	for i := 0; i < len(refMap); i++ {
		fileINode.FileBlocks = append(fileINode.FileBlocks, refMap[i])
	}

	fileInodeData, err := json.Marshal(fileINode)
	if err != nil {
		return nil, err
	}

	addr, err := f.client.UploadBlob(fileInodeData, true)
	if err != nil {
		return nil, err
	}

	meta.InodeAddress = addr
	fileMetaBytes, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}
	metaAddr, err := f.client.UploadBlob(fileMetaBytes, true)
	if err != nil {
		return nil, err
	}
	meta.MetaReference = metaAddr // the self address is stored to share this file easily
	f.AddToFileMap(filePath, &meta)
	return metaAddr, nil
}

func compress(dataToCompress []byte, compression string, blockSize uint32) ([]byte, error) {
	switch compression {
	case "gzip":
		var b bytes.Buffer
		w := pgzip.NewWriter(&b)
		block := int(blockSize / 10)
		err := w.SetConcurrency(block, 10)
		if err != nil {
			return nil, err
		}
		_, err = w.Write(dataToCompress)
		if err != nil {
			return nil, err
		}
		err = w.Close()
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case "snappy":
		return snappy.Encode(nil, dataToCompress), nil
	}
	return dataToCompress, nil
}
