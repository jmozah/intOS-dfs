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

package pod

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/jmozah/intOS-dfs/pkg/account"
	d "github.com/jmozah/intOS-dfs/pkg/dir"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	f "github.com/jmozah/intOS-dfs/pkg/file"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (p *Pod) LoadRootPod(dataDir string, passPhrase string) error {
	acc := account.New(utils.DefaultRoot, dataDir)
	fd := feed.New(acc, p.client)
	file := f.NewFile(utils.DefaultRoot, p.client, fd, acc)
	var dirInode *d.DirInode
	if acc.IsAlreadyInitialized() {
		err := acc.LoadRootAccount(passPhrase)
		if err != nil {
			return err
		}

		dir := d.NewDirectory(utils.DefaultRoot, p.client, fd, acc, file)
		_, dirInode, err = dir.GetDirNode(utils.DefaultRoot, fd, acc)
		if err != nil {
			if errors.Is(err, syscall.ECONNREFUSED) {
				return fmt.Errorf("connection refused")
			}
			return fmt.Errorf("root pod: %w", err)
		}
	} else {
		err := acc.CreateRootAccount(passPhrase)
		if err != nil {
			return err
		}

		dir := d.NewDirectory(utils.DefaultRoot, p.client, fd, acc, file)
		dirInode, _, err = dir.CreateDirINode(utils.DefaultRoot, "", nil)
		if err != nil {
			if errors.Is(err, syscall.ECONNREFUSED) {
				return fmt.Errorf("connection refused")
			}
			return fmt.Errorf("root pod: %w", err)
		}
	}

	p.rootFeed = fd
	p.rootDirInode = dirInode
	p.rootAccount = acc

	return nil
}

func (p *Pod) IsInitialized(dataDir string) bool {
	keyStore := filepath.Join(dataDir, account.KeyStoreDirectoryName)
	rootKeyName := filepath.Join(keyStore, utils.DefaultRoot+".key")

	fi, err := os.Stat(rootKeyName)
	if err == nil && !fi.IsDir() {
		return true
	}
	return false
}

func (p *Pod) RemoveRootKey(dataDir string) error {
	keyStore := filepath.Join(dataDir, account.KeyStoreDirectoryName)
	rootKeyName := filepath.Join(keyStore, utils.DefaultRoot+".key")
	err := os.Remove(rootKeyName)
	if err != nil {
		return err
	}
	p.rootFeed = nil
	p.rootDirInode = nil
	p.rootAccount = nil
	return nil
}
