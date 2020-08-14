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
	"bufio"
	"bytes"
	"fmt"
	"github.com/jmozah/intOS-dfs/pkg/account"
	"io"
	"strconv"
	"strings"
	"sync"

	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	f "github.com/jmozah/intOS-dfs/pkg/file"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) CreatePod(podName, dataDir string, passPhrase string) (*Info, error) {
	podName, err := CleanName(podName)
	if err != nil {
		return nil, fmt.Errorf("create pod: %w", err)
	}

	// check if pods is present and get free index
	pods, err := p.loadUserPods()
	if err != nil {
		return nil, fmt.Errorf("create pod: %w", err)
	}
	if p.checkIfPodPresent(pods, podName) {
		return nil, fmt.Errorf("create pod: pod already exist")
	}
	freeId, err := p.getFreeId(pods)
	if err != nil {
		return nil, err
	}

	// create a child account for the user and other data structures for the pod
	err = p.acc.CreatePodAccount(freeId, passPhrase)
	if err != nil {
		return nil, err
	}
	accountInfo := p.acc.GetAccountInfo(freeId)
	fd := feed.New(accountInfo, p.client)
	file := f.NewFile(podName, p.client, fd, accountInfo)
	dir := d.NewDirectory(podName, p.client, fd, accountInfo, file)

	// create the pod inode
	dirInode, _, err := dir.CreatePodINode(podName)
	if err != nil {
		return nil, fmt.Errorf("create pod: %w", err)
	}

	// store the pod file
	pods[freeId] = podName
	err = p.storeUserPods(pods)
	if err != nil {
		return nil, fmt.Errorf("create pod: %w", err)
	}

	// create the pod info and store it in the podMap
	podInfo := &Info{
		podName:         podName,
		dir:             dir,
		file:            file,
		accountInfo:     accountInfo,
		feed:            fd,
		currentPodInode: dirInode,
		curPodMu:        sync.RWMutex{},
		currentDirInode: dirInode,
		curDirMu:        sync.RWMutex{},
	}
	pods[freeId] = podName
	p.addPodToPodMap(podName, podInfo)

	return podInfo, nil
}

func (p *Pod) loadUserPods() (map[int]string, error) {
	// The user pod file topic should be in the name of the user account
	topic := utils.HashString(p.acc.GetAddress(account.UserAccountIndex).Hex())
	_, data, err := p.fd.GetFeedData(topic, p.acc.GetAddress(account.UserAccountIndex))
	if err != nil {
		if err.Error() != "no feed updates found" {
			return nil, fmt.Errorf("loading pods: %w", err)
		}
	}

	buf := bytes.NewBuffer(data)
	rd := bufio.NewReader(buf)
	pods := make(map[int]string)
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("loading pods: %w", err)
		}
		line = strings.Trim(line, "\n")
		lines := strings.Split(line, ",")
		index, err := strconv.ParseInt(lines[1], 10, 64)
		if err != nil {
			return pods, err
		}
		pods[int(index)] = lines[0]
	}
	return pods, nil
}

func (p *Pod) storeUserPods(pods map[int]string) error {
	buf := bytes.NewBuffer(nil)
	podLen := len(pods)
	for index, pod := range pods {
		pod := strings.Trim(pod, "\n")
		if podLen > 1 && pod == "" {
			continue
		}
		line := fmt.Sprintf("%s,%d", pod, index)
		buf.WriteString(line + "\n")
	}

	topic := utils.HashString(p.acc.GetAddress(account.UserAccountIndex).Hex())
	_, err := p.fd.UpdateFeed(topic, p.acc.GetAddress(account.UserAccountIndex), buf.Bytes())
	if err != nil {
		return fmt.Errorf("store pods: %w", err)
	}
	return nil
}

func (p *Pod) getFreeId(pods map[int]string) (int, error) {
	for i := 0; i < maxPodId; i++ {
		if _, ok := pods[i]; !ok {
			return i, nil
		}
	}
	return 0, fmt.Errorf("max pods exhausted")
}

func (p *Pod) checkIfPodPresent(pods map[int]string, podName string) bool {
	for _, pod := range pods {
		if strings.Trim(pod, "\n") == podName {
			return true
		}
	}
	return false
}
