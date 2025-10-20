package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
)
/*
原始 AsiaLife 代收(代退)編碼資料：
AsiaLife://D1=H001&D2=1&D3=3961&D4=H00101&D5=54043913731478&D
6=682101&D7=20210204235959&D8=1&D9=&D10=1&D11=1&D12=Howar
dYou&D13=A123456789&D14=0988057404&D15=2
l AsiaLife 代收編碼資料加密流程：
1. 取得特約夥伴秘密金鑰 [Key](共計 32 碼)和 [iV](共計 16 碼)，請洽詢本公司業務
窗口取得金鑰，金鑰範例如下：
金鑰[key]：1234567890123456789012345678H001
金鑰[iv]：123456789012H001
2. 使用[key]、[iv] 金鑰將 AsiaLife 代收編碼資料進行 AES-128-CBC 加密
AES-128-CBC 加密後範例如下：
zhGkfgQmDWhIxTTbTt3HfJTBpfSO5EvTzumoQtwanAlSr0XoGBt6Dn3phcJ9t
djgltQ1Ksj2qBmdwTsvco6kG24oQ42MNa03s449uTeCBU0q7s24UdDsVZGXjv
gckAwnPOEEbLFOOz2wyKI0D2SjI5FdiPJQz648ur0NmUj4LQjXf41Ae7kyJy9AY
dZK4kXbfapwT0Bw7S+lfM+MjomBPg==
*/

func main() {
	// 測試內容
	plaintext := `AsiaLife://D1=H001&D2=1&D3=3961&D4=H00101&D5=54043913731478&D
6=682101&D7=20210204235959&D8=1&D9=&D10=1&D11=1&D12=Howar
dYou&D13=A123456789&D14=0988057404&D15=2`
	
	// 金鑰和 IV（從字串轉成 byte array）
	key := []byte("1234567890123456789012345678H001")[:16] // 取前16位作為 AES-128 金鑰
	iv := []byte("123456789012H001")[:16]                  // 取前16位作為 IV

	// 加密
	ciphertext, err := encryptAES128CBC(plaintext, key, iv)
	if err != nil {
		log.Fatal("加密失敗:", err)
	}

	// 輸出結果
	fmt.Printf("原始文字: %s\n", plaintext)
	fmt.Printf("Hex 編碼: %s\n", hex.EncodeToString(ciphertext))
	fmt.Printf("Base64 編碼: %s\n", base64.StdEncoding.EncodeToString(ciphertext))

	// 解密測試
	decrypted, err := decryptAES128CBC(ciphertext, key, iv)
	if err != nil {
		log.Fatal("解密失敗:", err)
	}
	fmt.Printf("解密結果: %s\n", decrypted)
}

// AES-128-CBC 加密
func encryptAES128CBC(plaintext string, key, iv []byte) ([]byte, error) {
	// 創建 AES 密碼塊
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// PKCS7 填充
	plaintextBytes := []byte(plaintext)
	plaintextBytes = pkcs7Pad(plaintextBytes, aes.BlockSize)

	// 創建 CBC 模式加密器
	ciphertext := make([]byte, len(plaintextBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintextBytes)

	return ciphertext, nil
}

// AES-128-CBC 解密
func decryptAES128CBC(ciphertext, key, iv []byte) (string, error) {
	// 創建 AES 密碼塊
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 創建 CBC 模式解密器
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// PKCS7 去除填充
	plaintext = pkcs7Unpad(plaintext)

	return string(plaintext), nil
}

// PKCS7 填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7 去除填充
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}