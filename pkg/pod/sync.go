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
	"strings"
	"sync"

	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	d "github.com/jmozah/intOS-dfs/pkg/dir"
	m "github.com/jmozah/intOS-dfs/pkg/meta"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) SyncPod(podName string) error {
	podName, err := CleanName(podName)
	if err != nil {
		return fmt.Errorf("sync pod: %w", err)
	}

	if !p.isLoggedInToPod(podName) {
		return fmt.Errorf("sync pod: login to pod to do this operation")
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("sync pod: %w", err)
	}

	err = podInfo.SyncPod(podName, p.client)
	if err != nil {
		return fmt.Errorf("sync pod: %w", err)
	}
	return nil
}

func (pi *Info) SyncPod(podName string, client blockstore.Client) error {
	fd := pi.getFeed()
	accountInfo := pi.getAccountInfo()

	fmt.Println("Syncing pod", podName)
	var wg sync.WaitGroup
	for _, ref := range pi.currentPodInode.Hashes {
		wg.Add(1)
		go func(reference []byte) {
			defer wg.Done()
			_, data, err := fd.GetFeedData(reference, accountInfo.GetAddress())
			if err != nil {
				data, respCode, err := client.DownloadBlob(reference)
				if err != nil {
					fmt.Println("sync: download error: ", err)
					return
				}
				if respCode != http.StatusOK {
					fmt.Println("sync: download status not okay: ", respCode)
					return
				}
				var meta *m.FileMetaData
				err = json.Unmarshal(data, &meta)
				if err != nil {
					fmt.Println("sync: unmarshall error: ", err)
				}

				path := meta.Path + utils.PathSeperator + meta.Name
				pi.file.AddToFileMap(path, meta)
				path = strings.TrimPrefix(path, podName)
				fmt.Println(path)
				return
			}

			var dirInode *d.DirInode
			err = json.Unmarshal(data, &dirInode)
			if err != nil {
				fmt.Println("sync: unmarshall error: %w", err)
				return
			}

			path := dirInode.Meta.Path + utils.PathSeperator + dirInode.Meta.Name
			err = pi.getDirectory().LoadDirMeta(podName, dirInode, fd, accountInfo)
			if err != nil {
				fmt.Println("sync: load meta error: %w", err)
				return
			}
			pi.getDirectory().AddToDirectoryMap(path, dirInode)
			path = strings.TrimPrefix(path, podName)
			fmt.Println(path)
		}(ref)
	}
	wg.Wait()
	return nil
}
