module github.com/offchainlabs/arbitrum/tests/fibgo

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.10
	github.com/offchainlabs/arbitrum/packages/arb-provider-go v0.4.2
	github.com/offchainlabs/arbitrum/packages/arb-util v0.4.3
	github.com/offchainlabs/arbitrum/packages/arb-validator v0.4.3
	github.com/offchainlabs/arbitrum/packages/arb-validator-core v0.4.3
)

replace github.com/offchainlabs/arbitrum/packages/arb-provider-go => ../../packages/arb-provider-go

replace github.com/offchainlabs/arbitrum/packages/arb-validator => ../../packages/arb-validator

replace github.com/offchainlabs/arbitrum/packages/arb-validator-core => ../../packages/arb-validator-core

replace github.com/offchainlabs/arbitrum/packages/arb-avm-go => ../../packages/arb-avm-go

replace github.com/offchainlabs/arbitrum/packages/arb-avm-cpp => ../../packages/arb-avm-cpp

replace github.com/offchainlabs/arbitrum/packages/arb-util => ../../packages/arb-util
