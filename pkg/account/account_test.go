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

package account

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func TestAccount_CreateRootAccount(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "pod")
	if err != nil {
		t.Fatal(err)
	}

	podName := "test_pod1"
	password := "letmein"
	acc := New(podName, tempDir)
	_, err = acc.CreateUserAccount(password)
	if err != nil {
		t.Fatal(err)
	}

	if acc.wallet == nil || acc.wallet.encryptedmnemonic == "" {
		t.Fatal("wallet creation error")
	}

	plainMnemonic, err := acc.wallet.decryptMnemonic(password)
	if err != nil {
		t.Fatal(err)
	}

	words := strings.Split(plainMnemonic, " ")
	if len(words) != 12 {
		t.Fatal("mnemonic is not 12 words")
	}

	if !acc.IsAlreadyInitialized() {
		t.Fatal("mnemonic file does not exists")
	}

	if acc.userAcount.GetPrivateKey() == nil || acc.userAcount.GetPublicKey() == nil || len(acc.userAcount.address[:]) != utils.AddressLength {
		t.Fatalf("keys not intialised")
	}

	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadAndStoreMnemonic(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "pod")
	if err != nil {
		t.Fatal(err)
	}
	podName := "test_pod1"
	password := "letmein"
	acc := New(podName, tempDir)
	_, err = acc.CreateUserAccount(password)
	if err != nil {
		t.Fatal(err)
	}

	expectedMnemonic, err := acc.wallet.decryptMnemonic(password)
	if err != nil {
		t.Fatal(err)
	}

	acc.wallet.encryptedmnemonic = ""
	err = acc.LoadEncryptedMnemonicFromDisk(password)
	if err != nil {
		t.Fatal(err)
	}

	gotMnemonic, err := acc.wallet.decryptMnemonic(password)
	if err != nil {
		t.Fatal(err)
	}

	if gotMnemonic != expectedMnemonic {
		t.Fatalf("mnemonics does not match. expected %s and got %s", expectedMnemonic, gotMnemonic)
	}

	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Fatal(err)
	}
}
