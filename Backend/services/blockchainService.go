package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"

	"golang.org/x/crypto/sha3" // ✅ ใช้จาก external package	"encoding/json"

	// ✅ ใช้ไลบรารีที่ถูกต้อง	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	certification "finalyearproject/Backend/services/certification_event" // ✅ สำหรับ Raw Milk
	"finalyearproject/Backend/services/product"
	"finalyearproject/Backend/services/rawmilk"
	"finalyearproject/Backend/services/userregistry"
)

// BlockchainService - ใช้สำหรับเชื่อมต่อ Blockchain
type BlockchainService struct {
	client                *ethclient.Client
	auth                  *bind.TransactOpts
	userRegistryContract  *userregistry.Userregistry
	certificationContract *certification.Certification
	rawMilkContract       *rawmilk.Rawmilk
	productContract       *product.Product
}

func getChainID() *big.Int {
	chainIDStr := os.Getenv("GANACHE_CHAIN_ID")
	chainID, err := strconv.ParseInt(chainIDStr, 10, 64)
	if err != nil {
		chainID = 1337 // ✅ ค่า Default ถ้าไม่มีการกำหนดค่าใน .env
	}
	return big.NewInt(chainID)
}

// BlockchainServiceInstance - Global Instance
var BlockchainServiceInstance *BlockchainService

