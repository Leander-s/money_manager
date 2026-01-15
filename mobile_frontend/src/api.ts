const API_URL = 'http://localhost:8080';

export type AuthToken = {
  token: string;
  userID: number;
  expiry: string;
};

export type User = {
  id: number;
  username: string;
  email: string;
  password: string;
  created_at: string;
};

export type MoneyEntry = {
  id: number;
  balance: number;
  budget: number;
  ratio: number;
  created_at: string;
  user_id: number;
};

export type MoneyEntryRequest = {
  balance: number;
  ratio: number;
};

export type LoginRequest = {
  email: string;
  password: string;
};

export type RegisterRequest = {
  username: string | null;
  email: string;
  password: string;
};

async function requestJson<T>(path: string, options: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, options);
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed with status ${response.status}`);
  }
  return response.json() as Promise<T>;
}

async function requestNoContent(path: string, options: RequestInit): Promise<void> {
  const response = await fetch(`${API_URL}${path}`, options);
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed with status ${response.status}`);
  }
}

export async function login(request: LoginRequest): Promise<AuthToken> {
  return requestJson<AuthToken>('/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
}

export async function register(request: RegisterRequest): Promise<void> {
  return requestNoContent('/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
}

export async function fetchCurrentUser(token: AuthToken): Promise<User> {
  return requestJson<User>('/user/self', {
    method: 'GET',
    headers: { Authorization: `Bearer ${token.token}` },
  });
}

export async function fetchLatestMoneyEntry(token: AuthToken): Promise<MoneyEntry | null> {
  const entries = await requestJson<MoneyEntry[]>('/balance/1', {
    method: 'GET',
    headers: { Authorization: `Bearer ${token.token}` },
  });
  return entries[0] ?? null;
}

export async function createMoneyEntry(
  token: AuthToken,
  request: MoneyEntryRequest
): Promise<MoneyEntry> {
  return requestJson<MoneyEntry>('/balance', {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token.token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });
}
