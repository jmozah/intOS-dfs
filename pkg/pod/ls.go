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
)

func (p *Pod) ListPods() ([]string, error) {
	pods, err := p.loadUserPods()
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}

	var listPods []string
	for _, pod := range pods {
		listPods = append(listPods, pod)
	}
	return listPods, nil
}

func (p *Pod) ListEntiesInDir(podName string, dirName string) ([]string, []string, error) {
	if !p.isPodOpened(podName) {
		return nil, nil, ErrPodNotOpened
	}

	info, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return nil, nil, fmt.Errorf("ls: %w", err)
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
	}
	fl, dl := directory.ListDir(podName, path, printNames)
	return fl, dl, nil
}
