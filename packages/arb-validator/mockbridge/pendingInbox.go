/*
 * Copyright 2019-2020, Offchain Labs, Inc.
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
	"github.com/offchainlabs/arbitrum/packages/arb-validator/structures"
	"math/big"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-util/value"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/valprotocol"
)

type PendingInbox struct {
	//GlobalPendingInbox *globalpendinginbox.GlobalPendingInbox
	pending *structures.PendingInbox
	client  *MockArbClient
}

func newPendingInbox(address common.Address, client *MockArbClient) (*PendingInbox, error) {
	//globalPendingInboxContract, err := globalpendinginbox.NewGlobalPendingInbox(address, client)
	//if err != nil {
	//	return nil, errors2.Wrap(err, "Failed to connect to GlobalPendingInbox")
	//}
	client.MockEthClient.pending[address] = &PendingInbox{client: client}
	return &PendingInbox{structures.NewPendingInbox(), client}, nil
}

func (con *PendingInbox) SendMessage(
	ctx context.Context,
	msg valprotocol.Message,
) error {
	//var dataBuf bytes.Buffer
	//if err := value.MarshalValue(msg.Data, &dataBuf); err != nil {
	//	return err
	//}
	//tx, err := con.GlobalPendingInbox.SendMessage(
	//	auth,
	//	msg.Destination,
	//	msg.TokenType,
	//	msg.Currency,
	//	dataBuf.Bytes(),
	//)
	//if err != nil {
	//	return err
	//}
	//return waitForReceipt(auth.Context, con.client, auth.From, tx, "SendMessage")
	// if token type NFT
	if msg.TokenType[20] == 1 {
		//transferNFT()
	} else {
		//transferToken()
	}
	//
	// if send successful
	//
	//PendingInbox storage pendingInbox = pending[_destination];
	//if (pendingInbox.value != 0) {
	//	bytes32 dataHash = Value.deserializeHashed(_data);
	//	bytes32 txHash = keccak256(
	//		abi.encodePacked(
	//			_destination,
	//			dataHash,
	//			_value,
	//			_tokenType
	//	)
	//);
	//	Value.Data[] memory dataValues = new Value.Data[](3);
	//	dataValues[0] = Value.newHashOnly(dataHash);
	//	dataValues[1] = Value.newInt(block.number);
	//	dataValues[2] = Value.newInt(uint(txHash));
	//
	//	Value.Data[] memory values = new Value.Data[](4);
	//	values[0] = Value.newTuple(dataValues);
	//	values[1] = Value.newInt(uint256(_sender));
	//	values[2] = Value.newInt(_value);
	//	values[3] = Value.newInt(uint256(bytes32(_tokenType)));
	//	bytes32 messageHash =  Value.newTuple(values).hash().hash;
	//
	//	pendingInbox.value = Protocol.addMessageToPending(
	//		pendingInbox.value,
	//		messageHash
	//	);
	//	pendingInbox.count++;
	//}
	con.pending.DeliverMessage(msg.AsValue())
	//
	//emit IGlobalPendingInbox.MessageDelivered(
	//	_destination,
	//	_sender,
	//	_tokenType,
	//	_value,
	//	_data
	//)
	con.client.MockEthClient.pubMsg(arbbridge.MaybeEvent{
		Event: arbbridge.MessageDeliveredEvent{
			ChainInfo: arbbridge.ChainInfo{
				BlockId: con.client.MockEthClient.LatestBlock,
			},
			Msg: msg,
		},
	})
	return nil
}

func (con *PendingInbox) ForwardMessage(
	ctx context.Context,
	msg valprotocol.Message,
	sig []byte,
) error {
	//var dataBuf bytes.Buffer
	//if err := value.MarshalValue(msg.Data, &dataBuf); err != nil {
	//	return err
	//}
	//tx, err := con.GlobalPendingInbox.ForwardMessage(
	//	auth,
	//	msg.Destination,
	//	msg.TokenType,
	//	msg.Currency,
	//	dataBuf.Bytes(),
	//	sig,
	//)
	//if err != nil {
	//	return err
	//}
	//return waitForReceipt(auth.Context, con.client, auth.From, tx, "ForwardMessage")
	return nil
}

func (con *PendingInbox) SendEthMessage(
	ctx context.Context,
	data value.Value,
	destination common.Address,
	amount *big.Int,
) error {
	//var dataBuf bytes.Buffer
	//if err := value.MarshalValue(data, &dataBuf); err != nil {
	//	return 0, err
	//}
	//tx, err := con.GlobalPendingInbox.SendEthMessage(
	//	&bind.TransactOpts{
	//		From:     auth.From,
	//		Signer:   auth.Signer,
	//		GasLimit: auth.GasLimit,
	//		Value:    amount,
	//	},
	//	destination,
	//	dataBuf.Bytes(),
	//)
	//if err != nil {
	//	return 0, err
	//}
	//receipt, err := waitForReceiptWithResults(auth.Context, con.client, auth.From, tx, "SendEthMessage")
	//return receipt.Status, err

	//depositEth(_destination);
	con.SendMessage(ctx, valprotocol.Message{
		Data:        data,
		TokenType:   [21]byte{}, //???
		Currency:    amount,
		Destination: destination,
	})
	return nil
}

func (con *PendingInbox) DepositFunds(ctx context.Context, amount *big.Int, dest common.Address) error {
	//tx, err := con.GlobalPendingInbox.DepositEth(
	//	&bind.TransactOpts{
	//		From:     auth.From,
	//		Signer:   auth.Signer,
	//		GasLimit: auth.GasLimit,
	//		Value:    amount,
	//	},
	//	dest,
	//)
	//if err != nil {
	//	return err
	//}
	//return waitForReceipt(auth.Context, con.client, auth.From, tx, "DepositFunds")
	return nil
}

func (con *PendingInbox) GetTokenBalance(
	ctx context.Context,
	user common.Address,
	tokenContract common.Address,
) (*big.Int, error) {
	//return con.GlobalPendingInbox.GetTokenBalance(
	//	auth,
	//	tokenContract,
	//	user,
	//)
	return big.NewInt(0), nil
}
