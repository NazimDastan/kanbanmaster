<script setup lang="ts">
import { useConfirm } from '@/composables/useConfirm'

const { isOpen, title, message, confirmText, cancelText, isDanger, handleConfirm, handleCancel } = useConfirm()
</script>

<template>
  <v-dialog v-model="isOpen" max-width="400" persistent>
    <v-card class="pa-5" color="surface">
      <div class="flex items-start gap-3 mb-4">
        <div
          class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0"
          :class="isDanger ? 'bg-error/15' : 'bg-primary/15'"
        >
          <v-icon
            :icon="isDanger ? 'mdi-alert-circle-outline' : 'mdi-help-circle-outline'"
            :color="isDanger ? 'error' : 'primary'"
            size="22"
          />
        </div>
        <div>
          <h3 class="text-sm font-bold text-[var(--text)]">{{ title }}</h3>
          <p class="text-xs mt-1 text-[var(--text-secondary)]">{{ message }}</p>
        </div>
      </div>
      <div class="flex justify-end gap-2">
        <v-btn variant="text" size="small" @click="handleCancel">{{ cancelText }}</v-btn>
        <v-btn
          :color="isDanger ? 'error' : 'primary'"
          size="small"
          @click="handleConfirm"
        >
          {{ confirmText }}
        </v-btn>
      </div>
    </v-card>
  </v-dialog>
</template>
