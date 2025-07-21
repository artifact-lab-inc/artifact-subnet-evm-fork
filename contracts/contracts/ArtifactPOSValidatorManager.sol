// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./RelicToken.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";

contract ArtifactPOSValidatorManager is ReentrancyGuard {
    using Math for uint256;

    RelicToken public immutable relicToken;
    
    struct ValidatorStake {
        uint256 amount;
        uint256 lastStakeTimestamp;
    }

    mapping(address => ValidatorStake) public validatorStakes;
    uint256 public totalStaked;
    uint256 public immutable minimumStake;
    uint256 public constant UNSTAKE_DELAY = 14 days;

    event Staked(address indexed validator, uint256 amount);
    event Unstaked(address indexed validator, uint256 amount);
    event StakeIncreased(address indexed validator, uint256 additionalAmount);

    constructor(
        address _relicToken,
        uint256 _minimumStake
    ) {
        if(_relicToken == address(0)) revert("0");
        if(_minimumStake == 0) revert("2");
        
        relicToken = RelicToken(_relicToken);
        minimumStake = _minimumStake;
    }

    function stake(uint256 amount) external nonReentrant {
        if(amount == 0) revert("2");
        
        ValidatorStake storage validatorStake = validatorStakes[msg.sender];
        uint256 newStakeAmount = validatorStake.amount + amount;
        if(newStakeAmount < minimumStake) revert("3");

        // Transfer tokens from validator to contract
        if(!relicToken.transferFrom(msg.sender, address(this), amount)) revert("4");

        // Update validator stake info
        if (validatorStake.amount == 0) {
            validatorStake.lastStakeTimestamp = block.timestamp;
            emit Staked(msg.sender, amount);
        } else {
            emit StakeIncreased(msg.sender, amount);
        }
        
        validatorStake.amount = newStakeAmount;
        totalStaked += amount;
    }

    function unstake(uint256 amount) external nonReentrant {
        ValidatorStake storage validatorStake = validatorStakes[msg.sender];
        if(validatorStake.amount < amount) revert("5");
        if(block.timestamp < validatorStake.lastStakeTimestamp + UNSTAKE_DELAY) revert("6");

        uint256 remainingStake = validatorStake.amount - amount;
        if(remainingStake != 0 && remainingStake < minimumStake) revert("3");

        // Update validator stake info
        validatorStake.amount = remainingStake;
        totalStaked -= amount;

        // Transfer tokens back to validator
        if(!relicToken.transfer(msg.sender, amount)) revert("4");
        
        emit Unstaked(msg.sender, amount);
    }

    function getValidatorStake(address validator) external view returns (
        uint256 amount,
        uint256 lastStakeTime,
        bool canUnstake,
        uint256 unstakeAvailableTime
    ) {
        ValidatorStake memory validatorStake = validatorStakes[validator];
        return (
            validatorStake.amount,
            validatorStake.lastStakeTimestamp,
            block.timestamp >= validatorStake.lastStakeTimestamp + UNSTAKE_DELAY,
            validatorStake.lastStakeTimestamp + UNSTAKE_DELAY
        );
    }
} 