package services

import (
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
type RawMilkData struct {
	TankId           string `json:"tankId"`
	FarmWallet       string `json:"farmWallet"`
	PersonInCharge   string `json:"personInCharge"`
	QualityReportCID string `json:"qualityReportCid"` //
	QrCodeCID        string `json:"qrCodeCid"`
	Status           uint8  `json:"status"`
}

func (b *BlockchainService) CreateMilkTank(
	userWallet string,
	tankId string,
	personInCharge string, // ✅ เพิ่มพารามิเตอร์นี้
	qrCodeCID string, // ✅ คงไว้
) (string, error) {

	fmt.Println("📌 Creating Milk Tank on Blockchain for:", userWallet)

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

	// ✅ แปลง tankId เป็น bytes32
	tankIdBytes := common.BytesToHash([]byte(tankId))

	// ✅ ส่งธุรกรรมไปที่ Smart Contract
	tx, err := b.rawMilkContract.CreateMilkTank(
		auth,
		tankIdBytes,
		personInCharge, // ✅ ส่ง personInCharge ไปที่ Smart Contract
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

func (b *BlockchainService) ValidateMilkData(
	quantity uint64, temperature uint64, pH uint64, fat uint64, protein uint64, bacteria bool, contaminants bool,
) (bool, string) {

	// ✅ ตรวจสอบค่าตามกฎที่กำหนด
	if temperature < 200 || temperature > 600 {
		return false, "Error: Temperature out of range! (2.0C - 6.0C)"
	}
	if pH < 650 || pH > 680 {
		return false, "Error: pH out of range! (6.5 - 6.8)"
	}
	if fat < 300 || fat > 400 {
		return false, "Error: Fat percentage out of range! (3.0% - 4.0%)"
	}
	if protein < 300 || protein > 350 {
		return false, "Error: Protein percentage out of range! (3.0% - 3.5%)"
	}

	// ✅ ถ้าผ่านเงื่อนไขทั้งหมด
	return true, "Validated successfully."
}

func (b *BlockchainService) GetAllRawMilkTanks() ([]map[string]string, error) {
	fmt.Println("📌 Fetching all milk tanks from Blockchain...")

	// ✅ เรียก Smart Contract เพื่อดึง tankIds ทั้งหมด
	tankIds, err := b.rawMilkContract.GetAllMilkTanks(&bind.CallOpts{})
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tanks:", err)
		return nil, fmt.Errorf("❌ Failed to fetch milk tanks: %v", err)
	}

	var milkTanks []map[string]string

	// ✅ วนลูปดึงข้อมูลแท็งก์แต่ละอัน
	for _, id := range tankIds {
		tankId := common.BytesToHash(id[:]).Hex()

		// ✅ ดึงข้อมูลแท็งก์จาก Smart Contract ตาม tankId
		tankIdSC, _, personInCharge, status, _, qrCodeCID, err :=
			b.rawMilkContract.GetMilkTank(&bind.CallOpts{}, id)
		if err != nil {
			fmt.Printf("❌ Failed to fetch details for tank %s: %v\n", tankId, err)
			continue
		}

		// ✅ เพิ่มเข้าไปในรายการ
		milkTanks = append(milkTanks, map[string]string{
			"tankId":         common.BytesToHash(tankIdSC[:]).Hex(),
			"personInCharge": personInCharge,
			"status":         fmt.Sprintf("%d", status), // แปลง enum เป็น string
			"qrCodeCID":      qrCodeCID,
		})
	}

	fmt.Println("✅ Retrieved Milk Tanks:", milkTanks)
	return milkTanks, nil
}

func (b *BlockchainService) GetRawMilkTankDetails(tankId string) (*RawMilkData, error) {
	fmt.Println("📌 Fetching milk tank details for:", tankId)

	// ✅ แปลง tankId เป็น bytes32
	tankIdBytes := common.HexToHash(tankId)

	// ✅ ดึงข้อมูลแท็งก์จาก Smart Contract (ต้องรับค่า 6 ตัว)
	tankIdSC, farmWallet, personInCharge, status, qualityReportCID, qrCodeCID, err :=
		b.rawMilkContract.GetMilkTank(&bind.CallOpts{}, tankIdBytes)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tank details:", err)
		return nil, fmt.Errorf("❌ Failed to fetch milk tank details: %v", err)
	}

	// ✅ แปลงค่า tankIdSC เป็น string
	tankIdStr := common.BytesToHash(tankIdSC[:]).Hex()

	// ✅ แปลงข้อมูลจาก Smart Contract เป็นโครงสร้างที่ใช้ใน Go
	rawMilk := &RawMilkData{
		TankId:           tankIdStr,
		FarmWallet:       farmWallet.Hex(),
		PersonInCharge:   personInCharge,
		QualityReportCID: qualityReportCID,
		QrCodeCID:        qrCodeCID,
		Status:           uint8(status),
	}

	fmt.Println("✅ Milk Tank Details Retrieved:", rawMilk)
	return rawMilk, nil
}

func (b *BlockchainService) GetMilkTanksByFarmer(farmerAddress string) ([]map[string]string, error) {
	fmt.Println("📌 Fetching milk tanks for farmer:", farmerAddress)

	farmer := common.HexToAddress(farmerAddress)

	// ✅ เรียก Smart Contract เพื่อนำ Tank IDs ของฟาร์มมา
	tankIDs, err := b.rawMilkContract.GetMilkTanksByFarmer(&bind.CallOpts{}, farmer)
	if err != nil {
		fmt.Println("❌ Failed to fetch milk tanks for farmer:", err)
		return nil, err
	}

	var milkTanks []map[string]string

	// ✅ วนลูปดึงข้อมูลแท็งก์แต่ละอัน
	for _, id := range tankIDs {
		tankId := common.BytesToHash(id[:]).Hex()

		// ✅ ดึงข้อมูลแท็งก์จาก Smart Contract
		tankIdSC, _, personInCharge, status, _, qrCodeCID, err :=
			b.rawMilkContract.GetMilkTank(&bind.CallOpts{}, id)
		if err != nil {
			fmt.Printf("❌ Failed to fetch details for tank %s: %v\n", tankId, err)
			continue
		}

		// ✅ เพิ่มเข้าไปในรายการ
		milkTanks = append(milkTanks, map[string]string{
			"tankId":         common.BytesToHash(tankIdSC[:]).Hex(),
			"personInCharge": personInCharge,
			"status":         fmt.Sprintf("%d", status), // แปลง enum เป็น string
			"qrCodeCID":      qrCodeCID,
		})
	}

	fmt.Println("✅ Fetched milk tanks for farmer:", farmerAddress, milkTanks)
	return milkTanks, nil
}

func (b *BlockchainService) VerifyMilkQuality(userWallet string, tankID string, approved bool, qualityReportCID string) (string, error) {
	fmt.Println("📌 Verifying milk quality for Tank:", tankID, "Approved:", approved)

	// ✅ ดึง Private Key ของ Factory (หรือ User ที่มีสิทธิ์)
	privateKeyHex, err := b.getPrivateKeyForAddress(userWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	// ✅ สร้าง Transaction Auth โดยใช้ Private Key
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(userWallet)

	// ✅ ใช้ common.BytesToHash() แทน common.HexToHash()
	tankIDBytes := common.BytesToHash([]byte(tankID))

	// ✅ ส่ง Transaction ไปยัง Smart Contract
	tx, err := b.rawMilkContract.VerifyMilkQuality(auth, tankIDBytes, approved, qualityReportCID)
	if err != nil {
		fmt.Println("❌ Failed to verify milk quality:", err)
		return "", err
	}

	fmt.Println("✅ Transaction Sent:", tx.Hash().Hex())

	// ✅ รอให้ Transaction ถูกบันทึก
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		fmt.Println("❌ Transaction not mined:", err)
		return "", err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		fmt.Println("❌ Transaction failed!")
		return "", errors.New("Transaction failed")
	}

	fmt.Println("✅ Milk quality verified on Blockchain. TX Hash:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (b *BlockchainService) UpdateMilkTankStatus(userWallet string, tankID string, approved bool) (string, error) {
	fmt.Println("📌 Updating milk tank status for Tank:", tankID, "Approved:", approved)

	// ✅ ดึง Private Key ของ Factory (หรือ User ที่มีสิทธิ์)
	privateKeyHex, err := b.getPrivateKeyForAddress(userWallet)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get private key: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to parse private key: %v", err)
	}

	// ✅ สร้าง Transaction Auth โดยใช้ Private Key
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, getChainID())
	if err != nil {
		return "", fmt.Errorf("❌ Failed to create transactor: %v", err)
	}
	auth.From = common.HexToAddress(userWallet)

	// ✅ ใช้ common.BytesToHash() แทน common.HexToHash()
	tankIDBytes := common.BytesToHash([]byte(tankID))

	// ✅ ใช้ VerifyMilkQuality แทน SetTankStatus
	qualityReportCID := "" // ✅ ถ้าไม่มีการเปลี่ยน Quality Report ให้ใช้ค่าว่าง
	tx, err := b.rawMilkContract.VerifyMilkQuality(auth, tankIDBytes, approved, qualityReportCID)
	if err != nil {
		fmt.Println("❌ Failed to update milk tank status:", err)
		return "", err
	}

	fmt.Println("✅ Transaction Sent:", tx.Hash().Hex())

	// ✅ รอให้ Transaction ถูกบันทึก
	receipt, err := bind.WaitMined(context.Background(), b.client, tx)
	if err != nil {
		fmt.Println("❌ Transaction not mined:", err)
		return "", err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		fmt.Println("❌ Transaction failed!")
		return "", errors.New("Transaction failed")
	}

	fmt.Println("✅ Milk tank status updated on Blockchain. TX Hash:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}
