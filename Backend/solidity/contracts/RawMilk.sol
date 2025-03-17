// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "./UserRegistry.sol";  

contract RawMilk {
enum MilkStatus { Pending, Approved, Rejected, Used }

    struct MilkTank {
        bytes32 tankId;
        address farmer;
        bytes32 factoryId; // ✅ เปลี่ยนเป็น bytes32 เพื่อรองรับ FactoryID เป็น String
        string personInCharge;
        MilkStatus status;
        string qualityReportCID;
        string qrCodeCID;
    }

    struct MilkHistory {
        string personInCharge; // คนยืนยันการตรวจสอบ
        string qualityReportCID;
        MilkStatus status;
        uint256 timestamp; // เวลาที่อัปเดต
    }

    struct MilkTankWithHistory {
    bytes32 tankId;
    address farmer;
    bytes32 factoryId;
    string personInCharge;
    MilkStatus status;
    string qualityReportCID;
    string qrCodeCID;
    MilkHistory[] history; // ✅ คืนค่าประวัติเป็น Array
}


    UserRegistry public userRegistry;
    mapping(bytes32 => MilkTank) public milkTanks; 
    mapping(bytes32 => MilkHistory[]) public milkHistory; // ✅ ใช้ mapping เพื่อเก็บประวัติของแท็งก์นม
    bytes32[] public tankIds; 

    // ✅ Event สำหรับการสร้างแท็งก์นม
    event MilkTankCreated(
        bytes32 indexed tankId, 
        address indexed farmer, 
        bytes32 indexed factoryId, 
        string personInCharge, 
        string qualityReportCID,
        MilkStatus status // ✅ เพิ่ม status ตอนสร้าง
    );

    // ✅ Event สำหรับการอัปเดตแท็งก์นม
    event MilkTankUpdated(bytes32 indexed tankId);
    event MilkStatusUpdated(bytes32 indexed milkTankId, MilkStatus newStatus);

    // ✅ Event สำหรับการตรวจสอบคุณภาพนม
    event MilkQualityUpdated(
        bytes32 indexed tankId,
        string oldQualityReportCID,
        string newQualityReportCID,
        string oldPersonInCharge,
        string newPersonInCharge,
        MilkStatus status
    );

    // ✅ Event สำหรับการยืนยันคุณภาพนม
    event MilkQualityVerified(
        bytes32 indexed tankId,
        MilkStatus status,
        string qualityReportCID
    );

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

function markMilkTankAsUsed(bytes32 _milkTankId) public onlyFactory {
    require(milkTanks[_milkTankId].tankId != bytes32(0), "Milk tank does not exist");
    require(milkTanks[_milkTankId].status == MilkStatus.Approved, "Milk tank must be approved first");

    milkTanks[_milkTankId].status = MilkStatus.Used;

    emit MilkStatusUpdated(_milkTankId, MilkStatus.Used);
}

    // ✅ ฟังก์ชันสร้างแท็งก์นม
function createMilkTank(
    bytes32 _tankId,
    bytes32 _factoryId, 
    string memory _personInCharge,
    string memory _qualityReportCID,
    string memory _qrCodeCID
) public onlyFarmer {
    emit DebugLog("Creating Milk Tank", msg.sender, _tankId, _factoryId);

    require(milkTanks[_tankId].farmer == address(0), "Error: Tank ID already exists");

    // ✅ บันทึกข้อมูลแท็งก์นม
    milkTanks[_tankId] = MilkTank({
        tankId: _tankId,
        farmer: msg.sender,
        factoryId: _factoryId,
        personInCharge: _personInCharge,
        status: MilkStatus.Pending,
        qualityReportCID: _qualityReportCID,
        qrCodeCID: _qrCodeCID
    });

    // ✅ เก็บแท็งก์ในรายการ
    tankIds.push(_tankId);

    // ✅ บันทึกประวัติการสร้าง
    milkHistory[_tankId].push(MilkHistory({
        personInCharge: _personInCharge,
        qualityReportCID: _qualityReportCID,
        status: MilkStatus.Pending,
        timestamp: block.timestamp
    }));

    // ✅ Emit Event บันทึกการสร้าง
    emit MilkTankCreated(_tankId, msg.sender, _factoryId, _personInCharge, _qualityReportCID, MilkStatus.Pending);
}

// ✅ ฟังก์ชันตรวจสอบคุณภาพนม
function verifyMilkQuality(
    bytes32 _tankId, 
    bool _approved, 
    string memory _qualityReportCID,
    string memory _personInCharge
) public onlyFactory {
    require(milkTanks[_tankId].farmer != address(0), "Error: Tank ID does not exist");
    require(milkTanks[_tankId].status == MilkStatus.Pending, "Error: Milk already verified");

    // ✅ เก็บค่าก่อนอัปเดต
    string memory oldPersonInCharge = milkTanks[_tankId].personInCharge;
    string memory oldQualityReportCID = milkTanks[_tankId].qualityReportCID;

    // ✅ อัปเดตข้อมูลใหม่
    milkTanks[_tankId].status = _approved ? MilkStatus.Approved : MilkStatus.Rejected;
    milkTanks[_tankId].qualityReportCID = _qualityReportCID;
    milkTanks[_tankId].personInCharge = _personInCharge;

    // ✅ บันทึกประวัติการตรวจสอบ
    milkHistory[_tankId].push(MilkHistory({
        personInCharge: _personInCharge,
        qualityReportCID: _qualityReportCID,
        status: milkTanks[_tankId].status,
        timestamp: block.timestamp
    }));

    // ✅ Emit Event บันทึกการเปลี่ยนแปลง
    emit MilkQualityUpdated(
        _tankId, 
        oldQualityReportCID, 
        _qualityReportCID, 
        oldPersonInCharge, 
        _personInCharge, 
        milkTanks[_tankId].status
    );

    // ✅ Emit Event ยืนยันคุณภาพ
    emit MilkQualityVerified(_tankId, milkTanks[_tankId].status, _qualityReportCID);
}

    // ✅ ฟังก์ชันดึงข้อมูลแท็งก์นมดิบ (รวมประวัติ)
   function getMilkTank(bytes32 _tankId) public view returns (MilkTankWithHistory memory) {
    require(milkTanks[_tankId].farmer != address(0), "Error: Tank ID does not exist");

    return MilkTankWithHistory({
        tankId: milkTanks[_tankId].tankId,
        farmer: milkTanks[_tankId].farmer,
        factoryId: milkTanks[_tankId].factoryId,
        personInCharge: milkTanks[_tankId].personInCharge,
        status: milkTanks[_tankId].status,
        qualityReportCID: milkTanks[_tankId].qualityReportCID,
        qrCodeCID: milkTanks[_tankId].qrCodeCID,
        history: milkHistory[_tankId] // ✅ คืนค่าประวัติทั้งหมด
    });
}


    // ✅ ฟังก์ชันดึงรายการแท็งก์ตามฟาร์ม (รวมประวัติ)
    function getMilkTanksByFarmer(address _farmer) public view returns (bytes32[] memory, MilkHistory[][] memory) {
        uint count = 0;
        for (uint i = 0; i < tankIds.length; i++) {
            if (milkTanks[tankIds[i]].farmer == _farmer) {
                count++;
            }
        }

        bytes32[] memory farmerTanks = new bytes32[](count);
        MilkHistory[][] memory histories = new MilkHistory[][](count);
        uint index = 0;
        for (uint i = 0; i < tankIds.length; i++) {
            if (milkTanks[tankIds[i]].farmer == _farmer) {
                farmerTanks[index] = tankIds[i];
                histories[index] = milkHistory[tankIds[i]]; // ✅ คืนค่าประวัติ
                index++;
            }
        }
        return (farmerTanks, histories);
    }

    // ✅ ฟังก์ชันดึงรายการแท็งก์ตามโรงงาน (รวมประวัติ)
    function getMilkTanksByFactory(bytes32 _factoryId) public view returns (bytes32[] memory, MilkHistory[][] memory) {
        uint count = 0;
        for (uint i = 0; i < tankIds.length; i++) {
            if (milkTanks[tankIds[i]].factoryId == _factoryId) {
                count++;
            }
        }

        bytes32[] memory factoryTanks = new bytes32[](count);
        MilkHistory[][] memory histories = new MilkHistory[][](count);
        uint index = 0;
        for (uint i = 0; i < tankIds.length; i++) {
            if (milkTanks[tankIds[i]].factoryId == _factoryId) {
                factoryTanks[index] = tankIds[i];
                histories[index] = milkHistory[tankIds[i]]; // ✅ คืนค่าประวัติ
                index++;
            }
        }
        return (factoryTanks, histories);
    }

}