<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Priority } from '@/types/task'

const { t } = useI18n()

const props = defineProps<{ columnId: string; loading?: boolean }>()
const emit = defineEmits<{ submit: [data: { columnId: string; title: string; description?: string; priority: string; deadline?: string }]; cancel: [] }>()

const title = ref('')
const description = ref('')
const priority = ref<Priority>('medium')
const deadline = ref('')

const priorities: { value: Priority; color: string; icon: string }[] = [
  { value: 'urgent', color: '#dc2626', icon: 'mdi-alert-circle' },
  { value: 'high', color: '#06b6d4', icon: 'mdi-arrow-up-bold' },
  { value: 'medium', color: '#6366f1', icon: 'mdi-minus' },
  { value: 'low', color: '#64748b', icon: 'mdi-arrow-down-bold' },
]

function handleSubmit() {
  if (!title.value.trim()) return
  emit('submit', {
    columnId: props.columnId,
    title: title.value.trim(),
    description: description.value || undefined,
    priority: priority.value,
    deadline: deadline.value || undefined,
  })
}
</script>

<template>
  <v-card color="surface" rounded="xl" class="overflow-hidden">
    <!-- Header -->
    <div class="flex items-center gap-3 px-6 pt-6 pb-4">
      <div class="w-10 h-10 rounded-xl bg-primary/10 flex items-center justify-center">
        <v-icon icon="mdi-plus-circle-outline" size="20" color="primary" />
      </div>
      <div>
        <h3 class="text-base font-semibold text-[var(--text)]">{{ t('task.newTask') }}</h3>
        <p class="text-xs text-[var(--text-muted)]">{{ t('task.createTask') }}</p>
      </div>
    </div>

    <v-divider class="opacity-10" />

    <!-- Form -->
    <form class="px-6 py-5 space-y-4" @submit.prevent="handleSubmit">
      <!-- Title -->
      <div>
        <label class="text-xs font-semibold text-[var(--text-muted)] uppercase tracking-wider mb-1.5 block">{{ t('task.title') }} *</label>
        <input
          v-model="title"
          :placeholder="t('task.title')"
          class="w-full px-3.5 py-2.5 rounded-xl text-sm outline-none border border-[var(--border)] bg-[var(--bg-input)] text-[var(--text)] placeholder:text-[var(--text-muted)] focus:border-primary/40 transition-colors"
          autofocus
        />
      </div>

      <!-- Description -->
      <div>
        <label class="text-xs font-semibold text-[var(--text-muted)] uppercase tracking-wider mb-1.5 block">{{ t('task.description') }}</label>
        <textarea
          v-model="description"
          :placeholder="t('task.noDescription')"
          rows="2"
          class="w-full px-3.5 py-2.5 rounded-xl text-sm outline-none border border-[var(--border)] bg-[var(--bg-input)] text-[var(--text)] placeholder:text-[var(--text-muted)] focus:border-primary/40 transition-colors resize-none"
        />
      </div>

      <!-- Priority -->
      <div>
        <label class="text-xs font-semibold text-[var(--text-muted)] uppercase tracking-wider mb-2 block">{{ t('task.priority') }}</label>
        <div class="flex gap-2">
          <button
            v-for="p in priorities"
            :key="p.value"
            type="button"
            class="flex-1 flex items-center justify-center gap-1.5 py-2 rounded-lg text-xs font-medium border transition-all"
            :class="priority === p.value ? 'shadow-sm' : 'border-transparent opacity-50 hover:opacity-80'"
            :style="priority === p.value
              ? { backgroundColor: p.color + '15', color: p.color, borderColor: p.color + '40' }
              : {}"
            @click="priority = p.value"
          >
            <v-icon :icon="p.icon" size="14" />
            {{ t(`task.priorities.${p.value}`) }}
          </button>
        </div>
      </div>

      <!-- Deadline (full width, at bottom) -->
      <div>
        <label class="text-xs font-semibold text-[var(--text-muted)] uppercase tracking-wider mb-1.5 block">{{ t('task.deadline') }}</label>
        <div class="relative">
          <v-icon icon="mdi-calendar-outline" size="16" class="absolute left-3 top-1/2 -translate-y-1/2 text-[var(--text-muted)]" />
          <input
            v-model="deadline"
            type="date"
            class="w-full pl-9 pr-3.5 py-2.5 rounded-xl text-sm outline-none border border-[var(--border)] bg-[var(--bg-input)] text-[var(--text)] focus:border-primary/40 transition-colors"
          />
        </div>
      </div>
    </form>

    <v-divider class="opacity-10" />

    <!-- Actions -->
    <div class="flex justify-end gap-2 px-6 py-4">
      <v-btn variant="text" size="small" class="normal-case" @click="emit('cancel')">{{ t('common.cancel') }}</v-btn>
      <v-btn
        color="primary"
        size="small"
        :loading="loading"
        :disabled="!title.trim()"
        class="normal-case"
        prepend-icon="mdi-plus"
        @click="handleSubmit"
      >
        {{ t('task.createTask') }}
      </v-btn>
    </div>
  </v-card>
</template>
