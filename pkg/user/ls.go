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
	"github.com/jmozah/intOS-dfs/pkg/account"
	"os"
	"path/filepath"
	"strings"
)

func (u *Users) ListAllUsers(dataDir string) []string {
	var users []string
	keyFileDir := account.GetKeyFileDir(dataDir)
	err := filepath.Walk(keyFileDir,
		func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".key") {
				userName := strings.TrimSuffix(info.Name(), ".key")
				userName = "<User> " + userName
				users = append(users, userName)
			}
			return nil
		})
	if err != nil {
		return nil
	}
	return users
}
