import request from './request'

export const sendCodeApi = (email, type) =>
  request.post('/auth/verification-codes', { email, type })

export const registerApi = (data) =>
  request.post('/auth/register', data)

export const loginApi = (data) =>
  request.post('/auth/login', data)

export const resetPasswordApi = (data) =>
  request.put('/auth/password', data)

export const getProfileApi = () =>
  request.get('/users/me')
