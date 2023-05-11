package main

import (
	"encryption/api"
	"flag"
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

func main() {
	//判断参数是否正常
	if *en && *de {
		log.Fatalln("Cannot encrypt and decrypt at the same time~")
	}

	if *inFile == "" || *outFile == "" || *key == "" {
		flag.PrintDefaults()
		log.Fatalln("error params~")
	}

	if *en {
		// 走流处理
		err := api.StreamEncryptFiles(*inFile, *outFile, *key, *enLen)
		if err != nil {
			log.Fatalln("encryptFiles error! file is ", *inFile)
		}
	} else {
		err := api.StreamDecryptFiles(*inFile, *outFile, *key, *enLen)
		if err != nil {
			log.Fatalln("decryptFiles error! file is ", *inFile)
		}
	}
	os.Exit(0)
}
