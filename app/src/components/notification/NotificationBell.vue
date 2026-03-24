<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Notification } from '@/types/notification'
import { formatRelativeTime } from '@/utils/date'

const { t } = useI18n()
defineProps<{ count: number; notifications: Notification[] }>()
const emit = defineEmits<{ markRead: [id: string]; markAllRead: [] }>()

const typeIcons: Record<string, string> = {
  assigned: 'mdi-account-arrow-right', delegated: 'mdi-swap-horizontal', deadline: 'mdi-clock-alert-outline',
  comment: 'mdi-comment-outline', report_request: 'mdi-file-document-outline', completed: 'mdi-check-circle-outline', overdue: 'mdi-alert-circle-outline',
}
const typeColors: Record<string, string> = {
  assigned: '#6366f1', delegated: '#a855f7', deadline: '#f59e0b',
  comment: '#6366f1', report_request: '#a855f7', completed: '#10b981', overdue: '#ef4444',
}
</script>

<template>
  <v-menu :close-on-content-click="false" max-width="360">
    <template #activator="{ props }">
      <button v-bind="props" class="relative p-2 rounded-lg hover:bg-white/5 transition-colors">
        <v-icon icon="mdi-bell-outline" size="20" class="text-white/50" />
        <span v-if="count > 0" class="absolute top-1 right-1 w-4 h-4 rounded-full bg-error text-[9px] font-bold text-white flex items-center justify-center">{{ count }}</span>
      </button>
    </template>
    <v-card min-width="340" max-height="440" color="#161625" class="rounded-xl border border-white/5">
      <div class="flex items-center justify-between px-4 py-3 border-b border-white/5">
        <span class="text-xs font-semibold">{{ t('notifications.title') }}</span>
        <v-btn v-if="count > 0" variant="text" size="x-small" color="primary" @click="emit('markAllRead')">{{ t('notifications.markAllRead') }}</v-btn>
      </div>
      <div v-if="notifications.length > 0" class="overflow-y-auto" style="max-height: 340px">
        <div v-for="n in notifications.slice(0, 20)" :key="n.id" class="flex items-start gap-3 px-4 py-3 hover:bg-white/[0.02] transition-colors cursor-pointer" :class="{ 'bg-primary/5': !n.isRead }" @click="emit('markRead', n.id)">
          <div class="w-7 h-7 rounded-lg flex-shrink-0 flex items-center justify-center mt-0.5" :style="{ backgroundColor: (typeColors[n.type] ?? '#6366f1') + '12' }">
            <v-icon :icon="typeIcons[n.type] ?? 'mdi-bell-outline'" size="14" :style="{ color: typeColors[n.type] ?? '#6366f1' }" />
          </div>
          <div class="min-w-0">
            <p class="text-xs font-medium leading-snug">{{ n.title }}</p>
            <p class="text-[10px] text-white/25 mt-0.5">{{ formatRelativeTime(n.createdAt) }}</p>
          </div>
        </div>
      </div>
      <div v-else class="py-8 text-center">
        <v-icon icon="mdi-bell-off-outline" size="32" class="text-white/10" />
        <p class="text-xs text-white/25 mt-2">{{ t('notifications.noNotifications') }}</p>
      </div>
    </v-card>
  </v-menu>
</template>
