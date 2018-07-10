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
)

const yxh_version = byte(0x00)

const yxh_addressChecksumLen = 4

type Wallet struct {
	//1. 私钥
	Yxh_PrivateKey ecdsa.PrivateKey

	//2. 公钥
	Yxh_PublicKey  []byte
}

func Yxh_IsValidForAdress(adress []byte) bool {

	// 25
	version_public_checksumBytes := Yxh_Base58Decode(adress)

	fmt.Println(version_public_checksumBytes)

	//25
	//4
	//21
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes) - yxh_addressChecksumLen:]

	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes) - yxh_addressChecksumLen]

	//fmt.Println(len(checkSumBytes))
	//fmt.Println(len(version_ripemd160))

	checkBytes := Yxh_CheckSum(version_ripemd160)

	if bytes.Compare(checkSumBytes,checkBytes) == 0 {
		return true
	}

	return false
}

//查询地址
func (w *Wallet) Yxh_GetAddress() []byte  {

	//1. hash160
	// 20字节
	ripemd160Hash := Yxh_Ripemd160Hash(w.Yxh_PublicKey)

	// 21字节
	version_ripemd160Hash := append([]byte{yxh_version},ripemd160Hash...)

	// 两次的256 hash
	checkSumBytes := Yxh_CheckSum(version_ripemd160Hash)

	 //25
	bytes := append(version_ripemd160Hash,checkSumBytes...)

	return Yxh_Base58Encode(bytes)
}

// 校验
func Yxh_CheckSum(payload []byte) []byte {

	hash1 := sha256.Sum256(payload)

	hash2 := sha256.Sum256(hash1[:])

	return hash2[:yxh_addressChecksumLen]
}

// ripemd160hash
func Yxh_Ripemd160Hash(publicKey []byte) []byte {
	//1. 256

	hash256 := sha256.New()

	hash256.Write(publicKey)

	hash := hash256.Sum(nil)

	//2. 160

	ripemd160 := ripemd160.New()

	ripemd160.Write(hash)

	return ripemd160.Sum(nil)
}

// 创建钱包
func Yxh_NewWallet() *Wallet {

	privateKey,publicKey := yxh_newKeyPair()

	return &Wallet{privateKey,publicKey}
}

// 通过私钥产生公钥
func yxh_newKeyPair() (ecdsa.PrivateKey,[]byte) {

	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}