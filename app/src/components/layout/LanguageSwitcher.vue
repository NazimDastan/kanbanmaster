<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { setLocale } from '@/i18n'

const { locale, t } = useI18n()

const languages = [
  { code: 'en', flag: '🇬🇧' },
  { code: 'tr', flag: '🇹🇷' },
  { code: 'de', flag: '🇩🇪' },
  { code: 'ar', flag: '🇸🇦' },
]

function switchLanguage(code: string) {
  setLocale(code)
}
</script>

<template>
  <v-menu>
    <template #activator="{ props }">
      <button v-bind="props" class="flex items-center gap-1.5 px-2 py-1.5 rounded-lg transition-colors">
        <span class="text-sm">{{ languages.find(l => l.code === locale)?.flag }}</span>
        <v-icon icon="mdi-chevron-down" size="12" class="text-[var(--text-muted)]" />
      </button>
    </template>
    <v-list density="compact" min-width="140" class="rounded-xl bg-[var(--bg-card)] border border-[var(--border)]">
      <v-list-item
        v-for="lang in languages"
        :key="lang.code"
        :class="{ 'bg-primary/10': locale === lang.code }"
        @click="switchLanguage(lang.code)"
      >
        <div class="flex items-center gap-2">
          <span class="text-sm">{{ lang.flag }}</span>
          <span class="text-xs text-[var(--text)]">{{ t(`languages.${lang.code}`) }}</span>
        </div>
      </v-list-item>
    </v-list>
  </v-menu>
</template>
