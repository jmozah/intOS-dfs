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

package utils

import (
	"encoding/hex"
	"hash"
	"os"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/ethersphere/bee/pkg/swarm"
	bmtlegacy "github.com/ethersphere/bmt/legacy"
	"golang.org/x/crypto/sha3"
)

const (
	HashLength             = 32
	AddressLength          = 20
	SignatureLength        = 65
	MaxChunkLength         = 4096
	FileNameLength         = 25
	MaxDirectoryNameLength = 25
	PodNameLength          = 25
	SHA3Hash               = "SHA3"
	DefaultRoot            = "dfs_root"
	PodsInfoFile           = "dfs_pods_info"
	PathSeperator          = string(os.PathSeparator)
	MaxPodNameLength       = 25
	SpanLength             = 8
)

type Hash [HashLength]byte
type Address [AddressLength]byte

func (h Hash) Bytes() []byte        { return h[:] }
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

type decError struct{ msg string }

func (err decError) Error() string { return err.msg }

var (
	ErrEmptyString   = &decError{"empty hex string"}
	ErrMissingPrefix = &decError{"hex string without 0x prefix"}
	ErrSyntax        = &decError{"invalid hex string"}
	ErrOddLength     = &decError{"hex string of odd length"}
	ErrUint64Range   = &decError{"hex number > 64 bits"}
)

func (a Address) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	_, err := sha.Write([]byte(unchecksummed))
	if err != nil {
		return ""
	}
	sumHash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := sumHash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

func XORBytes(dst, a, b []byte) int {
	if supportsUnaligned {
		return fastXORBytes(dst, a, b)
	}
	return safeXORBytes(dst, a, b)
}

func fastXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] ^ bw[i]
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

func safeXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	enc := make([]byte, len(b)*2)
	hex.Encode(enc, b)
	return string(enc)
}

// Decode decodes a hex string with 0x prefix.
func Decode(input string) ([]byte, error) {
	if input == "" {
		return nil, ErrEmptyString
	}
	if !has0xPrefix(input) {
		return nil, ErrMissingPrefix
	}
	b, err := hex.DecodeString(input[2:])
	if err != nil {
		err = mapError(err)
	}
	return b, err
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func mapError(err error) error {
	if err, ok := err.(*strconv.NumError); ok {
		switch err.Err {
		case strconv.ErrRange:
			return ErrUint64Range
		case strconv.ErrSyntax:
			return ErrSyntax
		}
	}
	if _, ok := err.(hex.InvalidByteError); ok {
		return ErrSyntax
	}
	if err == hex.ErrLength {
		return ErrOddLength
	}
	return err
}

func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}
	copy(h[HashLength-len(b):], b)
}

func hashFunc() hash.Hash {
	return sha3.NewLegacyKeccak256()
}

func HashString(path string) []byte {
	p := bmtlegacy.NewTreePool(hashFunc, swarm.Branches, bmtlegacy.PoolSize)
	hasher := bmtlegacy.New(p)
	hasher.Reset()
	_, err := hasher.Write([]byte(path))
	if err != nil {
		return []byte{0}
	}
	return hasher.Sum(nil)
}
