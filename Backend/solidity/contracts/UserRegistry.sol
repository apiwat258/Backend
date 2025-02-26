// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract UserRegistry {
    enum UserRole { None, Farmer, Factory, Logistics, Retailer } // ✅ กำหนดประเภทของผู้ใช้

    struct User {
        address wallet;
        UserRole role;
        bool isRegistered;
    }

    mapping(address => User) public users;

    event UserRegistered(address indexed wallet, UserRole role);

    // ✅ ฟังก์ชันลงทะเบียนผู้ใช้
    function registerUser(UserRole role) public {
        require(users[msg.sender].isRegistered == false, "Error: User already registered");
        require(role != UserRole.None, "Error: Invalid role");

        users[msg.sender] = User({
            wallet: msg.sender,
            role: role,
            isRegistered: true
        });

        emit UserRegistered(msg.sender, role);
    }

    // ✅ ตรวจสอบว่าผู้ใช้ลงทะเบียนแล้วหรือไม่
    function isUserRegistered(address wallet) public view returns (bool) {
        return users[wallet].isRegistered;
    }

    // ✅ ตรวจสอบสิทธิ์ของผู้ใช้
    function getUserRole(address wallet) public view returns (UserRole) {
        return users[wallet].role;
    }
}
