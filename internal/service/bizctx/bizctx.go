// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package bizctx

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Service provides business context related logic.
type Service struct{}

const ContextKey = "ContextKey"

// Context defines the business context object injected into request context.
type Context struct {
	Session *ghttp.Session
	User    *User
}

// User defines the business user object injected into request context.
type User struct {
	Id    uint64
	Name  string
	Phone string
	Role  string // 权限角色：user-普通用户 admin-管理员
}

func New() *Service {
	return &Service{}
}

// Init initializes and injects custom business context into request context.
func (s *Service) Init(r *ghttp.Request, customCtx *Context) {
	r.SetCtxVar(ContextKey, customCtx)
}

// Get retrieves and returns the business context from context.
func (s *Service) Get(ctx context.Context) *Context {
	value := ctx.Value(ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*Context); ok {
		return localCtx
	}
	return nil
}

// SetUser injects business user object into context.
func (s *Service) SetUser(ctx context.Context, ctxUser *User) {
	s.Get(ctx).User = ctxUser
}
