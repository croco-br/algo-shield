<template>
  <v-container fluid class="pa-8">
    <div class="mb-10">
      <div class="d-flex align-center gap-3 mb-2">
        <v-icon icon="fa-users-cog" size="large" color="primary" />
        <h1 class="text-h4 font-weight-bold">Permissions Management</h1>
      </div>
      <p class="text-body-1 text-grey-darken-1">Manage users, roles, and permissions</p>
    </div>

    <ErrorMessage
      v-if="error"
      :message="error"
      class="mb-6"
      @dismiss="error = ''"
    />

    <LoadingSpinner v-if="loading" text="Loading..." :centered="false" />

    <BaseTable
      v-else
      :columns="tableColumns"
      :data="users"
      empty-text="No users found"
    >
      <template #cell-user="{ row }">
        <div class="d-flex align-center gap-3">
          <v-avatar size="32" color="primary">
            <v-img
              v-if="row.picture_url"
              :src="row.picture_url"
              :alt="row.name"
              cover
            />
            <span v-else class="text-white text-body-2 font-weight-medium">
              {{ row.name.charAt(0).toUpperCase() }}
            </span>
          </v-avatar>
          <span class="font-weight-medium text-grey-darken-3">{{ row.name }}</span>
        </div>
      </template>

      <template #cell-email="{ value }">
        <span class="text-body-2 text-grey-darken-2">{{ value }}</span>
      </template>

      <template #cell-roles="{ row }">
        <div class="d-flex flex-wrap gap-2 align-center">
          <BaseBadge
            v-for="role in (row.roles || [])"
            :key="role.id"
            variant="info"
            size="sm"
            closable
            @close="removeRole(row.id, role.id)"
          >
            {{ role.name }}
          </BaseBadge>
          <BaseButton
            size="sm"
            variant="ghost"
            @click="openRoleModal(row)"
            prepend-icon="fa-plus-circle"
          >
            Add Role
          </BaseButton>
        </div>
      </template>

      <template #cell-status="{ row }">
        <BaseBadge :variant="row.active ? 'success' : 'danger'" size="sm">
          {{ row.active ? 'Active' : 'Inactive' }}
        </BaseBadge>
      </template>

      <template #cell-actions="{ row }">
        <BaseButton
          size="sm"
          :variant="row.active ? 'danger' : 'secondary'"
          @click="toggleUserActive(row)"
          :prepend-icon="row.active ? 'fa-user-slash' : 'fa-user-check'"
        >
          {{ row.active ? 'Deactivate' : 'Activate' }}
        </BaseButton>
      </template>
    </BaseTable>

    <BaseModal
      v-model="showRoleModal"
      :title="`Assign Role to ${selectedUser?.name || ''}`"
      size="sm"
    >
      <v-list>
        <v-list-item
          v-for="role in availableRoles"
          :key="role.id"
          @click="assignRole(role.id)"
          class="cursor-pointer"
          prepend-icon="fa-shield"
        >
          <v-list-item-title class="font-weight-semibold mb-2">
            {{ role.name }}
          </v-list-item-title>
          <v-list-item-subtitle>
            {{ role.description }}
          </v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </BaseModal>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'
import BaseButton from '@/components/BaseButton.vue'
import BaseBadge from '@/components/BaseBadge.vue'
import BaseModal from '@/components/BaseModal.vue'
import BaseTable from '@/components/BaseTable.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

interface Role {
	id: string;
	name: string;
	description: string;
}

interface User {
	id: string;
	email: string;
	name: string;
	picture_url?: string;
	active: boolean;
	roles: Role[];
}

const router = useRouter()
const authStore = useAuthStore()

const tableColumns = [
	{ key: 'user', label: 'User' },
	{ key: 'email', label: 'Email' },
	{ key: 'roles', label: 'Roles' },
	{ key: 'status', label: 'Status' },
	{ key: 'actions', label: 'Actions' },
]

const users = ref<User[]>([])
const roles = ref<Role[]>([])
const loading = ref(true)
const error = ref('')
const selectedUser = ref<User | null>(null)
const showRoleModal = ref(false)

const availableRoles = computed(() => {
	if (!selectedUser.value) return []
	return roles.value.filter(role => !hasRole(selectedUser.value!, role.name))
})

onMounted(() => {
	if (authStore.user && !authStore.isAdmin) {
		router.push('/')
		return
	}
	if (authStore.user) {
		loadData()
	}
})

async function loadData() {
	try {
		loading.value = true
		error.value = ''
		const [usersResponse, rolesResponse] = await Promise.all([
			api.get<{ users: User[] }>('/api/v1/permissions/users'),
			api.get<Role[]>('/api/v1/roles'),
		])
		users.value = usersResponse?.users || []
		roles.value = rolesResponse || []
	} catch (e: any) {
		error.value = e.message || 'Failed to load data'
		console.error('Error loading data:', e)
	} finally {
		loading.value = false
	}
}

function hasRole(user: User, roleName: string): boolean {
	return user.roles?.some(role => role.name === roleName) || false
}

function openRoleModal(user: User) {
	selectedUser.value = user
	showRoleModal.value = true
}

async function assignRole(roleId: string) {
	if (!selectedUser.value) return
	
	try {
		await api.post(`/api/v1/permissions/users/${selectedUser.value.id}/roles`, { role_id: roleId })
		showRoleModal.value = false
		await loadData()
	} catch (e: any) {
		error.value = e.message || 'Failed to assign role'
	}
}

async function removeRole(userId: string, roleId: string) {
	try {
		await api.delete(`/api/v1/permissions/users/${userId}/roles/${roleId}`)
		await loadData()
	} catch (e: any) {
		error.value = e.message || 'Failed to remove role'
	}
}

async function toggleUserActive(user: User) {
	try {
		await api.put(`/api/v1/permissions/users/${user.id}/active`, { active: !user.active })
		await loadData()
	} catch (e: any) {
		error.value = e.message || 'Failed to toggle user status'
	}
}
</script>
