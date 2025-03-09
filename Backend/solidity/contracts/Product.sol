// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./UserRegistry.sol";

contract Product {
    UserRegistry public userRegistry;

    struct ProductInfo {
        bytes32 productId;
        address factoryWallet;
        string productName;
        string productCID;
        string category;
    }

    mapping(bytes32 => ProductInfo) public products;
    bytes32[] public productIds;

    event ProductCreated(
        bytes32 indexed productId,
        address indexed factoryWallet,
        string productName,
        string productCID,
        string category
    );

    event DebugLog(string message, address sender, bytes32 productId, string productName, string productCID, string category);

    //modifier onlyFactory() {
        //emit DebugLog("Checking Factory Role", msg.sender, bytes32(0), "", "", "");
        //require(
            //userRegistry.getUserRole(msg.sender) == UserRegistry.UserRole.Factory,
            //"Access denied: Only factories allowed"
        //);
        //emit DebugLog("Factory Role Verified", msg.sender, bytes32(0), "", "", "");
        //_;
    //}

    constructor(address _userRegistry) {
        userRegistry = UserRegistry(_userRegistry);
    }

    function createProduct(
        bytes32 _productId,
        string memory _productName,
        string memory _productCID,
        string memory _category
    ) public {
        emit DebugLog("Before Checks", msg.sender, _productId, _productName, _productCID, _category);

        //require(products[_productId].productId == bytes32(0), "Error: Product ID already exists");
        //emit DebugLog("Passed Product ID Check", msg.sender, _productId, _productName, _productCID, _category);

        //require(bytes(_category).length > 0, "Error: Category cannot be empty");
        //emit DebugLog("Passed Category Check", msg.sender, _productId, _productName, _productCID, _category);

        emit DebugLog("Proceeding to Create Product", msg.sender, _productId, _productName, _productCID, _category);

        products[_productId] = ProductInfo({
            productId: _productId,
            factoryWallet: msg.sender,
            productName: _productName,
            productCID: _productCID,
            category: _category
        });

        productIds.push(_productId);

        emit ProductCreated(_productId, msg.sender, _productName, _productCID, _category);
        emit DebugLog("Product Created Successfully", msg.sender, _productId, _productName, _productCID, _category);
    }
}
