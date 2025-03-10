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

    modifier onlyFactory() {
        emit DebugLog("Checking Factory Role", msg.sender, bytes32(0), "", "", "");
        require(
            userRegistry.getUserRole(msg.sender) == UserRegistry.UserRole.Factory,
            "Access denied: Only factories allowed"
        );
        emit DebugLog("Factory Role Verified", msg.sender, bytes32(0), "", "", "");
        _;
    }

    constructor(address _userRegistry) {
        userRegistry = UserRegistry(_userRegistry);
    }

    function createProduct(
        bytes32 _productId,
        string memory _productName,
        string memory _productCID,
        string memory _category
    ) public onlyFactory {
        emit DebugLog("Before Checks", msg.sender, _productId, _productName, _productCID, _category);

        // ✅ ตรวจสอบว่า `productId` ต้องไม่ซ้ำ
        emit DebugLog("Checking Existing Product ID", msg.sender, _productId, _productName, _productCID, _category);
        require(products[_productId].productId == bytes32(0), "Error: Product ID already exists");
        emit DebugLog("Passed Product ID Check", msg.sender, _productId, _productName, _productCID, _category);

        // ✅ ตรวจสอบว่า `category` ต้องไม่เป็นค่าว่าง
        require(bytes(_category).length > 0, "Error: Category cannot be empty");
        emit DebugLog("Passed Category Check", msg.sender, _productId, _productName, _productCID, _category);

        // ✅ บันทึกข้อมูล Product
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
// ✅ ดึงสินค้าทั้งหมดของ Factory ที่เรียกใช้งาน (ชื่อ, Category, ID)
function getProductsByFactory() public view returns (
    bytes32[] memory,
    string[] memory,
    string[] memory
) {
    uint count = 0;
    for (uint i = 0; i < productIds.length; i++) {
        if (products[productIds[i]].factoryWallet == msg.sender) {
            count++;
        }
    }

    bytes32[] memory ids = new bytes32[](count);
    string[] memory names = new string[](count);
    string[] memory categories = new string[](count);
    
    uint index = 0;
    for (uint i = 0; i < productIds.length; i++) {
        if (products[productIds[i]].factoryWallet == msg.sender) {
            ids[index] = productIds[i];
            names[index] = products[productIds[i]].productName;
            categories[index] = products[productIds[i]].category;
            index++;
        }
    }

    return (ids, names, categories);
}


// ✅ ดึงรายละเอียดทั้งหมดของ Product โดยใช้ productId
function getProductDetails(bytes32 _productId) public view returns (ProductInfo memory) {
    require(products[_productId].productId != bytes32(0), "Product does not exist");
    return products[_productId];
}


}
