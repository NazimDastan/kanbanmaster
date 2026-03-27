<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/useAuthStore'
import { useBoardStore } from '@/stores/useBoardStore'
import { taskService, type TaskFilter } from '@/services/taskService'
import { formatDeadline } from '@/utils/date'
import TaskPieChart from '@/components/dashboard/TaskPieChart.vue'
import { getInitials } from '@/utils/format'
import type { Task, Priority } from '@/types/task'
import { exportTasksToCSV } from '@/utils/export'
import { PRIORITY_CONFIG } from '@/types/task'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const boardStore = useBoardStore()

const allTasks = ref<Task[]>([])
const loading = ref(false)
const activeFilter = ref<TaskFilter>('')
const showTaskList = ref(false)

const totalCount = computed(() => allTasks.value.length)
const completedCount = computed(() => allTasks.value.filter(t => t.completedAt).length)
const overdueCount = computed(() => allTasks.value.filter(t => !t.completedAt && t.deadline && new Date(t.deadline) < new Date()).length)
const inProgressCount = computed(() => totalCount.value - completedCount.value - overdueCount.value)

const filteredTasks = ref<Task[]>([])

const stats = computed(() => [
  { key: 'all' as TaskFilter, title: t('dashboard.totalTasks'), value: totalCount.value, icon: 'mdi-clipboard-text-outline', gradient: 'linear-gradient(135deg, #6366f1, #818cf8)', filter: '' as TaskFilter },
  { key: 'completed' as TaskFilter, title: t('dashboard.completed'), value: completedCount.value, icon: 'mdi-check-circle-outline', gradient: 'linear-gradient(135deg, #10b981, #34d399)', filter: 'completed' as TaskFilter },
  { key: 'in_progress' as TaskFilter, title: t('dashboard.inProgress'), value: inProgressCount.value, icon: 'mdi-progress-clock', gradient: 'linear-gradient(135deg, #f59e0b, #fbbf24)', filter: 'in_progress' as TaskFilter },
  { key: 'overdue' as TaskFilter, title: t('dashboard.overdue'), value: overdueCount.value, icon: 'mdi-alert-circle-outline', gradient: 'linear-gradient(135deg, #ef4444, #f87171)', filter: 'overdue' as TaskFilter },
])

const filterTitle = computed(() => {
  const titles: Record<string, string> = {
    '': t('dashboard.allTasks'),
    completed: t('dashboard.completedTasks'),
    in_progress: t('dashboard.inProgress'),
    overdue: t('dashboard.overdueTasks'),
    assigned: t('dashboard.assignedToMe'),
  }
  return titles[activeFilter.value] ?? t('dashboard.tasks')
})

async function loadDashboard() {
  loading.value = true
  try {
    allTasks.value = await taskService.list()
    await boardStore.fetchBoards()
  } catch {
    // API not connected — empty state
  } finally {
    loading.value = false
  }
}

async function openFilter(filter: TaskFilter) {
  activeFilter.value = filter
  showTaskList.value = true
  try {
    filteredTasks.value = await taskService.list(filter || undefined)
  } catch {
    filteredTasks.value = []
  }
}

function navigateToBoard(boardId: string) {
  router.push(`/boards/${boardId}`)
}

function getPriorityConfig(priority: Priority) {
  return PRIORITY_CONFIG[priority]
}

onMounted(loadDashboard)
</script>

