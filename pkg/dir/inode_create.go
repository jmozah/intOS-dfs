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

package dir

import (
	"encoding/json"
	"fmt"
	"time"

	m "github.com/jmozah/intOS-dfs/pkg/meta"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (d *Directory) CreateDirINode(podName string, dirName string, parent *DirInode) (*DirInode, []byte, error) {
	// create the meta data
	parentPath := getPath(podName, parent)
	now := time.Now().Unix()
	meta := m.DirectoryMetaData{
		Version:          m.DirMetaVersion,
		Path:             parentPath,
		Name:             dirName,
		CreationTime:     now,
		ModificationTime: now,
		AccessTime:       now,
	}
	dirInode := &DirInode{
		Meta: &meta,
	}
	data, err := json.Marshal(dirInode)
	if err != nil {
		return nil, nil, fmt.Errorf("create inode: %w", err)
	}

	// create a feed for the directory and add data to it
	totalPath := parentPath + utils.PathSeperator + dirName
	topic := utils.HashString(totalPath)
	_, err = d.fd.CreateFeed(topic, d.acc.GetAddress(), data)
	if err != nil {
		return nil, nil, fmt.Errorf("create inode: %w", err)
	}

	d.AddToDirectoryMap(totalPath, dirInode)
	return dirInode, topic, nil
}

func (d *Directory) IsDirINodePresent(podName string, dirName string, parent *DirInode) bool {
	parentPath := getPath(podName, parent)
	totalPath := parentPath + utils.PathSeperator + dirName
	topic := utils.HashString(totalPath)
	_, _, err := d.fd.GetFeedData(topic, d.getAccount().GetAddress())
	if err != nil {
		return false
	}
	return true
}

func getPath(podName string, parent *DirInode) string {
	var path string
	if parent.Meta.Path == utils.PathSeperator {
		path = parent.Meta.Path + parent.Meta.Name
	} else {
		path = parent.Meta.Path + utils.PathSeperator + parent.Meta.Name
	}
	return path
}

func (d *Directory) CreatePodINode(podName string) (*DirInode, []byte, error) {
	// create the metadata
	now := time.Now().Unix()
	meta := m.DirectoryMetaData{
		Version:          m.DirMetaVersion,
		Path:             "/",
		Name:             podName,
		CreationTime:     now,
		ModificationTime: now,
		AccessTime:       now,
	}
	dirInode := &DirInode{
		Meta: &meta,
	}
	data, err := json.Marshal(dirInode)
	if err != nil {
		return nil, nil, fmt.Errorf("create pod inode: %w", err)
	}

	// create a feed and store the metadata of the pod
	totalPath := utils.PathSeperator + podName
	topic := utils.HashString(totalPath)
	_, err = d.fd.CreateFeed(topic, d.acc.GetAddress(), data)
	if err != nil {
		return nil, nil, fmt.Errorf("create pod inode: %w", err)
	}

	d.AddToDirectoryMap(totalPath, dirInode)
	return dirInode, topic, nil
}
