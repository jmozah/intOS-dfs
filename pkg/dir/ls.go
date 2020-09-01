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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	m "github.com/jmozah/intOS-dfs/pkg/meta"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

const (
	MineTypeDirectory = "inode/directory"
)

type DirOrFileEntry struct {
	Name             string `json:"name"`
	ContentType      string `json:"content_type"`
	Size             string `json:"size,omitempty"`
	CreationTime     string `json:"creation_time"`
	ModificationTime string `json:"modification_time"`
	AccessTime       string `json:"access_time"`
}

func (d *Directory) ListDir(podName, path string, printNames bool) []DirOrFileEntry {
	_, dirInode, err := d.GetDirNode(path, d.getFeed(), d.getAccount())
	if err != nil {
		return nil
	}

	var listEntries []DirOrFileEntry
	for _, ref := range dirInode.Hashes {
		// check if this is a directory
		_, data, err := d.getFeed().GetFeedData(ref, d.getAccount().GetAddress())
		if err != nil {
			// if it is not a dir, then treat this reference as a file
			data, _, err := d.getClient().DownloadBlob(ref)
			if err != nil {
				continue
			}
			var meta *m.FileMetaData
			err = json.Unmarshal(data, &meta)
			if err != nil {
				continue
			}

			entry := DirOrFileEntry{
				Name:             meta.Name,
				ContentType:      meta.ContentType,
				Size:             strconv.FormatUint(meta.FileSize, 10),
				CreationTime:     time.Unix(meta.CreationTime, 0).String(),
				AccessTime:       time.Unix(meta.AccessTime, 0).String(),
				ModificationTime: time.Unix(meta.ModificationTime, 0).String(),
			}
			listEntries = append(listEntries, entry)
			continue
		}

		var dirInode *DirInode
		err = json.Unmarshal(data, &dirInode)
		if err != nil {
			continue
		}
		entry := DirOrFileEntry{
			Name:             dirInode.Meta.Name,
			ContentType:      MineTypeDirectory, // per RFC2425
			CreationTime:     time.Unix(dirInode.Meta.CreationTime, 0).String(),
			AccessTime:       time.Unix(dirInode.Meta.AccessTime, 0).String(),
			ModificationTime: time.Unix(dirInode.Meta.ModificationTime, 0).String(),
		}
		listEntries = append(listEntries, entry)
	}
	return listEntries

}

func (d *Directory) ListDirOnlyNames(podName, path string, printNames bool) ([]string, []string) {
	d.dirMu.Lock()
	defer d.dirMu.Unlock()
	var fileListing []string
	var dirListing []string

	directory := ("<Dir>  : ")
	fl := ("<File> : ")
	for k := range d.dirMap {
		if strings.HasPrefix(k, path) {
			name := strings.TrimPrefix(k, path)
			if name != "" {
				if printNames {
					dirListing = append(dirListing, directory+name)
				} else {
					dirListing = append(dirListing, name)

				}
			}

			// Get the files inside the dir
			fileList := d.file.ListFiles(k)
			for _, file := range fileList {
				if strings.HasPrefix(file, path) {
					if filepath.Dir(file) != k {
						continue
					}
					var fileName string
					fileName = strings.TrimPrefix(file, path)
					fileName = strings.TrimSpace(fileName)
					fileName = strings.TrimPrefix(fileName, utils.PathSeperator)
					if strings.ContainsAny(fileName, utils.PathSeperator) {
						fileName = utils.PathSeperator + fileName
					}
					if printNames {
						fileListing = append(fileListing, fl+fileName)
					} else {
						fileListing = append(fileListing, fileName)
					}
				}
			}
		}
	}
	return fileListing, dirListing
}
