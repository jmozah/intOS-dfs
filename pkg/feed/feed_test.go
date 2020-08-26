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

package feed

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/blockstore/bee/mock"
)

func TestFeed(t *testing.T) {

	tempDir, err := ioutil.TempDir("", "feed_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	acc := account.New("feed_pod", tempDir)
	_, err = acc.CreateUserAccount("password")
	if err != nil {
		t.Fatal(err)
	}
	user := acc.GetAddress(account.UserAccountIndex)
	accountInfo := acc.GetAccountInfo(account.UserAccountIndex)
	client := mock.NewMockBeeClient()
	//client := bee.NewBeeClient("127.0.0.1", "8080")

	t.Run("create-feed", func(t *testing.T) {
		fd := New(accountInfo, client)
		topic := hashString("topic1")
		data := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		addr, err := fd.CreateFeed(topic, user, data)
		if err != nil {
			t.Fatal(err)
		}

		// check if the data and address is present and is same as stored
		rcvdAddr, rcvdData, err := fd.GetFeedData(topic, user)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(addr, rcvdAddr) {
			t.Fatal(err)
		}
		if !bytes.Equal(data, rcvdData) {
			t.Fatal(err)
		}
	})

	t.Run("read-feed-first-time", func(t *testing.T) {
		fd := New(accountInfo, client)
		topic := hashString("topic2")

		// check if the data and address is present and is same as stored
		_, _, err := fd.GetFeedData(topic, user)
		if err != nil && err.Error() != "no feed updates found" {
			t.Fatal(err)
		}

	})

	t.Run("update-feed", func(t *testing.T) {
		fd := New(accountInfo, client)
		topic := hashString("topic3")
		data := []byte{0}
		_, err := fd.CreateFeed(topic, user, data)
		if err != nil {
			t.Fatal(err)
		}

		for i := 1; i < 256; i++ {
			buf := make([]byte, 4)
			binary.LittleEndian.PutUint16(buf, uint16(i))
			_, err := fd.UpdateFeed(topic, user, buf)
			if err != nil {
				t.Fatal(err)
			}
			getAddr, rcvdData, err := fd.GetFeedData(topic, user)
			if err != nil {
				t.Fatal(err)
			}
			if getAddr == nil {
				t.Fatal("invalid update address")
			}
			if !bytes.Equal(buf, rcvdData) {
				t.Fatal(err)
			}
			fmt.Println("update ", i, " Done")
		}

	})

}
