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

package dfs

import (
	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	"github.com/jmozah/intOS-dfs/pkg/blockstore/bee"
	"github.com/jmozah/intOS-dfs/pkg/pod"
)

type DfsAPI struct {
	dataDir string
	client  blockstore.Client
	pod     *pod.Pod
}

func NewDfsAPI(dataDir, host, port string) *DfsAPI {
	c := bee.NewBeeClient(host, port)
	//c := mock.NewMockBeeClient()
	return &DfsAPI{
		dataDir: dataDir,
		client:  c,
		pod:     pod.NewPod(c),
	}
}

//
//  Pods related APIs
//
func (d *DfsAPI) IsInitialized(dataDir string) bool {
	return d.pod.IsInitialized(dataDir)
}

func (d *DfsAPI) RemoveRootKey(dataDir string) error {
	return d.pod.RemoveRootKey(dataDir)
}

func (d *DfsAPI) Init(dataDir, passPhrase string) error {
	err := d.pod.LoadRootPod(dataDir, passPhrase)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) CreatePod(podName string, passPhrase string) (*pod.Info, error) {
	po, err := d.pod.CreatePod(podName, d.dataDir, passPhrase)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (d *DfsAPI) DeletePod(podName string) error {
	err := d.pod.DeletePod(podName, d.dataDir)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) LoginPod(podName string, passPhrase string) (*pod.Info, error) {
	po, err := d.pod.LoginPod(podName, d.dataDir, passPhrase)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (d *DfsAPI) LogoutPod(podName string) error {
	err := d.pod.LogoutPod(podName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) PodStat(podName string) (*pod.PodStat, error) {
	podStat, err := d.pod.PodStat(podName)
	if err != nil {
		return nil, err
	}
	return podStat, nil
}

func (d *DfsAPI) SyncPod(podName string) error {
	err := d.pod.SyncPod(podName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ListPods() error {
	err := d.pod.ListPods()
	if err != nil {
		return err
	}
	return nil
}

//
//  Directory related APIs
//

func (d *DfsAPI) Mkdir(podName string, directoryName string) error {
	err := d.pod.MakeDir(podName, directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) RmDir(podName string, directoryName string) error {
	err := d.pod.RemoveDir(podName, directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ListDir(podName string) ([]string, error) {
	listing, err := d.pod.ListEntiesInDir(podName)
	if err != nil {
		return nil, err
	}
	return listing, nil
}

func (d *DfsAPI) DirectoryOrFileStat(podName string, directoryName string) error {
	err := d.pod.DirectoryOrFileStat(podName, directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ChangeDirectory(podName string, directoryName string) (*pod.Info, error) {
	podInfo, err := d.pod.ChangeDir(podName, directoryName)
	if err != nil {
		return nil, err
	}
	return podInfo, nil
}

//
// File related API's
//

func (d *DfsAPI) CopyFromLocal(podName, localFile string, podDir string, blockSize string) error {
	err := d.pod.CopyFromLocal(podName, localFile, podDir, blockSize)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) CopyToLocal(podName, localDir string, podFile string) error {
	err := d.pod.CopyToLocal(podName, localDir, podFile)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) Cat(podName string, fileName string) error {
	err := d.pod.Cat(podName, fileName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) RemoveFile(podName string, podFile string) error {
	err := d.pod.RemoveFile(podName, podFile)
	if err != nil {
		return err
	}
	return nil
}
