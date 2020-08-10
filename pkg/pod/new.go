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
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/jmozah/intOS-dfs/pkg/account"
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

	if p.rootDirInode == nil {
		return nil, fmt.Errorf("please format the dfs before using it")
	}

	index, err := p.checkIfPodPresent(podName)
	if err != nil {
		return nil, fmt.Errorf("create pod: %w", err)
	}

	if index != -1 {
		return nil, fmt.Errorf("create pod: pod already exist")
	}

	pods, err := p.getRootFileContents()
	if err != nil {
		return nil, fmt.Errorf("create pod: %w", err)
	}

	acc := account.New(podName, dataDir)
	freeId, err := p.getFreeId(pods)
	if err != nil {
		return nil, err
	}
	err = acc.CreateNormalAccount(freeId, passPhrase)
	if err != nil {
		return nil, err
	}
	fd := feed.New(acc, p.client)
	file := f.NewFile(podName, p.client, fd, acc)
	dir := d.NewDirectory(podName, p.client, fd, acc, file)

	dirInode, _, err := dir.CreateDirINode("", podName, nil)
	if err != nil {
		return nil, fmt.Errorf("create pod: %w", err)
	}

	podInfo := &Info{
		podName:         podName,
		dir:             dir,
		file:            file,
		account:         acc,
		feed:            fd,
		currentPodInode: dirInode,
		curPodMu:        sync.RWMutex{},
		currentDirInode: dirInode,
		curDirMu:        sync.RWMutex{},
	}
	pods[freeId] = podName
	err = p.storeAsRootFile(pods)
	if err != nil {
		return nil, fmt.Errorf("create pod: %w", err)
	}

	p.addPodToPodMap(podName, podInfo)
	return podInfo, nil
}

func (p *Pod) getRootFileContents() (map[int]string, error) {

	if p.rootAccount == nil {
		return nil, fmt.Errorf("loading pods: dfs not initialised")
	}

	topic := utils.HashString(utils.PodsInfoFile)
	_, data, err := p.rootFeed.GetFeedData(topic, p.rootAccount.GetAddress())
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

func (p *Pod) storeAsRootFile(pods map[int]string) error {
	if p.rootAccount == nil {
		return fmt.Errorf("store pods: dfs not initialised")
	}

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

	topic := utils.HashString(utils.PodsInfoFile)
	_, err := p.rootFeed.UpdateFeed(topic, p.rootAccount.GetAddress(), buf.Bytes())
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
