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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	m "github.com/jmozah/intOS-dfs/pkg/meta"
)

const (
	MaxRoutinesPerUpload = 128
)

func (f *File) CopyFromFile(podName, localFileName string, fileInfo os.FileInfo, blockSize uint32, filePath string) ([]byte, error) {
	now := time.Now().Unix()
	meta := m.FileMetaData{
		Version:          m.FileMetaVersion,
		Path:             filepath.Dir(filePath),
		Name:             fileInfo.Name(),
		FileSize:         uint64(fileInfo.Size()),
		BlockSize:        blockSize,
		CreationTime:     now,
		AccessTime:       now,
		ModificationTime: now,
	}

	fileINode := FileINode{}

	fl, err := os.Open(localFileName)
	if err != nil {
		return nil, fmt.Errorf("copyFromLocal: %w", err)
	}

	data := make([]byte, blockSize)
	var totalLength uint64
	i := 0
	for {
		r, err := fl.Read(data)
		totalLength += uint64(r)
		if err != nil {
			if err == io.EOF {
				if totalLength < uint64(fileInfo.Size()) {
					return nil, fmt.Errorf("copyFromLocal: invalid file length of file data received")
				}
				break
			} else {
				return nil, fmt.Errorf("copyFromLocal: %w", err)
			}
		}
		fmt.Printf("uploading block-%05d, ", i)

		addr, err := f.client.UploadBlob(data[:r])
		if err != nil {
			return nil, fmt.Errorf("copyFromLocal: %w", err)
		}

		fileBlock := &FileBlock{
			Name:    fmt.Sprintf("block-%05d", i),
			Size:    uint32(r),
			Address: addr,
		}

		fileINode.FileBlocks = append(fileINode.FileBlocks, fileBlock)
		fmt.Println(hex.EncodeToString(addr))
		i++
	}

	fileInodeData, err := json.Marshal(fileINode)
	if err != nil {
		return nil, fmt.Errorf("copyFromLocal: %v", fileInfo.Name())
	}

	addr, err := f.client.UploadBlob(fileInodeData)
	if err != nil {
		return nil, fmt.Errorf("copyFromLocal: %w", err)
	}

	meta.InodeAddress = addr
	fileMetaBytes, err := json.Marshal(meta)
	if err != nil {
		return nil, fmt.Errorf("copyFromLocal: %v", fileInfo.Name())
	}
	metaAddr, err := f.client.UploadBlob(fileMetaBytes)
	if err != nil {
		return nil, fmt.Errorf("copyFromLocal: %w", err)
	}

	f.AddToFileMap(filePath, &meta)
	return metaAddr, nil
}
