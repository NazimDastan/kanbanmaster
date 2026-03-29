<script setup lang="ts">
import { onMounted, ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useBoardStore } from '@/stores/useBoardStore'
import { useTaskStore } from '@/stores/useTaskStore'
import { columnService } from '@/services/columnService'
import { useToast } from '@/composables/useToast'
import BoardColumn from '@/components/board/BoardColumn.vue'
import TaskDetail from '@/components/task/TaskDetail.vue'
import TaskForm from '@/components/task/TaskForm.vue'
import AppEmptyState from '@/components/common/AppEmptyState.vue'
import DelegateModal from '@/components/task/DelegateModal.vue'
import CalendarView from '@/components/board/CalendarView.vue'
import ListView from '@/components/board/ListView.vue'
import { delegationService } from '@/services/reportService'
import { useKeyboard } from '@/composables/useKeyboard'
import draggable from 'vuedraggable'
import { useConfirm } from '@/composables/useConfirm'
import type { Task } from '@/types/task'
import type { Column } from '@/types/board'

const COLUMN_COLORS = [
  '#f59e0b', '#ef4444', '#6366f1', '#10b981',
  '#ec4899', '#8b5cf6', '#06b6d4', '#64748b',
]

const { t } = useI18n()
const toast = useToast()
const { confirm } = useConfirm()
const route = useRoute()
const boardStore = useBoardStore()
const taskStore = useTaskStore()

const showTaskDrawer = ref(false)
const showTaskFormModal = ref(false)
const activeColumnId = ref('')
const showAddColumnModal = ref(false)
const showDelegateModal = ref(false)
const showMoveTasksModal = ref(false)
const newColumnName = ref('')
const newColumnColor = ref('')
const taskFormLoading = ref(false)
const delegateLoading = ref(false)
const searchQuery = ref('')
const filterPriority = ref('')
const viewMode = ref<'kanban' | 'list' | 'calendar'>('kanban')
const deleteTargetColumnId = ref('')
const moveTargetColumnId = ref('')

// All tasks flat list (for list/calendar views)
const allTasks = computed(() => {
  if (!boardStore.currentBoard?.columns) return []
  return boardStore.currentBoard.columns.flatMap((col) => col.tasks)
})

// Filtered columns — search + priority filter applied
// Used by draggable when search/filter is active
const filteredColumns = computed(() => {
  if (!boardStore.currentBoard?.columns) return []
  const q = searchQuery.value.toLowerCase().trim()
  const p = filterPriority.value

  if (!q && !p) return boardStore.currentBoard.columns

  return boardStore.currentBoard.columns.map((col) => ({
    ...col,
    tasks: col.tasks.filter((task) => {
      const matchesSearch = !q || task.title.toLowerCase().includes(q) || (task.description?.toLowerCase().includes(q) ?? false)
      const matchesPriority = !p || task.priority === p
      return matchesSearch && matchesPriority
    }),
  }))
})

// Column being deleted (for the move tasks modal)
const deleteTargetColumn = computed((): Column | undefined => {
  return boardStore.currentBoard?.columns.find((c) => c.id === deleteTargetColumnId.value)
})

// Available columns to move tasks to (excludes the column being deleted)
const moveTargetColumns = computed((): Column[] => {
  if (!boardStore.currentBoard?.columns) return []
  return boardStore.currentBoard.columns.filter((c) => c.id !== deleteTargetColumnId.value)
})

// Keyboard shortcuts
useKeyboard([
  { key: 'n', handler: () => { activeColumnId.value = boardStore.currentBoard?.columns?.[0]?.id ?? ''; showTaskFormModal.value = true } },
  { key: 'f', handler: () => { const input = document.querySelector<HTMLInputElement>('input[placeholder]'); input?.focus() } },
  { key: '1', handler: () => { viewMode.value = 'kanban' } },
  { key: '2', handler: () => { viewMode.value = 'list' } },
  { key: '3', handler: () => { viewMode.value = 'calendar' } },
  { key: 'Escape', handler: () => { showTaskDrawer.value = false; showTaskFormModal.value = false } },
])

async function loadBoard() {
  await boardStore.fetchBoard(route.params.id as string)
}
onMounted(() => {
  if (window.innerWidth < 768) viewMode.value = 'list'
  loadBoard()
})
watch(() => route.params.id, loadBoard)

function handleTaskClick(task: Task) {
  taskStore.selectedTask = task
  showTaskDrawer.value = true
}

