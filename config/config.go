package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// BotConfig 机器人配置
type BotConfig struct {
	Bot struct {
		Nickname      []string `yaml:"nickname"`
		SuperUsers    []int64  `yaml:"super_users"`
		CommandPrefix string   `yaml:"command_prefix"`
	} `yaml:"bot"`

	Drivers []DriverConfig `yaml:"drivers"`

	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		Enabled  bool   `yaml:"enabled"`
	} `yaml:"redis"`

	Log struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"log"`

	Storage struct {
		Type string `yaml:"type"` // leveldb or memory
		Path string `yaml:"path"`
	} `yaml:"storage"`
}

// DriverConfig 驱动器配置
type DriverConfig struct {
	Type        string `yaml:"type"` // ws_reverse, ws_forward, http, http_post
	URL         string `yaml:"url"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	AccessToken string `yaml:"access_token"`

	// WebSocket 配置
	ReconnectInterval int `yaml:"reconnect_interval"` // 重连间隔（秒）
	MaxReconnect      int `yaml:"max_reconnect"`      // 最大重连次数，0 表示无限
	HeartbeatInterval int `yaml:"heartbeat_interval"` // 心跳间隔（秒）

	// 初始连接重试配置
	InitialRetries          int `yaml:"initial_retries"`            // 初始连接最大重试次数，默认 5
	InitialRetryInterval    int `yaml:"initial_retry_interval"`     // 初始连接重试间隔（秒），默认 2
	InitialMaxRetryInterval int `yaml:"initial_max_retry_interval"` // 初始连接最大重试间隔（秒），默认 30

	// API 配置
	Timeout int `yaml:"timeout"` // API 调用超时（秒）
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*BotConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config BotConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// 设置默认值
	setDefaults(&config)

	return &config, nil
}

// setDefaults 设置默认值
func setDefaults(config *BotConfig) {
	// Bot 默认值
	if config.Bot.CommandPrefix == "" {
		config.Bot.CommandPrefix = "/"
	}

	// 日志默认值
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}

	// 存储默认值
	if config.Storage.Type == "" {
		config.Storage.Type = "leveldb"
	}
	if config.Storage.Path == "" {
		config.Storage.Path = "data"
	}

	// 驱动器默认值
	for i := range config.Drivers {
		if config.Drivers[i].ReconnectInterval == 0 {
			config.Drivers[i].ReconnectInterval = 5
		}
		if config.Drivers[i].Timeout == 0 {
			config.Drivers[i].Timeout = 30
		}
		if config.Drivers[i].HeartbeatInterval == 0 {
			config.Drivers[i].HeartbeatInterval = 30
		}
		// 初始连接重试默认值
		if config.Drivers[i].InitialRetries == 0 {
			config.Drivers[i].InitialRetries = 5
		}
		if config.Drivers[i].InitialRetryInterval == 0 {
			config.Drivers[i].InitialRetryInterval = 2
		}
		if config.Drivers[i].InitialMaxRetryInterval == 0 {
			config.Drivers[i].InitialMaxRetryInterval = 30
		}
	}
}

// SaveConfig 保存配置到文件
func SaveConfig(path string, config *BotConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Watch 监听配置文件变化
func (c *BotConfig) Watch(path string, callback func(*BotConfig)) error {
	watcher, err := NewWatcher(path, func() {
		newConfig, err := LoadConfig(path)
		if err == nil {
			callback(newConfig)
		}
	})
	if err != nil {
		return err
	}

	return watcher.Start()
}
