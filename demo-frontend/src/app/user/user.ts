import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';

export interface User {
  id?: number;
  username: string;
  email: string;
  password: string;
  created_at?: string;
}

@Component({
  selector: 'app-user',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './user.html',
  styleUrl: './user.css',
})
export class UserComponent {
  @Input({ required: true }) user!: User;

  @Output() onDelete = new EventEmitter<User>();

  deleteUser() {
    this.onDelete.emit(this.user);
  }

}
