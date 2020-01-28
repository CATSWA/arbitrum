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

package rollup

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/offchainlabs/arbitrum/packages/arb-validator/challenges"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-util/protocol"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/structures"
)

type ChainListener interface {
	StakeCreated(context.Context, *ChainObserver, arbbridge.StakeCreatedEvent)
	StakeRemoved(context.Context, *ChainObserver, arbbridge.StakeRefundedEvent)
	StakeMoved(context.Context, *ChainObserver, arbbridge.StakeMovedEvent)
	StartedChallenge(context.Context, *ChainObserver, arbbridge.ChallengeStartedEvent, *Node, *Node)
	CompletedChallenge(context.Context, *ChainObserver, arbbridge.ChallengeCompletedEvent)
	SawAssertion(context.Context, *ChainObserver, arbbridge.AssertedEvent)
	ConfirmedNode(context.Context, *ChainObserver, arbbridge.ConfirmedEvent)
	PrunedLeaf(context.Context, *ChainObserver, arbbridge.PrunedEvent)
	MessageDelivered(context.Context, *ChainObserver, arbbridge.MessageDeliveredEvent)

	AssertionPrepared(context.Context, *ChainObserver, *preparedAssertion)
	ValidNodeConfirmable(context.Context, *ChainObserver, *confirmValidOpportunity)
	InvalidNodeConfirmable(context.Context, *ChainObserver, *confirmInvalidOpportunity)
	PrunableLeafs(context.Context, *ChainObserver, []pruneParams)
	MootableStakes(context.Context, *ChainObserver, []recoverStakeMootedParams)
	OldStakes(context.Context, *ChainObserver, []recoverStakeOldParams)

	AdvancedCalculatedValidNode(context.Context, *ChainObserver, common.Hash)
	AdvancedKnownAssertion(context.Context, *ChainObserver, *protocol.ExecutionAssertion, common.Hash)
}

type StakingKey struct {
	client   arbbridge.ArbAuthClient
	contract arbbridge.ArbRollup
}

type ValidatorChainListener struct {
	sync.Mutex
	actor                  arbbridge.ArbRollup
	rollupAddress          common.Address
	stakingKeys            map[common.Address]*StakingKey
	broadcastAssertions    map[common.Hash]*structures.AssertionParams
	broadcastConfirmations map[common.Hash]bool
	broadcastLeafPrunes    map[common.Hash]bool
	broadcastCreateStakes  map[common.Address]*common.TimeBlocks
}

func NewValidatorChainListener(rollupAddress common.Address, actor arbbridge.ArbRollup) *ValidatorChainListener {
	return &ValidatorChainListener{
		actor:                  actor,
		rollupAddress:          rollupAddress,
		stakingKeys:            make(map[common.Address]*StakingKey),
		broadcastAssertions:    make(map[common.Hash]*structures.AssertionParams),
		broadcastConfirmations: make(map[common.Hash]bool),
		broadcastLeafPrunes:    make(map[common.Hash]bool),
		broadcastCreateStakes:  make(map[common.Address]*common.TimeBlocks),
	}
}

func stakeLatestValid(ctx context.Context, chain *ChainObserver, stakingKey *StakingKey) error {
	location := chain.knownValidNode
	proof1 := GeneratePathProof(chain.nodeGraph.latestConfirmed, location)
	proof2 := GeneratePathProof(location, chain.nodeGraph.getLeaf(location))
	stakeAmount := chain.nodeGraph.params.StakeRequirement

	log.Println("Placing stake for", stakingKey.client.Address())
	return stakingKey.contract.PlaceStake(ctx, stakeAmount, proof1, proof2)
}

func (lis *ValidatorChainListener) AddStaker(client arbbridge.ArbAuthClient) error {
	contract, err := client.NewRollup(lis.rollupAddress)
	if err != nil {
		return err
	}

	address := client.Address()
	lis.stakingKeys[address] = &StakingKey{
		client:   client,
		contract: contract,
	}
	return nil
}

func makeAssertion(ctx context.Context, rollup arbbridge.ArbRollup, prepared *preparedAssertion, proof []common.Hash) error {
	return rollup.MakeAssertion(
		ctx,
		prepared.prevPrevLeafHash,
		prepared.prevDataHash,
		prepared.prevDeadline,
		prepared.prevChildType,
		prepared.beforeState,
		prepared.params,
		prepared.claim,
		proof,
	)
}

