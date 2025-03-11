// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./UserRegistry.sol";
import "./RawMilk.sol";
import "./Product.sol";

contract ProductLot {
    UserRegistry public userRegistry;
    RawMilk public rawMilkContract;
    Product public productContract;

    struct ProductLotInfo {
        bytes32 lotId;
        bytes32 productId;
        address factory;
        string inspector;
        uint256 inspectionDate;
        bool grade; // true = ผ่าน, false = ไม่ผ่าน
        string qualityAndNutritionCID;
        bytes32[] milkTankIds;
    }

    mapping(bytes32 => ProductLotInfo) public productLots;
    bytes32[] public productLotIds;

    event ProductLotCreated(
        bytes32 indexed lotId,
        bytes32 indexed productId,
        address indexed factory,
        string inspector,
        uint256 inspectionDate,
        bool grade,
        string qualityAndNutritionCID,
        bytes32[] milkTankIds
    );

    modifier onlyFactory() {
        require(
            userRegistry.getUserRole(msg.sender) == UserRegistry.UserRole.Factory,
            "Access denied: Only factories allowed"
        );
        _;
    }

    constructor(address _userRegistry, address _rawMilkContract, address _productContract) {
        userRegistry = UserRegistry(_userRegistry);
        rawMilkContract = RawMilk(_rawMilkContract);
        productContract = Product(_productContract);
    }

    function createProductLot(
        bytes32 _lotId,
        bytes32 _productId,
        string memory _inspector,
        bool _grade,
        string memory _qualityAndNutritionCID,
        bytes32[] memory _milkTankIds
    ) public onlyFactory {
        require(productLots[_lotId].lotId == bytes32(0), "Error: Product Lot ID already exists");
        require(productContract.getProductDetails(_productId).productId != bytes32(0), "Error: Product ID does not exist");

        for (uint i = 0; i < _milkTankIds.length; i++) {
            require(
                rawMilkContract.getMilkTank(_milkTankIds[i]).status == RawMilk.MilkStatus.Approved,
                "Error: Milk Tank must be approved"
            );
        }

        productLots[_lotId] = ProductLotInfo({
            lotId: _lotId,
            productId: _productId,
            factory: msg.sender,
            inspector: _inspector,
            inspectionDate: block.timestamp,
            grade: _grade,
            qualityAndNutritionCID: _qualityAndNutritionCID,
            milkTankIds: _milkTankIds
        });

        productLotIds.push(_lotId);

        emit ProductLotCreated(_lotId, _productId, msg.sender, _inspector, block.timestamp, _grade, _qualityAndNutritionCID, _milkTankIds);
    }

    function getProductLot(bytes32 _lotId) public view returns (ProductLotInfo memory) {
        require(productLots[_lotId].lotId != bytes32(0), "Error: Product Lot does not exist");
        return productLots[_lotId];
    }

    function isProductLotExists(bytes32 _lotId) public view returns (bool) {
        return productLots[_lotId].lotId != bytes32(0);
    }

    function getProductLotsByFactory(address _factory) public view returns (bytes32[] memory) {
        uint count = 0;
        for (uint i = 0; i < productLotIds.length; i++) {
            if (productLots[productLotIds[i]].factory == _factory) {
                count++;
            }
        }

        bytes32[] memory factoryLots = new bytes32[](count);
        uint index = 0;
        for (uint i = 0; i < productLotIds.length; i++) {
            if (productLots[productLotIds[i]].factory == _factory) {
                factoryLots[index] = productLotIds[i];
                index++;
            }
        }

        return factoryLots;
    }

    function getMilkTanksByProductLot(bytes32 _lotId) public view returns (bytes32[] memory) {
        require(productLots[_lotId].lotId != bytes32(0), "Error: Product Lot does not exist");
        return productLots[_lotId].milkTankIds;
    }
}
