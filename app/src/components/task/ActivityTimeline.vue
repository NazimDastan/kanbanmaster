<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { ActivityLog } from '@/services/reportService'
import { formatRelativeTime } from '@/utils/date'

const { t } = useI18n()
defineProps<{ activities: ActivityLog[] }>()

const actionIcons: Record<string, string> = { created: 'mdi-plus-circle-outline', moved: 'mdi-arrow-right-bold', assigned: 'mdi-account-arrow-right', delegated: 'mdi-swap-horizontal', completed: 'mdi-check-circle-outline', commented: 'mdi-comment-outline' }
const actionColors: Record<string, string> = { created: 'primary', moved: 'warning', assigned: 'secondary', delegated: 'secondary', completed: 'success', commented: 'primary' }
</script>

<template>
  <div>
    <p class="text-xs font-semibold uppercase tracking-wider text-text-secondary mb-3">{{ t('activity.title') }}</p>
    <div v-if="activities.length === 0" class="text-sm text-text-secondary text-center py-4">{{ t('task.noActivity') }}</div>
    <div v-else class="space-y-3">
      <div v-for="activity in activities" :key="activity.id" class="flex items-start gap-3">
        <v-avatar :color="actionColors[activity.action] ?? 'grey'" size="28">
          <v-icon :icon="actionIcons[activity.action] ?? 'mdi-circle-outline'" size="14" color="white" />
        </v-avatar>
        <div class="flex-1 min-w-0">
          <p class="text-sm">
            <span class="font-medium">{{ activity.user?.name ?? 'Unknown' }}</span>
            <span class="text-text-secondary"> {{ activity.action }} this task</span>
          </p>
          <p class="text-xs text-text-secondary">{{ formatRelativeTime(activity.createdAt) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
