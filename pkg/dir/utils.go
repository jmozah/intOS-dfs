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
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (d *DirInode) IsDirInodeRoot() bool {
	if d.Meta.Path == utils.PathSeperator {
		return true
	}
	return false
}

func (d *DirInode) GetDirInodePathAndNameForRoot() string {
	return d.Meta.Path + d.Meta.Name
}

func (d *DirInode) GetDirInodePathAndName() string {
	return d.Meta.Path + utils.PathSeperator + d.Meta.Name
}
