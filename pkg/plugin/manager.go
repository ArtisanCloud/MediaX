package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/ArtisanCloud/MediaX/pkg/plugin/core"
)

type PluginType string

const (
	Open     PluginType = "open"
	Closed   PluginType = "closed"
	External PluginType = "external"
)

type PluginManager struct {
	Plugins map[string]core.Provider
	mu      sync.Mutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		Plugins: make(map[string]core.Provider),
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
func (pm *PluginManager) Register(provider core.Provider) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.Plugins[provider.Name()]; exists {
		fmt.Printf("Provider %s already registered.\n", provider.Name())
		return
	}
	pm.Plugins[provider.Name()] = provider
	fmt.Printf("Provider %s registered successfully.\n", provider.Name())
}

// LoadPlugins 加载所有插件
func (pm *PluginManager) LoadPlugins() error {
	if err := pm.LoadPluginsFromDirectory(Open); err != nil {
		return err
	}
	if err := pm.LoadPluginsFromDirectory(Closed); err != nil {
		return err
	}
	return nil
}

// 从指定目录加载插件
// 从指定目录加载插件
func (pm *PluginManager) LoadPluginsFromDirectory(pluginType PluginType) error {
	// 获取插件目录的绝对路径
	pluginDir := filepath.Join("./plugins", string(pluginType))
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory %s: %v", pluginDir, err)
	}

	// 遍历插件目录中的所有文件夹
	for _, file := range files {
		// 仅加载文件夹，不加载单独的文件
		if !file.IsDir() {
			continue
		}

		// 构造子文件夹路径
		pluginFolderPath := filepath.Join(pluginDir, file.Name())

		// 查找该文件夹内的 plugin.json 文件
		pluginJsonPath := filepath.Join(pluginFolderPath, "plugin.json")
		if _, err := os.Stat(pluginJsonPath); os.IsNotExist(err) {
			// 如果没有 plugin.json 文件，跳过此文件夹
			continue
		}

		if err := pm.LoadPluginFromFile(pluginJsonPath); err != nil {
			return fmt.Errorf("failed to load open plugin %s: %v", pluginJsonPath, err)
		}
	}

	return nil
}

// LoadPluginFromFile 加载插件，支持从 JSON 描述文件加载
func (pm *PluginManager) LoadPluginFromFile(pluginFilePath string) error {
	pluginMetadata, err := pm.readPluginMetadata(pluginFilePath)
	if err != nil {
		return err
	}

	switch pluginMetadata.Type {
	case string(Open):
		return pm.loadOpenPlugin(pluginMetadata)
	case string(Closed):
		return pm.loadClosedPlugin(pluginMetadata)
	case string(External):
		return pm.loadExternalPlugin(pluginMetadata) // 未来扩展
	default:
		return fmt.Errorf("unknown plugin type %s", pluginMetadata.Type)
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
	if err := decoder.Decode(metadata); err != nil {
		return nil, fmt.Errorf("failed to decode plugin description file %s: %v", pluginFilePath, err)
	}

	// 增加验证逻辑
	if metadata.Name == "" || metadata.Version == "" || metadata.Type == "" {
		return nil, fmt.Errorf("plugin metadata is incomplete or invalid in %s", pluginFilePath)
	}

	return metadata, nil
}

// 编译开源插件到 .so 文件
func (pm *PluginManager) compileOpenPluginToSO(metadata *PluginMetadata) error {
	pluginDir := metadata.Path
	fmt.Printf("Compiling open plugin from directory: %s\n", pluginDir)

	pluginSOPath := filepath.Join(pluginDir, "plugin.so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginSOPath, pluginDir)

	// 增加更多的日志输出和错误处理
	output, err := cmd.CombinedOutput() // 获取标准输出和错误输出
	if err != nil {
		return fmt.Errorf("failed to compile open plugin: %v, output: %s", err, output)
	}

	return nil
}

// 加载开源插件
func (pm *PluginManager) loadOpenPlugin(metadata *PluginMetadata) error {
	// 读取插件描述文件以获取更多信息
	pluginMetadata, err := pm.readPluginMetadata(metadata.Path)
	if err != nil {
		return fmt.Errorf("failed to read plugin metadata for open plugin %s: %v", metadata.Name, err)
	}

	// 如果插件没有编译成 .so 文件，尝试编译
	pluginSOPath := filepath.Join(metadata.Path, "plugin.so")
	if _, err := os.Stat(pluginSOPath); os.IsNotExist(err) {
		// 尝试编译插件
		if err := pm.compileOpenPluginToSO(pluginMetadata); err != nil {
			return fmt.Errorf("failed to compile open plugin %s: %v", metadata.Name, err)
		}
	}

	// 现在尝试加载插件
	return pm.loadClosedPlugin(pluginMetadata)
}

// 加载闭源插件（.so 文件）
func (pm *PluginManager) loadClosedPlugin(metadata *PluginMetadata) error {
	pluginPath := metadata.Path
	fmt.Printf("Loading plugin from path: %s\n", pluginPath)

	// 确保路径是 .so 文件
	if !isSOFile(pluginPath) {
		return fmt.Errorf("plugin file %s is not a valid .so file", pluginPath)
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin file %s: %v", pluginPath, err)
	}

	symProvider, err := p.Lookup("Provider")
	if err != nil {
		return fmt.Errorf("symbol 'Provider' not found in plugin file %s: %v", pluginPath, err)
	}

	provider, ok := symProvider.(core.Provider)
	if !ok {
		return fmt.Errorf("symbol 'Provider' in plugin file %s is not of type 'core.Provider'", pluginPath)
	}

	pm.Register(provider)
	return nil
}

func (pm *PluginManager) loadExternalPlugin() error {
	return nil
}

// 获取插件
func (pm *PluginManager) GetPlugin(name string) (core.Provider, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	provider, exists := pm.Plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	return provider, nil
}

// 检查是否为 .so 文件
func isSOFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".so" || ext == ".dylib" || ext == ".dll" // 扩展支持
}
