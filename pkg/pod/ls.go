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

	"github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) ListPods() ([]string, error) {
	pods, err := p.loadUserPods()
	if err != nil {
		return nil, err
	}

	var listPods []string
	for _, pod := range pods {
		listPods = append(listPods, pod)
	}
	return listPods, nil
}

func (p *Pod) ListEntiesInDir(podName, dirName string) ([]dir.DirOrFileEntry, error) {
	if !p.isPodOpened(podName) {
		return nil, ErrPodNotOpened
	}

	info, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return nil, err
	}

	directory := info.getDirectory()
	printNames := false
	path := dirName // dirname is supplied in API, in REPL it is picked up from the current dir
	if path == "" {
		printNames = true
		path = info.GetCurrentDirPathAndName()
		if info.IsCurrentDirRoot() {
			path = info.GetCurrentPodPathAndName()
		}
	} else {
		path = utils.PathSeperator + podName + path
		path = strings.TrimSuffix(path, utils.PathSeperator)
	}

	return directory.ListDir(podName, path, printNames), nil
}
