// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package user

import (
	"context"
	"library-management-api/api/user/v1"
	usersvc "library-management-api/internal/service/user"
)

// Login 统一登录（支持邮箱或手机号）
func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	user, token, err := c.userSvc.Login(ctx, usersvc.LoginInput{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &v1.LoginRes{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Token: token,
	}, nil
}
