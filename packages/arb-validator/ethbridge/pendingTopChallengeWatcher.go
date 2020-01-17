/*
 * Copyright 2020, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ethbridge

import (
	"context"
	"errors"
	"log"
	"strings"

	errors2 "github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/ethbridge/pendingtopchallenge"
)

var pendingTopBisectedID ethcommon.Hash
var pendingTopOneStepProofCompletedID ethcommon.Hash

func init() {
	parsed, err := abi.JSON(strings.NewReader(pendingtopchallenge.PendingTopChallengeABI))
	if err != nil {
		panic(err)
	}
	pendingTopBisectedID = parsed.Events["Bisected"].ID()
	pendingTopOneStepProofCompletedID = parsed.Events["OneStepProofCompleted"].ID()
}

type pendingTopChallengeWatcher struct {
	*bisectionChallengeWatcher
	contract *pendingtopchallenge.PendingTopChallenge
	client   *ethclient.Client
	address  ethcommon.Address
}

func newPendingTopChallengeWatcher(address ethcommon.Address, client *ethclient.Client) (*pendingTopChallengeWatcher, error) {
	bisectionChallenge, err := newBisectionChallengeWatcher(address, client)
	if err != nil {
		return nil, err
	}
	pendingTopContract, err := pendingtopchallenge.NewPendingTopChallenge(address, client)
	if err != nil {
		return nil, errors2.Wrap(err, "Failed to connect to PendingTopChallenge")
	}
	return &pendingTopChallengeWatcher{
		bisectionChallengeWatcher: bisectionChallenge,
		contract:                  pendingTopContract,
		client:                    client,
		address:                   address,
	}, nil
}

func (c *pendingTopChallengeWatcher) topics() []ethcommon.Hash {
	tops := []ethcommon.Hash{
		pendingTopBisectedID,
		pendingTopOneStepProofCompletedID,
	}
	return append(tops, c.bisectionChallengeWatcher.topics()...)
}

func (c *pendingTopChallengeWatcher) StartConnection(ctx context.Context, startHeight *common.TimeBlocks, startLogIndex uint) (<-chan arbbridge.MaybeEvent, error) {
	filter := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{c.address},
		Topics:    [][]ethcommon.Hash{c.topics()},
	}

	logCtx, cancelFunc := context.WithCancel(ctx)
	maybeLogChan, err := getLogs(logCtx, c.client, filter, startHeight, startLogIndex)
	if err != nil {
		return nil, err
	}

	eventChan := make(chan arbbridge.MaybeEvent, 1024)
	go func() {
		defer close(eventChan)
		defer cancelFunc()
		for {
			select {
			case <-ctx.Done():
				log.Println("pendingTopChallengeWatcher context canceled")
				return
			case maybeLog, ok := <-maybeLogChan:
				if !ok {
					log.Println("maybeLogChan terminated")
					eventChan <- arbbridge.MaybeEvent{Err: errors.New("logChan terminated early")}
					return
				}
				if maybeLog.err != nil {
					log.Println("maybeLog had error", err)
					eventChan <- arbbridge.MaybeEvent{Err: err}
					return
				}
				header, err := c.client.HeaderByHash(ctx, maybeLog.log.BlockHash)
				if err != nil {
					log.Println("header had error", err)
					eventChan <- arbbridge.MaybeEvent{Err: err}
					return
				}
				chainInfo := getChainInfo(maybeLog.log, header)
				event, err := c.parsePendingTopEvent(chainInfo, maybeLog.log)
				if err != nil {
					log.Println("Failed parsing event", err)
					eventChan <- arbbridge.MaybeEvent{Err: err}
					return
				}
				eventChan <- arbbridge.MaybeEvent{Event: event}
			}
		}
	}()
	return eventChan, nil
}

func (c *pendingTopChallengeWatcher) parsePendingTopEvent(chainInfo arbbridge.ChainInfo, log types.Log) (arbbridge.Event, error) {
	if log.Topics[0] == pendingTopBisectedID {
		eventVal, err := c.contract.ParseBisected(log)
		if err != nil {
			return nil, err
		}
		return arbbridge.PendingTopBisectionEvent{
			ChainInfo:   chainInfo,
			ChainHashes: hashSliceToHashes(eventVal.ChainHashes),
			TotalLength: eventVal.TotalLength,
			Deadline:    common.TimeTicks{Val: eventVal.DeadlineTicks},
		}, nil
	} else if log.Topics[0] == pendingTopOneStepProofCompletedID {
		_, err := c.contract.ParseOneStepProofCompleted(log)
		if err != nil {
			return nil, err
		}
		return arbbridge.OneStepProofEvent{
			ChainInfo: chainInfo,
		}, nil
	}
	return c.bisectionChallengeWatcher.parseBisectionEvent(chainInfo, log)
}
