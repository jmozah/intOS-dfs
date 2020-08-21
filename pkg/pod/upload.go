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
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) UploadFile(podName string, fileName string, fileSize int64, reader io.ReadCloser, podDir string, blockSize string) (string, error) {
	if !p.isLoggedInToPod(podName) {
		return "", fmt.Errorf("upload: login to pod to do this operation")
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}
	dir := podInfo.getDirectory()

	bs, err := humanize.ParseBytes(blockSize)
	if err != nil {
		return "", fmt.Errorf("upload: block size parse error: %w", err)
	}

	path := p.getFilePath(podDir, podInfo)

	_, dirInode, err := dir.GetDirNode(path, podInfo.getFeed(), podInfo.getAccountInfo())
	if err != nil {
		return "", fmt.Errorf("error while fetching pod info: %w", err)
	}

	fpath := path + utils.PathSeperator + fileName
	if podInfo.file.IsFileAlreadyPResent(fpath) {
		return "", fmt.Errorf("file already present in the destination dir")
	}
	addr, err := podInfo.file.Upload(reader, fileName, fileSize, uint32(bs), fpath)
	if err != nil {
		return "", fmt.Errorf("error while copying file to pod: %w", err)
	}
	dirInode.Hashes = append(dirInode.Hashes, addr)
	return swarm.NewAddress(addr).String(), nil
}
