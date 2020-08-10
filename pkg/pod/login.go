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
	"github.com/jmozah/intOS-dfs/pkg/account"
	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	f "github.com/jmozah/intOS-dfs/pkg/file"
	"github.com/jmozah/intOS-dfs/pkg/utils"
	"sync"
)

func (p *Pod) LoginPod(podName string, dataDir string, passPhrase string) (*Info, error) {
	podName, err := CleanName(podName)
	if err != nil {
		return nil, fmt.Errorf("login pod: %w", err)
	}

	index, err := p.checkIfPodPresent(podName)
	if err != nil {
		return nil, fmt.Errorf("login pod: %w", err)
	}

	if index < 0 {
		return nil, fmt.Errorf("login pod: pod doesn't exist")
	}

	acc := account.New(podName, dataDir)
	err = acc.CreateNormalAccount(index, passPhrase)
	if err != nil {
		return nil, err
	}
	fd := feed.New(acc, p.client)
	file := f.NewFile(podName, p.client, fd, acc)
	dir := d.NewDirectory(podName, p.client, fd, acc, file)

	_, dirInode, err := dir.GetDirNode(utils.PathSeperator+podName, fd, acc)
	if err != nil {
		return nil, fmt.Errorf("login pod: %w", err)
	}

	podInfo := &Info{
		podName:         podName,
		account:         acc,
		feed:            fd,
		dir:             dir,
		file:            file,
		currentPodInode: dirInode,
		curPodMu:        sync.RWMutex{},
		currentDirInode: dirInode,
		curDirMu:        sync.RWMutex{},
	}

	p.addPodToPodMap(podName, podInfo)
	dir.AddToDirectoryMap(podName, dirInode)

	err = p.SyncPod(podName)
	if err != nil {
		return nil, fmt.Errorf("login pod: %s ", podName)
	}
	return podInfo, nil
}
