package xbot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"sync"
	"syscall"
	"time"

	"xbot/api"
	"xbot/config"
	"xbot/driver"
	"xbot/event"
	"xbot/logger"
	"xbot/session"
	"xbot/storage"

	"github.com/redis/go-redis/v9"
)

// Config 机器人配置
type Config struct {
	Nickname      []string
	SuperUsers    []int64
	CommandPrefix string
	Drivers       []driver.Driver
	DriverConfigs []config.DriverConfig // 保存原始驱动器配置
	Redis         *redis.Client
	Storage       storage.Storage
}

// Bot 机器人实例
type Bot struct {
	SelfID         int64
	Config         *Config
	API            *api.Client
	engines        []*Engine
	Storage        storage.Storage
	SessionManager *session.Manager
}

// BotManager 机器人管理器
type BotManager struct {
	config        *Config
	bots          sync.Map // selfID -> *Bot
	drivers       []driver.Driver
	driverConfigs []config.DriverConfig // 保存驱动器配置用于重试
	storage       storage.Storage
}

// Run 运行机器人
func Run(cfg *Config) (*BotManager, error) {
	manager := &BotManager{
		config:        cfg,
		drivers:       cfg.Drivers,
		driverConfigs: cfg.DriverConfigs,
		storage:       cfg.Storage,
	}

	// 如果没有提供存储，使用默认的内存存储
	if manager.storage == nil {
		manager.storage = storage.NewMemoryStorage()
	}

	// 设置事件处理器并连接驱动器
	for i, d := range manager.drivers {
		d.SetEventHandler(manager.handleEvent)

		// 获取驱动器配置（如果有的话）
		var drvCfg *config.DriverConfig
		if i < len(manager.driverConfigs) {
			drvCfg = &manager.driverConfigs[i]
		}

		// 连接驱动器（带重试）
		if err := manager.connectWithRetry(d, i, drvCfg); err != nil {
			logger.Error("驱动器连接失败，已达到最大重试次数", "driver", i, "error", err)
			continue
		}
	}

	logger.Info("机器人已启动")

	return manager, nil
}

// connectWithRetry 带重试的连接方法
func (bm *BotManager) connectWithRetry(d driver.Driver, driverIndex int, drvCfg *config.DriverConfig) error {
	// 默认值
	maxRetries := 5
	initialInterval := 2 * time.Second
	maxInterval := 30 * time.Second

	// 如果提供了配置，使用配置中的值
	if drvCfg != nil {
		if drvCfg.InitialRetries > 0 {
			maxRetries = drvCfg.InitialRetries
		}
		if drvCfg.InitialRetryInterval > 0 {
			initialInterval = time.Duration(drvCfg.InitialRetryInterval) * time.Second
		}
		if drvCfg.InitialMaxRetryInterval > 0 {
			maxInterval = time.Duration(drvCfg.InitialMaxRetryInterval) * time.Second
		}
	}

	var lastErr error
	retryInterval := initialInterval

	for attempt := 1; attempt <= maxRetries; attempt++ {
		logger.Info("正在连接驱动器",
			"driver", driverIndex,
			"attempt", attempt,
			"maxRetries", maxRetries,
		)

		err := d.Connect()
		if err == nil {
			logger.Info("驱动器连接成功", "driver", driverIndex, "attempt", attempt)
			return nil
		}

		lastErr = err
		logger.Warn("驱动器连接失败，准备重试",
			"driver", driverIndex,
			"attempt", attempt,
			"maxRetries", maxRetries,
			"error", err,
			"retryAfter", retryInterval,
		)

		// 如果还有重试机会，等待后重试
		if attempt < maxRetries {
			time.Sleep(retryInterval)

			// 使用指数退避策略，但不超过最大间隔
			retryInterval *= 2
			if retryInterval > maxInterval {
				retryInterval = maxInterval
			}
		}
	}

	return fmt.Errorf("连接失败，已重试 %d 次: %w", maxRetries, lastErr)
}

// RunAndListen 运行并阻塞
func RunAndListen(cfg *Config) error {
	manager, err := Run(cfg)
	if err != nil {
		return err
	}

	manager.Listen()
	return nil
}

// Listen 监听退出信号
func (bm *BotManager) Listen() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("机器人正在运行，按 Ctrl+C 退出")

	<-sigChan

	logger.Info("收到退出信号，正在关闭...")
	bm.Stop()
}

