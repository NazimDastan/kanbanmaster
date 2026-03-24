<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import AppNavbar from '@/components/layout/AppNavbar.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import ToastContainer from '@/components/common/ToastContainer.vue'

const route = useRoute()
const authStore = useAuthStore()
const isAuthPage = computed(() => route.meta.layout === 'auth')

onMounted(() => authStore.fetchUser())
</script>

<template>
  <v-app>
    <!-- Auth layout (login/register) — full screen, no sidebar -->
    <template v-if="isAuthPage">
      <router-view />
    </template>

    <!-- Main layout — fixed navbar top, fixed sidebar left, content scrolls -->
    <template v-else>
      <div class="h-screen w-screen flex flex-col overflow-hidden" style="background: var(--bg, #0a0a0f)">
        <!-- Top navbar — fixed height -->
        <AppNavbar />

        <!-- Below navbar: sidebar + content -->
        <div class="flex flex-1 min-h-0">
          <!-- Sidebar — fixed width, full height below navbar -->
          <AppSidebar />

          <!-- Main content — fills remaining space, scrolls independently -->
          <main class="flex-1 min-w-0 overflow-y-auto">
            <router-view v-slot="{ Component }">
              <transition name="fade-slide" mode="out-in">
                <component :is="Component" />
              </transition>
            </router-view>
          </main>
        </div>
      </div>
    </template>
    <ToastContainer />
  </v-app>
</template>
