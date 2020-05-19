package machine

import (
	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-util/value"
)

type CheckpointStorage interface {
	DeleteCheckpoint(machineHash common.Hash) bool
	CloseCheckpointStorage() bool
	GetInitialMachine() (Machine, error)
	GetMachine(machineHash common.Hash) (Machine, error)
	SaveValue(val value.Value) bool
	GetValue(hashValue common.Hash) value.Value
	DeleteValue(hashValue common.Hash) bool
	SaveData(key []byte, serializedValue []byte) bool
	GetData(key []byte) []byte
	DeleteData(key []byte) bool

	PutBlock(id *common.BlockId, data []byte) error
	DeleteBlock(id *common.BlockId) error
	GetBlock(id *common.BlockId) ([]byte, error)
	BlocksAtHeight(height *common.TimeBlocks) []*common.BlockId
	IsBlockStoreEmpty() bool
	MaxBlockStoreHeight() *common.TimeBlocks
	MinBlockStoreHeight() *common.TimeBlocks
}
