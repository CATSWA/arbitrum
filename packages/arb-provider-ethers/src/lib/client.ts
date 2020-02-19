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
/* eslint-env browser */
'use strict';

import * as ArbValue from './value';

import * as ethers from 'ethers';

// TODO remove this dep
const jaysonBrowserClient = require('jayson/lib/client/browser'); // eslint-disable-line @typescript-eslint/no-var-requires

export enum EVMCode {
    Revert = 0,
    Invalid = 1,
    Return = 2,
    Stop = 3,
    BadSequenceCode = 4,
}

function logValToLog(val: ArbValue.Value, index: number, orig: EthBridgeMessage): ethers.providers.Log {
    const value = val as ArbValue.TupleValue;
    return {
        blockNumber: orig.blockNumber.toNumber(),
        blockHash: orig.txHash,
        transactionIndex: 0,
        removed: false,
        transactionLogIndex: index,
        address: ethers.utils.hexlify((value.get(0) as ArbValue.IntValue).bignum),
        data: ethers.utils.hexlify(ArbValue.bytestackToBytes(value.get(1) as ArbValue.TupleValue)),
        topics: value.contents
            .slice(2)
            .map(rawTopic => ethers.utils.hexZeroPad(ethers.utils.hexlify((rawTopic as ArbValue.IntValue).bignum), 32)),
        transactionHash: orig.txHash,
        logIndex: index,
    };
}

function stackValueToList(value: ArbValue.TupleValue): ArbValue.Value[] {
    const values = [];
    while (value.contents.length !== 0) {
        values.push(value.get(1));
        value = value.get(0) as ArbValue.TupleValue;
    }
    return values;
}

class EthBridgeMessage {
    public typecode: number;
    public blockNumber: ethers.utils.BigNumber;
    public txHash: string;
    public sender: string;
    public message: ArbValue.TupleValue;
    public calldataHash: string;

    constructor(value: ArbValue.TupleValue) {
        this.blockNumber = (value.get(0) as ArbValue.IntValue).bignum;
        this.txHash = ethers.utils.hexZeroPad((value.get(1) as ArbValue.IntValue).bignum.toHexString(), 32);
        const restVal = value.get(2) as ArbValue.TupleValue;
        this.typecode = (restVal.get(0) as ArbValue.IntValue).bignum.toNumber();
        this.sender = ethers.utils.getAddress((restVal.get(1) as ArbValue.IntValue).bignum.toHexString());
        this.message = restVal.get(2) as ArbValue.TupleValue;
        this.calldataHash = restVal.hash();
    }

    getArbMessage(): ArbMessage {
        switch (this.typecode) {
            case 0:
                return new TxMessage(this.message);
            case 1:
                return new EthTransferMessage(this.message);
            case 2:
                return new TokenTransferMessage(this.message);
            case 3:
                return new TokenTransferMessage(this.message);
            case 4:
                return new ContractTxMessage(this.message);
            case 5:
                return new TxCall(this.message);
            default:
                throw 'Invalid arb message type';
        }
    }
}

export class TxCall {
    public to: string;
    public data: Uint8Array;

    constructor(value: ArbValue.TupleValue) {
        this.to = ethers.utils.getAddress(
            ethers.utils.hexZeroPad((value.get(0) as ArbValue.IntValue).bignum.toHexString(), 20),
        );
        this.data = ArbValue.bytestackToBytes(value.get(1) as ArbValue.TupleValue);
    }

    getDest(): string {
        return this.to;
    }
}

export class TxMessage {
    public to: string;
    public sequenceNum: ethers.utils.BigNumber;
    public amount: ethers.utils.BigNumber;
    public data: Uint8Array;

    constructor(value: ArbValue.TupleValue) {
        this.to = ethers.utils.getAddress(
            ethers.utils.hexZeroPad((value.get(0) as ArbValue.IntValue).bignum.toHexString(), 20),
        );
        this.sequenceNum = (value.get(1) as ArbValue.IntValue).bignum;
        this.amount = (value.get(2) as ArbValue.IntValue).bignum;
        this.data = ArbValue.bytestackToBytes(value.get(3) as ArbValue.TupleValue);
    }

