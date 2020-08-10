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

package mock

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"sync"

	"github.com/ethersphere/bee/pkg/swarm"
)

type MockBeeClient struct {
	storer   map[string][]byte
	storerMu sync.RWMutex
}

func NewMockBeeClient() *MockBeeClient {
	return &MockBeeClient{
		storer:   make(map[string][]byte),
		storerMu: sync.RWMutex{},
	}
}

func (m *MockBeeClient) UploadChunk(ch swarm.Chunk) (address []byte, err error) {
	m.storerMu.Lock()
	defer m.storerMu.Unlock()
	m.storer[ch.Address().String()] = ch.Data()
	fmt.Println("uploading ", ch.Address().Bytes())
	return ch.Address().Bytes(), nil
}

func (m *MockBeeClient) DownloadChunk(ctx context.Context, address []byte) (data []byte, err error) {
	m.storerMu.Lock()
	defer m.storerMu.Unlock()
	if data, ok := m.storer[swarm.NewAddress(address).String()]; ok {
		return data, nil
	}
	fmt.Println("downloading ", address)
	return nil, fmt.Errorf("error downloading data")
}

func (m *MockBeeClient) UploadBlob(data []byte) (address []byte, err error) {
	m.storerMu.Lock()
	defer m.storerMu.Unlock()
	address = make([]byte, 32)
	_, err = rand.Read(address)
	m.storer[swarm.NewAddress(address).String()] = data
	return address, nil
}

func (m *MockBeeClient) DownloadBlob(address []byte) (data []byte, respCode int, err error) {
	m.storerMu.Lock()
	defer m.storerMu.Unlock()
	if data, ok := m.storer[swarm.NewAddress(address).String()]; ok {
		return data, http.StatusOK, nil
	}
	return nil, http.StatusInternalServerError, fmt.Errorf("error downloading data")
}
