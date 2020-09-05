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
	"time"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/feed"
	"github.com/jmozah/intOS-dfs/pkg/pod"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

const (
	nameFile       = "Name"
	contactsFile   = "Contacts"
	inboxFileName  = "inbox"
	outboxFileName = "outbox"
)

type SystemFile struct {
	NameFile      string `json:"name,omitempty"`
	ContactsFile  string `json:"contacts,omitempty"`
	SharingInbox  string `json:"sharing_inbox,omitempty"`
	SharingOutbox string `json:"sharing_outbox,omitempty"`
}

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

type Inbox struct {
	Entries []InboxEntry `json:"entries"`
}

type InboxEntry struct {
	FileName     string `json:"name"`
	FileMetaHash string `json:"meta_hash"`
	Sender       string `json:"source_ref"`
	Receiver     string `json:"dest_ref"`
	SharedTime   string `json:"sent_time"`
}

type Outbox struct {
	Entries []OutboxEntry `json:"entries"`
}

type OutboxEntry struct {
	FileName     string `json:"name"`
	SourcePod    string `json:"pod_name"`
	FileMetaHash string `json:"meta_hash"`
	Sender       string `json:"source_ref"`
	Receiver     string `json:"dest_ref"`
	SharedTime   string `json:"shared_time"`
}

func (u *Users) CreateRootFeeds(userInfo *Info) error {
	rootReference := userInfo.GetAccount().GetAddress(account.UserAccountIndex)
	err := u.CreateSettingsFeeds(rootReference, userInfo)
	if err != nil {
		return err
	}
	err = u.CreateSharingFeeds(rootReference, userInfo)
	if err != nil {
		return err
	}
	return nil
}

func (u *Users) CreateSettingsFeeds(rootReference utils.Address, userInfo *Info) error {
	// create name file
	name := &Name{}
	data, err := json.Marshal(&name)
	if err != nil {
		return fmt.Errorf("create name feed: %w", err)
	}
	topic := utils.HashString(nameFile)
	_, err = userInfo.GetFeed().CreateFeed(topic, rootReference, data)
	if err != nil {
		return fmt.Errorf("create name feed: %w", err)
	}

	// create contacts file
	contacts := &Contacts{}
	data, err = json.Marshal(&contacts)
	if err != nil {
		return fmt.Errorf("create contacts feed: %w", err)
	}
	topic = utils.HashString(contactsFile)
	_, err = userInfo.GetFeed().CreateFeed(topic, rootReference, data)
	if err != nil {
		return fmt.Errorf("create contacts feed: %w", err)
	}
	return nil
}

func (u *Users) CreateSharingFeeds(rootReference utils.Address, userInfo *Info) error {
	// create inbox file
	inboxFile := &Inbox{Entries: make([]InboxEntry, 0)}
	data, err := json.Marshal(&inboxFile)
	if err != nil {
		return fmt.Errorf("create sharing inbox: %w", err)
	}
	topic := utils.HashString(inboxFileName)
	_, err = userInfo.GetFeed().CreateFeed(topic, rootReference, data)
	if err != nil {
		return fmt.Errorf("create sharing inbox: %w", err)
	}

	// create outbox file
	outFile := &Outbox{Entries: make([]OutboxEntry, 0)}
	data, err = json.Marshal(&outFile)
	if err != nil {
		return fmt.Errorf("create sharing outbox: %w", err)
	}
	topic = utils.HashString(outboxFileName)
	_, err = userInfo.GetFeed().CreateFeed(topic, rootReference, data)
	if err != nil {
		return fmt.Errorf("create sharing inbox: %w", err)
	}
	return nil
}

func (u *Users) ShareFileWithUser(podName, podFilePath, destinationRef string, userInfo *Info, pod *pod.Pod) (*OutboxEntry, error) {
	// Get the meta reference of the file to share
	metaRef, fileName, err := pod.GetMetaReferenceOfFile(podName, podFilePath)
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}

	// Create a outbox entry
	rootReference := userInfo.GetAccount().GetAddress(account.UserAccountIndex)
	now := time.Now().String()
	outEntry := OutboxEntry{
		FileName:     fileName,
		SourcePod:    userInfo.podName,
		FileMetaHash: utils.BytesToAddress(metaRef).Hex(),
		Sender:       rootReference.Hex(),
		Receiver:     destinationRef,
		SharedTime:   now,
	}

	// add the entry to outbox of the sender
	data, err := getFeedData(outboxFileName, rootReference, userInfo.GetFeed())
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}
	outbox := &Outbox{}
	err = json.Unmarshal(data, outbox)
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}
	outbox.Entries = append(outbox.Entries, outEntry)
	outData, err := json.Marshal(outbox)
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}
	err = putFeedData(outboxFileName, rootReference, outData, userInfo.GetFeed())
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}
	return &outEntry, nil
}

func (u *Users) ReceiveFileFromUser(podName string, outboxEntry OutboxEntry, userInfo *Info, pod *pod.Pod) error {
	// construct the inbox entry
	inboxEntry := InboxEntry{
		FileMetaHash: outboxEntry.FileMetaHash,
		Sender:       outboxEntry.Sender,
		Receiver:     outboxEntry.Receiver,
		SharedTime:   outboxEntry.SharedTime,
	}

	// add the file to the pod directory specified
	err := pod.ReceiveFileAndStore(podName, utils.PathSeperator, outboxEntry.FileName, outboxEntry.FileMetaHash)
	if err != nil {
		return fmt.Errorf("share: %w", err)
	}

	// add the inbox entry to inbox
	rootReference := userInfo.GetAccount().GetAddress(account.UserAccountIndex)
	data, err := getFeedData(outboxFileName, rootReference, userInfo.GetFeed())
	if err != nil {
		return fmt.Errorf("share: %w", err)
	}
	inbox := &Inbox{}
	err = json.Unmarshal(data, inbox)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}
	inbox.Entries = append(inbox.Entries, inboxEntry)
	inData, err := json.Marshal(inbox)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}
	return putFeedData(inboxFileName, rootReference, inData, userInfo.GetFeed())
}

func getFeedData(fileName string, rootReference utils.Address, fd *feed.API) ([]byte, error) {
	topic := utils.HashString(fileName)
	_, data, err := fd.GetFeedData(topic, rootReference)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func putFeedData(fileName string, rootReference utils.Address, data []byte, fd *feed.API) error {
	topic := utils.HashString(fileName)
	_, err := fd.UpdateFeed(topic, rootReference, data)
	if err != nil {
		return err
	}
	return nil
}
