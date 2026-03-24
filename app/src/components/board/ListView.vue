<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Task, Priority } from '@/types/task'
import { PRIORITY_CONFIG } from '@/types/task'
import { formatDeadline } from '@/utils/date'
import { getInitials } from '@/utils/format'

const { t } = useI18n()

const props = defineProps<{ tasks: Task[] }>()
const emit = defineEmits<{ taskClick: [task: Task] }>()

const sortedTasks = computed(() =>
  [...props.tasks].sort((a, b) => {
    const priorityOrder: Record<Priority, number> = { urgent: 0, high: 1, medium: 2, low: 3 }
    return priorityOrder[a.priority] - priorityOrder[b.priority]
  })
)

function priorityColor(p: Priority): string { return PRIORITY_CONFIG[p].color }
</script>

<template>
  <div class="rounded-xl border border-white/5 bg-[#0f0f1a] overflow-hidden">
    <!-- Header -->
    <div class="grid grid-cols-[1fr_100px_100px_120px_50px] gap-2 px-4 py-2.5 border-b border-white/5 text-[10px] font-semibold uppercase tracking-widest text-white/25">
      <span>{{ t('task.title') }}</span>
      <span>{{ t('task.priority') }}</span>
      <span>{{ t('task.assignee') }}</span>
      <span>{{ t('task.deadline') }}</span>
      <span></span>
    </div>

    <!-- Rows -->
    <div v-if="sortedTasks.length > 0" class="divide-y divide-white/[0.03] max-h-[500px] overflow-y-auto">
      <button
        v-for="task in sortedTasks"
        :key="task.id"
        class="w-full grid grid-cols-[1fr_100px_100px_120px_50px] gap-2 items-center px-4 py-2.5 hover:bg-white/[0.02] transition-colors text-left"
        @click="emit('taskClick', task)"
      >
        <!-- Title -->
        <div class="flex items-center gap-2 min-w-0">
          <div class="w-1.5 h-1.5 rounded-full flex-shrink-0" :style="{ backgroundColor: priorityColor(task.priority) }" />
          <span class="text-sm truncate">{{ task.title }}</span>
          <v-icon v-if="task.completedAt" icon="mdi-check-circle" size="14" color="#10b981" class="flex-shrink-0" />
        </div>

        <!-- Priority -->
        <span class="text-[11px] font-medium" :style="{ color: priorityColor(task.priority) }">
          {{ t(`task.priorities.${task.priority}`) }}
        </span>

        <!-- Assignee -->
        <div v-if="task.assignee" class="flex items-center gap-1.5">
          <div class="w-5 h-5 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center flex-shrink-0">
            <span class="text-[8px] text-white font-semibold">{{ getInitials(task.assignee.name) }}</span>
          </div>
          <span class="text-[11px] text-white/40 truncate">{{ task.assignee.name.split(' ')[0] }}</span>
        </div>
        <span v-else class="text-[11px] text-white/20">—</span>

        <!-- Deadline -->
        <span
          v-if="task.deadline"
          class="text-[11px]"
          :class="{
            'text-error': formatDeadline(task.deadline).isOverdue,
            'text-warning': formatDeadline(task.deadline).isUrgent,
            'text-white/30': !formatDeadline(task.deadline).isOverdue && !formatDeadline(task.deadline).isUrgent,
          }"
        >
          {{ formatDeadline(task.deadline).text }}
        </span>
        <span v-else class="text-[11px] text-white/15">—</span>

        <!-- Arrow -->
        <v-icon icon="mdi-chevron-right" size="14" class="text-white/10 justify-self-end" />
      </button>
    </div>

    <div v-else class="py-12 text-center">
      <p class="text-sm text-white/25">{{ t('dashboard.noTasksInCategory') }}</p>
    </div>
  </div>
</template>
