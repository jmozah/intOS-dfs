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
	"strings"

	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) RemoveDir(podName string, dirName string) error {
	directoryName, err := CleanName(dirName)
	if err != nil {
		return err
	}

	if !p.isPodOpened(podName) {
		return fmt.Errorf("rmdir: login to pod to do this operation")
	}

	info, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("rmdir: %w", err)
	}

	directory := info.getDirectory()

	dirInode, err := p.GetInodeFromName(directoryName, info.GetCurrentDirInode(), directory, info)
	if err != nil {
		return fmt.Errorf("rmdir: %w", err)
	}

	if dirInode == nil {
		return fmt.Errorf("rmdir: name is not a directory")
	}

	topic := info.GetCurrentDirPathAndName() + utils.PathSeperator + directoryName
	if info.IsCurrentDirRoot() {
		topic = info.GetCurrentPodPathAndName() + utils.PathSeperator + directoryName
	}
	topicBytes := utils.HashString(topic)
	err = p.UpdateTillThePod(podName, directory, topicBytes, false)
	if err != nil {
		return fmt.Errorf("error updating directory: %w", err)
	}
	directory.GetPrefixPodFromPathMap(topic)
	return nil
}

func (p *Pod) GetInodeFromName(nameToGetMeta string, curDirInode *d.DirInode, directory *d.Directory, info *Info) (*d.DirInode, error) {
	path := info.GetCurrentDirPathAndName() + utils.PathSeperator + nameToGetMeta
	if info.IsCurrentDirRoot() {
		if strings.HasPrefix(nameToGetMeta, utils.PathSeperator) {
			path = curDirInode.Meta.Path + curDirInode.Meta.Name + nameToGetMeta
		} else {
			path = curDirInode.Meta.Path + curDirInode.Meta.Name + utils.PathSeperator + nameToGetMeta
		}
	}
	return directory.GetDirFromDirectoryMap(path), nil
}
