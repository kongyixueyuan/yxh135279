package BLC

import "crypto/sha256"

type MerkelTree struct {
	Yxh_RootNode *MerkelNode
}

type MerkelNode struct {
	Yxh_Left  *MerkelNode
	Yxh_Right *MerkelNode
	Yxh_Data  []byte
}

func Yxh_NewMerkelTree(data [][]byte) *MerkelTree {
	var nodes []MerkelNode

	// 如果交易数据不是双数，将最后一个交易复制添加到最后
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	// 生成所有的一级节点，存储到node中
	for _, dataum := range data {
		node := Yxh_NewMerkelNode(nil, nil, dataum)
		nodes = append(nodes, *node)
	}

	// 遍历生成顶层节点
	for i := 0;i<len(data)/2 ;i++{
		var newLevel []MerkelNode
		for j:=0 ; j<len(nodes) ;j+=2  {
			node := Yxh_NewMerkelNode(&nodes[j],&nodes[j+1],nil)
			newLevel = append(newLevel,*node)
		}
		nodes = newLevel
	}

	//for ; len(nodes)==1 ;{
	//	var newLevel []Rwq_MerkelNode
	//	for j:=0 ; j<len(nodes) ;j+=2  {
	//		node := Rwq_NewMerkelNode(&nodes[j],&nodes[j+1],nil)
	//		newLevel = append(newLevel,*node)
	//	}
	//	nodes = newLevel
	//}
	mTree := MerkelTree{&nodes[0]}
	return &mTree
}

// 新叶节点
func Yxh_NewMerkelNode(left, right *MerkelNode, data []byte) *MerkelNode {
	mNode := MerkelNode{}

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
