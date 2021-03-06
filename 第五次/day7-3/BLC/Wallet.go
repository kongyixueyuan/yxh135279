package BLC

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"log"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"fmt"
	"bytes"
	"encoding/hex"
)

const version = byte(0x00) //占1位
const addressChecksumLen = 4 //校验码的长度，取前面4个

type Wallet struct {
	//1. 私钥
	PrivateKey ecdsa.PrivateKey

	//2. 公钥
	PublicKey  []byte
}

func IsValidForAdress(adress []byte) bool {

	version_public_checksumBytes := Base58Decode(adress)

	fmt.Println(hex.EncodeToString(version_public_checksumBytes))

	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes) - addressChecksumLen:]

	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes) - addressChecksumLen]

	//fmt.Println(len(checkSumBytes))
	//fmt.Println(len(version_ripemd160))

	checkBytes := CheckSum(version_ripemd160)

	if bytes.Compare(checkSumBytes,checkBytes) == 0 {
		return true
	}

	return false
}

// 取地址
func (w *Wallet) GetAddress() []byte  {

	//1. hash160

	ripemd160Hash := Ripemd160Hash(w.PublicKey)

	version_ripemd160Hash := append([]byte{version},ripemd160Hash...)

	checkSumBytes := CheckSum(version_ripemd160Hash)

	bytes := append(version_ripemd160Hash,checkSumBytes...)

	return Base58Encode(bytes)
}

//产生校验码(前4位)
func CheckSum(payload []byte) []byte {

	hash1 := sha256.Sum256(payload)

	hash2 := sha256.Sum256(hash1[:])

	return hash2[:addressChecksumLen]
}

// 计算pubhash
func Ripemd160Hash(publicKey []byte) []byte {
	//1. 256位hash
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	//2. 160位hash
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil)
}

// 创建钱包
func NewWallet() *Wallet {

	privateKey,publicKey := newKeyPair()

	return &Wallet{privateKey,publicKey}
}

// 通过私钥产生公钥
func newKeyPair() (ecdsa.PrivateKey,[]byte) {

	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}