func (lis *ValidatorChainListener) AssertionPrepared(ctx context.Context, chain *ChainObserver, prepared *preparedAssertion) {
	// Anyone confirm a node
	// No need to have your own stake
	fmt.Println("****************AssertionPrepared top")
	lis.Lock()
	prevParams, alreadySent := lis.broadcastAssertions[prepared.leafHash]
	lis.Unlock()
	if alreadySent && prevParams.Equals(prepared.params) {
		return
	}
	fmt.Println("****************AssertionPrepared not already sent")

	leaf, ok := chain.nodeGraph.nodeFromHash[prepared.leafHash]
	if !ok {
		log.Println("Prepared assertion on top of invalid node")
		return
	}

	for stakingAddress, stakingKey := range lis.stakingKeys {
		fmt.Println("****************AssertionPrepared staking keys loop")
		stakerPos := chain.nodeGraph.stakers.Get(stakingAddress)
		if stakerPos == nil {
			// stakingKey is not staked
			continue
		}
		proof := GeneratePathProof(stakerPos.location, leaf)
		if proof == nil {
			// staker can't move to new asertion
			continue
		}
		lis.Lock()
		lis.broadcastAssertions[prepared.leafHash] = prepared.params
		lis.Unlock()
		log.Printf("*******************%v is making an assertion\n", stakingAddress)
		go func() {
			err := makeAssertion(ctx, stakingKey.contract, prepared.Clone(), proof)
			if err != nil {
				log.Println("Error making assertion", err)
				lis.Lock()
				delete(lis.broadcastAssertions, prepared.leafHash)
				lis.Unlock()
			} else {
				log.Println("Successfully made assertion")
			}
		}()
		return
	}

	log.Println("Maybe putting down stake")
	for stakingAddress, stakingKey := range lis.stakingKeys {
		stakerPos := chain.nodeGraph.stakers.Get(stakingAddress)
		if stakerPos != nil {
			// stakingKey is already down
			continue
		}
		lis.Lock()
		stakeTime, placedStake := lis.broadcastCreateStakes[stakingAddress]
		if placedStake {
			log.Println("Thinking about placing stake", chain.latestBlockId.Height.AsInt(), new(big.Int).Add(stakeTime.AsInt(), big.NewInt(3)))
		}
		if !placedStake || chain.latestBlockId.Height.AsInt().Cmp(new(big.Int).Add(stakeTime.AsInt(), big.NewInt(3))) >= 0 {
			lis.broadcastCreateStakes[stakingAddress] = chain.latestBlockId.Height
			log.Println("No stake is currently down, so setting up a stake")
			lis.Unlock()
			// Put down new stake so that we can assert next time
			go func() {
				err := stakeLatestValid(ctx, chain, stakingKey)
				if err != nil {
					lis.Lock()
					delete(lis.broadcastCreateStakes, stakingAddress)
					lis.Unlock()
					log.Println("Error placing stake", err)
				}
			}()
			return
		} else {
			lis.Unlock()
		}
	}
}

func (lis *ValidatorChainListener) initiateChallenge(ctx context.Context, opp *challengeOpportunity) error {
	return lis.actor.StartChallenge(
		ctx,
		opp.asserter,
		opp.challenger,
		opp.prevNodeHash,
		opp.deadlineTicks.Val,
		opp.asserterNodeType,
		opp.challengerNodeType,
		opp.asserterVMProtoHash,
		opp.challengerVMProtoHash,
		opp.asserterProof,
		opp.challengerProof,
		opp.asserterNodeHash,
		opp.challengerDataHash,
		opp.challengerPeriodTicks,
	)
}

func (lis *ValidatorChainListener) StakeCreated(ctx context.Context, chain *ChainObserver, ev arbbridge.StakeCreatedEvent) {
	_, ok := lis.stakingKeys[ev.Staker]
	if ok {
		staker := chain.nodeGraph.stakers.Get(ev.Staker)
		if staker == nil {
			panic("Stake created but address is not in graph")
		}
		opp := chain.nodeGraph.checkChallengeOpportunityAny(staker)
		if opp != nil {
			go lis.initiateChallenge(ctx, opp)
		}
	} else {
		lis.challengeStakerIfPossible(ctx, chain, ev.Staker)
	}
}

func (lis *ValidatorChainListener) StakeMoved(ctx context.Context, chain *ChainObserver, ev arbbridge.StakeMovedEvent) {
	lis.challengeStakerIfPossible(ctx, chain, ev.Staker)
}