// InitBlockchainService - เชื่อมต่อ Blockchain และโหลดคอนแทรค
func InitBlockchainService() error {
	fmt.Println("🚀 Initializing Blockchain Service...")

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
	if privateKeyHex == "" {
		return fmt.Errorf("❌ PRIVATE_KEY is not set")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("❌ Invalid private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		return fmt.Errorf("❌ Failed to create transaction auth: %v", err)
	}

	// ✅ โหลด Smart Contract Address จาก ENV
	certContractAddress := os.Getenv("CERT_CONTRACT_ADDRESS")
	rawMilkContractAddress := os.Getenv("RAWMILK_CONTRACT_ADDRESS")
	userRegistryAddress := os.Getenv("USER_REGISTRY_CONTRACT_ADDRESS")
	productContractAddress := os.Getenv("PRODUCT_CONTRACT_ADDRESS") // ✅ เพิ่ม Product Contract

	if certContractAddress == "" || rawMilkContractAddress == "" || userRegistryAddress == "" || productContractAddress == "" {
		return fmt.Errorf("❌ Missing blockchain contract addresses")
	}

	// ✅ แปลง Address จาก String เป็น Ethereum Address
	certContractAddr := common.HexToAddress(certContractAddress)
	rawMilkContractAddr := common.HexToAddress(rawMilkContractAddress)
	userRegistryAddr := common.HexToAddress(userRegistryAddress)
	productContractAddr := common.HexToAddress(productContractAddress) // ✅ เพิ่ม Product Contract Address

	// ✅ โหลด Certification Contract
	certInstance, err := certification.NewCertification(certContractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load certification contract: %v", err)
	}

	// ✅ โหลด RawMilk Contract
	rawMilkInstance, err := rawmilk.NewRawmilk(rawMilkContractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load raw milk contract: %v", err)
	}

	// ✅ โหลด UserRegistry Contract
	userRegistryInstance, err := userregistry.NewUserregistry(userRegistryAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load user registry contract: %v", err)
	}

	// ✅ โหลด Product Contract
	productInstance, err := product.NewProduct(productContractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load product contract: %v", err)
	}

	BlockchainServiceInstance = &BlockchainService{
		client:                client,
		auth:                  auth,
		userRegistryContract:  userRegistryInstance,
		certificationContract: certInstance,
		rawMilkContract:       rawMilkInstance,
		productContract:       productInstance, // ✅ เพิ่ม Product Instance
	}

	fmt.Println("✅ Blockchain Service Initialized!")
	return nil
}
func (b *BlockchainService) getPrivateKeyForAddress(userWallet string) (string, error) {
	// ✅ กำหนด path ที่ถูกต้อง
	filePath := "services/private_keys.json"

	// Debug: เช็คว่าไฟล์อยู่ตรงไหน
	absPath, _ := os.Getwd()
	fmt.Println("📌 Looking for private_keys.json at:", absPath+"/"+filePath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("❌ Failed to load private keys file from:", absPath+"/"+filePath)
		return "", errors.New("Failed to load private keys file")
	}
	fmt.Println("✅ Loaded private_keys.json successfully")

	var privateKeys map[string]string
	err = json.Unmarshal(data, &privateKeys)
	if err != nil {
		return "", errors.New("Failed to parse private keys")
	}

	if key, exists := privateKeys[userWallet]; exists {
		return key, nil
	}
	return "", errors.New("Private key not found for address")
}

func (b *BlockchainService) RegisterUserOnBlockchain(userWallet string, role uint8) (string, error) {
	fmt.Println("📌 Registering User on Blockchain:", userWallet, "Role:", role)

	userAddress := common.HexToAddress(userWallet)

	// ✅ เช็คก่อนว่า User ลงทะเบียนไปแล้วหรือยัง
	fmt.Println("📌 Checking if user exists on blockchain:", userWallet)
	isRegistered, err := b.CheckUserOnBlockchain(userWallet)
	if err != nil {
		fmt.Println("❌ Error checking user registration:", err)
		return "", fmt.Errorf("❌ Failed to check user registration: %v", err)
	}
	if isRegistered {
		fmt.Println("✅ User is already registered on blockchain:", userWallet)
		return "", fmt.Errorf("❌ User is already registered")
	}

	// ✅ ดึง Private Key ของ Wallet ที่สุ่มมาให้ User
	fmt.Println("📌 Fetching Private Key for:", userWallet)
	privateKeyHex, err := b.getPrivateKeyForAddress(userWallet)
	if err != nil {
		fmt.Println("❌ Failed to get private key:", err)
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}
	fmt.Println("✅ Private Key Found:", privateKeyHex[:10]+"...") // โชว์แค่ 10 ตัวแรก

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		fmt.Println("❌ Failed to parse private key:", err)
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}
	fmt.Println("✅ Private Key Parsed Successfully")

	// ✅ สร้าง TransactOpts ใหม่ โดยใช้ Private Key ของ User
	fmt.Println("📌 Creating Transaction Auth")
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		fmt.Println("❌ Failed to create transactor:", err)
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = userAddress
	fmt.Println("✅ Transactor Created - From:", auth.From.Hex())

	// ✅ ลงทะเบียน User ใน Smart Contract `UserRegistry`
	fmt.Println("📌 Sending Transaction to Register User...")
	tx, err := b.userRegistryContract.RegisterUser(auth, role)
	if err != nil {
		fmt.Println("❌ Failed to register user on blockchain:", err)
		return "", err
	}
	fmt.Println("✅ Transaction Sent:", tx.Hash().Hex())

	// ✅ รอให้ Transaction ถูกบันทึก
	fmt.Println("📌 Waiting for Transaction to be Mined...")
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		fmt.Println("❌ Transaction not mined:", err)
		return "", err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		fmt.Println("❌ Transaction failed!")
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

func (b *BlockchainService) StoreCertificationOnBlockchain(walletAddress, eventID, entityType, entityID, certCID string, issuedDate, expiryDate *big.Int) (string, error) {
	fmt.Println("📌 Checking user registration before storing new certification...")

	// ✅ เช็คว่าผู้ใช้ลงทะเบียนในระบบแล้ว
	callOpts := &bind.CallOpts{Pending: false, Context: context.Background()}
	isRegistered, err := b.userRegistryContract.IsUserRegistered(callOpts, common.HexToAddress(walletAddress))
	if err != nil {
		fmt.Println("❌ Failed to check user registration:", err)
		return "", err
	}
	if !isRegistered {
		fmt.Println("❌ User is not registered in the system")
		return "", errors.New("User is not registered in the system")
	}

	fmt.Println("📌 Fetching Private Key for:", walletAddress)

	// ✅ ดึง Private Key ของ User จากไฟล์ JSON
	privateKeyHex, err := b.getPrivateKeyForAddress(walletAddress)
	if err != nil {
		fmt.Println("❌ Failed to get private key:", err)
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}
	fmt.Println("✅ Private Key Found:", privateKeyHex[:10]+"...")

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		fmt.Println("❌ Failed to parse private key:", err)
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}
	fmt.Println("✅ Private Key Parsed Successfully")

	// ✅ สร้าง `auth` ใหม่โดยใช้ Private Key ของ User
	fmt.Println("📌 Creating Transaction Auth for:", walletAddress)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		fmt.Println("❌ Failed to create transactor:", err)
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(walletAddress) // ✅ ใช้ Wallet Address ของ User
	fmt.Println("✅ Transactor Created - From:", auth.From.Hex())

	fmt.Println("📌 Storing new certification on Blockchain...")

	// ✅ ส่งธุรกรรมไปยัง Smart Contract
	tx, err := b.certificationContract.StoreCertificationEvent(
		auth,
		eventID,
		entityType,
		entityID,
		certCID,
		issuedDate,
		expiryDate,
	)
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
		return "", errors.New("Transaction failed")
	}

	fmt.Println("✅ Certification Event stored on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (b *BlockchainService) DeactivateCertificationOnBlockchain(walletAddress, eventID string) (string, error) {
	fmt.Println("📌 [Blockchain] Deactivating certification for Wallet:", walletAddress, "EventID:", eventID)

	// ✅ ตรวจสอบว่า `walletAddress` ลงทะเบียนใน Blockchain แล้ว
	callOpts := &bind.CallOpts{Pending: false, Context: context.Background()}
	isRegistered, err := b.userRegistryContract.IsUserRegistered(callOpts, common.HexToAddress(walletAddress))
	if err != nil {
		fmt.Println("❌ Failed to check user registration:", err)
		return "", err
	}
	if !isRegistered {
		fmt.Println("❌ User is not registered in the system")
		return "", errors.New("User is not registered in the system")
	}

	// ✅ ดึง Private Key ของ User จากไฟล์ JSON
	fmt.Println("📌 Fetching Private Key for:", walletAddress)
	privateKeyHex, err := b.getPrivateKeyForAddress(walletAddress)
	if err != nil {
		fmt.Println("❌ Failed to get private key:", err)
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}
	fmt.Println("✅ Private Key Found:", privateKeyHex[:10]+"...")

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		fmt.Println("❌ Failed to parse private key:", err)
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}
	fmt.Println("✅ Private Key Parsed Successfully")

	// ✅ สร้าง `auth` ใหม่โดยใช้ Private Key ของ User
	fmt.Println("📌 Creating Transaction Auth for:", walletAddress)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		fmt.Println("❌ Failed to create transactor:", err)
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(walletAddress) // ✅ ใช้ Wallet Address ของ User
	fmt.Println("✅ Transactor Created - From:", auth.From.Hex())

	// ✅ ส่งธุรกรรมไปยัง Blockchain
	tx, err := b.certificationContract.DeactivateCertificationEvent(auth, eventID)
	if err != nil {
		log.Println("❌ [Blockchain] Failed to deactivate certification event on blockchain:", err)
		return "", err
	}

	// ✅ รอให้ธุรกรรมถูก Mine
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		log.Println("❌ [Blockchain] Transaction not mined:", err)
		return "", err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		log.Println("❌ [Blockchain] Transaction failed!")
		return "", errors.New("transaction failed")
	}

	fmt.Println("✅ [Blockchain] Certification Event deactivated on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (b *BlockchainService) GetAllCertificationsForEntity(entityID string) ([]certification.CertificationEventCertEvent, error) {
	callOpts := &bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}

	fmt.Println("📌 Fetching active certifications for entity:", entityID)

	// ✅ เรียก Smart Contract
	certs, err := b.certificationContract.GetActiveCertificationsForEntity(callOpts, entityID)
	if err != nil {
		log.Println("❌ Failed to fetch certifications from blockchain:", err)
		return nil, err
	}

	// ✅ ถ้าไม่มีใบเซอร์เลย -> คืนค่าเป็น [] แทน nil เพื่อป้องกัน Panic
	if len(certs) == 0 {
		fmt.Println("📌 No certifications found for entity:", entityID)
		return []certification.CertificationEventCertEvent{}, nil
	}

	// ✅ กรองเฉพาะใบเซอร์ที่ยัง `isActive == true`
	var activeCerts []certification.CertificationEventCertEvent
	for _, cert := range certs {
		if cert.IsActive {
			activeCerts = append(activeCerts, cert)
		}
	}

	fmt.Println("✅ Retrieved active certifications from blockchain:", len(activeCerts))
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

