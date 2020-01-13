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

	"github.com/offchainlabs/arbitrum/packages/arb-validator/structures"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/offchainlabs/arbitrum/packages/arb-util/value"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/ethbridge/rollup"
)

type ArbRollup struct {
	Client    arbbridge.ArbClient
	ArbRollup *rollup.ArbRollup
	auth      *bind.TransactOpts
}

func NewRollup(address common.Address, client arbbridge.ArbClient, auth *bind.TransactOpts) (*ArbRollup, error) {
	//arbitrumRollupContract, err := rollup.NewArbRollup(address, client.(*MockArbClient).client)
	//if err != nil {
	//	return nil, errors2.Wrap(err, "Failed to connect to ArbRollup")
	//}
	//vm := &ArbRollup{Client: client.(*MockArbClient).client, ArbRollup: arbitrumRollupContract, auth: auth}
	return &ArbRollup{
		Client:    nil,
		ArbRollup: nil,
		auth:      nil,
	}, nil
}

func (vm *ArbRollup) PlaceStake(ctx context.Context, stakeAmount *big.Int, proof1 [][32]byte, proof2 [][32]byte) error {
	//call := &bind.TransactOpts{
	//	From:    vm.auth.From,
	//	Signer:  vm.auth.Signer,
	//	Value:   stakeAmount,
	//	Context: ctx,
	//}
	//tx, err := vm.ArbRollup.PlaceStake(
	//	call,
	//	proof1,
	//	proof2,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "PlaceStake")
	return nil
}

func (vm *ArbRollup) RecoverStakeConfirmed(ctx context.Context, proof [][32]byte) error {
	//vm.auth.Context = ctx
	//tx, err := vm.ArbRollup.RecoverStakeConfirmed(
	//	vm.auth,
	//	proof,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "RecoverStakeConfirmed")
	return nil
}

func (vm *ArbRollup) RecoverStakeOld(ctx context.Context, staker common.Address, proof [][32]byte) error {
	//vm.auth.Context = ctx
	//tx, err := vm.ArbRollup.RecoverStakeOld(
	//	vm.auth,
	//	staker,
	//	proof,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "RecoverStakeOld")
	return nil
}

func (vm *ArbRollup) RecoverStakeMooted(ctx context.Context, nodeHash [32]byte, staker common.Address, latestConfirmedProof [][32]byte, stakerProof [][32]byte) error {
	//vm.auth.Context = ctx
	//tx, err := vm.ArbRollup.RecoverStakeMooted(
	//	vm.auth,
	//	staker,
	//	nodeHash,
	//	latestConfirmedProof,
	//	stakerProof,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "RecoverStakeMooted")
	return nil
}

func (vm *ArbRollup) RecoverStakePassedDeadline(ctx context.Context, stakerAddress common.Address, deadlineTicks *big.Int, disputableNodeHashVal [32]byte, childType uint64, vmProtoStateHash [32]byte, proof [][32]byte) error {
	//vm.auth.Context = ctx
	//tx, err := vm.ArbRollup.RecoverStakePassedDeadline(
	//	vm.auth,
	//	stakerAddress,
	//	deadlineTicks,
	//	disputableNodeHashVal,
	//	new(big.Int).SetUint64(childType),
	//	vmProtoStateHash,
	//	proof,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "RecoverStakePassedDeadline")
	return nil
}

func (vm *ArbRollup) MoveStake(ctx context.Context, proof1 [][32]byte, proof2 [][32]byte) error {
	//vm.auth.Context = ctx
	//tx, err := vm.ArbRollup.MoveStake(
	//	vm.auth,
	//	proof1,
	//	proof2,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "MoveStake")
	return nil
}

func (vm *ArbRollup) PruneLeaf(ctx context.Context, from [32]byte, proof1 [][32]byte, proof2 [][32]byte) error {
	//vm.auth.Context = ctx
	//tx, err := vm.ArbRollup.PruneLeaf(
	//	vm.auth,
	//	from,
	//	proof1,
	//	proof2,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "PruneLeaf")
	return nil
}

func (vm *ArbRollup) MakeAssertion(
	ctx context.Context,

	prevPrevLeafHash [32]byte,
	prevDataHash [32]byte,
	prevDeadline structures.TimeTicks,
	prevChildType structures.ChildType,

	beforeState *structures.VMProtoData,
	assertionParams *structures.AssertionParams,
	assertionClaim *structures.AssertionClaim,
	stakerProof [][32]byte,
) error {
	//vm.auth.Context = ctx
	//tx, err := vm.ArbRollup.MakeAssertion(
	//	vm.auth,
	//	[9][32]byte{
	//		beforeState.MachineHash,
	//		beforeState.PendingTop,
	//		prevPrevLeafHash,
	//		prevDataHash,
	//		assertionClaim.AfterPendingTop,
	//		assertionClaim.ImportedMessagesSlice,
	//		assertionClaim.AssertionStub.AfterHashValue(),
	//		assertionClaim.AssertionStub.LastMessageHashValue(),
	//		assertionClaim.AssertionStub.LastLogHashValue(),
	//	},
	//	beforeState.PendingCount,
	//	prevDeadline.Val,
	//	uint32(prevChildType),
	//	assertionParams.NumSteps,
	//	assertionParams.TimeBounds.AsIntArray(),
	//	assertionParams.ImportedMessageCount,
	//	assertionClaim.AssertionStub.DidInboxInsn,
	//	assertionClaim.AssertionStub.NumGas,
	//	stakerProof,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "MakeAssertion")
	return nil
}

