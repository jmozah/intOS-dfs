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

package user

import (
	"fmt"
	"net/http"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	"github.com/jmozah/intOS-dfs/pkg/cookie"
	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	f "github.com/jmozah/intOS-dfs/pkg/file"
	"github.com/jmozah/intOS-dfs/pkg/pod"
)

func (u *Users) LoginUser(userName string, passPhrase string, dataDir string, client blockstore.Client, response http.ResponseWriter, sessionId string) error {
	if u.isUserPresentInMap(userName) {
		return ErrUserAlreadyLoggedIn
	}

	if !u.IsUsernameAvailable(userName, dataDir) {
		return ErrInvalidUserName
	}

	acc := account.New(userName, dataDir)
	accountInfo := acc.GetAccountInfo(account.UserAccountIndex)
	fd := feed.New(accountInfo, client)
	file := f.NewFile(userName, client, fd, accountInfo)
	err := acc.LoadUserAccount(passPhrase)
	if err != nil {
		if err.Error() == "mnemonic is invalid" {
			return ErrInvalidPassword
		}
		return fmt.Errorf("user login: %w", err)
	}
	dir := d.NewDirectory(userName, client, fd, accountInfo, file)

	if sessionId == "" {
		sessionId = cookie.GetUniqueSessionId()
	}

	ui := &Info{
		name:      userName,
		sessionId: sessionId,
		feedApi:   fd,
		account:   acc,
		file:      file,
		dir:       dir,
		pods:      pod.NewPod(u.client, fd, acc),
	}

	// set cookie and add user to map
	return u.Login(ui, response)
}

func (u *Users) Login(ui *Info, response http.ResponseWriter) error {
	if response != nil {
		err := cookie.SetSession(ui.GetUserName(), ui.GetSessionId(), response)
		if err != nil {
			return err
		}
	}
	u.addUserToMap(ui)

	return nil
}

func (u *Users) Logout(userName, sessionId string, response http.ResponseWriter) error {
	yes := u.isUserPresentInMap(userName)
	if !yes {
		return ErrUserNotLoggedIn
	}
	// get the user info and check if cookie id matches
	ui := u.getUserFromMap(userName)
	if ui.sessionId == sessionId {
		u.removeUserFromMap(userName)
	}
	if response != nil {
		cookie.ClearSession(response)
	}
	return nil
}

func (u *Users) IsUserLoggedIn(userName, sessionId string) bool {
	yes := u.isUserPresentInMap(userName)
	if !yes {
		return false
	}
	// get the user info and check if cookie id matches
	ui := u.getUserFromMap(userName)
	return ui.sessionId == sessionId
}

func (u *Users) GetLoggedInUserInfo(userName, sessionId string) *Info {
	ui := u.getUserFromMap(userName)
	if ui.GetSessionId() == sessionId {
		return ui
	}
	return nil
}
