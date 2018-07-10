package BLC

import (
	"math/big"
	"bytes"
)

// base64比base58多6个字符
//ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/
//0(零)，O(大写的 o)，I(大写的i)，l(小写的 L)，+，/

var yxh_b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// 字节数组转 Base58,加密
func Yxh_Base58Encode(input []byte) []byte {

	var result []byte

	//初始化除数
	x := big.NewInt(0).SetBytes(input)

	//初始化被除数
	base := big.NewInt(int64(len(yxh_b58Alphabet)))
	//初始化余数
	zero := big.NewInt(0)
	//初始化模
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, yxh_b58Alphabet[mod.Int64()])
	}

	Yxh_ReverseBytes(result)
	for b := range input {
		if b == 0x00 {
			//如果第一位是0，则变成1, 所以比特币地址的第一位都是1，//长度是34个字符
			result = append([]byte{yxh_b58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result
}

// Base58转字节数组，解密
func Yxh_Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(yxh_b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}
