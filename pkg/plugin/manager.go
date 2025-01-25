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
	//println(pluginDir)
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
		//println(pluginFolderPath)
		// 查找该文件夹内的 plugin.yaml 文件
		pluginYamlPath := filepath.Join(pluginFolderPath, "plugin/plugin.yaml")
		//println(pluginYamlPath)
		if _, err := os.Stat(pluginYamlPath); os.IsNotExist(err) {
			// 如果没有 plugin.yaml 文件，跳过此文件夹
			continue
		}
		//println(pluginYamlPath)
		if err := pm.LoadPluginFromFile(pluginFolderPath, pluginYamlPath); err != nil {
			return fmt.Errorf("failed to load open plugin %s: %v", pluginYamlPath, err)
		}
	}

	return nil
}

// LoadPluginFromFile 加载插件，支持从 JSON 描述文件加载
func (pm *PluginManager) LoadPluginFromFile(pluginDir string, pluginYamlFilePath string) error {
	pluginMetadata, err := plugin2.ReadPluginMetadata(pluginYamlFilePath)
	//fmt2.Dump(pluginMetadata)
	if err != nil {
		return err
	}

	switch pluginMetadata.Type {
	case core.Open:
		return pm.loadOpenPlugin(pluginDir, pluginMetadata)
	case core.Closed:
		return pm.loadClosedPlugin(pluginDir, pluginMetadata)
	case core.External:
		return pm.loadExternalPlugin(pluginDir, pluginMetadata) // 未来扩展
	default:
		return fmt.Errorf("unknown plugin type %s", pluginMetadata.Type)
	}
}

// 编译开源插件到 .so 文件
func (pm *PluginManager) compileOpenPluginToSO(pluginDir string, metadata *core.PluginMetadata) error {
	pluginSourceDir := filepath.Join(pluginDir, metadata.SourcePath)
	fmt.Printf("Compiling open plugin from directory: %s\n", pluginDir)

	pluginSOPath := filepath.Join(pluginDir, metadata.BuildPath)
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginSOPath, pluginSourceDir)

	// 增加更多的日志输出和错误处理
	output, err := cmd.CombinedOutput() // 获取标准输出和错误输出
	if err != nil {
		return fmt.Errorf("failed to compile open plugin: %v, output: %s", err, output)
	}
	fmt.Printf("Success to build open plugin to : %s\n", pluginSOPath)

	return nil
}

// 加载开源插件
func (pm *PluginManager) loadOpenPlugin(pluginDir string, metadata *core.PluginMetadata) error {

	// 如果插件没有编译成 .so 文件，尝试编译
	pluginSOPath := filepath.Join(pluginDir, metadata.BuildPath)
	if _, err := os.Stat(pluginSOPath); os.IsNotExist(err) {
		// 尝试编译插件
		if err := pm.compileOpenPluginToSO(pluginDir, metadata); err != nil {
			return fmt.Errorf("failed to compile open plugin %s: %v", metadata.Name, err)
		}
	}

	// 现在尝试加载插件
	return pm.loadClosedPlugin(pluginSOPath, metadata)
}

// 加载闭源插件（.so 文件）
func (pm *PluginManager) loadClosedPlugin(pluginSOPath string, metadata *core.PluginMetadata) error {

	fmt.Printf("Loading plugin from path: %s\n", pluginSOPath)

	// 确保路径是 .so 文件
	if !plugin2.IsSOFile(pluginSOPath) {
		return fmt.Errorf("plugin file %s is not a valid .so file", pluginSOPath)
	}

	p, err := plugin.Open(pluginSOPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin file %s: %v", pluginSOPath, err)
	}

	provider, err := plugin2.LookUpSymbol[contract.ProviderInterface](p, metadata.Name)
	if err != nil {
		return fmt.Errorf("failed to look up plugin  %s: %v", metadata.Name, err)
	}

	pm.Register(*provider)
	return nil
}

func (pm *PluginManager) loadExternalPlugin(pluginDir string, metadata *core.PluginMetadata) error {
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
