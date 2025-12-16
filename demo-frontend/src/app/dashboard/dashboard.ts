import { Component, ChangeDetectorRef, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { UserComponent, User } from '../user/user';
import { environment } from '../../environments/environment';
import { EntryComponent, Entry } from '../entry/entry';
import { Auth } from '../auth/auth';

@Component({
  selector: 'app-dashboard',
  imports: [CommonModule, FormsModule, UserComponent, EntryComponent],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.css',
})
export class Dashboard {
  private http = inject(HttpClient)
  private cdr = inject(ChangeDetectorRef)
  private auth = inject(Auth)
  users: User[] = []
  entries: Entry[] = []

  currentUser: User | null = null;

  loading = false
  error: string | null = null
  userError: string | null = null
  entryError: string | null = null

  newUsername = ''
  newEmail = ''
  newPassword = ''

  balance = 0
  ratio = 0.5

  constructor() {
    this.refresh()
  }

  addUser() {
    this.loading = true
    if (!this.newUsername || !this.newEmail || !this.newPassword) {
      this.error = 'All fields are required'
      this.loading = false
      this.cdr.markForCheck()
      return
    }

    const newUser: User = {
      id: undefined,
      username: this.newUsername,
      email: this.newEmail,
      password: this.newPassword
    }

    this.http.post<User>(`${environment.API_URL}/user`, newUser)
      .subscribe({
        next: () => {
          this.newUsername = ''
          this.newEmail = ''
          this.newPassword = ''
          this.refresh()
        }, error: (err) => {
          this.error = this.handleError(err)
          this.loading = false
          this.cdr.markForCheck()
        }
      })
  }

  refresh() {
    this.loading = true
    this.error = null
    this.loadCurrentUser()
    this.loadUsers()
    this.loadEntries()
  }

  deleteUser(user: User) {
    if (!user.id) return
    this.http.delete<User>(`${environment.API_URL}/user/${user.id}`).subscribe({
      next: () => {
        this.refresh()
        this.cdr.markForCheck()
      },
      error: (err) => {
        this.error = this.handleError(err)
        this.cdr.markForCheck()
      }
    });
  }

  loadUsers() {
    this.userError = null
    this.http.get<User[]>(`${environment.API_URL}/user`)
      .subscribe({
        next: (users) => {
          this.users = users
          this.loading = false
          this.cdr.markForCheck()
        }, error: (err) => {
          console.log(err)
          this.userError = this.handleError(err)
          this.loading = false
          this.cdr.markForCheck()
        }
      });

  }

  loadCurrentUser() {
    this.http.get<User>(`${environment.API_URL}/user/self`)
      .subscribe({
        next: (user) => {
          this.currentUser = user
          this.loading = false
          this.cdr.markForCheck()
        }, error: (err) => {
          console.log(err)
          this.userError = this.handleError(err)
          this.loading = false
          this.cdr.markForCheck()
        }
      });
  }

  addBalance() {
    this.error = null
    if (!this.isRatioValid()) {
      this.error = 'Ratio must be between 0 and 1'
      this.cdr.markForCheck()
      return
    }
    this.loading = true
    const newEntry: Entry = {
      balance: this.balance,
      ratio: this.ratio
    }

    this.http.post<Entry>(`${environment.API_URL}/balance`, newEntry)
      .subscribe({
        next: () => {
          this.balance = 0
          this.refresh()
        }, error: (err) => {
          this.entryError = this.handleError(err)
          this.loading = false
          this.cdr.markForCheck()
        }
      })
  }

  isRatioValid() {
    return Number.isFinite(this.ratio) && this.ratio >= 0 && this.ratio <= 1
  }

  loadEntries() {
    this.entryError = null
    this.http.get<Entry[]>(`${environment.API_URL}/balance`)
      .subscribe({
        next: (entries) => {
          this.entries = entries
          this.loading = false
          this.cdr.markForCheck()
        }, error: (err) => {
          console.log(err)
          this.entryError = this.handleError(err)
          this.loading = false
          this.cdr.markForCheck()
        }
      });
  }

  handleError(err: any): string {
    switch (err.status) {
      case 400:
        return 'Failed to load users'
      case 401:
        if (this.auth.isLoggedIn()) {
          this.auth.clearToken()
          return 'Session expired. Please log in again.'
        } else {
          return 'Please log in to continue'
        }
      case 403:
        return 'You are not authorized to perform this action'
      case 404:
        return 'Not Found'
      case 500:
        return 'Internal Server Error'
      default:
        return 'An unknown error occurred'
    }
  }
}
