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
	"strings"
	"time"

	// ✅ ใช้จาก external package	"encoding/json"

	// ✅ ใช้ไลบรารีที่ถูกต้อง	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	certification "finalyearproject/Backend/services/certification_event" // ✅ สำหรับ Raw Milk
	"finalyearproject/Backend/services/product"
	"finalyearproject/Backend/services/productlot"
	"finalyearproject/Backend/services/rawmilk"
	"finalyearproject/Backend/services/tracking"
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
	productLotContract    *productlot.Productlot
	trackingContract      *tracking.Tracking
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
	productContractAddress := os.Getenv("PRODUCT_CONTRACT_ADDRESS")
	productLotContractAddress := os.Getenv("PRODUCTLOT_CONTRACT_ADDRESS")
	trackingContractAddress := os.Getenv("TRACKING_CONTRACT_ADDRESS")

	if certContractAddress == "" || rawMilkContractAddress == "" || userRegistryAddress == "" || productContractAddress == "" || productLotContractAddress == "" {
		return fmt.Errorf("❌ Missing blockchain contract addresses")
	}

	// ✅ แปลง Address จาก String เป็น Ethereum Address
	certContractAddr := common.HexToAddress(certContractAddress)
	rawMilkContractAddr := common.HexToAddress(rawMilkContractAddress)
	userRegistryAddr := common.HexToAddress(userRegistryAddress)
	productContractAddr := common.HexToAddress(productContractAddress)
	productLotContractAddr := common.HexToAddress(productLotContractAddress)
	trackingContractAddr := common.HexToAddress(trackingContractAddress)

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

	// ✅ โหลด ProductLot Contract
	productLotInstance, err := productlot.NewProductlot(productLotContractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load product lot contract: %v", err)
	}

	// ✅ โหลด Tracking Contract
	trackingInstance, err := tracking.NewTracking(trackingContractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load tracking contract: %v", err)
	}

	BlockchainServiceInstance = &BlockchainService{
		client:                client,
		auth:                  auth,
		userRegistryContract:  userRegistryInstance,
		certificationContract: certInstance,
		rawMilkContract:       rawMilkInstance,
		productContract:       productInstance,
		productLotContract:    productLotInstance,
		trackingContract:      trackingInstance,
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

// ✅ ฟังก์ชันสร้างแท้งค์นมบนบล็อกเชน (อัปเดต Debug Log)
func (b *BlockchainService) CreateMilkTank(
	userWallet string,
	tankId string,
	factoryId string,
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

	// ✅ แปลง `tankId` และ `factoryId` เป็น `bytes32`
	tankIdBytes := common.BytesToHash([]byte(tankId))
	factoryIdBytes := common.BytesToHash([]byte(factoryId))

	// ✅ Debug Log ก่อนส่งไปยัง Blockchain
	fmt.Println("📌 Debug - Sending to Blockchain:")
	fmt.Println("   - Tank ID (Bytes32):", tankIdBytes)
	fmt.Println("   - Factory ID (Bytes32):", factoryIdBytes)
	fmt.Println("   - Person In Charge:", personInCharge)
	fmt.Println("   - Quality Report CID:", qualityReportCID)
	fmt.Println("   - QR Code CID:", qrCodeCID)

	// ✅ ส่งธุรกรรมไปที่ Smart Contract
	tx, err := b.rawMilkContract.CreateMilkTank(
		auth,
		tankIdBytes,
		factoryIdBytes,
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

		// ✅ personInCharge = จากฟาร์ม (entry แรกสุด)
		farmPersonInCharge := histories[i][0].PersonInCharge

		// ✅ oldPersonInCharge = จากโรงงาน (entry ที่สอง ถ้ามี)
		var factoryPersonInCharge string
		if len(histories[i]) > 1 {
			factoryPersonInCharge = histories[i][1].PersonInCharge
		} else {
			factoryPersonInCharge = "" // ยังไม่มีโรงงานรับ
		}

		// ✅ status = ล่าสุด (entry สุดท้าย)
		latestStatus := uint8(histories[i][len(histories[i])-1].Status)

		milkTank := map[string]interface{}{
			"tankId":            tankIdStr,
			"personInCharge":    farmPersonInCharge,
			"oldPersonInCharge": factoryPersonInCharge,
			"status":            latestStatus,
		}

		milkTanks = append(milkTanks, milkTank)
	}

	fmt.Println("✅ Fetched milk tanks for farmer (Farm & Factory PIC + Latest Status):", farmerAddress, milkTanks)
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

		// ✅ ดึง personInCharge ของฟาร์ม (entry แรก)
		farmPersonInCharge := histories[i][0].PersonInCharge

		// ✅ ดึง oldPersonInCharge ของโรงงาน (entry ที่สอง ถ้ามี)
		var factoryPersonInCharge string
		if len(histories[i]) > 1 {
			factoryPersonInCharge = histories[i][1].PersonInCharge
		} else {
			factoryPersonInCharge = "" // ยังไม่มีโรงงานรับ
		}

		// ✅ ดึง status ล่าสุด (entry สุดท้าย)
		latestStatus := uint8(histories[i][len(histories[i])-1].Status)

		// ✅ สร้าง JSON Response
		milkTank := map[string]interface{}{
			"tankId":            tankIdStr,
			"personInCharge":    farmPersonInCharge,
			"oldPersonInCharge": factoryPersonInCharge,
			"status":            latestStatus,
		}

		milkTanks = append(milkTanks, milkTank)
	}

	fmt.Println("✅ Fetched milk tanks for factory (Farm & Factory PIC + Latest Status):", factoryID, milkTanks)
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
// ✅ ฟังก์ชันสร้าง Product บนบล็อกเชน

// ✅ ฟังก์ชันสร้าง Product บนบล็อกเชน
func (b *BlockchainService) CreateProduct(
	userWallet string,
	productId string,
	productName string,
	productCID string,
	category string,
) (string, error) {

	fmt.Println("📌 Creating Product on Blockchain for Wallet:", userWallet)

	// ✅ ตรวจสอบค่าก่อนส่งธุรกรรม
	err := validateProductData(userWallet, productId, productName, productCID, category)
	if err != nil {
		return "", err
	}

	// ✅ ดึง Private Key ของ Wallet
	privateKeyHex, err := b.getPrivateKeyForAddress(userWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	// ✅ สร้าง Transaction Auth
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(userWallet)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(20000000000)

	// ✅ แปลง `productId` เป็น `[32]byte` แบบเดียวกับ `tankId`
	productIdBytes := common.BytesToHash([]byte(productId))

	// ✅ Debug Log ก่อนส่งไปยัง Blockchain
	fmt.Println("📌 Debug - Sending to Blockchain:")
	fmt.Println("   - Product ID (Bytes32):", productIdBytes) // ✅ ต้องออกมาเป็น 0x...
	fmt.Println("   - Product Name:", productName)
	fmt.Println("   - Product CID:", productCID)
	fmt.Println("   - Category:", category)

	// ✅ ส่งธุรกรรมไปที่ Smart Contract
	tx, err := b.productContract.CreateProduct(
		auth,
		productIdBytes, // ✅ แก้ให้เป็น `common.Hash`
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

func (b *BlockchainService) GetProductsByFactory(factoryAddress string) ([]map[string]interface{}, error) {
	fmt.Println("📌 Fetching products for factory:", factoryAddress)

	// ✅ แปลงที่อยู่ของโรงงานจาก string เป็น Ethereum Address
	factory := common.HexToAddress(factoryAddress)

	// ✅ ดึงรายการสินค้า จาก Smart Contract
	ids, names, categories, err := b.productContract.GetProductsByFactory(&bind.CallOpts{From: factory})
	if err != nil {
		fmt.Println("❌ Failed to fetch products:", err)
		return nil, err
	}

	var products []map[string]interface{}

	// ✅ วนลูปดึงข้อมูลของแต่ละ Product และลบ NULL Characters (`\x00`)
	for i := range ids {
		productIdStr := string(bytes.Trim(ids[i][:], "\x00"))

		product := map[string]interface{}{
			"productId":   productIdStr,
			"productName": names[i],
			"category":    categories[i],
			"detailsLink": fmt.Sprintf("/Factory/ProductDetails/%s", productIdStr),
		}
		products = append(products, product)
	}

	fmt.Println("✅ Fetched products for factory:", products)
	return products, nil
}

// ✅ ดึงรายละเอียดสินค้าตาม Product ID
func (b *BlockchainService) GetProductDetails(productId string) (map[string]interface{}, error) {
	fmt.Println("📌 Fetching product details:", productId)

	// ✅ ใช้ `common.BytesToHash([]byte(productId))` เหมือนตอนบันทึก
	productIdBytes := common.BytesToHash([]byte(productId))

	productData, err := b.productContract.GetProductDetails(&bind.CallOpts{}, productIdBytes)
	if err != nil {
		fmt.Println("❌ Failed to fetch product details:", err)
		return nil, err
	}

	product := map[string]interface{}{
		"productId":     strings.TrimRight(string(productData.ProductId[:]), "\x00"),
		"factoryWallet": productData.FactoryWallet.Hex(),
		"productName":   productData.ProductName,
		"productCID":    productData.ProductCID,
		"category":      productData.Category,
	}

	fmt.Println("✅ Product details fetched successfully:", product)
	return product, nil
}

// //////////////////////////////////////////////////////////// ProductLot /////////////////////////////////////////////////////////

type ProductLotInfo struct {
	LotID                  string
	ProductID              string
	Factory                string
	Inspector              string
	InspectionDate         time.Time
	Grade                  bool
	QualityAndNutritionCID string
	MilkTankIDs            []string
	Status                 uint8
}

func (b *BlockchainService) CreateProductLot(
	userWallet string,
	lotId string,
	productId string,
	inspector string,
	grade bool,
	qualityAndNutritionCID string,
	milkTankIds []string,
) (string, error) {

	fmt.Println("📌 Creating Product Lot on Blockchain for:", userWallet)

	// ✅ ตรวจสอบค่าก่อนส่งธุรกรรม
	err := validateProductLotData(userWallet, lotId, productId, inspector, strconv.FormatBool(grade), qualityAndNutritionCID, milkTankIds)
	if err != nil {
		return "", err
	}

	// ✅ ดึง Private Key
	privateKeyHex, err := b.getPrivateKeyForAddress(userWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	// ✅ Auth
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(userWallet)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(20000000000)

	// ✅ กลับมาใช้ common.BytesToHash → แบบเก่า!
	lotIdBytes := StringToBytes32(lotId)
	productIdBytes := common.BytesToHash([]byte(productId))

	// ✅ Milk Tanks
	var milkTankBytes [][32]byte
	for _, tankId := range milkTankIds {
		tankBytes := common.BytesToHash([]byte(tankId)) // แบบเดิม
		milkTankBytes = append(milkTankBytes, tankBytes)
	}

	// ✅ Debug
	fmt.Println("📌 Debug - Sending to Blockchain:")
	fmt.Println("   - Lot ID (Bytes32):", lotIdBytes)
	fmt.Println("   - Product ID (Bytes32):", productIdBytes)
	fmt.Println("   - Inspector:", inspector)
	fmt.Println("   - Inspection Date:", time.Now().Unix())
	fmt.Println("   - Grade:", grade)
	fmt.Println("   - Quality & Nutrition CID:", qualityAndNutritionCID)
	fmt.Println("   - Milk Tanks:", milkTankBytes)

	// ✅ ส่งธุรกรรม
	tx, err := b.productLotContract.CreateProductLot(
		auth,
		lotIdBytes,
		productIdBytes,
		inspector,
		grade,
		qualityAndNutritionCID,
		milkTankBytes,
	)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create product lot on blockchain: %v", err)
	}

	fmt.Println("✅ Transaction Sent:", tx.Hash().Hex())

	// ✅ รอ
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		return "", fmt.Errorf("❌ Transaction not mined: %v", err)
	}
	if receipt.Status == types.ReceiptStatusFailed {
		return "", errors.New("❌ Transaction failed")
	}

	fmt.Println("✅ Product Lot Created on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// ✅ ฟังก์ชันตรวจสอบค่าก่อนสร้างธุรกรรม
func validateProductLotData(userWallet, lotId, productId, inspector, grade, qualityAndNutritionCID string, milkTankIds []string) error {
	if userWallet == "" {
		return errors.New("❌ userWallet is required")
	}
	if !common.IsHexAddress(userWallet) {
		return errors.New("❌ userWallet is not a valid Ethereum address")
	}
	if lotId == "" {
		return errors.New("❌ lotId is required")
	}
	if productId == "" {
		return errors.New("❌ productId is required")
	}
	if inspector == "" {
		return errors.New("❌ inspector is required")
	}
	if grade == "" {
		return errors.New("❌ grade is required")
	}
	if qualityAndNutritionCID == "" {
		return errors.New("❌ qualityAndNutritionCID is required")
	}
	if len(milkTankIds) == 0 {
		return errors.New("❌ milkTankIds cannot be empty")
	}
	return nil
}

// ✅ ดึงข้อมูล Product Lot ตาม `productId`
func (b *BlockchainService) GetProductLotByLotID(lotId string) (*ProductLotInfo, error) {
	fmt.Println("📌 Fetching Product Lot for Lot ID:", lotId)

	// ✅ แปลง `lotId` เป็น `bytes32`
	lotIdBytes := StringToBytes32(lotId)

	// ✅ เรียก Smart Contract เพื่อนำข้อมูล Product Lot ออกมา
	productLotData, err := b.productLotContract.GetProductLot(nil, lotIdBytes)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch Product Lot: %v", err)
	}

	// ✅ แปลงข้อมูลจาก Smart Contract เป็น Struct
	result := &ProductLotInfo{
		LotID:                  string(bytes.Trim(productLotData.LotId[:], "\x00")),
		ProductID:              string(bytes.Trim(productLotData.ProductId[:], "\x00")),
		Factory:                productLotData.Factory.Hex(),
		Inspector:              productLotData.Inspector,
		InspectionDate:         time.Unix(productLotData.InspectionDate.Int64(), 0),
		Grade:                  productLotData.Grade,
		QualityAndNutritionCID: productLotData.QualityAndNutritionCID,
		MilkTankIDs:            convertBytes32ArrayToStrings(productLotData.MilkTankIds),
	}

	fmt.Println("✅ Product Lot Found:", result)
	return result, nil
}

// ✅ ดึง Product Lots ทั้งหมดของโรงงาน
func (b *BlockchainService) GetProductLotsByFactory(factoryAddress string) ([]map[string]string, error) {
	fmt.Println("📌 Fetching Product Lots for Factory:", factoryAddress)

	// ✅ แปลงที่อยู่โรงงานเป็น Address
	factoryAddr := common.HexToAddress(factoryAddress)

	// ✅ เรียก Smart Contract เพื่อดึง Lot IDs ทั้งหมดของโรงงาน
	lotIds, err := b.productLotContract.GetProductLotsByFactory(nil, factoryAddr)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch Product Lots: %v", err)
	}

	// ✅ แปลง `bytes32[]` เป็น `[]string`
	lotIdStrings := convertBytes32ArrayToStrings(lotIds)

	// ✅ เตรียมผลลัพธ์
	var productLots []map[string]string

	// ✅ ดึงข้อมูลของแต่ละ Product Lot
	for _, lotId := range lotIdStrings {
		// ✅ ดึงข้อมูล Product Lot จาก Blockchain
		productLotData, err := b.GetProductLotByLotID(lotId)
		if err != nil {
			fmt.Println("❌ Failed to fetch Product Lot:", lotId, err)
			continue // ข้ามอันที่ดึงไม่ได้
		}

		// ✅ ดึงข้อมูล Product Name จาก Smart Contract
		productID := productLotData.ProductID
		productData, err := b.GetProductDetails(productID)
		if err != nil {
			fmt.Println("❌ Failed to fetch Product Name for Product ID:", productID, err)
			continue // ข้ามอันที่ดึงไม่ได้
		}

		// ✅ แปลงสถานะจาก `uint8` เป็น `string`
		statusStr := strconv.Itoa(int(productLotData.Status))

		// ✅ เพิ่มข้อมูลเข้าไปในผลลัพธ์
		productLots = append(productLots, map[string]string{
			"Product Lot No":   lotId,
			"Product Name":     productData["productName"].(string),
			"Person In Charge": productLotData.Inspector, // ✅ ดึงชื่อ Inspector
			"Status":           statusStr,                // ✅ เพิ่มสถานะของ Product Lot
		})
	}

	fmt.Println("✅ Product Lots Fetched Successfully:", productLots)
	return productLots, nil
}

// ✅ Helper Function: แปลง `bytes32[]` เป็น `[]string`
func convertBytes32ArrayToStrings(arr [][32]byte) []string {
	var result []string
	for _, item := range arr {
		result = append(result, string(bytes.Trim(item[:], "\x00"))) // ลบ NULL Bytes
	}
	return result
}

// //////////////////////////////////////////////////////////// Tracking Event /////////////////////////////////////////////////////////
// CreateTrackingEvent - สร้างแทรคกิ้งอีเว้นต์
func StringToBytes32(s string) [32]byte {
	var b [32]byte
	copy(b[:], s)
	return b
}

func (b *BlockchainService) CreateTrackingEvent(
	userWallet string,
	trackingId string,
	productLotId string,
	retailerId string,
	qrCodeCID string,
) (string, error) {

	fmt.Println("📌 Creating Tracking Event on Blockchain for:", userWallet)

	// ✅ ดึง Private Key ของโรงงาน
	privateKeyHex, err := b.getPrivateKeyForAddress(userWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	// ✅ สร้าง Auth
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(userWallet)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(20000000000)

	// ✅ กลับไปใช้ BytesToHash (แบบเดิม)
	trackingIdBytes := StringToBytes32(trackingId)
	productLotIdBytes := StringToBytes32(productLotId)

	// ✅ Debug Log
	fmt.Println("📌 Debug - Sending to Blockchain:")
	fmt.Println("   - Tracking ID (Bytes32):", trackingIdBytes)
	fmt.Println("   - Product Lot ID (Bytes32):", productLotIdBytes)
	fmt.Println("   - Retailer ID:", retailerId)
	fmt.Println("   - QR Code CID:", qrCodeCID)

	// ✅ ส่งธุรกรรม
	tx, err := b.trackingContract.CreateTrackingEvent(
		auth,
		trackingIdBytes,
		productLotIdBytes,
		retailerId,
		qrCodeCID,
	)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create tracking event on blockchain: %v", err)
	}

	fmt.Println("✅ Transaction Sent:", tx.Hash().Hex())

	// ✅ รอ Confirm
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		return "", fmt.Errorf("❌ Transaction not mined: %v", err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		return "", errors.New("❌ Transaction failed")
	}

	fmt.Println("✅ Tracking Event Created on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (b *BlockchainService) GetTrackingByLotId(productLotId string) ([]string, []string, []string, error) {
	fmt.Println("📌 Fetching Tracking Events for Product Lot ID:", productLotId)

	// ✅ แปลง `productLotId` เป็น `bytes32`
	productLotIdBytes := StringToBytes32(productLotId)
	fmt.Println("✅ Converted ProductLotId to Bytes32:", productLotIdBytes)

	fmt.Println("📡 Calling Smart Contract...")
	result, err := b.trackingContract.GetTrackingByLotId(nil, productLotIdBytes)
	fmt.Println("✅ Smart Contract Call Completed!") // ❌ ถ้าไม่แสดงผล = Smart Contract มีปัญหา

	if err != nil {
		fmt.Println("❌ Failed to fetch tracking events:", err)
		return nil, nil, nil, fmt.Errorf("❌ Failed to fetch tracking events: %v", err)
	}
	// ✅ แปลง `[][32]byte` เป็น `[]string` โดยใช้ฟังก์ชันที่คุณให้มา
	trackingIds := convertBytes32ArrayToStrings(result.ResultTrackingIds)
	fmt.Println("✅ Smart Contract Returned Data:", result)
	return trackingIds, result.RetailerIds, result.QrCodeCIDs, nil
}

type TrackingResponse struct {
	TrackingId             string `json:"trackingId"`
	Status                 int    `json:"status"`
	ProductLotId           string `json:"productLotId"`
	PersonInChargePrevious string `json:"personInChargePrevious"`
	WalletAddressPrevious  string `json:"walletAddressPrevious"`
	SameLogistics          bool   `json:"sameLogistics"`
}

func (b *BlockchainService) GetAllTrackingIds(currentWallet string) ([]TrackingResponse, error) {
	fmt.Println("📌 Fetching All Tracking Events...")

	// ✅ Get all Tracking IDs
	trackingIds, err := b.trackingContract.GetAllTrackingIds(nil)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to fetch tracking events: %v", err)
	}

	trackingIdStrings := convertBytes32ArrayToStrings(trackingIds)
	var trackingList []TrackingResponse

	for _, trackingId := range trackingIdStrings {
		fmt.Println("📌 Processing Tracking ID:", trackingId)

		// ✅ ดึง ProductLotId
		productLotId, err := b.GetProductLotByTrackingId(trackingId)
		if err != nil {
			fmt.Println("❌ Failed to fetch Product Lot ID:", err)
			continue
		}
		fmt.Println("✅ Clean ProductLotId:", productLotId)

		// ✅ ดึง Tracking Event (ใช้ StringToBytes32)
		trackingEvent, err := b.trackingContract.TrackingEvents(nil, StringToBytes32(trackingId))
		if err != nil {
			fmt.Println("❌ Failed to fetch Tracking Event:", err)
			continue
		}
		status := int(trackingEvent.Status)

		personInCharge := ""
		walletAddress := ""
		sameLogistics := false

		if status == 0 {
			productLotIdBytes := StringToBytes32(productLotId)
			productLotDetails, err := b.productLotContract.GetProductLot(nil, productLotIdBytes)
			if err != nil {
				fmt.Println("❌ Failed to fetch Product Lot Details:", err)
			} else {
				personInCharge = productLotDetails.Inspector
				fmt.Println("✅ Inspector Name:", productLotDetails.Inspector)
			}
		}

		if status == 1 {
			checkpoints, err := b.trackingContract.GetLogisticsCheckpointsByTrackingId(nil, StringToBytes32(trackingId))
			if err != nil {
				fmt.Println("❌ Failed to fetch Checkpoints:", err)
			} else if len(checkpoints.AfterCheckpoints) > 0 {
				latest := checkpoints.AfterCheckpoints[len(checkpoints.AfterCheckpoints)-1]
				personInCharge = latest.PersonInCharge
				walletAddress = latest.LogisticsProvider.Hex()

				if walletAddress == currentWallet {
					sameLogistics = true
				}
			}
		}

		trackingList = append(trackingList, TrackingResponse{
			TrackingId:             trackingId,
			Status:                 status,
			ProductLotId:           productLotId,
			PersonInChargePrevious: personInCharge,
			WalletAddressPrevious:  walletAddress,
			SameLogistics:          sameLogistics,
		})
	}

	fmt.Println("✅ All Tracking Events Processed:", trackingList)
	return trackingList, nil
}

func (b *BlockchainService) UpdateLogisticsCheckpoint(
	logisticsWallet string,
	trackingId string,
	pickupTime uint64,
	deliveryTime uint64,
	quantity uint64,
	temperature int64,
	personInCharge string,
	checkType uint8, // ✅ 0 = Before, 1 = During, 2 = After
	receiverCID string, // ✅ บันทึกข้อมูลผู้รับสินค้า (IPFS CID)
) (string, error) {

	fmt.Println("📌 Updating Logistics Checkpoint on Blockchain for:", logisticsWallet)

	// ✅ ตรวจสอบค่าก่อนส่งธุรกรรม
	if logisticsWallet == "" || trackingId == "" || personInCharge == "" || receiverCID == "" {
		return "", fmt.Errorf("❌ Missing required fields")
	}

	// ✅ ดึง Private Key ของ Wallet ของโลจิสติกส์
	privateKeyHex, err := b.getPrivateKeyForAddress(logisticsWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	// ✅ สร้าง Transaction Auth โดยใช้ Private Key ของโลจิสติกส์
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(logisticsWallet)
	auth.GasLimit = uint64(3000000)         // ✅ กำหนด Gas Limit
	auth.GasPrice = big.NewInt(20000000000) // ✅ กำหนด Gas Price

	// ✅ แปลง `trackingId` เป็น `bytes32`
	trackingIdBytes := common.BytesToHash([]byte(trackingId))

	// ✅ Debug Log ก่อนส่งไปยัง Blockchain
	fmt.Println("📌 Debug - Sending to Blockchain:")
	fmt.Println("   - Tracking ID (Bytes32):", trackingIdBytes)
	fmt.Println("   - Pickup Time:", pickupTime)
	fmt.Println("   - Delivery Time:", deliveryTime)
	fmt.Println("   - Quantity:", quantity)
	fmt.Println("   - Temperature:", temperature)
	fmt.Println("   - Check Type:", checkType)
	fmt.Println("   - Receiver CID:", receiverCID)

	// ✅ ส่งธุรกรรมไปที่ Smart Contract
	tx, err := b.trackingContract.UpdateLogisticsCheckpoint(
		auth,
		trackingIdBytes, // ✅ ใช้ `bytes32`
		big.NewInt(int64(pickupTime)),
		big.NewInt(int64(deliveryTime)),
		big.NewInt(int64(quantity)),
		big.NewInt(temperature),
		personInCharge,
		uint8(checkType), // ✅ แปลง `enum` เป็น `uint8`
		receiverCID,      // ✅ บันทึกข้อมูลผู้รับสินค้า (IPFS CID)
	)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to update logistics checkpoint on blockchain: %v", err)
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

	fmt.Println("✅ Logistics Checkpoint Updated on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// ✅ Struct LogisticsCheckpoint ใน Go
type LogisticsCheckpoint struct {
	TrackingId        string `json:"trackingId"`
	LogisticsProvider string `json:"logisticsProvider"`
	PickupTime        uint64 `json:"pickupTime"`
	DeliveryTime      uint64 `json:"deliveryTime"`
	Quantity          uint64 `json:"quantity"`
	Temperature       int64  `json:"temperature"`
	PersonInCharge    string `json:"personInCharge"`
	CheckType         uint8  `json:"checkType"`
	ReceiverCID       string `json:"receiverCID"`
}

// ✅ ฟังก์ชันดึง Logistics Checkpoints จาก Blockchain
func (b *BlockchainService) GetLogisticsCheckpointsByTrackingId(trackingId string) ([]LogisticsCheckpoint, []LogisticsCheckpoint, []LogisticsCheckpoint, error) {
	fmt.Println("📌 Fetching Logistics Checkpoints for Tracking ID:", trackingId)

	// ✅ แปลง `trackingId` เป็น `bytes32`
	trackingIdBytes := common.BytesToHash([]byte(trackingId))
	fmt.Println("🛠 Debug - Tracking ID Before Query:", trackingId)
	fmt.Println("🛠 Debug - Tracking ID as Bytes32:", trackingIdBytes.Hex())

	// ✅ สร้างตัวแปรรับค่าผลลัพธ์จาก Smart Contract
	result, err := b.trackingContract.GetLogisticsCheckpointsByTrackingId(nil, trackingIdBytes)
	if err != nil {
		fmt.Println("❌ Failed to fetch logistics checkpoints:", err)
		return nil, nil, nil, fmt.Errorf("❌ Failed to fetch logistics checkpoints: %v", err)
	}

	// ✅ แปลงข้อมูลจาก Smart Contract เป็น Struct ของ Go
	beforeCheckpoints := convertToLogisticsCheckpointArray(result.BeforeCheckpoints)
	duringCheckpoints := convertToLogisticsCheckpointArray(result.DuringCheckpoints)
	afterCheckpoints := convertToLogisticsCheckpointArray(result.AfterCheckpoints)

	fmt.Println("✅ Logistics Checkpoints Retrieved Successfully")
	return beforeCheckpoints, duringCheckpoints, afterCheckpoints, nil
}

// ✅ แก้ไขฟังก์ชันให้รองรับ `tracking.TrackingLogisticsCheckpoint`
func convertToLogisticsCheckpointArray(data []tracking.TrackingLogisticsCheckpoint) []LogisticsCheckpoint {
	var checkpoints []LogisticsCheckpoint
	for _, d := range data {
		checkpoints = append(checkpoints, LogisticsCheckpoint{
			TrackingId:        bytes32ToString(d.TrackingId),
			LogisticsProvider: d.LogisticsProvider.Hex(),
			PickupTime:        d.PickupTime.Uint64(),
			DeliveryTime:      d.DeliveryTime.Uint64(),
			Quantity:          d.Quantity.Uint64(),
			Temperature:       d.Temperature.Int64(),
			PersonInCharge:    d.PersonInCharge,
			CheckType:         d.CheckType,
			ReceiverCID:       d.ReceiverCID,
		})
	}
	return checkpoints
}

// ✅ ฟังก์ชันช่วยแปลง bytes32 เป็น string
func bytes32ToString(data [32]byte) string {
	return strings.TrimRight(string(data[:]), "\x00")
}

func (b *BlockchainService) GetTrackingByRetailer(retailerID string) ([]map[string]interface{}, error) {
	fmt.Println("📌 Fetching tracking events for retailer:", retailerID)

	// ✅ ดึงรายการ Tracking IDs จาก Smart Contract
	trackingIDs, err := b.trackingContract.GetTrackingByRetailer(&bind.CallOpts{}, retailerID)
	if err != nil {
		fmt.Println("❌ Failed to fetch tracking events for retailer:", err)
		return nil, err
	}

	var trackingEvents []map[string]interface{}

	// ✅ วนลูปดึงข้อมูลของแต่ละ Tracking ID
	for _, id := range trackingIDs {
		trackingIDStr := string(bytes.Trim(id[:], "\x00"))

		// ✅ ดึงเฉพาะข้อมูล RetailerConfirmation
		_, _, retailerConfirmation, err := b.trackingContract.GetTrackingById(&bind.CallOpts{}, id)
		if err != nil {
			fmt.Println("❌ Failed to fetch tracking details for:", trackingIDStr, err)
			continue
		}

		// ✅ สร้าง JSON Response
		eventData := map[string]interface{}{
			"trackingId": trackingIDStr,
			"retailer":   retailerConfirmation,
		}

		trackingEvents = append(trackingEvents, eventData)
	}

	fmt.Println("✅ Fetched all tracking events for retailer:", retailerID, trackingEvents)
	return trackingEvents, nil
}

func (b *BlockchainService) RetailerReceiveProduct(
	userWallet string,
	trackingId string,
	retailerId string,
	qualityCID string,
	personInCharge string,
) (string, error) {

	fmt.Println("📌 Retailer Receiving Product on Blockchain for:", userWallet)

	// ✅ ดึง Private Key ของ Retailer
	privateKeyHex, err := b.getPrivateKeyForAddress(userWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	// ✅ สร้าง Transaction Auth โดยใช้ Private Key ของ Retailer
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(userWallet)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = big.NewInt(20000000000)

	// ✅ แปลง `trackingId` เป็น `bytes32`
	trackingIdBytes := common.BytesToHash([]byte(trackingId))

	// ✅ Debug Log ก่อนส่งไปยัง Blockchain
	fmt.Println("📌 Debug - Sending to Blockchain:")
	fmt.Println("   - Tracking ID (Bytes32):", trackingIdBytes)
	fmt.Println("   - Retailer ID:", retailerId)
	fmt.Println("   - Quality CID:", qualityCID)
	fmt.Println("   - Person In Charge:", personInCharge)

	// ✅ ส่งธุรกรรมไปที่ Smart Contract
	tx, err := b.trackingContract.RetailerReceiveProduct(
		auth,
		trackingIdBytes,
		retailerId,
		qualityCID,
		personInCharge,
	)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to execute retailerReceiveProduct on blockchain: %v", err)
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

	fmt.Println("✅ Retailer Received Product on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (b *BlockchainService) GetRetailerConfirmation(trackingId string) (map[string]interface{}, error) {
	fmt.Println("📌 Fetching Retailer Confirmation for Tracking ID:", trackingId)

	// ✅ แปลง Tracking ID เป็น bytes32
	trackingIdBytes := common.BytesToHash([]byte(trackingId))

	// ✅ ดึงเฉพาะข้อมูล Retailer Confirmation จาก Smart Contract
	_, _, retailerConfirmation, err := b.trackingContract.GetTrackingById(nil, trackingIdBytes)
	if err != nil {
		fmt.Println("❌ Failed to fetch retailer confirmation:", err)
		return nil, fmt.Errorf("Failed to fetch retailer confirmation: %v", err)
	}

	// ✅ ตรวจสอบว่า Retailer Confirmation มีอยู่จริง
	if retailerConfirmation.TrackingId == [32]byte{} {
		return nil, fmt.Errorf("No retailer confirmation found for Tracking ID: %s", trackingId)
	}

	// ✅ แปลงข้อมูลที่ได้เป็น Map (JSON-compatible)
	retailerData := map[string]interface{}{
		"trackingId":     trackingId,
		"retailerId":     retailerConfirmation.RetailerId,
		"receivedTime":   retailerConfirmation.ReceivedTime,
		"qualityCID":     retailerConfirmation.QualityCID,
		"personInCharge": retailerConfirmation.PersonInCharge,
	}

	fmt.Println("✅ Retailer Confirmation Data:", retailerData)
	return retailerData, nil
}

func Bytes32ToString(b [32]byte) string {
	n := bytes.IndexByte(b[:], 0) // หาตำแหน่ง Null byte แรก
	if n == -1 {
		n = len(b)
	}
	return string(b[:n])
}
func (b *BlockchainService) GetProductLotByTrackingId(trackingId string) (string, error) {
	fmt.Println("📌 Fetching Product Lot for Tracking ID:", trackingId)

	// ✅ ใช้ StringToBytes32
	trackingIdBytes := StringToBytes32(trackingId)

	// ✅ ดึงข้อมูลจาก Smart Contract
	trackingEvent, _, _, err := b.trackingContract.GetTrackingById(nil, trackingIdBytes)
	if err != nil {
		fmt.Println("❌ Failed to fetch product lot by tracking ID:", err)
		return "", fmt.Errorf("Failed to fetch product lot by tracking ID: %v", err)
	}

	// ✅ ตรวจสอบ
	if trackingEvent.TrackingId == [32]byte{} {
		fmt.Println("⚠️ No tracking event found for Tracking ID:", trackingId)
		return "", fmt.Errorf("No tracking event found for Tracking ID: %s", trackingId)
	}

	// ✅ Clean Null Bytes
	productLotId := Bytes32ToString(trackingEvent.ProductLotId)
	fmt.Println("✅ Clean Product Lot ID:", productLotId)

	return productLotId, nil
}

// /ฟังชั่นนี้ยังใช้ไม่ได้เพายังไม่อัปเดต
func (b *BlockchainService) GetOngoingShipmentsByLogistics(walletAddress string) ([]map[string]interface{}, error) {
	fmt.Println("📌 Fetching Ongoing Shipments for Logistics Wallet:", walletAddress)

	// ✅ Call Smart Contract function
	result, err := b.trackingContract.GetOngoingShipmentsByLogistics(nil)
	if err != nil {
		fmt.Println("❌ Failed to fetch ongoing shipments:", err)
		return nil, fmt.Errorf("❌ Failed to fetch ongoing shipments: %v", err)
	}

	trackingIds := result.TrackingIds
	personInChargeList := result.PersonInChargeList

	trackingIdStrings := convertBytes32ArrayToStrings(trackingIds)

	var shipmentList []map[string]interface{}

	for i, trackingId := range trackingIdStrings {
		shipmentList = append(shipmentList, map[string]interface{}{
			"trackingId":     trackingId,
			"personInCharge": personInChargeList[i],
			"walletAddress":  walletAddress, // ✅ เพิ่ม Wallet Address เดียวกับ msg.sender
		})
	}

	fmt.Println("✅ Ongoing Shipments Retrieved:", shipmentList)
	return shipmentList, nil
}

func (b *BlockchainService) GetLastLogisticsProvider(trackingID string) (string, error) {
	fmt.Println("📌 Fetching Last Logistics Provider for TrackingID:", trackingID)

	trackingIdBytes := common.HexToHash(trackingID)
	providerAddress, err := b.trackingContract.GetLastLogisticsProvider(nil, trackingIdBytes)
	if err != nil {
		fmt.Println("❌ Failed to fetch last logistics provider:", err)
		return "", err
	}

	return providerAddress.Hex(), nil
}

func (b *BlockchainService) GetRetailerInTransitTracking(retailerID string) ([]map[string]interface{}, error) {
	fmt.Println("📌 Fetching InTransit Tracking Events for retailer:", retailerID)

	trackingIDs, err := b.trackingContract.GetTrackingByRetailer(nil, retailerID)
	if err != nil {
		return nil, err
	}

	var trackingList []map[string]interface{}

	for _, id := range trackingIDs {
		trackingIDStr := string(bytes.Trim(id[:], "\x00"))

		trackingEvent, checkpoints, _, err := b.trackingContract.GetTrackingById(nil, id)
		if err != nil {
			continue
		}

		if int(trackingEvent.Status) != 1 { // ต้อง InTransit เท่านั้น
			continue
		}

		personInCharge := ""
		if len(checkpoints) > 0 {
			for i := len(checkpoints) - 1; i >= 0; i-- {
				if checkpoints[i].CheckType == 2 { // After
					personInCharge = checkpoints[i].PersonInCharge
					break
				}
			}
		}

		// ✅ ดึง Last Logistics Provider
		lastProvider, err := b.GetLastLogisticsProvider(trackingIDStr)
		if err != nil {
			lastProvider = "Unknown"
		}

		trackingList = append(trackingList, map[string]interface{}{
			"trackingId":          trackingIDStr,
			"personInCharge":      personInCharge,
			"lastLogisticsWallet": lastProvider,
			"status":              1,
		})
	}

	return trackingList, nil
}
