<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Task, Priority } from '@/types/task'
import type { User } from '@/types/user'
import { PRIORITY_CONFIG } from '@/types/task'
import { formatDeadline } from '@/utils/date'
import { getInitials } from '@/utils/format'
import { taskService } from '@/services/taskService'
import { delegationService, type ActivityLog } from '@/services/reportService'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'
import { useBoardStore } from '@/stores/useBoardStore'
import MemberPicker from './MemberPicker.vue'
import CommentSection from './CommentSection.vue'
import AttachmentSection from './AttachmentSection.vue'
import LabelManager from './LabelManager.vue'
import SubtaskList from './SubtaskList.vue'
import ActivityTimeline from './ActivityTimeline.vue'

const { t } = useI18n()
const toast = useToast()
const { confirm } = useConfirm()
const boardStore = useBoardStore()
const props = defineProps<{ task: Task; boardId: string }>()
const emit = defineEmits<{ close: []; delete: [id: string]; delegate: [taskId: string]; updated: [] }>()

const editTitle = ref(props.task.title)
const editDescription = ref(props.task.description ?? '')
const editPriority = ref<Priority>(props.task.priority)
const editDeadline = ref(props.task.deadline?.slice(0, 10) ?? '')
const saving = ref(false)
const hasChanges = ref(false)
const activities = ref<ActivityLog[]>([])
const activeTab = ref<'comments' | 'attachments' | 'activity'>('comments')
const showAddAssignee = ref(false)

const teamId = computed(() => boardStore.currentBoard?.teamId ?? '')
const priorityConfig = computed(() => PRIORITY_CONFIG[editPriority.value])
const deadlineInfo = computed(() => props.task.deadline ? formatDeadline(props.task.deadline) : null)
const totalSubtasks = computed(() => props.task.subtasks?.length ?? 0)
const assignees = computed<User[]>(() => props.task.assignees ?? [])

const priorities: { value: Priority; color: string; icon: string }[] = [
  { value: 'urgent', color: '#dc2626', icon: 'mdi-alert-circle' },
  { value: 'high', color: '#06b6d4', icon: 'mdi-arrow-up-bold' },
  { value: 'medium', color: '#6366f1', icon: 'mdi-minus' },
  { value: 'low', color: '#64748b', icon: 'mdi-arrow-down-bold' },
]

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
    toast.error(t('common.error'))
  } finally {
    saving.value = false
  }
}

async function handleClose() {
  if (hasChanges.value) {
    const ok = await confirm({
      title: t('common.save'),
      message: t('task.unsavedChanges'),
      confirmText: t('task.discardChanges'),
      danger: true,
    })
    if (!ok) return
  }
  emit('close')
}

async function handleAssign(userId: string, userName: string) {
  try {
    await taskService.assign(props.task.id, userId)
    toast.success(`${userName} ✓`)
    showAddAssignee.value = false
    emit('updated')
  } catch {
    toast.error(t('common.error'))
  }
}

async function handleAddAssignee(userId: string, userName: string) {
  try {
    await taskService.addAssignee(props.task.id, userId)
    toast.success(`${userName} ✓`)
    showAddAssignee.value = false
    emit('updated')
  } catch {
    toast.error(t('common.error'))
  }
}

async function handleRemoveAssignee(userId: string) {
  try {
    await taskService.removeAssignee(props.task.id, userId)
    emit('updated')
  } catch {
    toast.error(t('common.error'))
  }
}
</script>

