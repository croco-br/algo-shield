<template>
  <v-container fluid class="pa-8">
    <div class="mb-10">
      <h2 class="text-h4 font-weight-bold mb-2">Profile</h2>
      <p class="text-body-1 text-grey-darken-1">Manage your account information and preferences</p>
    </div>

    <v-card class="pa-8">
      <div class="d-flex align-center gap-3 mb-6">
        <v-avatar color="primary" size="48">
          <v-icon icon="fa-user" color="white" />
        </v-avatar>
        <div>
          <h3 class="text-h6 font-weight-bold">User Information</h3>
          <p class="text-body-2 text-grey-darken-1">View and update your profile details</p>
        </div>
      </div>

      <v-list class="bg-transparent">
        <v-list-item>
          <template #prepend>
            <v-icon icon="fa-user" color="grey-darken-1" />
          </template>
          <v-list-item-title class="font-weight-semibold">Name</v-list-item-title>
          <template #append>
            <span class="text-body-2 text-grey-darken-3">{{ user?.name || 'N/A' }}</span>
          </template>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon icon="fa-envelope" color="grey-darken-1" />
          </template>
          <v-list-item-title class="font-weight-semibold">Email</v-list-item-title>
          <template #append>
            <span class="text-body-2 text-grey-darken-3">{{ user?.email || 'N/A' }}</span>
          </template>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon icon="fa-key" color="grey-darken-1" />
          </template>
          <v-list-item-title class="font-weight-semibold">Authentication Type</v-list-item-title>
          <template #append>
            <span class="text-body-2 text-grey-darken-3 uppercase">{{ user?.auth_type || 'N/A' }}</span>
          </template>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon icon="fa-shield" color="grey-darken-1" />
          </template>
          <v-list-item-title class="font-weight-semibold">Roles</v-list-item-title>
          <template #append>
            <div class="d-flex gap-2 flex-wrap align-center">
              <template v-if="user?.roles && user.roles.length > 0">
                <v-chip
                  v-for="role in user.roles"
                  :key="role.id"
                  color="primary"
                  size="small"
                  variant="flat"
                >
                  {{ role.name }}
                </v-chip>
              </template>
              <span v-else class="text-body-2 text-grey-darken-1">
                No roles assigned
              </span>
            </div>
          </template>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon :icon="user?.active ? 'fa-check-circle' : 'fa-xmark-circle'" :color="user?.active ? 'success' : 'error'" />
          </template>
          <v-list-item-title class="font-weight-semibold">Status</v-list-item-title>
          <template #append>
            <v-chip
              :color="user?.active ? 'success' : 'error'"
              size="small"
              variant="flat"
            >
              {{ user?.active ? 'Active' : 'Inactive' }}
            </v-chip>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const user = computed(() => authStore.user)

// Refresh user data when component mounts to ensure roles are loaded
onMounted(async () => {
  if (authStore.user) {
    await authStore.refresh()
  }
})
</script>
