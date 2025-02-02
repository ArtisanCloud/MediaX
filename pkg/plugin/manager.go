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

func (pm *PluginManager) GetPluginKeyByMetadata(pluginMetadata *core.PluginMetadata) string {
	key := pluginMetadata.VendorName + "." + pluginMetadata.Name
	return key
}

func (pm *PluginManager) GetPluginAppKeyByMetadata(pluginMetadata *core.PluginMetadata, appId string) string {
	key := pm.GetPluginKeyByMetadata(pluginMetadata) + "." + appId
	return key
}

func (pm *PluginManager) GetPluginKey(vendorName contract.MediaVendor, pluginName contract.AppPlugin) string {
	key := string(vendorName) + "." + string(pluginName)
	return key
}
func (pm *PluginManager) GetPluginAppKey(vendorName contract.MediaVendor, pluginName contract.AppPlugin, appId string) string {
	key := pm.GetPluginKey(vendorName, pluginName) + "." + appId
	return key
}

// Register 插件注册
func (pm *PluginManager) Register(pluginMetadata *core.PluginMetadata) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	registeredName := pm.GetPluginKeyByMetadata(pluginMetadata)

	// 避免重复注册
	if _, exists := pm.PluginRegistration[registeredName]; exists {
		pm.Logger.InfoF("Plugin %s already registered with path: %s\n", registeredName, pluginMetadata.BuildPath)
		return
	}

	pm.PluginRegistration[registeredName] = pluginMetadata
	pm.Logger.InfoF("Plugin %s registered with path: %s\n", registeredName, pluginMetadata.BuildPath)
}

// LoadPlugins 扫描所有插件
func (pm *PluginManager) ScanPlugins() error {
	if err := pm.ScanPluginBundlesFromDirectory(core.Open); err != nil {
		return err
	}
	if err := pm.ScanPluginBundlesFromDirectory(core.Closed); err != nil {
		return err
	}
	return nil
}

// 从指定目录加载插件
func (pm *PluginManager) ScanPluginBundlesFromDirectory(pluginType core.PluginType) error {
	// 获取插件目录的绝对路径
	pluginDir := filepath.Join("./plugins", string(pluginType))
	//pm.Logger.InfoF("start to scan pluginDir: %s", pluginDir)

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
		//pm.Logger.InfoF("pluginFolderPath is: %s", pluginFolderPath)
		// 查找该文件夹内的 plugin.yaml 文件
		pluginYamlPath := filepath.Join(pluginFolderPath, "plugins/plugins.yaml")
		//pm.Logger.Info(pluginYamlPath)
		if _, err := os.Stat(pluginYamlPath); os.IsNotExist(err) {
			// 如果没有 plugin.yaml 文件，跳过此文件夹
			pm.Logger.WarnF("pluginYamlPath:%s does not exist", pluginYamlPath)
			continue
		}
		//pm.Logger.Info(pluginYamlPath)
		if err := pm.RegisterPluginBundleFromFile(pluginFolderPath, pluginYamlPath); err != nil {
			return fmt.Errorf("failed to load open plugin %s: %v", pluginYamlPath, err)
		}
	}

	return nil
}

// RegisterPluginBundleFromFile 加载插件，支持从 JSON 描述文件加载
func (pm *PluginManager) RegisterPluginBundleFromFile(pluginDir string, pluginYamlFilePath string) error {
	//pm.Logger.Info(pluginDir)
	//pm.Logger.Info(pluginYamlFilePath)
	// 读取插件元数据
	pluginsMetadata, err := plugin2.ReadPluginBundleMetadata(pluginYamlFilePath)
	//fmt2.Dump(pluginsMetadata)
	if err != nil {
		return err
	}

	if len(pluginsMetadata.Plugins) <= 0 {
		return fmt.Errorf("plugin %s doesn't have plugins", pluginsMetadata.Name)
	}
	if pluginsMetadata.Name == "" {
		return fmt.Errorf("plugin bundle path %s  doesn't have valid name", pluginDir)
	}

	for _, metadata := range pluginsMetadata.Plugins {
		//fmt2.Dump(metadata)
		// 更新当前元数据组件的工作目录
		metadata.WorkDir = pluginDir
		metadata.VendorName = pluginsMetadata.Name

		switch pluginsMetadata.Type {
		case core.Open:
			err = pm.registerOpenPlugin(pluginDir, &metadata)
			if err != nil {
				return err
			}
		case core.Closed:
			err = pm.registerClosedPlugin(pluginDir, &metadata)
			if err != nil {
				return err
			}
		case core.External:
			err = pm.registerExternalPlugin(pluginDir, &metadata) // 未来扩展
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown plugin type %s", pluginsMetadata.Type)
		}
	}
	return nil

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
	//pm.Logger.InfoF("register open plugin path: ", pluginSOPath)
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
	// 注册插件到管理器中
	pm.Register(metadata)
	return nil
}

func (pm *PluginManager) registerExternalPlugin(pluginDir string, metadata *core.PluginMetadata) error {
	pm.Logger.InfoF("Registering external plugin: %s\n in path %s\n", metadata.Name, pluginDir)
	return nil
}

// 获取插件
func (pm *PluginManager) GetPlugin(vendorName contract.MediaVendor, pluginName contract.AppPlugin, appId string, config *contract.PluginConfig) (contract.ProviderInterface, error) {
	if vendorName == "" {
		return nil, fmt.Errorf("vendorName cannot be empty")
	}
	if pluginName == "" {
		return nil, fmt.Errorf("pluginName cannot be empty")
	}
	if appId == "" {
		return nil, fmt.Errorf("appId cannot be empty")
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 1. 先检查是否已经加载
	pluginAppKey := pm.GetPluginAppKey(vendorName, pluginName, appId)
	if provider, exists := pm.Plugins[pluginAppKey]; exists {
		return provider, nil
	}

	// 2. 未加载，检查是否注册过路径
	pluginKey := pm.GetPluginKey(vendorName, pluginName)
	//fmt2.Dump(pm.PluginRegistration)
	pluginMetadata, exists := pm.PluginRegistration[pluginKey]
	if !exists {
		return nil, fmt.Errorf("plugin manager registerd plugin metadata %s not found", pluginKey)
	}

	// 3. 加载插件
	pluginSOPath := filepath.Join(pluginMetadata.WorkDir, pluginMetadata.BuildPath)
	pm.Logger.InfoF("Lazy loading plugin: %s from path: %s\n", pluginKey, pluginSOPath)
	p, err := plugin.Open(pluginSOPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin file %s: %v", pluginSOPath, err)
	}

	provider, err := plugin2.LookUpSymbol[contract.ProviderInterface](p, string(pluginName))
	if err != nil {
		return nil, fmt.Errorf("failed to look up plugin %s: %v", pluginAppKey, err)
	}

	// 4. 加载成功，存入缓存
	pm.Plugins[pluginAppKey] = *provider

	if config != nil {
		ctx := context.Background()
		err = pm.Plugins[pluginAppKey].Initialize(&ctx, config)
		if err != nil {
			return nil, err
		}
	}

	return *provider, nil
}
