<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/useAuthStore'
import { required, email as emailRule, minLength } from '@/utils/validation'

const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()

const form = ref(false)
const nameVal = ref('')
const emailVal = ref('')
const passwordVal = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const error = ref('')

const passwordMatch = (v: string) => v === passwordVal.value || 'Passwords do not match'

async function handleRegister() {
  if (!form.value) return
  error.value = ''
  try {
    await authStore.register(nameVal.value, emailVal.value, passwordVal.value)
    await router.push('/')
  } catch (err: unknown) {
    const axiosErr = err as { response?: { status?: number } }
    if (axiosErr.response?.status === 409) {
      error.value = t('auth.emailInUse')
    } else {
      error.value = t('auth.registerFailed')
    }
  }
}
</script>

<template>
  <div class="min-h-screen flex animated-gradient">
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-[30%] right-[25%] w-[28rem] h-[28rem] bg-secondary/8 rounded-full blur-[130px] animate-[float_9s_ease-in-out_infinite]" />
      <div class="absolute bottom-[25%] left-[20%] w-[24rem] h-[24rem] bg-accent/6 rounded-full blur-[110px] animate-[float_11s_ease-in-out_infinite_3s]" />
    </div>

    <!-- Branding -->
    <div class="hidden lg:flex flex-1 items-center justify-center relative">
      <div class="text-center px-10 max-w-md relative z-10">
        <div class="w-20 h-20 mx-auto mb-8 rounded-2xl bg-gradient-to-br from-secondary via-accent to-primary flex items-center justify-center glow-purple animate-[pulseRing_3s_ease-in-out_infinite]">
          <v-icon icon="mdi-account-group-outline" size="40" color="white" />
        </div>
        <h2 class="text-3xl font-bold mb-3 text-glow">{{ t('auth.brandJoinTitle') }}</h2>
        <p class="text-white/40 leading-relaxed">{{ t('auth.brandJoinDescription') }}</p>
      </div>
    </div>

    <!-- Form -->
    <div class="flex-1 flex items-center justify-center px-4 sm:px-6 py-8 sm:py-12 relative z-10">
      <div class="w-full max-w-sm">
        <div class="flex items-center gap-2.5 mb-10">
          <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center glow-primary">
            <v-icon icon="mdi-view-dashboard" color="white" size="20" />
          </div>
          <span class="text-lg font-bold">KanbanMaster</span>
        </div>

        <h1 class="text-2xl font-bold mb-1.5">{{ t('auth.createAccountTitle') }}</h1>
        <p class="text-sm text-white/40 mb-8">{{ t('auth.createAccountSubtitle') }}</p>

        <v-alert v-if="error" type="error" variant="tonal" closable class="mb-5 text-sm" @click:close="error = ''">{{ error }}</v-alert>

        <v-form v-model="form" @submit.prevent="handleRegister">
          <div class="space-y-3">
            <v-text-field v-model="nameVal" :label="t('auth.fullName')" prepend-inner-icon="mdi-account-outline" :rules="[required(t('auth.fullName')), minLength(2)]" />
            <v-text-field v-model="emailVal" :label="t('auth.email')" type="email" prepend-inner-icon="mdi-email-outline" :rules="[required(t('auth.email')), emailRule]" />
            <v-text-field v-model="passwordVal" :label="t('auth.password')" :type="showPassword ? 'text' : 'password'" prepend-inner-icon="mdi-lock-outline" :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'" :rules="[required(t('auth.password')), minLength(6)]" @click:append-inner="showPassword = !showPassword" />
            <v-text-field v-model="confirmPassword" :label="t('auth.confirmPassword')" :type="showPassword ? 'text' : 'password'" prepend-inner-icon="mdi-lock-check-outline" :rules="[required(t('auth.confirmPassword')), passwordMatch]" />
          </div>
          <v-btn type="submit" size="large" block :loading="authStore.loading" :disabled="!form" class="mt-5 font-semibold text-white bg-gradient-to-br from-[#a855f7] to-[#06b6d4] normal-case tracking-[0]">
            {{ t('auth.register') }}
          </v-btn>
        </v-form>

        <p class="text-center text-xs text-white/30 mt-8">
          {{ t('auth.hasAccount') }}
          <router-link to="/login" class="text-primary-light font-medium hover:underline">{{ t('auth.signIn') }}</router-link>
        </p>
      </div>
    </div>
  </div>
</template>
