<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
defineProps<{ saving: boolean }>()
defineEmits<{ save: []; reset: [] }>()
</script>

<template>
  <div class="px-5 py-3 flex-shrink-0 section-border-top bg-[var(--bg-card)]">
    <div class="flex items-center justify-between">
      <p class="text-xs text-warning flex items-center gap-1.5">
        <v-icon icon="mdi-circle-medium" size="16" />
        {{ t('common.save') }}
      </p>
      <div class="flex gap-2">
        <button
          class="px-3 py-1.5 rounded-lg text-xs text-[var(--text-muted)] hover:text-[var(--text-secondary)] hover:bg-[var(--bg-hover)] transition-all"
          @click="$emit('reset')"
        >
          {{ t('common.cancel') }}
        </button>
        <button
          class="px-4 py-1.5 rounded-lg text-xs font-medium text-white transition-all"
          :class="saving ? 'bg-primary/30' : 'bg-primary hover:bg-primary-dark'"
          :disabled="saving"
          @click="$emit('save')"
        >
          {{ saving ? t('common.loading') : t('common.save') }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.section-border-top {
  border-top: 1px solid var(--p-border, rgba(255,255,255,0.04));
  transition: border-color 0.3s;
}
</style>
