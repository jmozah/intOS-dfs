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

	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (d *Directory) UpdateDirectory(dirInode *DirInode) ([]byte, error) {
	dirName := dirInode.Meta.Name
	path := dirInode.Meta.Path
	meta := dirInode.Meta
	meta.ModificationTime = time.Now().Unix()
	dirInode.Meta = meta

	data, err := json.Marshal(dirInode)
	if err != nil {
		return nil, fmt.Errorf("could not marshall directory: %v", dirName)
	}

	curDir := path + utils.PathSeperator + dirName
	if path == utils.PathSeperator {
		curDir = path + dirName
	}
	topic := utils.HashString(curDir)
	_, err = d.getFeed().UpdateFeed(topic, d.getAccount().GetAddress(), data)
	if err != nil {
		return nil, fmt.Errorf("could not update feed for dir: %v", dirName)
	}

	d.AddToDirectoryMap(curDir, dirInode)
	return topic, nil
}
