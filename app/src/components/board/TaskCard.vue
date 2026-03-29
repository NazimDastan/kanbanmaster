<script setup lang="ts">
import type { Task } from '@/types/task'
import { formatDeadline } from '@/utils/date'
import { getInitials } from '@/utils/format'
import { computed } from 'vue'

const props = defineProps<{ task: Task }>()
const emit = defineEmits<{ click: [task: Task] }>()

const deadlineInfo = computed(() => props.task.deadline ? formatDeadline(props.task.deadline) : null)
const completedSubtasks = computed(() => props.task.subtasks?.filter((s) => s.isCompleted).length ?? 0)
const totalSubtasks = computed(() => props.task.subtasks?.length ?? 0)
const allLabels = computed(() => props.task.labels ?? [])
const allAssignees = computed(() => props.task.assignees ?? [])
const displayedAssignees = computed(() => allAssignees.value.slice(0, 3))
const extraAssigneeCount = computed(() => Math.max(0, allAssignees.value.length - 3))

function truncateLabel(name: string, max = 15): string {
  return name.length > max ? name.slice(0, max) + '…' : name
}

const priorityColors: Record<string, string> = {
  urgent: '#dc2626', high: '#06b6d4', medium: '#6366f1', low: '#64748b',
}
</script>

<template>
  <div
    class="rounded-lg bg-elevated-80 border border-white/[0.04] cursor-pointer transition-all duration-150 hover:border-white/10 hover:bg-card-hover group"
    @click="emit('click', task)"
  >
    <div class="h-[2px] rounded-t-lg" :style="{ background: priorityColors[task.priority] ?? '#6366f1' }" />
    <div class="p-2.5 space-y-2">
      <!-- Title -->
      <p class="text-[13px] font-medium leading-snug line-clamp-2 group-hover:text-primary-light transition-colors">{{ task.title }}</p>

      <!-- Middle row: deadline + subtasks + assignees -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <div v-if="deadlineInfo" class="flex items-center gap-0.5 text-[10px]" :class="{ 'text-error': deadlineInfo.isOverdue, 'text-warning': deadlineInfo.isUrgent, 'text-[var(--text-muted)]': !deadlineInfo.isOverdue && !deadlineInfo.isUrgent }">
            <v-icon icon="mdi-clock-outline" size="11" />
            {{ deadlineInfo.text }}
          </div>
          <div v-if="totalSubtasks > 0" class="flex items-center gap-0.5 text-[10px] text-[var(--text-muted)]">
            <v-icon icon="mdi-check-box-outline" size="11" />
            {{ completedSubtasks }}/{{ totalSubtasks }}
          </div>
        </div>
        <!-- Assignees -->
        <div v-if="displayedAssignees.length > 0" class="flex -space-x-1.5">
          <div
            v-for="assignee in displayedAssignees"
            :key="assignee.id"
            class="w-5 h-5 rounded-full bg-gradient-to-br from-primary/70 to-secondary/70 flex items-center justify-center flex-shrink-0 ring-1 ring-[var(--bg-card)]"
            :title="assignee.name"
          >
            <span class="text-[8px] text-white font-semibold">{{ getInitials(assignee.name) }}</span>
          </div>
          <div v-if="extraAssigneeCount > 0" class="w-5 h-5 rounded-full bg-white/10 flex items-center justify-center flex-shrink-0 ring-1 ring-[var(--bg-card)]">
            <span class="text-[8px] text-[var(--text-muted)] font-semibold">+{{ extraAssigneeCount }}</span>
          </div>
        </div>
        <div v-else-if="task.assignee" class="w-5 h-5 rounded-full bg-gradient-to-br from-primary/70 to-secondary/70 flex items-center justify-center flex-shrink-0">
          <span class="text-[8px] text-white font-semibold">{{ getInitials(task.assignee.name) }}</span>
        </div>
      </div>

      <!-- Labels (EN ALTTA) -->
      <div v-if="allLabels.length > 0" class="flex flex-wrap gap-1 pt-1 border-t border-white/[0.04]">
        <span
          v-for="label in allLabels"
          :key="label.id"
          class="px-1.5 py-0.5 rounded text-[9px] font-medium"
          :style="{ backgroundColor: label.color + '18', color: label.color }"
        >{{ truncateLabel(label.name) }}</span>
      </div>
    </div>
  </div>
</template>
