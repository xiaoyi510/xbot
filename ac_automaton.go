package xbot

import "sync"

// ACNode AC自动机节点
type ACNode struct {
	children map[rune]*ACNode // 子节点映射
	fail     *ACNode          // 失败指针
	output   bool             // 是否为某个模式串的结尾
}

// NewACNode 创建新的AC节点
func NewACNode() *ACNode {
	return &ACNode{
		children: make(map[rune]*ACNode),
		fail:     nil,
		output:   false,
	}
}

// ACMachine AC自动机
type ACMachine struct {
	root        *ACNode
	lastVersion int64
	mu          sync.RWMutex
}

// NewACMachine 创建新的AC自动机
func NewACMachine() *ACMachine {
	return &ACMachine{
		root:        NewACNode(),
		lastVersion: -1,
	}
}

// Build 构建AC自动机
// 参数：patterns - 模式串列表
func (ac *ACMachine) Build(patterns []string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	// 重新初始化root
	ac.root = NewACNode()

	// 第一步：构建Trie树
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}
		node := ac.root
		for _, ch := range []rune(pattern) {
			if _, exists := node.children[ch]; !exists {
				node.children[ch] = NewACNode()
			}
			node = node.children[ch]
		}
		node.output = true
	}

	// 第二步：构建失败指针（使用BFS）
	queue := make([]*ACNode, 0)

	// 初始化第一层节点的失败指针
	for _, child := range ac.root.children {
		child.fail = ac.root
		queue = append(queue, child)
	}

	// BFS构建其他层的失败指针
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for ch, child := range current.children {
			queue = append(queue, child)

			// 寻找失败指针
			failNode := current.fail
			for failNode != nil {
				if next, exists := failNode.children[ch]; exists {
					child.fail = next
					break
				}
				if failNode == ac.root {
					child.fail = ac.root
					break
				}
				failNode = failNode.fail
			}

			// 如果没有找到，指向root
			if child.fail == nil {
				child.fail = ac.root
			}

			// 继承失败节点的输出状态
			if child.fail.output {
				child.output = true
			}
		}
	}
}

// Match 匹配文本，返回是否找到任意一个模式串
// 参数：text - 待匹配文本
// 返回：是否匹配到模式串
func (ac *ACMachine) Match(text string) bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	if ac.root == nil {
		return false
	}

	node := ac.root
	for _, ch := range []rune(text) {
		// 沿着失败指针查找，直到找到匹配或回到root
		for node != ac.root && node.children[ch] == nil {
			node = node.fail
		}

		// 尝试转移
		if next, exists := node.children[ch]; exists {
			node = next
		}

		// 检查当前节点是否是某个模式的结尾
		if node.output {
			return true
		}
	}

	return false
}

// MatchAll 匹配文本，返回所有匹配到的模式串位置
// 参数：text - 待匹配文本
// 返回：匹配到的位置列表 (结束位置)
func (ac *ACMachine) MatchAll(text string) []int {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	if ac.root == nil {
		return []int{}
	}

	matches := make([]int, 0)
	node := ac.root

	runes := []rune(text)
	for i, ch := range runes {
		// 沿着失败指针查找
		for node != ac.root && node.children[ch] == nil {
			node = node.fail
		}

		// 尝试转移
		if next, exists := node.children[ch]; exists {
			node = next
		}

		// 检查当前节点及其失败链上的所有输出节点
		temp := node
		for temp != ac.root {
			if temp.output {
				matches = append(matches, i)
				break // 只记录一次位置
			}
			temp = temp.fail
		}
	}

	return matches
}

// GetVersion 获取版本号
func (ac *ACMachine) GetVersion() int64 {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.lastVersion
}

// SetVersion 设置版本号
func (ac *ACMachine) SetVersion(version int64) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.lastVersion = version
}
