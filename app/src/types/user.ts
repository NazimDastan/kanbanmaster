export interface User {
  id: string
  email: string
  name: string
  avatarUrl: string | null
  createdAt: string
  updatedAt: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface AuthResponse {
  accessToken: string
  refreshToken: string
  user: User
}

export type UserRole = 'owner' | 'leader' | 'member' | 'viewer'
