import { Component, inject, ChangeDetectorRef } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { environment } from '../../environments/environment';
import { Auth, AuthToken } from '../auth/auth';

type LoginRequest = {
  email: string,
  password: string
}

@Component({
  selector: 'app-login',
  imports: [CommonModule, FormsModule],
  templateUrl: './login.html',
  styleUrl: './login.css',
})
export class Login {
  constructor(
    private http: HttpClient,
    private auth: Auth
  ) {}

  private cdr = inject(ChangeDetectorRef)
  email: string = ""
  password: string = ""

  loading: boolean = false
  error: string | null = null

  submit() {
    this.loading = true
    var request: LoginRequest = {
      email: this.email, password: this.password
    }

    this.http.post<AuthToken>(environment.API_URL + "/login", request).subscribe({
      next: (token) => {
        this.auth.setToken(token)
        this.loading = false
        this.cdr.markForCheck()
      },
      error: (err) => {
        this.loading = false
        this.error = 'Failed to login:' + err
        this.cdr.markForCheck()
      }
    })
  }

  logout() {
    this.email = ""
    this.password = ""
    this.auth.clearToken()
    this.cdr.markForCheck()
  }

  isLoggedIn() {
    return this.auth.isLoggedIn()
  }
}
