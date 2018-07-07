package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"strconv"
)

//将int16转换为字节数组
func Int16ToBytes(num int16) []byte {
	buff := make([]byte,8)
	binary.BigEndian.PutUint16(buff, uint16(num))
	return buff
}

//将int32转换为字节数组
func Int32ToBytes(num int32) []byte {
	buff := make([]byte,8)
	binary.BigEndian.PutUint32(buff, uint32(num))
	return buff
}

//将int64转换为字节数组
func Int64ToBytes(num int64) []byte {
	buff := make([]byte,8)
	binary.BigEndian.PutUint64(buff, uint64(num))
	return buff
}

//将字节数组转换成16位整形
func BytesToInt16(buff []byte) int16 {
	num := binary.BigEndian.Uint16(buff)
	return int16(num)
}

//将字节数组转换成32位整形
func BytesToInt32(buff []byte) int32 {
	num := binary.BigEndian.Uint32(buff)
	return int32(num)
}

//将字节数组转换成64位整形
func BytesToInt64(buff []byte) int64 {
	num := binary.BigEndian.Uint64(buff)
	return int64(num)
}

//将字幕转换成字符串
func BytsToString(buff []byte) string {
	return string(buff)
}

// 字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

// 将int64转换为字节数组
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}


// 标准的JSON字符串转数组
func JSONToArray(jsonString string) []string {

	//json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}

// 求整数的二进制方法
func ConverToBin(num int) string {

	var result string
	for ; num > 0; num /=2 {
		lsb := num % 2
		result = strconv.Itoa(lsb) + result
	}
	return result

}