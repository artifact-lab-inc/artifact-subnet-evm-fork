// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Votes.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract RelicToken is ERC20, ERC20Votes, Ownable {
    address public controller;
    
    event ControllerUpdated(address indexed previousController, address indexed newController);

    constructor(
        uint256 initialSupply,
        address initialController
    ) ERC20("RELIC", "RELIC") ERC20Permit("RELIC") {
        if(initialController == address(0)) revert("0");
        controller = initialController;
        _mint(msg.sender, initialSupply);
        _transferOwnership(msg.sender);
    }

    modifier onlyController() {
        if(msg.sender != controller) revert("1");
        _;
    }

    function mint(address to, uint256 amount) external onlyController {
        if(to == address(0)) revert("0");
        _mint(to, amount);
    }

    function setController(address newController) external onlyOwner {
        if(newController == address(0)) revert("0");
        address oldController = controller;
        controller = newController;
        emit ControllerUpdated(oldController, newController);
    }

    // The following functions are overrides required by Solidity for ERC20Votes
    function _afterTokenTransfer(address from, address to, uint256 amount) internal override(ERC20, ERC20Votes) {
        super._afterTokenTransfer(from, to, amount);
    }

    function _mint(address to, uint256 amount) internal override(ERC20, ERC20Votes) {
        super._mint(to, amount);
    }

    function _burn(address account, uint256 amount) internal override(ERC20, ERC20Votes) {
        super._burn(account, amount);
    }
} 