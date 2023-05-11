# 加密流程
`假设一次加密n个字节
数据-> 判断是否能够被n整除(不能需要补齐到能够被n整除)-> 循环右移 -> 按位取反 -> 与key进行按位异或 -> 去除补齐的位数`

# 解密流程
`加密逆操作`


# 运行
`使用makefile`
- make 格式化go代码 并编译生成二进制文件
- make build 编译go代码生成二进制文件
- make clean 清理中间目标文件
- make test 执行测试case
- make check 格式化go代码
- make cover 检查测试覆盖率
- make run 直接运行程序
- make lint 执行代码检查
- make docker 构建docker镜像

````
make build
- 加密文件：
./encrypter -en -key yupyup -in test.xml -out test_en.xml
解密文件：
./encrypter -de -key yupyup -in test_en.xml -out test_de.xml
