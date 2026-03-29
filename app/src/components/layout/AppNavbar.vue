<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { useNotificationStore } from '@/stores/useNotificationStore'
import { getInitials } from '@/utils/format'
import NotificationBell from '@/components/notification/NotificationBell.vue'
import LanguageSwitcher from './LanguageSwitcher.vue'
import InvitationPanel from '@/components/team/InvitationPanel.vue'
import { useTheme } from '@/composables/useTheme'
import { useConfirm } from '@/composables/useConfirm'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const { theme, toggle: toggleTheme } = useTheme()
const { confirm } = useConfirm()
const initials = computed(() => getInitials(authStore.user?.name ?? ''))

// Fetch notifications on mount + when user changes
onMounted(() => {
  if (authStore.isAuthenticated) notificationStore.fetchNotifications()
})
watch(() => authStore.isAuthenticated, (authed) => {
  if (authed) notificationStore.fetchNotifications()
})

function handleNotificationNavigate(referenceId: string) {
  router.push('/')
}

async function handleLogout() {
  const ok = await confirm({
    title: t('navbar.logout'),
    message: t('navbar.logoutConfirm'),
    confirmText: t('navbar.logout'),
    danger: true,
  })
  if (ok) authStore.logout()
}
</script>

<template>
  <header
    class="h-14 flex-shrink-0 flex items-center justify-between px-4 md:px-6 border-b bg-[var(--bg-card)] border-[var(--border)]"
  >
    <div class="flex items-center gap-2">
      <!-- Mobile hamburger -->
      <button class="md:hidden p-1.5 rounded-lg text-[var(--text-secondary)]" @click="$emit('toggleSidebar')">
        <v-icon icon="mdi-menu" size="22" />
      </button>
      <router-link to="/" class="flex items-center gap-2 hover:opacity-80 transition-opacity">
        <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-secondary flex items-center justify-center">
          <v-icon icon="mdi-view-dashboard" color="white" size="17" />
        </div>
        <span class="text-sm font-bold tracking-tight hidden sm:block text-[var(--text)]">KanbanMaster</span>
      </router-link>
    </div>

    <div class="flex items-center gap-1">
      <!-- Theme toggle -->
      <button class="p-2 rounded-lg transition-colors text-[var(--text-secondary)]" aria-label="Toggle theme" @click="toggleTheme">
        <v-icon :icon="theme === 'dark' ? 'mdi-weather-sunny' : 'mdi-weather-night'" size="18" />
      </button>

      <LanguageSwitcher />
      <InvitationPanel />

      <NotificationBell
        :count="notificationStore.unreadCount"
        :notifications="notificationStore.notifications"
        @mark-read="notificationStore.markAsRead"
        @mark-all-read="notificationStore.markAllRead"
        @navigate="handleNotificationNavigate"
      />

      <v-menu>
        <template #activator="{ props }">
          <button v-bind="props" class="flex items-center gap-2 px-2 py-1.5 rounded-lg transition-colors">
            <div class="w-7 h-7 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center overflow-hidden">
              <img v-if="authStore.user?.avatarUrl" :src="authStore.user.avatarUrl" class="w-full h-full object-cover" />
              <span v-else class="text-[10px] font-semibold text-white">{{ initials }}</span>
            </div>
            <span class="text-xs font-medium hidden md:block text-[var(--text-secondary)]">{{ authStore.userName }}</span>
            <v-icon icon="mdi-chevron-down" size="14" class="text-[var(--text-muted)]" />
          </button>
        </template>
        <v-list density="compact" min-width="160" class="rounded-xl bg-[var(--bg-card)] border border-[var(--border)]">
          <v-list-item to="/profile" prepend-icon="mdi-account-outline" :title="t('navbar.profile')" />
          <v-divider class="opacity-10" />
          <v-list-item prepend-icon="mdi-logout" :title="t('navbar.logout')" @click="handleLogout" />
        </v-list>
      </v-menu>
    </div>
  </header>
</template>