func (lis *ValidatorChainListener) challengeStakerIfPossible(ctx context.Context, chain *ChainObserver, stakerAddr common.Address) {
	_, ok := lis.stakingKeys[stakerAddr]
	if ok {
		// Don't challenge yourself
		return
	}

	newStaker := chain.nodeGraph.stakers.Get(stakerAddr)
	if newStaker == nil {
		log.Fatalf("Nonexistant staker moved %v", stakerAddr)
	}

	// Search for an already staked staking key
	for myAddr, _ := range lis.stakingKeys {
		meAsStaker := chain.nodeGraph.stakers.Get(myAddr)
		if meAsStaker == nil {
			continue
		}
		opp := chain.nodeGraph.checkChallengeOpportunityPair(newStaker, meAsStaker)
		if opp != nil {
			go lis.initiateChallenge(ctx, opp)
			return
		}
	}
	opp := chain.nodeGraph.checkChallengeOpportunityAny(newStaker)
	if opp != nil {
		go lis.initiateChallenge(ctx, opp)
		return
	}
}

// All functions below are either only called if you have a stake down, or don't require a stake

func (lis *ValidatorChainListener) StartedChallenge(ctx context.Context, chain *ChainObserver, ev arbbridge.ChallengeStartedEvent, conflictNode *Node, challengerAncestor *Node) {
	// Must already be staked to be challenged
	startBlockId := ev.BlockId
	startLogIndex := ev.LogIndex - 1
	asserterKey, ok := lis.stakingKeys[ev.Asserter]
	if ok {
		switch conflictNode.linkType {
		case structures.InvalidPendingChildType:
			go challenges.DefendPendingTopClaim(
				ctx,
				asserterKey.client,
				ev.ChallengeContract,
				startBlockId,
				startLogIndex,
				chain.pendingInbox.MessageStack,
				conflictNode.disputable.AssertionClaim.AfterPendingTop,
				conflictNode.disputable.MaxPendingTop,
				100,
			)
		case structures.InvalidMessagesChildType:
			go challenges.DefendMessagesClaim(
				ctx,
				asserterKey.client,
				ev.ChallengeContract,
				startBlockId,
				startLogIndex,
				chain.pendingInbox.MessageStack,
				conflictNode.vmProtoData.PendingTop,
				conflictNode.disputable.AssertionClaim.AfterPendingTop,
				conflictNode.disputable.AssertionClaim.ImportedMessagesSlice,
				100,
			)
		case structures.InvalidExecutionChildType:
			go challenges.DefendExecutionClaim(
				ctx,
				asserterKey.client,
				ev.ChallengeContract,
				startBlockId,
				startLogIndex,
				chain.executionPrecondition(conflictNode),
				conflictNode.machine,
				conflictNode.disputable.AssertionParams.NumSteps,
				50,
			)
		}
	}

	challenger, ok := lis.stakingKeys[ev.Challenger]
	if ok {
		switch conflictNode.linkType {
		case structures.InvalidPendingChildType:
			go challenges.ChallengePendingTopClaim(
				ctx,
				challenger.client,
				ev.ChallengeContract,
				startBlockId,
				startLogIndex,
				chain.pendingInbox.MessageStack,
			)
		case structures.InvalidMessagesChildType:
			go challenges.ChallengeMessagesClaim(
				ctx,
				challenger.client,
				ev.ChallengeContract,
				startBlockId,
				startLogIndex,
				chain.pendingInbox.MessageStack,
				conflictNode.vmProtoData.PendingTop,
				conflictNode.disputable.AssertionClaim.AfterPendingTop,
			)
		case structures.InvalidExecutionChildType:
			go challenges.ChallengeExecutionClaim(
				ctx,
				challenger.client,
				ev.ChallengeContract,
				startBlockId,
				startLogIndex,
				chain.executionPrecondition(conflictNode),
				conflictNode.machine,
				false,
			)
		}
	}
}

func (lis *ValidatorChainListener) CompletedChallenge(ctx context.Context, chain *ChainObserver, ev arbbridge.ChallengeCompletedEvent) {
	// Must be staked to have challenge completed
	_, ok := lis.stakingKeys[ev.Winner]
	if ok {
		lis.wonChallenge(ev)
	}

	_, ok = lis.stakingKeys[ev.Loser]
	if ok {
		lis.lostChallenge(ev)
	}
	lis.challengeStakerIfPossible(ctx, chain, ev.Winner)
}

