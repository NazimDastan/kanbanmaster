<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Subtask } from '@/types/task'

const { t } = useI18n()
const props = defineProps<{ subtasks: Subtask[] }>()

const completedCount = computed(() => props.subtasks.filter((s) => s.isCompleted).length)
</script>

<template>
  <div class="px-5 py-4 section-border">
    <div class="flex items-center justify-between mb-3">
      <p class="text-[10px] font-semibold uppercase tracking-widest text-white/25">{{ t('task.subtasks') }}</p>
      <span class="text-[11px] text-white/30">{{ completedCount }}/{{ subtasks.length }}</span>
    </div>
    <v-progress-linear :model-value="(completedCount / subtasks.length) * 100" color="primary" height="4" rounded bg-color="surface" class="mb-3" />
    <div class="space-y-0.5">
      <div v-for="sub in subtasks" :key="sub.id" class="flex items-center gap-2.5 py-1.5 px-2 rounded-lg hover:bg-white/[0.02] transition-colors">
        <v-checkbox-btn :model-value="sub.isCompleted" density="compact" color="primary" />
        <span class="text-sm" :class="sub.isCompleted ? 'line-through text-white/25' : 'text-white/70'">{{ sub.title }}</span>
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
</style>
