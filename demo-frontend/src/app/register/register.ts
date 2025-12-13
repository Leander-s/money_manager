import { ChangeDetectorRef, Component, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { environment } from '../../environments/environment';

type RegisterRequest = {
  username: string | null,
  email: string,
  password: string
}

@Component({
  selector: 'app-register',
  imports: [CommonModule, FormsModule],
  templateUrl: './register.html',
  styleUrl: './register.css',
})
export class Register {
  constructor(
    private http: HttpClient,
  ) {}
  private cdc = inject(ChangeDetectorRef)

  username: string | null = null
  email: string = ""
  password: string = ""
  error: string = ""
  missingFields: string[] = []
  loading: boolean = false

  register() {
    this.loading = true
    if (this.email.length === 0) {
      this.missingFields.push("Email")
    }
    if (this.password.length === 0) {
      this.missingFields.push("Password")
    }

    if (this.missingFields.length !== 0) {
      this.error = "Missing fields:" + this.missingFields
      return
    }

    if (!this.username) {
      this.username = null
    }

    const request: RegisterRequest = {
      username: this.username,
      email: this.email,
      password: this.password
    }

    this.http.post(environment.API_URL + '/register', request).subscribe({
      next: () => {
        this.loading = false
        this.error = ""
        this.username = null
        this.email = ""
        this.password = ""
        this.cdc.markForCheck()
      },
      error: (err) => {
        this.loading = false
        this.error = err
        this.cdc.markForCheck()
      }
    })
  }
}
