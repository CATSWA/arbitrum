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

package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-util/hashing"
	"github.com/offchainlabs/arbitrum/packages/arb-util/value"
)

type TransactionBatch struct {
	Chain  common.Address
	TxData []byte
}

// BatchTxHash hashes the transaction data. This hash is signed by users
// who submit transactions as part of a batch
func BatchTxHash(
	rollupAddress common.Address,
	to common.Address,
	sequenceNum *big.Int,
	val *big.Int,
	data []byte,
) common.Hash {
	return hashing.SoliditySHA3(
		hashing.Address(rollupAddress),
		hashing.Address(to),
		hashing.Uint256(sequenceNum),
		hashing.Uint256(val),
		hashing.Bytes32(hashing.SoliditySHA3(data)),
	)
}

type BatchTx struct {
	To     common.Address
	SeqNum *big.Int
	Value  *big.Int
	Data   []byte
	Sig    [65]byte
}

func NewBatchTxFromData(data []byte, offset int) (BatchTx, error) {
	dataLength := int(binary.BigEndian.Uint16(data[offset : offset+2]))
	if offset+DataOffset+dataLength < len(data) {
		return BatchTx{}, errors.New("not enough data remaining")
	}
	offset += 2
	toRaw := data[offset : offset+20]
	var to common.Address
	copy(to[:], toRaw)
	offset += 20
	seqRaw := data[offset : offset+32]
	seq := new(big.Int).SetBytes(seqRaw)
	offset += 32
	valueRaw := data[offset : offset+32]
	val := new(big.Int).SetBytes(valueRaw)
	offset += 32
	sigRaw := data[offset : offset+65]
	var sig [65]byte
	copy(sig[:], sigRaw)
	offset += 65
	txData := data[offset : offset+dataLength]

	return BatchTx{
		To:     to,
		SeqNum: seq,
		Value:  val,
		Data:   txData,
		Sig:    sig,
	}, nil
}

func (b BatchTx) encodedLength() int {
	return DataOffset + len(b.Data)
}

func (b BatchTx) ToBytes() []byte {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data[:], uint16(len(b.Data)))
	data = append(data, b.To[:]...)
	data = append(data, math.U256Bytes(b.SeqNum)...)
	data = append(data, math.U256Bytes(b.Value)...)
	data = append(data, b.Sig[:]...)
	data = append(data, b.Data...)
	return data
}

var DataOffset = 151

func (m DeliveredTransactionBatch) getTransactions() []DeliveredTransaction {
	txes := make([]DeliveredTransaction, 0)
	offset := 0

	data := m.TxData
	for offset+DataOffset < len(data) {
		batch, err := NewBatchTxFromData(data, offset)
		if err != nil {
			break
		}

		batchTxHash := BatchTxHash(
			m.Chain,
			batch.To,
			batch.SeqNum,
			batch.Value,
			batch.Data,
		)
		messageHash := hashing.SoliditySHA3WithPrefix(batchTxHash[:])
		pubkey, err := crypto.SigToPub(messageHash.Bytes(), batch.Sig[:])
		if err != nil {
			// TODO: Is this possible? If so we need to handle it
			// What are the possible failure conditions and how do they relate
			// to ecrecover's behavior
			log.Fatalln("Invalid sig", err)
		}

		from := common.NewAddressFromEth(crypto.PubkeyToAddress(*pubkey))
		tx := Transaction{
			Chain:       m.Chain,
			To:          batch.To,
			From:        from,
			SequenceNum: batch.SeqNum,
			Value:       batch.Value,
			Data:        batch.Data,
		}

		txes = append(txes, DeliveredTransaction{
			Transaction: tx,
			BlockNum:    m.BlockNum,
			Timestamp:   m.Timestamp,
		})
		offset += batch.encodedLength()
	}
	return txes
}

func (m TransactionBatch) String() string {
	return fmt.Sprintf("TransactionBatch()")
}

// Equals check for equality between this object and any other message by
// checking for full equality of all members
func (m TransactionBatch) Equals(other Message) bool {
	o, ok := other.(TransactionBatch)
	if !ok {
		return false
	}
	return m.Chain != o.Chain || bytes.Equal(m.TxData, o.TxData)
}

func (m TransactionBatch) Type() MessageType {
	return TransactionBatchType
}

type DeliveredTransactionBatch struct {
	TransactionBatch
	BlockNum  *common.TimeBlocks
	Timestamp *big.Int
}

func (m DeliveredTransactionBatch) Equals(other Message) bool {
	o, ok := other.(DeliveredTransactionBatch)
	if !ok {
		return false
	}
	return m.TransactionBatch.Equals(o.TransactionBatch) &&
		m.BlockNum.Cmp(o.BlockNum) == 0 &&
		m.Timestamp.Cmp(o.Timestamp) == 0
}

func (m DeliveredTransactionBatch) deliveredHeight() *common.TimeBlocks {
	return m.BlockNum
}

func (m DeliveredTransactionBatch) deliveredTimestamp() *big.Int {
	return m.Timestamp
}

func (m DeliveredTransactionBatch) CommitmentHash() common.Hash {
	return hashing.SoliditySHA3(
		hashing.Uint8(uint8(m.Type())),
		m.TxData,
		hashing.Uint256(m.BlockNum.AsInt()),
		hashing.Uint256(m.Timestamp),
	)
}

func (m DeliveredTransactionBatch) CheckpointValue() value.Value {
	val, _ := value.NewTupleFromSlice([]value.Value{
		addressToIntValue(m.Chain),
		BytesToByteStack(m.TxData),
		value.NewIntValue(new(big.Int).Set(m.BlockNum.AsInt())),
		value.NewIntValue(new(big.Int).Set(m.Timestamp)),
	})
	return val
}

func UnmarshalTransactionBatchFromCheckpoint(v value.Value) (DeliveredTransactionBatch, error) {
	tup, ok := v.(value.TupleValue)
	failRet := DeliveredTransactionBatch{}
	if !ok || tup.Len() != 4 {
		return failRet, errors.New("tx val must be 7-tuple")
	}
	chain, _ := tup.GetByInt64(0)
	chainInt, ok := chain.(value.IntValue)
	if !ok {
		return failRet, errors.New("chain must be int")
	}
	data, _ := tup.GetByInt64(1)
	dataBytes, err := ByteStackToHex(data)
	if err != nil {
		return failRet, err
	}
	blockNum, _ := tup.GetByInt64(2)
	blockNumInt, ok := blockNum.(value.IntValue)
	if !ok {
		return failRet, errors.New("blockNum must be int")
	}
	timestamp, _ := tup.GetByInt64(3)
	timestampInt, ok := timestamp.(value.IntValue)
	if !ok {
		return failRet, errors.New("timestamp must be int")
	}

	return DeliveredTransactionBatch{
		TransactionBatch: TransactionBatch{
			Chain:  intValueToAddress(chainInt),
			TxData: dataBytes,
		},
		BlockNum:  common.NewTimeBlocks(blockNumInt.BigInt()),
		Timestamp: timestampInt.BigInt(),
	}, nil
}
