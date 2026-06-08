// config-cipher — config.yaml AES 加密/解密工具
//
// 用法：
//
//	config-cipher encrypt [-in config.yaml] [-out config.enc]   # 加密
//	config-cipher decrypt [-in config.enc]  [-out config.yaml]  # 解密（-out 留空则输出到 stdout）
//
// 密钥内置在 pkg/aes/cipher_key.go，加密后只需部署 config.enc。
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	pkgaes "github.com/v03413/bepusdt/pkg/aes"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	encCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
	encIn := encCmd.String("in", "config.yaml", "待加密的明文配置文件路径")
	encOut := encCmd.String("out", "config.enc", "加密后输出文件路径")

	decCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
	decIn := decCmd.String("in", "config.enc", "待解密的加密配置文件路径")
	decOut := decCmd.String("out", "", "解密后输出文件路径（留空则打印到 stdout）")

	switch os.Args[1] {
	case "encrypt":
		_ = encCmd.Parse(os.Args[2:])
		doEncrypt(*encIn, *encOut)
	case "decrypt":
		_ = decCmd.Parse(os.Args[2:])
		doDecrypt(*decIn, *decOut)
	default:
		printUsage()
		os.Exit(1)
	}
}

func doEncrypt(inFile, outFile string) {
	plain, err := os.ReadFile(inFile)
	if err != nil {
		fatalf("读取文件失败: %v", err)
	}
	cipherText, err := pkgaes.Encrypt(string(plain), pkgaes.ConfigCipherKey)
	if err != nil {
		fatalf("加密失败: %v", err)
	}
	if err := os.WriteFile(outFile, []byte(cipherText), 0o600); err != nil {
		fatalf("写入文件失败: %v", err)
	}
	fmt.Printf("加密成功: %s → %s\n", inFile, outFile)
	fmt.Println("部署时只需携带 config.enc，请勿将明文 config.yaml 放到服务器上。")
}

func doDecrypt(inFile, outFile string) {
	raw, err := os.ReadFile(inFile)
	if err != nil {
		fatalf("读取文件失败: %v", err)
	}
	plain, err := pkgaes.Decrypt(strings.TrimSpace(string(raw)), pkgaes.ConfigCipherKey)
	if err != nil {
		fatalf("解密失败: %v", err)
	}
	if outFile == "" {
		fmt.Print(plain)
		return
	}
	if err := os.WriteFile(outFile, []byte(plain), 0o600); err != nil {
		fatalf("写入文件失败: %v", err)
	}
	fmt.Printf("解密成功: %s → %s\n", inFile, outFile)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "错误: "+format+"\n", args...)
	os.Exit(1)
}

func printUsage() {
	fmt.Print(`config-cipher — config.yaml AES 加密/解密工具

用法:
  config-cipher encrypt [-in config.yaml] [-out config.enc]
  config-cipher decrypt [-in config.enc]  [-out config.yaml]

密钥定义在 pkg/aes/cipher_key.go，与服务启动解密使用同一密钥。

示例:
  config-cipher encrypt                          # 加密 config.yaml → config.enc
  config-cipher decrypt                          # 解密打印到 stdout（用于验证）
  config-cipher decrypt -out /tmp/plain.yaml     # 解密到文件
`)
}
