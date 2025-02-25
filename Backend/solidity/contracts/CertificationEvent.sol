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
        bool isActive; // เพิ่มฟิลด์สถานะ
    }

    mapping(string => CertEvent) public certificationEvents;
    event CertificationEventStored(
        string eventID,
        string entityType,
        string entityID,
        string certificationCID,
        uint256 issuedDate,
        uint256 expiryDate,
        uint256 createdOn,
        bool isActive
    );
    
    event CertificationEventDeactivated(string eventID, string entityID, uint256 timestamp);

    function storeCertificationEvent(
        string memory eventID,
        string memory entityType,
        string memory entityID,
        string memory certificationCID,
        uint256 issuedDate,
        uint256 expiryDate
    ) public {
        certificationEvents[eventID] = CertEvent(
            eventID,
            entityType,
            entityID,
            certificationCID,
            issuedDate,
            expiryDate,
            block.timestamp,
            true // ✅ ใบเซอร์เริ่มต้นจะเป็น active
        );

        emit CertificationEventStored(eventID, entityType, entityID, certificationCID, issuedDate, expiryDate, block.timestamp, true);
    }

    function getCertificationEvent(string memory eventID) public view returns (string memory, string memory, string memory, string memory, uint256, uint256, uint256, bool) {
        CertEvent memory certEvent = certificationEvents[eventID];
        return (certEvent.eventID, certEvent.entityType, certEvent.entityID, certEvent.certificationCID, certEvent.issuedDate, certEvent.expiryDate, certEvent.createdOn, certEvent.isActive);
    }

    function deactivateCertificationEvent(string memory eventID) public {
        require(bytes(certificationEvents[eventID].eventID).length != 0, "Certification event does not exist");
        certificationEvents[eventID].isActive = false; // ✅ เปลี่ยนสถานะเป็นไม่ใช้งาน
        emit CertificationEventDeactivated(eventID, certificationEvents[eventID].entityID, block.timestamp);
    }
}