// Stop 停止所有驱动器
func (bm *BotManager) Stop() {
	// 关闭所有驱动器
	for _, d := range bm.drivers {
		if err := d.Close(); err != nil {
			logger.Error("关闭驱动器失败", "error", err)
		}
	}

	// 关闭存储
	if bm.storage != nil {
		if err := bm.storage.Close(); err != nil {
			logger.Error("关闭存储失败", "error", err)
		}
	}

	logger.Info("机器人已停止")
}

// GetBot 获取指定机器人
func (bm *BotManager) GetBot(selfID int64) (*Bot, bool) {
	if b, ok := bm.bots.Load(selfID); ok {
		return b.(*Bot), true
	}
	return nil, false
}

// GetAllBots 获取所有机器人
func (bm *BotManager) GetAllBots() []*Bot {
	var bots []*Bot

	bm.bots.Range(func(key, value interface{}) bool {
		bots = append(bots, value.(*Bot))
		return true
	})

	return bots
}

// handleEvent 处理事件
func (bm *BotManager) handleEvent(evt event.Event) {
	selfID := evt.GetSelfID()

	// 获取或创建 Bot 实例
	bot, ok := bm.GetBot(selfID)
	if !ok {
		bot = bm.createBot(selfID)
		bm.bots.Store(selfID, bot)
		logger.Info("创建新的 Bot 实例", "selfID", selfID)
	}

	// 记录消息日志（只记录一次）
	logMessageEvent(evt, bot)

	// 通知会话管理器
	if msgEvt, ok := evt.(*event.PrivateMessageEvent); ok {
		bot.SessionManager.NotifyWaitSession(msgEvt.UserID, 0, NewContext(evt, bot))
	} else if msgEvt, ok := evt.(*event.GroupMessageEvent); ok {
		bot.SessionManager.NotifyWaitSession(msgEvt.UserID, msgEvt.GroupID, NewContext(evt, bot))
	}

	// 分发到所有引擎
	for _, engine := range bot.engines {
		go engine.HandleEvent(evt)
	}
}

// createBot 创建 Bot 实例
func (bm *BotManager) createBot(selfID int64) *Bot {
	// 为每个 Bot 创建 API 客户端
	var apiClient *api.Client
	if len(bm.drivers) > 0 {
		apiClient = api.NewClient(bm.drivers[0])
	}

	// 创建会话管理器
	sessionStore := session.NewMemoryStore()
	sessionManager := session.NewManager(sessionStore, 5*time.Minute)

	bot := &Bot{
		SelfID:         selfID,
		Config:         bm.config,
		API:            apiClient,
		engines:        GetEngines(),
		Storage:        bm.storage,
		SessionManager: sessionManager,
	}

	// 设置引擎的 Bot 引用
	for _, engine := range bot.engines {
		engine.SetBot(bot)
	}

	return bot
}

// ensureDirectories 确保必要的目录存在
func ensureDirectories() error {
	// 需要创建的目录列表
	dirs := []string{
		"data",
		"logs",
		"config",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
		logger.Debug("确保目录存在", "directory", dir)
	}

	return nil
}

