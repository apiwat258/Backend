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
    bytes32[] public trackingIds;

    event TrackingCreated(bytes32 indexed trackingId, bytes32 indexed productLotId, string retailerId, string qrCodeCID);
    event LogisticsUpdated(bytes32 indexed trackingId, address indexed logisticsProvider, LogisticsCheckType checkType);
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

        trackingIds.push(_trackingId);
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
        LogisticsCheckType _checkType
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
            checkType: _checkType
        }));
        emit LogisticsUpdated(_trackingId, msg.sender, _checkType);
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
        emit RetailerReceived(_trackingId, _retailerId, _qualityCID);
    }

    function getTrackingById(bytes32 _trackingId) public view returns (TrackingEvent memory, LogisticsCheckpoint[] memory, RetailerConfirmation memory) {
        require(trackingEvents[_trackingId].trackingId != bytes32(0), "Error: Tracking ID does not exist");
        return (trackingEvents[_trackingId], logisticsCheckpoints[_trackingId], retailerConfirmations[_trackingId]);
    }

    function getTrackingByRetailer(string memory _retailerId) public view returns (bytes32[] memory) {
        return retailerTracking[_retailerId];
    }
}