func (lis *ValidatorChainListener) ValidNodeConfirmable(ctx context.Context, observer *ChainObserver, conf *confirmValidOpportunity) {
	// Anyone confirm a node
	// No need to have your own stake
	_, alreadySent := lis.broadcastConfirmations[conf.nodeHash]
	if alreadySent {
		return
	}
	lis.broadcastConfirmations[conf.nodeHash] = true
	go func() {
		lis.actor.ConfirmValid(
			ctx,
			conf.deadlineTicks,
			conf.messages,
			conf.logsAcc,
			conf.vmProtoStateHash,
			conf.stakerAddresses,
			conf.stakerProofs,
			conf.stakerProofOffsets,
		)
	}()
}

func (lis *ValidatorChainListener) InvalidNodeConfirmable(ctx context.Context, observer *ChainObserver, conf *confirmInvalidOpportunity) {
	// Anyone confirm a node
	// No need to have your own stake
	_, alreadySent := lis.broadcastConfirmations[conf.nodeHash]
	if alreadySent {
		return
	}
	lis.broadcastConfirmations[conf.nodeHash] = true
	go func() {
		lis.actor.ConfirmInvalid(
			ctx,
			conf.deadlineTicks,
			conf.challengeNodeData,
			conf.branch,
			conf.vmProtoStateHash,
			conf.stakerAddresses,
			conf.stakerProofs,
			conf.stakerProofOffsets,
		)
	}()
}

func (lis *ValidatorChainListener) PrunableLeafs(ctx context.Context, observer *ChainObserver, params []pruneParams) {
	// Anyone can prune a leaf
	for _, prune := range params {
		_, alreadySent := lis.broadcastLeafPrunes[prune.leafHash]
		if alreadySent {
			continue
		}
		lis.broadcastLeafPrunes[prune.leafHash] = true
		pruneCopy := prune.Clone()
		go func() {
			lis.actor.PruneLeaf(
				ctx,
				pruneCopy.ancestorHash,
				pruneCopy.leafProof,
				pruneCopy.ancProof,
			)
		}()
	}
}

func (lis *ValidatorChainListener) MootableStakes(ctx context.Context, observer *ChainObserver, params []recoverStakeMootedParams) {
	// Anyone can moot any stake
	for _, moot := range params {
		go func() {
			lis.actor.RecoverStakeMooted(
				ctx,
				moot.ancestorHash,
				moot.addr,
				moot.lcProof,
				moot.stProof,
			)
		}()
	}
}

func (lis *ValidatorChainListener) OldStakes(ctx context.Context, observer *ChainObserver, params []recoverStakeOldParams) {
	// Anyone can remove an old stake
	for _, old := range params {
		go func() {
			lis.actor.RecoverStakeOld(
				ctx,
				old.addr,
				old.proof,
			)
		}()
	}
}

func (lis *ValidatorChainListener) AdvancedCalculatedValidNode(ctx context.Context, chain *ChainObserver, nodeHash common.Hash) {
	for stakingAddress, _ := range lis.stakingKeys {
		staker := chain.nodeGraph.stakers.idx[stakingAddress]
		newValidNode := chain.nodeGraph.nodeFromHash[nodeHash]
		if newValidNode.depth > staker.location.depth {
			proof1 := GeneratePathProof(staker.location, newValidNode)
			proof2 := GeneratePathProof(newValidNode, chain.nodeGraph.getLeaf(newValidNode))
			lis.actor.MoveStake(ctx, proof1, proof2)
		}
	}
}

func (lis *ValidatorChainListener) StakeRemoved(context.Context, *ChainObserver, arbbridge.StakeRefundedEvent) {
}
func (lis *ValidatorChainListener) lostChallenge(arbbridge.ChallengeCompletedEvent) {}
func (lis *ValidatorChainListener) wonChallenge(arbbridge.ChallengeCompletedEvent)  {}
func (lis *ValidatorChainListener) SawAssertion(context.Context, *ChainObserver, arbbridge.AssertedEvent) {
}
func (lis *ValidatorChainListener) ConfirmedNode(context.Context, *ChainObserver, arbbridge.ConfirmedEvent) {
}
func (lis *ValidatorChainListener) PrunedLeaf(context.Context, *ChainObserver, arbbridge.PrunedEvent) {
}
func (lis *ValidatorChainListener) MessageDelivered(context.Context, *ChainObserver, arbbridge.MessageDeliveredEvent) {
}
func (lis *ValidatorChainListener) AdvancedKnownAssertion(context.Context, *ChainObserver, *protocol.ExecutionAssertion, common.Hash) {
}
