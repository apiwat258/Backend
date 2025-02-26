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

	certification "finalyearproject/Backend/services/certification_event" // ✅ สำหรับ Certification Event
	"finalyearproject/Backend/services/rawmilk"                           // ✅ สำหรับ Raw Milk
	"finalyearproject/Backend/services/userregistry"
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
	userRegistryContract  *userregistry.Userregistry
	certificationContract *certification.Certification
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
	certInstance, err := certification.NewCertification(certContractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load certification contract: %v", err)
	}

	// ✅ โหลด RawMilk Smart Contract
	rawMilkInstance, err := rawmilk.NewRawmilk(rawMilkContractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load raw milk contract: %v", err)
	}

	userRegistryAddress := os.Getenv("USER_REGISTRY_CONTRACT_ADDRESS")
	if userRegistryAddress == "" {
		return fmt.Errorf("❌ USER_REGISTRY_CONTRACT_ADDRESS is not set")
	}

	userRegistryAddr := common.HexToAddress(userRegistryAddress)

	// ✅ โหลด UserRegistry Smart Contract
	userRegistryInstance, err := userregistry.NewUserregistry(userRegistryAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load user registry contract: %v", err)
	}

	BlockchainServiceInstance = &BlockchainService{
		client:                client,
		auth:                  auth,
		userRegistryContract:  userRegistryInstance,
		certificationContract: certInstance,
		rawMilkContract:       rawMilkInstance,
	}

	fmt.Println("✅ Blockchain Service Initialized!")
	return nil
}

func (b *BlockchainService) RegisterUserOnBlockchain(userWallet string, role uint8) (string, error) {
	fmt.Println("📌 Registering User on Blockchain:", userWallet, "Role:", role)

	userAddress := common.HexToAddress(userWallet)

	// ✅ เช็คก่อนว่า User ลงทะเบียนไปแล้วหรือยัง
	isRegistered, err := b.CheckUserOnBlockchain(userWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to check user registration: %v", err)
	}
	if isRegistered {
		fmt.Println("✅ User is already registered on blockchain:", userWallet)
		return "", fmt.Errorf("❌ User is already registered")
	}

	// ✅ ใช้ `userAddress` แทน `b.auth` ใน `From`
	opts := &bind.TransactOpts{
		From:   userAddress, // ✅ ให้ User เป็นคนส่ง Transaction เอง
		Signer: b.auth.Signer,
	}

	// ✅ ลงทะเบียน User ใน Smart Contract `UserRegistry`
	tx, err := b.userRegistryContract.RegisterUser(opts, role)
	if err != nil {
		log.Println("❌ Failed to register user on blockchain:", err)
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

	fmt.Println("✅ User registered on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (b *BlockchainService) CheckUserOnBlockchain(userWallet string) (bool, error) {
	fmt.Println("📌 Checking if user exists on blockchain:", userWallet)

	callOpts := &bind.CallOpts{Pending: false, Context: context.Background()}
	userAddress := common.HexToAddress(userWallet)

	isRegistered, err := b.userRegistryContract.IsUserRegistered(callOpts, userAddress)
	if err != nil {
		fmt.Println("❌ Failed to check user on blockchain:", err)
		return false, err
	}

	return isRegistered, nil
}

func (b *BlockchainService) StoreCertificationOnBlockchain(eventID, entityType, entityID, certCID string, issuedDate, expiryDate *big.Int) (string, error) {
	fmt.Println("📌 Checking existing certifications before storing new one...")

	// ✅ เช็คว่าผู้ใช้ลงทะเบียนในระบบแล้ว
	callOpts := &bind.CallOpts{Pending: false, Context: context.Background()}
	isRegistered, err := b.userRegistryContract.IsUserRegistered(callOpts, b.auth.From)
	if err != nil {
		fmt.Println("❌ Failed to check user registration:", err)
		return "", err
	}
	if !isRegistered {
		return "", errors.New("❌ User is not registered in the system")
	}

	// ✅ ดึงใบเซอร์ทั้งหมดของ entityID
	existingCerts, err := b.GetAllCertificationsForEntity(entityID)
	if err != nil {
		fmt.Println("❌ Failed to fetch existing certifications:", err)
		return "", err
	}

	// ✅ ตรวจสอบว่ามีใบเซอร์ที่ Active อยู่หรือไม่
	for _, cert := range existingCerts {
		if cert.IsActive {
			fmt.Println("📌 Found active certification, deactivating before storing new one:", cert.EventID)
			_, err := b.DeactivateCertificationOnBlockchain(cert.EventID)
			if err != nil {
				fmt.Println("❌ Failed to deactivate existing certification:", err)
				return "", err
			}
		}
	}

	fmt.Println("📌 Storing new certification on Blockchain...")

	opts := &bind.TransactOpts{
		From:     b.auth.From,
		Signer:   b.auth.Signer,
		Value:    big.NewInt(0),
		GasLimit: 800000,
	}

	// ✅ ส่งธุรกรรมไปยัง Smart Contract
	tx, err := b.certificationContract.StoreCertificationEvent(opts, eventID, entityType, entityID, certCID, issuedDate, expiryDate)
	if err != nil {
		fmt.Println("❌ Failed to store certification event on blockchain:", err)
		return "", err
	}

	// ✅ รอให้ธุรกรรมถูก Mine
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		fmt.Println("❌ Transaction not mined:", err)
		return "", err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		fmt.Println("❌ Transaction failed!")
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

	// ✅ รอให้ธุรกรรมถูก Mine
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		log.Println("❌ Transaction not mined:", err)
		return "", err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		log.Println("❌ Transaction failed!")
		return "", errors.New("transaction failed")
	}

	fmt.Println("✅ Certification Event deactivated on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (b *BlockchainService) GetAllCertificationsForEntity(entityID string) ([]certification.CertificationEventCertEvent, error) {
	callOpts := &bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}

	certs, err := b.certificationContract.GetActiveCertificationsForEntity(callOpts, entityID)
	if err != nil {
		log.Println("❌ Failed to fetch certifications from blockchain:", err)
		return nil, err
	}

	// ✅ กรองใบเซอร์ที่ `isActive == true` เท่านั้น
	var activeCerts []certification.CertificationEventCertEvent
	for _, cert := range certs {
		if cert.IsActive {
			activeCerts = append(activeCerts, cert)
		}
	}

	fmt.Println("✅ Retrieved active certifications from blockchain:", activeCerts)
	return activeCerts, nil
}

func (b *BlockchainService) CheckUserCertification(certCID string) (bool, error) {
	fmt.Println("📌 Checking if Certification CID is unique:", certCID)

	callOpts := &bind.CallOpts{Pending: false, Context: context.Background()}

	// ✅ ดึงข้อมูล "ทุกใบเซอร์" ที่เคยบันทึกไว้ใน Blockchain
	allCerts, err := b.certificationContract.GetAllCertifications(callOpts)
	if err != nil {
		fmt.Println("❌ Failed to fetch all certifications:", err)
		return false, err
	}

	// ✅ ตรวจสอบว่า CID นี้เคยถูกใช้มาก่อนหรือไม่
	for _, cert := range allCerts {
		if cert.CertificationCID == certCID {
			fmt.Println("❌ Certification CID already exists on blockchain:", cert.EventID)
			return false, nil
		}
	}

	fmt.Println("✅ Certification CID is unique, can be stored")
	return true, nil
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
