package plugin

import (
	"context"
	"fmt"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger"
	"github.com/ArtisanCloud/MediaXCore/pkg/logger/config"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"sync"

	plugin2 "github.com/ArtisanCloud/MediaXCore/pkg/plugin"
	"github.com/ArtisanCloud/MediaXCore/pkg/plugin/core"
	"github.com/ArtisanCloud/MediaXCore/pkg/plugin/core/contract"
)

type PluginManager struct {
	PluginRegistration map[string]*core.PluginMetadata
	Plugins            map[string]contract.ProviderInterface
	mu                 sync.Mutex
	Logger             *logger.Logger
}

func NewPluginManager(logConfig *config.LogConfig) *PluginManager {
	return &PluginManager{
		PluginRegistration: make(map[string]*core.PluginMetadata),
		Plugins:            make(map[string]contract.ProviderInterface),
		Logger:             logger.NewLogger(logConfig),
	}
}

// Register 插件注册
func (pm *PluginManager) Register(pluginMetadata *core.PluginMetadata) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	name := pluginMetadata.Name

	// 避免重复注册
	if _, exists := pm.PluginRegistration[name]; exists {
		pm.Logger.InfoF("Plugin %s already registered with path: %s\n", name, pluginMetadata.BuildPath)
		return
	}

	pm.PluginRegistration[name] = pluginMetadata
	pm.Logger.InfoF("Plugin %s registered with path: %s\n", name, pluginMetadata.BuildPath)
}

// LoadPlugins 扫描所有插件
func (pm *PluginManager) ScanPlugins() error {
	if err := pm.ScanPluginsFromDirectory(core.Open); err != nil {
		return err
	}
	if err := pm.ScanPluginsFromDirectory(core.Closed); err != nil {
		return err
	}
	return nil
}

// 从指定目录加载插件
// 从指定目录加载插件
func (pm *PluginManager) ScanPluginsFromDirectory(pluginType core.PluginType) error {
	// 获取插件目录的绝对路径
	pluginDir := filepath.Join("./plugins", string(pluginType))
	//pm.Logger.Info(pluginDir)
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
		//pm.Logger.Info(pluginFolderPath)
		// 查找该文件夹内的 plugin.yaml 文件
		pluginYamlPath := filepath.Join(pluginFolderPath, "plugin/plugin.yaml")
		//pm.Logger.Info(pluginYamlPath)
		if _, err := os.Stat(pluginYamlPath); os.IsNotExist(err) {
			// 如果没有 plugin.yaml 文件，跳过此文件夹
			continue
		}
		//pm.Logger.Info(pluginYamlPath)
		if err := pm.RegisterPluginFromFile(pluginFolderPath, pluginYamlPath); err != nil {
			return fmt.Errorf("failed to load open plugin %s: %v", pluginYamlPath, err)
		}
	}

	return nil
}

// RegisterPluginFromFile 加载插件，支持从 JSON 描述文件加载
func (pm *PluginManager) RegisterPluginFromFile(pluginDir string, pluginYamlFilePath string) error {
	// 读取插件元数据
	pluginMetadata, err := plugin2.ReadPluginMetadata(pluginYamlFilePath)
	//fmt2.Dump(pluginMetadata)
	if err != nil {
		return err
	}

	switch pluginMetadata.Type {
	case core.Open:
		return pm.registerOpenPlugin(pluginDir, pluginMetadata)
	case core.Closed:
		return pm.registerClosedPlugin(pluginDir, pluginMetadata)
	case core.External:
		return pm.registerExternalPlugin(pluginDir, pluginMetadata) // 未来扩展
	default:
		return fmt.Errorf("unknown plugin type %s", pluginMetadata.Type)
	}
}

// 编译开源插件到 .so 文件
func (pm *PluginManager) compileOpenPluginToSO(pluginDir string, metadata *core.PluginMetadata) error {
	pluginSourceDir := filepath.Join(pluginDir, metadata.SourcePath)
	pm.Logger.InfoF("Compiling open plugin from directory: %s\n", pluginDir)

	pluginSOPath := filepath.Join(pluginDir, metadata.BuildPath)
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginSOPath, pluginSourceDir)

	// 增加更多的日志输出和错误处理
	output, err := cmd.CombinedOutput() // 获取标准输出和错误输出
	if err != nil {
		return fmt.Errorf("failed to compile open plugin: %v, output: %s", err, output)
	}
	pm.Logger.InfoF("Success to build open plugin to : %s\n", pluginSOPath)

	return nil
}

// 加载开源插件
func (pm *PluginManager) registerOpenPlugin(pluginDir string, metadata *core.PluginMetadata) error {

	// 如果插件没有编译成 .so 文件，尝试编译
	pluginSOPath := filepath.Join(pluginDir, metadata.BuildPath)
	if _, err := os.Stat(pluginSOPath); os.IsNotExist(err) {
		// 尝试编译插件
		if err := pm.compileOpenPluginToSO(pluginDir, metadata); err != nil {
			return fmt.Errorf("failed to compile open plugin %s: %v", metadata.Name, err)
		}
	}

	// 现在尝试加载插件
	return pm.registerClosedPlugin(pluginDir, metadata)
}

// 加载闭源插件（.so 文件）
func (pm *PluginManager) registerClosedPlugin(pluginDir string, metadata *core.PluginMetadata) error {
	// 更新当前元数据组件的工作目录
	metadata.WorkDir = pluginDir
	// 注册插件到管理器中
	pm.Register(metadata)
	return nil
}

func (pm *PluginManager) registerExternalPlugin(pluginDir string, metadata *core.PluginMetadata) error {
	pm.Logger.InfoF("Registering external plugin: %s\n in path %s\n", metadata.Name, pluginDir)
	return nil
}

// 获取插件
func (pm *PluginManager) GetPlugin(name string, appId string, config *contract.PluginConfig) (contract.ProviderInterface, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if appId == "" {
		return nil, fmt.Errorf("appId cannot be empty")
	}
	pluginKey := fmt.Sprintf("%s.%s", name, appId)

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 1. 先检查是否已经加载
	if provider, exists := pm.Plugins[pluginKey]; exists {
		return provider, nil
	}

	// 2. 未加载，检查是否注册过路径
	pluginMetadata, exists := pm.PluginRegistration[name]
	if !exists {
		return nil, fmt.Errorf("plugin metadata %s not found", name)
	}

	// 3. 加载插件
	pluginSOPath := filepath.Join(pluginMetadata.WorkDir, pluginMetadata.BuildPath)
	pm.Logger.InfoF("Lazy loading plugin: %s from path: %s\n", name, pluginSOPath)
	p, err := plugin.Open(pluginSOPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin file %s: %v", pluginSOPath, err)
	}

	provider, err := plugin2.LookUpSymbol[contract.ProviderInterface](p, name)
	if err != nil {
		return nil, fmt.Errorf("failed to look up plugin %s: %v", name, err)
	}

	// 4. 加载成功，存入缓存
	pm.Plugins[pluginKey] = *provider

	if config != nil {
		ctx := context.Background()
		err = pm.Plugins[pluginKey].Initialize(&ctx, config)
		if err != nil {
			return nil, err
		}
	}

	return *provider, nil
}
