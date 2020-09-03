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
)

const (
	inboxFileName  = "inbox"
	outboxFileName = "outbox"
)

type Inbox struct {
	Entries []InboxEntry `json:"entries"`
}

type InboxEntry struct {
	FileMetaHash string `json:"meta_hash"`
	Sender       string `json:"sender_ref"`
	SentTime     string `json:"sent_time"`
}

type Outbox struct {
	Entries []OutboxEntry `json:"entries"`
}

type OutboxEntry struct {
	SourcePod    string `json:"pod_name"`
	FileMetaHash string `json:"meta_hash"`
	Sender       string `json:"sender_ref"`
	SentTime     string `json:"sent_time"`
}

func (p *Pod) CreateSharingFiles(podName, podDir string) error {
	// create inbox file
	inboxFile := &Inbox{Entries: make([]InboxEntry, 0)}
	data, err := json.Marshal(&inboxFile)
	if err != nil {
		return fmt.Errorf("sharing: %w", err)
	}
	reader := bytes.NewReader(data)
	_, err = p.UploadFile(podName, inboxFileName, int64(len(data)), reader, podDir, "10M")
	if err != nil {
		return fmt.Errorf("sharing: %w", err)
	}

	// create outbox file
	outFile := &Outbox{Entries: make([]OutboxEntry, 0)}
	data, err = json.Marshal(&outFile)
	if err != nil {
		return fmt.Errorf("sharing: %w", err)
	}
	reader = bytes.NewReader(data)
	_, err = p.UploadFile(podName, outboxFileName, int64(len(data)), reader, podDir, "10M")
	if err != nil {
		return fmt.Errorf("sharing: %w", err)
	}
	return nil
}
