// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package cmd

import (
	"context"
	"library-management-api/internal/controller/user"
	"library-management-api/internal/controller/book"
	"library-management-api/internal/controller/borrow"
	"library-management-api/internal/service/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = &gcmd.Command{
		Name:  "main",
		Brief: "start http server of library management system",
		Func:  mainFunc,
	}
)

func mainFunc(ctx context.Context, parser *gcmd.Parser) (err error) {
	s := g.Server()

	middlewareSvc := middleware.New()

	// 全局中间件
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Use(middlewareSvc.Ctx)
	s.Use(ghttp.MiddlewareCORS)

	userCtrl := user.NewV1()
	bookCtrl := book.NewV1()
	borrowCtrl := borrow.NewV1()

	// ========== RESTful API 路由设计 ==========
	// 原则：名词复数表示资源，HTTP 方法表示操作
	// POST   = 创建资源
	// GET    = 查询资源
	// PUT    = 更新资源
	// DELETE = 删除资源
	// ==========================================
	s.Group("/api", func(group *ghttp.RouterGroup) {

		// ─── 认证服务（无需认证）───
		// POST   /api/auth/verification-codes  发送验证码（register/reset）
		// POST   /api/auth/register            邮箱注册（含验证码）
		// POST   /api/auth/login               登录（邮箱或手机号）
		// DELETE /api/auth/logout              登出
		// PUT    /api/auth/password            重置密码
		group.Group("/auth", func(group *ghttp.RouterGroup) {
			group.POST("/verification-codes", userCtrl.SendVerificationCode)
			group.POST("/register", userCtrl.Register)
			group.POST("/login", userCtrl.Login)
			group.DELETE("/logout", userCtrl.SignOut)
			group.PUT("/password", userCtrl.ResetPassword)
		})

		// ─── 用户服务（无需认证）───
		// POST   /api/users                    手机号注册
		// GET    /api/users/me                 个人信息
		group.Group("/users", func(group *ghttp.RouterGroup) {
			group.POST("/", userCtrl.SignUp)
			group.GET("/me", userCtrl.Profile)
		})

		// ─── 图书服务（公开：查询）───
		// GET    /api/books                    图书列表
		// GET    /api/books/{id}               图书详情
		group.GET("/books", bookCtrl.List)
		group.GET("/books/{id}", bookCtrl.Get)

		// ─── 需认证的接口 ───
		group.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(middlewareSvc.Auth)

			// ─── 图书管理（需管理员权限）───
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middlewareSvc.Admin)

				// POST   /api/books                 创建图书
				// PUT    /api/books/{id}            更新图书
				// DELETE /api/books/{id}            删除图书
				group.POST("/books", bookCtrl.Create)
				group.PUT("/books/{id}", bookCtrl.Update)
				group.DELETE("/books/{id}", bookCtrl.Delete)
			})

			// ─── 借阅管理（登录即可）───
			// POST   /api/borrows               借阅图书
			// GET    /api/borrows               借阅列表
			// PUT    /api/borrows/{id}/return    归还图书
			group.POST("/borrows", borrowCtrl.Borrow)
			group.GET("/borrows", borrowCtrl.List)
			group.PUT("/borrows/{id}/return", borrowCtrl.Return)
		})
	})

	s.Run()
	return nil
}
