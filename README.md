# MediaX

MediaX 是一个支持多平台内容发布的轻量级 SDK，旨在帮助开发者快速集成到主流自媒体平台，实现统一的内容发布接口。

## 功能特点

1. **多平台支持**：支持主流的自媒体平台如抖音、小红书等。
2. **统一接口**：通过标准化的接口，实现对多个平台的内容发布。
3. **开源与闭源结合**：支持插件化架构，开源插件可自由扩展，闭源插件提供增强功能。
4. **跨语言支持**：通过 gRPC 实现跨语言调用，支持 Go 和其他语言集成。
5. **可扩展性**：灵活的 Provider 机制，轻松接入新平台。

## 快速开始

### 环境要求

- Go 1.18 或更高版本
- Git

### 安装

1. 初始化项目并下载 MediaX：

   ```bash
   go get github.com/ArtisanCloud/MediaX
   ```

2. 创建一个简单的主程序：

   ```go
   package main

   import (
       "fmt"
       "github.com/ArtisanCloud/MediaX"
   )

   func main() {
       // 初始化 MediaX
       pub := MediaX.New()

       // 添加抖音 Provider
       pub.RegisterProvider("douyin", NewDouyinProvider())

       // 添加小红书 Provider
       pub.RegisterProvider("redbook", NewRedbookProvider())

       // 发布内容到多个平台
       err := pub.Publish([]string{"douyin", "redbook"}, &MediaX.Content{
           Title:       "发布标题",
           SubTitle:    "发布副标题",
           Description: "描述内容",
           Media:       []string{"image1.jpg", "video1.mp4"},
       })

       if err != nil {
           fmt.Println("发布失败:", err)
       } else {
           fmt.Println("发布成功!")
       }
   }
   ```

3. 运行程序：

   ```bash
   go run main.go
   ```

## 目录结构

```
MediaX/
├── providers/           # 各平台 Provider 实现
│   ├── douyin.go        # 抖音 Provider
│   └── redbook.go       # 小红书 Provider
├── grpc/                # gRPC 服务实现
├── examples/            # 使用示例
├── LICENSE              # 许可证
├── README.md            # 项目说明
└── go.mod               # Go 模块
```

## 开发指南

### 添加新平台

1. 创建新的 Provider 实现：

   ```go
   type NewPlatformProvider struct {}

   func (p *NewPlatformProvider) Publish(content *MediaX.Content) error {
       // 实现具体的发布逻辑
       return nil
   }
   ```

2. 注册 Provider：

   ```go
   pub.RegisterProvider("new_platform", NewNewPlatformProvider())
   ```

### gRPC 支持

通过 `grpc` 文件夹中的服务定义，可以轻松扩展跨语言支持。

## 许可证

本项目采用 [MIT License](./LICENSE) 开源。

