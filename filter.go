package xbot

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"sync"

	"xbot/event"
	"xbot/logger"
	"xbot/message"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"
)

// Filter 过滤器类型
type Filter func(ctx *Context) bool

// OnlyPrivateMessage 只处理私聊消息
func OnlyPrivateMessage() Filter {
	return func(ctx *Context) bool {
		_, ok := ctx.Event.(*event.PrivateMessageEvent)
		return ok
	}
}

// OnlyGroupMessage 只处理群消息
func OnlyGroupMessage() Filter {
	return func(ctx *Context) bool {
		_, ok := ctx.Event.(*event.GroupMessageEvent)
		return ok
	}
}

// OnlyUsers 只处理指定用户的消息
func OnlyUsers(userIDs ...int64) Filter {
	userMap := make(map[int64]bool)
	for _, id := range userIDs {
		userMap[id] = true
	}

	return func(ctx *Context) bool {
		userID := ctx.GetUserID()
		return userMap[userID]
	}
}

// OnlyGroups 只处理指定群的消息
func OnlyGroups(groupIDs ...int64) Filter {
	groupMap := make(map[int64]bool)
	for _, id := range groupIDs {
		groupMap[id] = true
	}

	return func(ctx *Context) bool {
		groupID := ctx.GetGroupID()
		return groupMap[groupID]
	}
}

// OnlySuperUsers 只处理超级用户的消息
func OnlySuperUsers() Filter {
	return func(ctx *Context) bool {
		userID := ctx.GetUserID()
		for _, su := range ctx.Bot.Config.SuperUsers {
			if userID == su {
				return true
			}
		}
		return false
	}
}

// OnlyToMe 只处理 @ 机器人的消息
func OnlyToMe() Filter {
	return func(ctx *Context) bool {
		msg := ctx.GetMessage()
		if msg == nil {
			return false
		}

		// 检查是否有 @ 消息段
		for _, seg := range *msg {
			if seg.Type == "at" {
				qq := cast.ToInt64(seg.Data["qq"])
				if qq == ctx.Event.GetSelfID() {
					return true
				}
			}
		}

		// 检查是否包含昵称
		text := msg.GetPlainText()
		for _, nickname := range ctx.Bot.Config.Nickname {
			if len(text) > 0 && len(nickname) > 0 {
				matched, _ := regexp.MatchString(nickname, text)
				if matched {
					return true
				}
			}
		}

		return false
	}
}

// HasPermission 权限过滤器（群主或管理员）
func HasPermission() Filter {
	return func(ctx *Context) bool {
		if evt, ok := ctx.Event.(*event.GroupMessageEvent); ok {
			return evt.IsOwner() || evt.IsAdmin()
		}
		return false
	}
}

// And 组合过滤器（AND）
func And(filters ...Filter) Filter {
	return func(ctx *Context) bool {
		for _, f := range filters {
			if !f(ctx) {
				return false
			}
		}
		return true
	}
}

// Or 组合过滤器（OR）
func Or(filters ...Filter) Filter {
	return func(ctx *Context) bool {
		for _, f := range filters {
			if f(ctx) {
				return true
			}
		}
		return false
	}
}

// Not 反转过滤器
func Not(filter Filter) Filter {
	return func(ctx *Context) bool {
		return !filter(ctx)
	}
}

// ===== DFA 敏感词过滤器 =====

// DFANode DFA 树节点
type DFANode struct {
	children map[rune]*DFANode
	isEnd    bool
}

// SensitiveFilter DFA 敏感词过滤器
type SensitiveFilter struct {
	trie     *DFANode
	mu       sync.RWMutex
	filePath string
	watcher  *fsnotify.Watcher
}

// NewSensitiveFilter 创建敏感词过滤器
func NewSensitiveFilter(filePath string) (*SensitiveFilter, error) {
	sf := &SensitiveFilter{
		trie:     &DFANode{children: make(map[rune]*DFANode)},
		filePath: filePath,
	}

	if err := sf.Reload(); err != nil {
		return nil, err
	}

	return sf, nil
}

