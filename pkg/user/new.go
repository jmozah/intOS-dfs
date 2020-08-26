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

func (u *Users) CreateNewUser(userName, passPhrase, mnemonic, dataDir string, client blockstore.Client, response http.ResponseWriter, sessionId string) (string, string, error) {
	if u.IsUsernameAvailable(userName, dataDir) {
		return "", "", ErrUserAlreadyPresent
	}
	acc := account.New(userName, dataDir)
	accountInfo := acc.GetAccountInfo(account.UserAccountIndex)
	fd := feed.New(accountInfo, client)
	file := f.NewFile(userName, client, fd, accountInfo)

	mnemonic, err := acc.CreateUserAccount(passPhrase, mnemonic)
	if err != nil {
		return "", "", fmt.Errorf("user create:: %w", err)
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
	err = u.Login(ui, response)
	if err != nil {
		return "", "", err
	}

	return accountInfo.GetAddress().Hex(), mnemonic, nil
}
