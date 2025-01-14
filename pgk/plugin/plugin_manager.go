// Path: /pgk/plugin/plugin_manager.go
package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"
)

type PluginManager struct {
	Plugins map[string]Provider
	mu      sync.Mutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		Plugins: make(map[string]Provider),
	}
}

// 插件描述文件结构
type PluginMetadata struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Type    string                 `json:"type"`
	Path    string                 `json:"path"`
	Config  map[string]interface{} `json:"config"`
}

// Register 插件注册
func (pm *PluginManager) Register(provider Provider) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.Plugins[provider.Name()]; exists {
		fmt.Printf("Provider %s already registered.\n", provider.Name())
		return
	}
	pm.Plugins[provider.Name()] = provider
	fmt.Printf("Provider %s registered successfully.\n", provider.Name())
}

// LoadPluginFromFile 加载插件，支持从 JSON 描述文件加载
func (pm *PluginManager) LoadPluginFromFile(pluginFilePath string) error {
	// 读取插件描述文件
	pluginMetadata, err := pm.readPluginMetadata(pluginFilePath)
	if err != nil {
		return err
	}

	// 根据插件类型选择加载路径
	if pluginMetadata.Type == "open" {
		return pm.loadOpenPlugin(pluginMetadata)
	} else if pluginMetadata.Type == "closed" {
		return pm.loadClosedPlugin(pluginMetadata)
	} else {
		return errors.New("unknown plugin type")
	}
}

// 读取插件描述文件
func (pm *PluginManager) readPluginMetadata(pluginFilePath string) (*PluginMetadata, error) {
	file, err := os.Open(pluginFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin description file %s: %v", pluginFilePath, err)
	}
	defer file.Close()

	metadata := &PluginMetadata{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to decode plugin description file %s: %v", pluginFilePath, err)
	}

	return metadata, nil
}

// 加载开源插件
func (pm *PluginManager) loadOpenPlugin(metadata *PluginMetadata) error {
	// 假设 open 插件是一个 Go 包，路径可以通过 GOPATH 或 go.mod 管理
	// 这里只是一个示例，加载方式依赖于具体的构建和运行方式
	// 在实际场景中，可能需要通过插件系统来动态加载 Go 包（例如通过反射或自定义加载方式）

	// 在此示例中，我们只假设路径是 Go 包路径
	fmt.Printf("Loading open plugin from path: %s\n", metadata.Path)

	// 这里假设已经使用 go.mod 进行管理，且插件 Go 包路径已经被导入到项目中
	// 注册插件提供者
	var provider Provider
	// 获取 provider 实例的方式视具体实现而定
	pm.Register(provider)
	return nil
}

// 加载闭源插件
func (pm *PluginManager) loadClosedPlugin(metadata *PluginMetadata) error {
	pluginPath := filepath.Join("closed", metadata.Path)

	// 打开并加载 .so 插件文件
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin file %s: %v", pluginPath, err)
	}

	// 查找 Provider 符号
	symProvider, err := p.Lookup("Provider")
	if err != nil {
		return fmt.Errorf("symbol 'Provider' not found in plugin file %s: %v", pluginPath, err)
	}

	// 类型断言为 Provider
	provider, ok := symProvider.(Provider)
	if !ok {
		return fmt.Errorf("symbol 'Provider' in plugin file %s is not of type 'Provider'", pluginPath)
	}

	// 注册插件
	pm.Register(provider)
	return nil
}

// GetPlugin 获取插件
func (pm *PluginManager) GetPlugin(name string) (Provider, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	provider, exists := pm.Plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	return provider, nil
}
