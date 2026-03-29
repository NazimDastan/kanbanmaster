<script setup lang="ts">
import { computed, ref, onMounted, watch, provide } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { useWebSocket } from '@/composables/useWebSocket'
import AppNavbar from '@/components/layout/AppNavbar.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import ToastContainer from '@/components/common/ToastContainer.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

const route = useRoute()
const authStore = useAuthStore()
const { connect, disconnect } = useWebSocket()
const isAuthPage = computed(() => route.meta.layout === 'auth')
const sidebarOpen = ref(false)

provide('sidebarOpen', sidebarOpen)

function toggleSidebar() { sidebarOpen.value = !sidebarOpen.value }
provide('toggleSidebar', toggleSidebar)

onMounted(() => authStore.fetchUser())

// Connect WebSocket when authenticated, disconnect when logged out
watch(() => authStore.isAuthenticated, (authed) => {
  if (authed) connect()
  else disconnect()
}, { immediate: true })
</script>

<template>
  <v-app>
    <template v-if="isAuthPage">
      <router-view />
    </template>

    <template v-else>
      <div class="h-screen w-screen flex flex-col overflow-hidden bg-[var(--bg,#0a0a0f)]">
        <AppNavbar @toggle-sidebar="toggleSidebar" />

        <div class="flex flex-1 min-h-0 relative">
          <!-- Mobile overlay -->
          <div
            v-if="sidebarOpen"
            class="fixed inset-0 bg-black/50 z-40 md:hidden"
            @click="sidebarOpen = false"
          />

          <!-- Sidebar: hidden on mobile, shown via overlay -->
          <div
            class="fixed md:relative z-50 md:z-auto h-full md:h-auto transition-transform duration-200"
            :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'"
          >
            <AppSidebar @navigate="sidebarOpen = false" />
          </div>

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
    <ConfirmDialog />
  </v-app>
</template>
