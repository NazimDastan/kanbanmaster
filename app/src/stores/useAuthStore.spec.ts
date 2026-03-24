// @vitest-environment jsdom
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './useAuthStore'

// Mock authService
vi.mock('@/services/authService', () => ({
  authService: {
    login: vi.fn().mockResolvedValue({
      accessToken: 'test-access-token',
      refreshToken: 'test-refresh-token',
      user: { id: '1', name: 'Test User', email: 'test@example.com', avatarUrl: null, createdAt: '', updatedAt: '' },
    }),
    register: vi.fn().mockResolvedValue({
      accessToken: 'test-access-token',
      refreshToken: 'test-refresh-token',
      user: { id: '1', name: 'New User', email: 'new@example.com', avatarUrl: null, createdAt: '', updatedAt: '' },
    }),
    me: vi.fn().mockResolvedValue({
      id: '1', name: 'Test User', email: 'test@example.com', avatarUrl: null, createdAt: '', updatedAt: '',
    }),
  },
}))

describe('useAuthStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('starts with no user', () => {
    const store = useAuthStore()
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })

  it('login sets user and tokens', async () => {
    const store = useAuthStore()
    await store.login('test@example.com', 'password')

    expect(store.user).not.toBeNull()
    expect(store.user?.name).toBe('Test User')
    expect(store.isAuthenticated).toBe(true)
    expect(localStorage.getItem('access_token')).toBe('test-access-token')
    expect(localStorage.getItem('refresh_token')).toBe('test-refresh-token')
  })

  it('register sets user and tokens', async () => {
    const store = useAuthStore()
    await store.register('New User', 'new@example.com', 'password')

    expect(store.user).not.toBeNull()
    expect(store.user?.name).toBe('New User')
    expect(store.isAuthenticated).toBe(true)
  })

  it('fetchUser loads user from API', async () => {
    localStorage.setItem('access_token', 'some-token')
    const store = useAuthStore()
    await store.fetchUser()

    expect(store.user?.name).toBe('Test User')
  })

  it('userName returns empty string when no user', () => {
    const store = useAuthStore()
    expect(store.userName).toBe('')
  })

  it('userName returns name when user exists', async () => {
    const store = useAuthStore()
    await store.login('test@example.com', 'password')
    expect(store.userName).toBe('Test User')
  })
})
