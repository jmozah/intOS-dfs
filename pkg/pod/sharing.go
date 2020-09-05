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
	"path/filepath"
	"time"

	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) GetMetaReferenceOfFile(podName, filePath string) ([]byte, string, error) {
	if !p.isPodOpened(podName) {
		return nil, "", fmt.Errorf("share: login to pod to do this operation")
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return nil, "", fmt.Errorf("share: %w", err)
	}

	podDir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	path := p.getFilePath(podDir, podInfo)
	fpath := path + utils.PathSeperator + fileName

	return podInfo.getFile().GetFileReference(fpath)
}

func (p *Pod) ReceiveFileAndStore(podName, podDir, fileName, metaHexRef string) error {
	if !p.isPodOpened(podName) {
		return fmt.Errorf("receive: login to pod to do this operation")
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}

	path := p.getFilePath(podDir, podInfo)
	dir := podInfo.getDirectory()

	_, dirInode, err := dir.GetDirNode(path, podInfo.getFeed(), podInfo.getAccountInfo())
	if err != nil {
		return fmt.Errorf("receive: error while fetching pod info: %w", err)
	}

	// check if the file exists already
	fpath := path + utils.PathSeperator + fileName
	if podInfo.file.IsFileAlreadyPResent(fpath) {
		return fmt.Errorf("receive: file already present in the destination dir")
	}

	// append the file meta to the parent directory and update the directory feed
	metaReference, err := utils.ParseHexReference(metaHexRef)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}
	dirInode.Hashes = append(dirInode.Hashes, metaReference.Bytes())
	dirInode.Meta.ModificationTime = time.Now().Unix()
	topic, err := dir.UpdateDirectory(dirInode)
	if err != nil {
		return fmt.Errorf("receive: error updating directory: %w", err)
	}

	// if the directory path is not root.. then update all the parents too
	if path != podInfo.GetCurrentPodPathAndName() {
		err = p.UpdateTillThePod(podName, podInfo.getDirectory(), topic, path, true)
		if err != nil {
			return fmt.Errorf("receive: error updating directory: %w", err)
		}
	}

	// Add to file path map
	return podInfo.getFile().AddFileToPath(path, metaHexRef)
}
