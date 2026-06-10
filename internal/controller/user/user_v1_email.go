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
)

// SendVerificationCode 发送邮箱验证码（POST /api/auth/verification-codes）
// 支持类型：register-注册, reset-重置密码
func (c *ControllerV1) SendVerificationCode(ctx context.Context, req *v1.SendVerificationCodeReq) (res *v1.SendVerificationCodeRes, err error) {
	switch req.Type {
	case "register":
		_, err = c.userSvc.SendRegisterCode(ctx, req.Email)
		if err != nil {
			return nil, err
		}
	case "reset":
		_, err = c.userSvc.SendResetCode(ctx, req.Email)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	return &v1.SendVerificationCodeRes{
		Message: "验证码已发送到邮箱（若未配置SMTP，请查看服务器日志）",
	}, nil
}

// Register 邮箱注册（POST /api/auth/register）
func (c *ControllerV1) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	user, token, err := c.userSvc.EmailSignUp(ctx, req.Email, req.Code, req.Password, req.Name)
	if err != nil {
		return nil, err
	}

	return &v1.RegisterRes{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}

// ResetPassword 重置密码（PUT /api/auth/password）
func (c *ControllerV1) ResetPassword(ctx context.Context, req *v1.ResetPasswordReq) (res *v1.ResetPasswordRes, err error) {
	err = c.userSvc.ResetPassword(ctx, req.Email, req.Code, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &v1.ResetPasswordRes{
		Message: "密码重置成功",
	}, nil
}
