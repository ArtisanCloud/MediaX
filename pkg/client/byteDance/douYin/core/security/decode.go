package security

import (
	"crypto/sha1"
	"fmt"
)

// ## WebHook 解密文本内容 Demo
// 参考文档：https://developer.open-douyin.com/docs/resource/zh-CN/dop/develop/webhooks/summarize
// ```go
//	func Demo() string {
//		body := "{\"event\":\"verify_webhook\",\"client_key\":\"abc\",\"content\":{\"challenge\":12345}}" // 接收到的完整请求体
//		clientSecret := "***" // 替换成自己的 client_secret
//		return sign(clientSecret, []byte(body))
//	}
// ```
// ## sign 计算签名
func Sign(clientSecret string, body []byte) string {
	h := sha1.New()
	h.Write([]byte(clientSecret))
	h.Write(body)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