    getDest(): string {
        return this.to;
    }
}

export class ContractTxMessage {
    public to: string;
    public amount: ethers.utils.BigNumber;
    public data: Uint8Array;

    constructor(value: ArbValue.TupleValue) {
        this.to = ethers.utils.getAddress(
            ethers.utils.hexZeroPad((value.get(0) as ArbValue.IntValue).bignum.toHexString(), 20),
        );
        this.amount = (value.get(1) as ArbValue.IntValue).bignum;
        this.data = ArbValue.bytestackToBytes(value.get(2) as ArbValue.TupleValue);
    }

    getDest(): string {
        return this.to;
    }
}

class EthTransferMessage {
    public dest: string;
    public amount: ethers.utils.BigNumber;

    constructor(value: ArbValue.TupleValue) {
        this.dest = ethers.utils.getAddress((value.get(0) as ArbValue.IntValue).bignum.toHexString());
        this.amount = (value.get(1) as ArbValue.IntValue).bignum;
    }

    getDest(): string {
        return this.dest;
    }
}
class TokenTransferMessage {
    public tokenAddress: string;
    public dest: string;
    public amount: ethers.utils.BigNumber;

    constructor(value: ArbValue.TupleValue) {
        this.tokenAddress = ethers.utils.getAddress((value.get(0) as ArbValue.IntValue).bignum.toHexString());
        this.dest = ethers.utils.getAddress((value.get(1) as ArbValue.IntValue).bignum.toHexString());
        this.amount = (value.get(2) as ArbValue.IntValue).bignum;
    }

    getDest(): string {
        return this.dest;
    }
}

export type ArbMessage = TxMessage | EthTransferMessage | TokenTransferMessage | TxCall;

export type EVMResult = EVMReturn | EVMRevert | EVMStop | EVMBadSequenceCode | EVMInvalid;

export class EVMReturn {
    public bridgeData: EthBridgeMessage;
    public orig: ArbMessage;
    public data: Uint8Array;
    public logs: ethers.providers.Log[];
    public returnType: EVMCode.Return;

    constructor(value: ArbValue.TupleValue) {
        this.bridgeData = new EthBridgeMessage(value.get(0) as ArbValue.TupleValue);
        this.orig = this.bridgeData.getArbMessage();
        this.data = ArbValue.bytestackToBytes(value.get(2) as ArbValue.TupleValue);
        this.logs = stackValueToList(value.get(1) as ArbValue.TupleValue).map((val, index) => {
            return logValToLog(val, index, this.bridgeData);
        });
        this.returnType = EVMCode.Return;
    }
}

export class EVMRevert {
    public bridgeData: EthBridgeMessage;
    public orig: ArbMessage;
    public data: Uint8Array;
    public returnType: EVMCode.Revert;

    constructor(value: ArbValue.TupleValue) {
        this.bridgeData = new EthBridgeMessage(value.get(0) as ArbValue.TupleValue);
        this.orig = this.bridgeData.getArbMessage();
        this.data = ArbValue.bytestackToBytes(value.get(2) as ArbValue.TupleValue);
        this.returnType = EVMCode.Revert;
    }
}

export class EVMStop {
    public bridgeData: EthBridgeMessage;
    public orig: ArbMessage;
    public logs: ethers.providers.Log[];
    public returnType: EVMCode.Stop;

    constructor(value: ArbValue.TupleValue) {
        this.bridgeData = new EthBridgeMessage(value.get(0) as ArbValue.TupleValue);
        this.orig = this.bridgeData.getArbMessage();
        this.logs = stackValueToList(value.get(1) as ArbValue.TupleValue).map((val, index) => {
            return logValToLog(val, index, this.bridgeData);
        });
        this.returnType = EVMCode.Stop;
    }
}

