// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./RelicToken.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract ArtifactEmissionManager is ReentrancyGuard {
    RelicToken public immutable relicToken;
    address public governor;
    uint256 public emissionRatePerSecond;
    uint256 public totalEmitted;
    uint256 public immutable maxCap;
    uint256 public lastDripTimestamp;

    event EmissionRateUpdated(uint256 oldRate, uint256 newRate);
    event GovernanceTransferred(address indexed previousGovernor, address indexed newGovernor);
    event TokensDripped(address indexed recipient, uint256 amount);

    modifier onlyGovernor() {
        if(msg.sender != governor) revert("1");
        _;
    }

    constructor(
        address _relicToken,
        address _initialGovernor,
        uint256 _initialEmissionRate,
        uint256 _maxCap
    ) {
        if(_relicToken == address(0)) revert("0");
        if(_initialGovernor == address(0)) revert("0");
        if(_maxCap == 0) revert("2");

        relicToken = RelicToken(_relicToken);
        governor = _initialGovernor;
        emissionRatePerSecond = _initialEmissionRate;
        maxCap = _maxCap;
        lastDripTimestamp = block.timestamp;
    }

    function _calculateDripAmount() internal view returns (uint256) {
        if (totalEmitted >= maxCap) return 0;
        
        uint256 timePassed = block.timestamp - lastDripTimestamp;
        uint256 emissionAmount = timePassed * emissionRatePerSecond;

        // Cap the emission if it would exceed maxCap
        if (totalEmitted + emissionAmount > maxCap) {
            emissionAmount = maxCap - totalEmitted;
        }

        return emissionAmount;
    }

    function drip(address recipient) external nonReentrant returns (uint256) {
        if(recipient == address(0)) revert("0");
        
        uint256 emissionAmount = _calculateDripAmount();

        if (emissionAmount > 0) {
            totalEmitted += emissionAmount;
            lastDripTimestamp = block.timestamp;
            relicToken.mint(recipient, emissionAmount);
            emit TokensDripped(recipient, emissionAmount);
        }

        return emissionAmount;
    }

    function updateRate(uint256 newRate) external onlyGovernor {
        // Handle any pending emissions before updating the rate
        uint256 pendingAmount = _calculateDripAmount();
        if (pendingAmount > 0) {
            totalEmitted += pendingAmount;
            lastDripTimestamp = block.timestamp;
            relicToken.mint(governor, pendingAmount);
            emit TokensDripped(governor, pendingAmount);
        }
        
        uint256 oldRate = emissionRatePerSecond;
        emissionRatePerSecond = newRate;
        emit EmissionRateUpdated(oldRate, newRate);
    }

    function transferGovernance(address newGovernor) external onlyGovernor {
        if(newGovernor == address(0)) revert("0");
        address oldGovernor = governor;
        governor = newGovernor;
        emit GovernanceTransferred(oldGovernor, newGovernor);
    }

    function getEmissionInfo() external view returns (
        uint256 currentRate,
        uint256 emitted,
        uint256 remaining,
        uint256 cap
    ) {
        return (
            emissionRatePerSecond,
            totalEmitted,
            maxCap - totalEmitted,
            maxCap
        );
    }
} 