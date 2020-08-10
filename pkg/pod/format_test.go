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
	"testing"

	"github.com/jmozah/intOS-dfs/pkg/blockstore/bee/mock"
)

func TestFormat(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "pod")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	mockClient := mock.NewMockBeeClient()
	pod1 := NewPod(mockClient)

	t.Run("format-root", func(t *testing.T) {
		err := pod1.LoadRootPod(tempDir, "password")
		if err != nil {
			t.Fatal(err)
		}
		if pod1.rootAccount == nil || pod1.rootDirInode == nil || pod1.rootFeed == nil {
			t.Fatal(err)
		}

		_, err = pod1.CreatePod("zz", tempDir, "password")
		if err != nil {
			t.Fatal(err)
		}

	})
}
