<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import draggable from 'vuedraggable'
import type { Column } from '@/types/board'
import type { Task } from '@/types/task'
import TaskCard from './TaskCard.vue'

const { t } = useI18n()
const props = defineProps<{ column: Column }>()

const emit = defineEmits<{
  taskClick: [task: Task]
  addTask: [columnId: string]
  quickAdd: [data: { columnId: string; title: string }]
  taskMoved: [event: { taskId: string; toColumnId: string; newIndex: number }]
}>()

const showQuickAdd = ref(false)
const quickTitle = ref('')

function handleQuickSubmit() {
  if (!quickTitle.value.trim()) return
  emit('quickAdd', { columnId: props.column.id, title: quickTitle.value.trim() })
  quickTitle.value = ''
  showQuickAdd.value = false
}

function handleDragEnd(event: { item?: { dataset?: { taskId?: string } }; to?: { dataset?: { columnId?: string } }; newIndex?: number }) {
  const taskId = event.item?.dataset?.taskId
  const toColumnId = event.to?.dataset?.columnId ?? props.column.id
  const newIndex = event.newIndex ?? 0
  if (taskId) emit('taskMoved', { taskId, toColumnId, newIndex })
}
</script>

<template>
  <div class="flex flex-col min-w-[272px] max-w-[272px] h-full rounded-xl bg-card-60 border border-white/[0.04]">
    <!-- Header -->
    <div class="flex items-center justify-between px-3 py-2.5">
      <div class="flex items-center gap-2">
        <span class="text-xs font-semibold text-white/70">{{ column.name }}</span>
        <span class="px-1.5 py-0.5 rounded-md bg-white/5 text-[10px] font-medium text-white/30">{{ column.tasks.length }}</span>
      </div>
      <button class="text-white/20 hover:text-primary transition-colors" @click="showQuickAdd = true">
        <v-icon icon="mdi-plus" size="16" />
      </button>
    </div>

    <!-- Quick add inline -->
    <div v-if="showQuickAdd" class="px-2 pb-2">
      <div class="rounded-lg border border-primary/30 bg-elevated p-2.5">
        <input
          v-model="quickTitle"
          :placeholder="t('task.title')"
          class="w-full bg-transparent text-sm text-white placeholder:text-white/20 outline-none mb-2"
          autofocus
          @keyup.enter="handleQuickSubmit"
          @keyup.escape="showQuickAdd = false"
        />
        <div class="flex items-center justify-between">
          <p class="text-[10px] text-white/20">Enter ↵</p>
          <div class="flex gap-1.5">
            <button class="text-[11px] text-white/30 hover:text-white/60 transition-colors" @click="showQuickAdd = false">{{ t('common.cancel') }}</button>
            <button class="text-[11px] text-primary-light font-medium disabled:opacity-30" :disabled="!quickTitle.trim()" @click="handleQuickSubmit">{{ t('common.add') }}</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Draggable tasks -->
    <draggable
      :list="column.tasks"
      group="kanban-tasks"
      item-key="id"
      :data-column-id="column.id"
      class="flex-1 overflow-y-auto px-2 pb-2 space-y-1.5 min-h-[60px]"
      ghost-class="drag-ghost"
      drag-class="drag-active"
      :animation="200"
      @end="handleDragEnd"
    >
      <template #item="{ element }">
        <div :data-task-id="element.id">
          <TaskCard :task="element" @click="emit('taskClick', element)" />
        </div>
      </template>

      <template #footer>
        <div v-if="column.tasks.length === 0 && !showQuickAdd" class="flex justify-center py-8">
          <button class="flex items-center gap-1.5 px-3 py-2 rounded-lg border border-dashed border-white/[0.08] text-white/20 hover:text-primary-light hover:border-primary/30 transition-all text-[11px]" @click="showQuickAdd = true">
            <v-icon icon="mdi-plus" size="14" />
            {{ t('board.addTask') }}
          </button>
        </div>
      </template>
    </draggable>

    <!-- Bottom add button -->
    <div v-if="column.tasks.length > 0 && !showQuickAdd" class="px-2 pb-2">
      <button class="w-full flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-white/20 hover:text-white/40 hover:bg-white/[0.03] transition-all text-[11px]" @click="showQuickAdd = true">
        <v-icon icon="mdi-plus" size="14" />
        {{ t('board.addTask') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.drag-ghost { opacity: 0.3; border: 1px dashed rgba(99, 102, 241, 0.4) !important; border-radius: 8px; }
.drag-active { opacity: 0.8; transform: rotate(2deg); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4); }
</style>
