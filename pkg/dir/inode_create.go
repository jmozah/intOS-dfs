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

package datapod

import (
	"encoding/json"
	"fmt"
	"time"

	m "github.com/jmozah/intOS-dfs/pkg/meta"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (d *Directory) CreateDirINode(podName string, dirName string, parent *DirInode) (*DirInode, []byte, error) {
	path := getPath(podName, parent)
	now := time.Now().Unix()
	meta := m.DirectoryMetaData{
		Version:          m.DirMetaVersion,
		Path:             path,
		Name:             dirName,
		CreationTime:     now,
		ModificationTime: now,
		AccessTime:       now,
	}

	if podName == "" {
		meta.Path = utils.PathSeperator
	}

	dirInode := &DirInode{
		Meta: &meta,
	}

	data, err := json.Marshal(dirInode)
	if err != nil {
		return nil, nil, fmt.Errorf("create inode: %w", err)
	}

	totalPath := path + utils.PathSeperator + dirName
	topic := utils.HashString(totalPath)
	if podName == utils.DefaultRoot {
		topic = utils.HashString(utils.DefaultRoot)
	}

	_, err = d.fd.CreateFeed(topic, d.acc.GetAddress(), data)
	if err != nil {
		return nil, nil, fmt.Errorf("create inode: %w", err)
	}

	d.AddToDirectoryMap(totalPath, dirInode)
	return dirInode, topic, nil
}

func getPath(podName string, parent *DirInode) string {
	var path string
	if podName == "" || parent == nil {
		path = ""
	} else {
		if parent.Meta.Path == utils.PathSeperator {
			path = parent.Meta.Path + parent.Meta.Name
		} else {
			path = parent.Meta.Path + utils.PathSeperator + parent.Meta.Name
		}
	}
	return path
}
