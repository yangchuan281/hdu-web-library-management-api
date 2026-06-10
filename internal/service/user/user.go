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
	"errors"
	"crypto/md5"
	"fmt"
	"time"
	"library-management-api/internal/dao"
	"library-management-api/internal/model/do"
	"library-management-api/internal/model/entity"
	"library-management-api/internal/service/bizctx"
	"library-management-api/internal/service/email"
	"library-management-api/internal/service/jwt"
	"library-management-api/internal/service/session"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Service provides user-related business logic.
type Service struct {
	bizCtxSvc  *bizctx.Service
	sessionSvc *session.Service
	emailSvc   *email.Service   // 实验三新增：邮件发送服务
	jwtSvc     *jwt.Service     // 实验三新增：JWT认证服务
}

func New() *Service {
	return &Service{
		bizCtxSvc:  bizctx.New(),
		sessionSvc: session.New(),
		emailSvc:   email.New(),
		jwtSvc:     jwt.New(),
	}
}

// 使用MD5对密码进行加密
func encryptPassword(password string) string {
	data := []byte(password)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// CreateInput defines the input for creating a new user account.
type CreateInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

// Create registers a new user account.
func (s *Service) Create(ctx context.Context, in CreateInput) (userId uint64, err error) {
	// 检查手机号是否已注册
	count, err := dao.Users.Ctx(ctx).Where(do.Users{
		Phone: in.Phone,
	}).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("该手机号已被注册")
	}

	// 检查邮箱是否已注册
	if in.Email != "" {
		count, err = dao.Users.Ctx(ctx).Where(do.Users{
			Email: in.Email,
		}).Count()
		if err != nil {
			return 0, err
		}
		if count > 0 {
			return 0, errors.New("该邮箱已被注册")
		}
	}

	// 事务创建用户
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model("users").Data(g.Map{
			"name":     in.Name,
			"email":    in.Email,
			"phone":    in.Phone,
			"password": encryptPassword(in.Password),
			"status":   1,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		userId = uint64(id)
		return nil
	})
	return userId, err
}

// SignInInput defines the input for signing in.
type SignInInput struct {
	Phone    string
	Password string
}

// SignIn authenticates user and creates session.
func (s *Service) SignIn(ctx context.Context, in SignInInput) (user *entity.Users, err error) {
	err = dao.Users.Ctx(ctx).Where(do.Users{
		Phone:    in.Phone,
		Password: encryptPassword(in.Password),
		Status:   1,
	}).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("手机号或密码错误")
	}

	// 设置会话
	if err = s.sessionSvc.SetUser(ctx, user); err != nil {
		return nil, err
	}
	s.bizCtxSvc.SetUser(ctx, &bizctx.User{
		Id:    user.Id,
		Name:  user.Name,
		Phone: user.Phone,
	})
	return user, nil
}

// IsSignedIn checks whether current user is signed in.
func (s *Service) IsSignedIn(ctx context.Context) bool {
	if v := s.bizCtxSvc.Get(ctx); v != nil && v.User != nil {
		return true
	}
	return false
}

// SignOut removes session for current signed-in user.
func (s *Service) SignOut(ctx context.Context) error {
	return s.sessionSvc.RemoveUser(ctx)
}

// GetProfile retrieves current user info.
func (s *Service) GetProfile(ctx context.Context) (*entity.Users, error) {
	customCtx := s.bizCtxSvc.Get(ctx)
	if customCtx == nil || customCtx.User == nil {
		return nil, errors.New("用户未登录")
	}

	var user *entity.Users
	err := dao.Users.Ctx(ctx).Where("id", customCtx.User.Id).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	// 不返回密码
	user.Password = ""
	return user, nil
}

