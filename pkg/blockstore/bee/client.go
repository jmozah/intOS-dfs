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

package bee

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ethersphere/bee/pkg/swarm"
	bmtlegacy "github.com/ethersphere/bmt/legacy"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"

	"github.com/jmozah/intOS-dfs/pkg/logging"
)

const (
	MaxIdleConnections     int = 20
	RequestTimeout         int = 1
	LRUSize                    = 1024
	ChunkUploadDownloadUrl     = "/chunks/"
	BytesUploadDownloadUrl     = "/bytes"
)

type BeeClient struct {
	host   string
	port   string
	url    string
	client *http.Client
	hasher *bmtlegacy.Hasher
	cache  *lru.Cache
	logger logging.Logger
}

func hashFunc() hash.Hash {
	return sha3.NewLegacyKeccak256()
}

type bytesPostResponse struct {
	Reference swarm.Address `json:"reference"`
}

func NewBeeClient(host, port string, logger logging.Logger) *BeeClient {
	p := bmtlegacy.NewTreePool(hashFunc, swarm.Branches, bmtlegacy.PoolSize)
	cache, err := lru.New(LRUSize)
	if err != nil {
		fmt.Println("could not initialise cache. system will be slow")
	}
	return &BeeClient{
		host:   host,
		port:   port,
		url:    fmt.Sprintf("http://" + host + ":" + port),
		client: createHTTPClient(),
		hasher: bmtlegacy.New(p),
		cache:  cache,
		logger: logger,
	}
}

func (s *BeeClient) CheckConnection() bool {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		return false
	}

	response, err := s.client.Do(req)
	if err != nil {
		return false
	}

	if response.StatusCode != http.StatusOK {
		return false
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}
	if string(data) != "Ethereum Swarm Bee" {
		return false
	}
	return true
}

// upload a chunk in bee
func (s *BeeClient) UploadChunk(ch swarm.Chunk) (address []byte, err error) {
	to := time.Now()
	path := filepath.Join(ChunkUploadDownloadUrl, ch.Address().String())
	fullUrl := fmt.Sprintf(s.url + path)
	req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewBuffer(ch.Data()))
	if err != nil {
		return nil, err
	}

	response, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error uploading data")
	}

	if s.inCache(ch.Address().String()) {
		s.addToCache(ch.Address().String(), ch.Data())
	}
	fields := logrus.Fields{
		"reference": ch.Address().String(),
		"duration":  time.Since(to).String(),
	}
	s.logger.WithFields(fields).Log(logrus.DebugLevel, "upload chunk: ")
	return ch.Address().Bytes(), nil
}

// download a chunk from bee
func (s *BeeClient) DownloadChunk(ctx context.Context, address []byte) (data []byte, err error) {
	to := time.Now()
	addrString := swarm.NewAddress(address).String()
	if s.inCache(addrString) {
		return s.getFromCache(swarm.NewAddress(address).String()), nil
	}

	path := filepath.Join(ChunkUploadDownloadUrl, addrString)
	fullUrl := fmt.Sprintf(s.url + path)
	req, err := http.NewRequest(http.MethodGet, fullUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	response, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error downloading data")
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error downloading data")
	}

	s.addToCache(addrString, data)
	fields := logrus.Fields{
		"reference": addrString,
		"duration":  time.Since(to).String(),
	}
	s.logger.WithFields(fields).Log(logrus.DebugLevel, "download chunk: ")
	return data, nil
}

// upload a chunk in bee
func (s *BeeClient) UploadBlob(data []byte) (address []byte, err error) {
	to := time.Now()
	fullUrl := fmt.Sprintf(s.url + BytesUploadDownloadUrl)
	req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	response, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error uploading blob")
	}

	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error uploading blob")
	}

	var resp bytesPostResponse
	err = json.Unmarshal(respData, &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response")
	}
	fields := logrus.Fields{
		"reference": resp.Reference.String(),
		"size":      len(data),
		"duration":  time.Since(to).String(),
	}
	s.logger.WithFields(fields).Log(logrus.DebugLevel, "upload blob: ")
	return resp.Reference.Bytes(), nil
}

func (s *BeeClient) DownloadBlob(address []byte) ([]byte, int, error) {
	to := time.Now()
	addrString := swarm.NewAddress(address).String()
	if s.inCache(addrString) {
		return s.getFromCache(addrString), 200, nil
	}

	fullUrl := fmt.Sprintf(s.url + BytesUploadDownloadUrl + "/" + swarm.NewAddress(address).String())
	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	response, err := s.client.Do(req)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, errors.New("error downloading blob ")
	}

	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, errors.New("error downloading blob")
	}
	fields := logrus.Fields{
		"reference": addrString,
		"size":      len(respData),
		"duration":  time.Since(to).String(),
	}
	s.logger.WithFields(fields).Log(logrus.DebugLevel, "download blob: ")
	return respData, response.StatusCode, nil
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		//Timeout: time.Duration(RequestTimeout) * time.Second,
	}
	return client
}

func (s *BeeClient) addToCache(key string, value []byte) {
	if s.cache != nil {
		s.cache.Add(key, value)
	}
}

func (s *BeeClient) inCache(key string) bool {
	if s.cache != nil {
		return s.cache.Contains(key)
	}
	return false
}

func (s *BeeClient) getFromCache(key string) []byte {
	if s.cache != nil {
		value, ok := s.cache.Get(key)
		if ok {
			return value.([]byte)
		}
		return nil
	}
	return nil
}
