<template>
  <div class="max-w-7xl mx-auto px-12">
    <div class="mb-10">
      <h2 class="text-3xl font-bold text-slate-900 mb-2">Profile</h2>
      <p class="text-slate-600 font-medium">Manage your account information and preferences</p>
    </div>

    <div class="bg-white rounded-xl border-2 border-slate-200 shadow-sm p-8">
      <div class="flex items-center gap-3 mb-6">
        <div class="w-12 h-12 rounded-full bg-teal-100 flex items-center justify-center">
          <i class="fas fa-user text-teal-600 text-xl"></i>
        </div>
        <div>
          <h3 class="text-xl font-bold text-slate-900">User Information</h3>
          <p class="text-sm text-slate-600">View and update your profile details</p>
        </div>
      </div>

      <div class="space-y-5">
        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-user text-slate-600"></i>
            <span class="font-semibold text-slate-700">Name</span>
          </div>
          <span class="text-sm font-medium text-slate-900">{{ user?.name || 'N/A' }}</span>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-envelope text-slate-600"></i>
            <span class="font-semibold text-slate-700">Email</span>
          </div>
          <span class="text-sm font-medium text-slate-900">{{ user?.email || 'N/A' }}</span>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-key text-slate-600"></i>
            <span class="font-semibold text-slate-700">Authentication Type</span>
          </div>
          <span class="text-sm font-medium text-slate-900 uppercase">{{ user?.authType || 'N/A' }}</span>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-shield-alt text-slate-600"></i>
            <span class="font-semibold text-slate-700">Roles</span>
          </div>
          <div class="flex gap-2">
            <span
              v-for="role in user?.roles"
              :key="role.id"
              class="text-xs font-medium px-3 py-1 bg-teal-100 text-teal-700 rounded-full"
            >
              {{ role.name }}
            </span>
            <span v-if="!user?.roles || user.roles.length === 0" class="text-sm font-medium text-slate-500">No roles assigned</span>
          </div>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <i class="fas fa-circle text-slate-600"></i>
            <span class="font-semibold text-slate-700">Status</span>
          </div>
          <span
            :class="[
              'text-sm font-medium px-3 py-1 rounded-full',
              user?.active ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'
            ]"
          >
            {{ user?.active ? 'Active' : 'Inactive' }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const user = computed(() => authStore.user)
</script>
