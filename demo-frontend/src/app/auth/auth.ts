import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

const TOKEN_KEY = 'auth_token'

export type AuthToken = {
  token: string,
  userID: number,
  expiry: number
}

@Injectable({
  providedIn: 'root',
})
export class Auth {
  private tokenSubject = new BehaviorSubject<string | null>(
    localStorage.getItem(TOKEN_KEY)
  )

  token$ = this.tokenSubject

  setToken(token: AuthToken) {
    localStorage.setItem(TOKEN_KEY, JSON.stringify(token))
    this.tokenSubject.next(JSON.stringify(token))
  }

  getToken(): AuthToken | null {
    if (!this.tokenSubject.value) {
      return null
    }
    return JSON.parse(this.tokenSubject.value)
  }

  clearToken() {
    localStorage.removeItem(TOKEN_KEY)
    this.tokenSubject.next(null)
  }

  isLoggedIn(): boolean {
    return !!this.getToken()
  }
}
