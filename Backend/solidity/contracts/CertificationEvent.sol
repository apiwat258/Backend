// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract CertificationEvent {
    struct CertEvent {
        string eventID;
        string entityType; // Farmer, Factory, Retailer, Logistics
        string entityID;
        string certificationCID;
        uint256 issuedDate;
        uint256 expiryDate;
        uint256 createdOn;
    }

    mapping(string => CertEvent) public certificationEvents;
    event CertificationEventStored(string eventID, string entityType, string entityID, string certificationCID, uint256 issuedDate, uint256 expiryDate, uint256 createdOn);

    function storeCertificationEvent(
        string memory eventID,
        string memory entityType,
        string memory entityID,
        string memory certificationCID,
        uint256 issuedDate,
        uint256 expiryDate
    ) public {
        certificationEvents[eventID] = CertEvent(eventID, entityType, entityID, certificationCID, issuedDate, expiryDate, block.timestamp);
        emit CertificationEventStored(eventID, entityType, entityID, certificationCID, issuedDate, expiryDate, block.timestamp);
    }

    function getCertificationEvent(string memory eventID) public view returns (string memory, string memory, string memory, string memory, uint256, uint256, uint256) {
        CertEvent memory certEvent = certificationEvents[eventID];
        return (certEvent.eventID, certEvent.entityType, certEvent.entityID, certEvent.certificationCID, certEvent.issuedDate, certEvent.expiryDate, certEvent.createdOn);
    }
}
