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
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ethersphere/bee/pkg/crypto"
	"github.com/jmozah/intOS-dfs/pkg/utils"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"golang.org/x/crypto/ssh/terminal"
)

type Account struct {
	dataDir          string
	podName          string
	mnemonicFileName string
	wallet           *Wallet
	privateKey       *ecdsa.PrivateKey
	publicKey        *ecdsa.PublicKey
	address          utils.Address
}

const (
	KeyStoreDirectoryName = "keystore"
)

func New(podName, datadir string) *Account {
	destDir := filepath.Join(datadir, KeyStoreDirectoryName)
	destFile := filepath.Join(destDir, utils.DefaultRoot+".key")

	wallet := NewWallet("")
	return &Account{
		dataDir:          datadir,
		podName:          podName,
		wallet:           wallet,
		mnemonicFileName: destFile,
	}
}

func (a *Account) IsAlreadyInitialized() bool {
	info, err := os.Stat(a.mnemonicFileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (a *Account) CreateRootAccount(passPhrase string) error {
	if a.IsAlreadyInitialized() {
		var s string
		fmt.Println("dfs is already initialised")
		fmt.Println("reinitialising again will make all the your data inaccessible")
		fmt.Printf("do you still want to proceed (Y/N):")
		_, err := fmt.Scan(&s)
		if err != nil {
			return err
		}

		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		s = strings.Trim(s, "\n")

		if s == "n" || s == "no" {
			return nil
		}
		err = os.Remove(a.mnemonicFileName)
		if err != nil {
			return fmt.Errorf("could not remove root key: %w", err)
		}
	}

	wallet := NewWallet("")
	a.wallet = wallet
	acc, mnemonic, err := wallet.LoadMnemonicAndCreateRootAccount()
	if err != nil {
		return err
	}

	if passPhrase == "" {
		fmt.Println("Please store the following 24 words safely")
		fmt.Println("if can use this to import the wallet in another machine")
		fmt.Println("=============== Mnemonic ==========================")
		fmt.Println(mnemonic)
		fmt.Println("=============== Mnemonic ==========================")
	}

	hdw, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return err
	}

	// store publicKey, private key and user
	a.privateKey, err = hdw.PrivateKey(acc)
	if err != nil {
		return err
	}
	a.publicKey, err = hdw.PublicKey(acc)
	if err != nil {
		return err
	}
	addrBytes, err := crypto.NewEthereumAddress(a.privateKey.PublicKey)
	if err != nil {
		return err
	}
	a.address.SetBytes(addrBytes)

	// store the mnemonic
	encryptedMnemonic, err := a.storeAsEncryptedMnemonicToDisk(mnemonic, passPhrase)
	if err != nil {
		return err
	}
	a.wallet.encryptedmnemonic = encryptedMnemonic
	return nil
}

func (a *Account) LoadRootAccount(passPhrase string) error {
	password := passPhrase
	if password == "" {
		fmt.Print("Enter root password to unlock root account: ")
		password = a.getPassword()
	}

	err := a.LoadEncryptedMnemonicFromDisk(password)
	if err != nil {
		return nil
	}

	plainMnemonic, err := a.wallet.decryptMnemonic(password)
	if err != nil {
		return err
	}

	acc, err := a.wallet.CreateAccount(rootPath, plainMnemonic)
	if err != nil {
		return err
	}

	hdw, err := hdwallet.NewFromMnemonic(plainMnemonic)
	if err != nil {
		return err
	}
	a.privateKey, err = hdw.PrivateKey(acc)
	if err != nil {
		return err
	}
	a.publicKey, err = hdw.PublicKey(acc)
	if err != nil {
		return err
	}
	addrBytes, err := crypto.NewEthereumAddress(a.privateKey.PublicKey)
	if err != nil {
		return err
	}
	a.address.SetBytes(addrBytes)
	return nil
}

func (a *Account) CreateNormalAccount(accountId int, passPhrase string) error {
	if !a.IsAlreadyInitialized() {
		return fmt.Errorf("dfs not initalised. use the \"init\" command to intialise the system")
	}

	password := passPhrase
	if password == "" {
		fmt.Print("Enter root password to create a pod: ")
		password = a.getPassword()
	}

	plainMnemonic, err := a.wallet.decryptMnemonic(password)
	if err != nil {
		return err
	}

	path := genericPath + strconv.Itoa(accountId)
	acc, err := a.wallet.CreateAccount(path, plainMnemonic)
	if err != nil {
		return err
	}
	hdw, err := hdwallet.NewFromMnemonic(plainMnemonic)
	if err != nil {
		return err
	}

	a.privateKey, err = hdw.PrivateKey(acc)
	if err != nil {
		return err
	}
	a.publicKey, err = hdw.PublicKey(acc)
	if err != nil {
		return err
	}
	addrBytes, err := crypto.NewEthereumAddress(a.privateKey.PublicKey)
	if err != nil {
		return err
	}
	a.address.SetBytes(addrBytes)
	return nil
}

func (a *Account) LoadEncryptedMnemonicFromDisk(passPhrase string) error {
	if !a.IsAlreadyInitialized() {
		return fmt.Errorf("dfs not initalised. use the \"init\" command to intialise the system")
	}

	encryptedMessage, err := ioutil.ReadFile(a.mnemonicFileName)
	if err != nil {
		return nil
	}

	a.wallet.encryptedmnemonic = string(encryptedMessage)
	return nil
}

func (a *Account) storeAsEncryptedMnemonicToDisk(mnemonic string, passPhrase string) (string, error) {
	if a.IsAlreadyInitialized() {
		err := os.Remove(a.mnemonicFileName)
		if err != nil {
			return "", fmt.Errorf("could not remove old key file: %w", err)
		}
	}

	// get the password and hash it to 256 bits
	password := passPhrase
	if password == "" {
		fmt.Print("Enter root password to unlock root account: ")
		password = a.getPassword()
		password = strings.Trim(password, "\n")
	}
	aesKey := sha256.Sum256([]byte(password))

	// encrypt the mnemonic
	encryptedMessage, err := encrypt(aesKey[:], mnemonic)
	if err != nil {
		return "", fmt.Errorf("create root account: %w", err)
	}

	err = os.MkdirAll(filepath.Dir(a.mnemonicFileName), 0777)
	if err != nil {
		return "", err
	}

	// store the mnemonic in a file
	f, err := os.Create(a.mnemonicFileName)
	if err != nil {
		return "", err
	}
	n, err := f.WriteString(encryptedMessage)
	if err != nil {
		return "", err
	}
	if n != len(encryptedMessage) {
		return "", fmt.Errorf("file write error during encryption")
	}
	err = f.Sync()
	if err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}
	return encryptedMessage, nil
}

func (a *Account) GetPrivateKey() *ecdsa.PrivateKey {
	return a.privateKey
}

func (s *Account) GetAddress() utils.Address {
	return s.address
}

func (a *Account) getPassword() (password string) {
	// read the pass phrase
	bytePassword, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalf("error reading password")
		return
	}
	fmt.Println("")
	passwd := string(bytePassword)
	password = strings.TrimSpace(passwd)
	return password
}

// list accounts with balances

// withdraw eth from account

// import account

// export account
