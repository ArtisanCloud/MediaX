# MediaX

## 简介

MediaX 是一个支持多平台内容发布的轻量级 SDK，旨在帮助开发者快速集成到主流自媒体平台，实现统一的内容发布接口。

ArisanCloud 团队已经成功维护了 [PowerWechat](https://powerwechat.artisan-cloud.com) 开源项目，基本上 SDK 的框架接口趋于稳定，所以我们有延伸出这个 MediaX，同时也做了模块优化， 按需加载模块，支持更多平台的接口使用。

## MediaX 系列产品介绍

- MediaX：一个 开源的 Golang SDK，可以直接被其他项目引用。
- MediaX Studio：基于 MediaX 构建的商业化服务，包含 API（gRPC/HTTP）、高级功能和闭源插件，面向企业用户。

## 功能特点

1. **多平台支持**：支持主流的自媒体平台如抖音、小红书，youtube 等。
2. **统一接口**：通过标准化的接口，实现对多个平台的内容发布。
3. **开源与闭源结合**：Studio 支持插件化架构，开源插件可自由扩展，闭源插件提供增强功能。
4. **跨语言支持**：Studio 通过 gRPC 实现跨语言调用，支持 Go 和其他语言集成。
5. **可扩展性**：灵活的 Provider 机制，轻松接入新平台。

## MediaX 快速开始

### 环境要求

- Go 1.18 或更高版本

### 使用知识

- OAuth2.0 授权流程，可以参考[理解 OAuth 2.0](https://www.ruanyifeng.com/blog/2014/05/oauth_2_0.html)文章

### 安装

1. 初始化项目并下载 MediaX：

   ```bash
   go get github.com/ArtisanCloud/MediaX
   ```

   2. 创建一个简单的示例，本项目作者正在自己系统中使用，陆续会迭代版本：

      ```go

      import (
         "github.com/ArtisanCloud/MediaX/pkg/client"
         config2 "github.com/ArtisanCloud/MediaX/pkg/client/config"
         "github.com/ArtisanCloud/MediaX/pkg/utils"
         "github.com/ArtisanCloud/MediaXCore/pkg/cache"
         "github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
         "github.com/ArtisanCloud/MediaXCore/utils/fmt"
         "github.com/ArtisanCloud/MediaX/pkg/client/wechat/officialAccount/clientTokenClient/publish/schema"
         "github.com/redis/go-redis/v9"
      )

      // 配置Media Client实例的信息
      mediaXClient := client.NewMediaX(&config2.MediaXConfig{
        &config.LogConfig{
            Level:   "debug",
            Console: true,
            File: config.FileConfig{
                Enable: true,
            },
        },
      }, c)

      // 从MediaXClient实例中获取到微信平台中公众号的实例，该实例是Client Token模式，不需要用户授权
      wechatOAClient, err := mediaXClient.MediaXClient.NewWeChatOfficialAccountCTClient(&config2.WeChatOfficialAccountConfig{
        ClientConfig: &ClientConfig{
            BaseConfig: &BaseConfig{
               Timeout: 30,
               HttpDebug: true,
            },
            OAuthConfig: &OAuthConfig{
               ClientID: "your client/app id"
               ClientSecret: "your client/app secret"
            },
        },
      })
      if err != nil {
        panic(err)
      }

      // 调用 wechatOAClient 的方法
      ctx := context.Background()
      var reqData = &schema.DraftAddReq{}
      resData, err := oaClient.GetPublishClient().DraftAdd(ctx, reqData)
      if err != nil {
         return nil, err
      }

      ```

## 文档与接口说明

详细的功能接口文档和使用指南，请访问我们的[开发者中心](https://mediax.artisan-cloud.com)。在这里您可以找到:

- 完整的 API 参考文档
- 快速入门指南
- 最佳实践示例
- 常见问题解答

## 功能矩阵

| 平台     | 应用    | 授权类型          | 图文 | 视频 | 素材管理 | 评论管理 | 数据管理 |
| -------- | ------- | ----------------- | ---- | ---- | -------- | -------- | -------- |
| WeChat   | 公众号  | client_credential | ✔    | ✘    | ✔        | 🚧       | 🚧       |
| 小红书   | 聚光    | auth_code         | ✔    | ✔    | ✔        | ✘        | ✔        |
| 字节     | 抖音    | auth_code         | ✔    | ✔    | ✔        | ✔        | ✔        |
| 字节     | 抖音    | client_credential | ✔    | ✔    | ✔        | ✔        | ✔        |
| Bilibili | B 站    | auth_code         | 🚧   | ✔    | 🚧       | 🚧       | 🚧       |
| Google   | YouTube | auth_code         | ✔    | ✔    | ✔        | ✔        | ✔        |
| Google   | Blogger | auth_code         | 🚧   | ✘    | 🚧       | 🚧       | 🚧       |
| ......   |         |                   |      |      |          |          |          |

> 注：✔ 表示支持该功能，✘ 表示暂不支持该功能，🚧 表示功能开发中。

<!-- ## 功能介绍

\*\*\* [项目功能的开发安排](https://github.com/orgs/ArtisanCloud/projects/5/views/2) -->

## 许可证

MediaX SDK 项目 采用 [MIT License](./LICENSE) 开源。
