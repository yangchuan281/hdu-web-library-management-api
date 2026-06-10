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
	"errors"
	"library-management-api/internal/dao"
	"library-management-api/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

// Service provides book-related business logic.
type Service struct{}

func New() *Service {
	return &Service{}
}

// CreateInput defines input for creating a new book.
type CreateInput struct {
	Title       string
	Author      string
	Isbn        string
	PublishDate string
}

// Create adds a new book.
func (s *Service) Create(ctx context.Context, in CreateInput) (bookId uint64, err error) {
	// 检查ISBN是否已存在
	count, err := dao.Books.Ctx(ctx).Where("isbn", in.Isbn).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("该ISBN号已存在")
	}

	data := g.Map{
		"title":  in.Title,
		"author": in.Author,
		"isbn":   in.Isbn,
		"status": 1,
	}
	if in.PublishDate != "" {
		data["publish_date"] = in.PublishDate
	}

	id, err := dao.Books.Ctx(ctx).Data(data).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

// Get retrieves a book by ID.
func (s *Service) Get(ctx context.Context, id uint64) (*entity.Books, error) {
	var book *entity.Books
	err := dao.Books.Ctx(ctx).Where("id", id).Scan(&book)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("图书不存在")
	}
	return book, nil
}

// ListInput defines input for listing books.
type ListInput struct {
	Page   int
	Size   int
	Title  string
	Status int
}

// ListOutput defines output for listing books.
type ListOutput struct {
	List  []*entity.Books
	Total int
	Page  int
	Size  int
}

// List retrieves books with pagination and filters.
func (s *Service) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	m := dao.Books.Ctx(ctx)

	if in.Title != "" {
		m = m.WhereLike("title", "%"+in.Title+"%")
	}
	if in.Status >= 0 {
		m = m.Where("status", in.Status)
	}

	// 获取总数
	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	// 分页查询
	var books []*entity.Books
	err = m.Page(in.Page, in.Size).OrderDesc("id").Scan(&books)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		List:  books,
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}, nil
}

// UpdateInput defines input for updating a book.
type UpdateInput struct {
	Id          uint64
	Title       string
	Author      string
	Isbn        string
	PublishDate string
	Status      int8
}

// Update modifies book information.
func (s *Service) Update(ctx context.Context, in UpdateInput) error {
	data := g.Map{}
	if in.Title != "" {
		data["title"] = in.Title
	}
	if in.Author != "" {
		data["author"] = in.Author
	}
	if in.Isbn != "" {
		// 检查ISBN是否与其他书冲突
		count, err := dao.Books.Ctx(ctx).Where("isbn", in.Isbn).WhereNot("id", in.Id).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("该ISBN号已被其他图书使用")
		}
		data["isbn"] = in.Isbn
	}
	if in.PublishDate != "" {
		data["publish_date"] = in.PublishDate
	}
	if in.Status >= 0 {
		data["status"] = in.Status
	}

	if len(data) == 0 {
		return errors.New("没有需要更新的内容")
	}

	_, err := dao.Books.Ctx(ctx).Where("id", in.Id).Data(data).Update()
	return err
}

// Delete removes a book by ID.
func (s *Service) Delete(ctx context.Context, id uint64) error {
	// 检查是否有未归还的借阅记录
	count, err := dao.Borrows.Ctx(ctx).Where("book_id", id).Where("return_at IS NULL").Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该图书尚有未归还的借阅记录，无法删除")
	}

	_, err = dao.Books.Ctx(ctx).Where("id", id).Delete()
	return err
}
