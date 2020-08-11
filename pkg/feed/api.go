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

package feed

import (
	"context"
	"fmt"
	"time"

	"github.com/ethersphere/bee/pkg/content"
	"github.com/ethersphere/bee/pkg/crypto"
	"github.com/ethersphere/bee/pkg/soc"
	"github.com/ethersphere/bee/pkg/swarm"
	bmtlegacy "github.com/ethersphere/bmt/legacy"
	"github.com/jmozah/intOS-dfs/pkg/account"
	"github.com/jmozah/intOS-dfs/pkg/blockstore"
	"github.com/jmozah/intOS-dfs/pkg/feed/lookup"
	"github.com/jmozah/intOS-dfs/pkg/utils"
)

const (
	maxuint64 = ^uint64(0)
	idLength  = 32
)

var (
	// ErrInvalidTopicSize is returned when a topic is not equal to TopicLength
	ErrInvalidTopicSize = fmt.Errorf("Topic is not equal to %d", TopicLength)

	// ErrInvalidPayloadSize is returned when the payload is greater than the chunk size
	ErrInvalidPayloadSize = fmt.Errorf("payload is greater than %d", utils.MaxChunkLength)
)

type API struct {
	handler *Handler
	acc     *account.Account
}

type Request struct {
	ID
	//User   utils.Address
	idAddr swarm.Address // cached chunk address for the update (not serialized, for internal use)

	data       []byte     // actual data payload
	Signature  *Signature // Signature of the payload
	binaryData []byte     // cached serialized data (does not get serialized again!, for efficiency/internal use)
}

func New(account *account.Account, client blockstore.Client) *API {
	bmtPool := bmtlegacy.NewTreePool(hashFunc, swarm.Branches, bmtlegacy.PoolSize)
	return &API{
		handler: NewHandler(account, client, bmtPool),
		acc:     account,
	}
}

// create feed
func (a *API) CreateFeed(topic []byte, user utils.Address, data []byte) ([]byte, error) {
	var req Request

	if len(topic) != TopicLength {
		return nil, ErrInvalidTopicSize
	}

	if len(data) > utils.MaxChunkLength {
		return nil, ErrInvalidPayloadSize
	}

	// fill Feed and Epoc related details
	copy(req.ID.Topic[:], topic)
	req.ID.User = user
	req.Epoch.Level = 31
	req.Epoch.Time = uint64(time.Now().Unix())

	// Add initial feed data
	req.data = data

	// create the id, hash(topic, epoc)
	id, err := a.handler.getId(req.Topic, req.Time, req.Level)
	if err != nil {
		return nil, err
	}

	// get the payload id BMT(span, payload)
	payloadId, err := a.handler.getPayloadId(data)
	if err != nil {
		return nil, err
	}

	// create the signer and the content addressed chunk
	signer := crypto.NewDefaultSigner(a.acc.GetPrivateKey())
	ch, err := content.NewChunk(data)
	if err != nil {
		return nil, err
	}
	sch, err := soc.NewChunk(id, ch, signer)
	if err != nil {
		return nil, err
	}

	// set the address and the data for the soc chunk
	req.idAddr = sch.Address()
	req.binaryData = sch.Data()

	// set signature and binary data fields
	_, err = a.handler.toChunkContent(&req, id, payloadId)
	if err != nil {
		return nil, err
	}

	// send the updated soc chunk to bee
	address, err := a.handler.update(&req)
	if err != nil {
		return nil, err
	}
	return address, nil
}

func (a *API) GetFeedData(topic []byte, user utils.Address) ([]byte, []byte, error) {
	if len(topic) != TopicLength {
		return nil, nil, ErrInvalidTopicSize
	}

	ctx := context.Background()
	f := new(Feed)
	f.User = user
	copy(f.Topic[:], topic)

	// create the query from values
	q := &Query{Feed: *f}
	q.TimeLimit = 0
	q.Hint = lookup.NoClue
	_, err := a.handler.Lookup(ctx, q)
	if err != nil {
		return nil, nil, err
	}
	var data []byte
	addr, data, err := a.handler.GetContent(&q.Feed)
	if err != nil {
		return nil, nil, err
	}
	return addr.Bytes(), data, nil

}

func (a *API) UpdateFeed(topic []byte, user utils.Address, data []byte) ([]byte, error) {
	if len(topic) != TopicLength {
		return nil, ErrInvalidTopicSize
	}

	if len(data) > utils.MaxChunkLength {
		return nil, ErrInvalidPayloadSize
	}

	ctx := context.Background()
	f := new(Feed)
	f.User = user
	copy(f.Topic[:], topic)

	// get the existing request from DB
	req, err := a.handler.NewRequest(ctx, f)
	if err != nil {
		return nil, err
	}
	req.Time = uint64(time.Now().Unix())
	req.data = data

	// create the id, hash(topic, epoc)
	id, err := a.handler.getId(req.Topic, req.Time, req.Level)
	if err != nil {
		return nil, err
	}

	// get the payload id BMT(span, payload)
	payloadId, err := a.handler.getPayloadId(data)
	if err != nil {
		return nil, err
	}

	// create the signer and the content addressed chunk
	signer := crypto.NewDefaultSigner(a.acc.GetPrivateKey())
	ch, err := content.NewChunk(data)
	if err != nil {
		return nil, err
	}
	sch, err := soc.NewChunk(id, ch, signer)
	if err != nil {
		return nil, err
	}

	// set the address and the data for the soc chunk
	req.idAddr = sch.Address()
	req.binaryData = sch.Data()

	// set signature and binary data fields
	_, err = a.handler.toChunkContent(req, id, payloadId)
	if err != nil {
		return nil, err
	}

	address, err := a.handler.update(req)
	if err != nil {
		return nil, err
	}
	return address, nil
}


