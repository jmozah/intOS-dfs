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

	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) DeletePod(podName, dataDir string) error {
	podName, err := CleanName(podName)
	if err != nil {
		return fmt.Errorf("delete pod: %w", err)
	}

	pods, err := p.getRootFileContents()
	if err != nil {
		return fmt.Errorf("delete pod: %w", err)
	}
	found := false
	for index, pod := range pods {
		if strings.Trim(pod, "\n") == podName {
			delete(pods, index)
			found = true
		}
	}
	if !found {
		return fmt.Errorf("delete pod: pod not found")
	}

	// if last pod is deleted.. something should be there to update the feed
	if pods == nil {
		pods = make(map[int]string)
		pods[0] = ""
	}

	err = p.storeAsRootFile(pods)
	if err != nil {
		return fmt.Errorf("delete pod: %w", err)
	}

	if p.isLoggedInToPod(podName) {
		return p.LogoutPod(podName)
	} else {
		podInfo, err := p.GetPodInfoFromPodMap(podName)
		if err != nil {
			return fmt.Errorf("delete pod: %w", err)
		}
		podInfo.dir.RemoveFromDirectoryMap(podName)
		p.removePodFromPodMap(podName)
	}

	keyStore := filepath.Join(dataDir, "keystore")
	err = os.Remove(keyStore + string(utils.PathSeperator) + podName + ".key")
	if err != nil {
		return fmt.Errorf("delete pod: %w", err)
	}

	return nil
}
