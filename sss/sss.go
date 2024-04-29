package sss

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
)

const prime = 251

var (
	SECRET         uint8  //uint8秘密
	SERIALIZEDDATA []byte //序列化后的[]byte秘密
	T              int    //秘密字节数
)

// 打印以及加密
func Encrypt(secret []byte, t int) ([][]byte, error) {
	T = t
	secrets, coeffs, err := ShareSecret(secret, t)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Secret polynomial and points:\n")
	for _, point := range secrets {
		fmt.Printf("Part x = %d, y = %d\n", point[0], point[1])
	}

	expr := ""
	for i, coeff := range coeffs {
		if i > 0 {
			expr += fmt.Sprintf(" + %dx^%d", coeff, i)
		} else {
			expr = fmt.Sprintf("%d", coeff)
		}
	}
	fmt.Printf("Polynomial: %s = y\n", expr)
	return secrets, nil
}

// ShareSecret 分割秘密成为多个小秘密片段，需要至少t个片段来恢复原始秘密。
// 返回秘密的分享点和多项式系数。
func ShareSecret(secret []byte, t int) ([][]byte, []uint8, error) {
	n := len(secret) // n 由秘密的字节数组长度确定
	if t > n {
		return nil, nil, errors.New("错误：需要的最小片段数不能大于秘密的长度")
	}

	// 序列化
	serializedData, err := serialize(secret)
	if err != nil {
		fmt.Println("Serialization error:", err)
	}
	SERIALIZEDDATA = serializedData
	fmt.Println("serializedData:", serializedData)

	// 转换为大整数
	bignint := bytesToBigInt(serializedData)

	fmt.Println("Big integer:", bignint)

	// 对大整数取模
	// 作用：当做秘密m提供给我们进行计算
	m := uint8(bignint.Mod(bignint, big.NewInt(251)).Uint64())
	SECRET = m
	fmt.Println("modbigintTOuint8 - secret - uint8:", m)

	//多项式
	coefficients := make([]uint8, t)
	coefficients[0] = m // 秘密m作为多项式的常数项系数
	for i := 1; i < t; i++ {
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			return nil, nil, err
		}
		coefficients[i] = randomByte[0] % prime
	}

	points := make([][]byte, n)
	for x := 1; x <= n; x++ {
		fx := uint64(0) // 使用更大的数据类型以避免溢出
		for power, coeff := range coefficients {
			term := uint64(coeff) * uint64(PowMod(uint8(x), uint8(power), prime))
			fx += term
		}
		fx %= uint64(prime) // 模运算在所有加法操作完成后执行
		points[x-1] = []byte{uint8(x), uint8(fx)}
	}

	return points, coefficients, nil
}

// calculateLagrangePolynomial 用于计算拉格朗日基本多项式l_i(x)
func calculateLagrangePolynomial(x int, i int, points [][]byte, prime uint8) *big.Int {
	num := big.NewInt(1)
	den := big.NewInt(1)
	xBig := big.NewInt(int64(x))

	for j, point := range points {
		if j != i {
			xj := big.NewInt(int64(point[0]))
			xi := big.NewInt(int64(points[i][0]))

			// 计算分子 (x - x_j)
			num.Mul(num, new(big.Int).Sub(xBig, xj))

			// 计算分母 (x_i - x_j)
			den.Mul(den, new(big.Int).Sub(xi, xj))
		}
	}

	// 分母的逆元
	denInv := new(big.Int).ModInverse(den, big.NewInt(int64(prime)))
	if denInv == nil {
		return big.NewInt(0) // 如果不存在逆元，直接返回0
	}

	// l_i(x) = (分子 * 分母的逆元) mod p
	return new(big.Int).Mod(num.Mul(num, denInv), big.NewInt(int64(prime)))
}

// Decrypt 用拉格朗日插值法恢复秘密
func Decrypt(points [][]byte) (string, error) {
	if len(points) < T { //TODO
		return "", errors.New("提供的point不足以解开秘密")
	}

	prime := uint8(251) // 设置模数，与加密时相同
	secret := big.NewInt(0)

	for i, point := range points {
		li := calculateLagrangePolynomial(0, i, points, prime)                // 计算 l_i(0)
		secret.Add(secret, new(big.Int).Mul(li, big.NewInt(int64(point[1])))) // 累加 y_i * l_i(0)
	}

	//取模
	DecryptSecret := uint8(secret.Mod(secret, big.NewInt(int64(prime))).Uint64()) // 最终秘密 s = L(0) mod p

	//判断秘密与加密前是否一致
	fmt.Println("Decrypted secret:", DecryptSecret)
	if DecryptSecret != SECRET {
		fmt.Printf("Fatal: 解密后的秘密与原先的不一致\n")
		os.Exit(-1)
	}

	// 反序列化
	var deserializedData []byte
	err := deserialize(SERIALIZEDDATA, &deserializedData)
	if err != nil {
		fmt.Println("Deserialization error:", err)
		os.Exit(-1)
	}

	return string(deserializedData), nil
}
