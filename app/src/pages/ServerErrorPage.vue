<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center">
      <h1 class="text-6xl font-bold text-error mb-4">500</h1>
      <h2 class="text-xl font-semibold text-text-primary mb-2">{{ t('errors.serverError') }}</h2>
      <p class="text-text-secondary mb-6">{{ t('errors.serverErrorDesc') }}</p>
      <div class="flex gap-3 justify-center">
        <v-btn color="primary" prepend-icon="mdi-refresh" @click="$router.go(0)">{{ t('errors.tryAgain') }}</v-btn>
        <v-btn variant="outlined" to="/" prepend-icon="mdi-home">{{ t('notFound.goHome') }}</v-btn>
      </div>
    </div>
  </div>
</template>
