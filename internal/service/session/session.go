// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package session

import (
	"context"
	"library-management-api/internal/model/entity"
	"library-management-api/internal/service/bizctx"
)

const UserSessionKey = "UserSessionKey"

// Service provides session management logic.
type Service struct {
	bizCtxSvc *bizctx.Service
}

func New() *Service {
	return &Service{
		bizCtxSvc: bizctx.New(),
	}
}

// SetUser stores user info into session.
func (s *Service) SetUser(ctx context.Context, user *entity.Users) error {
	return s.bizCtxSvc.Get(ctx).Session.Set(UserSessionKey, user)
}

// GetUser retrieves and returns user from session.
func (s *Service) GetUser(ctx context.Context) *entity.Users {
	customCtx := s.bizCtxSvc.Get(ctx)
	if customCtx != nil {
		if v := customCtx.Session.MustGet(UserSessionKey); !v.IsNil() {
			var user *entity.Users
			_ = v.Struct(&user)
			return user
		}
	}
	return nil
}

// RemoveUser removes user from session.
func (s *Service) RemoveUser(ctx context.Context) error {
	customCtx := s.bizCtxSvc.Get(ctx)
	if customCtx != nil {
		return customCtx.Session.Remove(UserSessionKey)
	}
	return nil
}