// //////////////////////////////////////////////////////////// RawMilk /////////////////////////////////////////////////////////
// ✅ อัปเดต Struct `RawMilkData` ให้รองรับ `qualityReportCID`
type RawMilkData struct {
	TankId           string `json:"tankId"`
	FarmWallet       string `json:"farmWallet"`
	FactoryId        string `json:"factoryId"`
	PersonInCharge   string `json:"personInCharge"`
	QualityReportCID string `json:"qualityReportCid"` // ✅ เพิ่ม
	QrCodeCID        string `json:"qrCodeCid"`
	Status           uint8  `json:"status"`
}

// ✅ ฟังก์ชันสร้างแท้งค์นมบนบล็อกเชน (แก้ `factoryId` เป็น `bytes32`)
func (b *BlockchainService) CreateMilkTank(
	userWallet string,
	tankId string,
	factoryId string, // ✅ รับ FactoryID แบบ string (UUID หรือ Custom ID)
	personInCharge string,
	qualityReportCID string,
	qrCodeCID string,
) (string, error) {

	fmt.Println("📌 Creating Milk Tank on Blockchain for:", userWallet)

	// ✅ ตรวจสอบค่าก่อนส่งธุรกรรม
	err := validateMilkTankData(userWallet, tankId, factoryId, personInCharge, qualityReportCID, qrCodeCID)
	if err != nil {
		return "", err
	}

	// ✅ ดึง Private Key ของ Wallet ของเกษตรกร
	privateKeyHex, err := b.getPrivateKeyForAddress(userWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	// ✅ สร้าง Transaction Auth โดยใช้ Private Key ของเกษตรกร
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(userWallet)
	auth.GasLimit = uint64(3000000)         // ✅ เพิ่ม Gas Limit
	auth.GasPrice = big.NewInt(20000000000) // ✅ กำหนด Gas Price

	// ✅ แปลง `tankId` เป็น `bytes32`
	tankIdBytes := common.BytesToHash([]byte(tankId))

	// ✅ แปลง `factoryId` เป็น `bytes32`
	factoryIdBytes := common.BytesToHash([]byte(factoryId))

	// ✅ Debug Log ก่อนส่งไปยัง Blockchain
	fmt.Println("📌 Debug - Sending to Blockchain:")
	fmt.Println("   - Tank ID (Bytes32):", tankIdBytes)
	fmt.Println("   - Factory ID (Bytes32):", factoryIdBytes)
	fmt.Println("   - Person In Charge:", personInCharge)

	// ✅ ส่งธุรกรรมไปที่ Smart Contract
	tx, err := b.rawMilkContract.CreateMilkTank(
		auth,
		tankIdBytes,    // ✅ ใช้ [32]byte
		factoryIdBytes, // ✅ ใช้ [32]byte
		personInCharge,
		qualityReportCID,
		qrCodeCID,
	)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create milk tank on blockchain: %v", err)
	}

	fmt.Println("✅ Transaction Sent:", tx.Hash().Hex())

	// ✅ รอให้ Transaction ถูกบันทึก
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		return "", fmt.Errorf("❌ Transaction not mined: %v", err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return "", errors.New("❌ Transaction failed")
	}

	fmt.Println("✅ Milk Tank Created on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil

}

// ✅ ฟังก์ชันตรวจสอบค่าก่อนสร้างธุรกรรม (แก้ `factoryId` ให้เป็น `bytes32`)
func validateMilkTankData(userWallet, tankId, factoryId, personInCharge, qualityReportCID, qrCodeCID string) error {
	if userWallet == "" {
		return errors.New("❌ userWallet is required")
	}
	if !common.IsHexAddress(userWallet) {
		return errors.New("❌ userWallet is not a valid Ethereum address")
	}
	if tankId == "" {
		return errors.New("❌ tankId is required")
	}
	if factoryId == "" {
		return errors.New("❌ factoryId is required")
	}
	if personInCharge == "" {
		return errors.New("❌ personInCharge is required")
	}
	if qualityReportCID == "" {
		return errors.New("❌ qualityReportCID is required")
	}
	if qrCodeCID == "" {
		return errors.New("❌ qrCodeCID is required")
	}
	return nil
}

func (b *BlockchainService) GetMilkTanksByFarmer(farmerAddress string) ([]map[string]interface{}, error) {
	fmt.Println("📌 Fetching milk tanks for farmer:", farmerAddress)

	// ✅ แปลงที่อยู่ของฟาร์มจาก string เป็น Ethereum Address
	farmer := common.HexToAddress(farmerAddress)

	// ✅ ดึงรายการ Tank IDs และประวัติจาก Smart Contract
	tankIDs, histories, err := b.rawMilkContract.GetMilkTanksByFarmer(&bind.CallOpts{}, farmer)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tanks for farmer:", err)
		return nil, err
	}

	var milkTanks []map[string]interface{}

	// ✅ วนลูปดึงข้อมูลของแต่ละแท็งก์
	for i, id := range tankIDs {
		tankIdStr := string(bytes.Trim(id[:], "\x00"))

		// ✅ ใช้ข้อมูลจากประวัติล่าสุด
		latestEntry := histories[i][len(histories[i])-1]

		// ✅ หาค่า OLDPERSONINCHARGE (ประวัติรองสุดท้าย ถ้ามี)
		var oldPersonInCharge string
		if len(histories[i]) > 1 {
			oldPersonInCharge = histories[i][len(histories[i])-2].PersonInCharge
		} else {
			oldPersonInCharge = latestEntry.PersonInCharge // ถ้าไม่มีข้อมูลเก่า ให้ใช้ค่าปัจจุบัน
		}

		milkTank := map[string]interface{}{
			"tankId":            tankIdStr,
			"personInCharge":    latestEntry.PersonInCharge,
			"oldPersonInCharge": oldPersonInCharge,
			"status":            uint8(latestEntry.Status),
		}

		milkTanks = append(milkTanks, milkTank)
	}

	fmt.Println("✅ Fetched milk tanks for farmer (All statuses):", farmerAddress, milkTanks)
	return milkTanks, nil
}

