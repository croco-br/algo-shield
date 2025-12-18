'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/stores/auth';
import { api } from '@/lib/api';

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

export default function PermissionsPage() {
	const { user: currentUser } = useAuth();
	const router = useRouter();
	const [users, setUsers] = useState<User[]>([]);
	const [roles, setRoles] = useState<Role[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState('');
	const [selectedUser, setSelectedUser] = useState<User | null>(null);
	const [showRoleModal, setShowRoleModal] = useState(false);

	useEffect(() => {
		const isAdmin = currentUser?.roles?.some((r) => r.name === 'admin');
		if (currentUser && !isAdmin) {
			router.push('/');
			return;
		}
		if (currentUser) {
			loadData();
		}
	}, [currentUser, router]);

	async function loadData() {
		try {
			setLoading(true);
			setError('');
			const [usersRes, rolesRes] = await Promise.all([
				api.get<{ users: User[] }>('/api/v1/permissions/users'),
				api.get<{ roles: Role[] }>('/api/v1/roles'),
			]);
			setUsers(usersRes.users);
			setRoles(rolesRes.roles);
		} catch (e: any) {
			setError(e.message || 'Failed to load data');
		} finally {
			setLoading(false);
		}
	}

	async function toggleUserActive(user: User) {
		try {
			await api.put(`/api/v1/permissions/users/${user.id}/active`, {
				active: !user.active,
			});
			await loadData();
		} catch (e: any) {
			setError(e.message || 'Failed to update user');
		}
	}

	function openRoleModal(user: User) {
		setSelectedUser(user);
		setShowRoleModal(true);
	}

	function closeRoleModal() {
		setSelectedUser(null);
		setShowRoleModal(false);
	}

	async function assignRole(roleId: string) {
		if (!selectedUser) return;
		try {
			await api.post(`/api/v1/permissions/users/${selectedUser.id}/roles`, {
				role_id: roleId,
			});
			await loadData();
			closeRoleModal();
		} catch (e: any) {
			setError(e.message || 'Failed to assign role');
		}
	}

	async function removeRole(userId: string, roleId: string) {
		try {
			await api.delete(`/api/v1/permissions/users/${userId}/roles/${roleId}`);
			await loadData();
		} catch (e: any) {
			setError(e.message || 'Failed to remove role');
		}
	}

	function hasRole(user: User, roleName: string): boolean {
		return user.roles?.some((r) => r.name === roleName) || false;
	}

	return (
		<div className="max-w-7xl mx-auto px-8">
			<div className="mb-8">
				<h1 className="text-3xl font-semibold mb-2">Permissions Management</h1>
				<p className="text-gray-500">Manage users, roles, and permissions</p>
			</div>

			{error && (
				<div className="bg-red-50 border border-red-200 rounded-md p-4 mb-6 text-red-600 flex justify-between items-center">
					<span>{error}</span>
					<button onClick={() => setError('')} className="text-2xl leading-none">
						×
					</button>
				</div>
			)}

			{loading ? (
				<div className="text-center py-12 text-gray-500">Loading...</div>
			) : (
				<div className="bg-white rounded-lg border border-gray-200 overflow-hidden shadow-sm">
					<table className="w-full border-collapse">
						<thead className="bg-gray-50">
							<tr>
								<th className="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">User</th>
								<th className="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">Email</th>
								<th className="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">Roles</th>
								<th className="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">Status</th>
								<th className="px-4 py-4 text-left font-semibold text-gray-900 border-b-2 border-gray-200">Actions</th>
							</tr>
						</thead>
						<tbody>
							{users.map((user) => (
								<tr key={user.id} className="hover:bg-gray-50">
									<td className="px-4 py-4 border-b border-gray-200">
										<div className="flex items-center gap-3">
											{user.picture_url ? (
												<img
													src={user.picture_url}
													alt={user.name}
													className="w-8 h-8 rounded-full object-cover"
												/>
											) : (
												<div className="w-8 h-8 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 text-white flex items-center justify-center font-semibold text-sm">
													{user.name.charAt(0).toUpperCase()}
												</div>
											)}
											<span>{user.name}</span>
										</div>
									</td>
									<td className="px-4 py-4 border-b border-gray-200 text-gray-900">{user.email}</td>
									<td className="px-4 py-4 border-b border-gray-200">
										<div className="flex flex-wrap gap-2 items-center">
											{(user.roles || []).map((role) => (
												<span
													key={role.id}
													className="inline-flex items-center gap-2 px-3 py-1.5 bg-indigo-600 text-white rounded text-sm font-medium"
												>
													{role.name}
													<button
														onClick={() => removeRole(user.id, role.id)}
														className="bg-white bg-opacity-30 rounded-full w-4.5 h-4.5 flex items-center justify-center text-xs hover:bg-opacity-50 transition-colors"
														title="Remove role"
													>
														×
													</button>
												</span>
											))}
											<button
												onClick={() => openRoleModal(user)}
												className="px-3 py-1.5 border border-dashed border-gray-300 bg-white rounded text-sm text-gray-500 hover:border-indigo-500 hover:text-indigo-600 transition-colors"
											>
												+ Add Role
											</button>
										</div>
									</td>
									<td className="px-4 py-4 border-b border-gray-200">
										<span
											className={`px-3 py-1.5 rounded text-sm font-medium ${
												user.active
													? 'bg-green-100 text-green-800'
													: 'bg-red-100 text-red-800'
											}`}
										>
											{user.active ? 'Active' : 'Inactive'}
										</span>
									</td>
									<td className="px-4 py-4 border-b border-gray-200">
										<button
											onClick={() => toggleUserActive(user)}
											className={`px-4 py-2 border rounded text-sm transition-colors ${
												user.active
													? 'text-red-600 border-red-600 hover:bg-red-50'
													: 'text-gray-700 border-gray-300 hover:bg-gray-50'
											}`}
										>
											{user.active ? 'Deactivate' : 'Activate'}
										</button>
									</td>
								</tr>
							))}
						</tbody>
					</table>
				</div>
			)}

			{showRoleModal && selectedUser && (
				<div
					className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
					onClick={closeRoleModal}
				>
					<div
						className="bg-white rounded-lg max-w-md w-full mx-4 max-h-[80vh] overflow-auto shadow-xl"
						onClick={(e) => e.stopPropagation()}
					>
						<div className="flex justify-between items-center p-6 border-b border-gray-200">
							<h2 className="text-2xl font-semibold">Assign Role to {selectedUser.name}</h2>
							<button
								onClick={closeRoleModal}
								className="text-3xl leading-none text-gray-500 hover:text-gray-700 w-8 h-8 flex items-center justify-center"
							>
								×
							</button>
						</div>
						<div className="p-6">
							<div className="flex flex-col gap-3">
								{roles.map((role) => {
									if (hasRole(selectedUser, role.name)) return null;
									return (
										<button
											key={role.id}
											onClick={() => assignRole(role.id)}
											className="p-4 border border-gray-200 rounded-md bg-white text-left hover:border-indigo-500 hover:bg-gray-50 transition-all"
										>
											<div className="font-semibold text-gray-900 mb-1">{role.name}</div>
											<div className="text-sm text-gray-500">{role.description}</div>
										</button>
									);
								})}
							</div>
						</div>
					</div>
				</div>
			)}
		</div>
	);
}