// LoadConfigFile 从文件加载配置
func LoadConfigFile(path string) (*Config, error) {
	// 首先确保必要的目录存在
	if err := ensureDirectories(); err != nil {
		return nil, err
	}

	cfg, err := config.LoadConfig(path)
	if err != nil {
		return nil, err
	}

	// 转换为 Bot Config
	botCfg := &Config{
		Nickname:      cfg.Bot.Nickname,
		SuperUsers:    cfg.Bot.SuperUsers,
		CommandPrefix: cfg.Bot.CommandPrefix,
	}

	// 设置日志级别
	logger.SetLevel(logger.ParseLevel(cfg.Log.Level))

	// 如果配置了日志文件，同时输出到文件
	if cfg.Log.File != "" {
		file := logger.MustCreateFile(cfg.Log.File)
		multiWriter := logger.NewMultiWriter(os.Stdout, file)
		logger.SetDefault(logger.NewDefaultLogger(multiWriter, logger.ParseLevel(cfg.Log.Level)))
	}

	// 创建存储
	if cfg.Storage.Type == "leveldb" {
		storePath := filepath.Join("data", "storage")
		if _, err := os.Stat(storePath); os.IsNotExist(err) {
			if err := os.MkdirAll(storePath, 0755); err != nil {
				return nil, fmt.Errorf("创建存储目录失败: %w", err)
			}
		}
		db, err := storage.NewLevelDB(storePath)
		if err != nil {
			return nil, fmt.Errorf("创建 LevelDB 失败: %w", err)
		}
		botCfg.Storage = db
	} else {
		botCfg.Storage = storage.NewMemoryStorage()
	}

	// 创建 Redis 客户端
	if cfg.Redis.Enabled {
		botCfg.Redis = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := botCfg.Redis.Ping(ctx).Err(); err != nil {
			logger.Warn("Redis 连接失败", "error", err)
			botCfg.Redis = nil
			os.Exit(0)
		} else {
			logger.Info("Redis 连接成功")
		}
	}

	// 创建驱动器
	for _, drvCfg := range cfg.Drivers {
		var drv driver.Driver

		switch drvCfg.Type {
		case "ws_reverse":
			drv = driver.NewWSReverseDriver(driver.Config{
				Type:              drvCfg.Type,
				URL:               drvCfg.URL,
				AccessToken:       drvCfg.AccessToken,
				ReconnectInterval: drvCfg.ReconnectInterval,
				MaxReconnect:      drvCfg.MaxReconnect,
				Timeout:           drvCfg.Timeout,
			})
		case "ws", "websocket":
			drv = driver.NewWebSocketDriver(driver.Config{
				Type:              drvCfg.Type,
				URL:               drvCfg.URL,
				Host:              drvCfg.Host,
				Port:              drvCfg.Port,
				AccessToken:       drvCfg.AccessToken,
				ReconnectInterval: drvCfg.ReconnectInterval,
				MaxReconnect:      drvCfg.MaxReconnect,
				HeartbeatInterval: drvCfg.HeartbeatInterval,
				Timeout:           drvCfg.Timeout,
			})
		case "http":
			drv = driver.NewHTTPDriver(driver.Config{
				Type:        drvCfg.Type,
				URL:         drvCfg.URL,
				Host:        drvCfg.Host,
				Port:        drvCfg.Port,
				AccessToken: drvCfg.AccessToken,
				Timeout:     drvCfg.Timeout,
			})
		case "http_post":
			drv = driver.NewHTTPPostDriver(driver.Config{
				Type:        drvCfg.Type,
				URL:         drvCfg.URL,
				Host:        drvCfg.Host,
				Port:        drvCfg.Port,
				AccessToken: drvCfg.AccessToken,
				Timeout:     drvCfg.Timeout,
			})
		default:
			logger.Warn("不支持的驱动器类型", "type", drvCfg.Type)
			continue
		}

		if drv != nil {
			botCfg.Drivers = append(botCfg.Drivers, drv)
			botCfg.DriverConfigs = append(botCfg.DriverConfigs, drvCfg)
		}
	}

	return botCfg, nil
}

// GetStorage 获取插件专用存储
func GetStorage(pluginName string) storage.Storage {
	// 创建插件数据目录
	dataPath := filepath.Join("data", pluginName)
	os.MkdirAll(dataPath, 0755)

	// 创建 LevelDB 存储
	db, err := storage.NewLevelDB(dataPath)
	if err != nil {
		logger.Error("创建插件存储失败", "plugin", pluginName, "error", err)
		return storage.NewMemoryStorage()
	}

	return db
}

// logMessageEvent 记录消息事件详细日志
func logMessageEvent(evt event.Event, bot *Bot) {
	switch e := evt.(type) {
	case *event.PrivateMessageEvent:
		// 记录私聊消息
		logger.Info("收到私聊消息",
			"发送者", e.Sender.Nickname,
			"用户ID", e.UserID,
			"消息类型", "私聊",
			"子类型", e.SubType,
			"消息内容", e.ParsedMessage.GetPlainText(),
			"消息ID", e.MessageID,
		)
	case *event.GroupMessageEvent:
		// 记录群消息
		logger.Info("收到群消息",
			"发送者", e.Sender.Nickname,
			"用户ID", e.UserID,
			"群号", e.GroupID,
			"消息类型", "群聊",
			"角色", e.Sender.Role,
			"消息内容", e.ParsedMessage.GetRawMessage(),
			"消息ID", e.MessageID,
		)
	case *event.HeartbeatMetaEvent:
		logger.Info("收到心跳事件", "心跳类型", e.MetaEventType, "心跳数据", e.Status)
	case *event.LifecycleMetaEvent:
		logger.Info("收到生命周期事件", "生命周期类型", e.SubType, "生命周期数据", e.MetaEventType)
	default:
		if evt.GetPostType() == "notice" {
			return
		}
		logger.Info("收到未知事件", "事件类型", reflect.TypeOf(evt), evt.GetPostType())
	}
}
