'use client';

import { useEffect, useState } from 'react';
import { api } from '@/lib/api';

interface Rule {
	id: string;
	name: string;
	description: string;
	type: string;
	action: string;
	priority: number;
	enabled: boolean;
	conditions: any;
	score: number;
	created_at: string;
	updated_at: string;
}

const ruleTypes = [
	{ value: 'amount', label: 'Amount Threshold' },
	{ value: 'velocity', label: 'Transaction Velocity' },
	{ value: 'blocklist', label: 'Blocklist' },
	{ value: 'pattern', label: 'Pattern Match' },
	{ value: 'custom', label: 'Custom' },
];

const ruleActions = [
	{ value: 'allow', label: 'Allow' },
	{ value: 'block', label: 'Block' },
	{ value: 'review', label: 'Review' },
	{ value: 'score', label: 'Add Score' },
];

export default function RulesPage() {
	const [rules, setRules] = useState<Rule[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const [showModal, setShowModal] = useState(false);
	const [editingRule, setEditingRule] = useState<Partial<Rule>>({});
	const [isEditing, setIsEditing] = useState(false);

	useEffect(() => {
		loadRules();
	}, []);

	async function loadRules() {
		setLoading(true);
		setError(null);
		try {
			const data = await api.get<{ rules: Rule[] }>('/api/v1/rules');
			setRules(data.rules || []);
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : 'Failed to load rules';
			console.error('Failed to load rules:', error);
			setError(errorMessage);
		} finally {
			setLoading(false);
		}
	}

	function openCreateModal() {
		setEditingRule({
			name: '',
			description: '',
			type: 'amount',
			action: 'score',
			priority: 10,
			enabled: true,
			conditions: {},
			score: 0,
		});
		setIsEditing(false);
		setShowModal(true);
	}

	function openEditModal(rule: Rule) {
		setEditingRule({ ...rule });
		setIsEditing(true);
		setShowModal(true);
	}

	function closeModal() {
		setShowModal(false);
		setEditingRule({});
	}

	async function saveRule() {
		try {
			if (isEditing) {
				await api.put(`/api/v1/rules/${editingRule.id}`, editingRule);
			} else {
				await api.post('/api/v1/rules', editingRule);
			}
			closeModal();
			loadRules();
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : 'Failed to save rule';
			console.error('Failed to save rule:', errorMessage);
			alert(`Error: ${errorMessage}`);
		}
	}

	async function deleteRule(id: string) {
		if (!confirm('Are you sure you want to delete this rule?')) return;

		try {
			await api.delete(`/api/v1/rules/${id}`);
			loadRules();
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : 'Failed to delete rule';
			console.error('Failed to delete rule:', errorMessage);
			alert(`Error: ${errorMessage}`);
		}
	}

	async function toggleRule(rule: Rule) {
		try {
			await api.put(`/api/v1/rules/${rule.id}`, { ...rule, enabled: !rule.enabled });
			loadRules();
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : 'Failed to toggle rule';
			console.error('Failed to toggle rule:', errorMessage);
			alert(`Error: ${errorMessage}`);
		}
	}

	function getActionBadgeClass(action: string) {
		switch (action) {
			case 'block':
				return 'bg-red-100 text-red-800';
			case 'review':
				return 'bg-yellow-100 text-yellow-800';
			case 'allow':
				return 'bg-green-100 text-green-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	}

	return (
		<div className="max-w-7xl mx-auto px-8">
			<div className="flex justify-between items-center mb-8">
				<div>
					<h2 className="text-3xl font-semibold mb-2">Rules Management</h2>
					<p className="text-gray-500">Configure custom rules for fraud detection and AML</p>
				</div>
				<button
					onClick={openCreateModal}
					className="px-6 py-3 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition-colors"
				>
					+ Create Rule
				</button>
			</div>

			{loading ? (
				<div className="bg-white rounded-lg border border-gray-200 p-8">
					<p>Loading rules...</p>
				</div>
			) : error ? (
				<div className="bg-red-50 border border-red-200 rounded-lg p-8">
					<div className="flex items-start">
						<div className="flex-shrink-0">
							<svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
								<path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
							</svg>
						</div>
						<div className="ml-3">
							<h3 className="text-sm font-medium text-red-800">Error loading rules</h3>
							<div className="mt-2 text-sm text-red-700">
								<p>{error}</p>
							</div>
							<div className="mt-4">
								<button
									onClick={loadRules}
									className="text-sm font-medium text-red-800 hover:text-red-900 underline"
								>
									Try again
								</button>
							</div>
						</div>
					</div>
				</div>
			) : rules.length === 0 ? (
				<div className="bg-white rounded-lg border border-gray-200 p-8">
					<p>No rules configured. Create your first rule to get started.</p>
				</div>
			) : (
				<div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
					<table className="w-full border-collapse">
						<thead>
							<tr className="bg-gray-50">
								<th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
									Name
								</th>
								<th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
									Type
								</th>
								<th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
									Action
								</th>
								<th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
									Score
								</th>
								<th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
									Priority
								</th>
								<th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
									Status
								</th>
								<th className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
									Actions
								</th>
							</tr>
						</thead>
						<tbody>
							{rules.map((rule) => (
								<tr key={rule.id} className="hover:bg-gray-50">
									<td className="px-4 py-4 border-b border-gray-200">
										<div className="font-semibold text-gray-900">{rule.name}</div>
										<div className="text-sm text-gray-500">{rule.description}</div>
									</td>
									<td className="px-4 py-4 border-b border-gray-200 text-gray-900">{rule.type}</td>
									<td className="px-4 py-4 border-b border-gray-200">
										<span className={`px-3 py-1 rounded-full text-xs font-medium ${getActionBadgeClass(rule.action)}`}>
											{rule.action}
										</span>
									</td>
									<td className="px-4 py-4 border-b border-gray-200 text-gray-900">{rule.score}</td>
									<td className="px-4 py-4 border-b border-gray-200 text-gray-900">{rule.priority}</td>
									<td className="px-4 py-4 border-b border-gray-200">
										<button
											onClick={() => toggleRule(rule)}
											className={`px-3 py-1 rounded text-sm border ${
												rule.enabled
													? 'bg-green-500 text-white border-green-500'
													: 'bg-white text-gray-700 border-gray-300'
											}`}
										>
											{rule.enabled ? 'Enabled' : 'Disabled'}
										</button>
									</td>
									<td className="px-4 py-4 border-b border-gray-200">
										<button
											onClick={() => openEditModal(rule)}
											className="px-4 py-2 mr-2 bg-indigo-600 text-white rounded text-sm hover:bg-indigo-700 transition-colors"
										>
											Edit
										</button>
										<button
											onClick={() => deleteRule(rule.id)}
											className="px-4 py-2 bg-red-600 text-white rounded text-sm hover:bg-red-700 transition-colors"
										>
											Delete
										</button>
									</td>
								</tr>
							))}
						</tbody>
					</table>
				</div>
			)}

			{showModal && (
				<div
					className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
					onClick={closeModal}
				>
					<div
						className="bg-white rounded-lg p-8 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto"
						onClick={(e) => e.stopPropagation()}
					>
						<h3 className="text-2xl font-semibold mb-6">{isEditing ? 'Edit Rule' : 'Create New Rule'}</h3>

						<div className="space-y-6">
							<div>
								<label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
									Name
								</label>
								<input
									id="name"
									type="text"
									value={editingRule.name || ''}
									onChange={(e) => setEditingRule({ ...editingRule, name: e.target.value })}
									placeholder="Rule name"
									className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								/>
							</div>

							<div>
								<label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
									Description
								</label>
								<input
									id="description"
									type="text"
									value={editingRule.description || ''}
									onChange={(e) => setEditingRule({ ...editingRule, description: e.target.value })}
									placeholder="Description"
									className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								/>
							</div>

							<div>
								<label htmlFor="type" className="block text-sm font-medium text-gray-700 mb-2">
									Type
								</label>
								<select
									id="type"
									value={editingRule.type || 'amount'}
									onChange={(e) => setEditingRule({ ...editingRule, type: e.target.value })}
									className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								>
									{ruleTypes.map((type) => (
										<option key={type.value} value={type.value}>
											{type.label}
										</option>
									))}
								</select>
							</div>

							<div>
								<label htmlFor="action" className="block text-sm font-medium text-gray-700 mb-2">
									Action
								</label>
								<select
									id="action"
									value={editingRule.action || 'score'}
									onChange={(e) => setEditingRule({ ...editingRule, action: e.target.value })}
									className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								>
									{ruleActions.map((action) => (
										<option key={action.value} value={action.value}>
											{action.label}
										</option>
									))}
								</select>
							</div>

							<div>
								<label htmlFor="score" className="block text-sm font-medium text-gray-700 mb-2">
									Risk Score
								</label>
								<input
									id="score"
									type="number"
									value={editingRule.score || 0}
									onChange={(e) => setEditingRule({ ...editingRule, score: parseInt(e.target.value) || 0 })}
									placeholder="0-100"
									className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								/>
							</div>

							<div>
								<label htmlFor="priority" className="block text-sm font-medium text-gray-700 mb-2">
									Priority
								</label>
								<input
									id="priority"
									type="number"
									value={editingRule.priority || 10}
									onChange={(e) => setEditingRule({ ...editingRule, priority: parseInt(e.target.value) || 10 })}
									placeholder="Lower number = higher priority"
									className="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
								/>
							</div>
						</div>

						<div className="flex justify-end gap-4 mt-8">
							<button
								onClick={closeModal}
								className="px-6 py-3 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 transition-colors"
							>
								Cancel
							</button>
							<button
								onClick={saveRule}
								className="px-6 py-3 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors"
							>
								{isEditing ? 'Update' : 'Create'}
							</button>
						</div>
					</div>
				</div>
			)}
		</div>
	);
}
