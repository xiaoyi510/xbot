package middleware

// HandlerFunc 处理函数类型（占位符，实际在 context 中定义）
type HandlerFunc func(ctx interface{})

// Middleware 中间件类型
type Middleware func(next HandlerFunc) HandlerFunc

// Chain 中间件链
type Chain struct {
	middlewares []Middleware
}

// NewChain 创建中间件链
func NewChain(middlewares ...Middleware) *Chain {
	return &Chain{
		middlewares: middlewares,
	}
}

// Use 添加中间件
func (c *Chain) Use(m ...Middleware) *Chain {
	c.middlewares = append(c.middlewares, m...)
	return c
}

// Then 应用中间件链到处理函数
func (c *Chain) Then(handler HandlerFunc) HandlerFunc {
	// 从后向前应用中间件
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	return handler
}

// Len 返回中间件数量
func (c *Chain) Len() int {
	return len(c.middlewares)
}
