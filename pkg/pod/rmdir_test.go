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
	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/jmozah/intOS-dfs/pkg/blockstore/bee/mock"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func TestPod_RemoveDir(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "pod")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	mockClient := mock.NewMockBeeClient()
	acc := account.New("user1", tempDir)
	err = acc.CreateUserAccount("password")
	if err != nil {
		t.Fatal(err)
	}
	fd := feed.New(acc.GetAccountInfo(account.UserAccountIndex), mockClient)
	pod1 := NewPod(mockClient, fd, acc)

	podName1 := "test1"
	firstDir := "dir1"
	secondDir := "dir2"
	thirdAndFourthDir := "dir3/dir4"
	fifthDir := "/dir5"
	t.Run("rmdir-on-root-of-pod", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		err = pod1.MakeDir(podName1, firstDir)
		if err != nil {
			t.Fatalf("error creating directory %s", firstDir)
		}
		dirPath := utils.PathSeperator + podName1 + utils.PathSeperator + firstDir
		dirInode := info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}

		err = pod1.RemoveDir(podName1, firstDir)
		if err != nil {
			t.Fatalf("error removing directory")
		}
		dirPath = utils.PathSeperator + podName1 + utils.PathSeperator + firstDir
		dirInode = info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode != nil {
			t.Fatalf("directory not removed")
		}
		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("rmdir-second-dir-from-first-dir", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		err = pod1.MakeDir(podName1, firstDir)
		if err != nil {
			t.Fatalf("error creating directory %s", firstDir)
		}
		_, err = pod1.ChangeDir(podName1, firstDir)
		if err != nil {
			t.Fatalf("error changing directory")
		}
		err = pod1.MakeDir(podName1, secondDir)
		if err != nil {
			t.Fatalf("error creating directory %s", secondDir)
		}
		dirPath := utils.PathSeperator + podName1 + utils.PathSeperator + firstDir + utils.PathSeperator + secondDir
		dirInode := info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}

		err = pod1.RemoveDir(podName1, secondDir)
		if err != nil {
			t.Fatalf("error removing directory")
		}
		dirPath = utils.PathSeperator + podName1 + utils.PathSeperator + firstDir + utils.PathSeperator + secondDir
		dirInode = info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode != nil {
			t.Fatalf("directory not removed")
		}

		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("rmdir-second-dir-from-pod", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		err = pod1.MakeDir(podName1, firstDir)
		if err != nil {
			t.Fatalf("error creating directory %s", err)
		}
		time.Sleep(1 * time.Second)
		err = pod1.MakeDir(podName1, firstDir+utils.PathSeperator+secondDir)
		if err != nil {
			t.Fatalf("error creating directory %s", err)
		}
		dirPath := utils.PathSeperator + podName1 + utils.PathSeperator + firstDir + utils.PathSeperator + secondDir
		dirInode := info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}

		err = pod1.RemoveDir(podName1, firstDir+utils.PathSeperator+secondDir)
		if err != nil {
			t.Fatalf("error removing directory")
		}
		dirInode = info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode != nil {
			t.Fatalf("directory not removed")
		}

		dirPath = utils.PathSeperator + podName1 + utils.PathSeperator + firstDir
		dirInode = info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory deleted")
		}

		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("rmdir-multiple-dirs-from-pod", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		err = pod1.MakeDir(podName1, thirdAndFourthDir)
		if err != nil {
			t.Fatalf("error creating directory %s", thirdAndFourthDir)
		}

		// check /test/dir3
		dirPath := utils.PathSeperator + podName1 + utils.PathSeperator + "dir3"
		dirInode := info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}
		// check /test/dir3/dir4
		dirPath = utils.PathSeperator + podName1 + utils.PathSeperator + thirdAndFourthDir
		dirInode = info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}

		err = pod1.RemoveDir(podName1, "dir3")
		if err != nil {
			t.Fatalf("error removing directory")
		}
		dirPath = utils.PathSeperator + podName1 + utils.PathSeperator + "dir3"
		dirInode = info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode != nil {
			t.Fatalf("directory not removed")
		}
		dirPath = utils.PathSeperator + podName1 + utils.PathSeperator + thirdAndFourthDir
		dirInode = info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode != nil {
			t.Fatalf("directory not removed")
		}

		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("rmdir-with-slash-on-pod", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		err = pod1.MakeDir(podName1, fifthDir)
		if err != nil {
			t.Fatalf("error creating directory %s", fifthDir)
		}
		dirPath := utils.PathSeperator + podName1 + fifthDir
		dirInode := info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}

		err = pod1.RemoveDir(podName1, fifthDir)
		if err != nil {
			t.Fatalf("error removing directory")
		}
		dirPath = utils.PathSeperator + podName1 + fifthDir
		dirInode = info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode != nil {
			t.Fatalf("directory not deleted")
		}

		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})
}
