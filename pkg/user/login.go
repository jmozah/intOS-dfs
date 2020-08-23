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

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	f "github.com/jmozah/intOS-dfs/pkg/file"
	"github.com/jmozah/intOS-dfs/pkg/pod"
)

func (u *Users) LoginUser(userName string, passPhrase string, dataDir string, client blockstore.Client) error {
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

	ui := &Info{
		name:    userName,
		feedApi: fd,
		account: acc,
		file:    file,
		dir:     dir,
		pods:    pod.NewPod(u.client, fd, acc),
	}
	u.addUserToMap(ui)
	return nil
}

func (u *Users) IsUserLoggedIn(userName string) bool {
	return u.isUserPresentInMap(userName)
}

func (u *Users) GetLoggedInUserInfo(userName string) *Info {
	return u.getUserFromMap(userName)
}
