<template>
  <div class="max-w-7xl mx-auto px-12">
    <div class="mb-10">
      <h1 class="text-3xl font-bold text-slate-900 mb-2">Permissions Management</h1>
      <p class="text-slate-600 font-medium">Manage users, roles, and permissions</p>
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
        <div class="flex items-center gap-3">
          <img
            v-if="row.picture_url"
            :src="row.picture_url"
            :alt="row.name"
            class="w-8 h-8 rounded-full object-cover"
          />
          <div
            v-else
            class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-600 to-blue-700 text-white flex items-center justify-center font-semibold text-sm"
          >
            {{ row.name.charAt(0).toUpperCase() }}
          </div>
          <span class="font-medium text-slate-900">{{ row.name }}</span>
        </div>
      </template>

      <template #cell-email="{ value }">
        <span class="text-slate-700">{{ value }}</span>
      </template>

      <template #cell-roles="{ row }">
        <div class="flex flex-wrap gap-2 items-center">
          <BaseBadge
            v-for="role in (row.roles || [])"
            :key="role.id"
            variant="info"
            class="cursor-pointer"
          >
            {{ role.name }}
            <button
              @click.stop="removeRole(row.id, role.id)"
              class="ml-2 bg-white bg-opacity-30 rounded-full w-4 h-4 flex items-center justify-center text-xs hover:bg-opacity-50 transition-colors"
              title="Remove role"
            >
              ×
            </button>
          </BaseBadge>
          <BaseButton
            size="sm"
            variant="ghost"
            @click="openRoleModal(row)"
            class="border-2 border-dashed border-slate-300 hover:border-slate-400"
          >
            + Add Role
          </BaseButton>
        </div>
      </template>

      <template #cell-status="{ row }">
        <BaseBadge :variant="row.active ? 'success' : 'danger'">
          {{ row.active ? 'Active' : 'Inactive' }}
        </BaseBadge>
      </template>

      <template #cell-actions="{ row }">
        <BaseButton
          size="sm"
          :variant="row.active ? 'danger' : 'secondary'"
          @click="toggleUserActive(row)"
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
      <div class="flex flex-col gap-4">
        <button
          v-for="role in availableRoles"
          :key="role.id"
          @click="assignRole(role.id)"
          class="p-5 border-2 border-slate-200 rounded-lg bg-white text-left hover:border-blue-500 hover:bg-blue-50 transition-all duration-200"
        >
          <div class="font-semibold text-slate-900 mb-2">{{ role.name }}</div>
          <div class="text-sm text-slate-600">{{ role.description }}</div>
        </button>
      </div>
    </BaseModal>
  </div>
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
		const [usersRes, rolesRes] = await Promise.all([
			api.get<{ users: User[] }>('/api/v1/permissions/users'),
			api.get<{ roles: Role[] }>('/api/v1/roles'),
		])
		users.value = usersRes.users
		roles.value = rolesRes.roles
	} catch (e: any) {
		error.value = e.message || 'Failed to load data'
	} finally {
		loading.value = false
	}
}

async function toggleUserActive(user: User) {
	try {
		await api.put(`/api/v1/permissions/users/${user.id}/active`, {
			active: !user.active,
		})
		await loadData()
	} catch (e: any) {
		error.value = e.message || 'Failed to update user'
	}
}

function openRoleModal(user: User) {
	selectedUser.value = user
	showRoleModal.value = true
}

function closeRoleModal() {
	selectedUser.value = null
	showRoleModal.value = false
}

async function assignRole(roleId: string) {
	if (!selectedUser.value) return
	try {
		await api.post(`/api/v1/permissions/users/${selectedUser.value.id}/roles`, {
			role_id: roleId,
		})
		await loadData()
		closeRoleModal()
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

function hasRole(user: User, roleName: string): boolean {
	return user.roles?.some((r) => r.name === roleName) || false
}
</script>
