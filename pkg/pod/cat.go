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

	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) Cat(podName string, fileName string) error {

	if !p.isPodOpened(podName) {
		return fmt.Errorf("copyFromLocal: login to pod to do this operation")
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("cat: %w", err)
	}

	var fname string
	if podInfo.IsCurrentDirRoot() {
		fname = podInfo.GetCurrentPodPathAndName() + fileName
	} else {
		fname = podInfo.GetCurrentDirPathAndName() + utils.PathSeperator + fileName
	}

	return podInfo.getFile().Cat(fname)
}
