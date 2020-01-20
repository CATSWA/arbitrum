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
	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
	"math/big"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/ethbridge/challengefactory"
)

type challengeFactory struct {
	contract *challengefactory.ChallengeFactory
	client   arbbridge.ArbClient
	auth     *TransOpts
}

func newChallengeFactory(address common.Address, client arbbridge.ArbClient, auth *TransOpts) (*challengeFactory, error) {
	//vmCreatorContract, err := challengefactory.NewChallengeFactory(address, client)
	//if err != nil {
	//	return nil, errors2.Wrap(err, "Failed to connect to arbFactory")
	//}
	return &challengeFactory{nil, client, auth}, nil
}

func (con *challengeFactory) CreateChallenge(
	ctx context.Context,
	asserter common.Address,
	challenger common.Address,
	challengePeriod common.TimeTicks,
	challengeHash common.Hash,
	challengeType *big.Int,
) (common.Address, error) {
	//con.auth.Context = ctx
	//tx, err := con.contract.CreateChallenge(
	//	con.auth,
	//	asserter.ToEthAddress(),
	//	challenger.ToEthAddress(),
	//	challengePeriod.Val,
	//	challengeHash,
	//	challengeType,
	//)
	//if err != nil {
	//	return common.Address{}, errors2.Wrap(err, "Failed to call to challengeFactory.CreateChallenge")
	//}
	//
	//receipt, err := WaitForReceiptWithResults(con.auth.Context, con.client, con.auth.From, tx, "CreateChallenge")
	//if err != nil {
	//	return common.Address{}, err
	//}
	//
	//if len(receipt.Logs) != 1 {
	//	return common.Address{}, errors2.New("Wrong receipt count")
	//}
	//
	//return common.NewAddressFromEth(receipt.Logs[0].Address), nil
	return *new(common.Address), nil
}
