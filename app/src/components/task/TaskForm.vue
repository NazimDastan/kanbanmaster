<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { required } from '@/utils/validation'
import type { Priority } from '@/types/task'

const { t } = useI18n()

const props = defineProps<{ columnId: string; loading?: boolean }>()
const emit = defineEmits<{ submit: [data: { columnId: string; title: string; description?: string; priority: string; deadline?: string }]; cancel: [] }>()

const form = ref(false)
const title = ref('')
const description = ref('')
const priority = ref<Priority>('medium')
const deadline = ref('')

const priorityOptions = [
  { value: 'urgent', title: t('task.priorities.urgent'), color: '#D50000' },
  { value: 'high', title: t('task.priorities.high'), color: '#E53935' },
  { value: 'medium', title: t('task.priorities.medium'), color: '#1E88E5' },
  { value: 'low', title: t('task.priorities.low'), color: '#9E9E9E' },
]

function handleSubmit() {
  if (!form.value) return
  emit('submit', { columnId: props.columnId, title: title.value, description: description.value || undefined, priority: priority.value, deadline: deadline.value || undefined })
}
</script>

<template>
  <v-card class="pa-6" max-width="480">
    <h3 class="text-lg font-bold mb-4">{{ t('task.newTask') }}</h3>
    <v-form v-model="form" @submit.prevent="handleSubmit">
      <v-text-field v-model="title" :label="t('task.title')" :rules="[required(t('task.title'))]" prepend-inner-icon="mdi-text" class="mb-2" />
      <v-textarea v-model="description" :label="t('task.description')" rows="3" auto-grow variant="outlined" rounded="lg" class="mb-2" />
      <v-select v-model="priority" :items="priorityOptions" item-value="value" item-title="title" :label="t('task.priority')" class="mb-2">
        <template #item="{ item, props: itemProps }">
          <v-list-item v-bind="itemProps">
            <template #prepend>
              <div class="w-3 h-3 rounded-full mr-2" :style="{ backgroundColor: item.raw.color }" />
            </template>
          </v-list-item>
        </template>
      </v-select>
      <v-text-field v-model="deadline" :label="t('task.deadline')" type="date" prepend-inner-icon="mdi-calendar" class="mb-4" />
      <div class="flex justify-end gap-2">
        <v-btn variant="text" @click="emit('cancel')">{{ t('common.cancel') }}</v-btn>
        <v-btn type="submit" color="primary" :disabled="!form" :loading="loading">{{ t('task.createTask') }}</v-btn>
      </div>
    </v-form>
  </v-card>
</template>
