package util

import "github.com/gofrs/uuid"

// 生成随机UUID
func GenerateUUID() string {
	u1, _ := uuid.NewV4()
	return u1.String()
}
