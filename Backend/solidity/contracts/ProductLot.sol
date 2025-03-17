// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./UserRegistry.sol";
import "./RawMilk.sol";
import "./Product.sol";
import "./Tracking.sol"; // ✅ Import Tracking.sol

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
    bool grade;
    string qualityAndNutritionCID;
    bytes32[] milkTankIds;
    ProductLotStatus status; // ✅ เพิ่มฟิลด์สถานะ
}

    enum ProductLotStatus { Created, InTransit, Received } // ✅ เพิ่มสถานะ
    mapping(bytes32 => ProductLotInfo) public productLots;
    bytes32[] public productLotIds;
    event ProductLotStatusUpdated(bytes32 indexed lotId, ProductLotStatus newStatus);


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
ProductLotStatus initialStatus = ProductLotStatus.Created; // Default
        productLots[_lotId] = ProductLotInfo({
    lotId: _lotId,
    productId: _productId,
    factory: msg.sender,
    inspector: _inspector,
    inspectionDate: block.timestamp,
    grade: _grade,
    qualityAndNutritionCID: _qualityAndNutritionCID,
    milkTankIds: _milkTankIds,
    status: initialStatus
    });


        productLotIds.push(_lotId);
        // หลังจากบันทึก ProductLot เสร็จ
    for (uint i = 0; i < _milkTankIds.length; i++) {
        rawMilkContract.markMilkTankAsUsed(_milkTankIds[i]); // ✅ Mark เป็น Used
    }

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
    Tracking public trackingContract; // ✅ Tracking Contract (เชื่อมภายหลัง)
function setTrackingContract(address _trackingContract) public {
    require(address(trackingContract) == address(0), "Tracking contract already set"); // ✅ ป้องกันการเปลี่ยน Tracking Contract
    trackingContract = Tracking(_trackingContract);
}
function updateProductLotStatus(bytes32 _lotId) public {
    require(productLots[_lotId].lotId != bytes32(0), "Error: Product Lot does not exist");
    require(address(trackingContract) != address(0), "Tracking contract not set");

    // ✅ รับค่าจาก getTrackingByLotId (เฉพาะ trackingIds เท่านั้น)
    (bytes32[] memory trackingIds, , ) = trackingContract.getTrackingByLotId(_lotId);
    require(trackingIds.length > 0, "Error: No tracking data found for this product lot");

    bool allReceived = true;
    bool hasInTransit = false;

    for (uint i = 0; i < trackingIds.length; i++) {
        // ✅ รับค่าจาก getTrackingById (เฉพาะ TrackingEvent เท่านั้น)
        (Tracking.TrackingEvent memory trackingEvent, , ) = trackingContract.getTrackingById(trackingIds[i]);

        if (trackingEvent.status == Tracking.TrackingStatus.InTransit) {
            hasInTransit = true;
        }
        if (trackingEvent.status != Tracking.TrackingStatus.Received) {
            allReceived = false;
        }
    }

    if (hasInTransit) {
        productLots[_lotId].status = ProductLotStatus.InTransit;
        emit ProductLotStatusUpdated(_lotId, ProductLotStatus.InTransit);
        return;
    }

    if (allReceived) {
        productLots[_lotId].status = ProductLotStatus.Received;
        emit ProductLotStatusUpdated(_lotId, ProductLotStatus.Received);
    }
}
function forceUpdateStatus(bytes32 _lotId, ProductLotStatus _status) public onlyFactory {
    require(productLots[_lotId].lotId != bytes32(0), "Error: Product Lot does not exist");
    productLots[_lotId].status = _status;
    emit ProductLotStatusUpdated(_lotId, _status);
}
function getProductLotStatus(bytes32 _lotId) public view returns (ProductLotStatus) {
    require(productLots[_lotId].lotId != bytes32(0), "Error: Product Lot does not exist");
    return productLots[_lotId].status;
}


}
