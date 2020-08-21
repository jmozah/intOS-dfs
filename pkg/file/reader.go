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

package file

import (
	"fmt"

	"github.com/jmozah/intOS-dfs/pkg/blockstore"
)

type Reader struct {
	totalBytes  int64
	offset      int64
	client      blockstore.Client
	fileInode   FileINode
	fileC       chan []byte
	lastBlock   []byte
	fileSize    uint64
	blockSize   uint32
	blockCursor uint32
}

func NewReader(fileInode FileINode, client blockstore.Client, fileSize uint64, blockSize uint32) *Reader {
	r := &Reader{
		fileInode: fileInode,
		client:    client,
		fileC:     make(chan []byte),
		fileSize:  fileSize,
		blockSize: blockSize,
	}
	return r
}

func (r *Reader) Read(b []byte) (n int, err error) {
	sizeRemaining := uint32(len(b))
	sizeRead := 0
	if r.lastBlock != nil {
		if len(b) <= len(r.lastBlock[r.blockCursor:]) {
			copy(b, r.lastBlock[r.blockCursor:sizeRemaining])
			r.blockCursor += sizeRemaining
			r.offset += int64(sizeRemaining)
			sizeRead = int(sizeRemaining)
			return sizeRead, nil
		} else {
			remblockSize := r.blockSize - r.blockCursor
			copy(b[:remblockSize], r.lastBlock[r.blockCursor:r.blockSize])
			r.lastBlock = nil
			r.blockCursor = 0
			r.offset += int64(remblockSize)
			sizeRemaining -= remblockSize
			sizeRead += int(remblockSize)
		}
	}

	if r.lastBlock == nil {
		noOfBlocks := int((sizeRemaining / r.blockSize) + 1)
		for i := 0; i < noOfBlocks; i++ {
			blockIndex := (r.offset / int64(r.blockSize)) + 1
			if blockIndex > int64(len(r.fileInode.FileBlocks)) {
				return sizeRead, fmt.Errorf("asking past EOF")
			}
			r.lastBlock, err = r.getBlock(r.fileInode.FileBlocks[blockIndex].Address)
			if err != nil {
				return sizeRead, err
			}
			copySize := r.blockSize
			if uint32(len(b))-sizeRemaining < r.blockSize {
				copySize = uint32(len(b)) - sizeRemaining
			}
			copy(b[sizeRead:copySize], r.lastBlock[:copySize])
			if copySize == r.blockSize {
				r.lastBlock = nil
				r.blockCursor = 0
			}
			r.offset += int64(copySize)
			sizeRemaining -= copySize
			sizeRead += int(copySize)

			if sizeRemaining <= 0 {
				return sizeRead, nil
			}
		}
	}
	return 0, nil
}

func (r *Reader) getBlock(addr []byte) ([]byte, error) {
	stdoutBytes, _, err := r.client.DownloadBlob(addr)
	if err != nil {
		return nil, err
	}
	return stdoutBytes, nil
}

func (r *Reader) Close() error {
	return nil
}
