// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./UserRegistry.sol";
import "./ProductLot.sol";

contract Tracking {
    UserRegistry public userRegistry;
    ProductLot public productLotContract;

    enum TrackingStatus { Pending, InTransit, Received }
    enum LogisticsCheckType { Before, During, After }

    struct TrackingEvent {
        bytes32 trackingId;
        bytes32 productLotId;
        string retailerId; // ✅ เก็บไอดีของร้านค้าแทน Address
        string qrCodeCID;
        TrackingStatus status;
    }

    struct LogisticsCheckpoint {
    bytes32 trackingId;
    address logisticsProvider;
    uint256 pickupTime;
    uint256 deliveryTime;
    uint256 quantity;
    int256 temperature;
    string personInCharge;
    LogisticsCheckType checkType;
    string receiverCID; // ✅ บันทึกข้อมูลผู้รับสินค้า (IPFS CID)
}

    struct RetailerConfirmation {
        bytes32 trackingId;
        string retailerId; // ✅ ใช้ Retailer ID แทน Address
        uint256 receivedTime;
        string qualityCID;
        string personInCharge;
    }

    mapping(bytes32 => TrackingEvent) public trackingEvents;
    mapping(bytes32 => LogisticsCheckpoint[]) public logisticsCheckpoints;
    mapping(bytes32 => RetailerConfirmation) public retailerConfirmations;
    mapping(string => bytes32[]) public retailerTracking; // ✅ ใช้ Retailer ID เป็น key
   bytes32[] public allTrackingIds; // ✅ เปลี่ยนชื่อให้ไม่ชนกับตัวแปรในฟังก์ชัน



    event TrackingCreated(bytes32 indexed trackingId, bytes32 indexed productLotId, string retailerId, string qrCodeCID);
    event LogisticsUpdated(
    bytes32 indexed trackingId,
    address indexed logisticsProvider,
    LogisticsCheckType checkType,
    string receiverCID // ✅ เพิ่มในอีเวนต์
);
    event RetailerReceived(bytes32 indexed trackingId, string retailerId, string qualityCID);

    modifier onlyLogistics() {
        require(userRegistry.getUserRole(msg.sender) == UserRegistry.UserRole.Logistics, "Access denied: Only logistics allowed");
        _;
    }

    modifier onlyRetailer(string memory _retailerId) {
        require(bytes(_retailerId).length > 0, "Access denied: Invalid retailer ID");
        _;
    }

    constructor(address _userRegistry, address _productLotContract) {
        userRegistry = UserRegistry(_userRegistry);
        productLotContract = ProductLot(_productLotContract);
    }

    function createTrackingEvent(
        bytes32 _trackingId,
        bytes32 _productLotId,
        string memory _retailerId,
        string memory _qrCodeCID
    ) public {
        require(productLotContract.isProductLotExists(_productLotId), "Error: Product Lot does not exist");
        require(trackingEvents[_trackingId].trackingId == bytes32(0), "Error: Tracking ID already exists");
        require(bytes(_retailerId).length > 0, "Error: Invalid Retailer ID");
        require(bytes(_qrCodeCID).length > 0, "Error: QR Code CID cannot be empty");

        trackingEvents[_trackingId] = TrackingEvent({
            trackingId: _trackingId,
            productLotId: _productLotId,
            retailerId: _retailerId,
            qrCodeCID: _qrCodeCID,
            status: TrackingStatus.Pending
        });

        allTrackingIds.push(_trackingId);
        retailerTracking[_retailerId].push(_trackingId);
        emit TrackingCreated(_trackingId, _productLotId, _retailerId, _qrCodeCID);
    }

    function updateLogisticsCheckpoint(
    bytes32 _trackingId,
    uint256 _pickupTime,
    uint256 _deliveryTime,
    uint256 _quantity,
    int256 _temperature,
    string memory _personInCharge,
    LogisticsCheckType _checkType,
    string memory _receiverCID
) public onlyLogistics {
    require(trackingEvents[_trackingId].trackingId != bytes32(0), "Error: Tracking ID does not exist");

    logisticsCheckpoints[_trackingId].push(LogisticsCheckpoint({
        trackingId: _trackingId,
        logisticsProvider: msg.sender,
        pickupTime: _pickupTime,
        deliveryTime: _deliveryTime,
        quantity: _quantity,
        temperature: _temperature,
        personInCharge: _personInCharge,
        checkType: _checkType,
        receiverCID: _receiverCID
    }));

    // ✅ อัปเดตสถานะเป็น InTransit ถ้าสถานะเดิมเป็น Pending
    if (trackingEvents[_trackingId].status == TrackingStatus.Pending) {
        trackingEvents[_trackingId].status = TrackingStatus.InTransit;
        // ✅ อัปเดต Product Lot ให้เป็น InTransit
        productLotContract.updateProductLotStatus(trackingEvents[_trackingId].productLotId);
    }

    emit LogisticsUpdated(_trackingId, msg.sender, _checkType, _receiverCID);
}



    function retailerReceiveProduct(
    bytes32 _trackingId,
    string memory _retailerId,
    string memory _qualityCID,
    string memory _personInCharge
) public onlyRetailer(_retailerId) {
    require(trackingEvents[_trackingId].trackingId != bytes32(0), "Error: Tracking ID does not exist");
    require(keccak256(abi.encodePacked(trackingEvents[_trackingId].retailerId)) == keccak256(abi.encodePacked(_retailerId)), "Error: Only assigned retailer can receive");

    trackingEvents[_trackingId].status = TrackingStatus.Received;
    retailerConfirmations[_trackingId] = RetailerConfirmation({
        trackingId: _trackingId,
        retailerId: _retailerId,
        receivedTime: block.timestamp,
        qualityCID: _qualityCID,
        personInCharge: _personInCharge
    });

    // ✅ อัปเดตสถานะของ Product Lot ถ้าทุก Tracking ถูก "Received"
    productLotContract.updateProductLotStatus(trackingEvents[_trackingId].productLotId);

    emit RetailerReceived(_trackingId, _retailerId, _qualityCID);
}


    function getTrackingById(bytes32 _trackingId) public view returns (TrackingEvent memory, LogisticsCheckpoint[] memory, RetailerConfirmation memory) {
        require(trackingEvents[_trackingId].trackingId != bytes32(0), "Error: Tracking ID does not exist");
        return (trackingEvents[_trackingId], logisticsCheckpoints[_trackingId], retailerConfirmations[_trackingId]);
    }

    function getTrackingByRetailer(string memory _retailerId) public view returns (bytes32[] memory) {
        return retailerTracking[_retailerId];
    }
function getTrackingByLotId(bytes32 _productLotId) 
    public view 
    returns (bytes32[] memory resultTrackingIds, string[] memory retailerIds, string[] memory qrCodeCIDs) 
{
    uint count = 0;

    // ✅ ใช้ `allTrackingIds` เพื่อนับจำนวน Tracking Events
    for (uint i = 0; i < allTrackingIds.length; i++) {
        if (trackingEvents[allTrackingIds[i]].productLotId == _productLotId) {
            count++;
        }
    }

    // ✅ สร้างอาเรย์สำหรับเก็บผลลัพธ์
    resultTrackingIds = new bytes32[](count);
    retailerIds = new string[](count);
    qrCodeCIDs = new string[](count);

    uint index = 0;
    for (uint i = 0; i < allTrackingIds.length; i++) {
        if (trackingEvents[allTrackingIds[i]].productLotId == _productLotId) {
            resultTrackingIds[index] = allTrackingIds[i];
            retailerIds[index] = trackingEvents[allTrackingIds[i]].retailerId;
            qrCodeCIDs[index] = trackingEvents[allTrackingIds[i]].qrCodeCID;
            index++;
        }
    }

    return (resultTrackingIds, retailerIds, qrCodeCIDs);
}
function getOngoingShipmentsByLogistics() 
    public view onlyLogistics 
    returns (bytes32[] memory trackingIds, string[] memory personInChargeList) 
{
    uint count = 0;

    // ✅ นับจำนวน Tracking ที่อยู่ในสถานะ InTransit และอัปเดตล่าสุดโดยผู้เรียก
    for (uint i = 0; i < allTrackingIds.length; i++) {
        bytes32 trackingId = allTrackingIds[i];

        if (
            trackingEvents[trackingId].status == TrackingStatus.InTransit && 
            logisticsCheckpoints[trackingId].length > 0 &&
            logisticsCheckpoints[trackingId][logisticsCheckpoints[trackingId].length - 1].logisticsProvider == msg.sender
        ) {
            count++;
        }
    }

    // ✅ สร้างอาร์เรย์สำหรับเก็บผลลัพธ์
    trackingIds = new bytes32[](count);
    personInChargeList = new string[](count);
    uint index = 0;

    for (uint i = 0; i < allTrackingIds.length; i++) {
        bytes32 trackingId = allTrackingIds[i];

        if (
            trackingEvents[trackingId].status == TrackingStatus.InTransit && 
            logisticsCheckpoints[trackingId].length > 0 &&
            logisticsCheckpoints[trackingId][logisticsCheckpoints[trackingId].length - 1].logisticsProvider == msg.sender
        ) {
            trackingIds[index] = trackingId;
            personInChargeList[index] = logisticsCheckpoints[trackingId][logisticsCheckpoints[trackingId].length - 1].personInCharge;
            index++;
        }
    }

    return (trackingIds, personInChargeList);
}

function getAllTrackingIds() public view returns (bytes32[] memory) {
    return allTrackingIds;
}
function getLogisticsCheckpointsByTrackingId(bytes32 _trackingId) 
    public view 
    returns (
        LogisticsCheckpoint[] memory beforeCheckpoints, 
        LogisticsCheckpoint[] memory duringCheckpoints, 
        LogisticsCheckpoint[] memory afterCheckpoints
    ) 
{
    require(trackingEvents[_trackingId].trackingId != bytes32(0), "Error: Tracking ID does not exist");

    uint beforeCount = 0;
    uint duringCount = 0;
    uint afterCount = 0;

    // ✅ นับจำนวนแต่ละประเภทก่อนสร้างอาเรย์
    for (uint i = 0; i < logisticsCheckpoints[_trackingId].length; i++) {
        if (logisticsCheckpoints[_trackingId][i].checkType == LogisticsCheckType.Before) {
            beforeCount++;
        } else if (logisticsCheckpoints[_trackingId][i].checkType == LogisticsCheckType.During) {
            duringCount++;
        } else if (logisticsCheckpoints[_trackingId][i].checkType == LogisticsCheckType.After) {
            afterCount++;
        }
    }

    // ✅ สร้างอาเรย์สำหรับแต่ละประเภท
    beforeCheckpoints = new LogisticsCheckpoint[](beforeCount);
    duringCheckpoints = new LogisticsCheckpoint[](duringCount);
    afterCheckpoints = new LogisticsCheckpoint[](afterCount);

    uint beforeIndex = 0;
    uint duringIndex = 0;
    uint afterIndex = 0;

    // ✅ แยกข้อมูลเข้าไปในแต่ละอาเรย์
    for (uint i = 0; i < logisticsCheckpoints[_trackingId].length; i++) {
        if (logisticsCheckpoints[_trackingId][i].checkType == LogisticsCheckType.Before) {
            beforeCheckpoints[beforeIndex] = logisticsCheckpoints[_trackingId][i];
            beforeIndex++;
        } else if (logisticsCheckpoints[_trackingId][i].checkType == LogisticsCheckType.During) {
            duringCheckpoints[duringIndex] = logisticsCheckpoints[_trackingId][i];
            duringIndex++;
        } else if (logisticsCheckpoints[_trackingId][i].checkType == LogisticsCheckType.After) {
            afterCheckpoints[afterIndex] = logisticsCheckpoints[_trackingId][i];
            afterIndex++;
        }
    }

    return (beforeCheckpoints, duringCheckpoints, afterCheckpoints);
}


}