// GetUserById retrieves user by ID.
func (s *Service) GetUserById(ctx context.Context, id uint64) (*entity.Users, error) {
	var user *entity.Users
	err := dao.Users.Ctx(ctx).Where("id", id).Scan(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetCurrentUserId returns current logged-in user's ID.
func (s *Service) GetCurrentUserId(ctx context.Context) (uint64, error) {
	customCtx := s.bizCtxSvc.Get(ctx)
	if customCtx == nil || customCtx.User == nil {
		return 0, errors.New("用户未登录")
	}
	return customCtx.User.Id, nil
}

// ==================== 实验三新增：RESTful 邮箱注册/登录/重置密码 ====================

// LoginInput 统一登录参数（支持邮箱或手机号）
type LoginInput struct {
	Email    string
	Phone    string
	Password string
}

// Login 统一登录（支持邮箱或手机号，返回JWT）
func (s *Service) Login(ctx context.Context, in LoginInput) (*entity.Users, string, error) {
	var user *entity.Users
	var err error

	if in.Email != "" {
		// 邮箱登录
		err = dao.Users.Ctx(ctx).Where(do.Users{
			Email:    in.Email,
			Password: encryptPassword(in.Password),
			Status:   1,
		}).Scan(&user)
		if err != nil {
			return nil, "", err
		}
		if user == nil {
			return nil, "", errors.New("邮箱或密码错误")
		}
	} else if in.Phone != "" {
		// 手机号登录
		err = dao.Users.Ctx(ctx).Where(do.Users{
			Phone:    in.Phone,
			Password: encryptPassword(in.Password),
			Status:   1,
		}).Scan(&user)
		if err != nil {
			return nil, "", err
		}
		if user == nil {
			return nil, "", errors.New("手机号或密码错误")
		}
	} else {
		return nil, "", errors.New("请输入邮箱或手机号")
	}

	// 生成JWT Token（统一使用JWT认证）
	token, err := s.jwtSvc.GenerateToken(user.Id, user.Phone, user.Email, user.Role)
	if err != nil {
		return nil, "", errors.New("Token生成失败")
	}

	return user, token, nil
}

// SendRegisterCode 发送注册验证码到邮箱
func (s *Service) SendRegisterCode(ctx context.Context, emailAddr string) (string, error) {
	// 检查邮箱是否已被注册
	count, err := dao.Users.Ctx(ctx).Where(do.Users{Email: emailAddr}).Count()
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New("该邮箱已被注册")
	}

	// 生成4位验证码
	code := s.emailSvc.GenerateCode()

	// 保存验证码到数据库（验证码有效期5分钟）
	_, err = g.DB().Model("verification_codes").Insert(g.Map{
		"email":      emailAddr,
		"code":       code,
		"type":       "register",
		"used":       0,
		"expires_at": time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	// 发送验证码邮件
	err = s.emailSvc.SendVerificationCode(emailAddr, code, "register")
	if err != nil {
		g.Log().Warning(ctx, "邮件发送失败: %v", err)
	}

	return code, nil
}

// EmailSignUp 邮箱注册（含验证码验证）
func (s *Service) EmailSignUp(ctx context.Context, emailAddr, code, password, name string) (*entity.Users, string, error) {
	// 验证验证码有效性
	record, err := g.DB().Model("verification_codes").
		Where("email = ? AND code = ? AND type = ? AND used = 0", emailAddr, code, "register").
		OrderDesc("id").
		One()
	if err != nil {
		return nil, "", errors.New("验证码查询失败")
	}
	if record == nil {
		return nil, "", errors.New("验证码无效或已使用")
	}

	// 检查过期时间
	if time.Now().After(record["expires_at"].Time()) {
		return nil, "", errors.New("验证码已过期")
	}
	vcId := record["id"].Int64()

	// 标记验证码为已使用
	_, err = g.DB().Model("verification_codes").
		Where("id", vcId).
		Update(g.Map{"used": 1})
	if err != nil {
		return nil, "", err
	}

	randPhone := ""

	// 创建用户
	var userId uint64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model("users").Data(g.Map{
			"name":     name,
			"email":    emailAddr,
			"phone":    randPhone,
			"password": encryptPassword(password),
			"status":   1,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		userId = uint64(id)
		return nil
	})
	if err != nil {
		return nil, "", errors.New("用户创建失败: " + err.Error())
	}

	// 生成JWT Token
	token, err := s.jwtSvc.GenerateToken(userId, randPhone, emailAddr, "user")
	if err != nil {
		return nil, "", errors.New("Token生成失败")
	}

	// 返回用户信息
	user := &entity.Users{
		Id:    userId,
		Name:  name,
		Email: emailAddr,
		Phone: randPhone,
	}

	return user, token, nil
}

// SendResetCode 发送密码重置验证码到邮箱
func (s *Service) SendResetCode(ctx context.Context, emailAddr string) (string, error) {
	// 检查邮箱是否存在
	count, err := dao.Users.Ctx(ctx).Where(do.Users{Email: emailAddr}).Count()
	if err != nil {
		return "", err
	}
	if count == 0 {
		return "", errors.New("该邮箱未注册")
	}

	// 生成4位验证码
	code := s.emailSvc.GenerateCode()

	// 保存验证码到数据库
	_, err = g.DB().Model("verification_codes").Insert(g.Map{
		"email":      emailAddr,
		"code":       code,
		"type":       "reset",
		"used":       0,
		"expires_at": time.Now().Add(5 * time.Minute),
	})
	if err != nil {
		return "", err
	}

	// 发送验证码邮件
	err = s.emailSvc.SendVerificationCode(emailAddr, code, "reset")
	if err != nil {
		g.Log().Warning(ctx, "邮件发送失败: %v", err)
	}

	return code, nil
}

// ResetPassword 重置密码（验证验证码后更新密码）
func (s *Service) ResetPassword(ctx context.Context, emailAddr, code, newPassword string) error {
	// 验证重置码
	record, err := g.DB().Model("verification_codes").
		Where("email = ? AND code = ? AND type = ? AND used = 0", emailAddr, code, "reset").
		OrderDesc("id").
		One()
	if err != nil {
		return errors.New("重置码查询失败")
	}
	if record == nil {
		return errors.New("重置码无效或已使用")
	}

	// 检查过期时间
	if time.Now().After(record["expires_at"].Time()) {
		return errors.New("重置码已过期")
	}
	vcId := record["id"].Int64()

	// 标记验证码为已使用
	_, err = g.DB().Model("verification_codes").
		Where("id", vcId).
		Update(g.Map{"used": 1})
	if err != nil {
		return err
	}

	// 更新密码
	_, err = dao.Users.Ctx(ctx).
		Where("email", emailAddr).
		Update(g.Map{"password": encryptPassword(newPassword)})
	if err != nil {
		return errors.New("密码重置失败")
	}

	return nil
}
