<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';
	import { authStore } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	interface User {
		id: string;
		email: string;
		name: string;
		picture_url?: string;
		active: boolean;
		roles: Role[];
	}

	interface Role {
		id: string;
		name: string;
		description: string;
	}

	let users: User[] = [];
	let roles: Role[] = [];
	let loading = true;
	let error = '';
	let selectedUser: User | null = null;
	let showRoleModal = false;

	onMount(async () => {
		// Check if user is admin
		const currentUser = await new Promise<User | null>((resolve) => {
			const unsubscribe = authStore.subscribe((user) => {
				unsubscribe();
				resolve(user);
			});
		});

		const isAdmin = currentUser?.roles?.some((r) => r.name === 'admin');
		if (!isAdmin) {
			goto('/');
			return;
		}

		await loadData();
	});

	async function loadData() {
		try {
			loading = true;
			error = '';
		const [usersRes, rolesRes] = await Promise.all([
			api.get<{ users: User[] }>('/api/v1/permissions/users'),
			api.get<{ roles: Role[] }>('/api/v1/roles')
		]);
			users = usersRes.users;
			roles = rolesRes.roles;
		} catch (e: any) {
			error = e.message || 'Failed to load data';
		} finally {
			loading = false;
		}
	}

	async function toggleUserActive(user: User) {
		try {
			await api.put(`/api/v1/permissions/users/${user.id}/active`, {
				active: !user.active
			});
			await loadData();
		} catch (e: any) {
			error = e.message || 'Failed to update user';
		}
	}

	function openRoleModal(user: User) {
		selectedUser = user;
		showRoleModal = true;
	}

	function closeRoleModal() {
		selectedUser = null;
		showRoleModal = false;
	}

	async function assignRole(roleId: string) {
		if (!selectedUser) return;
		try {
			await api.post(`/api/v1/permissions/users/${selectedUser.id}/roles`, {
				role_id: roleId
			});
			await loadData();
			closeRoleModal();
		} catch (e: any) {
			error = e.message || 'Failed to assign role';
		}
	}

	async function removeRole(userId: string, roleId: string) {
		try {
			await api.delete(`/api/v1/permissions/users/${userId}/roles/${roleId}`);
			await loadData();
		} catch (e: any) {
			error = e.message || 'Failed to remove role';
		}
	}

	function hasRole(user: User, roleName: string): boolean {
		return user.roles?.some((r) => r.name === roleName) || false;
	}
</script>

