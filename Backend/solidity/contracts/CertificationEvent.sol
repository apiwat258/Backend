// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./UserRegistry.sol"; // ✅ เชื่อมกับระบบลงทะเบียนกลาง

contract CertificationEvent {
    UserRegistry private userRegistry; // ✅ อ้างถึง UserRegistry

    struct CertEvent {
        string eventID;
        string entityType; 
        string entityID;
        string certificationCID;
        uint256 issuedDate;
        uint256 expiryDate;
        uint256 createdOn;
        bool isActive;
    }

    mapping(string => CertEvent) public certificationEvents;
    mapping(string => string[]) private entityCertifications; // ✅ เก็บ eventID ทั้งหมดของ entityID
    mapping(string => bool) private existingCIDs; // ✅ ป้องกัน CID ซ้ำในระบบ
    string[] private allEventIDs; // ✅ เก็บ eventID ทั้งหมด (ใช้ใน getAllCertifications)

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

    // ✅ Constructor เชื่อมกับ UserRegistry
    constructor(address _userRegistryAddress) {
        userRegistry = UserRegistry(_userRegistryAddress);
    }

    // ✅ Modifier เช็คว่า msg.sender ต้องเป็นผู้ใช้ที่ลงทะเบียนแล้ว
    modifier onlyRegisteredUser() {
        require(userRegistry.isUserRegistered(msg.sender), "Error: User must register first");
        _;
    }

    // ✅ บันทึกใบเซอร์ใหม่ → ทุก Role สามารถใช้ได้
    function storeCertificationEvent(
        string memory eventID,
        string memory entityType,
        string memory entityID,
        string memory certificationCID,
        uint256 issuedDate,
        uint256 expiryDate
    ) public onlyRegisteredUser {
        require(bytes(certificationEvents[eventID].eventID).length == 0, "Error: Certification already exists");

        // ✅ ตรวจสอบว่า CID นี้เคยถูกใช้แล้วหรือไม่ (Global Unique Check)
        require(existingCIDs[certificationCID] == false, "Error: Certification CID already exists in system");

        // ✅ บันทึกข้อมูลใบเซอร์ใหม่
        certificationEvents[eventID] = CertEvent(
            eventID,
            entityType,
            entityID,
            certificationCID,
            issuedDate,
            expiryDate,
            block.timestamp,
            true
        );

        entityCertifications[entityID].push(eventID); // ✅ เพิ่ม eventID เข้าไปใน entityID
        existingCIDs[certificationCID] = true; // ✅ บันทึกว่า CID นี้ถูกใช้แล้ว
        allEventIDs.push(eventID); // ✅ บันทึก eventID ไว้เพื่อใช้ดึงข้อมูลทั้งหมด

        emit CertificationEventStored(eventID, entityType, entityID, certificationCID, issuedDate, expiryDate, block.timestamp, true);
    }

    // ✅ ปิดใช้งานใบเซอร์
    function deactivateCertificationEvent(string memory eventID) public onlyRegisteredUser {
        require(bytes(certificationEvents[eventID].eventID).length != 0, "Error: Certification event does not exist");
        require(certificationEvents[eventID].isActive == true, "Error: Certification is already inactive");

        certificationEvents[eventID].isActive = false;
        emit CertificationEventDeactivated(eventID, certificationEvents[eventID].entityID, block.timestamp);
    }

    // ✅ ดึง "ใบเซอร์ที่ Active อยู่เท่านั้น" ของ entityID
    function getActiveCertificationsForEntity(string memory entityID) public view returns (CertEvent[] memory) {
        string[] memory certs = entityCertifications[entityID];
        uint256 activeCount = 0;

        // ✅ นับจำนวนใบเซอร์ที่ Active
        for (uint256 i = 0; i < certs.length; i++) {
            if (certificationEvents[certs[i]].isActive) {
                activeCount++;
            }
        }

        // ✅ สร้างอาร์เรย์เก็บใบเซอร์ที่ Active
        CertEvent[] memory activeCerts = new CertEvent[](activeCount);
        uint256 index = 0;
        for (uint256 i = 0; i < certs.length; i++) {
            if (certificationEvents[certs[i]].isActive) {
                activeCerts[index] = certificationEvents[certs[i]];
                index++;
            }
        }

        return activeCerts;
    }

    // ✅ ดึง "ใบเซอร์ทั้งหมดในระบบ" (ทุก entity)
    function getAllCertifications() public view returns (CertEvent[] memory) {
        uint256 totalCerts = allEventIDs.length;

        // ✅ สร้างอาร์เรย์สำหรับเก็บทุกใบเซอร์
        CertEvent[] memory allCerts = new CertEvent[](totalCerts);
        
        for (uint256 i = 0; i < totalCerts; i++) {
            allCerts[i] = certificationEvents[allEventIDs[i]];
        }

        return allCerts;
    }

    // ✅ ตรวจสอบว่า CID นี้มีอยู่ในระบบหรือไม่
    function isCertificationCIDExists(string memory certificationCID) public view returns (bool) {
        return existingCIDs[certificationCID];
    }
}
