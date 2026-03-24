import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Notification } from '@/types/notification'
import api from '@/services/api'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])
  const loading = ref(false)

  const unreadCount = computed(
    () => notifications.value.filter((n) => !n.isRead).length,
  )

  async function fetchNotifications() {
    loading.value = true
    try {
      const { data } = await api.get<Notification[]>('/notifications')
      notifications.value = data
    } finally {
      loading.value = false
    }
  }

  async function markAsRead(id: string) {
    await api.patch(`/notifications/${id}/read`)
    const notif = notifications.value.find((n) => n.id === id)
    if (notif) notif.isRead = true
  }

  async function markAllRead() {
    await api.patch('/notifications/read-all')
    notifications.value.forEach((n) => (n.isRead = true))
  }

  function addNotification(notification: Notification) {
    notifications.value.unshift(notification)
  }

  return { notifications, loading, unreadCount, fetchNotifications, markAsRead, markAllRead, addNotification }
})
