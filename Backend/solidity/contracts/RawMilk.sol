// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./UserRegistry.sol";  

contract RawMilk {
    enum MilkStatus { Pending, Approved, Rejected }

    struct MilkTank {
        bytes32 tankId;
        address farmer;
        bytes32 factoryId; // ✅ เปลี่ยนเป็น bytes32 เพื่อรองรับ FactoryID เป็น String
        string personInCharge;
        MilkStatus status;
        string qualityReportCID;
        string qrCodeCID;
    }

    UserRegistry public userRegistry;
    mapping(bytes32 => MilkTank) public milkTanks; 
    bytes32[] public tankIds; 

    event MilkTankCreated(
        bytes32 indexed tankId, 
        address indexed farmer, 
        bytes32 indexed factoryId, // ✅ เปลี่ยน factoryId เป็น bytes32
        string personInCharge, 
        string qualityReportCID
    );

    event MilkTankUpdated(bytes32 indexed tankId);
    event MilkQualityVerified(bytes32 indexed tankId, MilkStatus status, string qualityReportCID);

    event DebugLog(string message, address sender, bytes32 tankId, bytes32 factoryId);

modifier onlyFarmer() {
    emit DebugLog("Checking User Role", msg.sender, 0x0, 0x0);
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

    function createMilkTank(
    bytes32 _tankId,
    bytes32 _factoryId, // ✅ เก็บเป็น bytes32 ตามเดิม
    string memory _personInCharge,
    string memory _qualityReportCID,
    string memory _qrCodeCID
) public onlyFarmer {
    emit DebugLog("Creating Milk Tank", msg.sender, _tankId, _factoryId);

    require(milkTanks[_tankId].farmer == address(0), "Error: Tank ID already exists");

    milkTanks[_tankId] = MilkTank({
        tankId: _tankId,
        farmer: msg.sender,
        factoryId: _factoryId, // ✅ ใช้ Factory ID เป็น `bytes32`
        personInCharge: _personInCharge,
        status: MilkStatus.Pending,
        qualityReportCID: _qualityReportCID,
        qrCodeCID: _qrCodeCID
    });

    tankIds.push(_tankId);
    emit MilkTankCreated(_tankId, msg.sender, _factoryId, _personInCharge, _qualityReportCID);
}


    // ✅ ฟังก์ชันอัปเดต QR Code หรือ Quality Report CID
    function updateMilkTank(
        bytes32 _tankId,
        string memory _qrCodeCID
    ) public onlyFarmer {
        require(milkTanks[_tankId].farmer == msg.sender, "Error: Unauthorized");
        require(milkTanks[_tankId].status == MilkStatus.Pending, "Error: Cannot update approved/rejected milk");

        milkTanks[_tankId].qrCodeCID = _qrCodeCID;

        emit MilkTankUpdated(_tankId);
    }

    // ✅ ฟังก์ชันตรวจสอบคุณภาพนม
    function verifyMilkQuality(
        bytes32 _tankId, 
        bool _approved, 
        string memory _qualityReportCID
    ) public onlyFactory {
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
        bytes32,
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
            tank.factoryId, // ✅ คืนค่า FactoryID (bytes32)
            tank.personInCharge,
            tank.status,
            tank.qualityReportCID,
            tank.qrCodeCID
        );
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

    // ✅ ฟังก์ชันดึงรายการแท็งก์ตามโรงงาน (ใช้ `bytes32` แทน `address`)
    function getMilkTanksByFactory(bytes32 _factoryId) public view returns (bytes32[] memory) {
        uint count = 0;
        for (uint i = 0; i < tankIds.length; i++) {
            if (milkTanks[tankIds[i]].factoryId == _factoryId) {
                count++;
            }
        }

        bytes32[] memory factoryTanks = new bytes32[](count);
        uint index = 0;
        for (uint i = 0; i < tankIds.length; i++) {
            if (milkTanks[tankIds[i]].factoryId == _factoryId) {
                factoryTanks[index] = tankIds[i];
                index++;
            }
        }
        return factoryTanks;
    }
}