function handleAddTask(columnId: string) {
  activeColumnId.value = columnId
  showTaskFormModal.value = true
}

async function handleQuickAdd(data: { columnId: string; title: string }) {
  try {
    await taskStore.createTask({ columnId: data.columnId, title: data.title, priority: 'medium' })
    toast.success(t('task.createTask') + ' ✓')
    await boardStore.fetchBoard(route.params.id as string)
  } catch {
    toast.error(t('common.error'))
  }
}

async function handleTaskMoved(event: { taskId: string; toColumnId: string; newIndex: number }) {
  try {
    await taskStore.moveTask(event.taskId, event.toColumnId, event.newIndex)
  } catch {
    toast.error(t('common.error'))
    await boardStore.fetchBoard(route.params.id as string)
  }
}

async function handleTaskSubmit(data: { columnId: string; title: string; description?: string; priority: string; deadline?: string }) {
  taskFormLoading.value = true
  try {
    await taskStore.createTask(data)
    toast.success(t('task.createTask') + ' ✓')
    showTaskFormModal.value = false
    await boardStore.fetchBoard(route.params.id as string)
  } finally {
    taskFormLoading.value = false
  }
}

async function handleDeleteTask(taskId: string) {
  const ok = await confirm({ title: t('task.deleteTask'), message: t('common.confirm') + '?', confirmText: t('common.delete'), danger: true })
  if (!ok) return
  await taskStore.deleteTask(taskId)
  toast.success(t('task.deleteTask') + ' ✓')
  showTaskDrawer.value = false
  await boardStore.fetchBoard(route.params.id as string)
}

function handleDelegate() {
  showTaskDrawer.value = false
  showDelegateModal.value = true
}

async function handleDelegateSubmit(data: { toUserId: string; reason: string }) {
  if (!taskStore.selectedTask) return
  delegateLoading.value = true
  try {
    await delegationService.delegate(taskStore.selectedTask.id, data.toUserId, data.reason)
    toast.success(t('task.delegateTask') + ' ✓')
    showDelegateModal.value = false
    await boardStore.fetchBoard(route.params.id as string)
  } finally {
    delegateLoading.value = false
  }
}

async function handleColumnReorder() {
  if (!boardStore.currentBoard) return
  const items = boardStore.currentBoard.columns.map((col, i) => ({ columnId: col.id, position: i }))
  try {
    await columnService.reorder(boardStore.currentBoard.id, items)
  } catch {
    toast.error(t('common.error'))
    await boardStore.fetchBoard(boardStore.currentBoard.id)
  }
}

async function handleAddColumn() {
  if (!newColumnName.value || !boardStore.currentBoard) return
  await columnService.create(boardStore.currentBoard.id, newColumnName.value, newColumnColor.value)
  toast.success(t('board.addColumn') + ' ✓')
  showAddColumnModal.value = false
  newColumnName.value = ''
  newColumnColor.value = ''
  await boardStore.fetchBoard(boardStore.currentBoard.id)
}

function handleDeleteColumn(columnId: string) {
  const columns = boardStore.currentBoard?.columns
  if (!columns) return

  const colIndex = columns.findIndex((c) => c.id === columnId)
  if (colIndex < 3) {
    toast.error(t('board.cannotDeleteDefault'))
    return
  }

  const column = columns[colIndex]
  if (column.tasks.length > 0) {
    // Show warning modal to move tasks first
    deleteTargetColumnId.value = columnId
    moveTargetColumnId.value = columns[0]?.id ?? ''
    showMoveTasksModal.value = true
    return
  }

  // Column is empty — confirm and delete
  deleteColumnDirectly(columnId)
}

async function deleteColumnDirectly(columnId: string) {
  const ok = await confirm({
    title: t('board.deleteColumn'),
    message: t('common.confirm') + '?',
    confirmText: t('common.delete'),
    danger: true,
  })
  if (!ok) return
  await columnService.delete(columnId)
  toast.success(t('board.deleteColumn') + ' ✓')
  await boardStore.fetchBoard(route.params.id as string)
}

