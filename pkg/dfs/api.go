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
	"fmt"

	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	"github.com/jmozah/intOS-dfs/pkg/blockstore/bee"
	"github.com/jmozah/intOS-dfs/pkg/pod"
	"github.com/jmozah/intOS-dfs/pkg/user"
)

type DfsAPI struct {
	dataDir string
	client  blockstore.Client
	users   *user.Users
}

func NewDfsAPI(dataDir, host, port string) *DfsAPI {
	c := bee.NewBeeClient(host, port)
	users := user.NewUsers(dataDir, c)
	return &DfsAPI{
		dataDir: dataDir,
		client:  c,
		users:   users,
	}
}

//
//  User related APIs
//
func (d *DfsAPI) CreateUser(userName, passPhrase string) (string, string, error) {
	return d.users.CreateNewUser(userName, passPhrase, d.dataDir, d.client)
}

func (d *DfsAPI) LoginUser(userName, passPhrase string) error {
	return d.users.LoginUser(userName, passPhrase, d.dataDir, d.client)
}

func (d *DfsAPI) LogoutUser(userName string) error {
	return d.users.LogoutUser(userName, d.dataDir)
}

func (d *DfsAPI) DeleteUser(userName, passPhrase string) error {
	return d.users.DeleteUser(userName, d.dataDir, passPhrase)
}

func (d *DfsAPI) IsUserNameAvailable(userName string) bool {
	return d.users.IsUsernameAvailable(userName, d.dataDir)
}

func (d *DfsAPI) ListAllUsers() []string {
	return d.users.ListAllUsers(d.dataDir)
}

//
//  Pods related APIs
//
func (d *DfsAPI) CreatePod(userName string, podName string, passPhrase string) (*pod.Info, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return nil, fmt.Errorf("create pod: login as a user to execute this command")
	}

	// create the pod
	pi, err := ui.GetPod().CreatePod(podName, d.dataDir, passPhrase)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (d *DfsAPI) DeletePod(userName string, podName string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("delete pod: login as a user to execute this command")
	}

	// delete the pod
	err := ui.GetPod().DeletePod(podName, d.dataDir)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) OpenPod(userName string, podName string, passPhrase string) (*pod.Info, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return nil, fmt.Errorf("open pod: login as a user to execute this command")
	}

	// open the pod
	po, err := ui.GetPod().OpenPod(podName, d.dataDir, passPhrase)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (d *DfsAPI) ClosePod(userName string, podName string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("close pod: login as a user to execute this command")
	}

	// close the pod
	err := ui.GetPod().ClosePod(podName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) PodStat(userName string, podName string) (*pod.PodStat, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return nil, fmt.Errorf("pod stat: login as a user to execute this command")
	}

	// get the pod stat
	podStat, err := ui.GetPod().PodStat(podName)
	if err != nil {
		return nil, err
	}
	return podStat, nil
}

func (d *DfsAPI) SyncPod(userName string, podName string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("sync pod: login as a user to execute this command")
	}

	// sync the pod
	err := ui.GetPod().SyncPod(podName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ListPods(userName string, print bool) ([]string, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return nil, fmt.Errorf("sync pod: login as a user to execute this command")
	}

	// list pods of a user
	pods, err := ui.GetPod().ListPods(print)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

//
//  Directory related APIs
//

func (d *DfsAPI) Mkdir(userName string, podName string, directoryName string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("mkdir: login as a user to execute this command")
	}

	// make dir
	err := ui.GetPod().MakeDir(podName, directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) RmDir(userName string, podName string, directoryName string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("rmdir: login as a user to execute this command")
	}

	err := ui.GetPod().RemoveDir(podName, directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ListDir(userName string, podName string) ([]string, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return nil, fmt.Errorf("ls dir: login as a user to execute this command")
	}

	listing, err := ui.GetPod().ListEntiesInDir(podName)
	if err != nil {
		return nil, err
	}
	return listing, nil
}

func (d *DfsAPI) DirectoryOrFileStat(userName string, podName string, directoryName string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("stat: login as a user to execute this command")
	}

	err := ui.GetPod().DirectoryOrFileStat(podName, directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ChangeDirectory(userName string, podName string, directoryName string) (*pod.Info, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return nil, fmt.Errorf("cd: login as a user to execute this command")
	}

	podInfo, err := ui.GetPod().ChangeDir(podName, directoryName)
	if err != nil {
		return nil, err
	}
	return podInfo, nil
}

//
// File related API's
//

func (d *DfsAPI) CopyFromLocal(userName string, podName, localFile string, podDir string, blockSize string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("copyFromLocal: login as a user to execute this command")
	}

	err := ui.GetPod().CopyFromLocal(podName, localFile, podDir, blockSize)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) CopyToLocal(userName string, podName, localDir string, podFile string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("copyToLocal: login as a user to execute this command")
	}

	err := ui.GetPod().CopyToLocal(podName, localDir, podFile)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) Cat(userName string, podName string, fileName string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("cat: login as a user to execute this command")
	}

	err := ui.GetPod().Cat(podName, fileName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) RemoveFile(userName string, podName string, podFile string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName)
	if ui == nil {
		return fmt.Errorf("cat: login as a user to execute this command")
	}

	err := ui.GetPod().RemoveFile(podName, podFile)
	if err != nil {
		return err
	}
	return nil
}