<template>
  <div class="p-4 md:p-6 lg:p-8 space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-xl md:text-2xl font-bold">{{ t('dashboard.hello', { name: authStore.userName || 'User' }) }}</h1>
      <p class="text-sm text-white/40 mt-0.5">{{ t('dashboard.subtitle') }}</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <v-progress-circular indeterminate color="primary" size="40" />
    </div>

    <template v-else>
      <!-- Onboarding: No boards yet -->
      <div v-if="boardStore.boards.length === 0 && allTasks.length === 0" class="rounded-2xl border border-white/5 bg-[#0f0f1a] p-8 md:p-12 text-center">
        <div class="w-16 h-16 mx-auto mb-5 rounded-2xl bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center">
          <v-icon icon="mdi-view-column-outline" size="32" class="text-primary-light" />
        </div>
        <h2 class="text-xl font-bold mb-2">{{ t('dashboard.welcomeTitle') }}</h2>
        <p class="text-sm text-white/40 mb-6 max-w-md mx-auto">
          {{ t('dashboard.welcomeDescription') }}
        </p>
        <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
          <div class="flex gap-2 items-center">
            <span class="w-7 h-7 rounded-full bg-primary/20 text-primary-light text-xs font-bold flex items-center justify-center">1</span>
            <span class="text-sm text-white/50">{{ t('dashboard.step1') }}</span>
          </div>
          <v-icon icon="mdi-arrow-right" size="16" class="text-white/15 hidden sm:block" />
          <div class="flex gap-2 items-center">
            <span class="w-7 h-7 rounded-full bg-primary/20 text-primary-light text-xs font-bold flex items-center justify-center">2</span>
            <span class="text-sm text-white/50">{{ t('dashboard.step2') }}</span>
          </div>
          <v-icon icon="mdi-arrow-right" size="16" class="text-white/15 hidden sm:block" />
          <div class="flex gap-2 items-center">
            <span class="w-7 h-7 rounded-full bg-primary/20 text-primary-light text-xs font-bold flex items-center justify-center">3</span>
            <span class="text-sm text-white/50">{{ t('dashboard.step3') }}</span>
          </div>
        </div>
        <p class="text-xs text-white/25 mt-6">{{ t('dashboard.sidebarHint') }}</p>
      </div>

      <!-- Stat cards — CLICKABLE -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
        <button
          v-for="s in stats"
          :key="s.key"
          class="rounded-xl border border-white/5 bg-[#0f0f1a] p-4 transition-all duration-200 hover:border-primary/30 hover:bg-[#12121f] text-left group"
          :class="{ '!border-primary/40 glow-primary': activeFilter === s.filter && showTaskList }"
          @click="openFilter(s.filter)"
        >
          <div class="flex items-start justify-between">
            <div>
              <p class="text-[10px] font-semibold uppercase tracking-widest text-white/30 mb-1">{{ s.title }}</p>
              <p class="text-2xl md:text-3xl font-bold">{{ s.value }}</p>
              <p class="text-[10px] text-white/20 mt-1 group-hover:text-primary-light transition-colors">{{ t('common.clickToView') }}</p>
            </div>
            <div class="w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0" :style="{ background: s.gradient }">
              <v-icon :icon="s.icon" color="white" size="18" />
            </div>
          </div>
        </button>
      </div>

      <!-- Task list panel (shows when stat clicked) -->
      <div v-if="showTaskList" class="rounded-2xl border border-white/5 bg-[#0f0f1a] overflow-hidden">
        <div class="flex items-center justify-between px-4 md:px-5 py-3 border-b border-white/5">
          <div class="flex items-center gap-3">
            <h2 class="text-sm font-semibold">{{ filterTitle }}</h2>
            <span class="px-2 py-0.5 rounded-md bg-white/5 text-[11px] font-medium text-white/40">{{ filteredTasks.length }}</span>
          </div>
          <button class="text-white/30 hover:text-white/60 transition-colors" @click="showTaskList = false">
            <v-icon icon="mdi-close" size="18" />
          </button>
        </div>

        <!-- Task rows -->
        <div v-if="filteredTasks.length > 0" class="divide-y divide-white/[0.03] max-h-[400px] overflow-y-auto">
          <div
            v-for="task in filteredTasks"
            :key="task.id"
            class="flex items-center gap-3 px-4 md:px-5 py-3 hover:bg-white/[0.02] transition-colors cursor-pointer"
          >
            <!-- Priority dot -->
            <div class="w-2 h-2 rounded-full flex-shrink-0" :style="{ backgroundColor: getPriorityConfig(task.priority).color }" />

            <!-- Task info -->
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium truncate">{{ task.title }}</p>
              <div class="flex items-center gap-3 mt-0.5">
                <span class="text-[10px] text-white/25">{{ task.priority }}</span>
                <span v-if="task.deadline" class="text-[10px]" :class="{ 'text-error': formatDeadline(task.deadline).isOverdue, 'text-warning': formatDeadline(task.deadline).isUrgent, 'text-white/25': !formatDeadline(task.deadline).isOverdue && !formatDeadline(task.deadline).isUrgent }">
                  {{ formatDeadline(task.deadline).text }}
                </span>
              </div>
            </div>

            <!-- Assignee -->
            <div v-if="task.assignee" class="w-6 h-6 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center flex-shrink-0">
              <span class="text-[8px] text-white font-semibold">{{ getInitials(task.assignee?.name ?? '') }}</span>
            </div>

            <!-- Completed badge -->
            <v-icon v-if="task.completedAt" icon="mdi-check-circle" size="16" color="#10b981" />
          </div>
        </div>

        <!-- Empty -->
        <div v-else class="py-10 text-center">
          <v-icon icon="mdi-checkbox-marked-circle-outline" size="36" class="text-white/10 mb-2" />
          <p class="text-sm text-white/30">{{ t('dashboard.noTasksInCategory') }}</p>
          <p class="text-[11px] text-white/15 mt-1">{{ t('dashboard.createTasksHint') }}</p>
        </div>
      </div>

      <!-- Chart + Boards -->
      <div v-if="totalCount > 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="rounded-2xl border border-white/5 bg-[#0f0f1a] p-5">
          <h2 class="text-sm font-semibold mb-4">{{ t('dashboard.totalTasks') }}</h2>
          <TaskPieChart :completed="completedCount" :in-progress="inProgressCount" :overdue="overdueCount" />
          <div class="flex justify-center gap-4 mt-4">
            <div class="flex items-center gap-1.5"><div class="w-2.5 h-2.5 rounded-full bg-success" /><span class="text-[11px] text-white/40">{{ t('dashboard.completed') }}</span></div>
            <div class="flex items-center gap-1.5"><div class="w-2.5 h-2.5 rounded-full bg-primary" /><span class="text-[11px] text-white/40">{{ t('dashboard.inProgress') }}</span></div>
            <div class="flex items-center gap-1.5"><div class="w-2.5 h-2.5 rounded-full bg-error" /><span class="text-[11px] text-white/40">{{ t('dashboard.overdue') }}</span></div>
          </div>
        </div>

        <!-- Completion rate large -->
        <div class="rounded-2xl border border-white/5 bg-[#0f0f1a] p-5 flex flex-col justify-center">
          <h2 class="text-sm font-semibold mb-6">{{ t('dashboard.completionRate') }}</h2>
          <div class="text-center mb-4">
            <p class="text-5xl font-bold" :class="completedCount / Math.max(totalCount, 1) >= 0.7 ? 'text-success' : completedCount / Math.max(totalCount, 1) >= 0.4 ? 'text-primary-light' : 'text-warning'">
              {{ totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0 }}%
            </p>
            <p class="text-xs text-white/25 mt-1">{{ completedCount }} / {{ totalCount }} {{ t('dashboard.totalTasks').toLowerCase() }}</p>
          </div>
          <v-progress-linear :model-value="totalCount > 0 ? (completedCount / totalCount) * 100 : 0" color="#10b981" height="6" rounded bg-color="surface" />
        </div>
      </div>

      <!-- Content: Boards + Quick Actions -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 md:gap-6">
        <!-- My Boards -->
        <div class="lg:col-span-2 rounded-2xl border border-white/5 bg-[#0f0f1a] overflow-hidden">
          <div class="px-4 md:px-5 py-3.5 border-b border-white/5">
            <h2 class="text-sm font-semibold">{{ t('dashboard.myBoards') }}</h2>
          </div>

          <div v-if="boardStore.boards.length > 0" class="grid grid-cols-1 sm:grid-cols-2 gap-3 p-4">
            <button
              v-for="board in boardStore.boards"
              :key="board.id"
              class="flex items-center gap-3 p-3 rounded-xl border border-white/5 bg-[#161625]/50 hover:border-primary/20 hover:bg-[#161625] transition-all text-left"
              @click="navigateToBoard(board.id)"
            >
              <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center flex-shrink-0">
                <v-icon icon="mdi-view-column-outline" size="20" class="text-primary-light" />
              </div>
              <div class="min-w-0">
                <p class="text-sm font-medium truncate">{{ board.name }}</p>
                <p class="text-[11px] text-white/25">{{ t('dashboard.openBoard') }}</p>
              </div>
              <v-icon icon="mdi-chevron-right" size="16" class="text-white/15 ml-auto" />
            </button>
          </div>

          <div v-else class="py-10 text-center">
            <v-icon icon="mdi-view-column-outline" size="36" class="text-white/10 mb-2" />
            <p class="text-sm text-white/30">{{ t('dashboard.noBoardsYet') }}</p>
            <p class="text-[11px] text-white/15 mt-1">{{ t('dashboard.createBoardHint') }}</p>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="rounded-2xl border border-white/5 bg-[#0f0f1a] p-4 md:p-5">
          <h2 class="text-sm font-semibold mb-4">{{ t('dashboard.quickActions') }}</h2>
          <div class="space-y-2">
            <v-btn block prepend-icon="mdi-view-column-outline" variant="tonal" color="primary" class="justify-start" style="text-transform: none" to="/teams">{{ t('dashboard.goToBoards') }}</v-btn>
            <v-btn block prepend-icon="mdi-account-group-outline" variant="tonal" color="secondary" class="justify-start" style="text-transform: none" to="/teams">{{ t('dashboard.manageTeams') }}</v-btn>
            <v-btn block prepend-icon="mdi-chart-bar" variant="tonal" class="justify-start" style="text-transform: none" to="/reports">{{ t('dashboard.viewReports') }}</v-btn>
            <v-btn v-if="allTasks.length > 0" block prepend-icon="mdi-download-outline" variant="tonal" class="justify-start" style="text-transform: none" @click="exportTasksToCSV(allTasks)">Export CSV</v-btn>
          </div>

          <div class="mt-5 pt-4 border-t border-white/5">
            <p class="text-[10px] font-semibold uppercase tracking-widest text-white/25 mb-3">{{ t('dashboard.completionRate') }}</p>
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-xs text-white/40">{{ t('dashboard.taskCount', { completed: completedCount, total: totalCount }) }}</span>
              <span class="text-xs font-semibold text-success">{{ totalCount > 0 ? Math.round((completedCount / totalCount) * 100) : 0 }}%</span>
            </div>
            <v-progress-linear :model-value="totalCount > 0 ? (completedCount / totalCount) * 100 : 0" color="#10b981" height="4" rounded bg-color="surface" />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
