package main

import (
	"fmt"
	"math/rand"
	"os"
	"shamir/sss"
	"time"
)

func main() {
	// 提示用户输入数据
	var data string
	fmt.Print("输入你要加密的秘密: ")
	_, err := fmt.Scanln(&data)
	if err != nil {
		fmt.Println("输入有误请重试:", err)
		return
	}

	// 转换成字节数组的形式
	secret, _ := sss.ToByteArray(data)

	// 询问需要多少份秘密份额来恢复秘密
	var t int
	fmt.Printf("请输入需要的最小片段数(1-%d): ", len(secret))
	_, err = fmt.Scanln(&t)
	if err != nil || t < 1 || t > len(secret) {
		fmt.Println("无效的输入，请输入一个1到秘密长度之间的整数")
		return
	}

	points, err := sss.Encrypt(secret, t)
	if err != nil {
		fmt.Println("加密过程中出错:", err)
		return
	}

	// 随机选择t份秘密份额
	rand.Seed(time.Now().UnixNano())
	selectedPoints := make([][]byte, t)
	for i, index := range rand.Perm(len(points))[:t] {
		selectedPoints[i] = points[index]
	}

	fmt.Println("随机选择的解密对↓:")
	for _, pt := range selectedPoints {
		fmt.Printf("(%d, %d)\n", pt[0], pt[1])
	}

	// 恢复秘密
	combined, err := sss.Decrypt(selectedPoints)
	if err != nil {
		fmt.Printf("解密失败: %v\n", err)
		os.Exit(-1)
	}
	fmt.Printf("解密后的秘密: %s\n", string(combined))
}
