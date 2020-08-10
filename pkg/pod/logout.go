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

func (p *Pod) LogoutPod(podName string) error {
	podName, err := CleanName(podName)
	if err != nil {
		return fmt.Errorf("login pod: %w", err)
	}

	if !p.isLoggedInToPod(podName) {
		return fmt.Errorf("logout pod: login to pod to do this operation")
	}

	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("logout pod: %w", err)
	}

	p.removePodFromPodMap(podName)
	podInfo.dir.RemoveFromDirectoryMap(podName)
	return nil
}
