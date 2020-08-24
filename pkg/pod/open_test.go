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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/feed"

	"github.com/jmozah/intOS-dfs/pkg/blockstore/bee/mock"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func TestPod_LoginPod(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "pod")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	mockClient := mock.NewMockBeeClient()
	acc := account.New("user1", tempDir)
	_, err = acc.CreateUserAccount("password", "")
	if err != nil {
		t.Fatal(err)
	}
	fd := feed.New(acc.GetAccountInfo(account.UserAccountIndex), mockClient)
	pod1 := NewPod(mockClient, fd, acc)

	podName1 := "test1"
	firstDir := "dir1"
	t.Run("simple-login-to-pod", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}
		err = pod1.ClosePod(podName1)
		if err != nil {
			t.Fatalf("could not logout")
		}

		infoLogin, err := pod1.OpenPod(podName1, "password")
		if err != nil {
			t.Fatalf("login failed")
		}
		if info.podName != infoLogin.podName {
			t.Fatalf("invalid podname")
		}
		if info.GetCurrentPodPathAndName() != infoLogin.GetCurrentPodPathAndName() {
			t.Fatalf("invalid podname path and name")
		}

		err = pod1.DeletePod(podName1)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("login-with-sync-contents", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		//Make a dir
		err = pod1.MakeDir(podName1, firstDir)
		if err != nil {
			t.Fatalf("error creating directory %s", firstDir)
		}

		dirPath := utils.PathSeperator + podName1 + utils.PathSeperator + firstDir
		dirInode := info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}

		// create a file
		localFile, clean := createRandomFile(t, 540)
		defer clean()
		podDir := utils.PathSeperator + firstDir
		err = pod1.CopyFromLocal(podName1, localFile, podDir, "100")
		if err != nil {
			t.Fatalf("copyFromlocal failed: %s", err.Error())
		}
		if !info.getFile().IsFileAlreadyPResent(dirPath + utils.PathSeperator + filepath.Base(localFile)) {
			t.Fatalf("file not copied in pod")
		}

		err = pod1.ClosePod(podName1)
		if err != nil {
			t.Fatalf("could not logout")
		}

		// Now login and check if the dir and file exists
		infoLogin, err := pod1.OpenPod(podName1, "password")
		if err != nil {
			t.Fatalf("login failed")
		}
		if info.podName != infoLogin.podName {
			t.Fatalf("invalid podname")
		}
		if info.GetCurrentPodPathAndName() != infoLogin.GetCurrentPodPathAndName() {
			t.Fatalf("invalid podname path and name")
		}
		dirInodeLogin := infoLogin.dir.GetDirFromDirectoryMap(dirPath)
		if dirInodeLogin == nil {
			t.Fatalf("dir not synced")
		}
		if dirInodeLogin.Meta.Path != info.GetCurrentPodPathAndName() {
			t.Fatalf("dir not synced")
		}
		if dirInodeLogin.Meta.Name != firstDir {
			t.Fatalf("dir not synced")
		}
		fileMeta := infoLogin.getFile().GetFromFileMap(dirPath + utils.PathSeperator + filepath.Base(localFile))
		if fileMeta == nil {
			t.Fatalf("file not synced")
		}

		err = pod1.DeletePod(podName1)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

}