func (b *BlockchainService) GetRawMilkTankDetails(tankId string) (*RawMilkData, []map[string]interface{}, error) {
	fmt.Println("📌 Fetching milk tank details for:", tankId)

	// ✅ แปลง tankId เป็น bytes32
	tankIdBytes := common.BytesToHash([]byte(tankId))

	// ✅ ดึงข้อมูลแท็งก์จาก Smart Contract
	milkTankData, err := b.rawMilkContract.GetMilkTank(&bind.CallOpts{}, tankIdBytes)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tank details:", err)
		return nil, nil, fmt.Errorf("❌ Failed to fetch milk tank details: %v", err)
	}

	// ✅ แปลงค่า tankId และ factoryId เป็น string
	tankIdStr := string(bytes.Trim(milkTankData.TankId[:], "\x00"))
	factoryIdStr := string(bytes.Trim(milkTankData.FactoryId[:], "\x00"))

	// ✅ แปลงข้อมูลแท็งก์เป็นโครงสร้าง `RawMilkData`
	rawMilk := &RawMilkData{
		TankId:           tankIdStr,
		FarmWallet:       milkTankData.Farmer.Hex(),
		FactoryId:        factoryIdStr,
		PersonInCharge:   milkTankData.PersonInCharge,
		QualityReportCID: milkTankData.QualityReportCID,
		QrCodeCID:        milkTankData.QrCodeCID,
		Status:           uint8(milkTankData.Status),
	}

	// ✅ ดึงประวัติของแท็งก์จาก Smart Contract
	historyData := milkTankData.History // `History` มาจาก `MilkTankWithHistory`

	// ✅ สร้างอาร์เรย์เก็บประวัติการเปลี่ยนแปลง
	var historyList []map[string]interface{}
	for _, entry := range historyData {
		historyList = append(historyList, map[string]interface{}{
			"personInCharge":   entry.PersonInCharge,
			"qualityReportCID": entry.QualityReportCID,
			"status":           uint8(entry.Status),
			"timestamp":        entry.Timestamp,
		})
	}

	fmt.Println("✅ Milk Tank Details Retrieved:", rawMilk)
	fmt.Println("✅ Milk Tank History Retrieved:", historyList)
	return rawMilk, historyList, nil
}

