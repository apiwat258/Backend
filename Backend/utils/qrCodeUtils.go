package utils

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/common"
)

// GenerateHash - แปลง 16-char ID เป็น bytes32 Hash (สำหรับ Blockchain)
func GenerateHash(rawMilkID string) common.Hash {
	hash := sha256.Sum256([]byte(rawMilkID))
	return common.BytesToHash(hash[:])
}
