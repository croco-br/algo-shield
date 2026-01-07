<template>
  <v-container fluid class="pa-8">
    <div class="mb-10">
      <h2 class="text-h4 font-weight-bold mb-2">{{ $t('views.profile.title') }}</h2>
      <p class="text-body-1 text-grey-darken-1">{{ $t('views.profile.subtitle') }}</p>
    </div>

    <v-card class="pa-8">
      <div class="d-flex align-center gap-3 mb-6">
        <v-avatar color="primary" size="48">
          <v-icon icon="fa-user" color="white" />
        </v-avatar>
        <div>
          <h3 class="text-h6 font-weight-bold">{{ $t('views.profile.userInfo') }}</h3>
          <p class="text-body-2 text-grey-darken-1">{{ $t('views.profile.userInfoSubtitle') }}</p>
        </div>
      </div>

      <v-list class="bg-transparent">
        <v-list-item>
          <template #prepend>
            <v-icon icon="fa-user" color="grey-darken-1" />
          </template>
          <v-list-item-title class="font-weight-semibold">{{ $t('views.profile.name') }}</v-list-item-title>
          <template #append>
            <span class="text-body-2 text-grey-darken-3">{{ user?.name || $t('views.profile.notAvailable') }}</span>
          </template>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon icon="fa-envelope" color="grey-darken-1" />
          </template>
          <v-list-item-title class="font-weight-semibold">{{ $t('views.profile.email') }}</v-list-item-title>
          <template #append>
            <span class="text-body-2 text-grey-darken-3">{{ user?.email || $t('views.profile.notAvailable') }}</span>
          </template>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon icon="fa-key" color="grey-darken-1" />
          </template>
          <v-list-item-title class="font-weight-semibold">{{ $t('views.profile.authType') }}</v-list-item-title>
          <template #append>
            <span class="text-body-2 text-grey-darken-3 uppercase">{{ user?.auth_type || $t('views.profile.notAvailable') }}</span>
          </template>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon icon="fa-shield" color="grey-darken-1" />
          </template>
          <v-list-item-title class="font-weight-semibold">{{ $t('views.profile.roles') }}</v-list-item-title>
          <template #append>
            <div class="d-flex gap-2 flex-wrap align-center">
              <template v-if="user?.roles && user.roles.length > 0">
                <BaseBadge
                  v-for="role in user.roles"
                  :key="role.id"
                  variant="info"
                  size="sm"
                >
                  {{ role.name }}
                </BaseBadge>
              </template>
              <span v-else class="text-body-2 text-grey-darken-1">
                {{ $t('views.profile.noRoles') }}
              </span>
            </div>
          </template>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon :icon="user?.active ? 'fa-check-circle' : 'fa-xmark-circle'" :color="user?.active ? 'success' : 'error'" />
          </template>
          <v-list-item-title class="font-weight-semibold">{{ $t('views.profile.status') }}</v-list-item-title>
          <template #append>
            <BaseBadge
              :variant="user?.active ? 'success' : 'danger'"
              size="sm"
            >
              {{ user?.active ? $t('common.active') : $t('common.inactive') }}
            </BaseBadge>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import BaseBadge from '@/components/BaseBadge.vue'

const authStore = useAuthStore()

const user = computed(() => authStore.user)

// Refresh user data when component mounts to ensure roles are loaded
onMounted(async () => {
  if (authStore.user) {
    await authStore.refresh()
  }
})
</script>
