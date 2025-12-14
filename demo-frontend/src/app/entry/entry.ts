import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

export interface Entry {
  id: number;
  balance: number;
  budget: number;
  ratio: number;
  created_at: string;
  user_id: number;
  user_email: string;
}

@Component({
  selector: 'app-entry',
  imports: [CommonModule],
  templateUrl: './entry.html',
  styleUrl: './entry.css',
})
export class EntryComponent {
  @Input({ required: true }) entry!: Entry;
}
