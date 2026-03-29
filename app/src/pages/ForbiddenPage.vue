<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center">
      <h1 class="text-6xl font-bold text-warning mb-4">403</h1>
      <h2 class="text-xl font-semibold text-text-primary mb-2">{{ t('errors.forbidden') }}</h2>
      <p class="text-text-secondary mb-6">{{ t('errors.forbiddenDesc') }}</p>
      <div class="flex gap-3 justify-center">
        <v-btn color="primary" prepend-icon="mdi-arrow-left" @click="$router.back()">{{ t('errors.goBack') }}</v-btn>
        <v-btn variant="outlined" to="/" prepend-icon="mdi-home">{{ t('notFound.goHome') }}</v-btn>
      </div>
    </div>
  </div>
</template>
