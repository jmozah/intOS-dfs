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
	"net/http"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	"github.com/jmozah/intOS-dfs/pkg/cookie"
	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	f "github.com/jmozah/intOS-dfs/pkg/file"
	"github.com/jmozah/intOS-dfs/pkg/pod"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (u *Users) ImportUsingAddress(userName, passPhrase, addressString, dataDir string, client blockstore.Client, response http.ResponseWriter, sessionId string) error {
	if u.IsUsernameAvailable(userName, dataDir) {
		return ErrUserAlreadyPresent
	}

	acc := account.New(u.logger)
	accountInfo := acc.GetUserAccountInfo()
	fd := feed.New(accountInfo, client, u.logger)
	file := f.NewFile(userName, client, fd, accountInfo, u.logger)

	address := utils.HexToAddress(addressString)

	// load the encrypted mnemonic and see if it is valid
	encryptedMnemonic, err := u.getEncryptedMnemonic(userName, address, fd)
	if err != nil {
		return err
	}
	err = acc.LoadUserAccount(passPhrase, encryptedMnemonic)
	if err != nil {
		if err.Error() == "mnemonic is invalid" {
			return ErrInvalidPassword
		}
		return err
	}

	// store the username -> address mapping locally
	err = u.storeUserNameToAddressFileMapping(userName, dataDir, accountInfo.GetAddress())
	if err != nil {
		return err
	}
	dir := d.NewDirectory(userName, client, fd, accountInfo, file, u.logger)

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
		pods:      pod.NewPod(u.client, fd, acc, u.logger),
	}

	// set cookie and add user to map
	err = u.Login(ui, response)
	if err != nil {
		return err
	}

	return nil
}
