<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/useAuthStore'
import { required, email as emailRule, minLength } from '@/utils/validation'
import { useToast } from '@/composables/useToast'

const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()
const toast = useToast()

const form = ref(false)
const emailVal = ref('')
const passwordVal = ref('')
const showPassword = ref(false)
const error = ref('')

async function handleLogin() {
  if (!form.value) return
  error.value = ''
  try {
    await authStore.login(emailVal.value, passwordVal.value)
    toast.success(t('auth.welcomeBack') + '!')
    await router.push('/')
  } catch (err: unknown) {
    const axiosErr = err as { response?: { status?: number } }
    error.value = axiosErr.response?.status === 401
      ? t('auth.invalidCredentials')
      : t('auth.loginFailed')
  }
}
</script>

<template>
  <div class="min-h-screen flex animated-gradient">
    <!-- Orbs -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-[20%] left-[20%] w-[30rem] h-[30rem] bg-primary/8 rounded-full blur-[140px]" style="animation: float 8s ease-in-out infinite" />
      <div class="absolute bottom-[20%] right-[20%] w-[25rem] h-[25rem] bg-secondary/8 rounded-full blur-[120px]" style="animation: float 10s ease-in-out infinite 2s" />
    </div>

    <!-- Form side -->
    <div class="flex-1 flex items-center justify-center px-6 py-12 relative z-10">
      <div class="w-full max-w-sm">
        <!-- Logo -->
        <div class="flex items-center gap-2.5 mb-10">
          <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center glow-primary">
            <v-icon icon="mdi-view-dashboard" color="white" size="20" />
          </div>
          <span class="text-lg font-bold">KanbanMaster</span>
        </div>

        <h1 class="text-2xl font-bold mb-1.5">{{ t('auth.welcomeBack') }}</h1>
        <p class="text-sm text-white/40 mb-8">{{ t('auth.signInSubtitle') }}</p>

        <v-alert v-if="error" type="error" variant="tonal" closable class="mb-5 text-sm" @click:close="error = ''">{{ error }}</v-alert>

        <v-form v-model="form" @submit.prevent="handleLogin">
          <div class="space-y-3">
            <v-text-field v-model="emailVal" :label="t('auth.email')" type="email" prepend-inner-icon="mdi-email-outline" :rules="[required(t('auth.email')), emailRule]" />
            <v-text-field v-model="passwordVal" :label="t('auth.password')" :type="showPassword ? 'text' : 'password'" prepend-inner-icon="mdi-lock-outline" :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'" :rules="[required(t('auth.password')), minLength(6)]" @click:append-inner="showPassword = !showPassword" />
          </div>
          <v-btn type="submit" size="large" block :loading="authStore.loading" :disabled="!form" class="mt-5 font-semibold text-white" style="background: linear-gradient(135deg, #6366f1, #a855f7); text-transform: none; letter-spacing: 0;">
            {{ t('auth.login') }}
          </v-btn>
        </v-form>

        <p class="text-center text-xs text-white/30 mt-8">
          {{ t('auth.noAccount') }}
          <router-link to="/register" class="text-primary-light font-medium hover:underline">{{ t('auth.signUp') }}</router-link>
        </p>
      </div>
    </div>

    <!-- Branding side -->
    <div class="hidden lg:flex flex-1 items-center justify-center relative">
      <div class="text-center px-10 max-w-md relative z-10">
        <div class="w-20 h-20 mx-auto mb-8 rounded-2xl bg-gradient-to-br from-primary via-secondary to-accent flex items-center justify-center glow-primary" style="animation: pulseRing 3s ease-in-out infinite">
          <v-icon icon="mdi-view-column-outline" size="40" color="white" />
        </div>
        <h2 class="text-3xl font-bold mb-3 text-glow">{{ t('auth.brandTitle') }}</h2>
        <p class="text-white/40 leading-relaxed">{{ t('auth.brandDescription') }}</p>
        <div class="flex flex-wrap justify-center gap-1.5 mt-6">
          <span class="glass px-2.5 py-1 rounded-full text-[10px] text-white/40">{{ t('auth.featureBoards') }}</span>
          <span class="glass px-2.5 py-1 rounded-full text-[10px] text-white/40">{{ t('auth.featureTeams') }}</span>
          <span class="glass px-2.5 py-1 rounded-full text-[10px] text-white/40">{{ t('auth.featureAnalytics') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
