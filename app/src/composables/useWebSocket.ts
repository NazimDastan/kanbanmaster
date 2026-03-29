import { ref, onUnmounted } from 'vue'
import { useNotificationStore } from '@/stores/useNotificationStore'
import type { Notification } from '@/types/notification'

const socket = ref<WebSocket | null>(null)
const connected = ref(false)
let reconnectTimer: ReturnType<typeof setTimeout> | null = null

function getWsUrl(): string {
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = import.meta.env.VITE_API_URL
    ? new URL(import.meta.env.VITE_API_URL).host
    : window.location.host
  return `${proto}//${host}/ws/notifications`
}

export function useWebSocket() {
  const notificationStore = useNotificationStore()

  function connect() {
    const token = localStorage.getItem('access_token')
    if (!token || socket.value?.readyState === WebSocket.OPEN) return

    const url = `${getWsUrl()}?token=${encodeURIComponent(token)}`
    const ws = new WebSocket(url)

    ws.onopen = () => {
      connected.value = true
      if (reconnectTimer) {
        clearTimeout(reconnectTimer)
        reconnectTimer = null
      }
    }

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        if (msg.type === 'notification') {
          notificationStore.addNotification(msg.payload as Notification)
        }
      } catch {
        // ignore malformed messages
      }
    }

    ws.onclose = () => {
      connected.value = false
      socket.value = null
      scheduleReconnect()
    }

    ws.onerror = () => {
      ws.close()
    }

    socket.value = ws
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    socket.value?.close()
    socket.value = null
    connected.value = false
  }

  function scheduleReconnect() {
    if (!reconnectTimer) {
      reconnectTimer = setTimeout(() => {
        reconnectTimer = null
        connect()
      }, 3000)
    }
  }

  onUnmounted(() => {
    disconnect()
  })

  return { connected, connect, disconnect }
}
