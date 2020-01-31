# Copyright 2019-2020, Offchain Labs, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from .. import std, value
from ..annotation import modifies_stack
from . import os, accounts, call_frame, call_finish
from .types import (
    local_exec_state,
    eth_transfer_message,
    token_transfer_message,
    message,
    ethbridge_message,
)

WITHDRAW_ETH_TYPECODE = 1
WITHDRAW_ERC20_TYPECODE = 2
WITHDRAW_ERC721_TYPECODE = 3


@modifies_stack([local_exec_state.typ], [value.IntType()])
def perform_precompile_call(vm):
    os.get_call_frame(vm)
    vm.dup0()
    call_frame.call_frame.set_val("parent_frame")(vm)
    std.sized_byterange.new(vm)
    vm.swap1()
    call_frame.call_frame.set_val("memory")(vm)
    os.set_call_frame(vm)
    # local_exec_state
    vm.dup0()
    local_exec_state.get("data")(vm)
    vm.push(0)
    vm.swap1()
    std.sized_byterange.get(vm)
    vm.push(224)
    vm.swap1()
    std.bitwise.shift_right(vm)
    vm.dup0()
    vm.push(0x1B9A91A4)
    vm.eq()
    vm.ifelse(
        lambda vm: [vm.pop(), withdraw_eth_interrupt(vm)],
        lambda vm: [
            vm.dup0(),
            vm.push(0xA1DB9782),
            vm.eq(),
            vm.ifelse(
                lambda vm: [vm.pop(), withdraw_erc20_interrupt(vm)],
                lambda vm: [
                    vm.dup0(),
                    vm.push(0xF3E414F8),
                    vm.eq(),
                    vm.ifelse(
                        lambda vm: [vm.pop(), withdraw_erc721_interrupt(vm)],
                        lambda vm: [
                            vm.dup0(),
                            vm.push(0xBDE19776),
                            vm.eq(),
                            vm.ifelse(
                                lambda vm: [vm.pop(), arbsys_time_upper_bound(vm)],
                                lambda vm: [
                                    vm.dup0(),
                                    vm.push(0x44F50653),
                                    vm.eq(),
                                    vm.ifelse(
                                        lambda vm: [
                                            vm.pop(),
                                            arbsys_current_message_time(vm),
                                        ],
                                        lambda vm: [
                                            vm.dup0(),
                                            vm.push(0x23CA0CD2),
                                            vm.eq(),
                                            vm.ifelse(
                                                lambda vm: [
                                                    vm.pop(),
                                                    transaction_count_interrupt(vm),
                                                ],
                                                lambda vm: [
                                                    vm.dup0(),
                                                    vm.push(0x474ED9C0),
                                                    vm.eq(),
                                                    vm.ifelse(
                                                        lambda vm: [
                                                            vm.pop(),
                                                            clone_contract_interrupt(
                                                                vm
                                                            ),
                                                        ],
                                                        lambda vm: [
                                                            vm.pop(),
                                                            vm.pop(),
                                                            vm.push(0),
                                                        ],
                                                    ),
                                                ],
                                            ),
                                        ],
                                    ),
                                ],
                            ),
                        ],
                    ),
                ],
            ),
        ],
    )


@modifies_stack([local_exec_state.typ], [value.IntType()] * 3)
def parse_withdraw_call(vm):
    # local_exec_state
    vm.dup0()
    local_exec_state.get("caller")(vm)
    vm.swap1()
    local_exec_state.get("data")(vm)
    vm.dup0()
    vm.push(4)
    vm.swap1()
    std.sized_byterange.get(vm)
    # dest data sender
    vm.swap1()
    vm.push(36)
    vm.swap1()
    std.sized_byterange.get(vm)
    # amount dest sender


