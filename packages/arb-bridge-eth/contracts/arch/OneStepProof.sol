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

pragma solidity ^0.5.3;

import "./Value.sol";
import "./Machine.sol";
import "../libraries/DebugPrint.sol";

// Sourced from https://github.com/leapdao/solEVM-enforcer/tree/master


library OneStepProof {
    using Machine for Machine.Data;
    using Value for Value.Data;

    struct ValidateProofData {
        bytes32 beforeHash;
        uint128[2] timeBoundsBlocks;
        bytes32 beforeInbox;
        bytes32 afterHash;
        bool    didInboxInsn;
        bytes32 firstMessage;
        bytes32 lastMessage;
        bytes32 firstLog;
        bytes32 lastLog;
        uint64  gas;
        bytes proof;
    }

    function validateProof(
        bytes32 beforeHash,
        uint128[2] memory timeBoundsBlocks,
        bytes32 beforeInbox,
        bytes32 afterHash,
        bool    didInboxInsn,
        bytes32 firstMessage,
        bytes32 lastMessage,
        bytes32 firstLog,
        bytes32 lastLog,
        uint64  gas,
        bytes memory proof
    )
        public
        pure
        returns(uint)
    {
        return checkProof(
            ValidateProofData(
                beforeHash,
                timeBoundsBlocks,
                beforeInbox,
                afterHash,
                didInboxInsn,
                firstMessage,
                lastMessage,
                firstLog,
                lastLog,
                gas,
                proof
            )
        );
    }

    // Arithmetic

    function executeAddInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := add(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeMulInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := mul(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeSubInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := sub(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeDivInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        if (b == 0) {
            return false;
        }
        uint c;
        assembly {
            c := div(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeSdivInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        if (b == 0) {
            return false;
        }
        uint c;
        assembly {
            c := sdiv(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeModInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        if (b == 0) {
            return false;
        }
        uint c;
        assembly {
            c := mod(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeSmodInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        if (b == 0) {
            return false;
        }
        uint c;
        assembly {
            c := smod(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeAddmodInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2,
        Value.Data memory val3
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint m = val3.intVal;
        if (m == 0) {
            return false;
        }
        uint c;
        assembly {
            c := addmod(a, b, m)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeMulmodInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2,
        Value.Data memory val3
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint m = val3.intVal;
        if (m == 0) {
            return false;
        }
        uint c;
        assembly {
            c := mulmod(a, b, m)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeExpInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := exp(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    // Comparison

    function executeLtInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := lt(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeGtInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := gt(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeSltInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := slt(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeSgtInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := sgt(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeEqInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackValue(Value.newBoolean(val1.hash().hash == val2.hash().hash));
        return true;
    }

    function executeIszeroInsn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt()) {
            machine.addDataStackInt(0);
        } else {
            uint a = val1.intVal;
            uint c;
            assembly {
                c := iszero(a)
            }
            machine.addDataStackInt(c);
        }
        return true;
    }

    function executeAndInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := and(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeOrInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := or(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeXorInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        uint c;
        assembly {
            c := xor(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeNotInsn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint c;
        assembly {
            c := not(a)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeByteInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint x = val1.intVal;
        uint n = val2.intVal;
        uint c;
        assembly {
            c := byte(n, x)
        }
        machine.addDataStackInt(c);
        return true;
    }

    function executeSignextendInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint b = val1.intVal;
        uint a = val2.intVal;
        uint c;
        assembly {
            c := signextend(a, b)
        }
        machine.addDataStackInt(c);
        return true;
    }

    // Hash

    function executeSha3Insn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackInt(uint256(val1.hash().hash));
        return true;
    }

    function executeTypeInsn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackValue(val1.typeCodeVal());
        return true;
    }

    function executeEthhash2Insn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    ) 
        internal 
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isInt()) {
            return false;
        }
        uint a = val1.intVal;
        uint b = val2.intVal;
        bytes32 res = keccak256(abi.encodePacked(a, b));
        machine.addDataStackInt(uint256(res));
        return true;
    }


    // Stack ops

    function executePopInsn(
        Machine.Data memory,
        Value.Data memory
    )
        internal
        pure
        returns (bool)
    {
        return true;
    }

    function executeSpushInsn(
        Machine.Data memory machine
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackHashValue(machine.staticHash);
        return true;
    }

    function executeRpushInsn(
        Machine.Data memory machine
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackHashValue(machine.registerHash);
        return true;
    }

    function executeRsetInsn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        machine.registerHash = val1.hash();
        return true;
    }

    function executeJumpInsn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        machine.instructionStackHash = val1.hash();
        return true;
    }

    function executeCjumpInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isCodePoint()) {
            return false;
        }
        if (!val2.isInt()) {
            return false;
        }
        if (val2.intVal != 0) {
            machine.instructionStackHash = val1.hash();
        }
        return true;
    }

    function executeStackemptyInsn(
        Machine.Data memory machine
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackValue(
            Value.newBoolean(machine.dataStackHash.hash == Value.newNone().hash().hash)
        );
        return true;
    }

    function executePcpushInsn(
        Machine.Data memory machine,
        Value.HashOnly memory pc
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackHashValue(pc);
        return true;
    }

    function executeAuxpushInsn(
        Machine.Data memory machine,
        Value.Data memory val
    )
        internal
        pure
        returns (bool)
    {
        machine.addAuxStackValue(val);
        return true;
    }

    function executeAuxstackemptyInsn(
        Machine.Data memory machine
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackValue(
            Value.newBoolean(machine.auxStackHash.hash == Value.newNone().hash().hash)
        );
        return true;
    }

    function executeErrpushInsn(
        Machine.Data memory machine
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackHashValue(machine.errHandler);
        return true;
    }

    function executeErrsetInsn(
        Machine.Data memory machine,
        Value.Data memory val
    )
        internal
        pure
        returns (bool)
    {
        if (!val.isCodePoint()) {
            return false;
        }
        machine.errHandler = val.hash();
        return true;
    }

    // Dup ops

    function executeDup0Insn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackValue(val1);
        machine.addDataStackValue(val1);
        return true;
    }

    function executeDup1Insn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackValue(val2);
        machine.addDataStackValue(val1);
        machine.addDataStackValue(val2);
        return true;
    }

    function executeDup2Insn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2,
        Value.Data memory val3
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackValue(val3);
        machine.addDataStackValue(val2);
        machine.addDataStackValue(val1);
        machine.addDataStackValue(val3);
        return true;
    }

    // Swap ops

    function executeSwap1Insn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackValue(val1);
        machine.addDataStackValue(val2);
        return true;
    }

    function executeSwap2Insn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2,
        Value.Data memory val3
    )
        internal
        pure
        returns (bool)
    {
        machine.addDataStackValue(val1);
        machine.addDataStackValue(val2);
        machine.addDataStackValue(val3);
        return true;
    }

    // Tuple ops

    function executeTgetInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt() || !val2.isTuple()) {
            return false;
        }

        if (val1.intVal >= val2.valLength()) {
            return false;
        }

        machine.addDataStackValue(val2.tupleVal[val1.intVal]);
        return true;
    }

    function executeTsetInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.Data memory val2,
        Value.Data memory val3
    )
        internal
        pure
        returns (bool)
    {
        if (!val2.isTuple() || !val1.isInt()) {
            return false;
        }

        if (val1.intVal >= val2.valLength()) {
            return false;
        }
        val2.tupleVal[val1.intVal] = val3;
        machine.addDataStackValue(val2);
        return true;
    }

    function executeTlenInsn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isTuple()) {
            return false;
        }
        machine.addDataStackInt(val1.valLength());
        return true;
    }

    // Memory insns
    function executeMnewInsn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        if (!val1.isInt()) {
            return false;
        } 
        uint iv256 = val1.intVal;
        require(iv256 < (1<<32));
        uint64 iv64 = uint64(iv256);       
        machine.addDataStackHashValue(Value.HashOnly(hashOfEmptyMem(iv64)));
        return true;
    }

    function executeMgetInsn(
        Machine.Data memory machine, 
        Value.Data memory val1, 
        Value.Data memory val2, 
        bytes memory proof, 
        uint offset
    ) 
        internal
        pure
        returns(bool, uint)  // success, offset
    {
        if (!val1.isInt()) {
            return (false, offset);
        }
        uint index = val1.intVal;
        bytes32 memHash = Value.hash(val2).hash;
        bool success;
        uint size;
        (success, offset, size) = Value.deserializeInt(proof, offset);
        if (!success){
            return (false, offset);
        }
        Value.HashOnly memory returnHash;
        (success, offset, returnHash) = Value.deserializeHashOnly(proof, offset);
        bytes32 resultHash = returnHash.hash;
        if (!success){
            return (false, offset);
        }
        for (uint chunkSize=1; chunkSize<size; chunkSize=chunkSize*2) {
            Value.HashOnly memory h;
            (success, offset, h) = Value.deserializeHashOnly(proof, offset);
            if ((index & chunkSize) == 0) {
                resultHash = keccak256(abi.encodePacked(
                    resultHash, 
                    h.hash
                ));
            } else {
                resultHash = keccak256(abi.encodePacked(
                    h.hash, 
                    resultHash
                ));
            }
        }
        machine.addDataStackHashValue(returnHash);
        return (resultHash==memHash, offset);
    }

    function hashOfEmptyMem(uint64 size) 
        internal
        pure
        returns(bytes32)
    {
        uint64 chunkSize = 1;
        bytes32 h = Value.hashEmptyTuple();
        h = keccak256(abi.encodePacked(h, h));
        while (2*chunkSize < size) {
            h = keccak256(abi.encodePacked(h, h));
            chunkSize = 2*chunkSize;
        }

        return keccak256(
            abi.encodePacked(
                uint8(12),
                size,
                h
            )
        );
    }

    // Logging

    function executeBreakpointInsn(Machine.Data memory) internal pure returns (bool) {
        return true;
    }

    function executeLogInsn(
        Machine.Data memory,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool, bytes32)
    {
        Value.HashOnly memory hashVal = val1.hash();
        return (true, hashVal.hash);
    }

    // System operations

    function executeSendInsn(
        Machine.Data memory,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool, bytes32)
    {
        return (true, val1.hash().hash);
    }

    function executeInboxInsn(
        Machine.Data memory machine,
        Value.Data memory val1,
        Value.HashOnly memory beforeInbox,
        uint lowerTimeBound
    )
        internal
        pure
        returns (bool)
    {
        if (! val1.isInt()) {
            return false;
        }
        require(lowerTimeBound<val1.intVal && beforeInbox.hash==Value.hashEmptyTuple(),
            "Inbox instruction was blocked");
        machine.addDataStackHashValue(beforeInbox);
        return true;
    }

    function executeSetgasInsn(
        Machine.Data memory machine,
        Value.Data memory val1
    )
        internal
        pure
        returns (bool)
    {
        if (! val1.isInt()) {
            return false;
        }
        machine.arbGasRemaining = val1.intVal;
        return true;
    }

    function executePushgasInsn(
        Machine.Data memory machine
    )
        internal
        pure
        returns (bool) 
    {
        machine.addDataStackInt(machine.arbGasRemaining);
        return true;   
    }

    // Stop and arithmetic ops
    uint8 constant internal OP_ADD = 0x01;
    uint8 constant internal OP_MUL = 0x02;
    uint8 constant internal OP_SUB = 0x03;
    uint8 constant internal OP_DIV = 0x04;
    uint8 constant internal OP_SDIV = 0x05;
    uint8 constant internal OP_MOD = 0x06;
    uint8 constant internal OP_SMOD = 0x07;
    uint8 constant internal OP_ADDMOD = 0x08;
    uint8 constant internal OP_MULMOD = 0x09;
    uint8 constant internal OP_EXP = 0x0a;

    // Comparison & bitwise logic
    uint8 constant internal OP_LT = 0x10;
    uint8 constant internal OP_GT = 0x11;
    uint8 constant internal OP_SLT = 0x12;
    uint8 constant internal OP_SGT = 0x13;
    uint8 constant internal OP_EQ = 0x14;
    uint8 constant internal OP_ISZERO = 0x15;
    uint8 constant internal OP_AND = 0x16;
    uint8 constant internal OP_OR = 0x17;
    uint8 constant internal OP_XOR = 0x18;
    uint8 constant internal OP_NOT = 0x19;
    uint8 constant internal OP_BYTE = 0x1a;
    uint8 constant internal OP_SIGNEXTEND = 0x1b;

    // SHA3
    uint8 constant internal OP_SHA3 = 0x20;
    uint8 constant internal OP_TYPE = 0x21;
    uint8 constant internal OP_ETHHASH2 = 0x22;

    // Stack, Memory, Storage and Flow Operations
    uint8 constant internal OP_POP = 0x30;
    uint8 constant internal OP_SPUSH = 0x31;
    uint8 constant internal OP_RPUSH = 0x32;
    uint8 constant internal OP_RSET = 0x33;
    uint8 constant internal OP_JUMP = 0x34;
    uint8 constant internal OP_CJUMP = 0x35;
    uint8 constant internal OP_STACKEMPTY = 0x36;
    uint8 constant internal OP_PCPUSH = 0x37;
    uint8 constant internal OP_AUXPUSH = 0x38;
    uint8 constant internal OP_AUXPOP = 0x39;
    uint8 constant internal OP_AUXSTACKEMPTY = 0x3a;
    uint8 constant internal OP_NOP = 0x3b;
    uint8 constant internal OP_ERRPUSH = 0x3c;
    uint8 constant internal OP_ERRSET = 0x3d;

    // Duplication and Exchange operations
    uint8 constant internal OP_DUP0 = 0x40;
    uint8 constant internal OP_DUP1 = 0x41;
    uint8 constant internal OP_DUP2 = 0x42;
    uint8 constant internal OP_SWAP1 = 0x43;
    uint8 constant internal OP_SWAP2 = 0x44;

    // Tuple opertations
    uint8 constant internal OP_TGET = 0x50;
    uint8 constant internal OP_TSET = 0x51;
    uint8 constant internal OP_TLEN = 0x52;
    uint8 constant internal OP_MNEW = 0x53;
    uint8 constant internal OP_MGET = 0x54;

    // Logging opertations
    uint8 constant internal OP_BREAKPOINT = 0x60;
    uint8 constant internal OP_LOG = 0x61;

    // System operations
    uint8 constant internal OP_SEND = 0x70;
    uint8 constant internal OP_GETTIME = 0x71;
    uint8 constant internal OP_INBOX = 0x72;
    uint8 constant internal OP_ERROR = 0x73;
    uint8 constant internal OP_STOP = 0x74;
    uint8 constant internal OP_SETGAS = 0x75;
    uint8 constant internal OP_PUSHGAS = 0x76;

    function opInfo(uint opCode) internal pure returns (uint, uint) {
        if (opCode == OP_ADD) {
            return (2, 1);
        } else if (opCode == OP_MUL) {
            return (2, 1);
        } else if (opCode == OP_SUB) {
            return (2, 1);
        } else if (opCode == OP_DIV) {
            return (2, 1);
        } else if (opCode == OP_SDIV) {
            return (2, 1);
        } else if (opCode == OP_MOD) {
            return (2, 1);
        } else if (opCode == OP_SMOD) {
            return (2, 1);
        } else if (opCode == OP_ADDMOD) {
            return (3, 1);
        } else if (opCode == OP_MULMOD) {
            return (3, 1);
        } else if (opCode == OP_EXP) {
            return (2, 1);
        } else if (opCode == OP_LT) {
            return (2, 1);
        } else if (opCode == OP_GT) {
            return (2, 1);
        } else if (opCode == OP_SLT) {
            return (2, 1);
        } else if (opCode == OP_SGT) {
            return (2, 1);
        } else if (opCode == OP_EQ) {
            return (2, 1);
        } else if (opCode == OP_ISZERO) {
            return (1, 1);
        } else if (opCode == OP_AND) {
            return (2, 1);
        } else if (opCode == OP_OR) {
            return (2, 1);
        } else if (opCode == OP_XOR) {
            return (2, 1);
        } else if (opCode == OP_NOT) {
            return (1, 1);
        } else if (opCode == OP_BYTE) {
            return (2, 1);
        } else if (opCode == OP_SIGNEXTEND) {
            return (2, 1);
        } else if (opCode == OP_SHA3) {
            return (1, 1);
        } else if (opCode == OP_TYPE) {
            return (1, 1);
        } else if (opCode == OP_ETHHASH2) {
            return (2, 1);
        } else if (opCode == OP_POP) {
            return (1, 0);
        } else if (opCode == OP_SPUSH) {
            return (0, 1);
        } else if (opCode == OP_RPUSH) {
            return (0, 1);
        } else if (opCode == OP_RSET) {
            return (1, 0);
        } else if (opCode == OP_JUMP) {
            return (1, 0);
        } else if (opCode == OP_CJUMP) {
            return (2, 0);
        } else if (opCode == OP_STACKEMPTY) {
            return (0, 1);
        } else if (opCode == OP_PCPUSH) {
            return (0, 1);
        } else if (opCode == OP_AUXPUSH) {
            return (1, 0);
        } else if (opCode == OP_AUXPOP) {
            return (0, 1);
        } else if (opCode == OP_AUXSTACKEMPTY) {
            return (0, 1);
        } else if (opCode == OP_NOP) {
            return (0, 0);
        } else if (opCode == OP_ERRPUSH) {
            return (0, 1);
        } else if (opCode == OP_ERRSET) {
            return (1, 0);
        } else if (opCode == OP_DUP0) {
            return (1, 2);
        } else if (opCode == OP_DUP1) {
            return (2, 3);
        } else if (opCode == OP_DUP2) {
            return (3, 4);
        } else if (opCode == OP_SWAP1) {
            return (2, 2);
        } else if (opCode == OP_SWAP2) {
            return (3, 3);
        } else if (opCode == OP_TGET) {
            return (2, 1);
        } else if (opCode == OP_TSET) {
            return (3, 1);
        } else if (opCode == OP_TLEN) {
            return (1, 1);
        } else if (opCode == OP_MNEW) {
            return (1, 1);
        } else if (opCode == OP_MGET) {
            return (2, 1);
        } else if (opCode == OP_BREAKPOINT) {
            return (0, 0);
        } else if (opCode == OP_LOG) {
            return (1, 0);
        } else if (opCode == OP_SEND) {
            return (1, 0);
        } else if (opCode == OP_GETTIME) {
            return (0, 1);
        } else if (opCode == OP_INBOX) {
            return (1, 1);
        } else if (opCode == OP_ERROR) {
            return (0, 0);
        } else if (opCode == OP_STOP) {
            return (0, 0);
        } else if (opCode == OP_SETGAS) {
            return (1, 0);
        } else if (opCode == OP_PUSHGAS) {
            return (0, 1);
        } else {
            require(false, "Invalid opcode");
        }
    }

    function opPopCount(uint8 opCode) internal pure returns(uint) {
        uint popCount;
        uint pushCount;
        (popCount, pushCount) = opInfo(opCode);
        return popCount;
    }

    function opGasCost(uint8 opCode) internal pure returns(uint64) {
       if (opCode == OP_ADD) {
            return 3;
        } else if (opCode == OP_MUL) {
            return 3;
        } else if (opCode == OP_SUB) {
            return 3;
        } else if (opCode == OP_DIV) {
            return 4;
        } else if (opCode == OP_SDIV) {
            return 7;
        } else if (opCode == OP_MOD) {
            return 4;
        } else if (opCode == OP_SMOD) {
            return 7;
        } else if (opCode == OP_ADDMOD) {
            return 4;
        } else if (opCode == OP_MULMOD) {
            return 4;
        } else if (opCode == OP_EXP) {
            return 25;
        } else if (opCode == OP_LT) {
            return 2;
        } else if (opCode == OP_GT) {
            return 2;
        } else if (opCode == OP_SLT) {
            return 2;
        } else if (opCode == OP_SGT) {
            return 2;
        } else if (opCode == OP_EQ) {
            return 2;
        } else if (opCode == OP_ISZERO) {
            return 1;
        } else if (opCode == OP_AND) {
            return 2;
        } else if (opCode == OP_OR) {
            return 2;
        } else if (opCode == OP_XOR) {
            return 2;
        } else if (opCode == OP_NOT) {
            return 1;
        } else if (opCode == OP_BYTE) {
            return 4;
        } else if (opCode == OP_SIGNEXTEND) {
            return 7;
        } else if (opCode == OP_SHA3) {
            return 7;
        } else if (opCode == OP_TYPE) {
            return 3;
       } else if (opCode == OP_ETHHASH2) {
            return 8;
        } else if (opCode == OP_POP) {
            return 1;
        } else if (opCode == OP_SPUSH) {
            return 1;
        } else if (opCode == OP_RPUSH) {
            return 1;
        } else if (opCode == OP_RSET) {
            return 2;
        } else if (opCode == OP_JUMP) {
            return 4;
        } else if (opCode == OP_CJUMP) {
            return 4;
        } else if (opCode == OP_STACKEMPTY) {
            return 2;
        } else if (opCode == OP_PCPUSH) {
            return 1;
        } else if (opCode == OP_AUXPUSH) {
            return 1;
        } else if (opCode == OP_AUXPOP) {
            return 1;
        } else if (opCode == OP_AUXSTACKEMPTY) {
            return 2;
        } else if (opCode == OP_NOP) {
            return 1;
        } else if (opCode == OP_ERRPUSH) {
            return 1;
        } else if (opCode == OP_ERRSET) {
            return 1;
        } else if (opCode == OP_DUP0) {
            return 1;
        } else if (opCode == OP_DUP1) {
            return 1;
        } else if (opCode == OP_DUP2) {
            return 1;
        } else if (opCode == OP_SWAP1) {
            return 1;
        } else if (opCode == OP_SWAP2) {
            return 1;
        } else if (opCode == OP_TGET) {
            return 2;
        } else if (opCode == OP_TSET) {
            return 40;
        } else if (opCode == OP_TLEN) {
            return 2;
        } else if (opCode == OP_MNEW) {
            return 40;
        } else if (opCode == OP_MGET) {
            return 100;
        } else if (opCode == OP_BREAKPOINT) {
            return 100;
        } else if (opCode == OP_LOG) {
            return 100;
        } else if (opCode == OP_SEND) {
            return 100;
        } else if (opCode == OP_GETTIME) {
            return 40;
        } else if (opCode == OP_INBOX) {
            return 40;
        } else if (opCode == OP_ERROR) {
            return 5;
        } else if (opCode == OP_STOP) {
            return 10;
        } else if (opCode == OP_SETGAS) {
            return 0;
        } else if (opCode == OP_PUSHGAS) {
            return 1;
        } else {
            require(false, "Invalid opcode");
        }
    }

    function loadMachine(
        ValidateProofData memory _data
    )
        internal
        pure
        returns (
            uint8,
            Value.Data[] memory,
            Machine.Data memory,
            Machine.Data memory,
            uint
        )
    {
        uint offset = 0;
        bool valid;
        Machine.Data memory startMachine;
        startMachine.setExtensive();
        (valid, offset, startMachine) = Machine.deserializeMachine(_data.proof, offset);
        require(valid, "Not a valid serialized machine");
        Machine.Data memory endMachine = startMachine.clone();
        uint8 immediate = uint8(_data.proof[offset]);
        uint8 opCode = uint8(_data.proof[offset + 1]);
        uint popCount = opPopCount(opCode);
        Value.Data[] memory stackVals = new Value.Data[](popCount);
        offset += 2;

        require(immediate == 0 || immediate == 1, "Proof had bad operation type");
        if (immediate == 0) {
            startMachine.instructionStackHash = Value.HashOnly(Value.hashCodePointBasic(
                uint8(opCode),
                startMachine.instructionStackHash.hash
            ));
        } else {
            Value.Data memory immediateVal;
            (valid, offset, immediateVal) = Value.deserialize(_data.proof, offset);
            // string(abi.encodePacked("Proof had bad immediate value ", uint2str(valid)))
            require(valid, "Proof had bad immediate value");
            if (popCount > 0) {
                stackVals[0] = immediateVal;
            } else {
                endMachine.addDataStackValue(immediateVal);
            }

            startMachine.instructionStackHash = Value.HashOnly(Value.hashCodePointImmediate(
                uint8(opCode),
                immediateVal.hash().hash,
                startMachine.instructionStackHash.hash
            ));
        }

        uint i = 0;
        for (i = immediate; i < popCount; i++) {
            (valid, offset, stackVals[i]) = Value.deserialize(_data.proof, offset);
            require(valid, "Proof had bad stack value");
        }
        if (stackVals.length > 0) {
            for (i = 0; i < stackVals.length - immediate; i++) {
                startMachine.addDataStackValue(stackVals[stackVals.length - 1 - i]);
            }
        }
        return (opCode, stackVals, startMachine, endMachine, offset);
    }

    uint8 constant CODE_POINT_TYPECODE = 1;
    bytes32 constant CODE_POINT_ERROR = keccak256(
        abi.encodePacked(
            CODE_POINT_TYPECODE,
            uint8(0),
            bytes32(0)
        )
    );

    function checkProof(ValidateProofData memory _data) internal pure returns(uint) {
        uint8 opCode;
        bool valid;
        uint offset;
        Value.Data[] memory stackVals;
        Machine.Data memory startMachine;
        Machine.Data memory endMachine;
        (opCode, stackVals, startMachine, endMachine, offset) = loadMachine(_data);
        bool correct = true;
        bytes32 messageHash;
        require(_data.gas == opGasCost(opCode), "Invalid gas in proof");
        require((_data.didInboxInsn && opCode==OP_INBOX) || (!_data.didInboxInsn && opCode!=OP_INBOX),
            "Invalid didInboxInsn claim");
        if (startMachine.arbGasRemaining <= opGasCost(opCode)) {
            endMachine.arbGasRemaining = ((1<<128)+1)*((1<<128)-1); // = MaxUint256
            correct = false;
        } else {
            if (opCode == OP_ADD) {
                correct = executeAddInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_MUL) {
                correct = executeMulInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_SUB) {
                correct = executeSubInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_DIV) {
                correct = executeDivInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_SDIV) {
                correct = executeSdivInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_MOD) {
                correct = executeModInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_SMOD) {
                correct = executeSmodInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_ADDMOD) {
                correct = executeAddmodInsn(
                    endMachine,
                    stackVals[0],
                    stackVals[1],
                    stackVals[2]
                );
            } else if (opCode == OP_MULMOD) {
                correct = executeMulmodInsn(
                    endMachine,
                    stackVals[0],
                    stackVals[1],
                    stackVals[2]
                );
            } else if (opCode == OP_EXP) {
                correct = executeExpInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_LT) {
                correct = executeLtInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_GT) {
                correct = executeGtInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_SLT) {
                correct = executeSltInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_SGT) {
                correct = executeSgtInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_EQ) {
                correct = executeEqInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_ISZERO) {
                correct = executeIszeroInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_AND) {
                correct = executeAndInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_OR) {
                correct = executeOrInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_XOR) {
                correct = executeXorInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_NOT) {
                correct = executeNotInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_BYTE) {
                correct = executeByteInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_SIGNEXTEND) {
                correct = executeSignextendInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_SHA3) {
                correct = executeSha3Insn(endMachine, stackVals[0]);
            } else if (opCode == OP_TYPE) {
                correct = executeTypeInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_ETHHASH2) {
                correct = executeEthhash2Insn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_POP) {
                correct = executePopInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_SPUSH) {
                correct = executeSpushInsn(endMachine);
            } else if (opCode == OP_RPUSH) {
                correct = executeRpushInsn(endMachine);
            } else if (opCode == OP_RSET) {
                correct = executeRsetInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_JUMP) {
                correct = executeJumpInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_CJUMP) {
                correct = executeCjumpInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_STACKEMPTY) {
                correct = executeStackemptyInsn(endMachine);
            } else if (opCode == OP_PCPUSH) {
                correct = executePcpushInsn(endMachine, startMachine.instructionStackHash);
            } else if (opCode == OP_AUXPUSH) {
                correct = executeAuxpushInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_AUXPOP) {
                Value.Data memory auxVal;
                (valid, offset, auxVal) = Value.deserialize(_data.proof, offset);
                require(valid, "Proof of auxpop had bad aux value");
                startMachine.addAuxStackValue(auxVal);
                endMachine.addDataStackValue(auxVal);
            } else if (opCode == OP_AUXSTACKEMPTY) {
                correct = executeAuxstackemptyInsn(endMachine);
            } else if (opCode == OP_NOP) {
            } else if (opCode == OP_ERRPUSH) {
                correct = executeErrpushInsn(endMachine);
            } else if (opCode == OP_ERRSET) {
                correct = executeErrsetInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_DUP0) {
                correct = executeDup0Insn(endMachine, stackVals[0]);
            } else if (opCode == OP_DUP1) {
                correct = executeDup1Insn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_DUP2) {
                correct = executeDup2Insn(
                    endMachine,
                    stackVals[0],
                    stackVals[1],
                    stackVals[2]
                );
            } else if (opCode == OP_SWAP1) {
                correct = executeSwap1Insn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_SWAP2) {
                correct = executeSwap2Insn(
                    endMachine,
                    stackVals[0],
                    stackVals[1],
                    stackVals[2]
                );
            } else if (opCode == OP_TGET) {
                correct = executeTgetInsn(endMachine, stackVals[0], stackVals[1]);
            } else if (opCode == OP_TSET) {
                correct = executeTsetInsn(
                    endMachine,
                    stackVals[0],
                    stackVals[1],
                    stackVals[2]
                );
            } else if (opCode == OP_TLEN) {
                correct = executeTlenInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_MNEW) {
                correct = executeMnewInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_MGET) {
                (correct, offset) = executeMgetInsn(endMachine, stackVals[0], stackVals[1], _data.proof, offset);
            } else if (opCode == OP_BREAKPOINT) {
                correct = executeBreakpointInsn(endMachine);
            } else if (opCode == OP_LOG) {
                (correct, messageHash) = executeLogInsn(endMachine, stackVals[0]);
                if (correct) {
                    require(
                        keccak256(
                            abi.encodePacked(
                                _data.firstLog,
                                messageHash
                            )
                        ) == _data.lastLog,
                        "Logged value doesn't match output log"
                    );
                    require(_data.firstMessage == _data.lastMessage, "Send not called, but message is nonzero");
                } else {
                    messageHash = 0;
                }

            } else if (opCode == OP_SEND) {
                (correct, messageHash) = executeSendInsn(endMachine, stackVals[0]);
                if (correct) {
                    require(
                        keccak256(
                            abi.encodePacked(
                                _data.firstMessage,
                                messageHash
                            )
                        ) == _data.lastMessage,
                        "sent message doesn't match output message"
                    );

                    require(_data.firstLog == _data.lastLog, "Log not called, but message is nonzero");
                } else {
                    messageHash = 0;
                }
            } else if (opCode == OP_GETTIME) {
                Value.Data[] memory contents = new Value.Data[](2);
                contents[0] = Value.newInt(_data.timeBoundsBlocks[0]);
                contents[1] = Value.newInt(_data.timeBoundsBlocks[1]);
                endMachine.addDataStackValue(Value.newTuple(contents));
            } else if (opCode == OP_INBOX) {
                correct = executeInboxInsn(endMachine, stackVals[0], Value.HashOnly(_data.beforeInbox), _data.timeBoundsBlocks[0]);
            } else if (opCode == OP_ERROR) {
                correct = false;
            } else if (opCode == OP_STOP) {
                endMachine.setHalt();
            } else if (opCode == OP_SETGAS) {
                correct = executeSetgasInsn(endMachine, stackVals[0]);
            } else if (opCode == OP_PUSHGAS) {
                correct = executePushgasInsn(endMachine);
            }

            if (messageHash == 0) {
                require(_data.firstMessage == _data.lastMessage, "Send not called, but message is nonzero");
                require(_data.firstLog == _data.lastLog, "Log not called, but message is nonzero");
            }
            endMachine.arbGasRemaining = startMachine.arbGasRemaining - opGasCost(opCode);

        }

        if (!correct) {
            if (endMachine.errHandler.hash == CODE_POINT_ERROR) {
                endMachine.setErrorStop();
            } else {
                endMachine.instructionStackHash = endMachine.errHandler;
            }
        }

        // require(
        //     _data.beforeHash == startMachine.hash(),
        //     string(abi.encodePacked("Proof had non matching start state: ", startMachine.toString()))
        // );
        require(_data.beforeHash == startMachine.hash(), "Proof had non matching start state");
        // require(
        //     _data.afterHash == endMachine.hash(),
        //     string(abi.encodePacked("Proof had non matching end state: ", endMachine.toString(),
        //     
        require(_data.afterHash == endMachine.hash(), DebugPrint.uint2str(endMachine.arbGasRemaining));

        //require(_data.afterHash == endMachine.hash(), "Proof had non matching end state");

        return 0;
    }
}
