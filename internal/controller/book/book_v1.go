// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package book

import (
	"context"
	"library-management-api/api/book/v1"
	booksvc "library-management-api/internal/service/book"
)

func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.bookSvc.Create(ctx, booksvc.CreateInput{
		Title:       req.Title,
		Author:      req.Author,
		Isbn:        req.Isbn,
		PublishDate: req.PublishDate,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{
		Id:    id,
		Title: req.Title,
		Isbn:  req.Isbn,
	}, nil
}

func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := c.bookSvc.List(ctx, booksvc.ListInput{
		Page:   req.Page,
		Size:   req.Size,
		Title:  req.Title,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	items := make([]v1.BookItem, 0, len(out.List))
	for _, book := range out.List {
		item := v1.BookItem{
			Id:     book.Id,
			Title:  book.Title,
			Author: book.Author,
			Isbn:   book.Isbn,
			Status: book.Status,
		}
		if book.PublishDate != nil {
			item.PublishDate = book.PublishDate.String()
		}
		if book.CreatedAt != nil {
			item.CreatedAt = book.CreatedAt.String()
		}
		items = append(items, item)
	}

	return &v1.ListRes{
		List:  items,
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	book, err := c.bookSvc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	item := &v1.BookItem{
		Id:     book.Id,
		Title:  book.Title,
		Author: book.Author,
		Isbn:   book.Isbn,
		Status: book.Status,
	}
	if book.PublishDate != nil {
		item.PublishDate = book.PublishDate.String()
	}
	if book.CreatedAt != nil {
		item.CreatedAt = book.CreatedAt.String()
	}

	return &v1.GetRes{BookItem: item}, nil
}

func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = c.bookSvc.Update(ctx, booksvc.UpdateInput{
		Id:          req.Id,
		Title:       req.Title,
		Author:      req.Author,
		Isbn:        req.Isbn,
		PublishDate: req.PublishDate,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateRes{
		Id:    req.Id,
		Title: req.Title,
	}, nil
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	err = c.bookSvc.Delete(ctx, req.Id)
	return
}
