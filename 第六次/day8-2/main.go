package main

import "gostudy/blockchain/day8-2/BLC"

func main()  {

	//CLI入口
	cli := BLC.CLI{}
	cli.Yxh_Run()

	//var yxh_b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	//var result []byte
	//x := big.NewInt(0).SetBytes([]byte("123"))
	//
	//base := big.NewInt(int64(58))
	//zero := big.NewInt(0)
	//mod := &big.Int{}
	//
	//for x.Cmp(zero) != 0 {
	//	x.DivMod(x, base, mod)
	//	result = append(result, yxh_b58Alphabet[mod.Int64()])
	//}
	//
	//fmt.Printf("result: %x", result)
	//fmt.Println()
	//
	//b := BLC.Yxh_Base58Encode([]byte("123"))
	//fmt.Printf("b: %x", b)


}