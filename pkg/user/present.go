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
	"os"
	"regexp"

	"github.com/jmozah/intOS-dfs/pkg/account"
)

func (u *Users) IsUsernameAvailable(userName string, dataDir string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	if !re.MatchString(userName) {
		return false
	}

	userKeyFileName := account.ConstructUserKeyFile(userName, dataDir)
	info, err := os.Stat(userKeyFileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
