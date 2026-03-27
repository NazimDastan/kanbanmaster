<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Task, Priority } from '@/types/task'
import { PRIORITY_CONFIG } from '@/types/task'
import { formatDate, formatDeadline } from '@/utils/date'
import { getInitials } from '@/utils/format'
import { taskService } from '@/services/taskService'
import { delegationService, type ActivityLog } from '@/services/reportService'
import { useToast } from '@/composables/useToast'
import ActivityTimeline from './ActivityTimeline.vue'
import CommentSection from './CommentSection.vue'
import AttachmentSection from './AttachmentSection.vue'
import LabelManager from './LabelManager.vue'

const { t } = useI18n()
const toast = useToast()
const props = defineProps<{ task: Task; boardId: string }>()
const emit = defineEmits<{ close: []; delete: [id: string]; delegate: [taskId: string]; updated: [] }>()

// Editable fields
const editTitle = ref(props.task.title)
const editDescription = ref(props.task.description ?? '')
const editPriority = ref<Priority>(props.task.priority)
const editDeadline = ref(props.task.deadline?.slice(0, 10) ?? '')
const saving = ref(false)
const hasChanges = ref(false)

const activities = ref<ActivityLog[]>([])

const priorityOptions = [
  { value: 'urgent', title: t('task.priorities.urgent'), color: '#dc2626' },
  { value: 'high', title: t('task.priorities.high'), color: '#06b6d4' },
  { value: 'medium', title: t('task.priorities.medium'), color: '#6366f1' },
  { value: 'low', title: t('task.priorities.low'), color: '#9E9E9E' },
]

// Priority-based theming
const priorityTheme = computed(() => {
  const themes: Record<Priority, { border: string; borderSoft: string; glow: string; bg: string }> = {
    urgent: { border: 'rgba(220, 38, 38, 0.5)', borderSoft: 'rgba(220, 38, 38, 0.18)', glow: '0 0 24px rgba(220, 38, 38, 0.2)', bg: 'rgba(220, 38, 38, 0.04)' },
    high: { border: 'rgba(6, 182, 212, 0.4)', borderSoft: 'rgba(6, 182, 212, 0.15)', glow: '0 0 18px rgba(6, 182, 212, 0.12)', bg: 'rgba(6, 182, 212, 0.03)' },
    medium: { border: 'rgba(99, 102, 241, 0.2)', borderSoft: 'rgba(99, 102, 241, 0.08)', glow: 'none', bg: 'transparent' },
    low: { border: 'rgba(255, 255, 255, 0.04)', borderSoft: 'rgba(255, 255, 255, 0.02)', glow: 'none', bg: 'transparent' },
  }
  return themes[editPriority.value]
})

// Track changes
watch([editTitle, editDescription, editPriority, editDeadline], () => {
  hasChanges.value =
    editTitle.value !== props.task.title ||
    editDescription.value !== (props.task.description ?? '') ||
    editPriority.value !== props.task.priority ||
    editDeadline.value !== (props.task.deadline?.slice(0, 10) ?? '')
})

async function loadActivities() {
  try { activities.value = await delegationService.getActivity(props.task.id) }
  catch { /* Not available */ }
}
onMounted(loadActivities)

async function handleSave() {
  if (!hasChanges.value) return
  saving.value = true
  try {
    await taskService.update(props.task.id, {
      title: editTitle.value,
      description: editDescription.value || undefined,
      priority: editPriority.value,
      deadline: editDeadline.value ? editDeadline.value + 'T00:00:00Z' : undefined,
    } as Partial<Task>)
    toast.success(t('common.save') + ' ✓')
    hasChanges.value = false
    emit('updated')
  } catch {
    toast.error('Failed to save')
  } finally {
    saving.value = false
  }
}

const rootStyle = computed(() => ({
  'box-shadow': priorityTheme.value.glow,
  '--p-border': priorityTheme.value.border,
  '--p-border-soft': priorityTheme.value.borderSoft,
  '--p-bg': priorityTheme.value.bg,
  '--p-color': PRIORITY_CONFIG[editPriority.value].color,
}))

const priorityConfig = computed(() => PRIORITY_CONFIG[editPriority.value])
const deadlineInfo = computed(() => props.task.deadline ? formatDeadline(props.task.deadline) : null)
const completedSubtasks = computed(() => props.task.subtasks?.filter((s) => s.isCompleted).length ?? 0)
const totalSubtasks = computed(() => props.task.subtasks?.length ?? 0)
</script>

