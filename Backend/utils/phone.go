package utils

import "strings"

// ExtractAreaCodeAndPhone แยกรหัสพื้นที่และหมายเลขโทรศัพท์
func ExtractAreaCodeAndPhone(fullPhone string) (string, string) {
	parts := strings.SplitN(fullPhone, " ", 2) // แยกเป็น 2 ส่วน: [areaCode, phoneNumber]
	if len(parts) == 2 {
		return parts[0], parts[1] // ✅ คืนค่า areaCode และ phoneNumber
	}
	return "", fullPhone // ✅ ถ้าไม่มีรหัสพื้นที่ ให้คืนค่าเป็น "" และใช้เลขที่เหลือเป็น phoneNumber
}
