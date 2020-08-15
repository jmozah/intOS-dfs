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
	"fmt"
	"strings"
	"time"

	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (d *Directory) DirStat(podName, dirName string, dirInode *DirInode, account string, podAddr string) error {

	meta := dirInode.Meta
	listing := d.ListDir(podName, dirName)

	files := 0
	dirs := 0
	for _, list := range listing {
		if strings.HasPrefix(list, "<Dir>") {
			dirs++
		} else {
			files++
		}
	}
	path := meta.Path
	if meta.Path == podName {
		path = utils.PathSeperator
	}

	fmt.Println("Account 	: ", account)
	fmt.Println("Pod Address	: ", podAddr)
	fmt.Println("PodName 	: ", podName)
	fmt.Println("Dir Path	: ", path)
	fmt.Println("Dir Name	: ", meta.Name)
	fmt.Println("Cr. Time	: ", time.Unix(meta.CreationTime, 0).String())
	fmt.Println("Mo. Time	: ", time.Unix(meta.ModificationTime, 0).String())
	fmt.Println("Ac. Time	: ", time.Unix(meta.AccessTime, 0).String())
	fmt.Println("Child Dirs	: ", dirs)
	fmt.Println("Child files	: ", files)
	return nil
}
