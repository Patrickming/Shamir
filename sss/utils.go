package sss

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
)

func ToByteArray(input interface{}) ([]byte, error) {
	// 根据输入类型进行转换
	switch v := input.(type) {
	case int:
		// 转换整数为字符串，然后转换为字节数组
		return []byte(strconv.Itoa(v)), nil
	case float64:
		// 转换浮点数为字符串，然后转换为字节数组
		return []byte(strconv.FormatFloat(v, 'f', -1, 64)), nil
	case string:
		// 直接转换字符串为字节数组
		return []byte(v), nil
	default:
		// 如果输入的类型不是上述任何一种，返回错误
		return nil, fmt.Errorf("unsupported type: %v", reflect.TypeOf(input))
	}
}

// 将字节数组转换为一个大整数
func bytesToBigInt(data []byte) *big.Int {
	hash := sha256.Sum256(data)
	bigInt := new(big.Int)
	bigInt.SetBytes(hash[:])
	return bigInt
}

// 将任意数据编码为字节
func serialize(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 将字节解码为原始数据
func deserialize(data []byte, target interface{}) error {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(target)
	return err
}

// powMod 计算a的b次方对p取模的结果
func PowMod(a, b, p uint8) uint8 {
	result := uint8(1)
	base := a % p

	for b > 0 {
		if b&1 == 1 { // 检查b的最低位是否为1
			result = (result * base) % p
		}
		base = (base * base) % p
		b >>= 1 // 右移一位
	}

	return result
}
