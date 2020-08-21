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
	"path/filepath"
	"strings"

	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (d *Directory) ListDir(podName, path string, printNames bool) ([]string, []string) {
	d.dirMu.Lock()
	defer d.dirMu.Unlock()
	var fileListing []string
	var dirListing []string

	directory := ("<Dir>  : ")
	f := ("<File> : ")
	for k := range d.dirMap {
		if strings.HasPrefix(k, path) {
			if k != podName {
				name := strings.TrimPrefix(k, path)
				name = strings.TrimSpace(name)
				name = strings.TrimPrefix(name, utils.PathSeperator)
				if strings.ContainsAny(name, utils.PathSeperator) {
					name = utils.PathSeperator + name
				}
				if name != "" {
					if printNames {
						dirListing = append(dirListing, directory+name)
					} else {
						dirListing = append(dirListing, name)
					}
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
						fileListing = append(fileListing, f+fileName)
					} else {
						fileListing = append(fileListing, fileName)
					}
				}
			}
		}
	}
	return fileListing, dirListing
}
