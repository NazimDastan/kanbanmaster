<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Priority } from '@/types/task'

const { t } = useI18n()

export interface FilterState {
  search: string
  priority: Priority | ''
  assigneeId: string
  labelId: string
  hasDeadline: boolean | null
}

const emit = defineEmits<{ filter: [filters: FilterState]; clear: [] }>()

const search = ref('')
const priority = ref<Priority | ''>('')
const assigneeId = ref('')
const labelId = ref('')
const hasDeadline = ref<boolean | null>(null)

const priorityOptions = [
  { value: '', title: t('filter.allPriorities') },
  { value: 'urgent', title: t('task.priorities.urgent') },
  { value: 'high', title: t('task.priorities.high') },
  { value: 'medium', title: t('task.priorities.medium') },
  { value: 'low', title: t('task.priorities.low') },
]

function applyFilters() {
  emit('filter', { search: search.value, priority: priority.value, assigneeId: assigneeId.value, labelId: labelId.value, hasDeadline: hasDeadline.value })
}

function clearFilters() {
  search.value = ''; priority.value = ''; assigneeId.value = ''; labelId.value = ''; hasDeadline.value = null; emit('clear')
}
</script>

<template>
  <v-card class="pa-4" elevation="2">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-sm font-semibold">{{ t('filter.title') }}</h3>
      <v-btn variant="text" size="x-small" color="primary" @click="clearFilters">{{ t('filter.clearAll') }}</v-btn>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
      <v-text-field v-model="search" :placeholder="t('filter.searchTasks')" prepend-inner-icon="mdi-magnify" density="compact" variant="outlined" rounded="lg" hide-details @input="applyFilters" />
      <v-select v-model="priority" :items="priorityOptions" item-value="value" item-title="title" density="compact" variant="outlined" rounded="lg" hide-details @update:model-value="applyFilters" />
      <v-text-field v-model="assigneeId" :placeholder="t('filter.assigneeId')" prepend-inner-icon="mdi-account-outline" density="compact" variant="outlined" rounded="lg" hide-details @input="applyFilters" />
      <v-btn-toggle v-model="hasDeadline" density="compact" color="primary" class="rounded-lg" @update:model-value="applyFilters">
        <v-btn :value="null" size="small">{{ t('filter.all') }}</v-btn>
        <v-btn :value="true" size="small">{{ t('filter.hasDeadline') }}</v-btn>
        <v-btn :value="false" size="small">{{ t('filter.noDeadline') }}</v-btn>
      </v-btn-toggle>
    </div>
  </v-card>
</template>
