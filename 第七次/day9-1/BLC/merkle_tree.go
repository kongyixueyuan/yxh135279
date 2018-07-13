package BLC

import (
	"crypto/sha256"
)





type MerkleTree struct {
	Yxh_RootNode *MerkleNode
}


// Block  [tx1 tx2 tx3]


//MerkleNode{nil,nil,tx1Bytes} node1
//MerkleNode{nil,nil,tx2Bytes} node2
//MerkleNode{nil,nil,tx3Bytes} node3
//MerkleNode{nil,nil,tx3Bytes} node4
//
//

//MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes+tx2Bytes)}
//
//MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes+tx3Bytes)}


//
//MerkleNode:
//	left: MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
//
//	right: MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
//
//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))





type MerkleNode struct {
	Yxh_Left  *MerkleNode
	Yxh_Right *MerkleNode
	Yxh_Data  []byte
}


func Yxh_NewMerkleTree(data [][]byte) *MerkleTree {

	//[tx1,tx2,tx3]

	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
		//[tx1,tx2,tx3,tx3]
	}

	// 创建叶子节点
	for _, datum := range data {
		node := Yxh_NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}
///------

	//MerkleNode{nil,nil,tx1Bytes}
	//MerkleNode{nil,nil,tx2Bytes}
	//MerkleNode{nil,nil,tx3Bytes}
	//MerkleNode{nil,nil,tx3Bytes}



	//MerkleNode{nil,nil,tx1Bytes}
	//MerkleNode{nil,nil,tx2Bytes}

	//MerkleNode{nil,nil,tx3Bytes}
	//MerkleNode{nil,nil,tx4Bytes}

	//MerkleNode{nil,nil,tx5Bytes}
	//MerkleNode{nil,nil,tx6Bytes}

	//MerkleNode{nil,nil,tx5Bytes}
	//MerkleNode{nil,nil,tx6Bytes}


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


		if len(newLevel) % 2 != 0 {

			newLevel = append(newLevel,newLevel[len(newLevel) - 1])
		}


		nodes = newLevel
	}

	//MerkleNode:
	//	left: MerkleNode{MerkleNode{nil,nil,tx1Bytes},MerkleNode{nil,nil,tx2Bytes},sha256(tx1Bytes,tx2Bytes)}
	//
	//	right: MerkleNode{MerkleNode{nil,nil,tx3Bytes},MerkleNode{nil,nil,tx3Bytes},sha256(tx3Bytes,tx3Bytes)}
	//
	//	sha256(sha256(tx1Bytes,tx2Bytes)+sha256(tx3Bytes,tx3Bytes))

	mTree := MerkleTree{&nodes[0]}

	return &mTree
}


func Yxh_NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {

	mNode := MerkleNode{}

	// 创建叶子节点
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Yxh_Data = hash[:]
	// 非叶子节点
	} else {
		prevHashes := append(left.Yxh_Data, right.Yxh_Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Yxh_Data = hash[:]
	}

	mNode.Yxh_Left = left
	mNode.Yxh_Right = right

	return &mNode
}