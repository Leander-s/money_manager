import { Component, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';

const API = 'http://localhost:8080';

interface User {
  id: number;
  username: string;
  email: string;
}

@Component({
  selector: 'app-root',
  imports: [],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  private http = inject(HttpClient);
  users: User[] = []
  loading = false;
  error: string | null = null;
  refresh() {
    this.loading = true;
    this.http.get<User[]>(API + '/users')
      .subscribe({
        next: (users) => {
          this.users = users;
          this.loading = false;
        }, error: () => {
          this.error = 'Failed to load users';
          this.loading = false;
        }
      });
  }
}
