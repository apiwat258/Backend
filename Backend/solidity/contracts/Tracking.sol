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
        address retailer;
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
        address retailer;
        uint256 receivedTime;
        string qualityCID;
        string personInCharge;
    }

    mapping(bytes32 => TrackingEvent) public trackingEvents;
    mapping(bytes32 => LogisticsCheckpoint[]) public logisticsCheckpoints;
    mapping(bytes32 => RetailerConfirmation) public retailerConfirmations;
    mapping(address => bytes32[]) public retailerTracking;
    bytes32[] public trackingIds;

    event TrackingCreated(bytes32 indexed trackingId, bytes32 indexed productLotId, address indexed retailer, string qrCodeCID);
    event LogisticsUpdated(bytes32 indexed trackingId, address indexed logisticsProvider, LogisticsCheckType checkType);
    event RetailerReceived(bytes32 indexed trackingId, address indexed retailer, string qualityCID);

    modifier onlyLogistics() {
        require(userRegistry.getUserRole(msg.sender) == UserRegistry.UserRole.Logistics, "Access denied: Only logistics allowed");
        _;
    }

    modifier onlyRetailer() {
        require(userRegistry.getUserRole(msg.sender) == UserRegistry.UserRole.Retailer, "Access denied: Only retailers allowed");
        _;
    }

    constructor(address _userRegistry, address _productLotContract) {
        userRegistry = UserRegistry(_userRegistry);
        productLotContract = ProductLot(_productLotContract);
    }

    function createTrackingEvent(
        bytes32 _trackingId,
        bytes32 _productLotId,
        address _retailer,
        string memory _qrCodeCID
    ) public {
        require(productLotContract.isProductLotExists(_productLotId), "Error: Product Lot does not exist");
        require(trackingEvents[_trackingId].trackingId == bytes32(0), "Error: Tracking ID already exists");

        trackingEvents[_trackingId] = TrackingEvent({
            trackingId: _trackingId,
            productLotId: _productLotId,
            retailer: _retailer,
            qrCodeCID: _qrCodeCID,
            status: TrackingStatus.Pending
        });

        trackingIds.push(_trackingId);
        retailerTracking[_retailer].push(_trackingId);
        emit TrackingCreated(_trackingId, _productLotId, _retailer, _qrCodeCID);
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
        string memory _qualityCID,
        string memory _personInCharge
    ) public onlyRetailer {
        require(trackingEvents[_trackingId].trackingId != bytes32(0), "Error: Tracking ID does not exist");
        require(trackingEvents[_trackingId].retailer == msg.sender, "Error: Only assigned retailer can receive");

        trackingEvents[_trackingId].status = TrackingStatus.Received;
        retailerConfirmations[_trackingId] = RetailerConfirmation({
            trackingId: _trackingId,
            retailer: msg.sender,
            receivedTime: block.timestamp,
            qualityCID: _qualityCID,
            personInCharge: _personInCharge
        });
        emit RetailerReceived(_trackingId, msg.sender, _qualityCID);
    }

    function getTrackingById(bytes32 _trackingId) public view returns (TrackingEvent memory, LogisticsCheckpoint[] memory, RetailerConfirmation memory) {
        require(trackingEvents[_trackingId].trackingId != bytes32(0), "Error: Tracking ID does not exist");
        return (trackingEvents[_trackingId], logisticsCheckpoints[_trackingId], retailerConfirmations[_trackingId]);
    }

    function getTrackingByRetailer(address _retailer) public view returns (bytes32[] memory) {
        return retailerTracking[_retailer];
    }
}
