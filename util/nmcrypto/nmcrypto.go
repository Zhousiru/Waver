package nmcrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"

	"github.com/Zhousiru/Waver/util/tool"
)

const (
	iv      = "0102030405060708"
	modulus = "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7"
	nonce   = "0CoJUm6Qyw8W8jud"
	pubKey  = "010001"
)

func aesEncrypt(data, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	padding := blockSize - len(data)%blockSize
	padText := string(bytes.Repeat([]byte{byte(padding)}, padding))
	data += padText

	encrypter := cipher.NewCBCEncrypter(block, []byte(iv))
	ciphertext := make([]byte, len(data))
	encrypter.CryptBlocks(ciphertext, []byte(data))

	result := base64.StdEncoding.EncodeToString(ciphertext)

	return result, nil
}

func rsaEncrypt(data string) string {
	rData := []byte(data)
	for i, j := 0, len(rData)-1; i < j; i, j = i+1, j-1 {
		rData[i], rData[j] = rData[j], rData[i]
	}

	rData10 := new(big.Int)
	pubKey10 := new(big.Int)
	modulus10 := new(big.Int)

	rData10.SetString(hex.EncodeToString(rData), 16)
	pubKey10.SetString(pubKey, 16)
	modulus10.SetString(modulus, 16)

	pow := rData10.Exp(rData10, pubKey10, nil)
	ciphertext := fmt.Sprintf("%x", pow.Mod(pow, modulus10))
	ciphertextLen := len(ciphertext)
	if ciphertextLen < 256 {
		padding := 256 - ciphertextLen
		padText := string(bytes.Repeat([]byte("0"), padding))
		ciphertext = padText + ciphertext
	}

	return ciphertext
}

// EncryptRequest - 加密网易云音乐 POST 参数
func EncryptRequest(data map[string]string) (string, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	key := tool.GetRandomStr(16)

	firstEncrypt, err := aesEncrypt(string(dataJSON), nonce)
	if err != nil {
		return "", err
	}
	params, err := aesEncrypt(firstEncrypt, key)
	if err != nil {
		return "", err
	}
	encSecKey := rsaEncrypt(key)

	postData := url.Values{}
	postData.Add("params", params)
	postData.Add("encSecKey", encSecKey)

	return postData.Encode(), nil
}