func (vm *ArbRollup) ConfirmValid(
	ctx context.Context,
	deadline structures.TimeTicks,
	outMsgs []value.Value,
	logsAccHash [32]byte,
	protoHash [32]byte,
	stakerAddresses []common.Address,
	stakerProofs [][32]byte,
	stakerProofOffsets []*big.Int,
) error {
	//vm.auth.Context = ctx
	//messages := hashing.CombineMessages(outMsgs)
	//tx, err := vm.ArbRollup.ConfirmValid(
	//	vm.auth,
	//	deadline.Val,
	//	messages,
	//	logsAccHash,
	//	protoHash,
	//	stakerAddresses,
	//	stakerProofs,
	//	stakerProofOffsets,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "ConfirmValid")
	return nil
}

func (vm *ArbRollup) ConfirmInvalid(
	ctx context.Context,
	deadline structures.TimeTicks,
	challengeNodeData [32]byte,
	branch structures.ChildType,
	protoHash [32]byte,
	stakerAddresses []common.Address,
	stakerProofs [][32]byte,
	stakerProofOffsets []*big.Int,
) error {
	//vm.auth.Context = ctx
	//tx, err := vm.ArbRollup.ConfirmInvalid(
	//	vm.auth,
	//	deadline.Val,
	//	challengeNodeData,
	//	new(big.Int).SetUint64(uint64(branch)),
	//	protoHash,
	//	stakerAddresses,
	//	stakerProofs,
	//	stakerProofOffsets,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "ConfirmInvalid")
	return nil
}

func (vm *ArbRollup) StartChallenge(ctx context.Context, asserterAddress common.Address, challengerAddress common.Address, prevNode [32]byte, disputableDeadline *big.Int, asserterPosition structures.ChildType, challengerPosition structures.ChildType, asserterVMProtoHash [32]byte, challengerVMProtoHash [32]byte, asserterProof [][32]byte, challengerProof [][32]byte, asserterDataHash [32]byte, asserterPeriodTicks structures.TimeTicks, challengerNodeHash [32]byte) error {
	//vm.auth.Context = ctx
	//tx, err := vm.ArbRollup.StartChallenge(
	//	vm.auth,
	//	asserterAddress,
	//	challengerAddress,
	//	prevNode,
	//	disputableDeadline,
	//	[2]*big.Int{
	//		new(big.Int).SetUint64(uint64(asserterPosition)),
	//		new(big.Int).SetUint64(uint64(challengerPosition)),
	//	},
	//	[2][32]byte{
	//		asserterVMProtoHash,
	//		challengerVMProtoHash,
	//	},
	//	asserterProof,
	//	challengerProof,
	//	asserterDataHash,
	//	asserterPeriodTicks.Val,
	//	challengerNodeHash,
	//)
	//if err != nil {
	//	return err
	//}
	//return vm.waitForReceipt(ctx, tx, "StartExecutionChallenge")
	return nil
}

//func (vm *ArbRollup) VerifyVM(
//	auth *bind.CallOpts,
//	config *valmessage.VMConfiguration,
//	machine [32]byte,
//) error {
//	//code, err := vm.contract.Client.CodeAt(auth.Context, vm.address, nil)
//	// Verify that VM has correct code
//	vmInfo, err := vm.ArbRollup.Vm(auth)
//	if err != nil {
//		return err
//	}
//
//	if vmInfo.MachineHash != machine {
//		return errors.New("VM has different machine hash")
//	}
//
//	if config.GracePeriod != uint64(vmInfo.GracePeriod) {
//		return errors.New("VM has different grace period")
//	}
//
//	if value.NewBigIntFromBuf(config.EscrowRequired).Cmp(vmInfo.EscrowRequired) != 0 {
//		return errors.New("VM has different escrow required")
//	}
//
//	if config.MaxExecutionStepCount != vmInfo.MaxExecutionSteps {
//		return errors.New("VM has different mxa steps")
//	}
//
//	owner, err := vm.ArbRollup.Owner(auth)
//	if err != nil {
//		return err
//	}
//	if protocol.NewAddressFromBuf(config.Owner) != owner {
//		return errors.New("VM has different owner")
//	}
//	return nil
//}

//func (vm *ArbRollup) waitForReceipt(ctx context.Context, tx *types.Transaction, methodName string) error {
//	return waitForReceipt(ctx, vm.Client, vm.auth.From, tx, methodName)
//}
