// @vitest-environment jsdom
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createVuetify } from 'vuetify'
import { h } from 'vue'
import NotificationBell from './NotificationBell.vue'
import type { Notification } from '@/types/notification'

vi.mock('@/utils/date', () => ({
  formatRelativeTime: vi.fn(() => '5m ago'),
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'notifications.title': 'Notifications',
        'notifications.markAllRead': 'Mark all read',
        'notifications.noNotifications': 'No notifications',
      }
      return translations[key] ?? key
    },
  }),
}))

const vuetify = createVuetify()

function createNotification(overrides: Partial<Notification> = {}): Notification {
  return {
    id: 'n1',
    userId: 'u1',
    type: 'assigned',
    title: 'You were assigned a task',
    message: 'Task: Build login page',
    referenceId: 't1',
    isRead: false,
    createdAt: '2026-03-28T10:00:00Z',
    ...overrides,
  }
}

function mountBell(props: { count: number; notifications: Notification[] }) {
  return mount(NotificationBell, {
    props,
    global: {
      plugins: [vuetify],
      stubs: {
        'v-menu': {
          template: '<div><slot name="activator" v-bind="{ props: {} }" /><slot /></div>',
        },
        'v-card': {
          template: '<div><slot /></div>',
        },
        'v-btn': {
          template: '<button><slot /></button>',
        },
        'v-icon': {
          template: '<i />',
        },
      },
    },
  })
}

describe('NotificationBell', () => {
  it('shows unread count badge when count > 0', () => {
    // Arrange
    const notifications = [createNotification()]

    // Act
    const wrapper = mountBell({ count: 3, notifications })

    // Assert
    const badge = wrapper.find('span.absolute')
    expect(badge.exists()).toBe(true)
    expect(badge.text()).toBe('3')
  })

  it('hides badge when count is 0', () => {
    // Arrange & Act
    const wrapper = mountBell({ count: 0, notifications: [] })

    // Assert
    const badge = wrapper.find('span.absolute')
    expect(badge.exists()).toBe(false)
  })

  it('renders notification items in the list', () => {
    // Arrange
    const notifications = [
      createNotification({ id: 'n1', title: 'Task assigned to you' }),
      createNotification({ id: 'n2', title: 'Deadline approaching' }),
    ]

    // Act
    const wrapper = mountBell({ count: 2, notifications })

    // Assert
    expect(wrapper.text()).toContain('Task assigned to you')
    expect(wrapper.text()).toContain('Deadline approaching')
  })
})
