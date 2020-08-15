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
	"crypto/rand"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/feed"

	"github.com/jmozah/intOS-dfs/pkg/blockstore/bee/mock"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func TestPod_CopyFromLocal(t *testing.T) {
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
	t.Run("copy-file-to-root-of-pod", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		localFile, clean := createRandomFile(t, 540)
		defer clean()

		err = pod1.CopyFromLocal(podName1, localFile, info.GetCurrentPodPathAndName(), "100")
		if err != nil {
			t.Fatalf("copyFromlocal failed: %s", err.Error())
		}

		if !info.getFile().IsFileAlreadyPResent(info.GetCurrentPodPathAndName() + utils.PathSeperator + filepath.Base(localFile)) {
			t.Fatalf("file not copied in pod")
		}

		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("copy-file-to-root-of-pod-with-dot", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		localFile, clean := createRandomFile(t, 540)
		defer clean()

		err = pod1.CopyFromLocal(podName1, localFile, ".", "100")
		if err != nil {
			t.Fatalf("copyFromlocal failed: %s", err.Error())
		}

		if !info.getFile().IsFileAlreadyPResent(info.GetCurrentPodPathAndName() + utils.PathSeperator + filepath.Base(localFile)) {
			t.Fatalf("file not copied in pod")
		}

		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("copy-file-to-first-dir-from-root", func(t *testing.T) {
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

		localFile, clean := createRandomFile(t, 540)
		defer clean()

		podDir := utils.PathSeperator + firstDir
		err = pod1.CopyFromLocal(podName1, localFile, podDir, "100")
		if err != nil {
			t.Fatalf("copyFromlocal failed")
		}

		if !info.getFile().IsFileAlreadyPResent(info.GetCurrentPodPathAndName() + podDir + utils.PathSeperator + filepath.Base(localFile)) {
			t.Fatalf("file not copied in pod")
		}

		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("copy-file-to-first-dir-from-firstdir", func(t *testing.T) {
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

		_, err = pod1.ChangeDir(podName1, firstDir)
		if err != nil {
			t.Fatalf("error changing directory")
		}

		localFile, clean := createRandomFile(t, 540)
		defer clean()

		podDir := utils.PathSeperator + firstDir
		err = pod1.CopyFromLocal(podName1, localFile, podDir, "100")
		if err != nil {
			t.Fatalf("copyFromlocal failed")
		}

		if !info.getFile().IsFileAlreadyPResent(info.GetCurrentPodPathAndName() + podDir + utils.PathSeperator + filepath.Base(localFile)) {
			t.Fatalf("file not copied in pod")
		}

		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

}

func createRandomFile(t *testing.T, size int) (string, func()) {
	file, err := ioutil.TempFile("/tmp", "intos")
	if err != nil {
		t.Fatal(err)
	}
	bytes := make([]byte, size)
	_, err = rand.Read(bytes)
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write(bytes)
	if err != nil {
		t.Fatal(err)
	}
	clean := func() { os.Remove(file.Name()) }
	fileName := file.Name()
	return fileName, clean
}
