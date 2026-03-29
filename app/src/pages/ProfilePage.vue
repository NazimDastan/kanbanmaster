<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/useAuthStore'
import { getInitials } from '@/utils/format'
import { required, minLength, email as emailRule } from '@/utils/validation'
import { setLocale } from '@/i18n'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'
import { useTheme } from '@/composables/useTheme'
import type { AccentColor } from '@/composables/useTheme'

const { t, locale } = useI18n()
const toast = useToast()
const { confirm } = useConfirm()
const authStore = useAuthStore()
const { accentColor, accentColors, setAccent } = useTheme()

const accentOptions: { key: AccentColor; gradient: string }[] = [
  { key: 'purple', gradient: 'linear-gradient(135deg, #6366f1, #a855f7)' },
  { key: 'blue', gradient: 'linear-gradient(135deg, #3b82f6, #06b6d4)' },
  { key: 'green', gradient: 'linear-gradient(135deg, #10b981, #14b8a6)' },
  { key: 'orange', gradient: 'linear-gradient(135deg, #f59e0b, #f97316)' },
]

function handleAccentChange(color: AccentColor) {
  setAccent(color)
  toast.info(t(`profile.${color}`))
}

// Profile
const name = ref(authStore.user?.name ?? '')
const emailVal = ref(authStore.user?.email ?? '')
const profileForm = ref(false)
const savingProfile = ref(false)
const profileSaved = ref(false)

// Avatar
const avatarFile = ref<File | null>(null)
const avatarPreview = ref<string | null>(null)
const showPhotoPreview = ref(false)

const hasPhoto = computed(() => !!(avatarPreview.value || authStore.user?.avatarUrl))
const photoSrc = computed(() => avatarPreview.value ?? authStore.user?.avatarUrl ?? '')

// Password
const currentPassword = ref('')
const newPassword = ref('')
const confirmNewPassword = ref('')
const showCurrentPw = ref(false)
const showNewPw = ref(false)
const passwordForm = ref(false)
const savingPassword = ref(false)
const passwordSaved = ref(false)

const passwordMatch = (v: string) => v === newPassword.value || t('validation.passwordMismatch')

const initials = computed(() => getInitials(authStore.user?.name ?? 'U'))

const languages = [
  { code: 'en', label: 'English', flag: '🇬🇧' },
  { code: 'tr', label: 'Türkçe', flag: '🇹🇷' },
  { code: 'de', label: 'Deutsch', flag: '🇩🇪' },
  { code: 'ar', label: 'العربية', flag: '🇸🇦' },
]

async function handleAvatarSelect(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  // Validate file size (max 2MB)
  if (file.size > 2 * 1024 * 1024) {
    toast.error('Max 2MB')
    return
  }

  avatarFile.value = file
  const reader = new FileReader()
  reader.onload = async (e) => {
    const base64 = e.target?.result as string
    avatarPreview.value = base64
    try {
      await authStore.updateAvatar(base64)
      toast.success(t('profile.changeAvatar') + ' ✓')
    } catch {
      toast.error(t('common.error'))
    }
  }
  reader.readAsDataURL(file)
}

function triggerFileInput() {
  const input = document.getElementById('avatar-input') as HTMLInputElement
  input?.click()
}

async function handleSaveProfile() {
  savingProfile.value = true
  profileSaved.value = false
  try {
    await authStore.updateProfile(name.value, emailVal.value)
    toast.success(t('profile.saveChanges') + ' ✓')
    profileSaved.value = true
    setTimeout(() => { profileSaved.value = false }, 3000)
  } catch {
    toast.error(t('common.error'))
  } finally {
    savingProfile.value = false
  }
}

async function handleChangePassword() {
  savingPassword.value = true
  passwordSaved.value = false
  try {
    await authStore.changePassword(currentPassword.value, newPassword.value)
    toast.success(t('profile.updatePassword') + ' ✓')
    currentPassword.value = ''
    newPassword.value = ''
    confirmNewPassword.value = ''
    passwordSaved.value = true
    setTimeout(() => { passwordSaved.value = false }, 3000)
  } catch (err: unknown) {
    const axiosErr = err as { response?: { status?: number } }
    toast.error(axiosErr.response?.status === 401 ? t('auth.invalidCredentials') : t('common.error'))
  } finally {
    savingPassword.value = false
  }
}

async function handleDeleteAccount() {
  const ok = await confirm({
    title: t('profile.deleteAccount'),
    message: t('profile.deleteWarning'),
    confirmText: t('profile.deleteAccount'),
    danger: true,
  })
  if (!ok) return
  toast.info('This feature is coming soon.')
}

function handleLanguageChange(code: string) {
  setLocale(code)
  const lang = languages.find(l => l.code === code)
  toast.info(`${lang?.flag} ${lang?.label}`)
}
</script>

