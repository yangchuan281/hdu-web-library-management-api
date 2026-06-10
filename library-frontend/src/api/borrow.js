import request from './request'

export const getBorrowsApi = (params) =>
  request.get('/borrows', { params })

export const borrowBookApi = (bookId) =>
  request.post('/borrows', { book_id: bookId })

export const returnBookApi = (id) =>
  request.put(`/borrows/${id}/return`)
