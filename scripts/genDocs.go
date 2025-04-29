package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 使用 godocdown 工具生成 Markdown 文档
func generateMarkdownWithGodocdown(goFilePath, docsDir string) (markdownFile string, err error) {
	// 生成对应的 Markdown 文件名（去掉 pkg/client/ 前缀）
	relPath := strings.TrimPrefix(goFilePath, "pkg/client/")
	relPath = strings.TrimSuffix(relPath, ".go") + ".md"
	mdFilePath := filepath.Join(docsDir, relPath)

	// 创建输出目录
	if err := os.MkdirAll(filepath.Dir(mdFilePath), 0755); err != nil {
		return "", err
	}

	// 使用 godocdown 处理文件所在的目录（而不是文件本身）
	dir := filepath.Dir(goFilePath)
	log.Printf("生成文档: godocdown %s -> %s\n", dir, mdFilePath)

	cmd := exec.Command("godocdown", dir)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("godocdown 执行失败: %w", err)
	}

	// 写入 Markdown 文件
	if err := os.WriteFile(mdFilePath, output, 0644); err != nil {
		return "", err
	}

	return mdFilePath, nil
}

// 遍历目录并调用 godocdown 生成 Markdown
func traverseDirectoryAndGenerateMD(directory string, docsDir string, excludeDirs map[string]bool) ([]string, error) {
	visited := map[string]bool{} // 避免重复处理同一目录
	markdownFiles := []string{}
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 排除目录
		if info.IsDir() {
			if excludeDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// 只处理 .go 文件
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		// 避免对同一目录多次执行 godocdown
		dir := filepath.Dir(path)
		if visited[dir] {
			return nil
		}
		visited[dir] = true

		mdFilePath, err := generateMarkdownWithGodocdown(path, docsDir)
		if err != nil {
			return err
		}
		markdownFiles = append(markdownFiles, mdFilePath)

		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Println("✅ Markdown 文档生成完毕！")
	return markdownFiles, nil
}

type SidebarGroup struct {
	Text  string        `json:"text"`
	Items []SidebarItem `json:"items"`
}
type SidebarItem struct {
	Text string `json:"text"`
	Link string `json:"link"`
}

// 假设 markdownFiles 是 []string，例如 ["docs/mediax/foo.md"]
func buildSidebar(markdownFiles []string) ([]SidebarGroup, error) {
	groups := make(map[string]*SidebarGroup)

	for _, path := range markdownFiles {
		rel, _ := filepath.Rel("docs", path) // mediax/foo.md
		parts := strings.Split(rel, string(os.PathSeparator))
		if len(parts) < 2 {
			continue
		}
		groupName := strings.Title(strings.ReplaceAll(parts[0], "-", " "))
		link := "/" + strings.TrimSuffix(rel, ".md")

		item := SidebarItem{
			Text: strings.TrimSuffix(filepath.Base(rel), ".md"),
			Link: link,
		}

		if groups[groupName] == nil {
			groups[groupName] = &SidebarGroup{
				Text:  groupName,
				Items: []SidebarItem{item},
			}
		} else {
			groups[groupName].Items = append(groups[groupName].Items, item)
		}
	}

	var result []SidebarGroup
	for _, g := range groups {
		result = append(result, *g)
	}
	return result, nil
}

func writeSidebarToFile(sidebar []SidebarGroup, outputPath string) error {
	// 序列化为 JSON 字符串
	jsonBytes, err := json.MarshalIndent(sidebar, "", "  ")
	if err != nil {
		return err
	}

	// 构造最终的 TypeScript 文件内容
	content := fmt.Sprintf("export const sidebar = %s;\n", string(jsonBytes))

	// 写入文件
	err = os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var targetDir string
	var docsDir string

	flag.StringVar(&targetDir, "target", "pkg/client", "指定遍历的 Go 源代码目录")
	flag.StringVar(&docsDir, "docs", "docs", "指定 Markdown 文档输出目录")
	flag.Parse()

	// 排除某些子目录
	excludeDirs := map[string]bool{
		"core":   true,
		"schema": true,
		"config": true,
	}
	_, err := traverseDirectoryAndGenerateMD(targetDir, docsDir, excludeDirs)
	if err != nil {
		log.Fatal(err)
	}

	//sidebarGroup, err := buildSidebar(markdownFiles)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//vitePressSidebarPath := filepath.Join(docsDir, "../../.vitepress", "sidebar.mts")
	//println(vitePressSidebarPath)
	//err = writeSidebarToFile(sidebarGroup, vitePressSidebarPath)
	//if err != nil {
	//	log.Fatal(err)
	//}

}
