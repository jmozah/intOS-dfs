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
	"bytes"
	"fmt"
	gopath "path"
	"strings"
	"time"

	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) MakeDir(podName string, dirName string) error {
	directoryName, err := CleanName(dirName)
	if err != nil {
		return fmt.Errorf("mkdir: error cleaning directory Name")
	}

	if !p.isLoggedInToPod(podName) {
		return fmt.Errorf("mkdir: login to pod to do this operation")
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	directory := podInfo.getDirectory()

	var firstTopic []byte
	var topic []byte
	var dirInode *d.DirInode
	var previousDirINode *d.DirInode
	addToPod := false

	dirs := strings.Split(directoryName, utils.PathSeperator)
	for _, dirName := range dirs {
		if len(dirName) > utils.DirectoryNameLength {
			return fmt.Errorf("mkdir: directory Name length is > %v", utils.DirectoryNameLength)
		}
	}

	// ex: mkdir make/all/this/dir
	if len(dirs) > 1 {
		for i, dirName := range dirs {
			path := p.buildPath(podInfo, dirs, i)
			_, dirInode, err = directory.GetDirNode(path, podInfo.getFeed(), podInfo.getAccountInfo())
			if err != nil {
				if previousDirINode == nil {
					if podInfo.IsCurrentDirRoot() {
						addToPod = true
					}
					dirInode, topic, err = directory.CreateDirINode(podName, dirName, podInfo.GetCurrentDirInode())
					fmt.Println("created dir ", dirName)
				} else {
					dirInode, topic, err = directory.CreateDirINode(podName, dirName, previousDirINode)
					fmt.Println("created dir ", dirName)
				}
				if err != nil {
					return fmt.Errorf("mkdir: %w", err)
				}
				if i == 0 {
					firstTopic = topic
				}

				if previousDirINode != nil {
					found := false
					for _, hash := range previousDirINode.Hashes {
						if bytes.Equal(hash, topic) {
							found = true
						}
					}
					if !found {
						previousDirINode.Hashes = append(previousDirINode.Hashes, topic)
						dirInode.Meta.Path = previousDirINode.Meta.Path + utils.PathSeperator + previousDirINode.Meta.Name
						previousDirINode.Meta.ModificationTime = time.Now().Unix()
						_, err = directory.UpdateDirectory(previousDirINode)
						if err != nil {
							return fmt.Errorf("mkdir : %w", err)
						}
					}
				}
			} else {
				fmt.Println("not creating ", dirName, path, dirInode.Meta.Path, dirInode.Meta.Name)
			}
			previousDirINode = dirInode
		}
		topic = firstTopic
	} else {
		_, topic, err = directory.CreateDirINode(podName, directoryName, podInfo.GetCurrentDirInode())
		if err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
		addToPod = true
	}

	if addToPod {
		err = p.UpdateTillThePod(podName, directory, topic, true)
		if err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
	}
	return nil
}

// Assumption is that the d.currentDirInode is the newly updated one
func (p *Pod) UpdateTillThePod(podName string, directory *d.Directory, topic []byte, isAddHash bool) error {
	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	var path string
	if podInfo.IsCurrentDirRoot() {
		path = podInfo.GetCurrentPodPathAndName()
	} else {
		path = podInfo.GetCurrentDirPathAndName()
	}

	var dirInode *d.DirInode
	for path != utils.PathSeperator {
		_, dirInode, err = directory.GetDirNode(path, podInfo.getFeed(), podInfo.getAccountInfo())
		if err != nil {
			return fmt.Errorf("update directory: %w", err)
		}
		if isAddHash {
			// Add or update a hash
			found := false
			for i, hash := range dirInode.Hashes {
				if bytes.Equal(hash, topic) {
					found = true
					dirInode.Hashes[i] = topic
				}
			}
			// ignore if it is the current dir, otherwise there will be a loop
			pathTopic := utils.HashString(path)
			if bytes.Equal(pathTopic, topic) {
				path = gopath.Dir(path)
				continue
			}
			if !found {
				dirInode.Hashes = append(dirInode.Hashes, topic)
			}
		} else {
			// remove hash
			var newHashes [][]byte
			for _, hash := range dirInode.Hashes {
				if !bytes.Equal(hash, topic) {
					newHashes = append(newHashes, hash)
				}
			}
			dirInode.Hashes = newHashes
		}
		dirInode.Meta.ModificationTime = time.Now().Unix()

		topic, err = directory.UpdateDirectory(dirInode)
		if err != nil {
			return fmt.Errorf("update directory: %w", err)
		}
		path = gopath.Dir(path)
	}
	podInfo.SetCurrentPodInode(dirInode)
	p.addPodToPodMap(podName, podInfo)
	return nil
}

func (p *Pod) buildPath(podInfo *Info, dirs []string, index int) string {
	var path string
	i := 0
	if podInfo.IsCurrentDirRoot() {
		path = podInfo.GetCurrentPodPathAndName()
	}
	for ; i <= index; i++ {
		path = path + utils.PathSeperator + dirs[i]
	}
	return path
}
