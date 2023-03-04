package BLC

import (
	"bytes"
	"math/big"
)

// base64的全字符包含:
// ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/
// base58是在base64的基础上, 去掉0, O(大写的o), I(大写的i), l(小写的L), +, / 构成的
var b58Alphabet = []byte("ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz123456789")

// 字节数组转Base58, 加密过程
func Base58Encode(input []byte) []byte {
	var result []byte
	x := big.NewInt(0).SetBytes(input)
	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	// ↓ x大于0返回+1; 等于返回0; 小于返回-1
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	ReverseBytes(result)
	for b := range input {
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...) // 这个...是什么意思?
			// ↑ ...是展开运算符, 如果result是个[a, b, c], 那么result...就是不把result当作列表传入而是分别传入a, b, c
			// ↑↑ 即append(A, result...) = append(A, a, b, c)
		} else {
			break
		}
	}

	return result
}

// Base58转字节数组, 解密
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}
