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
	"strings"
	"testing"

	"github.com/jmozah/intOS-dfs/pkg/blockstore/bee/mock"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func TestPod_MakeDir(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "pod")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	mockClient := mock.NewMockBeeClient()
	pod1 := NewPod(mockClient)
	err = pod1.LoadRootPod(tempDir, "password")
	if err != nil {
		t.Fatal(err)
	}
	podName1 := "test1"
	firstDir := "dir1"
	secondDir := "dir2"
	thirdDir := "dir3/dir4"
	fourthDir := "/dir5"

	t.Run("mkdir-on-root-of-pod", func(t *testing.T) {
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
		if dirInode.Meta.Path != utils.PathSeperator+podName1 {
			t.Fatalf("invalid path in meta")
		}
		if dirInode.Meta.Name != firstDir {
			t.Fatalf("invalid name in meta")
		}
		if dirInode.GetDirInodePathAndName() != dirPath {
			t.Fatalf("invalid path or name")
		}
		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("mkdir-second-dir-from-first-dir", func(t *testing.T) {
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
		if dirInode.Meta.Path != utils.PathSeperator+podName1+utils.PathSeperator+firstDir {
			t.Fatalf("invalid path in meta")
		}
		if dirInode.Meta.Name != secondDir {
			t.Fatalf("invalid name in meta")
		}
		if dirInode.GetDirInodePathAndName() != dirPath {
			t.Fatalf("invalid path or name")
		}
		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("mkdir-second-dir-from-pod", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		err = pod1.MakeDir(podName1, firstDir)
		if err != nil {
			t.Fatalf("error creating directory %s", err)
		}
		err = pod1.MakeDir(podName1, firstDir+utils.PathSeperator+secondDir)
		if err != nil {
			t.Fatalf("error creating directory %s", err)
		}

		dirPath := utils.PathSeperator + podName1 + utils.PathSeperator + firstDir + utils.PathSeperator + secondDir
		dirInode := info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}
		if dirInode.Meta.Path != utils.PathSeperator+podName1+utils.PathSeperator+"dir1" {
			t.Fatalf("invalid path in meta")
		}
		if dirInode.Meta.Name != "dir2" {
			t.Fatalf("invalid name in meta")
		}
		if dirInode.GetDirInodePathAndName() != dirPath {
			t.Fatalf("invalid path or name")
		}
		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("mkdir-multiple-dirs-from-pod", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		err = pod1.MakeDir(podName1, thirdDir)
		if err != nil {
			t.Fatalf("error creating directory %s", thirdDir)
		}

		// check /test/dir3
		dirPath := utils.PathSeperator + podName1 + utils.PathSeperator + "dir3"
		dirInode := info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}
		if dirInode.Meta.Path != utils.PathSeperator+podName1 {
			t.Fatalf("invalid path in meta")
		}
		if dirInode.Meta.Name != "dir3" {
			t.Fatalf("invalid name in meta")
		}
		if dirInode.GetDirInodePathAndName() != dirPath {
			t.Fatalf("invalid path or name")
		}

		// check /test/dir3/dir4
		dirPath = utils.PathSeperator + podName1 + utils.PathSeperator + thirdDir
		dirInode = info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}
		if dirInode.Meta.Path != utils.PathSeperator+podName1+utils.PathSeperator+"dir3" {
			t.Fatalf("invalid path in meta")
		}
		if dirInode.Meta.Name != "dir4" {
			t.Fatalf("invalid name in meta")
		}
		if dirInode.GetDirInodePathAndName() != dirPath {
			t.Fatalf("invalid path or name")
		}
		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})

	t.Run("mkdir-with-slash-on-pod", func(t *testing.T) {
		info, err := pod1.CreatePod(podName1, tempDir, "password")
		if err != nil {
			t.Fatalf("error creating pod %s", podName1)
		}

		err = pod1.MakeDir(podName1, fourthDir)
		if err != nil {
			t.Fatalf("error creating directory %s", fourthDir)
		}

		dirPath := utils.PathSeperator + podName1 + fourthDir
		dirInode := info.getDirectory().GetDirFromDirectoryMap(dirPath)
		if dirInode == nil {
			t.Fatalf("directory not created")
		}
		if dirInode.Meta.Path != utils.PathSeperator+podName1 {
			t.Fatalf("invalid path in meta")
		}
		if dirInode.Meta.Name != strings.TrimPrefix(fourthDir, "/") {
			t.Fatalf("invalid name in meta")
		}
		if dirInode.GetDirInodePathAndName() != dirPath {
			t.Fatalf("invalid path or name")
		}
		err = pod1.DeletePod(podName1, tempDir)
		if err != nil {
			t.Fatalf("could not delete pod")
		}
	})
}
