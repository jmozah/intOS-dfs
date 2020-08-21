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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type FileStats struct {
	Account          string `json:"account"`
	PodName          string `json:"pod_name"`
	FilePath         string `json:"file_path"`
	FileName         string `json:"file_name"`
	FileSize         string `json:"file_size"`
	BlockSize        string `json:"block_size"`
	CreationTime     string `json:"creation_time"`
	ModificationTime string `json:"modification_time"`
	AccessTime       string `json:"access_time"`
	Blocks           []Blocks
}

type Blocks struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
	Size      string `json:"size"`
}

func (f *File) FileStat(podName string, fileName string, account string) (*FileStats, error) {
	meta := f.GetFromFileMap(fileName)
	if meta == nil {
		return nil, fmt.Errorf("file stat: file not found")
	}

	fileInodeBytes, _, err := f.getClient().DownloadBlob(meta.InodeAddress)
	if err != nil {
		return nil, fmt.Errorf("file stat: could not find file Inode")
	}
	var fileInode FileINode
	err = json.Unmarshal(fileInodeBytes, &fileInode)
	if err != nil {
		return nil, fmt.Errorf("stat: file Inode unmarshall error: %w", err)
	}

	var fileBlocks []Blocks
	for _, b := range fileInode.FileBlocks {
		fb := Blocks{
			Name:      b.Name,
			Reference: hex.EncodeToString(b.Address),
			Size:      strconv.Itoa(int(b.Size)),
		}
		fileBlocks = append(fileBlocks, fb)
	}
	return &FileStats{
		Account:          account,
		PodName:          podName,
		FilePath:         meta.Path,
		FileName:         meta.Name,
		FileSize:         strconv.FormatUint(meta.FileSize, 10),
		BlockSize:        strconv.Itoa(int(meta.BlockSize)),
		CreationTime:     time.Unix(meta.CreationTime, 0).String(),
		ModificationTime: time.Unix(meta.ModificationTime, 0).String(),
		AccessTime:       time.Unix(meta.AccessTime, 0).String(),
		Blocks:           fileBlocks,
	}, nil

	//fmt.Println("Account 	: ", account)
	//fmt.Println("PodName 	: ", podName)
	//fmt.Println("File Path	: ", meta.Path)
	//fmt.Println("File Name	: ", meta.Name)
	//fmt.Println("File Size	: ", meta.FileSize, " Bytes")
	//fmt.Println("Block Size	: ", meta.BlockSize, " Bytes")
	//fmt.Println("Cr. Time	: ", time.Unix(meta.CreationTime, 0).String())
	//fmt.Println("Mo. Time	: ", time.Unix(meta.ModificationTime, 0).String())
	//fmt.Println("Ac. Time	: ", time.Unix(meta.AccessTime, 0).String())
	//fmt.Println("----- Blocks -------")
	//for _, fb := range fileInode.FileBlocks {
	//	blkStr := fmt.Sprintf("%s, 0x%s, %d bytes", fb.Name, hex.EncodeToString(fb.Address), fb.Size)
	//	fmt.Println(blkStr)
	//}
}
