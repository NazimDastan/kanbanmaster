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
      <button v-bind="props" class="flex items-center gap-1.5 px-2 py-1.5 rounded-lg hover:bg-white/5 transition-colors">
        <span class="text-sm">{{ languages.find(l => l.code === locale)?.flag }}</span>
        <v-icon icon="mdi-chevron-down" size="12" class="text-white/40" />
      </button>
    </template>
    <v-list density="compact" min-width="140" bg-color="#1e1e32" class="rounded-xl border border-white/5">
      <v-list-item
        v-for="lang in languages"
        :key="lang.code"
        :class="{ 'bg-primary/10': locale === lang.code }"
        @click="switchLanguage(lang.code)"
      >
        <div class="flex items-center gap-2">
          <span class="text-sm">{{ lang.flag }}</span>
          <span class="text-xs">{{ t(`languages.${lang.code}`) }}</span>
        </div>
      </v-list-item>
    </v-list>
  </v-menu>
</template>
