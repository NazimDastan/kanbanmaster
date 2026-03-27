<script setup lang="ts">
defineProps<{
  title: string
  value: number | string
  icon: string
  gradient: string
  change?: number
}>()
</script>

<template>
  <div class="rounded-xl border border-white/5 bg-card p-4 transition-all duration-200 hover:border-white/10 hover:bg-card-hover">
    <div class="flex items-start justify-between">
      <div>
        <p class="text-[10px] font-semibold uppercase tracking-widest text-white/30 mb-1">{{ title }}</p>
        <p class="text-2xl md:text-3xl font-bold">{{ value }}</p>
        <div v-if="change !== undefined" class="flex items-center gap-1 mt-1.5">
          <v-icon :icon="change >= 0 ? 'mdi-trending-up' : 'mdi-trending-down'" :color="change >= 0 ? '#10b981' : '#ef4444'" size="13" />
          <span class="text-[11px] font-medium" :class="change >= 0 ? 'text-success' : 'text-error'">{{ Math.abs(change) }}%</span>
        </div>
      </div>
      <div class="w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0" :style="{ background: gradient }">
        <v-icon :icon="icon" color="white" size="18" />
      </div>
    </div>
  </div>
</template>