<div class="permissions-container">
	<div class="header">
		<h1>Permissions Management</h1>
		<p>Manage users, roles, and permissions</p>
	</div>

	{#if error}
		<div class="error-banner">
			{error}
			<button on:click={() => (error = '')}>×</button>
		</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else}
		<div class="users-table">
			<table>
				<thead>
					<tr>
						<th>User</th>
						<th>Email</th>
						<th>Roles</th>
						<th>Status</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each users as user}
						<tr>
							<td>
								<div class="user-cell">
									{#if user.picture_url}
										<img src={user.picture_url} alt={user.name} class="user-avatar-small" />
									{:else}
										<div class="user-avatar-small">{user.name.charAt(0).toUpperCase()}</div>
									{/if}
									<span>{user.name}</span>
								</div>
							</td>
							<td>{user.email}</td>
							<td>
								<div class="roles-cell">
									{#each user.roles || [] as role}
										<span class="role-tag">
											{role.name}
											<button
												class="remove-role"
												on:click={() => removeRole(user.id, role.id)}
												title="Remove role"
											>
												×
											</button>
										</span>
									{/each}
									<button class="add-role-btn" on:click={() => openRoleModal(user)}>
										+ Add Role
									</button>
								</div>
							</td>
							<td>
								<span class="status-badge" class:active={user.active} class:inactive={!user.active}>
									{user.active ? 'Active' : 'Inactive'}
								</span>
							</td>
							<td>
								<button
									class="toggle-btn"
									on:click={() => toggleUserActive(user)}
									class:deactivate={user.active}
								>
									{user.active ? 'Deactivate' : 'Activate'}
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

{#if showRoleModal && selectedUser}
	<div class="modal-overlay" on:click={closeRoleModal}>
		<div class="modal" on:click|stopPropagation>
			<div class="modal-header">
				<h2>Assign Role to {selectedUser.name}</h2>
				<button class="close-btn" on:click={closeRoleModal}>×</button>
			</div>
			<div class="modal-body">
				<div class="roles-list">
					{#each roles as role}
						{#if !hasRole(selectedUser, role.name)}
							<button class="role-option" on:click={() => assignRole(role.id)}>
								<div class="role-name">{role.name}</div>
								<div class="role-description">{role.description}</div>
							</button>
						{/if}
					{/each}
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.permissions-container {
		max-width: 1400px;
		margin: 0 auto;
		padding: 2rem;
	}

	.header {
		margin-bottom: 2rem;
	}

	.header h1 {
		font-size: 2rem;
		margin-bottom: 0.5rem;
		color: var(--text);
	}

	.header p {
		color: var(--text-secondary);
	}

	.error-banner {
		background: #fee;
		border: 1px solid #fcc;
		border-radius: 6px;
		padding: 1rem;
		margin-bottom: 1.5rem;
		color: #c33;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.error-banner button {
		background: none;
		border: none;
		font-size: 1.5rem;
		cursor: pointer;
		color: #c33;
	}

	.loading {
		text-align: center;
		padding: 3rem;
		color: var(--text-secondary);
	}

	.users-table {
		background: white;
		border-radius: 8px;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
		overflow: hidden;
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	thead {
		background: var(--surface);
	}

	th {
		padding: 1rem;
		text-align: left;
		font-weight: 600;
		color: var(--text);
		border-bottom: 2px solid var(--border);
	}

	td {
		padding: 1rem;
		border-bottom: 1px solid var(--border);
	}

	tr:hover {
		background: var(--surface);
	}

	.user-cell {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.user-avatar-small {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		font-size: 0.875rem;
		overflow: hidden;
	}

	.user-avatar-small img {
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.roles-cell {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
		align-items: center;
	}

	.role-tag {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.375rem 0.75rem;
		background: var(--primary);
		color: white;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.remove-role {
		background: rgba(255, 255, 255, 0.3);
		border: none;
		border-radius: 50%;
		width: 18px;
		height: 18px;
		cursor: pointer;
		color: white;
		font-size: 1rem;
		line-height: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0;
	}

	.remove-role:hover {
		background: rgba(255, 255, 255, 0.5);
	}

	.add-role-btn {
		padding: 0.375rem 0.75rem;
		border: 1px dashed var(--border);
		background: white;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
		color: var(--text-secondary);
		transition: all 0.2s;
	}

	.add-role-btn:hover {
		border-color: var(--primary);
		color: var(--primary);
	}

	.status-badge {
		padding: 0.375rem 0.75rem;
		border-radius: 4px;
		font-size: 0.875rem;
		font-weight: 500;
	}

	.status-badge.active {
		background: #d4edda;
		color: #155724;
	}

	.status-badge.inactive {
		background: #f8d7da;
		color: #721c24;
	}

	.toggle-btn {
		padding: 0.5rem 1rem;
		border: 1px solid var(--border);
		background: white;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
		transition: all 0.2s;
	}

	.toggle-btn:hover {
		background: var(--surface);
	}

	.toggle-btn.deactivate {
		color: var(--danger);
		border-color: var(--danger);
	}

	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal {
		background: white;
		border-radius: 8px;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
		max-width: 500px;
		width: 90%;
		max-height: 80vh;
		overflow: auto;
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1.5rem;
		border-bottom: 1px solid var(--border);
	}

	.modal-header h2 {
		margin: 0;
		font-size: 1.5rem;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: 2rem;
		cursor: pointer;
		color: var(--text-secondary);
		line-height: 1;
		padding: 0;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.modal-body {
		padding: 1.5rem;
	}

	.roles-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.role-option {
		padding: 1rem;
		border: 1px solid var(--border);
		border-radius: 6px;
		background: white;
		cursor: pointer;
		text-align: left;
		transition: all 0.2s;
	}

	.role-option:hover {
		border-color: var(--primary);
		background: var(--surface);
	}

	.role-name {
		font-weight: 600;
		color: var(--text);
		margin-bottom: 0.25rem;
	}

	.role-description {
		font-size: 0.875rem;
		color: var(--text-secondary);
	}
</style>