@modifies_stack([local_exec_state.typ], [value.IntType()])
def withdraw_eth_interrupt(vm):
    # local_exec_state
    parse_withdraw_call(vm)
    # amount dest sender
    vm.dup0()
    os.process_eth_withdraw(vm)
    vm.ifelse(
        lambda vm: [
            vm.push(eth_transfer_message.make()),
            vm.cast(eth_transfer_message.typ),
            eth_transfer_message.set_val("amount")(vm),
            eth_transfer_message.set_val("dest")(vm),
            # token_transfer_message sender
            vm.push(WITHDRAW_ETH_TYPECODE),
            vm.push(message.make()),
            vm.cast(message.typ),
            message.set_val("type")(vm),
            message.set_val("message")(vm),
            message.set_val("sender")(vm),
            vm.cast(value.TupleType()),
            vm.send(),
            call_finish.stop_impl(vm),
        ],
        lambda vm: [
            vm.pop(),
            vm.pop(),
            vm.pop(),
            vm.push(0),
            vm.push(0),
            call_finish.revert_impl(vm),
        ],
    )


def withdraw_token_interrupt(vm, token_type):
    # local_exec_state
    parse_withdraw_call(vm)
    # amount dest token_address
    vm.push(token_transfer_message.make())
    vm.cast(token_transfer_message.typ)
    token_transfer_message.set_val("amount")(vm)
    token_transfer_message.set_val("dest")(vm)
    token_transfer_message.set_val("token_address")(vm)
    os.message_caller(vm)
    vm.push(token_type)
    vm.push(message.make())
    vm.cast(message.typ)
    message.set_val("type")(vm)
    message.set_val("sender")(vm)
    message.set_val("message")(vm)
    vm.cast(value.TupleType())
    vm.send()
    call_finish.stop_impl(vm)


@modifies_stack([local_exec_state.typ], [value.IntType()])
def withdraw_erc20_interrupt(vm):
    # local_exec_state
    withdraw_token_interrupt(vm, WITHDRAW_ERC20_TYPECODE)


@modifies_stack([local_exec_state.typ], [value.IntType()])
def withdraw_erc721_interrupt(vm):
    # local_exec_state
    withdraw_token_interrupt(vm, WITHDRAW_ERC721_TYPECODE)


@modifies_stack([value.IntType()], [value.IntType()])
def return_one_uint_to_solidity_caller(vm):
    vm.push(0)
    os.get_call_frame(vm)
    call_frame.call_frame.get("memory")(vm)
    std.sized_byterange.set_val(vm)
    os.get_call_frame(vm)
    call_frame.call_frame.set_val("memory")(vm)
    os.set_call_frame(vm)
    vm.push(32)
    vm.push(0)
    call_finish.ret_impl(vm)


@modifies_stack([local_exec_state.typ], [value.IntType()])
def arbsys_time_upper_bound(vm):
    vm.pop()
    vm.gettime()
    vm.tgetn(1)
    return_one_uint_to_solidity_caller(vm)


@modifies_stack([local_exec_state.typ], [value.IntType()])
def arbsys_current_message_time(vm):
    vm.pop()
    os.get_chain_state(vm)
    os.chain_state.get("global_exec_state")(vm)
    os.global_exec_state.get("current_msg")(vm)
    ethbridge_message.get("block_number")(vm)
    return_one_uint_to_solidity_caller(vm)


@modifies_stack([local_exec_state.typ], [value.IntType()])
def transaction_count_interrupt(vm):
    local_exec_state.get("data")(vm)
    vm.push(4)
    vm.swap1()
    std.sized_byterange.get(vm)
    # address
    os.get_call_frame(vm)
    os.call_frame.call_frame.get("accounts")(vm)
    accounts.account_store.get(vm)
    accounts.account_state.get("nextSeqNum")(vm)
    return_one_uint_to_solidity_caller(vm)


@modifies_stack([local_exec_state.typ], [value.IntType()])
def clone_contract_interrupt(vm):
    # local_exec_state
    local_exec_state.get("data")(vm)
    vm.push(4)
    vm.swap1()
    std.sized_byterange.get(vm)
    # address
    os.create_clone_contract(vm)
    return_one_uint_to_solidity_caller(vm)
