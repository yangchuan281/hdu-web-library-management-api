// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package v1

import "github.com/gogf/gf/v2/frame/g"

// BorrowReq 借阅图书请求
type BorrowReq struct {
	g.Meta `path:"/borrows" method:"post" tags:"借阅服务" summary:"借阅图书"`
	BookId uint64 `v:"required" dc:"图书ID"`
}

// BorrowRes 借阅图书响应
type BorrowRes struct {
	Id       uint64 `json:"id" dc:"借阅记录ID"`
	BookId   uint64 `json:"book_id" dc:"图书ID"`
	BookName string `json:"book_name" dc:"图书名称"`
	BorrowAt string `json:"borrow_at" dc:"借阅时间"`
}

// ReturnReq 归还图书请求
type ReturnReq struct {
	g.Meta `path:"/borrows/{id}/return" method:"put" tags:"借阅服务" summary:"归还图书"`
	Id     uint64 `v:"required" dc:"借阅记录ID"`
}

// ReturnRes 归还图书响应
type ReturnRes struct {
	Id       uint64 `json:"id" dc:"借阅记录ID"`
	ReturnAt string `json:"return_at" dc:"归还时间"`
}

// ListReq 借阅记录列表请求
type ListReq struct {
	g.Meta `path:"/borrows" method:"get" tags:"借阅服务" summary:"获取借阅记录列表"`
	Page   int `d:"1" dc:"页码"`
	Size   int `d:"10" dc:"每页条数"`
	Status int `d:"-1" dc:"状态筛选（-1=全部，0=借阅中，1=已归还）"`
}

// ListRes 借阅记录列表响应
type ListRes struct {
	List  []BorrowItem `json:"list" dc:"借阅记录列表"`
	Total int          `json:"total" dc:"总条数"`
	Page  int          `json:"page" dc:"当前页码"`
	Size  int          `json:"size" dc:"每页条数"`
}

type BorrowItem struct {
	Id        uint64 `json:"id" dc:"借阅记录ID"`
	UserId    uint64 `json:"user_id" dc:"用户ID"`
	UserName  string `json:"user_name" dc:"用户姓名"`
	BookId    uint64 `json:"book_id" dc:"图书ID"`
	BookName  string `json:"book_name" dc:"图书名称"`
	BorrowAt  string `json:"borrow_at" dc:"借阅时间"`
	ReturnAt  string `json:"return_at" dc:"归还时间"`
}
