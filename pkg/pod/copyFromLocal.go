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
	"os"
	"path/filepath"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) CopyFromLocal(podName string, localFile string, podDir string, blockSize string) error {
	if !p.isPodOpened(podName) {
		return fmt.Errorf("copyFromLocal: login to pod to do this operation")
	}

	if len(filepath.Base(localFile)) > utils.FileNameLength {
		return fmt.Errorf("copyFromLocal: file Name length is > %v", utils.FileNameLength)
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("copyFromLocal: %w", err)
	}
	dir := podInfo.getDirectory()

	bs, err := humanize.ParseBytes(blockSize)
	if err != nil {
		return fmt.Errorf("copyFromLocal: block size parse error: %w", err)
	}

	path := p.getFilePath(podDir, podInfo)

	_, dirInode, err := dir.GetDirNode(path, podInfo.getFeed(), podInfo.getAccountInfo())
	if err != nil {
		return fmt.Errorf("error while fetching pod info: %w", err)
	}

	var localDir string
	var suffixRegex string
	if strings.HasSuffix(localFile, "*") {
		localDir = filepath.Dir(localFile)
		suffixRegex = filepath.Base(localFile)
	}

	fileInfo, err := os.Stat(localFile)
	if err != nil {
		if localDir != "" {
			fileInfo, err = os.Stat(localDir)
			if err != nil {
				return fmt.Errorf("local dir not present : %w", err)
			}
		} else {
			return fmt.Errorf("local file not present : %w", err)
		}
	}

	if fileInfo.IsDir() {
		suffix := strings.TrimSuffix(suffixRegex, "*")
		err = filepath.Walk(localDir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if strings.HasPrefix(info.Name(), suffix) {
					if podInfo.file.IsFileAlreadyPResent(path + utils.PathSeperator + info.Name()) {
						return fmt.Errorf("file already present in the destination dir")
					}

					localFile = localDir + utils.PathSeperator + info.Name()
					var fpath string
					if podInfo.GetCurrentDirInode().IsDirInodeRoot() {
						fpath = podInfo.GetCurrentPodPathAndName() + utils.PathSeperator + info.Name()
					} else {
						fpath = podInfo.GetCurrentDirPathAndName() + utils.PathSeperator + info.Name()
					}
					addr, err := podInfo.file.CopyFromFile(podName, localFile, info, uint32(bs), fpath)
					if err != nil {
						return fmt.Errorf("error while copying file to pod: %w", err)
					}
					dirInode.Hashes = append(dirInode.Hashes, addr)
				}

				return nil
			})
		if err != nil {
			return err
		}
	} else {
		fpath := path + utils.PathSeperator + fileInfo.Name()
		if podInfo.file.IsFileAlreadyPResent(fpath) {
			return fmt.Errorf("file already present in the destination dir")
		}
		addr, err := podInfo.file.CopyFromFile(podName, localFile, fileInfo, uint32(bs), fpath)
		if err != nil {
			return fmt.Errorf("error while copying file to pod: %w", err)
		}
		dirInode.Hashes = append(dirInode.Hashes, addr)

	}

	dirInode.Meta.ModificationTime = time.Now().Unix()
	topic, err := dir.UpdateDirectory(dirInode)
	if err != nil {
		return fmt.Errorf("error updating directory: %w", err)
	}

	if path != podInfo.GetCurrentPodPathAndName() {
		err = p.UpdateTillThePod(podName, podInfo.getDirectory(), topic, path, true)
		if err != nil {
			return fmt.Errorf("error updating directory: %w", err)
		}
	}
	return nil
}

func (p *Pod) getFilePath(podDir string, podInfo *Info) string {
	var path string

	if podDir == podInfo.GetCurrentPodPathAndName() {
		return podInfo.GetCurrentPodPathAndName()
	}

	// this is a full path.. so use it as it is
	if strings.HasPrefix(podDir, "/") {
		return podInfo.GetCurrentPodPathAndName() + podDir
	}

	if podInfo.IsCurrentDirRoot() {
		if podDir == "." {
			path = podInfo.GetCurrentPodPathAndName()
		} else {
			path = podInfo.GetCurrentPodPathAndName() + utils.PathSeperator + podDir
		}
	} else {
		if podDir == "." {
			path = podInfo.GetCurrentDirPathAndName()
		} else {
			path = podInfo.GetCurrentDirPathAndName() + utils.PathSeperator + podDir
		}
	}
	return path
}
