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

import "fmt"

func (p *Pod) ListPods() error {
	pods, err := p.loadUserPods()
	if err != nil {
		return fmt.Errorf("list pods: %w", err)
	}
	for _, pod := range pods {
		fmt.Print("Pod: ", pod)
	}
	fmt.Println("")
	return nil
}

func (p *Pod) ListEntiesInDir(podName string) ([]string, error) {
	if !p.isLoggedInToPod(podName) {
		return nil, fmt.Errorf("ls: login to pod to do this operation")
	}

	info, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return nil, fmt.Errorf("ls: %w", err)
	}

	directory := info.getDirectory()

	path := info.GetCurrentDirPathAndName()
	if info.IsCurrentDirRoot() {
		path = info.GetCurrentPodPathAndName()
	}

	return directory.ListDir(podName, path), nil
}
