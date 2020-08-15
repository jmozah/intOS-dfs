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
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

const (
	rootPath    = "m/44'/60'/0'/0/0"
	genericPath = "m/44'/60'/0'/0/"
)

var (
	wallet *Wallet
	once   sync.Once
)

type Wallet struct {
	encryptedmnemonic string
}

func NewWallet(mnemonic string) *Wallet {
	once.Do(func() {
		wallet = &Wallet{
			encryptedmnemonic: mnemonic,
		}
	})
	return wallet
}

func (w *Wallet) LoadMnemonicAndCreateRootAccount() (accounts.Account, string, error) {
	// Generate a mnemonic for memorization or user-friendly seeds
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return accounts.Account{}, "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return accounts.Account{}, "", err
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return accounts.Account{}, "", err
	}
	path := hdwallet.MustParseDerivationPath(rootPath)
	acc, err := wallet.Derive(path, false)
	if err != nil {
		return accounts.Account{}, "", err
	}
	return acc, mnemonic, nil

}

func (w *Wallet) CreateAccount(walletPath string, plainMnemonic string) (accounts.Account, error) {
	wallet, err := hdwallet.NewFromMnemonic(plainMnemonic)
	if err != nil {
		return accounts.Account{}, err
	}
	path := hdwallet.MustParseDerivationPath(walletPath)
	acc, err := wallet.Derive(path, false)
	if err != nil {
		return accounts.Account{}, err
	}
	return acc, nil
}

func (w *Wallet) decryptMnemonic(password string) (string, error) {
	if w.encryptedmnemonic == "" {
		return "", fmt.Errorf("invalid encrypted mnemonic")
	}

	aesKey := sha256.Sum256([]byte(password))

	//decrypt the message
	mnemonic, err := decrypt(aesKey[:], w.encryptedmnemonic)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