func (b *BlockchainService) GetMilkTanksByFactory(factoryID string) ([]map[string]interface{}, error) {
	fmt.Println("📌 Fetching milk tanks for factory:", factoryID)

	// ✅ แปลง FactoryID เป็น bytes32
	factoryIDBytes32 := common.BytesToHash([]byte(factoryID))
	fmt.Println("🔍 [Fixed] Converted FactoryID to Bytes32:", factoryIDBytes32)

	// ✅ ดึงรายการ Tank IDs และประวัติจาก Smart Contract
	tankIDs, histories, err := b.rawMilkContract.GetMilkTanksByFactory(&bind.CallOpts{}, factoryIDBytes32)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tanks for factory:", err)
		return nil, err
	}

	var milkTanks []map[string]interface{}

	// ✅ วนลูปดึงข้อมูลของแต่ละแท็งก์
	for i, id := range tankIDs {
		tankIdStr := string(bytes.Trim(id[:], "\x00"))

		// ✅ ใช้ข้อมูลจากประวัติล่าสุด (อันสุดท้ายใน Array)
		latestEntry := histories[i][len(histories[i])-1]

		// ✅ สร้าง JSON Response ที่มี `tankId`, `personInCharge`, `status` (ทุกสถานะ)
		milkTank := map[string]interface{}{
			"tankId":         tankIdStr,
			"personInCharge": latestEntry.PersonInCharge,
			"status":         uint8(latestEntry.Status),
		}

		milkTanks = append(milkTanks, milkTank)
	}

	fmt.Println("✅ Fetched all milk tanks for factory:", factoryID, milkTanks)
	return milkTanks, nil
}

