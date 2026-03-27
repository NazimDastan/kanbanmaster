<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/useAuthStore'
import { useNotificationStore } from '@/stores/useNotificationStore'
import { getInitials } from '@/utils/format'
import NotificationBell from '@/components/notification/NotificationBell.vue'
import LanguageSwitcher from './LanguageSwitcher.vue'
import { useTheme } from '@/composables/useTheme'

const { t } = useI18n()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const { theme, toggle: toggleTheme } = useTheme()
const initials = computed(() => getInitials(authStore.user?.name ?? ''))
</script>

<template>
  <header
    class="h-14 flex-shrink-0 flex items-center justify-between px-4 md:px-6 border-b"
    :style="{ background: 'var(--bg-card)', borderColor: 'var(--border)' }"
  >
    <div class="flex items-center gap-3">
      <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-secondary flex items-center justify-center">
        <v-icon icon="mdi-view-dashboard" color="white" size="17" />
      </div>
      <span class="text-sm font-bold tracking-tight hidden sm:block" :style="{ color: 'var(--text)' }">KanbanMaster</span>
    </div>

    <div class="flex items-center gap-1">
      <!-- Theme toggle -->
      <button class="p-2 rounded-lg transition-colors" :style="{ color: 'var(--text-secondary)' }" @click="toggleTheme">
        <v-icon :icon="theme === 'dark' ? 'mdi-weather-sunny' : 'mdi-weather-night'" size="18" />
      </button>

      <LanguageSwitcher />

      <NotificationBell
        :count="notificationStore.unreadCount"
        :notifications="notificationStore.notifications"
        @mark-read="notificationStore.markAsRead"
        @mark-all-read="notificationStore.markAllRead"
      />

      <v-menu>
        <template #activator="{ props }">
          <button v-bind="props" class="flex items-center gap-2 px-2 py-1.5 rounded-lg transition-colors">
            <div class="w-7 h-7 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center overflow-hidden">
              <img v-if="authStore.user?.avatarUrl" :src="authStore.user.avatarUrl" class="w-full h-full object-cover" />
              <span v-else class="text-[10px] font-semibold text-white">{{ initials }}</span>
            </div>
            <span class="text-xs font-medium hidden md:block" :style="{ color: 'var(--text-secondary)' }">{{ authStore.userName }}</span>
            <v-icon icon="mdi-chevron-down" size="14" :style="{ color: 'var(--text-muted)' }" />
          </button>
        </template>
        <v-list density="compact" min-width="160" class="rounded-xl" :style="{ background: 'var(--bg-card)', border: '1px solid var(--border)' }">
          <v-list-item to="/profile" prepend-icon="mdi-account-outline" :title="t('navbar.profile')" />
          <v-divider class="opacity-10" />
          <v-list-item prepend-icon="mdi-logout" :title="t('navbar.logout')" @click="authStore.logout" />
        </v-list>
      </v-menu>
    </div>
  </header>
</template>