async function handleMoveAndDelete() {
  if (!deleteTargetColumn.value || !moveTargetColumnId.value) return

  // Move all tasks to the target column
  const tasks = deleteTargetColumn.value.tasks
  for (const task of tasks) {
    await taskStore.moveTask(task.id, moveTargetColumnId.value, 0)
  }

  // Now delete the empty column
  await columnService.delete(deleteTargetColumnId.value)
  toast.success(t('board.deleteColumn') + ' ✓')
  showMoveTasksModal.value = false
  deleteTargetColumnId.value = ''
  moveTargetColumnId.value = ''
  await boardStore.fetchBoard(route.params.id as string)
}
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Board header -->
    <div class="flex items-center justify-between px-4 md:px-6 py-3 border-b border-white/5 flex-shrink-0">
      <h1 class="text-base font-bold">{{ boardStore.currentBoard?.name ?? 'Board' }}</h1>
      <div class="flex items-center gap-2">
        <!-- Search -->
        <div class="relative hidden sm:block">
          <v-icon icon="mdi-magnify" size="16" class="absolute left-2.5 top-1/2 -translate-y-1/2 text-white/25" />
          <input
            v-model="searchQuery"
            :placeholder="t('board.searchTasks')"
            class="w-44 pl-8 pr-3 py-1.5 rounded-lg bg-white/5 border border-white/5 text-xs text-white placeholder:text-white/20 outline-none focus:border-primary/30 transition-colors"
          />
        </div>

        <!-- Priority filter -->
        <select
          v-model="filterPriority"
          class="hidden sm:block px-2.5 py-1.5 rounded-lg bg-white/5 border border-white/5 text-xs text-white/60 outline-none focus:border-primary/30 transition-colors"
        >
          <option value="">{{ t('filter.allPriorities') }}</option>
          <option value="urgent">{{ t('task.priorities.urgent') }}</option>
          <option value="high">{{ t('task.priorities.high') }}</option>
          <option value="medium">{{ t('task.priorities.medium') }}</option>
          <option value="low">{{ t('task.priorities.low') }}</option>
        </select>

        <!-- View mode toggle -->
        <div class="flex bg-white/5 rounded-lg p-0.5">
          <button
            v-for="mode in (['kanban', 'list', 'calendar'] as const)"
            :key="mode"
            class="px-2 py-1 rounded-md transition-all"
            :class="viewMode === mode ? 'bg-primary/20 text-primary-light' : 'text-white/30 hover:text-white/50'"
            @click="viewMode = mode"
          >
            <v-icon :icon="mode === 'kanban' ? 'mdi-view-column-outline' : mode === 'list' ? 'mdi-format-list-bulleted' : 'mdi-calendar-month-outline'" size="16" />
          </button>
        </div>

        <v-btn color="primary" prepend-icon="mdi-plus" size="small" class="normal-case" :disabled="!boardStore.currentBoard?.columns?.length" @click="activeColumnId = boardStore.currentBoard?.columns?.[0]?.id ?? ''; showTaskFormModal = true">
          {{ t('board.newTask') }}
        </v-btn>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="boardStore.loading" class="flex-1 flex items-center justify-center">
      <v-progress-circular indeterminate color="primary" size="40" />
    </div>

    <!-- Empty -->
    <AppEmptyState v-else-if="!boardStore.currentBoard?.columns?.length" icon="mdi-view-column-outline" :title="t('board.noColumnsYet')" :description="t('board.noColumnsDescription')" :action-label="t('board.addColumn')" @action="showAddColumnModal = true" />

    <!-- Kanban View -->
    <div v-else-if="viewMode === 'kanban'" class="flex-1 flex p-4 md:p-5 overflow-x-auto min-h-0">
      <draggable
        :list="filteredColumns"
        item-key="id"
        group="kanban-columns"
        direction="horizontal"
        ghost-class="column-ghost"
        handle=".column-drag-handle"
        :animation="200"
        class="flex gap-3"
        @end="handleColumnReorder"
      >
        <template #item="{ element: col, index }">
          <div class="column-drag-handle cursor-grab active:cursor-grabbing">
            <BoardColumn
              :column="col"
              :column-index="index"
              :can-delete="index >= 3"
              @task-click="handleTaskClick"
              @add-task="handleAddTask"
              @quick-add="handleQuickAdd"
              @task-moved="handleTaskMoved"
              @delete-column="handleDeleteColumn"
            />
          </div>
        </template>
      </draggable>

      <button class="flex flex-col items-center justify-center min-w-[272px] max-w-[272px] h-28 border border-dashed border-white/[0.06] rounded-xl hover:border-primary/30 hover:bg-primary/5 transition-all ml-3 flex-shrink-0" @click="showAddColumnModal = true">
        <v-icon icon="mdi-plus" size="20" class="text-white/20" />
        <span class="text-[11px] text-white/20 mt-1">{{ t('board.addColumn') }}</span>
      </button>
    </div>

    <!-- List View -->
    <div v-else-if="viewMode === 'list'" class="flex-1 overflow-y-auto p-4 md:p-5 min-h-0">
      <ListView :tasks="allTasks" @task-click="handleTaskClick" />
    </div>

    <!-- Calendar View -->
    <div v-else-if="viewMode === 'calendar'" class="flex-1 overflow-y-auto p-4 md:p-5 min-h-0">
      <CalendarView :tasks="allTasks" @task-click="handleTaskClick" />
    </div>

    <!-- Task drawer -->
    <v-navigation-drawer v-model="showTaskDrawer" location="right" temporary width="440" color="surface">
      <TaskDetail v-if="taskStore.selectedTask" :task="taskStore.selectedTask" :board-id="boardStore.currentBoard?.id ?? ''" @close="showTaskDrawer = false" @delete="handleDeleteTask" @delegate="handleDelegate" @updated="boardStore.fetchBoard(route.params.id as string)" />
    </v-navigation-drawer>

    <!-- Modals -->
    <v-dialog v-model="showDelegateModal" max-width="440">
      <DelegateModal :task-title="taskStore.selectedTask?.title ?? ''" :loading="delegateLoading" @submit="handleDelegateSubmit" @cancel="showDelegateModal = false" />
    </v-dialog>

    <v-dialog v-model="showTaskFormModal" max-width="440">
      <TaskForm :column-id="activeColumnId" :loading="taskFormLoading" @submit="handleTaskSubmit" @cancel="showTaskFormModal = false" />
    </v-dialog>

    <!-- Add Column Modal (with color picker) -->
    <v-dialog v-model="showAddColumnModal" max-width="360">
      <v-card class="pa-5" color="surface">
        <h3 class="text-sm font-bold mb-3">{{ t('board.addColumn') }}</h3>
        <v-text-field v-model="newColumnName" :label="t('board.columnName')" prepend-inner-icon="mdi-view-column-outline" autofocus @keyup.enter="handleAddColumn" />
        <div class="mt-1 mb-3">
          <p class="text-xs text-white/50 mb-2">{{ t('board.columnColor') }}</p>
          <div class="flex gap-2 flex-wrap">
            <button
              v-for="color in COLUMN_COLORS"
              :key="color"
              class="w-7 h-7 rounded-full border-2 transition-all hover:scale-110"
              :class="newColumnColor === color ? 'border-white scale-110' : 'border-transparent'"
              :style="{ backgroundColor: color }"
              @click="newColumnColor = newColumnColor === color ? '' : color"
            />
          </div>
        </div>
        <div class="flex justify-end gap-2 mt-3">
          <v-btn variant="text" size="small" @click="showAddColumnModal = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="primary" size="small" :disabled="!newColumnName" class="normal-case" @click="handleAddColumn">{{ t('common.add') }}</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Move Tasks Warning Modal -->
    <v-dialog v-model="showMoveTasksModal" max-width="420">
      <v-card class="pa-5" color="surface">
        <div class="flex items-center gap-2 mb-3">
          <v-icon icon="mdi-alert-outline" color="warning" size="20" />
          <h3 class="text-sm font-bold">{{ t('board.deleteColumn') }}</h3>
        </div>
        <p class="text-xs text-white/60 mb-4">{{ t('board.moveTasksFirst') }}</p>
        <div class="mb-4">
          <label class="text-xs text-white/50 mb-1 block">{{ t('board.moveTo') }}</label>
          <select
            v-model="moveTargetColumnId"
            class="w-full px-3 py-2 rounded-lg bg-white/5 border border-white/10 text-sm text-white outline-none focus:border-primary/30 transition-colors"
          >
            <option v-for="col in moveTargetColumns" :key="col.id" :value="col.id">
              {{ col.name }} ({{ col.tasks.length }})
            </option>
          </select>
        </div>
        <div class="flex justify-end gap-2">
          <v-btn variant="text" size="small" @click="showMoveTasksModal = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="error" size="small" class="normal-case" :disabled="!moveTargetColumnId" @click="handleMoveAndDelete">
            {{ t('board.deleteColumn') }}
          </v-btn>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>
