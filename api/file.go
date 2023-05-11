package api

import (
	"encryption/core"
	"fmt"
	"io"
	"log"
	"os"
)

// READBYTE 分片大小
var READBYTE = 1024

// StreamEncryptFiles 分片流式加密
// 按照每次1024个字节进行分片加密 (针对大文件加密)
func StreamEncryptFiles(infile, outfile, key string, enLen int) error {

	pass := []byte(key)

	inFile, err := os.OpenFile(infile, os.O_RDONLY, 0644)
	if err != nil {
		log.Println("open file error", err)
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("open file error", err)
		return err
	}
	defer outFile.Close()

	for {
		buffer := make([]byte, READBYTE)
		count, err := inFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		// 加密
		endata, _ := core.Encrypt(buffer[:count], pass, enLen)

		outFile.Write(endata)
	}
	return nil
}

// StreamDecryptFiles  分片流式解密
// 按照每次1024个字节进行分片解密 (针对大文件解密)
func StreamDecryptFiles(infile, outfile, key string, enLen int) error {

	pass := []byte(key)

	inFile, err := os.OpenFile(infile, os.O_RDONLY, 0644)
	if err != nil {
		log.Println("open file error", err)
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("open file error", err)
		return err
	}
	defer outFile.Close()

	for {
		buffer := make([]byte, READBYTE)
		count, err := inFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		// 解密
		endata, _ := core.Decrypt(buffer[:count], pass, enLen)
		outFile.Write(endata)
	}
	return nil
}

func EncryptFiles(infile, outfile, key string, enLen int) error {
	pass := []byte(key)

	data, err := os.ReadFile(infile)
	if err != nil {
		log.Println("read file failed~ ", err)
		return err
	}
	out, err := core.Encrypt(data, pass, enLen)
	if err != nil {
		log.Println("Encrypt file failed~ ", err)
		return err
	}
	log.Println("Encrypt file success~")

	err = os.WriteFile(outfile, out, 0777)
	if err != nil {
		log.Println("write file failed~ ", err)
		return err
	}
	return nil

}

func DecryptFiles(infile, outfile, key string, enLen int) error {
	pass := []byte(key)

	data, err := os.ReadFile(infile)
	if err != nil {
		log.Println("read file failed~ ", err)
		return err
	}
	out, err := core.Decrypt(data, pass, enLen)
	if err != nil {
		log.Println("Decrypt file failed~ ", err)
		return err
	}
	log.Println("Decrypt file success~")

	err = os.WriteFile(outfile, out, 0777)
	if err != nil {
		log.Println("write file failed~ ", err)
		return err
	}
	return nil
}
