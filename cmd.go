package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

//				循环移位加密算法 加密文件
//	        -en file       需要加密的文件
//	        -de file       需要解密的文件
//	        -out file      结果输出文件
//	        -key str       密码
//	        -test file     如果是加密操作，尝试把加密结果文件解密
//	        -help str      查看帮助

var (
	inFile  = flag.String("in", "", "输入文件")
	outFile = flag.String("out", "", "输出文件")
	key     = flag.String("key", "", "密钥")
	enLen   = flag.Int("n", 4, "n字节为一组")
	en      = flag.Bool("en", false, "加密文件")
	de      = flag.Bool("de", false, "解密文件")
)

func init() {
	// 设置日志
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	flag.Parse()
}

// READBYTE 分片大小
var READBYTE = 256

// NEEDSHARD 分片阈值
var NEEDSHARD int64 = 1024 // 1mb

// 分片加密解密 大小会变

func main() {
	//判断参数是否正常
	if *en && *de {
		log.Fatalln("Cannot encrypt and decrypt at the same time~")
	}

	if *inFile == "" || *outFile == "" || *key == "" {
		flag.PrintDefaults()
		log.Fatalln("error params~")
	}

	size, err := fileSize(*inFile)
	if err != nil {
		log.Fatalln("file error! file is ", *inFile)
	}

	if *en {
		// 如果文件过大则使用分片
		if size >= NEEDSHARD {
			err = streamencryptFiles(*inFile, *outFile, *key, *enLen)
			if err != nil {
				log.Fatalln("encryptFiles error! file is ", *inFile)
			}
		} else {
			// 加密
			log.Println("正常加密")
			err = encryptFiles(*inFile, *outFile, *key, *enLen)
			if err != nil {
				log.Fatalln("encryptFiles error! file is ", *inFile)
			}
		}
	} else {
		if size >= NEEDSHARD {
			err = streamdecryptFiles(*inFile, *outFile, *key, *enLen)
			if err != nil {
				log.Fatalln("decryptFiles error! file is ", *inFile)
			}
		} else {
			// 解密
			log.Println("正常解密")
			err := decryptFiles(*inFile, *outFile, *key, *enLen)
			if err != nil {
				log.Fatalln("decryptFiles error! file is ", *inFile)
			}
		}
	}
	os.Exit(0)
}

func fileSize(filename string) (int64, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Println("file err, ", err)
		return 0, err
	}
	return fileInfo.Size(), nil
}

// 分片流式加密
// 按照每次1024个字节进行分片加密 (针对大文件加密)
func streamencryptFiles(infile, outfile, key string, enLen int) error {

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
		endata, _ := Encrypt(buffer[:count], pass, enLen)

		outFile.Write(endata)
	}
	return nil
}

// 分片流式解密
// 按照每次1024个字节进行分片解密 (针对大文件解密)
func streamdecryptFiles(infile, outfile, key string, enLen int) error {

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
		endata, _ := Decrypt(buffer[:count], pass, enLen)
		outFile.Write(endata)
	}
	return nil
}

func encryptFiles(infile, outfile, key string, enLen int) error {
	pass := []byte(key)

	data, err := os.ReadFile(infile)
	if err != nil {
		log.Println("read file failed~ ", err)
		return err
	}
	out, err := Encrypt(data, pass, enLen)
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

func decryptFiles(infile, outfile, key string, enLen int) error {
	pass := []byte(key)

	data, err := os.ReadFile(infile)
	if err != nil {
		log.Println("read file failed~ ", err)
		return err
	}
	out, err := Decrypt(data, pass, enLen)
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
