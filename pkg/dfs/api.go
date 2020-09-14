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
	"github.com/jmozah/intOS-dfs/pkg/logging"
	"github.com/jmozah/intOS-dfs/pkg/pod"
	"github.com/jmozah/intOS-dfs/pkg/user"
)

type DfsAPI struct {
	dataDir string
	client  blockstore.Client
	users   *user.Users
	logger  logging.Logger
}

func NewDfsAPI(dataDir, host, port string, logger logging.Logger) *DfsAPI {
	c := bee.NewBeeClient(host, port, logger)
	users := user.NewUsers(dataDir, c, logger)
	return &DfsAPI{
		dataDir: dataDir,
		client:  c,
		users:   users,
		logger:  logger,
	}
}

//
//  User related APIs
//
func (d *DfsAPI) CreateUser(userName, passPhrase, mnemonic string, response http.ResponseWriter, sessionId string) (string, string, error) {
	reference, rcvdMnemonic, userInfo, err := d.users.CreateNewUser(userName, passPhrase, mnemonic, d.dataDir, d.client, response, sessionId)
	if err != nil {
		return reference, rcvdMnemonic, err
	}

	// TODO: check if the connection is there before creating user
	err = d.users.CreateRootFeeds(userInfo)
	if err != nil {
		return reference, rcvdMnemonic, err
	}
	return reference, rcvdMnemonic, nil
}

func (d *DfsAPI) LoginUser(userName, passPhrase string, response http.ResponseWriter, sessionId string) error {
	return d.users.LoginUser(userName, passPhrase, d.dataDir, d.client, response, sessionId)
}

func (d *DfsAPI) LogoutUser(sessionId string, response http.ResponseWriter) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	return d.users.LogoutUser(ui.GetUserName(), d.dataDir, sessionId, response)
}

func (d *DfsAPI) DeleteUser(passPhrase, sessionId string, response http.ResponseWriter) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	return d.users.DeleteUser(ui.GetUserName(), d.dataDir, passPhrase, sessionId, response)
}

func (d *DfsAPI) IsUserNameAvailable(userName string) bool {
	return d.users.IsUsernameAvailable(userName, d.dataDir)
}

func (d *DfsAPI) IsUserLoggedIn(userName string) bool {
	// check if a given user is logged in
	return d.users.IsUserNameLoggedIn(userName)
}

func (d *DfsAPI) ListAllUsers() []string {
	return d.users.ListAllUsers(d.dataDir)
}

func (d *DfsAPI) SaveAvatar(sessionId string, data []byte) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	return d.users.SaveAvatar(data, ui)
}

func (d *DfsAPI) GetAvatar(sessionId string) ([]byte, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	return d.users.GetAvatar(ui)
}

func (d *DfsAPI) SaveName(firstName, lastName, middleName, surname, sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}
	return d.users.SaveName(firstName, lastName, middleName, surname, ui)
}

func (d *DfsAPI) GetName(sessionId string) (*user.Name, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}
	return d.users.GetName(ui)
}

func (d *DfsAPI) SaveContact(phone, mobile string, address *user.Address, sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}
	return d.users.SaveContacts(phone, mobile, address, ui)
}

func (d *DfsAPI) GetContact(sessionId string) (*user.Contacts, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}
	return d.users.GetContacts(ui)
}

func (d *DfsAPI) GetUserStat(sessionId string) (*user.Stat, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}
	return d.users.GetUserStat(ui)
}

func (d *DfsAPI) GetUserSharingInbox(sessionId string) (*user.Inbox, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}
	return d.users.GetSharingInbox(ui)
}

func (d *DfsAPI) GetUserSharingOutbox(sessionId string) (*user.Outbox, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}
	return d.users.GetSharingOutbox(ui)
}

//
//  Pods related APIs
//
func (d *DfsAPI) CreatePod(podName, passPhrase, sessionId string) (*pod.Info, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// create the pod
	pi, err := ui.GetPod().CreatePod(podName, passPhrase)
	if err != nil {
		return nil, err
	}

	// open the pod
	_, err = ui.GetPod().OpenPod(podName, passPhrase)
	if err != nil {
		return nil, err
	}

	// Add podName in the login user session
	ui.SetPodName(podName)
	return pi, nil
}

func (d *DfsAPI) DeletePod(podName, sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// delete the pod
	err := ui.GetPod().DeletePod(podName)
	if err != nil {
		return err
	}

	// close the pod and delete it from login user session, if the delete is for a opened pod
	if ui.GetPodName() != "" && podName == ui.GetPodName() {
		// close the pod
		err = ui.GetPod().ClosePod(ui.GetPodName())
		if err != nil {
			return err
		}

		// remove from the login session
		ui.RemovePodName()
	}

	return nil
}

