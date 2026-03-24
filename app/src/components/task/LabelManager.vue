<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Label } from '@/types/task'
import { labelService } from '@/services/labelService'

const { t } = useI18n()
const props = defineProps<{ boardId: string; taskId: string; currentLabels: Label[] }>()
const emit = defineEmits<{ updated: [] }>()

const allLabels = ref<Label[]>([])
const showCreate = ref(false)
const newLabelName = ref('')
const newLabelColor = ref('#1E88E5')
const loading = ref(false)
const colorOptions = ['#1E88E5', '#7C4DFF', '#43A047', '#FB8C00', '#E53935', '#00ACC1', '#8E24AA', '#F4511E']

async function loadLabels() {
  try { allLabels.value = await labelService.listByBoard(props.boardId) } catch { /* */ }
}
onMounted(loadLabels)

function isAttached(labelId: string): boolean { return props.currentLabels.some((l) => l.id === labelId) }

async function toggleLabel(labelId: string) {
  loading.value = true
  try {
    if (isAttached(labelId)) await labelService.removeFromTask(props.taskId, labelId)
    else await labelService.addToTask(props.taskId, labelId)
    emit('updated')
  } finally { loading.value = false }
}

async function handleCreateLabel() {
  if (!newLabelName.value) return
  try {
    const label = await labelService.create(props.boardId, newLabelName.value, newLabelColor.value)
    allLabels.value.push(label)
    await labelService.addToTask(props.taskId, label.id)
    emit('updated'); newLabelName.value = ''; showCreate.value = false
  } catch { /* */ }
}
</script>

<template>
  <div>
    <p class="text-xs font-semibold uppercase tracking-wider text-text-secondary mb-2">{{ t('task.labels') }}</p>
    <div class="flex flex-wrap gap-1 mb-3">
      <v-chip v-for="label in currentLabels" :key="label.id" :color="label.color" size="small" variant="flat" label closable @click:close="toggleLabel(label.id)">{{ label.name }}</v-chip>
      <v-chip size="small" variant="outlined" label @click="showCreate = !showCreate"><v-icon icon="mdi-plus" size="14" /></v-chip>
    </div>
    <div v-if="showCreate" class="bg-background rounded-lg p-3 space-y-2">
      <div v-if="allLabels.length > 0" class="flex flex-wrap gap-1 mb-2">
        <v-chip v-for="label in allLabels" :key="label.id" :color="label.color" size="small" :variant="isAttached(label.id) ? 'flat' : 'outlined'" label class="cursor-pointer" @click="toggleLabel(label.id)">
          <v-icon v-if="isAttached(label.id)" icon="mdi-check" size="12" class="mr-1" />{{ label.name }}
        </v-chip>
      </div>
      <v-divider class="my-2" />
      <div class="flex gap-2 items-end">
        <v-text-field v-model="newLabelName" :label="t('label.newLabel')" density="compact" variant="outlined" rounded="lg" hide-details class="flex-1" />
        <div class="flex gap-1">
          <div v-for="color in colorOptions" :key="color" class="w-6 h-6 rounded-full cursor-pointer border-2 transition-all" :class="newLabelColor === color ? 'border-text-primary scale-110' : 'border-transparent'" :style="{ backgroundColor: color }" @click="newLabelColor = color" />
        </div>
        <v-btn icon="mdi-check" size="small" color="primary" :disabled="!newLabelName" @click="handleCreateLabel" />
      </div>
    </div>
  </div>
</template>
