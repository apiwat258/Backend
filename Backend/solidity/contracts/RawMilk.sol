// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract RawMilkSupplyChain {
    address public owner;

    enum MilkStatus { Pending, Approved, Rejected }

    struct RawMilk {
        bytes32 id;
        address farmWallet;
        uint256 temperature;
        uint256 pH;
        uint256 fat;
        uint256 protein; 
        string ipfsCid;
        MilkStatus status;
        uint256 timestamp;
    }

    mapping(bytes32 => RawMilk) public rawMilkRecords;
    mapping(address => bool) public farms;

    event RawMilkAdded(bytes32 indexed rawMilkID, address indexed farmWallet, MilkStatus status, string ipfsCid);
    event RawMilkStatusUpdated(bytes32 indexed rawMilkID, MilkStatus newStatus);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only contract owner can call this");
        _;
    }

    modifier onlyFarm() {
        require(farms[msg.sender], "Only registered farms can add raw milk");
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    event FarmRegistered(address indexed farmWallet);

function registerFarm(address _farmWallet) external onlyOwner {
    farms[_farmWallet] = true;
    emit FarmRegistered(_farmWallet);  // ✅ เพิ่ม Event
}


function addRawMilk(
    bytes32 _rawMilkID, // 🆕 ให้เกษตรกรส่งค่า ID มาเอง
    uint256 _temperature,
    uint256 _pH,
    uint256 _fat,
    uint256 _protein,
    string memory _ipfsCid
) external onlyFarm {
    require(_temperature >= 200 && _temperature <= 600, "Temperature out of range!");
    require(_pH >= 650 && _pH <= 680, "pH out of range!");
    require(_fat >= 300 && _fat <= 400, "Fat percentage out of range!");
    require(_protein >= 300 && _protein <= 350, "Protein percentage out of range!");

    require(rawMilkRecords[_rawMilkID].farmWallet == address(0), "Milk ID already exists!"); // 🛑 ห้ามใช้ ID ซ้ำ

    rawMilkRecords[_rawMilkID] = RawMilk({
        id: _rawMilkID,
        farmWallet: msg.sender,
        temperature: _temperature,
        pH: _pH,
        fat: _fat,
        protein: _protein,
        ipfsCid: _ipfsCid,
        status: MilkStatus.Pending,
        timestamp: block.timestamp
    });

    emit RawMilkAdded(_rawMilkID, msg.sender, MilkStatus.Pending, _ipfsCid);
}

    function getRawMilk(bytes32 _rawMilkID) external view returns (
        address farmWallet,
        uint256 temperature,
        uint256 pH,
        uint256 fat,
        uint256 protein,
        string memory ipfsCid,
        MilkStatus status,
        uint256 timestamp
    ) {
        RawMilk memory milk = rawMilkRecords[_rawMilkID];
        require(milk.farmWallet != address(0), "Milk record not found");

        return (
            milk.farmWallet,
            milk.temperature,
            milk.pH,
            milk.fat,
            milk.protein,
            milk.ipfsCid,
            milk.status,
            milk.timestamp
        );
    }

    function updateRawMilkStatus(bytes32 _rawMilkID, MilkStatus _newStatus) external onlyOwner {
        require(rawMilkRecords[_rawMilkID].farmWallet != address(0), "Milk record not found");

        rawMilkRecords[_rawMilkID].status = _newStatus;

        emit RawMilkStatusUpdated(_rawMilkID, _newStatus);
    }
}
