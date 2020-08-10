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

package datapod

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

func (f *File) FileStat(podName string, fileName string, account string) error {
	meta := f.GetFromFileMap(fileName)
	if meta == nil {
		return fmt.Errorf("file stat: file not found")
	}

	fileInodeBytes, _, err := f.getClient().DownloadBlob(meta.InodeAddress)
	if err != nil {
		return fmt.Errorf("file stat: could not find file Inode")
	}
	var fileInode FileINode
	err = json.Unmarshal(fileInodeBytes, &fileInode)
	if err != nil {
		return fmt.Errorf("stat: file Inode unmarshall error: %w", err)
	}

	fmt.Println("Account 	: ", account)
	fmt.Println("PodName 	: ", podName)
	fmt.Println("File Path	: ", meta.Path)
	fmt.Println("File Name	: ", meta.Name)
	fmt.Println("File Size	: ", meta.FileSize, " Bytes")
	fmt.Println("Block Size	: ", meta.BlockSize, " Bytes")
	fmt.Println("Cr. Time	: ", time.Unix(meta.CreationTime, 0).String())
	fmt.Println("Mo. Time	: ", time.Unix(meta.ModificationTime, 0).String())
	fmt.Println("Ac. Time	: ", time.Unix(meta.AccessTime, 0).String())
	fmt.Println("----- Blocks -------")
	for _, fb := range fileInode.FileBlocks {
		blkStr := fmt.Sprintf("%s, 0x%s, %d bytes", fb.Name, hex.EncodeToString(fb.Address), fb.Size)
		fmt.Println(blkStr)
	}
	return nil
}
