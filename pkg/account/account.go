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
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	UserAccountIndex = -1
)

type Account struct {
	dataDir          string
	podName          string
	mnemonicFileName string
	wallet           *Wallet
	userAcount       *AccountInfo
	podAccounts      map[int]*AccountInfo
}

type AccountInfo struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    utils.Address
}

const (
	KeyStoreDirectoryName = "keystore"
)

func New(podName, dataDir string) *Account {
	destFile := ConstructUserKeyFile(podName, dataDir)
	wallet := NewWallet("")
	return &Account{
		dataDir:          dataDir,
		podName:          podName,
		wallet:           wallet,
		mnemonicFileName: destFile,
		userAcount:       &AccountInfo{},
		podAccounts:      make(map[int]*AccountInfo),
	}
}

func (a *Account) IsAlreadyInitialized() bool {
	info, err := os.Stat(a.mnemonicFileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (a *Account) CreateUserAccount(passPhrase, mnemonic string) (string, error) {
	if passPhrase == "" {
		if a.IsAlreadyInitialized() {
			var s string
			fmt.Println("user is already initialised")
			fmt.Println("reinitialising again will make all the your data inaccessible")
			fmt.Printf("do you still want to proceed (Y/N):")
			_, err := fmt.Scan(&s)
			if err != nil {
				return "", err
			}

			s = strings.TrimSpace(s)
			s = strings.ToLower(s)
			s = strings.Trim(s, "\n")

			if s == "n" || s == "no" {
				return "", nil
			}
			err = os.Remove(a.mnemonicFileName)
			if err != nil {
				return "", fmt.Errorf("could not remove user key: %w", err)
			}
		}
	} else {
		if a.IsAlreadyInitialized() {
			return "", fmt.Errorf("user already present")
		}
	}

	wallet := NewWallet("")
	a.wallet = wallet
	acc, mnemonic, err := wallet.LoadMnemonicAndCreateRootAccount(mnemonic)
	if err != nil {
		return "", err
	}

	hdw, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	// store publicKey, private key and user
	a.userAcount.privateKey, err = hdw.PrivateKey(acc)
	if err != nil {
		return "", err
	}
	a.userAcount.publicKey, err = hdw.PublicKey(acc)
	if err != nil {
		return "", err
	}
	addrBytes, err := crypto.NewEthereumAddress(a.userAcount.privateKey.PublicKey)
	if err != nil {
		return "", err
	}
	a.userAcount.address.SetBytes(addrBytes)

	// store the mnemonic
	encryptedMnemonic, err := a.storeAsEncryptedMnemonicToDisk(mnemonic, passPhrase)
	if err != nil {
		return "", err
	}
	a.wallet.encryptedmnemonic = encryptedMnemonic

	return mnemonic, nil
}

func (a *Account) LoadUserAccount(passPhrase string) error {
	password := passPhrase
	if password == "" {
		fmt.Print("Enter password to unlock user account: ")
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
	a.userAcount.privateKey, err = hdw.PrivateKey(acc)
	if err != nil {
		return err
	}
	a.userAcount.publicKey, err = hdw.PublicKey(acc)
	if err != nil {
		return err
	}
	addrBytes, err := crypto.NewEthereumAddress(a.userAcount.privateKey.PublicKey)
	if err != nil {
		return err
	}
	a.userAcount.address.SetBytes(addrBytes)
	return nil
}

func (a *Account) Authorise(password string) bool {
	if password == "" {
		fmt.Print("Enter user password to create a pod: ")
		password = a.getPassword()
	}
	plainMnemonic, err := a.wallet.decryptMnemonic(password)
	if err != nil {
		return false
	}
	// check the validity of the mnemonic
	if plainMnemonic == "" {
		return false
	}
	words := strings.Split(plainMnemonic, " ")
	if len(words) != 12 {
		return false
	}
	if !bip39.IsMnemonicValid(plainMnemonic) {
		return false
	}
	return true
}

func (a *Account) CreatePodAccount(accountId int, passPhrase string) error {
	if !a.IsAlreadyInitialized() {
		return fmt.Errorf("user not created")
	}

	if _, ok := a.podAccounts[accountId]; ok {
		return nil
	}

	password := passPhrase
	if password == "" {
		fmt.Print("Enter user password to create a pod: ")
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

	accountInfo := &AccountInfo{}

	accountInfo.privateKey, err = hdw.PrivateKey(acc)
	if err != nil {
		return err
	}
	accountInfo.publicKey, err = hdw.PublicKey(acc)
	if err != nil {
		return err
	}
	addrBytes, err := crypto.NewEthereumAddress(accountInfo.privateKey.PublicKey)
	if err != nil {
		return err
	}
	accountInfo.address.SetBytes(addrBytes)
	a.podAccounts[accountId] = accountInfo
	return nil
}

func (a *Account) DeletePodAccount(accountId int) {
	delete(a.podAccounts, accountId)
}

func (a *Account) LoadEncryptedMnemonicFromDisk(passPhrase string) error {
	if !a.IsAlreadyInitialized() {
		return fmt.Errorf("dfs not initialised. use the \"init\" command to intialise the system")
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

func (a *Account) GetUserPrivateKey(index int) *ecdsa.PrivateKey {
	if index == UserAccountIndex {
		return a.userAcount.privateKey
	} else {
		return a.podAccounts[index].privateKey
	}
}

func (a *Account) GetAddress(index int) utils.Address {
	if index == UserAccountIndex {
		return a.userAcount.address
	} else {
		return a.podAccounts[index].address
	}
}

func (a *Account) GetAccountInfo(index int) *AccountInfo {
	if index == UserAccountIndex {
		return a.userAcount
	} else {
		return a.podAccounts[index]
	}
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

func (ai *AccountInfo) GetAddress() utils.Address {
	return ai.address
}

func (ai *AccountInfo) GetPrivateKey() *ecdsa.PrivateKey {
	return ai.privateKey
}

func (ai *AccountInfo) GetPublicKey() *ecdsa.PublicKey {
	return ai.publicKey
}

// list accounts with balances

// withdraw eth from account

// import account

// export account
