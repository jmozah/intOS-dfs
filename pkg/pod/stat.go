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
	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/jmozah/intOS-dfs/pkg/utils"
	"strconv"
	"time"
)

type PodStat struct {
	Version          string
	PodName          string
	PodPath          string
	CreationTime     string
	AccessTime       string
	ModificationTime string
}

func (p *Pod) PodStat(podName string) (*PodStat, error) {
	podInfo, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return nil, fmt.Errorf("pod stat: %w", err)
	}
	podInode := podInfo.GetCurrentPodInode()
	return &PodStat{
		Version:          strconv.Itoa(int(podInode.Meta.Version)),
		PodName:          podInode.Meta.Name,
		PodPath:          podInode.Meta.Path,
		CreationTime:     time.Unix(podInode.Meta.CreationTime, 0).String(),
		AccessTime:       time.Unix(podInode.Meta.AccessTime, 0).String(),
		ModificationTime: time.Unix(podInode.Meta.AccessTime, 0).String(),
	}, nil
}

func (p *Pod) DirectoryOrFileStat(podName, podFileOrDir string) error {
	if !p.isLoggedInToPod(podName) {
		return fmt.Errorf("login to pod to do this operation")
	}

	info, err := p.GetPodInfoFromPodMap(podName)
	if err != nil {
		return fmt.Errorf("rmdir: %w", err)
	}

	acc := info.getAccountInfo().GetAddress()
	account := swarm.NewAddress(acc[:]).String()

	path := p.getDirectoryPath(podFileOrDir, info)
	dirInode := info.getDirectory().GetDirFromDirectoryMap(path)
	if dirInode != nil {
		meta := dirInode.Meta
		addr, dirInode, err := info.getDirectory().GetDirNode(meta.Path+utils.PathSeperator+meta.Name, info.getFeed(), info.getAccountInfo())
		if err != nil {
			return fmt.Errorf("could not get dirnode: %w", err)
		}
		podAddress := swarm.NewAddress(addr).String()
		return info.getDirectory().DirStat(podName, path, dirInode, account, podAddress)
	}

	if !info.file.IsFileAlreadyPResent(path) {
		return fmt.Errorf("file not present in pod")
	}
	return info.file.FileStat(podName, path, account)
}
