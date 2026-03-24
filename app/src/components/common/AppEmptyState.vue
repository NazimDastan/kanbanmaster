<script setup lang="ts">
defineProps<{ icon: string; title: string; description: string; actionLabel?: string }>()
const emit = defineEmits<{ action: [] }>()
</script>

<template>
  <div class="flex flex-col items-center justify-center py-16 px-8 text-center">
    <v-icon :icon="icon" size="56" class="text-white/10 mb-4" />
    <h3 class="text-base font-semibold mb-1.5">{{ title }}</h3>
    <p class="text-sm text-white/30 mb-6 max-w-xs">{{ description }}</p>
    <v-btn v-if="actionLabel" color="primary" prepend-icon="mdi-plus" style="text-transform: none" @click="emit('action')">{{ actionLabel }}</v-btn>
  </div>
</template>
