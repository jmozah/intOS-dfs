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
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/ethersphere/bee/pkg/swarm"
)

func (f *File) Download(podFile string) (io.ReadCloser, string, string, error) {
	//TODO: need to change the access time for podFile

	meta := f.GetFromFileMap(podFile)
	if meta == nil {
		return nil, "", "", fmt.Errorf("file not found in dfs")
	}

	fileInodeBytes, _, err := f.getClient().DownloadBlob(meta.InodeAddress)
	if err != nil {
		return nil, "", "", err
	}
	var fileInode FileINode
	err = json.Unmarshal(fileInodeBytes, &fileInode)
	if err != nil {
		return nil, "", "", err
	}

	reader := NewReader(fileInode, f.getClient(), meta.FileSize, meta.BlockSize, meta.Compression)
	ref := swarm.NewAddress(meta.InodeAddress).String()
	size := strconv.FormatUint(meta.FileSize, 10)
	return reader, ref, size, nil
}
