// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"
	"library-management-api/api/user/v1"
)

// IUserV1 defines the interface for user API version 1.
type IUserV1 interface {
	// POST /api/users - 手机号注册
	SignUp(ctx context.Context, req *v1.SignUpReq) (res *v1.SignUpRes, err error)
	// POST /api/auth/login - 登录（邮箱或手机号）
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	// DELETE /api/auth/logout - 登出
	SignOut(ctx context.Context, req *v1.SignOutReq) (res *v1.SignOutRes, err error)
	// GET /api/users/me - 个人信息
	Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error)
	// POST /api/auth/verification-codes - 发送验证码
	SendVerificationCode(ctx context.Context, req *v1.SendVerificationCodeReq) (res *v1.SendVerificationCodeRes, err error)
	// POST /api/auth/register - 邮箱注册
	Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error)
	// PUT /api/auth/password - 重置密码
	ResetPassword(ctx context.Context, req *v1.ResetPasswordReq) (res *v1.ResetPasswordRes, err error)
}
