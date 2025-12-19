<template>
  <div class="max-w-7xl mx-auto px-8">
    <div class="mb-8">
      <h1 class="text-3xl font-semibold mb-2">Permissions Management</h1>
      <p class="text-gray-500">Manage users, roles, and permissions</p>
    </div>

    <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mb-6 text-red-600 flex justify-between items-center">
      <span>{{ error }}</span>
      <button @click="error = ''" class="text-2xl leading-none">×</button>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">Loading...</div>

    <div v-else class="bg-white rounded-lg border border-gray-200 overflow-hidden shadow-sm">
      <table class="w-full border-collapse">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">User</th>
            <th class="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">Email</th>
            <th class="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">Roles</th>
            <th class="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">Status</th>
            <th class="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50">
            <td class="px-4 py-4 border-b border-gray-200">
              <div class="flex items-center gap-3">
                <img
                  v-if="user.picture_url"
                  :src="user.picture_url"
                  :alt="user.name"
                  class="w-8 h-8 rounded-full object-cover"
                />
                <div
                  v-else
                  class="w-8 h-8 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 text-white flex items-center justify-center font-semibold text-sm"
                >
                  {{ user.name.charAt(0).toUpperCase() }}
                </div>
                <span>{{ user.name }}</span>
              </div>
            </td>
            <td class="px-4 py-4 border-b border-gray-200 text-gray-900">{{ user.email }}</td>
            <td class="px-4 py-4 border-b border-gray-200">
              <div class="flex flex-wrap gap-2 items-center">
                <span
                  v-for="role in (user.roles || [])"
                  :key="role.id"
                  class="inline-flex items-center gap-2 px-3 py-1.5 bg-indigo-600 text-white rounded text-sm font-medium"
                >
                  {{ role.name }}
                  <button
                    @click="removeRole(user.id, role.id)"
                    class="bg-white bg-opacity-30 rounded-full w-4.5 h-4.5 flex items-center justify-center text-xs hover:bg-opacity-50 transition-colors"
                    title="Remove role"
                  >
                    ×
                  </button>
                </span>
                <button
                  @click="openRoleModal(user)"
                  class="px-3 py-1.5 border border-dashed border-gray-300 bg-white rounded text-sm text-gray-500 hover:border-indigo-500 hover:text-indigo-600 transition-colors"
                >
                  + Add Role
                </button>
              </div>
            </td>
            <td class="px-4 py-4 border-b border-gray-200">
              <span
                :class="[
                  'px-3 py-1.5 rounded text-sm font-medium',
                  user.active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                ]"
              >
                {{ user.active ? 'Active' : 'Inactive' }}
              </span>
            </td>
            <td class="px-4 py-4 border-b border-gray-200">
              <button
                @click="toggleUserActive(user)"
                :class="[
                  'px-4 py-2 border rounded text-sm transition-colors',
                  user.active
                    ? 'text-red-600 border-red-600 hover:bg-red-50'
                    : 'text-gray-700 border-gray-300 hover:bg-gray-50'
                ]"
              >
                {{ user.active ? 'Deactivate' : 'Activate' }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      v-if="showRoleModal && selectedUser"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="closeRoleModal"
    >
      <div
        class="bg-white rounded-lg max-w-md w-full mx-4 max-h-[80vh] overflow-auto shadow-xl"
        @click.stop
      >
        <div class="flex justify-between items-center p-6 border-b border-gray-200">
          <h2 class="text-2xl font-semibold">Assign Role to {{ selectedUser.name }}</h2>
          <button
            @click="closeRoleModal"
            class="text-3xl leading-none text-gray-500 hover:text-gray-700 w-8 h-8 flex items-center justify-center"
          >
            ×
          </button>
        </div>
        <div class="p-6">
          <div class="flex flex-col gap-3">
            <button
              v-for="role in availableRoles"
              :key="role.id"
              @click="assignRole(role.id)"
              class="p-4 border border-gray-200 rounded-md bg-white text-left hover:border-indigo-500 hover:bg-gray-50 transition-all"
            >
              <div class="font-semibold text-gray-900 mb-1">{{ role.name }}</div>
              <div class="text-sm text-gray-500">{{ role.description }}</div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'

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
