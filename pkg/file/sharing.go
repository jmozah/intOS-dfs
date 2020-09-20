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
	"encoding/json"
	"fmt"
	"net/http"

	m "github.com/jmozah/intOS-dfs/pkg/meta"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (f *File) GetFileReference(podFile string) ([]byte, string, error) {
	// Get the meta of the file to share
	meta := f.GetFromFileMap(podFile)
	if meta == nil {
		return nil, "", fmt.Errorf("file not found in dfs")
	}
	return meta.MetaReference, meta.Name, nil
}

func (f *File) AddFileToPath(filePath, metaHexRef string) error {
	metaReferenace, err := utils.ParseHexReference(metaHexRef)
	if err != nil {
		return err
	}
	data, respCode, err := f.getClient().DownloadBlob(metaReferenace.Bytes())
	if err != nil || respCode != http.StatusOK {
		return err
	}
	meta := &m.FileMetaData{}
	err = json.Unmarshal(data, meta)
	if err != nil {
		return err
	}
	meta.MetaReference = metaReferenace.Bytes()
	f.AddToFileMap(filePath, meta)
	return nil
}
