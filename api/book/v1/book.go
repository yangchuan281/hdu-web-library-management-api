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

// CreateReq 创建图书请求
type CreateReq struct {
	g.Meta       `path:"/books" method:"post" tags:"图书服务" summary:"新增图书"`
	Title        string `v:"required" dc:"图书名称"`
	Author       string `v:"required" dc:"作者"`
	Isbn         string `v:"required" dc:"ISBN号"`
	PublishDate  string `dc:"出版日期（格式：2006-01-02）"`
}

// CreateRes 创建图书响应
type CreateRes struct {
	Id     uint64 `json:"id" dc:"图书ID"`
	Title  string `json:"title" dc:"图书名称"`
	Isbn   string `json:"isbn" dc:"ISBN号"`
}

// ListReq 图书列表请求
type ListReq struct {
	g.Meta `path:"/books" method:"get" tags:"图书服务" summary:"获取图书列表"`
	Page   int `d:"1" dc:"页码"`
	Size   int `d:"10" dc:"每页条数"`
	Title  string `dc:"图书名称（模糊搜索）"`
	Status int    `d:"-1" dc:"状态筛选（-1=全部，1=可借阅，0=已借出）"`
}

// ListRes 图书列表响应
type ListRes struct {
	List     []BookItem `json:"list" dc:"图书列表"`
	Total    int        `json:"total" dc:"总条数"`
	Page     int        `json:"page" dc:"当前页码"`
	Size     int        `json:"size" dc:"每页条数"`
}

type BookItem struct {
	Id          uint64 `json:"id" dc:"图书ID"`
	Title       string `json:"title" dc:"图书名称"`
	Author      string `json:"author" dc:"作者"`
	Isbn        string `json:"isbn" dc:"ISBN号"`
	Status      int8   `json:"status" dc:"状态：1-可借阅，0-已借出，2-下架"`
	PublishDate string `json:"publish_date" dc:"出版日期"`
	CreatedAt   string `json:"created_at" dc:"创建时间"`
}

// GetReq 获取图书详情请求
type GetReq struct {
	g.Meta `path:"/books/{id}" method:"get" tags:"图书服务" summary:"获取图书详情"`
	Id     uint64 `v:"required" dc:"图书ID"`
}

// GetRes 获取图书详情响应
type GetRes struct {
	*BookItem
}

// UpdateReq 更新图书信息请求
type UpdateReq struct {
	g.Meta       `path:"/books/{id}" method:"put" tags:"图书服务" summary:"更新图书信息"`
	Id           uint64 `v:"required" dc:"图书ID"`
	Title        string `dc:"图书名称"`
	Author       string `dc:"作者"`
	Isbn         string `dc:"ISBN号"`
	PublishDate  string `dc:"出版日期"`
	Status       int8   `d:"-1" dc:"状态"`
}

// UpdateRes 更新图书信息响应
type UpdateRes struct {
	Id     uint64 `json:"id" dc:"图书ID"`
	Title  string `json:"title" dc:"图书名称"`
}

// DeleteReq 删除图书请求
type DeleteReq struct {
	g.Meta `path:"/books/{id}" method:"delete" tags:"图书服务" summary:"删除图书"`
	Id     uint64 `v:"required" dc:"图书ID"`
}

// DeleteRes 删除图书响应
type DeleteRes struct {
}