export class EVMBadSequenceCode {
    public bridgeData: EthBridgeMessage;
    public orig: ArbMessage;
    public returnType: EVMCode.BadSequenceCode;

    constructor(value: ArbValue.TupleValue) {
        this.bridgeData = new EthBridgeMessage(value.get(0) as ArbValue.TupleValue);
        this.orig = this.bridgeData.getArbMessage();
        this.returnType = EVMCode.BadSequenceCode;
    }
}

export class EVMInvalid {
    public bridgeData: EthBridgeMessage;
    public orig: ArbMessage;
    public returnType: EVMCode.Invalid;

    constructor(value: ArbValue.TupleValue) {
        this.bridgeData = new EthBridgeMessage(value.get(0) as ArbValue.TupleValue);
        this.orig = this.bridgeData.getArbMessage();
        this.returnType = EVMCode.Invalid;
    }
}

function processLog(value: ArbValue.TupleValue): EVMResult {
    const returnCode = value.get(3) as ArbValue.IntValue;
    switch (returnCode.bignum.toNumber()) {
        case EVMCode.Return:
            return new EVMReturn(value);
        case EVMCode.Revert:
            return new EVMRevert(value);
        case EVMCode.Stop:
            return new EVMStop(value);
        case EVMCode.BadSequenceCode:
            return new EVMBadSequenceCode(value);
        case EVMCode.Invalid:
            return new EVMInvalid(value);
        default:
            throw Error('processLogs Invalid EVM return code');
    }
}

interface GetVMInfoReply {
    vmID: string;
}

interface GetValidatorListReply {
    validators: string[];
}

interface GetAssertionCountReply {
    assertionCount: number;
}

interface GetMessageResultReply {
    found: boolean;
    rawVal: string;
    logPreHash: string;
    logPostHash: string;
    logValHashes: string[];
    validatorSigs: string[];
    partialHash: string;
    onChainTxHash: string;
}

interface SendMessageReply {
    txHash: string;
}

interface CallMessageReply {
    rawVal: string;
}

interface LogInfo {
    address: string;
    blockHash: string;
    blockNumber: string;
    data: string;
    logIndex: string;
    topics: string[];
    transactionIndex: string;
    transactionHash: string;
}

interface FindLogsReply {
    logs: LogInfo[];
}

function _arbClient(managerAddress: string): any {
    const callServer = (request: any, callback: any): void => {
        const options = {
            body: request, // request is a string
            headers: {
                'Content-Type': 'application/json',
            },
            method: 'POST',
        };

        fetch(managerAddress, options)
            .then((res: any) => {
                return res.text();
            })
            .then((text: string) => {
                callback(null, text);
            })
            .catch((err: Error) => {
                callback(err);
            });
    };

    return jaysonBrowserClient(callServer, {});
}

interface MessageResult {
    logPostHash: string;
    logPreHash: string;
    logValHashes: string[];
    onChainTxHash: string;
    partialHash: string;
    val: ArbValue.Value;
    validatorSigs: string[];
    vmId: string;
    evmVal: EVMResult;
}

export class ArbClient {
    public client: any;

    constructor(managerUrl: string) {
        this.client = _arbClient(managerUrl);
    }

    public async getMessageResult(txHash: string): Promise<MessageResult | null> {
        const messageResult = await new Promise<GetMessageResultReply>((resolve, reject): void => {
            this.client.request(
                'Validator.GetMessageResult',
                [
                    {
                        txHash,
                    },
                ],
                (err: Error, error: Error, result: GetMessageResultReply) => {
                    if (err) {
                        reject(err);
                    } else if (error) {
                        reject(error);
                    } else {
                        resolve(result);
                    }
                },
            );
        });
        if (messageResult.found) {
            const vmId = await this.getVmID();
            const val = ArbValue.unmarshal(messageResult.rawVal);
            const evmVal = processLog(val as ArbValue.TupleValue);
            let logValHashes = messageResult.logValHashes;
            if (!logValHashes) {
                logValHashes = [];
            }

            return {
                logPostHash: messageResult.logPostHash,
                logPreHash: messageResult.logPreHash,
                logValHashes,
                onChainTxHash: messageResult.onChainTxHash,
                partialHash: messageResult.partialHash,
                val,
                validatorSigs: messageResult.validatorSigs,
                vmId,
                evmVal,
            };
        } else {
            return null;
        }
    }