func (b *BlockchainService) UpdateMilkTankStatus(
	factoryWallet string,
	tankId string,
	approved bool, // ✅ true = Approved, false = Rejected
	personInCharge string,
	qualityReportCID string,
) (string, error) {

	fmt.Println("📌 Updating Milk Tank Status on Blockchain for Factory:", factoryWallet)

	// ✅ ตรวจสอบค่าก่อนส่งธุรกรรม
	if factoryWallet == "" || tankId == "" || personInCharge == "" || qualityReportCID == "" {
		return "", fmt.Errorf("❌ Missing required fields")
	}

	// ✅ ดึง Private Key ของ Wallet ของโรงงาน
	privateKeyHex, err := b.getPrivateKeyForAddress(factoryWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	// ✅ สร้าง Transaction Auth โดยใช้ Private Key ของโรงงาน
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(factoryWallet)
	auth.GasLimit = uint64(3000000)         // ✅ เพิ่ม Gas Limit
	auth.GasPrice = big.NewInt(20000000000) // ✅ กำหนด Gas Price

	// ✅ แปลง `tankId` เป็น `bytes32`
	tankIdBytes := common.BytesToHash([]byte(tankId))

	// ✅ Debug Log ก่อนส่งไปยัง Blockchain
	fmt.Println("📌 Debug - Sending to Blockchain:")
	fmt.Println("   - Tank ID (Bytes32):", tankIdBytes)
	fmt.Println("   - Person In Charge:", personInCharge)
	fmt.Println("   - Approved:", approved)

	// ✅ ส่งธุรกรรมไปที่ Smart Contract
	tx, err := b.rawMilkContract.VerifyMilkQuality(
		auth,
		tankIdBytes, // ✅ ใช้ [32]byte
		approved,    // ✅ อัปเดตเป็น Approved หรือ Rejected
		qualityReportCID,
		personInCharge,
	)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to update milk tank status on blockchain: %v", err)
	}

	fmt.Println("✅ Transaction Sent:", tx.Hash().Hex())

	// ✅ รอให้ Transaction ถูกบันทึก
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		return "", fmt.Errorf("❌ Transaction not mined: %v", err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return "", errors.New("❌ Transaction failed")
	}

	fmt.Println("✅ Milk Tank Status Updated on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// //////////////////////////////////////////////////////////// Product /////////////////////////////////////////////////////////
// ✅ ฟังก์ชันแปลง string เป็น bytes32 แบบตรงไปตรงมา (ไม่ใช้ hash)
func stringToBytes32(str string) [32]byte {
	hash := sha3.NewLegacyKeccak256() // ✅ ใช้ Keccak-256
	hash.Write([]byte(str))
	var bytes32 [32]byte
	copy(bytes32[:], hash.Sum(nil)) // ✅ ใช้ค่า Hash 32 Byte
	return bytes32
}

// ✅ ฟังก์ชันสร้างสินค้า
func (b *BlockchainService) CreateProduct(
	factoryWallet string,
	productId string,
	productName string,
	productCID string,
	category string,
) (string, error) {

	fmt.Println("📌 Creating Product on Blockchain for Wallet:", factoryWallet)

	// ✅ ตรวจสอบค่าก่อนส่งธุรกรรม
	err := validateProductData(factoryWallet, productId, productName, productCID, category)
	if err != nil {
		return "", err
	}

	// ✅ ดึง Private Key ของ Factory Wallet
	privateKeyHex, err := b.getPrivateKeyForAddress(factoryWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	role, err := b.userRegistryContract.GetUserRole(nil, common.HexToAddress(factoryWallet))
	fmt.Println("✅ User Role on Blockchain:", role)

	// ✅ สร้าง Transaction Auth โดยใช้ Private Key ของ Factory
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(factoryWallet) // ✅ ต้องให้ตรงกับ Wallet ของ Factory ที่ Register
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(20000000000)

	// ✅ แปลง `productId` เป็น `bytes32`
	productIDBytes := stringToBytes32(productId)

	// ✅ Debug Log
	fmt.Println("📌 Debug - Sending to Blockchain:")
	fmt.Println("   - Product ID (Bytes32):", productIDBytes)
	fmt.Println("   - Factory Wallet:", factoryWallet)
	fmt.Println("   - Product Name:", productName)
	fmt.Println("   - Product CID:", productCID)
	fmt.Println("   - Category:", category)

	// ✅ ส่งธุรกรรมไปที่ Smart Contract
	tx, err := b.productContract.CreateProduct(
		auth,
		productIDBytes,
		productName,
		productCID,
		category,
	)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create product on blockchain: %v", err)
	}

	fmt.Println("✅ Transaction Sent:", tx.Hash().Hex())

	// ✅ รอให้ Transaction ถูกบันทึก
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		return "", fmt.Errorf("❌ Transaction not mined: %v", err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return "", errors.New("❌ Transaction failed")
	}

	fmt.Println("✅ Product Created on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func validateProductData(factoryWallet, productId, productName, productCID, category string) error {
	if factoryWallet == "" {
		return errors.New("❌ factoryWallet is required")
	}
	if !common.IsHexAddress(factoryWallet) {
		return errors.New("❌ factoryWallet is not a valid Ethereum address")
	}
	if productId == "" {
		return errors.New("❌ productId is required")
	}
	if productName == "" {
		return errors.New("❌ productName is required")
	}
	if productCID == "" {
		return errors.New("❌ productCID is required")
	}
	if category == "" {
		return errors.New("❌ category is required")
	}
	return nil
}