// Reload 重新加载敏感词库
func (sf *SensitiveFilter) Reload() error {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	// 读取文件
	data, err := readLines(sf.filePath)
	if err != nil {
		return err
	}

	// 重建 DFA 树
	sf.trie = &DFANode{children: make(map[rune]*DFANode)}
	for _, word := range data {
		if word == "" {
			continue
		}
		sf.addWord(word)
	}

	logger.Info("敏感词库已重新加载", "count", len(data))
	return nil
}

// addWord 添加敏感词到 DFA 树
func (sf *SensitiveFilter) addWord(word string) {
	node := sf.trie
	for _, char := range []rune(word) {
		if _, ok := node.children[char]; !ok {
			node.children[char] = &DFANode{children: make(map[rune]*DFANode)}
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// Contains 检查消息是否包含敏感词
func (sf *SensitiveFilter) Contains(text string) bool {
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		node := sf.trie
		j := i

		for j < len(runes) {
			if child, ok := node.children[runes[j]]; ok {
				node = child
				j++
				if node.isEnd {
					return true
				}
			} else {
				break
			}
		}
	}

	return false
}

// Filter 过滤敏感词，返回替换后的文本
func (sf *SensitiveFilter) Filter(text string, replace rune) string {
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	runes := []rune(text)
	result := make([]rune, len(runes))
	copy(result, runes)

	for i := 0; i < len(runes); i++ {
		node := sf.trie
		j := i
		start := i

		for j < len(runes) {
			if child, ok := node.children[runes[j]]; ok {
				node = child
				j++
				if node.isEnd {
					// 替换敏感词
					for k := start; k < j; k++ {
						result[k] = replace
					}
				}
			} else {
				break
			}
		}
	}

	return string(result)
}

// FindAll 获取消息中的所有敏感词
func (sf *SensitiveFilter) FindAll(text string) []string {
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	var words []string
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		node := sf.trie
		j := i

		for j < len(runes) {
			if child, ok := node.children[runes[j]]; ok {
				node = child
				j++
				if node.isEnd {
					words = append(words, string(runes[i:j]))
				}
			} else {
				break
			}
		}
	}

	return words
}

// AsFilter 作为 Filter 使用（阻止包含敏感词的消息）
func (sf *SensitiveFilter) AsFilter() Filter {
	return func(ctx *Context) bool {
		msg := ctx.GetMessage()
		if msg == nil {
			return true
		}
		text := msg.GetPlainText()
		return !sf.Contains(text)
	}
}

// AsReplacer 作为 Filter 使用（自动替换敏感词）
func (sf *SensitiveFilter) AsReplacer(replace rune) Filter {
	return func(ctx *Context) bool {
		msg := ctx.GetMessage()
		if msg == nil {
			return true
		}

		text := msg.GetPlainText()
		if sf.Contains(text) {
			filtered := sf.Filter(text, replace)
			// 替换消息内容
			*msg = message.Message{message.Text(filtered)}
		}

		return true
	}
}

// EnableHotReload 启用热重载
func (sf *SensitiveFilter) EnableHotReload() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	sf.watcher = watcher

	if err := watcher.Add(sf.filePath); err != nil {
		return err
	}

	go sf.watchFile()
	return nil
}

// watchFile 监听文件变化
func (sf *SensitiveFilter) watchFile() {
	for {
		select {
		case event, ok := <-sf.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				logger.Info("敏感词文件已更新", "path", sf.filePath)
				if err := sf.Reload(); err != nil {
					logger.Error("重新加载敏感词失败", "error", err)
				}
			}
		case err, ok := <-sf.watcher.Errors:
			if !ok {
				return
			}
			logger.Error("文件监控错误", "error", err)
		}
	}
}

// StopHotReload 停止热重载
func (sf *SensitiveFilter) StopHotReload() {
	if sf.watcher != nil {
		sf.watcher.Close()
	}
}

// readLines 读取文件的所有行
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}
