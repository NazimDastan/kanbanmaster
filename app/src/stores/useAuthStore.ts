import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/user'
import { authService } from '@/services/authService'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!user.value)
  const userName = computed(() => user.value?.name ?? '')

  async function login(email: string, password: string) {
    loading.value = true
    try {
      const data = await authService.login({ email, password })
      localStorage.setItem('access_token', data.accessToken)
      localStorage.setItem('refresh_token', data.refreshToken)
      user.value = data.user
    } finally {
      loading.value = false
    }
  }

  async function register(name: string, email: string, password: string) {
    loading.value = true
    try {
      const data = await authService.register({ name, email, password })
      localStorage.setItem('access_token', data.accessToken)
      localStorage.setItem('refresh_token', data.refreshToken)
      user.value = data.user
    } finally {
      loading.value = false
    }
  }

  async function fetchUser() {
    const token = localStorage.getItem('access_token')
    if (!token) return
    try {
      user.value = await authService.me()
    } catch {
      logout()
    }
  }

  async function updateProfile(name: string, email: string) {
    const updated = await authService.updateProfile(name, email)
    user.value = updated
  }

  async function updateAvatar(avatarUrl: string) {
    const updated = await authService.updateAvatar(avatarUrl)
    user.value = updated
  }

  async function changePassword(currentPassword: string, newPassword: string) {
    await authService.changePassword(currentPassword, newPassword)
  }

  function logout() {
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    window.location.href = '/login'
  }

  return { user, loading, isAuthenticated, userName, login, register, fetchUser, updateProfile, updateAvatar, changePassword, logout }
})