    public sendMessage(
        to: string,
        sequenceNum: ethers.utils.BigNumberish,
        value: ethers.utils.BigNumberish,
        data: string,
        signature: string,
        pubkey: string,
    ): Promise<string> {
        return new Promise((resolve, reject): void => {
            this.client.request(
                'Validator.SendMessage',
                [
                    {
                        to,
                        sequenceNum,
                        value,
                        data,
                        signature,
                        pubkey,
                    },
                ],
                (err: Error, error: Error, result: SendMessageReply) => {
                    if (err) {
                        reject(err);
                    } else if (error) {
                        reject(error);
                    } else {
                        resolve(result.txHash);
                    }
                },
            );
        });
    }

    public call(contractAddress: string, sender: string, data: string): Promise<Uint8Array> {
        return new Promise((resolve, reject): void => {
            this.client.request(
                'Validator.CallMessage',
                [
                    {
                        contractAddress,
                        data,
                        sender,
                    },
                ],
                (err: Error, error: Error, result: CallMessageReply) => {
                    if (err) {
                        reject(err);
                    } else if (error) {
                        reject(error);
                    } else {
                        const val = ArbValue.unmarshal(result.rawVal);
                        const evmVal = processLog(val as ArbValue.TupleValue);
                        switch (evmVal.returnType) {
                            case EVMCode.Return:
                                resolve(evmVal.data);
                                break;
                            case EVMCode.Stop:
                                resolve(new Uint8Array());
                                break;
                            default:
                                reject(new Error('Call was reverted'));
                                break;
                        }
                    }
                },
            );
        });
    }

    public findLogs(fromBlock: number, toBlock: number, address: string, topics: string[]): Promise<LogInfo[]> {
        return new Promise((resolve, reject): void => {
            return this.client.request(
                'Validator.FindLogs',
                [
                    {
                        address,
                        fromHeight: fromBlock,
                        toHeight: toBlock,
                        topics,
                    },
                ],
                (err: Error, error: Error, result: FindLogsReply) => {
                    if (err) {
                        reject(err);
                    } else if (error) {
                        reject(error);
                    } else {
                        resolve(result.logs);
                    }
                },
            );
        });
    }

    public getVmID(): Promise<string> {
        return new Promise((resolve, reject): void => {
            this.client.request('Validator.GetVMInfo', [], (err: Error, error: Error, result: GetVMInfoReply) => {
                if (err) {
                    reject(err);
                } else if (error) {
                    reject(error);
                } else {
                    resolve(result.vmID);
                }
            });
        });
    }

    public getAssertionCount(): Promise<number> {
        return new Promise((resolve, reject): void => {
            this.client.request(
                'Validator.GetAssertionCount',
                [],
                (err: Error, error: Error, result: GetAssertionCountReply) => {
                    if (err) {
                        reject(err);
                    } else if (error) {
                        reject(error);
                    } else {
                        resolve(result.assertionCount);
                    }
                },
            );
        });
    }

    public getValidatorList(): Promise<string[]> {
        return new Promise((resolve, reject): void => {
            this.client.request(
                'Validator.GetValidatorList',
                [],
                (err: Error, error: Error, result: GetValidatorListReply) => {
                    if (err) {
                        reject(err);
                    } else if (error) {
                        reject(error);
                    } else {
                        resolve(result.validators);
                    }
                },
            );
        });
    }
}
