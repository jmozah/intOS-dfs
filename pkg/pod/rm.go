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
	"encoding/json"
	"fmt"
	"net/http"
	gopath "path"
	"time"

	"github.com/ethersphere/bee/pkg/swarm"

	m "github.com/jmozah/intOS-dfs/pkg/meta"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) RemoveFile(podName string, podFile string) error {
	if !p.isLoggedInToPod(podName) {
		return fmt.Errorf("rm: login to pod to do this operation")
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("rm: %w", err)
	}
	dir := podInfo.getDirectory()

	var path string
	if podInfo.IsCurrentDirRoot() {
		path = podInfo.GetCurrentPodPathAndName() + podFile
	} else {
		path = podInfo.GetCurrentDirPathAndName() + utils.PathSeperator + podFile
	}

	if !podInfo.getFile().IsFileAlreadyPResent(path) {
		return fmt.Errorf("rm: file not present in pod")
	}

	_, dirInode, err := dir.GetDirNode(gopath.Dir(path), podInfo.getFeed(), podInfo.getAccountInfo())
	if err != nil {
		return fmt.Errorf("error while fetching pod info: %w", err)
	}

	// remove the file
	var newHashes [][]byte
	for _, hash := range dirInode.Hashes {
		_, _, err := podInfo.getFeed().GetFeedData(hash, podInfo.getAccountInfo().GetAddress())
		if err != nil {
			data, respCode, err := p.GetClient().DownloadBlob(hash)
			if err != nil || respCode != http.StatusOK {
				fmt.Println("could not load address ", swarm.NewAddress(hash).String())
				continue
			}
			var meta *m.FileMetaData
			err = json.Unmarshal(data, &meta)
			if err != nil {
				fmt.Println("could not unmarshall data in address ", swarm.NewAddress(hash).String())
				continue
			}
			if meta.Name != gopath.Base(path) {
				newHashes = append(newHashes, hash)
			} else {
				podInfo.getFile().RemoveFromFileMap(path)
			}
		}
	}
	dirInode.Hashes = newHashes

	dirInode.Meta.ModificationTime = time.Now().Unix()
	topic, err := dir.UpdateDirectory(dirInode)
	if err != nil {
		return fmt.Errorf("rm: error updating directory: %w", err)
	}

	if path != podInfo.GetCurrentPodPathAndName() {
		err = p.UpdateTillThePod(podName, podInfo.getDirectory(), topic, true)
		if err != nil {
			return fmt.Errorf("rm: error updating directory: %w", err)
		}
	}
	return nil

}
