package BLC

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"bytes"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	Yxh_PrivateKey ecdsa.PrivateKey
	Yxh_PublicKey  []byte
}

// 重成新的钱包
func Yxh_NewWallet() *Wallet {
	private, public := Yxh_newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

// 获取钱包地址
func (w Wallet) Yxh_GetAddress() []byte {
	pubKeyHash := Yxh_HashPubKey(w.Yxh_PublicKey)
	versionPayload := append([]byte{version}, pubKeyHash...)
	checksum := Yxh_checksum(versionPayload)
	fullPayload := append(versionPayload, checksum...)
	return Yxh_Base58Encode(fullPayload)
}

// 将公钥先进行Hash256 再进行 RIPEMD160 Hash
func Yxh_HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	return RIPEMD160Hasher.Sum(nil)
}

// 验证钱包地址是否有效
// 地址解码后，将前21个字符两次hash后
// 取前四位跟 解码后的地址后四位对比
func Yxh_ValidateAddress(address string) bool {
	pubKeyHash := Yxh_Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := Yxh_checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// 两次hash256获取校验值，hash的前4个
func Yxh_checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:addressChecksumLen]
}

// 生成新的私钥和公钥
func Yxh_newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}
