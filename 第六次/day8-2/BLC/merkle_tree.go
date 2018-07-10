package BLC

import (
	"crypto/sha256"
)

//MerkleTree结构
type MerkleTree struct {
	Yxh_RootNode *MerkleNode
}

//MerkleNode节点(二叉树，节点里有节点)
type MerkleNode struct {
	Yxh_Left  *MerkleNode
	Yxh_Right *MerkleNode
	Yxh_Data  []byte
}

// Block  [tx1 tx2 tx3 tx3]

//MerkleNode{nil,nil,tx1Bytes}
//MerkleNode{nil,nil,tx2Bytes}
//MerkleNode{nil,nil,tx3Bytes}
//MerkleNode{nil,nil,tx3Bytes}
//
//

//
//MerkleNode:
//	left: MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
//
//	right: MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
//
//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))

//创建MerkleTree
func Yxh_NewMerkleTree(data [][]byte) *MerkleTree {

	var nodes []MerkleNode

	// 如果节点是单数，则将最后一个节点复制一个变成双数
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	// 创建叶子节点
	for _, datum := range data {
		node := Yxh_NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	//MerkleNode{nil,nil,tx1Bytes}
	//MerkleNode{nil,nil,tx2Bytes}
	//MerkleNode{nil,nil,tx3Bytes}
	//MerkleNode{nil,nil,tx3Bytes}

	// 　循环两次
	for i := 0; i < len(data)/2; i++ {

		var newLevel []MerkleNode

		for j := 0; j < len(nodes); j += 2 {

			node := Yxh_NewMerkleNode(&nodes[j], &nodes[j+1], nil)

			newLevel = append(newLevel, *node)
		}

		//MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
		//
		//MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
		//

		nodes = newLevel
	}

	//MerkleNode:
	//	left: MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
	//
	//	right: MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
	//
	//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))

	//只存储第一个节点
	mTree := MerkleTree{&nodes[0]}

	return &mTree
}

//创建新节点
func Yxh_NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {

	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Yxh_Data = hash[:]
	} else {
		prevHashes := append(left.Yxh_Data, right.Yxh_Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Yxh_Data = hash[:]
	}

	mNode.Yxh_Left = left
	mNode.Yxh_Right = right

	return &mNode
}