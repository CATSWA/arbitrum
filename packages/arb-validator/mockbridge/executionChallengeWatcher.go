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

package mockbridge

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/valprotocol"
	errors2 "github.com/pkg/errors"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/ethbridge/executionchallenge"
)

var bisectedAssertionID ethcommon.Hash
var oneStepProofCompletedID ethcommon.Hash

//func init() {
//	parsed, err := abi.JSON(strings.NewReader(executionchallenge.ExecutionChallengeABI))
//	if err != nil {
//		panic(err)
//	}
//	bisectedAssertionID = parsed.Events["BisectedAssertion"].ID()
//	oneStepProofCompletedID = parsed.Events["OneStepProofCompleted"].ID()
//}

type executionChallengeWatcher struct {
	*bisectionChallengeWatcher
	challenge *executionchallenge.ExecutionChallenge
	client    arbbridge.ArbClient
	address   ethcommon.Address
}

func newExecutionChallengeWatcher(address ethcommon.Address, client arbbridge.ArbClient) (*executionChallengeWatcher, error) {
	bisectionChallenge, err := newBisectionChallengeWatcher(address, client)
	if err != nil {
		return nil, err
	}
	//executionContract, err := executionchallenge.NewExecutionChallenge(address, client)
	//if err != nil {
	//	return nil, errors2.Wrap(err, "Failed to connect to ChallengeManager")
	//}
	return &executionChallengeWatcher{
		bisectionChallengeWatcher: bisectionChallenge,
		challenge:                 nil,
		client:                    client,
		address:                   address,
	}, nil
}

func (c *executionChallengeWatcher) topics() []ethcommon.Hash {
	tops := []ethcommon.Hash{
		bisectedAssertionID,
		oneStepProofCompletedID,
	}
	return append(tops, c.bisectionChallengeWatcher.topics()...)
}

func (c *executionChallengeWatcher) StartConnection(ctx context.Context, startHeight *common.TimeBlocks, startLogIndex uint, errChan chan error, outChan chan arbbridge.Notification) error {
	//headers := make(chan *types.Header)
	//headersSub, err := c.client.SubscribeNewHead(ctx, headers)
	//if err != nil {
	//	return err
	//}

	//filter := ethereum.FilterQuery{
	//	Addresses: []ethcommon.Address{c.address},
	//	Topics:    [][]ethcommon.Hash{c.topics()},
	//}

	logChan := make(chan types.Log, 1024)
	logErrChan := make(chan error, 10)

	//if err := getLogs(ctx, c.client, filter, big.NewInt(0), logChan, logErrChan); err != nil {
	//	return err
	//}

	go func() {
		//defer headersSub.Unsubscribe()

		for {
			select {
			case <-ctx.Done():
				break
			case evmLog, ok := <-logChan:
				if !ok {
					errChan <- errors.New("logChan terminated early")
					return
				}
				if err := c.processEvents(ctx, evmLog, outChan); err != nil {
					errChan <- err
					return
				}
			case err := <-logErrChan:
				errChan <- err
				return
				//case err := <-headersSub.Err():
				//	errChan <- err
				//	return
			}
		}
	}()
	return nil
}

func (c *executionChallengeWatcher) processEvents(ctx context.Context, log types.Log, outChan chan arbbridge.Notification) error {
	event, err := func() (arbbridge.Event, error) {
		if log.Topics[0] == bisectedAssertionID {
			bisectChal, err := c.challenge.ParseBisectedAssertion(log)
			if err != nil {
				return nil, err
			}
			bisectionCount := len(bisectChal.MachineHashes) - 1
			assertions := make([]*valprotocol.ExecutionAssertionStub, 0, bisectionCount)
			for i := 0; i < bisectionCount; i++ {
				assertion := &valprotocol.ExecutionAssertionStub{
					AfterHash:        bisectChal.MachineHashes[i+1],
					DidInboxInsn:     bisectChal.DidInboxInsns[i],
					NumGas:           bisectChal.Gases[i],
					FirstMessageHash: bisectChal.MessageAccs[i],
					LastMessageHash:  bisectChal.MessageAccs[i+1],
					FirstLogHash:     bisectChal.LogAccs[i],
					LastLogHash:      bisectChal.LogAccs[i+1],
				}
				assertions = append(assertions, assertion)
			}
			return arbbridge.ExecutionBisectionEvent{
				Assertions: assertions,
				TotalSteps: bisectChal.TotalSteps,
				Deadline:   common.TimeTicks{Val: bisectChal.DeadlineTicks},
			}, nil
		} else if log.Topics[0] == oneStepProofCompletedID {
			_, err := c.challenge.ParseOneStepProofCompleted(log)
			if err != nil {
				return nil, err
			}
			return arbbridge.OneStepProofEvent{}, nil
		} else {
			event, err := c.bisectionChallengeWatcher.parseBisectionEvent(log)
			if event != nil || err != nil {
				return event, err
			}
		}
		return nil, errors2.New("unknown arbitrum event type")
	}()

	if err != nil {
		return err
	}

	if event == nil {
		return nil
	}
	//header, err := c.client.HeaderByHash(ctx, log.BlockHash)
	//if err != nil {
	//	return err
	//}
	outChan <- arbbridge.Notification{
		//BlockId: common.Hash{},
		Event:  event,
		TxHash: log.TxHash,
	}
	return nil
}
