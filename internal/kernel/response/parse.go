package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ParseResponseToObject 解析 HTTP 响应的 body 到目标对象，并将响应体写回到 res.Body
func ParseResponseToObject(res *http.Response, obj interface{}) error {
	// 确保响应体被正确关闭
	defer res.Body.Close()

	// 创建一个缓冲区用于存储读取的数据
	var buf bytes.Buffer

	// 使用 io.TeeReader 既读取响应体数据到 buf 中，又将其传递给原始的 res.Body 以保持其内容不变
	teeReader := io.TeeReader(res.Body, &buf)

	// 读取响应体的内容
	bodyBytes, err := io.ReadAll(teeReader)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// 打印响应内容 (可选)
	fmt.Printf("Response Body: %s\n", string(bodyBytes))

	// 将响应体内容解析到目标对象
	err = json.Unmarshal(bodyBytes, obj)
	if err != nil {
		return fmt.Errorf("error unmarshaling response: %v", err)
	}

	// 将缓冲区的数据重新设置回 res.Body，这样外部调用不会受到影响
	res.Body = io.NopCloser(&buf)

	return nil
}
