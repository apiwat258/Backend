// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./UserRegistry.sol";  

contract RawMilk {
    enum MilkStatus { Pending, Approved, Rejected }

    struct MilkTank {
        bytes32 tankId;
        address farmer;
        string personInCharge;
        MilkStatus status;
        string qualityReportCID; // ✅ รวม bacteriaInfo, contaminantInfo, abnormalType
        string qrCodeCID;
    }

    UserRegistry public userRegistry;
    mapping(bytes32 => MilkTank) public milkTanks; 
    bytes32[] public tankIds; 

    event MilkTankCreated(bytes32 indexed tankId, address indexed farmer);
    event MilkTankUpdated(bytes32 indexed tankId);
    event MilkQualityVerified(bytes32 indexed tankId, MilkStatus status, string qualityReportCID);

    modifier onlyFarmer() {
        require(userRegistry.getUserRole(msg.sender) == UserRegistry.UserRole.Farmer, "Access denied: Only farmers allowed");
        _;
    }

    modifier onlyFactory() {
        require(userRegistry.getUserRole(msg.sender) == UserRegistry.UserRole.Factory, "Access denied: Only factories allowed");
        _;
    }

    constructor(address _userRegistry) {
        userRegistry = UserRegistry(_userRegistry);
    }

    // ✅ ฟังก์ชันตรวจสอบค่าก่อนบันทึก (ไม่เก็บบนบล็อกเชน)
    function validateMilkData(
        uint256 _quantity,
        uint256 _temperature,
        uint256 _pH,
        uint256 _fat,
        uint256 _protein,
        bool _bacteria,
        bool _contaminants
    ) public pure returns (bool valid, string memory message) {
        if (_temperature < 200 || _temperature > 600) return (false, "Error: Temperature out of range! (2.0C - 6.0C)");
        if (_pH < 650 || _pH > 680) return (false, "Error: pH out of range! (6.5 - 6.8)");
        if (_fat < 300 || _fat > 400) return (false, "Error: Fat percentage out of range! (3.0% - 4.0%)");
        if (_protein < 300 || _protein > 350) return (false, "Error: Protein percentage out of range! (3.0% - 3.5%)");

        return (true, "Validated successfully.");
    }

    // ✅ ฟังก์ชันสร้างแท็งก์ (เฉพาะค่าที่ต้องเก็บบนบล็อกเชน)
    function createMilkTank(
        bytes32 _tankId,
        string memory _personInCharge,
        string memory _qrCodeCID
    ) public onlyFarmer {
        require(milkTanks[_tankId].farmer == address(0), "Error: Tank ID already exists");

        milkTanks[_tankId] = MilkTank({
            tankId: _tankId,
            farmer: msg.sender,
            personInCharge: _personInCharge,
            status: MilkStatus.Pending,
            qualityReportCID: "",
            qrCodeCID: _qrCodeCID
        });

        tankIds.push(_tankId);
        emit MilkTankCreated(_tankId, msg.sender);
    }

    // ✅ ฟังก์ชันอัปเดต QR Code หรือ Quality Report CID (ห้ามแก้ข้อมูลอื่น)
    function updateMilkTank(
        bytes32 _tankId,
        string memory _qrCodeCID
    ) public onlyFarmer {
        require(milkTanks[_tankId].farmer == msg.sender, "Error: Unauthorized");
        require(milkTanks[_tankId].status == MilkStatus.Pending, "Error: Cannot update approved/rejected milk");

        milkTanks[_tankId].qrCodeCID = _qrCodeCID;

        emit MilkTankUpdated(_tankId);
    }

    // ✅ ฟังก์ชันตรวจสอบคุณภาพนม (โรงงานเป็นผู้อนุมัติ)
    function verifyMilkQuality(bytes32 _tankId, bool _approved, string memory _qualityReportCID) public onlyFactory {
        require(milkTanks[_tankId].farmer != address(0), "Error: Tank ID does not exist");
        require(milkTanks[_tankId].status == MilkStatus.Pending, "Error: Milk already verified");

        milkTanks[_tankId].status = _approved ? MilkStatus.Approved : MilkStatus.Rejected;
        milkTanks[_tankId].qualityReportCID = _qualityReportCID;

        emit MilkQualityVerified(_tankId, milkTanks[_tankId].status, _qualityReportCID);
    }

    // ✅ ฟังก์ชันดึงข้อมูลแท็งก์นมดิบ
    function getMilkTank(bytes32 _tankId) public view returns (
        bytes32,
        address,
        string memory,
        MilkStatus,
        string memory,
        string memory
    ) {
        MilkTank memory tank = milkTanks[_tankId];
        require(tank.farmer != address(0), "Error: Tank ID does not exist");

        return (
            tank.tankId,
            tank.farmer,
            tank.personInCharge,
            tank.status,
            tank.qualityReportCID,
            tank.qrCodeCID
        );
    }

    // ✅ ฟังก์ชันดึงรายการแท็งก์ทั้งหมด
    function getAllMilkTanks() public view returns (bytes32[] memory) {
        return tankIds;
    }

    // ✅ ฟังก์ชันดึงรายการแท็งก์ตามฟาร์ม
    function getMilkTanksByFarmer(address _farmer) public view returns (bytes32[] memory) {
        uint count = 0;
        for (uint i = 0; i < tankIds.length; i++) {
            if (milkTanks[tankIds[i]].farmer == _farmer) {
                count++;
            }
        }

        bytes32[] memory farmerTanks = new bytes32[](count);
        uint index = 0;
        for (uint i = 0; i < tankIds.length; i++) {
            if (milkTanks[tankIds[i]].farmer == _farmer) {
                farmerTanks[index] = tankIds[i];
                index++;
            }
        }
        return farmerTanks;
    }
}
