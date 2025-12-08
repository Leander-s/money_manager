import { Component, ChangeDetectorRef, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { UserComponent, User } from './user/user';

export const API = 'http://localhost:8080';

@Component({
  selector: 'app-root',
  imports: [CommonModule, FormsModule, UserComponent],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  private http = inject(HttpClient);
  private cdr = inject(ChangeDetectorRef);
  users: User[] = []
  loading = false;
  error: string | null = null;

  newUsername = '';
  newEmail = '';
  newPassword = '';

  constructor() {
    this.refresh();
  }

  addUser() {
    this.loading = true;
    if (!this.newUsername || !this.newEmail || !this.newPassword) {
      this.error = 'All fields are required';
      this.loading = false;
      this.cdr.markForCheck();
      return;
    }

    const newUser: User = {
      ID: undefined,
      Username: this.newUsername,
      Email: this.newEmail,
      Password: this.newPassword
    };

    this.http.post<User>(API + '/user', newUser)
      .subscribe({
        next: () => {
          this.newUsername = '';
          this.newEmail = '';
          this.newPassword = '';
          this.refresh();
        }, error: () => {
          this.error = 'Failed to add user';
          this.loading = false;
          this.cdr.markForCheck();
        }
      });
  }

  refresh() {
    this.loading = true;
    this.error = null;
    this.http.get<User[]>(API + '/user')
      .subscribe({
        next: (users) => {
          this.users = users;
          this.loading = false;
          this.cdr.markForCheck();
        }, error: () => {
          this.error = 'Failed to load users';
          this.loading = false;
          this.cdr.markForCheck();
        }
      });
  }

  deleteUser(user: User) {
    if(!user.ID) return;
    this.http.delete<User>(`${API}/user/${user.ID}`).subscribe({
      next: () => {
        this.refresh();
        this.cdr.markForCheck();
      },
      error: () => {
        this.error = 'Failed to delete user';
        this.cdr.markForCheck();
      }
    });
  }
}
