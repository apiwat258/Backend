package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"finalyearproject/Backend/services/certification_event"
	"finalyearproject/Backend/models"
)

// CertificateService - สำหรับจัดการใบ Certification บน Blockchain
type CertificateService struct {
	client                *ethclient.Client
	auth                  *bind.TransactOpts
	certificationContract *certification_event.CertificationEvent
}

// ✅ Global Instance ของ Certification Service
var CertificateServiceInstance *CertificateService

// InitCertificateService - ใช้เชื่อมต่อกับ Blockchain และโหลด Smart Contract
func InitCertificateService(client *ethclient.Client, auth *bind.TransactOpts, contractAddr common.Address) error {
	certInstance, err := certification_event.NewCertificationEvent(contractAddr, client)
	if err != nil {
		return fmt.Errorf("❌ Failed to load certification contract: %v", err)
	}

	CertificateServiceInstance = &CertificateService{
		client:                client,
		auth:                  auth,
		certificationContract: certInstance,
	}

	fmt.Println("✅ Certification Service Initialized!")
	return nil
}

// StoreCertificationOnBlockchain - บันทึกใบเซอร์ลง Blockchain
func (cs *CertificateService) StoreCertificationOnBlockchain(eventID, entityType, entityID, certCID string, issuedDate, expiryDate *big.Int) (string, error) {
	tx, err := cs.certificationContract.StoreCertificationEvent(cs.auth, eventID, entityType, entityID, certCID, issuedDate, expiryDate)
	if err != nil {
		log.Println("❌ Failed to store certification event on blockchain:", err)
		return "", err
	}

	fmt.Println("✅ Certification Event stored on Blockchain:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// GetCertificationFromBlockchain - ดึงข้อมูลใบเซอร์จาก Blockchain
func (cs *CertificateService) GetCertificationFromBlockchain(eventID string) (*models.Certification, error) {
	// ✅ เรียก Smart Contract เพื่อนำข้อมูลมา
	certEvent, err := cs.certificationContract.GetCertificationEvent(&bind.CallOpts{}, eventID)
	if err != nil {
		log.Println("❌ [Blockchain] Failed to fetch certification:", err)
		return nil, err
	}

	// ✅ แปลงค่า timestamp เป็น `time.Time`
	issuedDate := time.Unix(certEvent.IssuedDate.Int64(), 0)
	expiryDate := time.Unix(certEvent.ExpiryDate.Int64(), 0)

	// ✅ คืนค่าเป็น Struct ที่ใช้ในระบบ
	return &models.Certification{
		CertificationID:   certEvent.EventID,
		EntityType:        certEvent.EntityType,
		EntityID:          certEvent.EntityID,
		CertificationCID:  certEvent.CertificationCID,
		IssuedDate:        issuedDate,
		EffectiveDate:     expiryDate,
		BlockchainTxHash:  "", // ไม่มีค่าจาก Smart Contract
		CreatedOn:         time.Unix(certEvent.CreatedOn.Int64(), 0),
	}, nil
}