<template>
  <div class="h-full flex flex-col bg-[var(--bg-card)]">
    <!-- Priority stripe -->
    <div class="h-1 flex-shrink-0" :style="{ background: `linear-gradient(90deg, ${priorityConfig.color}, transparent)` }" />

    <!-- Header -->
    <div class="px-5 pt-4 pb-2 flex-shrink-0">
      <div class="flex items-start justify-between gap-2 mb-3">
        <input
          v-model="editTitle"
          class="flex-1 bg-transparent text-lg font-bold outline-none text-[var(--text)] placeholder:text-[var(--text-muted)]"
          :placeholder="t('task.title')"
        />
        <div class="flex items-center gap-0.5 flex-shrink-0">
          <v-menu>
            <template #activator="{ props: menuProps }">
              <button v-bind="menuProps" class="p-1.5 rounded-lg hover:bg-[var(--bg-input)] transition-colors">
                <v-icon icon="mdi-dots-horizontal" size="18" class="text-[var(--text-muted)]" />
              </button>
            </template>
            <v-list density="compact" min-width="180" class="rounded-xl bg-[var(--bg-card)] border border-[var(--border)]">
              <v-list-item prepend-icon="mdi-swap-horizontal" :title="t('task.delegateTask')" @click="emit('delegate', task.id)" />
              <v-divider class="opacity-10" />
              <v-list-item prepend-icon="mdi-delete-outline" :title="t('task.deleteTask')" class="text-error" @click="emit('delete', task.id)" />
            </v-list>
          </v-menu>
          <button class="p-1.5 rounded-lg hover:bg-[var(--bg-input)] transition-colors" @click="handleClose">
            <v-icon icon="mdi-close" size="18" class="text-[var(--text-muted)]" />
          </button>
        </div>
      </div>

      <!-- Save bar -->
      <div v-if="hasChanges" class="flex items-center justify-between mb-3 px-3 py-2 rounded-lg bg-primary/10 border border-primary/20">
        <span class="text-xs text-primary">{{ t('task.unsavedChanges').split('.')[0] }}</span>
        <button
          class="flex items-center gap-1 px-3 py-1 rounded-md text-xs font-medium bg-primary text-white hover:bg-primary-dark transition-all"
          :disabled="saving"
          @click="handleSave"
        >
          <v-icon :icon="saving ? 'mdi-loading mdi-spin' : 'mdi-check'" size="14" />
          {{ t('common.save') }}
        </button>
      </div>
    </div>

    <!-- Scrollable content -->
    <div class="flex-1 overflow-y-auto px-5 space-y-5 pb-5">

      <!-- === SECTION: Properties === -->
      <div class="space-y-4">
        <!-- Priority -->
        <div>
          <label class="text-[10px] font-semibold uppercase tracking-widest text-[var(--text-muted)] mb-2 block">{{ t('task.priority') }}</label>
          <div class="flex gap-2">
            <button
              v-for="p in priorities"
              :key="p.value"
              class="flex-1 flex items-center justify-center gap-1.5 py-2 rounded-lg text-xs font-medium border transition-all"
              :class="editPriority === p.value ? 'shadow-sm' : 'border-transparent opacity-40 hover:opacity-70'"
              :style="editPriority === p.value
                ? { backgroundColor: p.color + '15', color: p.color, borderColor: p.color + '40' }
                : {}"
              @click="editPriority = p.value"
            >
              <v-icon :icon="p.icon" size="14" />
              {{ t(`task.priorities.${p.value}`) }}
            </button>
          </div>
        </div>

        <!-- Assignees -->
        <div>
          <label class="text-[10px] font-semibold uppercase tracking-widest text-[var(--text-muted)] mb-2 block">{{ t('task.assignee') }}</label>

          <!-- Assignee list -->
          <div v-if="assignees.length > 0 || task.assignee" class="space-y-2 mb-2">
            <!-- Show all assignees -->
            <div v-for="member in assignees" :key="member.id" class="flex items-center gap-3 py-1">
              <div class="w-7 h-7 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center flex-shrink-0">
                <span class="text-[9px] text-white font-semibold">{{ getInitials(member.name) }}</span>
              </div>
              <p class="flex-1 text-sm text-[var(--text)]">{{ member.name }}</p>
              <button class="p-1 rounded hover:bg-error/10 transition-colors" @click="handleRemoveAssignee(member.id)">
                <v-icon icon="mdi-close" size="14" class="text-[var(--text-muted)] hover:text-error" />
              </button>
            </div>
            <!-- Fallback: primary assignee if no multi-assign data -->
            <div v-if="assignees.length === 0 && task.assignee" class="flex items-center gap-3 py-1">
              <div class="w-7 h-7 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center flex-shrink-0">
                <span class="text-[9px] text-white font-semibold">{{ getInitials(task.assignee.name) }}</span>
              </div>
              <p class="flex-1 text-sm text-[var(--text)]">{{ task.assignee.name }}</p>
            </div>
          </div>

          <!-- Add assignee button / picker -->
          <button
            v-if="!showAddAssignee"
            class="w-full flex items-center gap-2 px-3 py-2.5 rounded-xl border border-dashed transition-all text-sm border-[var(--border)] text-[var(--text-muted)] hover:border-primary/30 hover:bg-primary/5"
            @click="showAddAssignee = true"
          >
            <v-icon icon="mdi-account-plus-outline" size="18" />
            {{ t('task.assignTask') }}
          </button>
          <MemberPicker v-else :team-id="teamId" :current-assignee-id="null" @select="handleAddAssignee" />
        </div>

        <!-- Description -->
        <div>
          <label class="text-[10px] font-semibold uppercase tracking-widest text-[var(--text-muted)] mb-2 block">{{ t('task.description') }}</label>
          <textarea
            v-model="editDescription"
            :placeholder="t('task.noDescription')"
            rows="3"
            class="w-full bg-[var(--bg-input)] rounded-xl px-3.5 py-2.5 text-sm text-[var(--text-secondary)] leading-relaxed outline-none transition-all resize-none border border-[var(--border)] focus:border-primary/40"
          />
        </div>

        <!-- Labels -->
        <LabelManager :board-id="boardId" :task-id="task.id" :current-labels="task.labels ?? []" @updated="emit('updated')" />

        <!-- Subtasks -->
        <SubtaskList v-if="totalSubtasks > 0" :subtasks="task.subtasks ?? []" />

        <!-- Deadline (full width, at bottom of properties) -->
        <div>
          <label class="text-[10px] font-semibold uppercase tracking-widest text-[var(--text-muted)] mb-2 block">{{ t('task.deadline') }}</label>
          <div class="relative">
            <v-icon icon="mdi-calendar-outline" size="16" class="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--text-muted)]" />
            <input
              v-model="editDeadline"
              type="date"
              class="w-full pl-9 pr-3.5 py-2.5 rounded-xl text-sm outline-none border border-[var(--border)] bg-[var(--bg-input)] text-[var(--text)] focus:border-primary/40 transition-colors"
            />
          </div>
          <div
            v-if="deadlineInfo"
            class="mt-1.5 text-xs px-2 py-1 rounded-lg inline-flex items-center gap-1"
            :class="{
              'bg-error/10 text-error': deadlineInfo.isOverdue,
              'bg-warning/10 text-warning': deadlineInfo.isUrgent,
              'bg-[var(--bg-input)] text-[var(--text-muted)]': !deadlineInfo.isOverdue && !deadlineInfo.isUrgent
            }"
          >
            <v-icon :icon="deadlineInfo.isOverdue ? 'mdi-alert-circle' : 'mdi-clock-outline'" size="13" />
            {{ deadlineInfo.text }}
          </div>
        </div>
      </div>

      <!-- === SECTION: Tabs (Comments / Attachments / Activity) === -->
      <div class="border-t border-[var(--border)] pt-4">
        <div class="flex gap-4 mb-4">
          <button
            v-for="tab in (['comments', 'attachments', 'activity'] as const)"
            :key="tab"
            class="pb-1.5 text-xs font-medium border-b-2 transition-all"
            :class="activeTab === tab ? 'border-primary text-primary-light' : 'border-transparent text-[var(--text-muted)] hover:text-[var(--text-secondary)]'"
            @click="activeTab = tab"
          >
            {{ tab === 'comments' ? t('task.comments') : tab === 'attachments' ? t('task.labels') : t('task.activity') }}
          </button>
        </div>

        <CommentSection v-if="activeTab === 'comments'" :task-id="task.id" />
        <AttachmentSection v-else-if="activeTab === 'attachments'" :task-id="task.id" />
        <ActivityTimeline v-else :activities="activities" />
      </div>
    </div>
  </div>
</template>