func (d *DfsAPI) OpenPod(podName, passPhrase, sessionId string) (*pod.Info, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// close the already open pod
	if ui.GetPodName() != "" {
		err := ui.GetPod().ClosePod(ui.GetPodName())
		if err != nil {
			return nil, err
		}
	}

	// open the pod
	po, err := ui.GetPod().OpenPod(podName, passPhrase)
	if err != nil {
		return nil, err
	}

	// Add podName in the login user session
	ui.SetPodName(podName)
	return po, nil
}

func (d *DfsAPI) ClosePod(sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return ErrPodNotOpen
	}

	// close the pod
	err := ui.GetPod().ClosePod(ui.GetPodName())
	if err != nil {
		return err
	}

	// delete podName in the login user session
	ui.RemovePodName()
	return nil
}

func (d *DfsAPI) PodStat(podName, sessionId string) (*pod.PodStat, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
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

func (d *DfsAPI) SyncPod(sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return ErrPodNotOpen
	}

	// sync the pod
	err := ui.GetPod().SyncPod(ui.GetPodName())
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ListPods(sessionId string) ([]string, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
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

func (d *DfsAPI) Mkdir(directoryName, sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return ErrPodNotOpen
	}

	// make dir
	err := ui.GetPod().MakeDir(ui.GetPodName(), directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) RmDir(directoryName, sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return ErrPodNotOpen
	}

	err := ui.GetPod().RemoveDir(ui.GetPodName(), directoryName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) ListDir(currentDir, sessionId string) ([]dir.DirOrFileEntry, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return nil, ErrPodNotOpen
	}

	entries, err := ui.GetPod().ListEntiesInDir(ui.GetPodName(), currentDir)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (d *DfsAPI) DirectoryStat(directoryName, sessionId string) (*dir.DirStats, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return nil, ErrPodNotOpen
	}

	ds, err := ui.GetPod().DirectoryStat(ui.GetPodName(), directoryName)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func (d *DfsAPI) ChangeDirectory(directoryName, sessionId string) (*pod.Info, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return nil, ErrPodNotOpen
	}

	podInfo, err := ui.GetPod().ChangeDir(ui.GetPodName(), directoryName)
	if err != nil {
		return nil, err
	}
	return podInfo, nil
}

//
// File related API's
//

func (d *DfsAPI) CopyFromLocal(localFile, podDir, blockSize, sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return ErrPodNotOpen
	}

	err := ui.GetPod().CopyFromLocal(ui.GetPodName(), localFile, podDir, blockSize)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) CopyToLocal(localDir, podFile, sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return ErrPodNotOpen
	}

	err := ui.GetPod().CopyToLocal(ui.GetPodName(), localDir, podFile)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) Cat(fileName, sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return ErrPodNotOpen
	}

	err := ui.GetPod().Cat(ui.GetPodName(), fileName)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) DeleteFile(podFile, sessionId string) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return ErrPodNotOpen
	}

	err := ui.GetPod().RemoveFile(ui.GetPodName(), podFile)
	if err != nil {
		return err
	}
	return nil
}

func (d *DfsAPI) FileStat(fileName, sessionId string) (*file.FileStats, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return nil, ErrPodNotOpen
	}

	ds, err := ui.GetPod().FileStat(ui.GetPodName(), fileName)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func (d *DfsAPI) UploadFile(fileName, sessionId string, fileSize int64, fd multipart.File, podDir, blockSize string, compression string) (string, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return "", ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return "", ErrPodNotOpen
	}

	ref, err := ui.GetPod().UploadFile(ui.GetPodName(), fileName, fileSize, fd, podDir, blockSize, compression)
	if err != nil {
		return "", err
	}
	return ref, nil
}

func (d *DfsAPI) DownloadFile(podFile, sessionId string) (io.ReadCloser, string, string, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, "", "", ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return nil, "", "", ErrPodNotOpen
	}

	reader, ref, size, err := ui.GetPod().DownloadFile(ui.GetPodName(), podFile)
	if err != nil {
		return nil, "", "", err
	}
	return reader, ref, size, nil
}

func (d *DfsAPI) ShareFile(podFile, destinationUser, sessionId string) (*user.OutboxEntry, error) {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return nil, ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return nil, ErrPodNotOpen
	}

	outEntry, err := d.users.ShareFileWithUser(ui.GetPodName(), podFile, destinationUser, ui, ui.GetPod())
	if err != nil {
		return nil, err
	}
	return outEntry, nil
}

func (d *DfsAPI) ReceiveFile(sessionId string, inboxEntry user.InboxEntry) error {
	// get the logged in user information
	ui := d.users.GetLoggedInUserInfo(sessionId)
	if ui == nil {
		return ErrUserNotLoggedIn
	}

	// check if pod open
	if ui.GetPodName() == "" {
		return ErrPodNotOpen
	}

	return d.users.ReceiveFileFromUser(ui.GetPodName(), inboxEntry, ui, ui.GetPod())
}
