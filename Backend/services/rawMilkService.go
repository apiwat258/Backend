package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"finalyearproject/Backend/services/rawmilk"
)

// ✅ Struct ใช้เก็บข้อมูล Raw Milk
type RawMilkData struct {
	FarmWallet  string  `json:"farmWallet"`
	Temperature float64 `json:"temperature"`
	PH          float64 `json:"pH"`
	Fat         float64 `json:"fat"`
	Protein     float64 `json:"protein"`
	IPFSCid     string  `json:"ipfsCid"`
	Status      uint8   `json:"status"`
	Timestamp   int64   `json:"timestamp"`
}

// RawMilkService - สำหรับจัดการข้อมูลน้ำนมดิบบน Blockchain
type RawMilkService struct {
	client          *ethclient.Client
	auth            *bind.TransactOpts
	rawMilkContract *rawmilk.Rawmilk
}

// ✅ Global Instance ของ Raw Milk Service
var RawMilkServiceInstance *RawMilkService

// InitRawMilkService - ใช้เชื่อมต่อกับ Blockchain และโหลด Smart Contract
func InitRawMilkService(client *ethclient.Client, auth *bind.TransactOpts, contractAddr common.Address) error {
	rawMilkInstance, err := rawmilk.NewRawmilk(contractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load raw milk contract: %v", err)
	}

	RawMilkServiceInstance = &RawMilkService{
		client:          client,
		auth:            auth,
		rawMilkContract: rawMilkInstance,
	}

	fmt.Println("✅ Raw Milk Service Initialized!")
	return nil
}

// StoreRawMilkOnBlockchain - บันทึกข้อมูลน้ำนมดิบลง Blockchain
func (rms *RawMilkService) StoreRawMilkOnBlockchain(
	rawMilkHash [32]byte, // ✅ ใช้ bytes32
	farmWallet string,
	temperature, pH, fat, protein float64,
	ipfsCid string,
) (string, error) {
	tempBigInt := big.NewInt(int64(temperature * 100))
	pHBigInt := big.NewInt(int64(pH * 100))
	fatBigInt := big.NewInt(int64(fat * 100))
	proteinBigInt := big.NewInt(int64(protein * 100))

	tx, err := rms.rawMilkContract.AddRawMilk(
		rms.auth,
		rawMilkHash, // ✅ ใช้ bytes32 ตาม Smart Contract
		tempBigInt,
		pHBigInt,
		fatBigInt,
		proteinBigInt,
		ipfsCid,
	)
	if err != nil {
		log.Println("❌ Failed to store raw milk on blockchain:", err)
		return "", err
	}

	fmt.Println("✅ Raw Milk stored on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// GetRawMilkFromBlockchain - ดึงข้อมูลน้ำนมดิบจาก Blockchain
func (rms *RawMilkService) GetRawMilkFromBlockchain(rawMilkID common.Hash) (*RawMilkData, error) {
	milk, err := rms.rawMilkContract.GetRawMilk(&bind.CallOpts{}, rawMilkID)
	if err != nil {
		log.Println("❌ Failed to fetch raw milk data from blockchain:", err)
		return nil, err
	}

	// ✅ แปลงค่าจาก BigInt → float64 และ uint8
	rawMilk := &RawMilkData{
		FarmWallet:  milk.FarmWallet.Hex(),
		Temperature: float64(milk.Temperature.Int64()) / 100,
		PH:          float64(milk.PH.Int64()) / 100,
		Fat:         float64(milk.Fat.Int64()) / 100,
		Protein:     float64(milk.Protein.Int64()) / 100,
		IPFSCid:     milk.IpfsCid,
		Status:      uint8(milk.Status), // ✅ ใช้ uint8 ตรง ๆ
		Timestamp:   milk.Timestamp.Int64(),
	}

	return rawMilk, nil
}
