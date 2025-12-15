import { inject } from '@angular/core';
import { HttpInterceptorFn } from '@angular/common/http';
import { Auth, AuthToken } from './auth'

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  if (req.url.includes('/login') || req.url.includes('/register')) {
    return next(req)
  }

  const auth = inject(Auth)
  const token: AuthToken | null = auth.getToken()

  if (!token) {
    return next(req)
  }

  const authReq = req.clone({
    setHeaders: { Authorization: `Bearer ${token.token}` }
  })

  return next(authReq)
};
