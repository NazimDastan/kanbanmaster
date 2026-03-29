<script setup lang="ts">
import { formatDate } from '@/utils/date'

defineProps<{ createdAt: string; updatedAt: string }>()
</script>

<template>
  <div class="px-5 py-4">
    <div class="flex items-center gap-4 text-[11px] text-white/20">
      <div class="flex items-center gap-1.5">
        <v-icon icon="mdi-calendar-plus" size="13" />
        {{ formatDate(createdAt) }}
      </div>
      <div class="flex items-center gap-1.5">
        <v-icon icon="mdi-calendar-edit" size="13" />
        {{ formatDate(updatedAt) }}
      </div>
    </div>
  </div>
</template>
