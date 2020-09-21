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
	"strings"

	"github.com/jmozah/intOS-dfs/pkg/account"
	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) isPodOpened(podName string) bool {
	p.podMu.Lock()
	defer p.podMu.Unlock()
	name1 := utils.PathSeperator + podName
	if podInfo, ok := p.podMap[name1]; ok {
		if podInfo.currentPodInode != nil {
			return true
		}
	}
	return false
}

func (p *Pod) GetPath(inode *d.DirInode) string {
	if inode != nil {
		return inode.Meta.Path
	}
	return ""
}

func (p *Pod) GetName(inode *d.DirInode) string {
	if inode != nil {
		return inode.Meta.Name
	}
	return ""
}

func (p *Pod) GetAccountInfo(podName string) (*account.AccountInfo, error) {
	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return nil, err
	}
	return podInfo.getAccountInfo(), nil
}

func CleanPodName(podName string) (string, error) {
	if podName == "" {
		return "", ErrInvalidPodName
	}
	if len(podName) > utils.MaxPodNameLength {
		return "", ErrTooLongPodName
	}
	podName = strings.TrimSpace(podName)
	podName = strings.Trim(podName, "\\/,\t ")
	return podName, nil
}

func CleanDirName(dirName string) ([]string, error) {
	if dirName == "" {
		return nil, ErrInvalidDirectory
	}

	var cleanedDirs []string
	if dirName == utils.PathSeperator {
		cleanedDirs = append(cleanedDirs, dirName)
		return cleanedDirs, nil
	}

	dirs := strings.Split(dirName, utils.PathSeperator)

	for _, dir := range dirs {
		if len(dir) > utils.MaxDirectoryNameLength {
			return nil, ErrTooLongDirectoryName
		}
		dir = strings.TrimSpace(dir)
		dir = strings.Trim(dir, "\\/,\t ")
		if dir != "" {
			cleanedDirs = append(cleanedDirs, dir)
		}
	}
	return cleanedDirs, nil
}
