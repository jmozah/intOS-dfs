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
	"io"

	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) DownloadFile(podName string, podFile string) (io.ReadCloser, string, string, error) {
	if !p.isLoggedInToPod(podName) {
		return nil, "", "", fmt.Errorf("copyToLocal: login to pod to do this operation")
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return nil, "", "", fmt.Errorf("copyToLocal: %w", err)
	}

	var path string
	if podInfo.IsCurrentDirRoot() {
		path = podInfo.GetCurrentPodPathAndName() + podFile
	} else {
		path = podInfo.GetCurrentDirPathAndName() + utils.PathSeperator + podFile
	}

	if !podInfo.getFile().IsFileAlreadyPResent(path) {
		return nil, "", "", fmt.Errorf("copyToLocal: file not present in pod")
	}

	reader, ref, size, err := podInfo.getFile().Download(path)
	if err != nil {
		return nil, "", "", fmt.Errorf("copyToLocal: %w", err)
	}
	return reader, ref, size, nil
}
