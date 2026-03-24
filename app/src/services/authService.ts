import api from './api'
import type { AuthResponse, LoginRequest, RegisterRequest, User } from '@/types/user'

export const authService = {
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/auth/login', data)
    return response.data
  },

  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/auth/register', data)
    return response.data
  },

  async refresh(refreshToken: string): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/auth/refresh', { refreshToken })
    return response.data
  },

  async me(): Promise<User> {
    const response = await api.get<User>('/auth/me')
    return response.data
  },

  async updateAvatar(avatarUrl: string): Promise<User> {
    const response = await api.put<User>('/auth/avatar', { avatarUrl })
    return response.data
  },

  async updateProfile(name: string, email: string): Promise<User> {
    const response = await api.put<User>('/auth/profile', { name, email })
    return response.data
  },

  async changePassword(currentPassword: string, newPassword: string): Promise<void> {
    await api.put('/auth/password', { currentPassword, newPassword })
  },
}
