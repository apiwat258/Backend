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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	certification "finalyearproject/Backend/services/certification_event" // ✅ สำหรับ Raw Milk
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

	// ✅ ดึงรายการ Tank IDs จาก Smart Contract
	tankIDs, err := b.rawMilkContract.GetMilkTanksByFarmer(&bind.CallOpts{}, farmer)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tanks for farmer:", err)
		return nil, err
	}

	// ✅ เก็บข้อมูลนมดิบของฟาร์ม
	var milkTanks []map[string]interface{}

	// ✅ วนลูปดึงข้อมูลของแต่ละแท็งก์
	for _, id := range tankIDs {
		// ✅ แปลง `bytes32` เป็น `string`
		tankIdStr := string(bytes.Trim(id[:], "\x00")) // ✅ ดึงรายละเอียดแท็งก์จาก Smart Contract
		_, farmerAddr, factoryIdSC, personInCharge, status, qualityReportCID, qrCodeCID, err :=
			b.rawMilkContract.GetMilkTank(&bind.CallOpts{}, id)
		if err != nil {
			fmt.Printf("❌ Failed to fetch details for tank %s: %v\n", tankIdStr, err)
			continue
		}

		// ✅ แปลงค่าที่ได้จาก Smart Contract
		milkTank := map[string]interface{}{
			"tankId":           tankIdStr, // ✅ ใช้ค่า string ที่แปลงมา
			"farmerAddress":    farmerAddr.Hex(),
			"factoryId":        string(factoryIdSC[:]), // ✅ แปลง `factoryId` ด้วย
			"personInCharge":   personInCharge,
			"status":           uint8(status), // ✅ Convert Enum เป็น Number
			"qualityReportCID": qualityReportCID,
			"qrCodeCID":        qrCodeCID,
		}

		// ✅ เพิ่มเข้าไปในรายการ
		milkTanks = append(milkTanks, milkTank)
	}

	fmt.Println("✅ Fetched milk tanks for farmer:", farmerAddress, milkTanks)
	return milkTanks, nil
}

func (b *BlockchainService) GetRawMilkTankDetails(tankId string) (*RawMilkData, error) {
	fmt.Println("📌 Fetching milk tank details for:", tankId)

	// ✅ แปลง tankId เป็น bytes32 (ต้องใช้ Hash เพราะตอน Create ใช้ Hash)
	tankIdBytes := common.BytesToHash([]byte(tankId))

	// ✅ ดึงข้อมูลแท็งก์จาก Smart Contract
	tankIdSC, farmWallet, factoryIdSC, personInCharge, status, qualityReportCID, qrCodeCID, err :=
		b.rawMilkContract.GetMilkTank(&bind.CallOpts{}, tankIdBytes)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tank details:", err)
		return nil, fmt.Errorf("❌ Failed to fetch milk tank details: %v", err)
	}

	// ✅ แปลงค่า tankIdSC และ factoryIdSC เป็น string
	tankIdStr := string(tankIdSC[:])       // ✅ คืนค่า Tank ID เป็น String
	factoryIdStr := string(factoryIdSC[:]) // ✅ คืนค่า Factory ID เป็น String

	// ✅ แปลงข้อมูลจาก Smart Contract เป็นโครงสร้างที่ใช้ใน Go
	rawMilk := &RawMilkData{
		TankId:           tankIdStr,
		FarmWallet:       farmWallet.Hex(),
		FactoryId:        factoryIdStr,
		PersonInCharge:   personInCharge,
		QualityReportCID: qualityReportCID,
		QrCodeCID:        qrCodeCID,
		Status:           uint8(status), // ✅ แปลง Enum เป็น Number
	}

	fmt.Println("✅ Milk Tank Details Retrieved:", rawMilk)
	return rawMilk, nil
}

func (b *BlockchainService) GetMilkTanksByFactory(factoryID string) ([]map[string]interface{}, error) {
	fmt.Println("📌 Fetching milk tanks for factory:", factoryID)

	factoryIDBytes32 := common.BytesToHash([]byte(factoryID))
	fmt.Println("🔍 [Fixed] Converted FactoryID to Bytes32:", factoryIDBytes32)

	// ✅ ดึงรายการ Tank IDs จาก Smart Contract
	tankIDs, err := b.rawMilkContract.GetMilkTanksByFactory(&bind.CallOpts{}, factoryIDBytes32)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tanks for factory:", err)
		return nil, err
	}

	// ✅ เก็บข้อมูลนมดิบของโรงงาน
	var milkTanks []map[string]interface{}

	// ✅ วนลูปดึงข้อมูลของแต่ละแท็งก์
	for _, id := range tankIDs {
		// ✅ แปลง `bytes32` เป็น `string`
		tankIdStr := string(bytes.Trim(id[:], "\x00"))

		// ✅ ดึงรายละเอียดแท็งก์จาก Smart Contract
		_, farmerAddr, factoryIdSC, personInCharge, status, qualityReportCID, qrCodeCID, err :=
			b.rawMilkContract.GetMilkTank(&bind.CallOpts{}, id)
		if err != nil {
			fmt.Printf("❌ Failed to fetch details for tank %s: %v\n", tankIdStr, err)
			continue
		}

		// ✅ แปลงค่าที่ได้จาก Smart Contract
		milkTank := map[string]interface{}{
			"tankId":           tankIdStr, // ✅ ใช้ค่า string ที่แปลงมา
			"farmerAddress":    farmerAddr.Hex(),
			"factoryId":        string(factoryIdSC[:]), // ✅ แปลง `factoryId` ด้วย
			"personInCharge":   personInCharge,
			"status":           uint8(status), // ✅ Convert Enum เป็น Number
			"qualityReportCID": qualityReportCID,
			"qrCodeCID":        qrCodeCID,
		}

		// ✅ เพิ่มเข้าไปในรายการ
		milkTanks = append(milkTanks, milkTank)
	}

	fmt.Println("✅ Fetched milk tanks for factory:", factoryID, milkTanks)
	return milkTanks, nil
}
