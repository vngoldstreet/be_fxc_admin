package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// Map lưu trữ mã liên kết dựa trên ID
var idToCodeMap map[string]string
var idToContest map[int]string

// Độ dài của mã liên kết
const codeLength = 16
const codeLengthContest = 8

// GenerateCode tạo mã liên kết dựa trên ID và lưu trữ nó trong map
func GenerateCode(id string) string {
	code := GenerateSecureCode(id)
	idToCodeMap[id] = code
	return code
}

func GenerateCodeContest(id int) string {
	code := GenerateSecureCodeContest(id)
	idToContest[id] = code
	return code
}

// Hàm này tạo mã liên kết an toàn dựa trên ID sử dụng ngẫu nhiên và mã hóa base64
func GenerateSecureCode(id string) string {
	randomBytes := make([]byte, codeLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	code := base64.URLEncoding.EncodeToString(randomBytes)
	code = code[:codeLength]
	return code
}

func GenerateSecureCodeContest(id int) string {
	randomBytes := make([]byte, codeLengthContest)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println(err.Error())
	}
	code := base64.URLEncoding.EncodeToString(randomBytes)
	code = code[:codeLengthContest]
	return code
}