<template>
  <div class="p-4 md:p-6 lg:p-8 max-w-2xl mx-auto space-y-5">
    <!-- Page title -->
    <h1 class="text-xl font-bold">{{ t('profile.title') }}</h1>

    <!-- Avatar & Info Card -->
    <div class="rounded-2xl border border-white/5 bg-card p-5 md:p-6">
      <div class="flex flex-col sm:flex-row items-center sm:items-start gap-5">
        <!-- Avatar -->
        <div class="relative group">
          <div
            class="w-24 h-24 rounded-2xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center overflow-hidden cursor-pointer glow-primary"
            @click="hasPhoto ? showPhotoPreview = true : triggerFileInput()"
          >
            <img
              v-if="hasPhoto"
              :src="photoSrc"
              class="w-full h-full object-cover"
            />
            <span v-else class="text-3xl font-bold text-white">{{ initials }}</span>

            <!-- Hover overlay (only when no photo) -->
            <div
              v-if="!hasPhoto"
              class="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity rounded-2xl"
            >
              <v-icon icon="mdi-camera" size="24" color="white" />
            </div>
          </div>

          <!-- Camera icon overlay to change photo -->
          <button
            v-if="hasPhoto"
            class="absolute -bottom-1 -right-1 w-7 h-7 rounded-full bg-primary flex items-center justify-center shadow-lg border-2 border-[#1e1e2e] hover:scale-110 transition-transform"
            @click.stop="triggerFileInput"
          >
            <v-icon icon="mdi-camera" size="14" color="white" />
          </button>
          <input
            id="avatar-input"
            type="file"
            accept="image/*"
            class="hidden"
            @change="handleAvatarSelect"
          />
          <p class="text-[10px] text-white/25 text-center mt-2">{{ t('profile.changeAvatar') }}</p>
        </div>

        <!-- User info -->
        <div class="flex-1 text-center sm:text-left">
          <h2 class="text-lg font-bold">{{ authStore.user?.name }}</h2>
          <p class="text-sm text-white/40 mt-0.5">{{ authStore.user?.email }}</p>

          <!-- Role badge -->
          <div class="flex items-center gap-2 mt-3 justify-center sm:justify-start">
            <span class="px-2.5 py-1 rounded-lg bg-primary/15 text-primary-light text-[11px] font-medium">
              <v-icon icon="mdi-shield-account" size="12" class="mr-1" />
              {{ t('profile.member') }}
            </span>
            <span class="px-2.5 py-1 rounded-lg bg-white/5 text-white/40 text-[11px]">
              <v-icon icon="mdi-calendar" size="12" class="mr-1" />
              {{ t('task.created') }}: {{ authStore.user?.createdAt?.slice(0, 10) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Personal Info Form -->
    <div class="rounded-2xl border border-white/5 bg-card p-5 md:p-6">
      <div class="flex items-center justify-between mb-5">
        <h3 class="text-sm font-semibold">{{ t('profile.personalInfo') }}</h3>
        <v-icon v-if="profileSaved" icon="mdi-check-circle" color="#10b981" size="18" />
      </div>

      <v-form v-model="profileForm" @submit.prevent="handleSaveProfile">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <v-text-field
            v-model="name"
            :label="t('auth.fullName')"
            :rules="[required(t('auth.fullName')), minLength(2)]"
            prepend-inner-icon="mdi-account-outline"
          />
          <v-text-field
            v-model="emailVal"
            :label="t('auth.email')"
            type="email"
            :rules="[required(t('auth.email')), emailRule]"
            prepend-inner-icon="mdi-email-outline"
          />
        </div>

        <div class="flex justify-end mt-4">
          <v-btn
            type="submit"
            :disabled="!profileForm"
            :loading="savingProfile"
            class="font-medium normal-case bg-gradient-to-br from-[#6366f1] to-[#a855f7] text-white"
          >
            <v-icon icon="mdi-content-save-outline" size="16" class="mr-1.5" />
            {{ t('profile.saveChanges') }}
          </v-btn>
        </div>
      </v-form>
    </div>

    <!-- Password Form -->
    <div class="rounded-2xl border border-white/5 bg-card p-5 md:p-6">
      <div class="flex items-center justify-between mb-5">
        <h3 class="text-sm font-semibold">{{ t('profile.changePassword') }}</h3>
        <v-icon v-if="passwordSaved" icon="mdi-check-circle" color="#10b981" size="18" />
      </div>

      <v-form v-model="passwordForm" @submit.prevent="handleChangePassword">
        <div class="space-y-4">
          <v-text-field
            v-model="currentPassword"
            :label="t('profile.currentPassword')"
            :type="showCurrentPw ? 'text' : 'password'"
            :rules="[required(t('profile.currentPassword'))]"
            prepend-inner-icon="mdi-lock-outline"
            :append-inner-icon="showCurrentPw ? 'mdi-eye-off' : 'mdi-eye'"
            @click:append-inner="showCurrentPw = !showCurrentPw"
          />
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <v-text-field
              v-model="newPassword"
              :label="t('profile.newPassword')"
              :type="showNewPw ? 'text' : 'password'"
              :rules="[required(t('profile.newPassword')), minLength(6)]"
              prepend-inner-icon="mdi-lock-plus-outline"
              :append-inner-icon="showNewPw ? 'mdi-eye-off' : 'mdi-eye'"
              @click:append-inner="showNewPw = !showNewPw"
            />
            <v-text-field
              v-model="confirmNewPassword"
              :label="t('profile.confirmNewPassword')"
              :type="showNewPw ? 'text' : 'password'"
              :rules="[required(t('profile.confirmNewPassword')), passwordMatch]"
              prepend-inner-icon="mdi-lock-check-outline"
            />
          </div>
        </div>

        <div class="flex justify-end mt-4">
          <v-btn
            type="submit"
            :disabled="!passwordForm"
            :loading="savingPassword"
            class="font-medium normal-case bg-gradient-to-br from-[#6366f1] to-[#a855f7] text-white"
          >
            <v-icon icon="mdi-lock-reset" size="16" class="mr-1.5" />
            {{ t('profile.updatePassword') }}
          </v-btn>
        </div>
      </v-form>
    </div>

    <!-- Language Preference -->
    <div class="rounded-2xl border border-white/5 bg-card p-5 md:p-6">
      <h3 class="text-sm font-semibold mb-4">{{ t('profile.language') }}</h3>

      <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
        <button
          v-for="lang in languages"
          :key="lang.code"
          class="flex items-center gap-2.5 px-4 py-3 rounded-xl border transition-all text-left"
          :class="locale === lang.code
            ? 'border-primary/40 bg-primary/10'
            : 'border-white/5 bg-white/[0.02] hover:border-white/10 hover:bg-white/[0.04]'"
          @click="handleLanguageChange(lang.code)"
        >
          <span class="text-lg">{{ lang.flag }}</span>
          <div>
            <p class="text-xs font-medium" :class="locale === lang.code ? 'text-primary-light' : 'text-white/70'">{{ lang.label }}</p>
          </div>
          <v-icon
            v-if="locale === lang.code"
            icon="mdi-check-circle"
            size="14"
            color="#6366f1"
            class="ml-auto"
          />
        </button>
      </div>
    </div>

    <!-- Theme Color -->
    <div class="rounded-2xl border border-white/5 bg-card p-5 md:p-6">
      <h3 class="text-sm font-semibold mb-4">{{ t('profile.themeColor') }}</h3>

      <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
        <button
          v-for="opt in accentOptions"
          :key="opt.key"
          class="flex items-center gap-2.5 px-4 py-3 rounded-xl border transition-all text-left"
          :class="accentColor === opt.key
            ? 'border-primary/40 bg-primary/10'
            : 'border-white/5 bg-white/[0.02] hover:border-white/10 hover:bg-white/[0.04]'"
          @click="handleAccentChange(opt.key)"
        >
          <span
            class="w-6 h-6 rounded-full shrink-0 flex items-center justify-center"
            :style="{ background: opt.gradient }"
          >
            <v-icon
              v-if="accentColor === opt.key"
              icon="mdi-check"
              size="14"
              color="white"
            />
          </span>
          <p
            class="text-xs font-medium"
            :class="accentColor === opt.key ? 'text-primary-light' : 'text-white/70'"
          >
            {{ t(`profile.${opt.key}`) }}
          </p>
        </button>
      </div>
    </div>

    <!-- Danger Zone -->
    <div class="rounded-2xl border border-error/20 bg-error/5 p-5 md:p-6">
      <h3 class="text-sm font-semibold text-error mb-2">{{ t('profile.dangerZone') }}</h3>
      <p class="text-xs text-white/30 mb-4">{{ t('profile.deleteWarning') }}</p>
      <v-btn variant="outlined" color="error" size="small" class="normal-case" prepend-icon="mdi-delete-outline" @click="handleDeleteAccount">
        {{ t('profile.deleteAccount') }}
      </v-btn>
    </div>

    <!-- Photo Preview Dialog -->
    <v-dialog v-model="showPhotoPreview" max-width="440">
      <v-card class="rounded-2xl bg-card pa-4 text-center">
        <div class="flex justify-end mb-2">
          <v-btn icon variant="text" size="small" @click="showPhotoPreview = false">
            <v-icon icon="mdi-close" />
          </v-btn>
        </div>
        <img
          :src="photoSrc"
          class="max-w-[400px] max-h-[400px] w-full rounded-xl object-contain mx-auto"
        />
      </v-card>
    </v-dialog>
  </div>
</template>
