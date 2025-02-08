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
   go get github.com/ArtisanCloud/MediaX/v1
   ```

2. 创建一个简单的示例，本项目作者正在自己系统中使用，陆续会迭代版本：

   ```go
   wechatClient, err := mediaX.CreateWechat(wechatConfig)
	if err != nil {
		panic(err)
	}

	// 调用 WeChatClient 的方法
	ctx := context.Background()
   var content = ...
	wechatClient.Publish(ctx,content)
   
   ```

## 功能介绍
*** [项目功能的开发安排](https://github.com/orgs/ArtisanCloud/projects/5/views/2)



## 许可证

本项目采用 [MIT License](./LICENSE) 开源。

