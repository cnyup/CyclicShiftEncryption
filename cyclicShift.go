package main

import (
	"bytes"
)

type body struct {
	data          []byte
	key           []byte
	enLen         int //一次按照enLen个字节处理
	dLen          int //数据长度
	kLen          int //key长度
	complementNum int // 补齐的位数
}

// 封装 对外使用
func Encrypt(data, key []byte, enLen int) ([]byte, error) {
	b := &body{
		data:          data,
		key:           key,
		enLen:         enLen,
		dLen:          len(data),
		kLen:          len(key),
		complementNum: 0,
	}
	b.check()
	return b.encryption(), nil
}

// 封装 对外使用
func Decrypt(data, key []byte, enLen int) ([]byte, error) {
	b := &body{
		data:          data,
		key:           key,
		enLen:         enLen,
		dLen:          len(data),
		kLen:          len(key),
		complementNum: 0,
	}
	b.check()
	return b.decryption(), nil
}

// 检查
func (b *body) check() {
	if b.enLen == 0 {
		b.enLen = b.kLen
	}
}

// 对数据进行补全处理  根据enLen进行补全
func (b *body) complement() {
	newData := make([]byte, b.dLen)
	if b.dLen%b.enLen != 0 {
		b.complementNum = b.enLen - (b.dLen % b.enLen)
		newData = append(b.data, bytes.Repeat([]byte{0}, b.complementNum)...)
		b.data = newData
		b.dLen = len(newData)
	}

}

// 按位取反
// 参数：需要处理的数据
// 返回值： 处理后的数据
func (b *body) bitInvert(data []byte) {
	for i := 0; i < b.dLen; i++ {
		data[i] = ^data[i]
	}
}

// 循环右移
// 先右移8-n位(一个字节位数为8)，然后按位或上 左移n位
// 参数：需要处理的数据, 移动的位数
// 返回值： 处理后的数据
func (b *body) shiftRight(data []byte) {
	res := make([]byte, len(data))
	shift := (b.enLen + b.kLen) % 8
	for i := 0; i < len(data); i++ {
		res[i] = data[i]>>shift | data[i]<<(8-shift)
	}
}

// 循环左移 与右移相反
func (b *body) shiftLeft(data []byte) {
	res := make([]byte, len(data))
	shift := (b.enLen + b.kLen) % 8
	for i := 0; i < len(data); i++ {
		res[i] = data[i]<<shift | data[i]>>(8-shift)
	}
}

func (b *body) encryption() []byte {

	// 判断传进来的字节数组是不是单次处理的倍数(4的倍数)，不够的话需要低位补齐
	// 加密后需要对补全的字符去除
	// 解密之后需要去掉补齐的位数， 所以需要记录补齐了多少位
	b.complement()

	// 循环右移   以右移加密
	b.shiftRight(b.data)
	// 取反
	b.bitInvert(b.data)

	//与key进行异或
	//data一次取enLen个字节 与key中的某个字节(长度取余循环次数)进行异或
	res := make([]byte, b.dLen)
	for i := 0; i < b.dLen; i += b.enLen {
		for j := 0; j < b.enLen; j++ {
			if i+j < b.dLen {
				res[i+j] = b.data[i+j] ^ b.key[j%b.kLen] // len(key)%j 会越界
			}
		}
	}

	// 去除补全的字符
	return res[:(len(res) - b.complementNum)]
}

func (b *body) decryption() []byte {
	b.complement()
	//加密的逆操作
	res := make([]byte, b.dLen)

	//与key异或
	for i := 0; i < b.dLen; i += b.enLen {
		for j := 0; j < b.enLen; j++ {
			// 需要判断 否则会存在越界情况
			if i+j < b.dLen {
				res[i+j] = b.data[i+j] ^ b.key[j%b.kLen]
			}
		}
	}
	b.data = res
	b.bitInvert(b.data)
	b.shiftLeft(b.data)

	return res[:(len(res) - b.complementNum)]
}
