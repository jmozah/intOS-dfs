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

package user

import (
	"fmt"

	"github.com/jmozah/intOS-dfs/pkg/feed"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

const (
	AvatarFileSuffix = ".avatar"
)

func (u *Users) StoreSettingsFile(rootReference utils.Address, fileType string, fd *feed.API, data []byte) error {
	fileName := rootReference.Hex() + fileType
	topic := utils.HashString(fileName)

	// If the avatar not already present, create it
	_, _, err := fd.GetFeedData(topic, rootReference)
	if err != nil {
		if err.Error() == "no feed updates found" {
			_, err := fd.CreateFeed(topic, rootReference, data)
			if err != nil {
				return fmt.Errorf("store settings: %w", err)
			}
			return nil
		}
	}

	// if the avatar is already present, update it
	_, err = fd.UpdateFeed(topic, rootReference, data)
	if err != nil {
		return fmt.Errorf("store settings: %w", err)
	}
	return nil
}

func (u *Users) LoadSettingsFile(rootReference utils.Address, fileType string, fd *feed.API) ([]byte, error) {
	fileName := rootReference.Hex() + fileType
	topic := utils.HashString(fileName)

	_, data, err := fd.GetFeedData(topic, rootReference)
	if err != nil {
		if err.Error() != "no feed updates found" {
			return nil, fmt.Errorf("load settings: %w", err)
		}
	}
	return data, nil
}
