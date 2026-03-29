<script setup lang="ts">
import { useToast } from '@/composables/useToast'

const { toasts } = useToast()

const iconMap: Record<string, string> = {
  success: 'mdi-check-circle',
  error: 'mdi-alert-circle',
  info: 'mdi-information',
  warning: 'mdi-alert',
}
const colorMap: Record<string, string> = {
  success: '#10b981',
  error: '#ef4444',
  info: '#6366f1',
  warning: '#f59e0b',
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-[9999] flex flex-col gap-2 pointer-events-none max-w-[360px]">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="pointer-events-auto flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg backdrop-blur-md"
          :style="{
            backgroundColor: colorMap[toast.type] + '18',
            borderColor: colorMap[toast.type] + '40',
            border: '1px solid ' + colorMap[toast.type] + '40',
            color: 'var(--text)',
          }"
        >
          <v-icon :icon="iconMap[toast.type]" size="18" :style="{ color: colorMap[toast.type] }" />
          <p class="text-sm font-medium flex-1 text-[var(--text)]">{{ toast.message }}</p>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active { transition: all 0.3s ease-out; }
.toast-leave-active { transition: all 0.2s ease-in; }
.toast-enter-from { opacity: 0; transform: translateX(40px); }
.toast-leave-to { opacity: 0; transform: translateX(40px); }
</style>
