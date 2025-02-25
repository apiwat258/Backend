package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"finalyearproject/Backend/services/certification_event" // ✅ สำหรับ Certification Event
	"finalyearproject/Backend/services/rawmilk"             // ✅ สำหรับ Raw Milk
)

// ✅ เพิ่ม struct นี้ก่อนฟังก์ชัน
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

// BlockchainService - ใช้สำหรับเชื่อมต่อ Blockchain
type BlockchainService struct {
	client                *ethclient.Client
	auth                  *bind.TransactOpts
	certificationContract *certification_event.CertificationEvent
	rawMilkContract       *rawmilk.Rawmilk // ✅ ใช้ struct ที่ถูกต้อง// ✅ ใช้ Smart Contract ของ Raw Milk
}

// BlockchainServiceInstance - Global Instance
var BlockchainServiceInstance *BlockchainService

// InitBlockchainService - เชื่อมต่อ Blockchain
func InitBlockchainService() error {
	// ✅ Debug: ตรวจสอบค่า ENV ที่โหลดเข้ามา
	fmt.Println("📌 DEBUG - BLOCKCHAIN_RPC_URL:", os.Getenv("BLOCKCHAIN_RPC_URL"))
	fmt.Println("📌 DEBUG - PRIVATE_KEY:", os.Getenv("PRIVATE_KEY"))
	fmt.Println("📌 DEBUG - CERT_CONTRACT_ADDRESS:", os.Getenv("CERT_CONTRACT_ADDRESS"))
	fmt.Println("📌 DEBUG - RAWMILK_CONTRACT_ADDRESS:", os.Getenv("RAWMILK_CONTRACT_ADDRESS"))

	// ✅ โหลดค่า RPC URL จาก ENV
	rpcURL := os.Getenv("BLOCKCHAIN_RPC_URL")
	if rpcURL == "" {
		return fmt.Errorf("❌ BLOCKCHAIN_RPC_URL is not set")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return fmt.Errorf("❌ Failed to connect to blockchain: %v", err)
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	certContractAddress := os.Getenv("CERT_CONTRACT_ADDRESS")
	rawMilkContractAddress := os.Getenv("RAWMILK_CONTRACT_ADDRESS")

	if privateKeyHex == "" || certContractAddress == "" || rawMilkContractAddress == "" {
		return fmt.Errorf("❌ Missing blockchain env variables")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("❌ Invalid private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337)) // ✅ ใช้ Chain ID 1337 (Ganache)
	if err != nil {
		return fmt.Errorf("❌ Failed to create transaction auth: %v", err)
	}

	certContractAddr := common.HexToAddress(certContractAddress)
	rawMilkContractAddr := common.HexToAddress(rawMilkContractAddress)

	// ✅ โหลด Certification Smart Contract
	certInstance, err := certification_event.NewCertificationEvent(certContractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load certification contract: %v", err)
	}

	// ✅ โหลด RawMilk Smart Contract
	rawMilkInstance, err := rawmilk.NewRawmilk(rawMilkContractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load raw milk contract: %v", err)
	}

	BlockchainServiceInstance = &BlockchainService{
		client:                client,
		auth:                  auth,
		certificationContract: certInstance,
		rawMilkContract:       rawMilkInstance,
	}

	fmt.Println("✅ Blockchain Service Initialized!")
	return nil
}

// RegisterFarmOnBlockchain - ลงทะเบียนฟาร์มบน Blockchain
func (b *BlockchainService) RegisterFarmOnBlockchain(farmWallet string) (string, error) {
	fmt.Println("📌 DEBUG - Registering Farm on Blockchain:", farmWallet)

	farmAddress := common.HexToAddress(farmWallet)

	tx, err := b.rawMilkContract.RegisterFarm(b.auth, farmAddress)
	if err != nil {
		log.Println("❌ Failed to register farm on blockchain:", err)
		return "", err
	}

	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		log.Println("❌ Transaction not mined:", err)
		return "", err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		log.Println("❌ Transaction failed!")
		return "", errors.New("Transaction failed")
	}

	fmt.Println("✅ Farm registered on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// StoreCertificationOnBlockchain - ฟังก์ชันบันทึกใบเซอร์ลง Blockchain
func (b *BlockchainService) StoreCertificationOnBlockchain(eventID, entityType, entityID, certCID string, issuedDate, expiryDate *big.Int) (string, error) {
	tx, err := b.certificationContract.StoreCertificationEvent(b.auth, eventID, entityType, entityID, certCID, issuedDate, expiryDate)
	if err != nil {
		log.Println("❌ Failed to store certification event on blockchain:", err)
		return "", err
	}

	// รอให้ธุรกรรมถูก Mine
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		log.Println("❌ Transaction not mined:", err)
		return "", err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		log.Println("❌ Transaction failed!")
		return "", errors.New("transaction failed")
	}

	fmt.Println("✅ Certification Event stored on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// DeactivateCertificationOnBlockchain - ปิดใช้งานใบเซอร์บน Blockchain
func (b *BlockchainService) DeactivateCertificationOnBlockchain(eventID string) (string, error) {
	tx, err := b.certificationContract.DeactivateCertificationEvent(b.auth, eventID)
	if err != nil {
		log.Println("❌ Failed to deactivate certification event on blockchain:", err)
		return "", err
	}

	fmt.Println("✅ Certification Event deactivated on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// GetCertificationFromBlockchain - ดึงข้อมูลใบเซอร์จาก Blockchain
func (b *BlockchainService) GetCertificationFromBlockchain(eventID string) (*models.Certification, error) {
	certEvent, err := b.certificationContract.GetCertificationEvent(&bind.CallOpts{}, eventID)
	if err != nil {
		log.Println("❌ [Blockchain] Failed to fetch certification:", err)
		return nil, err
	}

	issuedDate := time.Unix(certEvent.IssuedDate.Int64(), 0)
	expiryDate := time.Unix(certEvent.ExpiryDate.Int64(), 0)

	return &models.Certification{
		CertificationID:   certEvent.EventID,
		EntityType:        certEvent.EntityType,
		EntityID:          certEvent.EntityID,
		CertificationCID:  certEvent.CertificationCID,
		IssuedDate:        issuedDate,
		EffectiveDate:     expiryDate,
		BlockchainTxHash:  "",
		CreatedOn:         time.Unix(certEvent.CreatedOn.Int64(), 0),
		IsActive:          certEvent.IsActive, // ✅ เพิ่มฟิลด์ isActive
	}, nil
}




// StoreRawMilkOnBlockchain - บันทึกข้อมูลน้ำนมดิบลง Blockchain
func (b *BlockchainService) StoreRawMilkOnBlockchain(
	rawMilkHash [32]byte, // ✅ ใช้ bytes32
	farmWallet string,
	temperature, pH, fat, protein float64,
	ipfsCid string,
) (string, error) {
	tempBigInt := big.NewInt(int64(temperature * 100))
	pHBigInt := big.NewInt(int64(pH * 100))
	fatBigInt := big.NewInt(int64(fat * 100))
	proteinBigInt := big.NewInt(int64(protein * 100))

	tx, err := b.rawMilkContract.AddRawMilk(
		b.auth,
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
func (b *BlockchainService) GetRawMilkFromBlockchain(rawMilkID common.Hash) (*RawMilkData, error) {
	milk, err := b.rawMilkContract.GetRawMilk(&bind.CallOpts{}, rawMilkID)
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

// UpdateRawMilkStatus - อัปเดตสถานะน้ำนมดิบ
func (b *BlockchainService) UpdateRawMilkStatus(rawMilkID string, newStatus uint8) (string, error) {
	rawMilkBytes := common.HexToHash(rawMilkID)

	tx, err := b.rawMilkContract.UpdateRawMilkStatus(b.auth, rawMilkBytes, newStatus) // ✅ ใช้ uint8 ตรง ๆ
	if err != nil {
		log.Println("❌ Failed to update raw milk status on blockchain:", err)
		return "", err
	}

	fmt.Println("✅ Raw Milk status updated on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}
