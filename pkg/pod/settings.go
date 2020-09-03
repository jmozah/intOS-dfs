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

package pod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	NameFile     = "Name"
	ContactsFile = "Contacts"
)

type Name struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	SurName    string `json:"surname"`
}

type Contacts struct {
	Phone  string  `json:"phone_number"`
	Mobile string  `json:"mobile"`
	Addr   Address `json:"address"`
}

type Address struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	State        string `json:"state/Province/Region"`
	ZipCode      string `json:"zip_code"`
}

func (p *Pod) CreateSettingsFiles(podName, podDir string) error {
	// create name file
	name := &Name{}
	data, err := json.Marshal(&name)
	if err != nil {
		return fmt.Errorf("setting: %w", err)
	}
	reader := bytes.NewReader(data)
	_, err = p.UploadFile(podName, NameFile, int64(len(data)), reader, podDir, "1M")
	if err != nil {
		return fmt.Errorf("setting: %w", err)
	}

	// create contacts file
	contacts := &Contacts{}
	data, err = json.Marshal(&contacts)
	if err != nil {
		return fmt.Errorf("setting: %w", err)
	}
	reader = bytes.NewReader(data)
	_, err = p.UploadFile(podName, ContactsFile, int64(len(data)), reader, podDir, "1M")
	if err != nil {
		return fmt.Errorf("setting: %w", err)
	}

	return nil
}

func (p *Pod) SaveNameFile(podName, podDir string, data []byte) error {
	reader := bytes.NewReader(data)
	_, err := p.UploadFile(podName, NameFile, int64(len(data)), reader, podDir, "1M")
	if err != nil {
		return fmt.Errorf("name: %w", err)
	}
	return nil
}

func (p *Pod) GetNameFile(podName, podFile string) (*Name, error) {
	reader, _, _, err := p.DownloadFile(podName, podFile)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	name := &Name{}
	err = json.Unmarshal(data, name)
	if err != nil {
		return nil, err
	}

	return name, nil
}
