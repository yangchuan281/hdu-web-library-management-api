// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package borrow

import (
	"context"
	"library-management-api/api/borrow/v1"
	borrowsvc "library-management-api/internal/service/borrow"
)

func (c *ControllerV1) Borrow(ctx context.Context, req *v1.BorrowReq) (res *v1.BorrowRes, err error) {
	// 获取当前用户ID
	userId, err := c.borrowSvc.GetCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	out, err := c.borrowSvc.Borrow(ctx, borrowsvc.BorrowInput{
		UserId: userId,
		BookId: req.BookId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.BorrowRes{
		Id:       out.Id,
		BookId:   req.BookId,
		BookName: out.BookName,
		BorrowAt: out.BorrowAt,
	}, nil
}

func (c *ControllerV1) Return(ctx context.Context, req *v1.ReturnReq) (res *v1.ReturnRes, err error) {
	returnAt, err := c.borrowSvc.Return(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.ReturnRes{
		Id:       req.Id,
		ReturnAt: returnAt,
	}, nil
}

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	// 获取当前用户ID（普通用户只看自己的，管理员可以看全部）
	userId, err := c.borrowSvc.GetCurrentUserId(ctx)
	if err != nil {
		return nil, err
	}

	total, records, err := c.borrowSvc.List(ctx, borrowsvc.ListInput{
		UserId: userId,
		Page:   req.Page,
		Size:   req.Size,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	items := make([]v1.BorrowItem, 0, len(records))
	for _, r := range records {
		items = append(items, v1.BorrowItem{
			Id:        r.Id,
			UserId:    r.UserId,
			UserName:  r.UserName,
			BookId:    r.BookId,
			BookName:  r.BookName,
			BorrowAt:  r.BorrowAt,
			ReturnAt:  r.ReturnAt,
		})
	}

	return &v1.ListRes{
		List:  items,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}
