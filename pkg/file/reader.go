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
	"io"

	"github.com/jmozah/intOS-dfs/pkg/blockstore"
)

type Reader struct {
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
	bytesToRead := uint32(len(b))
	bytesRead := 0
	if r.lastBlock != nil {
		remDataSize := r.blockSize - r.blockCursor
		if bytesToRead <= remDataSize {
			copy(b, r.lastBlock[r.blockCursor:r.blockCursor+bytesToRead])
			r.blockCursor += bytesToRead
			r.offset += int64(bytesToRead)
			bytesRead = int(bytesToRead)
			//bytesToRead = 0
			if r.blockCursor == r.blockSize {
				r.lastBlock = nil
				r.blockCursor = 0
			}
			return bytesRead, nil
		} else {
			copy(b, r.lastBlock[r.blockCursor:r.blockSize])
			r.lastBlock = nil
			r.blockCursor = 0
			r.offset += int64(remDataSize)
			bytesRead += int(remDataSize)
			bytesToRead -= remDataSize
			// read spans across block.. so flow down and read the next block
		}
	}

	if r.lastBlock == nil {
		noOfBlocks := int((bytesToRead / r.blockSize) + 1)
		for i := 0; i < noOfBlocks; i++ {
			if r.lastBlock == nil {
				blockIndex := (r.offset / int64(r.blockSize))
				if blockIndex > int64(len(r.fileInode.FileBlocks)) {
					return bytesRead, fmt.Errorf("asking past EOF")
				}
				if blockIndex >= int64(len(r.fileInode.FileBlocks)) {
					return 0, io.EOF
				}
				r.lastBlock, err = r.getBlock(r.fileInode.FileBlocks[blockIndex].Address)
				if err != nil {
					return bytesRead, err
				}
			}

			// if length of bytes to read is greater than block size
			if bytesToRead > r.blockSize {
				bytesToRead = r.blockSize
			}

			copy(b[bytesRead:bytesToRead], r.lastBlock[:bytesToRead])
			if bytesToRead == r.blockSize {
				r.lastBlock = nil
				r.blockCursor = 0
			} else {
				r.blockCursor += bytesToRead
			}
			r.offset += int64(bytesToRead)
			bytesRead += int(bytesToRead)
			bytesToRead -= bytesToRead

			if bytesToRead <= 0 {
				return bytesRead, nil
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
