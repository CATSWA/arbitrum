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

#include "config.hpp"
#include "helper.hpp"

#include <avm/machine.hpp>
#include <data_storage/checkpoint/checkpointstorage.hpp>

#include <catch2/catch.hpp>

#include <boost/filesystem.hpp>

auto execution_path = boost::filesystem::current_path();

void checkpointState(CheckpointStorage& storage, Machine& machine) {
    auto results = machine.checkpoint(storage);

    std::vector<unsigned char> hash_vector;
    marshal_uint256_t(machine.hash(), hash_vector);

    REQUIRE(results.status.ok());
    REQUIRE(results.reference_count == 1);
    REQUIRE(results.storage_key == hash_vector);
}

void checkpointStateTwice(CheckpointStorage& storage, Machine& machine) {
    auto results = machine.checkpoint(storage);
    auto results2 = machine.checkpoint(storage);

    std::vector<unsigned char> hash_vector;
    marshal_uint256_t(machine.hash(), hash_vector);

    REQUIRE(results2.status.ok());
    REQUIRE(results2.reference_count == 2);
    REQUIRE(results2.storage_key == hash_vector);
}

void deleteCheckpoint(CheckpointStorage& storage,
                      Machine& machine,
                      const std::vector<unsigned char>& checkpoint_key) {
    auto results = machine.deleteCheckpoint(storage);
    REQUIRE(results.status.ok());
    REQUIRE(results.reference_count == 0);
}

void restoreCheckpoint(CheckpointStorage& storage,
                       Machine& expected_machine,
                       const std::vector<unsigned char>& checkpoint_key) {
    Machine machine;
    machine.initializeMachine(test_contract_path);
    auto success = machine.restoreCheckpoint(storage, checkpoint_key);

    REQUIRE(success);
    REQUIRE(machine.hash() == expected_machine.hash());
}

TEST_CASE("Checkpoint State") {
    SECTION("default") {
        DBDeleter deleter;
        TuplePool pool;
        CheckpointStorage storage(dbpath, test_contract_path);
        Machine machine;

        bool initialized = machine.initializeMachine(test_contract_path);
        REQUIRE(initialized);

        checkpointState(storage, machine);
    }
    SECTION("save twice") {
        DBDeleter deleter;
        TuplePool pool;
        CheckpointStorage storage(dbpath, test_contract_path);
        Machine machine;
        machine.initializeMachine(test_contract_path);

        checkpointStateTwice(storage, machine);
    }
    SECTION("assert machine hash") {
        DBDeleter deleter;
        Machine machine;
        machine.initializeMachine(test_contract_path);

        CheckpointStorage storage(dbpath, test_contract_path);
        auto initial_machine = storage.getInitialVmValues();
        MachineState machine_state(initial_machine.code,
                                   initial_machine.staticVal, storage.pool);

        auto machine2 = new Machine();
        machine2->initializeMachine(machine_state);

        auto hash1 = machine.hash();
        auto hash2 = machine2->hash();

        machine.checkpoint(storage);

        std::vector<unsigned char> hash_vector;
        marshal_uint256_t(hash1, hash_vector);

        Machine machine3;
        machine3.restoreCheckpoint(storage, hash_vector);

        auto hash3 = machine3.hash();

        REQUIRE(hash1 == hash2);
        REQUIRE(hash3 == hash2);
    }
}

TEST_CASE("Delete machine checkpoint") {
    SECTION("default") {
        DBDeleter deleter;
        TuplePool pool;
        CheckpointStorage storage(dbpath, test_contract_path);
        Machine machine;
        machine.initializeMachine(test_contract_path);
        auto results = machine.checkpoint(storage);

        deleteCheckpoint(storage, machine, results.storage_key);
    }
}

TEST_CASE("Restore checkpoint") {
    SECTION("default") {
        DBDeleter deleter;
        TuplePool pool;
        CheckpointStorage storage(dbpath, test_contract_path);
        Machine machine;
        machine.initializeMachine(test_contract_path);
        auto results = machine.checkpoint(storage);

        restoreCheckpoint(storage, machine, results.storage_key);
    }
}
