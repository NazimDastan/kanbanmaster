<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center">
      <h1 class="text-6xl font-bold text-primary mb-4">404</h1>
      <h2 class="text-xl font-semibold text-text-primary mb-2">{{ t('notFound.title') }}</h2>
      <p class="text-text-secondary mb-6">{{ t('notFound.description') }}</p>
      <v-btn color="primary" to="/" prepend-icon="mdi-home">{{ t('notFound.goHome') }}</v-btn>
    </div>
  </div>
</template>
