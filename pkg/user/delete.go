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
	"os"

	"github.com/jmozah/intOS-dfs/pkg/account"
)

func (u *Users) DeleteUser(userName, dataDir, password string) error {

	if !u.IsUsernameAvailable(userName, dataDir) {
		return ErrInvalidUserName
	}

	if !u.isUserPresentInMap(userName) {
		return ErrUserNotLoggedIn
	}

	// check for valid password
	userInfo := u.getUserFromMap(userName)
	acc := userInfo.account
	if !acc.Authorise(password) {
		return ErrInvalidPassword
	}

	// Logout user
	if u.IsUserLoggedIn(userName) {
		u.removeUserFromMap(userName)
	}

	// remove the user mnemonic file
	userKeyFileName := account.ConstructUserKeyFile(userName, dataDir)
	err := os.Remove(userKeyFileName)
	if err != nil {
		return fmt.Errorf("user del: could not remove user key: %w", err)
	}
	return nil
}
