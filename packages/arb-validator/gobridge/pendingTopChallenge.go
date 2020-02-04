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

package gobridge

import (
	"context"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/structures"
	"math/big"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
)

type pendingTopChallenge struct {
	*bisectionChallenge
}

func NewPendingTopChallenge(address common.Address, client *MockArbAuthClient) (*pendingTopChallenge, error) {
	bisectionChallenge, err := newBisectionChallenge(address, client)
	if err != nil {
		return nil, err
	}
	vm := &pendingTopChallenge{bisectionChallenge: bisectionChallenge}
	err = vm.setupContracts()
	return vm, err
}

func (c *pendingTopChallenge) setupContracts() error {
	//challengeManagerContract, err := pendingtopchallenge.NewPendingTopChallenge(c.address, c.Client)
	//if err != nil {
	//	return errors2.Wrap(err, "Failed to connect to MessagesChallenge")
	//}
	//
	//c.challenge = challengeManagerContract
	return nil
}

func (vm *pendingTopChallenge) GetEvents(ctx context.Context, blockId *structures.BlockId) ([]arbbridge.Event, error) {
	return nil, nil
}

func (c *pendingTopChallenge) Bisect(
	ctx context.Context,
	chainHashes []common.Hash,
	chainLength *big.Int,
) error {
	//c.auth.Context = ctx
	//tx, err := c.challenge.Bisect(
	//	c.auth,
	//	chainHashes,
	//	chainLength,
	//)
	//if err != nil {
	//	return err
	//}
	//return c.waitForReceipt(ctx, tx, "Bisect")
	return nil
}

func (c *pendingTopChallenge) OneStepProof(
	ctx context.Context,
	lowerHashA common.Hash,
	value common.Hash,
) error {
	//c.auth.Context = ctx
	//tx, err := c.challenge.OneStepProof(
	//	c.auth,
	//	lowerHashA,
	//	topHashA,
	//	value,
	//)
	//if err != nil {
	//	return err
	//}
	//return c.waitForReceipt(ctx, tx, "OneStepProof")
	return nil
}

func (c *pendingTopChallenge) ChooseSegment(
	ctx context.Context,
	assertionToChallenge uint16,
	chainHashes []common.Hash,
	chainLength uint64,
) error {
	bisectionCount := uint64(len(chainHashes) - 1)
	bisectionHashes := make([]common.Hash, 0, bisectionCount)
	for i := uint64(0); i < bisectionCount; i++ {
		stepCount := structures.CalculateBisectionStepCount(i, bisectionCount, chainLength)
		bisectionHashes = append(
			bisectionHashes,
			structures.PendingTopChallengeDataHash(
				chainHashes[i],
				chainHashes[i+1],
				new(big.Int).SetUint64(uint64(stepCount)),
			),
		)
	}
	return c.bisectionChallenge.chooseSegment(
		ctx,
		assertionToChallenge,
		bisectionHashes,
	)
	return nil
}
