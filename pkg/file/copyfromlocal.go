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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	m "github.com/jmozah/intOS-dfs/pkg/meta"
)

const (
	MaxRoutinesPerUpload = 128
)

func (f *File) CopyFromFile(podName, localFileName string, fileInfo os.FileInfo, blockSize uint32, filePath string) ([]byte, error) {
	fl, err := os.Open(localFileName)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	reader := bufio.NewReader(fl)
	now := time.Now().Unix()
	meta := m.FileMetaData{
		Version:          m.FileMetaVersion,
		Path:             filepath.Dir(filePath),
		Name:             fileInfo.Name(),
		FileSize:         uint64(fileInfo.Size()),
		BlockSize:        blockSize,
		ContentType:      f.GetContentType(reader),
		CreationTime:     now,
		AccessTime:       now,
		ModificationTime: now,
	}

	fileINode := FileINode{}
	data := make([]byte, blockSize)
	var totalLength uint64
	i := 0
	for {
		r, err := reader.Read(data)
		totalLength += uint64(r)
		if err != nil {
			if err == io.EOF {
				if totalLength < uint64(fileInfo.Size()) {
					return nil, fmt.Errorf("invalid file length of file data received")
				}
				break
			} else {
				return nil, err
			}
		}
		fmt.Printf("uploading block-%05d, ", i)

		addr, err := f.client.UploadBlob(data[:r], true)
		if err != nil {
			return nil, err
		}

		fileBlock := &FileBlock{
			Name:    fmt.Sprintf("block-%05d", i),
			Size:    uint32(r),
			Address: addr,
		}

		fileINode.FileBlocks = append(fileINode.FileBlocks, fileBlock)
		f.logger.Infof(hex.EncodeToString(addr))
		i++
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

	meta.MetaReference = metaAddr // to get the address for sharing
	f.AddToFileMap(filePath, &meta)
	return metaAddr, nil
}

func (f *File) GetContentType(bufferReader *bufio.Reader) string {
	buffer, err := bufferReader.Peek(512)
	if err != nil && err != io.EOF {
		return ""
	}
	return http.DetectContentType(buffer)
}
