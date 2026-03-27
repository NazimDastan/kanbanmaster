<script setup lang="ts">
import type { Task } from '@/types/task'
import type { Priority } from '@/types/task'
import { formatDeadline } from '@/utils/date'
import { getInitials } from '@/utils/format'
import { computed } from 'vue'

const props = defineProps<{ task: Task }>()
const emit = defineEmits<{ click: [task: Task] }>()

const deadlineInfo = computed(() => props.task.deadline ? formatDeadline(props.task.deadline) : null)
const completedSubtasks = computed(() => props.task.subtasks?.filter((s) => s.isCompleted).length ?? 0)
const totalSubtasks = computed(() => props.task.subtasks?.length ?? 0)
const visibleLabels = computed(() => props.task.labels?.slice(0, 3) ?? [])
const extraLabelCount = computed(() => Math.max(0, (props.task.labels?.length ?? 0) - 3))

const priorityColors: Record<string, string> = {
  urgent: '#dc2626', high: '#06b6d4', medium: '#6366f1', low: '#64748b',
}
</script>

<template>
  <div
    class="rounded-lg bg-[#161625]/80 border border-white/[0.04] cursor-pointer transition-all duration-150 hover:border-white/10 hover:bg-[#1a1a2e] group"
    @click="emit('click', task)"
  >
    <div class="h-[2px] rounded-t-lg" :style="{ background: priorityColors[task.priority] ?? '#6366f1' }" />
    <div class="p-2.5">
      <!-- Labels -->
      <div v-if="visibleLabels.length > 0" class="flex flex-wrap gap-1 mb-2">
        <span v-for="label in visibleLabels" :key="label.id" class="px-1.5 py-0.5 rounded text-[9px] font-medium" :style="{ backgroundColor: label.color + '18', color: label.color }">{{ label.name }}</span>
        <span v-if="extraLabelCount > 0" class="px-1.5 py-0.5 rounded text-[9px] text-white/25 bg-white/5">+{{ extraLabelCount }}</span>
      </div>

      <p class="text-[13px] font-medium leading-snug mb-2 line-clamp-2 group-hover:text-primary-light transition-colors">{{ task.title }}</p>

      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <div v-if="deadlineInfo" class="flex items-center gap-0.5 text-[10px]" :class="{ 'text-error': deadlineInfo.isOverdue, 'text-warning': deadlineInfo.isUrgent, 'text-white/25': !deadlineInfo.isOverdue && !deadlineInfo.isUrgent }">
            <v-icon icon="mdi-clock-outline" size="11" />
            {{ deadlineInfo.text }}
          </div>
          <div v-if="totalSubtasks > 0" class="flex items-center gap-0.5 text-[10px] text-white/25">
            <v-icon icon="mdi-check-box-outline" size="11" />
            {{ completedSubtasks }}/{{ totalSubtasks }}
          </div>
        </div>
        <div v-if="task.assignee" class="w-5 h-5 rounded-full bg-gradient-to-br from-primary/70 to-secondary/70 flex items-center justify-center">
          <span class="text-[8px] text-white font-semibold">{{ getInitials(task.assignee.name) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
