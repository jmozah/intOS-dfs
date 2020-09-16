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
	"encoding/json"
	"fmt"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

func (u *Users) CreateRootFeeds(userInfo *Info) error {
	rootAddress := userInfo.GetAccount().GetAddress(account.UserAccountIndex)
	err := u.CreateSettingsFeeds(rootAddress, userInfo)
	if err != nil {
		return err
	}
	err = u.CreateSharingFeeds(rootAddress, userInfo)
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) CreateSettingsFeeds(rootAddress utils.Address, userInfo *Info) error {
	// create name feed
	name := &Name{}
	data, err := json.Marshal(&name)
	if err != nil {
		return fmt.Errorf("create name feed: %w", err)
	}
	topic := utils.HashString(nameFeedName)
	_, err = userInfo.GetFeed().CreateFeed(topic, rootAddress, data)
	if err != nil {
		return fmt.Errorf("create name feed: %w", err)
	}

	// create contacts feed
	contacts := &Contacts{}
	data, err = json.Marshal(&contacts)
	if err != nil {
		return fmt.Errorf("create contacts feed: %w", err)
	}
	topic = utils.HashString(contactsFeedName)
	_, err = userInfo.GetFeed().CreateFeed(topic, rootAddress, data)
	if err != nil {
		return fmt.Errorf("create contacts feed: %w", err)
	}

	// create avatar feed
	topic = utils.HashString(avatarFeedName)
	data = make([]byte, 0)
	_, err = userInfo.GetFeed().CreateFeed(topic, rootAddress, data)
	if err != nil {
		return fmt.Errorf("create avatar feed: %w", err)
	}

	return nil
}

func (u *Users) CreateSharingFeeds(rootAddress utils.Address, userInfo *Info) error {
	// create inbox feed data
	inboxFile := &Inbox{Entries: make([]SharingEntry, 0)}
	inboxFileBytes, err := json.Marshal(&inboxFile)
	if err != nil {
		return fmt.Errorf("create sharing inbox: %w", err)
	}

	// store the new inbox file data
	newInboxRef, err := u.client.UploadBlob(inboxFileBytes, true)
	if err != nil {
		return fmt.Errorf("create sharing inbox: %w", err)
	}

	// store the inbox reference in to inbox feed
	topic := utils.HashString(inboxFeedName)
	_, err = userInfo.GetFeed().CreateFeed(topic, rootAddress, newInboxRef)
	if err != nil {
		return fmt.Errorf("create sharing inbox: %w", err)
	}

	// create outbox feed data
	outFile := &Outbox{Entries: make([]SharingEntry, 0)}
	outboxFileBytes, err := json.Marshal(&outFile)
	if err != nil {
		return fmt.Errorf("create sharing inbox: %w", err)
	}

	// store the new outbox file data
	newOutboxRef, err := u.client.UploadBlob(outboxFileBytes, true)
	if err != nil {
		return fmt.Errorf("create sharing outbox: %w", err)
	}

	// store the outbox reference in to ourbox feed
	topic = utils.HashString(outboxFeedName)
	_, err = userInfo.GetFeed().CreateFeed(topic, rootAddress, newOutboxRef)
	if err != nil {
		return fmt.Errorf("create sharing inbox: %w", err)
	}
	return nil
}
