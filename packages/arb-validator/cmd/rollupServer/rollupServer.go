/*
 * Copyright 2019, Offchain Labs, Inc.
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

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/offchainlabs/arbitrum/packages/arb-validator/rollupmanager"

	"github.com/offchainlabs/arbitrum/packages/arb-validator/rollupvalidator"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/ethbridge"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/loader"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/rollup"
)

func main() {
	// Check number of args
	flag.Parse()
	switch os.Args[1] {
	case "create":
		createRollupChain()
	case "validate":
		if err := validateRollupChain(); err != nil {
			log.Fatal(err)
		}
	default:
	}
}

func createRollupChain() {
	// Check number of args
	flag.Parse()
	if flag.NArg() != 5 {
		log.Fatalln("usage: rollupServer create <contract.ao> <private_key.txt> <ethURL> <factoryAddress>")
	}

	// 1) Compiled Arbitrum bytecode
	mach, err := loader.LoadMachineFromFile(flag.Arg(1), true, "cpp")
	if err != nil {
		log.Fatal("Loader Error: ", err)
	}

	// 2) Private key
	keyFile, err := os.Open(flag.Arg(2))
	if err != nil {
		log.Fatalln(err)
	}
	byteValue, err := ioutil.ReadAll(keyFile)
	if err != nil {
		log.Fatalln(err)
	}
	if err := keyFile.Close(); err != nil {
		log.Fatalln(err)
	}
	rawKey := strings.TrimPrefix(strings.TrimSpace(string(byteValue)), "0x")
	key, err := crypto.HexToECDSA(rawKey)
	if err != nil {
		log.Fatal("HexToECDSA private key error: ", err)
	}

	// 3) URL
	ethURL := flag.Arg(3)

	// 4) Rollup factory address
	addressString := flag.Arg(4)
	factoryAddress := common.HexToAddress(addressString)

	// Rollup creation
	auth := bind.NewKeyedTransactor(key)
	client, err := ethbridge.NewEthAuthClient(ethURL, auth)
	if err != nil {
		log.Fatal(err)
	}

	factory, err := client.NewArbFactory(factoryAddress)
	if err != nil {
		log.Fatal(err)
	}

	address, err := factory.CreateRollup(
		context.Background(),
		mach.Hash(),
		rollup.DefaultChainParams(),
		common.Address{},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(address.Hex())
}

func validateRollupChain() error {
	// Check number of args

	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	rpcEnable := validateCmd.Bool("rpc", false, "rpc")
	blocktime := validateCmd.Int64("blocktime", 2, "blocktime=NumSeconds")
	gasPrice := validateCmd.Float64("gasprice", 4.5, "gasprice=FloatInGwei")
	err := validateCmd.Parse(os.Args[2:])
	if err != nil {
		return err
	}

	if validateCmd.NArg() != 3 {
		return errors.New("usage: rollupServer validate [--rpc] [--blocktime=NumSeconds] [--gasprice==FloatInGwei] <validator_folder> <ethURL> <rollup_address>")
	}

	common.SetDurationPerBlock(time.Duration(*blocktime) * time.Second)

	validatorFolder := validateCmd.Arg(0)
	ethURL := validateCmd.Arg(1)
	addressString := validateCmd.Arg(2)
	address := common.HexToAddress(addressString)

	keyFilePath := filepath.Join(validatorFolder, "private_key.txt")
	keyFile, err := os.Open(keyFilePath)
	if err != nil {
		return err
	}
	byteValue, err := ioutil.ReadAll(keyFile)
	if err != nil {
		return err
	}
	if err := keyFile.Close(); err != nil {
		return err
	}
	rawKey := strings.TrimPrefix(strings.TrimSpace(string(byteValue)), "0x")
	key, err := crypto.HexToECDSA(rawKey)
	if err != nil {
		return fmt.Errorf("HexToECDSA private key error: %v", err)
	}

	// Rollup creation
	auth := bind.NewKeyedTransactor(key)
	gasPriceAsFloat := 1e9 * (*gasPrice)
	if gasPriceAsFloat < math.MaxInt64 {
		auth.GasPrice = big.NewInt(int64(gasPriceAsFloat))
	}
	client, err := ethbridge.NewEthAuthClient(ethURL, auth)
	if err != nil {
		return err
	}

	rollupActor, err := client.NewRollup(address)
	if err != nil {
		return err
	}

	validatorListener := rollup.NewValidatorChainListener(context.Background(), address, rollupActor)
	err = validatorListener.AddStaker(client)
	if err != nil {
		return err
	}

	contractFile := filepath.Join(validatorFolder, "contract.ao")
	dbPath := filepath.Join(validatorFolder, "checkpoint_db")
	manager, err := rollupmanager.CreateManager(address, client, contractFile, dbPath)

	if err != nil {
		return err
	}
	manager.AddListener(&rollup.AnnouncerListener{})
	manager.AddListener(validatorListener)

	if *rpcEnable {
		if err := rollupvalidator.LaunchRPC(manager, "1235"); err != nil {
			log.Fatal(err)
		}
	} else {
		wait := make(chan bool)
		<-wait
	}
	return nil
}
