# MediaX

MediaX 是一个支持多平台内容发布的轻量级 SDK，旨在帮助开发者快速集成到主流自媒体平台，实现统一的内容发布接口。

MediaX系列产品介绍
- MediaX：一个 开源的Golang SDK，可以直接被其他项目引用。
- MediaX Pro：基于 MediaX 构建的商业化服务，包含 API（gRPC/HTTP）、高级功能和闭源插件，面向企业用户。

## 功能特点

1. **多平台支持**：支持主流的自媒体平台如抖音、小红书，youtube等。
2. **统一接口**：通过标准化的接口，实现对多个平台的内容发布。
3. **开源与闭源结合**：Pro支持插件化架构，开源插件可自由扩展，闭源插件提供增强功能。
4. **跨语言支持**：Pro通过 gRPC 实现跨语言调用，支持 Go 和其他语言集成。
5. **可扩展性**：灵活的 Provider 机制，轻松接入新平台。

## MediaX快速开始

### 环境要求

- Go 1.18 或更高版本
- Git

### 安装

1. 初始化项目并下载 MediaX：

   ```bash
   go get github.com/ArtisanCloud/MediaX
   ```

   2. 创建一个简单的示例，本项目作者正在自己系统中使用，陆续会迭代版本：

      ```go
   
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
   
      // 从MediaXClient实例中获取到微信平台中公众号到实例
      wechatOAClient, err := mediaXClient.MediaXClient.CreateWechatOfficialAccount(&config2.WeChatOfficialAccountConfig{
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
      var content = ...
      resData, err := oaClient.GetPublishClient().DraftAdd(ctx, &reqData)
      if err != nil {
         return nil, err
      }
   
      ```

[//]: # (## 更多文档介绍)

[//]: # ([文档]&#40;./docs&#41;)

功能矩阵
平台 | 应用 | 授权类型 | 文章 |  |  |  | 视频 |  |  |  | 媒体管理
|  |  | 新增 | 删除 | 修改 | 查询 | 新增 | 删除 | 修改 | 查询 |
google | youtube | OAuth2 | ✔ |  |  | ✔ | ✔ |  | ✔ | ✔ | ✔
wechat | official account | 网页授权 | ✔ | ✔ | ✔ | ✔ |  |  |  |  | ✔

注：✔ 表示支持该功能，✘ 表示暂不支持该功能。

## 功能介绍
*** [项目功能的开发安排](https://github.com/orgs/ArtisanCloud/projects/5/views/2)



## 许可证

MediaX SDK项目 采用 [MIT License](./LICENSE) 开源。

