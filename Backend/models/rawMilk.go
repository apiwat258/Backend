package models

type RawMilk struct {
	TankID          string       `json:"tankId" gorm:"primaryKey"` // รหัสแท็งก์นมดิบ
	FarmID          string       `json:"farmId"`                   // รหัสฟาร์ม
	PersonInCharge  string       `json:"personInCharge"`           // ผู้รับผิดชอบ
	Quantity        float64      `json:"quantity"`                 // ปริมาณนมดิบ (ลิตร)
	Temperature     float64      `json:"temperature"`              // อุณหภูมิขณะเก็บ
	PH              float64      `json:"pH"`                       // ค่า pH
	Fat             float64      `json:"fat"`                      // ไขมัน (%)
	Protein         float64      `json:"protein"`                  // โปรตีน (%)
	Status          string       `json:"status"`                   // Pending / Ready for Shipping
	Bacteria        bool         `json:"bacteria"`                 // ตรวจพบน้ำปนเปื้อนหรือไม่
	BacteriaInfo    string       `json:"bacteriaInfo"`             // รายละเอียดการปนเปื้อน (ถ้ามี)
	Contaminants    bool         `json:"contaminants"`             // มีสารปนเปื้อนหรือไม่
	ContaminantInfo string       `json:"contaminantInfo"`          // รายละเอียดสารปนเปื้อน (ถ้ามี)
	AbnormalChar    bool         `json:"abnormalChar"`             // มีคุณลักษณะผิดปกติหรือไม่
	AbnormalData    AbnormalType `json:"abnormalType"`             // ข้อมูลลักษณะผิดปกติของนมดิบ
	QualityCheck    QualityCheck `json:"qualityVerification"`      // การตรวจสอบคุณภาพ
	QRCodeCID       string       `json:"qrCodeCID"`                // CID ของ QR Code ที่เก็บข้อมูลแท็งก์
}

type AbnormalType struct {
	SmellBad      bool `json:"smellBad"`      // มีกลิ่นเสีย
	SmellNotFresh bool `json:"smellNotFresh"` // มีกลิ่นไม่สด
	AbnormalColor bool `json:"abnormalColor"` // สีผิดปกติ
	Sour          bool `json:"sour"`          // รสเปรี้ยว
	Bitter        bool `json:"bitter"`        // รสขม
	Cloudy        bool `json:"cloudy"`        // ขุ่น
	Lumpy         bool `json:"lumpy"`         // เป็นก้อน
	Separation    bool `json:"separation"`    // แยกชั้น
}

type QualityCheck struct {
	Approved        bool   `json:"approved"`         // ผ่านการตรวจสอบหรือไม่
	RejectionReason string `json:"rejectionReason"`  // เหตุผลที่ถูกปฏิเสธ (ถ้ามี)
	ReportCID       string `json:"qualityReportCID"` // CID ของรายงานใน IPFS
}
