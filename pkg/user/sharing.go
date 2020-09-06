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
	"net/http"
	"path/filepath"
	"time"

	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/pod"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

const (
	inboxFeedName  = "inbox"
	outboxFeedName = "outbox"
)

type Inbox struct {
	Entries []InboxEntry `json:"entries"`
}

type InboxEntry struct {
	FilePath     string `json:"file_path"`
	PodName      string `json:"pod_name"`
	FileMetaHash string `json:"meta_ref"`
	Sender       string `json:"source_address"`
	Receiver     string `json:"dest_address"`
	SharedTime   string `json:"shared_time"`
}

type Outbox struct {
	Entries []OutboxEntry `json:"entries"`
}

type OutboxEntry struct {
	FileName     string `json:"name"`
	PodName      string `json:"pod_name"`
	FileMetaHash string `json:"meta_ref"`
	Sender       string `json:"source_address"`
	Receiver     string `json:"dest_address"`
	SharedTime   string `json:"shared_time"`
}

func (u *Users) ShareFileWithUser(podName, podFilePath, destinationRef string, userInfo *Info, pod *pod.Pod) (*OutboxEntry, error) {
	// Get the meta reference of the file to share
	metaRef, fileName, err := pod.GetMetaReferenceOfFile(podName, podFilePath)
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}

	// Create a outbox entry
	rootReference := userInfo.GetAccount().GetAddress(account.UserAccountIndex)
	now := time.Now()
	nowTimeStr := now.Format(time.RFC3339)
	outEntry := OutboxEntry{
		FileName:     fileName,
		PodName:      userInfo.podName,
		FileMetaHash: utils.NewReference(metaRef).String(),
		Sender:       rootReference.String(),
		Receiver:     destinationRef,
		SharedTime:   nowTimeStr,
	}

	// get the outbox reference from outbox feed
	outboxRef, err := getFeedData(outboxFeedName, rootReference, userInfo.GetFeed())
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}

	// download the entire outbox file
	outboxFileBytes, respCode, err := u.client.DownloadBlob(outboxRef)
	if err != nil && respCode != http.StatusOK {
		return nil, fmt.Errorf("share: %w", err)
	}

	// unmarshall, add a new outbox entry, marshall the data again
	outbox := &Outbox{}
	err = json.Unmarshal(outboxFileBytes, outbox)
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}
	outbox.Entries = append(outbox.Entries, outEntry)
	outData, err := json.Marshal(outbox)
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}

	// store the new outbox file data
	newOutboxRef, err := u.client.UploadBlob(outData)
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}

	// update the outbox feed with the new outbox file reference
	err = putFeedData(outboxFeedName, rootReference, newOutboxRef, userInfo.GetFeed())
	if err != nil {
		return nil, fmt.Errorf("share: %w", err)
	}
	return &outEntry, nil
}

func (u *Users) ReceiveFileFromUser(podName string, inboxEntry InboxEntry, userInfo *Info, pod *pod.Pod) error {
	// add the file to the pod directory specified
	podDir := filepath.Dir(inboxEntry.FilePath)
	fileName := filepath.Base(inboxEntry.FilePath)
	inboxEntry.PodName = podName
	err := pod.ReceiveFileAndStore(podName, podDir, fileName, inboxEntry.FileMetaHash)
	if err != nil {
		return fmt.Errorf("share: %w", err)
	}

	// get the inbox reference from inbox feed
	rootReference := userInfo.GetAccount().GetAddress(account.UserAccountIndex)
	inboxRef, err := getFeedData(inboxFeedName, rootReference, userInfo.GetFeed())
	if err != nil {
		return fmt.Errorf("share: %w", err)
	}

	// download the entire inbox file
	inboxFileBytes, respCode, err := u.client.DownloadBlob(inboxRef)
	if err != nil && respCode != http.StatusOK {
		return fmt.Errorf("share: %w", err)
	}

	// unmarshall, add a new inbox entry, marshall the data again
	inbox := &Inbox{}
	err = json.Unmarshal(inboxFileBytes, inbox)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}
	inbox.Entries = append(inbox.Entries, inboxEntry)
	inData, err := json.Marshal(inbox)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}

	// store the new inbox file data
	newInboxRef, err := u.client.UploadBlob(inData)
	if err != nil {
		return fmt.Errorf("share: %w", err)
	}

	// update the inbox feed with the new inbox file reference
	return putFeedData(inboxFeedName, rootReference, newInboxRef, userInfo.GetFeed())
}

func (u *Users) GetSharingInbox(userInfo *Info) (*Inbox, error) {
	// get the inbox reference from the inbox feed
	rootReference := userInfo.GetAccount().GetAddress(account.UserAccountIndex)
	inboxRef, err := getFeedData(inboxFeedName, rootReference, userInfo.GetFeed())
	if err != nil {
		return nil, fmt.Errorf("get sharing inbox: %w", err)
	}

	if len(inboxRef) < utils.ReferenceLength {
		return nil, fmt.Errorf("get sharing inbox: empty inbox")
	}

	// download the entire inbox file
	inboxFileBytes, respCode, err := u.client.DownloadBlob(inboxRef)
	if err != nil && respCode != http.StatusOK {
		return nil, fmt.Errorf("share: %w", err)
	}

	// unmarshall it and return the entire structure
	inbox := &Inbox{}
	err = json.Unmarshal(inboxFileBytes, inbox)
	if err != nil {
		return nil, fmt.Errorf("get sharing inbox: %w", err)
	}
	return inbox, nil
}

func (u *Users) GetSharingOutbox(userInfo *Info) (*Outbox, error) {
	// get the outbox reference from the inbox feed
	rootReference := userInfo.GetAccount().GetAddress(account.UserAccountIndex)
	outboxRef, err := getFeedData(outboxFeedName, rootReference, userInfo.GetFeed())
	if err != nil {
		return nil, fmt.Errorf("get sharing outbox: %w", err)
	}

	if len(outboxRef) < utils.ReferenceLength {
		return nil, fmt.Errorf("get sharing outbox: empty outbox")
	}

	// download the entire outbox file
	outboxFileBytes, respCode, err := u.client.DownloadBlob(outboxRef)
	if err != nil && respCode != http.StatusOK {
		return nil, fmt.Errorf("share: %w", err)
	}

	outbox := &Outbox{}
	err = json.Unmarshal(outboxFileBytes, outbox)
	if err != nil {
		return nil, fmt.Errorf("get sharing outbox: %w", err)
	}
	return outbox, nil
}
