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

import "../libraries/CloneFactory.sol";

import "./IChallengeFactory.sol";
import "./IBisectionChallenge.sol";
import "./ChallengeType.sol";


contract ChallengeFactory is CloneFactory, ChallengeType, IChallengeFactory {

    // Invalid challenge type
    string constant INVALID_TYPE = "INVALID_TYPE";

    address public messagesChallengeTemplate;
    address public pendingTopChallengeTemplate;
    address public executionChallengeTemplate;

    constructor(
        address _messagesChallengeTemplate,
        address _pendingTopChallengeTemplate,
        address _executionChallengeTemplate
    ) public {
        messagesChallengeTemplate = _messagesChallengeTemplate;
        pendingTopChallengeTemplate = _pendingTopChallengeTemplate;
        executionChallengeTemplate = _executionChallengeTemplate;
    }

    function generateCloneAddress(
        address asserter,
        address challenger,
        uint256 challengeType
    )
        public
        view
        returns(address)
    {
        return address(
            uint160(
                uint256(
                    keccak256(
                        abi.encodePacked(
                            bytes1(0xff),
                            address(this),
                            generateNonce(asserter, challenger),
                            cloneCodeHash(getChallengeTemplate(challengeType))
                        )
                    )
                )
            )
        );
    }

    function createChallenge(
        address payable _asserter,
        address payable _challenger,
        uint256 _challengePeriodTicks,
        bytes32 _challengeHash,
        uint256 challengeType
    )
        external
        returns(address)
    {
        address challengeTemplate = getChallengeTemplate(challengeType);
        address clone = createClone(challengeTemplate);
        IBisectionChallenge(clone).initializeBisection(
            msg.sender,
            _asserter,
            _challenger,
            _challengePeriodTicks,
            _challengeHash
        );
        return address(clone);
    }

    function generateNonce(address asserter, address challenger) private view returns(uint) {
        return uint(
            keccak256(
                abi.encodePacked(
                    asserter,
                    challenger,
                    msg.sender
                )
            )
        );
    }

    function getChallengeTemplate(uint256 challengeType) private view returns(address) {
        if (challengeType == INVALID_PENDING_TOP_TYPE) {
            return pendingTopChallengeTemplate;
        } else if (challengeType == INVALID_MESSAGES_TYPE) {
            return messagesChallengeTemplate;
        } else if (challengeType == INVALID_EXECUTION_TYPE) {
            return executionChallengeTemplate;
        } else {
            require(false, INVALID_TYPE);
        }
    }
}