<template>
  <div
    class="h-full flex flex-col transition-all duration-300"
    :style="[rootStyle, { background: 'var(--bg-card)' }]"
  >
    <!-- Priority top bar indicator -->
    <div class="h-[3px] flex-shrink-0 transition-all duration-300" :style="{ background: `linear-gradient(90deg, ${priorityConfig.color}, transparent)` }" />

    <!-- Header -->
    <div class="flex items-center justify-between px-5 py-3 flex-shrink-0" :style="{ borderBottom: `1px solid ${priorityTheme.border}` }">
      <div class="flex items-center gap-2 flex-1 mr-2">
        <!-- Save indicator -->
        <div v-if="hasChanges" class="w-2 h-2 rounded-full bg-warning animate-pulse flex-shrink-0" />
        <span class="text-xs" :style="{ color: 'var(--text-muted)' }">{{ t('task.title') }}</span>
      </div>
      <div class="flex items-center gap-1">
        <!-- Save button -->
        <button
          v-if="hasChanges"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium transition-all"
          :class="saving ? 'bg-white/5 text-white/30' : 'bg-primary/20 text-primary-light hover:bg-primary/30'"
          :disabled="saving"
          @click="handleSave"
        >
          <v-icon :icon="saving ? 'mdi-loading mdi-spin' : 'mdi-content-save-outline'" size="14" />
          {{ t('common.save') }}
        </button>

        <v-menu>
          <template #activator="{ props: menuProps }">
            <button v-bind="menuProps" class="p-1.5 rounded-lg transition-colors">
              <v-icon icon="mdi-dots-horizontal" size="18" :style="{ color: 'var(--text-muted)' }" />
            </button>
          </template>
          <v-list density="compact" min-width="180" class="rounded-xl" :style="{ background: 'var(--bg-card)', border: '1px solid var(--border)' }">
            <v-list-item prepend-icon="mdi-swap-horizontal" :title="t('task.delegateTask')" @click="emit('delegate', task.id)" />
            <v-divider class="opacity-10" />
            <v-list-item prepend-icon="mdi-delete-outline" :title="t('task.deleteTask')" class="text-error" @click="emit('delete', task.id)" />
          </v-list>
        </v-menu>
        <button class="p-1.5 rounded-lg transition-colors" @click="emit('close')">
          <v-icon icon="mdi-close" size="18" :style="{ color: 'var(--text-muted)' }" />
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto">
      <!-- Title (editable) -->
      <div class="px-5 py-4 section-border">
        <input
          v-model="editTitle"
          class="w-full bg-transparent text-base font-bold outline-none"
          :style="{ color: 'var(--text)', '--placeholder-color': 'var(--text-muted)' }"
          :placeholder="t('task.title')"
        />
      </div>

      <!-- Priority & Deadline (editable) -->
      <div class="px-5 py-4 section-border">
        <div class="grid grid-cols-2 gap-3">
          <!-- Priority -->
          <div>
            <p class="text-[10px] font-semibold uppercase tracking-widest text-white/25 mb-2">{{ t('task.priority') }}</p>
            <div class="flex flex-wrap gap-1.5">
              <button
                v-for="opt in priorityOptions"
                :key="opt.value"
                class="flex items-center gap-1 px-2.5 py-1.5 rounded-lg text-[11px] font-medium transition-all border"
                :class="editPriority === opt.value ? 'border-white/20' : 'border-transparent hover:bg-white/5'"
                :style="editPriority === opt.value ? { backgroundColor: opt.color + '20', color: opt.color } : { color: 'rgba(255,255,255,0.4)' }"
                @click="editPriority = opt.value as Priority"
              >
                {{ opt.title }}
              </button>
            </div>
          </div>

          <!-- Deadline -->
          <div>
            <p class="text-[10px] font-semibold uppercase tracking-widest text-white/25 mb-2">{{ t('task.deadline') }}</p>
            <input
              v-model="editDeadline"
              type="date"
              class="w-full bg-white/5 rounded-lg px-3 py-1.5 text-sm text-white/70 outline-none transition-all duration-300 field-border"
            />
            <div
              v-if="deadlineInfo"
              class="mt-1.5 text-[11px]"
              :class="{ 'text-error': deadlineInfo.isOverdue, 'text-warning': deadlineInfo.isUrgent, 'text-white/25': !deadlineInfo.isOverdue && !deadlineInfo.isUrgent }"
            >
              {{ deadlineInfo.text }}
            </div>
          </div>
        </div>
      </div>

      <!-- Assignee -->
      <div class="px-5 py-4 section-border">
        <p class="text-[10px] font-semibold uppercase tracking-widest text-white/25 mb-2">{{ t('task.assignee') }}</p>
        <div v-if="task.assignee" class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center">
            <span class="text-[10px] text-white font-semibold">{{ getInitials(task.assignee.name) }}</span>
          </div>
          <p class="text-sm font-medium">{{ task.assignee.name }}</p>
        </div>
        <p v-else class="text-sm text-white/30">{{ t('task.unassigned') }}</p>
      </div>

      <!-- Description (editable) -->
      <div class="px-5 py-4 section-border">
        <p class="text-[10px] font-semibold uppercase tracking-widest text-white/25 mb-2">{{ t('task.description') }}</p>
        <textarea
          v-model="editDescription"
          :placeholder="t('task.noDescription')"
          rows="3"
          class="w-full bg-white/[0.03] rounded-lg px-3 py-2.5 text-sm text-white/60 leading-relaxed outline-none transition-all duration-300 resize-none field-border"
        />
      </div>

      <!-- Subtasks -->
      <div v-if="totalSubtasks > 0" class="px-5 py-4 section-border">
        <div class="flex items-center justify-between mb-3">
          <p class="text-[10px] font-semibold uppercase tracking-widest text-white/25">{{ t('task.subtasks') }}</p>
          <span class="text-[11px] text-white/30">{{ completedSubtasks }}/{{ totalSubtasks }}</span>
        </div>
        <v-progress-linear :model-value="(completedSubtasks / totalSubtasks) * 100" color="primary" height="4" rounded bg-color="surface" class="mb-3" />
        <div class="space-y-0.5">
          <div v-for="sub in task.subtasks" :key="sub.id" class="flex items-center gap-2.5 py-1.5 px-2 rounded-lg hover:bg-white/[0.02] transition-colors">
            <v-checkbox-btn :model-value="sub.isCompleted" density="compact" color="primary" />
            <span class="text-sm" :class="sub.isCompleted ? 'line-through text-white/25' : 'text-white/70'">{{ sub.title }}</span>
          </div>
        </div>
      </div>

      <!-- Labels -->
      <div class="px-5 py-4 section-border">
        <LabelManager :board-id="boardId" :task-id="task.id" :current-labels="task.labels ?? []" @updated="emit('updated')" />
      </div>

      <!-- Comments -->
      <div class="px-5 py-4 section-border">
        <CommentSection :task-id="task.id" />
      </div>

      <!-- Attachments -->
      <div class="px-5 py-4 section-border">
        <AttachmentSection :task-id="task.id" />
      </div>

      <!-- Activity -->
      <div class="px-5 py-4 section-border">
        <ActivityTimeline :activities="activities" />
      </div>

      <!-- Meta -->
      <div class="px-5 py-4">
        <div class="flex items-center gap-4 text-[11px] text-white/20">
          <div class="flex items-center gap-1.5">
            <v-icon icon="mdi-calendar-plus" size="13" />
            {{ formatDate(task.createdAt) }}
          </div>
          <div class="flex items-center gap-1.5">
            <v-icon icon="mdi-calendar-edit" size="13" />
            {{ formatDate(task.updatedAt) }}
          </div>
        </div>
      </div>
    </div>

    <!-- Sticky save bar (when changes exist) -->
    <div v-if="hasChanges" class="px-5 py-3 flex-shrink-0 section-border-top" :style="{ background: 'var(--bg-card)' }">
      <div class="flex items-center justify-between">
        <p class="text-xs text-warning flex items-center gap-1.5">
          <v-icon icon="mdi-circle-medium" size="16" />
          {{ t('common.save') }}
        </p>
        <div class="flex gap-2">
          <button
            class="px-3 py-1.5 rounded-lg text-xs text-white/40 hover:text-white/60 hover:bg-white/5 transition-all"
            @click="editTitle = task.title; editDescription = task.description ?? ''; editPriority = task.priority; editDeadline = task.deadline?.slice(0, 10) ?? ''; hasChanges = false"
          >
            {{ t('common.cancel') }}
          </button>
          <button
            class="px-4 py-1.5 rounded-lg text-xs font-medium text-white transition-all"
            :class="saving ? 'bg-primary/30' : 'bg-primary hover:bg-primary-dark'"
            :disabled="saving"
            @click="handleSave"
          >
            {{ saving ? t('common.loading') : t('common.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.section-border {
  border-bottom: 1px solid var(--p-border, rgba(255,255,255,0.04));
  background: var(--p-bg, transparent);
  transition: border-color 0.3s, background 0.3s;
}
.section-border-top {
  border-top: 1px solid var(--p-border, rgba(255,255,255,0.04));
  transition: border-color 0.3s;
}
.field-border {
  border: 1px solid var(--p-border-soft, rgba(255,255,255,0.05));
  transition: border-color 0.3s;
}
.field-border:focus {
  border-color: var(--p-border, rgba(99, 102, 241, 0.4));
}
</style>
