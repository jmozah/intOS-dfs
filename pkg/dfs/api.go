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
	"io"
	"mime/multipart"
	"net/http"

	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	"github.com/jmozah/intOS-dfs/pkg/blockstore/bee"
	"github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/file"
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
func (d *DfsAPI) CreateUser(userName, passPhrase, mnemonic string, response http.ResponseWriter, sessionId string) (string, string, error) {
	return d.users.CreateNewUser(userName, passPhrase, mnemonic, d.dataDir, d.client, response, sessionId)
}

func (d *DfsAPI) LoginUser(userName, passPhrase string, response http.ResponseWriter, sessionId string) error {
	return d.users.LoginUser(userName, passPhrase, d.dataDir, d.client, response, sessionId)
}

func (d *DfsAPI) LogoutUser(userName, sessionId string, response http.ResponseWriter) error {
	return d.users.LogoutUser(userName, d.dataDir, sessionId, response)
}

func (d *DfsAPI) DeleteUser(userName, passPhrase, sessionId string, response http.ResponseWriter) error {
	return d.users.DeleteUser(userName, d.dataDir, passPhrase, sessionId, response)
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
func (d *DfsAPI) CreatePod(userName, podName, passPhrase, sessionId string, response http.ResponseWriter, request *http.Request) (*pod.Info, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return nil, ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// create the pod
	pi, err := ui.GetPod().CreatePod(podName, passPhrase, response, request)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func (d *DfsAPI) DeletePod(userName, podName, sessionId string, response http.ResponseWriter, request *http.Request) error {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// delete the pod
	err := ui.GetPod().DeletePod(podName, response, request)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) OpenPod(userName, podName, passPhrase, sessionId string, response http.ResponseWriter, request *http.Request) (*pod.Info, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return nil, ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// open the pod
	po, err := ui.GetPod().OpenPod(podName, passPhrase, response, request)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (d *DfsAPI) ClosePod(userName, podName, sessionId string, response http.ResponseWriter, request *http.Request) error {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// close the pod
	err := ui.GetPod().ClosePod(podName, response, request)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) PodStat(userName, podName, sessionId string) (*pod.PodStat, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return nil, ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// get the pod stat
	podStat, err := ui.GetPod().PodStat(podName)
	if err != nil {
		return nil, err
	}
	return podStat, nil
}

func (d *DfsAPI) SyncPod(userName, podName, sessionId string) error {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// sync the pod
	err := ui.GetPod().SyncPod(podName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ListPods(userName, sessionId string) ([]string, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return nil, ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// list pods of a user
	pods, err := ui.GetPod().ListPods()
	if err != nil {
		return nil, err
	}
	return pods, nil
}

//
//  Directory related APIs
//

func (d *DfsAPI) Mkdir(userName, podName, directoryName, sessionId string) error {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// make dir
	err := ui.GetPod().MakeDir(podName, directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) RmDir(userName, podName, directoryName, sessionId string) error {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	err := ui.GetPod().RemoveDir(podName, directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ListDir(userName, podName, currentDir, sessionId string) ([]dir.DirOrFileEntry, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return nil, ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	entries, err := ui.GetPod().ListEntiesInDir(podName, currentDir)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (d *DfsAPI) DirectoryStat(userName, podName, directoryName, sessionId string) (*dir.DirStats, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return nil, ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	ds, err := ui.GetPod().DirectoryStat(podName, directoryName)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func (d *DfsAPI) ChangeDirectory(userName, podName, directoryName, sessionId string) (*pod.Info, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return nil, ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
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

func (d *DfsAPI) CopyFromLocal(userName, podName, localFile, podDir, blockSize, sessionId string) error {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	err := ui.GetPod().CopyFromLocal(podName, localFile, podDir, blockSize)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) CopyToLocal(userName, podName, localDir, podFile, sessionId string) error {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	err := ui.GetPod().CopyToLocal(podName, localDir, podFile)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) Cat(userName, podName, fileName, sessionId string) error {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	err := ui.GetPod().Cat(podName, fileName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) DeleteFile(userName, podName, podFile, sessionId string) error {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	err := ui.GetPod().RemoveFile(podName, podFile)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) FileStat(userName, podName, fileName, sessionId string) (*file.FileStats, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return nil, ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	ds, err := ui.GetPod().FileStat(podName, fileName)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func (d *DfsAPI) UploadFile(userName, podName, fileName, sessionId string, fileSize int64, fd multipart.File, podDir, blockSize string) (string, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return "", ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return "", ErrUserNotLoggedIn
	}

	ref, err := ui.GetPod().UploadFile(podName, fileName, fileSize, fd, podDir, blockSize)
	if err != nil {
		return "", err
	}
	return ref, nil
}

func (d *DfsAPI) DownloadFile(userName, podName, podFile, sessionId string) (io.ReadCloser, string, string, error) {
	// check if the user is valid
	if !d.users.IsUsernameAvailable(userName, d.dataDir) {
		return nil, "", "", ErrInvalidUserName
	}

	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(userName, sessionId)
	if ui == nil {
		return nil, "", "", ErrUserNotLoggedIn
	}

	reader, ref, size, err := ui.GetPod().DownloadFile(podName, podFile)
	if err != nil {
		return nil, "", "", err
	}
	return reader, ref, size, nil
}
