package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, os.ModePerm)
}

func copyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		log.Printf("source file: %s\n", relPath)
		dstPath := filepath.Join(dstDir, relPath)
		log.Printf("target file: %s\n", relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, os.ModePerm)
		}

		return copyFile(path, dstPath)
	})
}

func main() {
	var srcDir string
	var dstDir string

	flag.StringVar(&srcDir, "src", "manifest/protobuf/proto/", "指定遍历的 Go 源代码目录")
	flag.StringVar(&dstDir, "dst", "../MediaXProtobuf", "指定 protobuf 文档输出目录")
	flag.Parse()

	err := copyDir(srcDir, dstDir)
	if err != nil {
		panic(err)
	}
}
