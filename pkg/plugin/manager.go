package plugin

import (
	"fmt"
	plugin2 "github.com/ArtisanCloud/MediaXCore/pkg/plugin"
	"github.com/ArtisanCloud/MediaXCore/pkg/plugin/core"
	"github.com/ArtisanCloud/MediaXCore/pkg/plugin/core/contract"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"sync"
)

type PluginManager struct {
	Plugins map[string]contract.ProviderInterface
	mu      sync.Mutex
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		Plugins: make(map[string]contract.ProviderInterface),
	}
}

// Register 插件注册
func (pm *PluginManager) Register(provider contract.ProviderInterface) {
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
	if err := pm.LoadPluginsFromDirectory(core.Open); err != nil {
		return err
	}
	if err := pm.LoadPluginsFromDirectory(core.Closed); err != nil {
		return err
	}
	return nil
}

// 从指定目录加载插件
// 从指定目录加载插件
func (pm *PluginManager) LoadPluginsFromDirectory(pluginType core.PluginType) error {
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

		// 查找该文件夹内的 plugin.yaml 文件
		pluginJsonPath := filepath.Join(pluginFolderPath, "plugin.yaml")
		if _, err := os.Stat(pluginJsonPath); os.IsNotExist(err) {
			// 如果没有 plugin.yaml 文件，跳过此文件夹
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
	pluginMetadata, err := plugin2.ReadPluginMetadata(pluginFilePath)
	if err != nil {
		return err
	}

	switch pluginMetadata.Type {
	case core.Open:
		return pm.loadOpenPlugin(pluginMetadata)
	case core.Closed:
		return pm.loadClosedPlugin(pluginMetadata)
	case core.External:
		return pm.loadExternalPlugin(pluginMetadata) // 未来扩展
	default:
		return fmt.Errorf("unknown plugin type %s", pluginMetadata.Type)
	}
}

// 编译开源插件到 .so 文件
func (pm *PluginManager) compileOpenPluginToSO(metadata *core.PluginMetadata) error {
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
func (pm *PluginManager) loadOpenPlugin(metadata *core.PluginMetadata) error {
	// 读取插件描述文件以获取更多信息
	pluginMetadata, err := plugin2.ReadPluginMetadata(metadata.Path)
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
func (pm *PluginManager) loadClosedPlugin(metadata *core.PluginMetadata) error {
	pluginPath := metadata.Path
	fmt.Printf("Loading plugin from path: %s\n", pluginPath)

	// 确保路径是 .so 文件
	if !plugin2.IsSOFile(pluginPath) {
		return fmt.Errorf("plugin file %s is not a valid .so file", pluginPath)
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin file %s: %v", pluginPath, err)
	}

	provider, err := plugin2.LookUpSymbol[contract.ProviderInterface](p, metadata.Name)
	if err != nil {
		return fmt.Errorf("failed to look up plugin  %s: %v", metadata.Name, err)
	}

	pm.Register(*provider)
	return nil
}

func (pm *PluginManager) loadExternalPlugin(metadata *core.PluginMetadata) error {
	return nil
}

// 获取插件
func (pm *PluginManager) GetPlugin(name string) (contract.ProviderInterface, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	provider, exists := pm.Plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	return provider, nil
}
