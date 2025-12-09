<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api';

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

	let rules: Rule[] = [];
	let loading = true;
	let showModal = false;
	let editingRule: Partial<Rule> = {};
	let isEditing = false;

	const ruleTypes = [
		{ value: 'amount', label: 'Amount Threshold' },
		{ value: 'velocity', label: 'Transaction Velocity' },
		{ value: 'blocklist', label: 'Blocklist' },
		{ value: 'pattern', label: 'Pattern Match' },
		{ value: 'custom', label: 'Custom' }
	];

	const ruleActions = [
		{ value: 'allow', label: 'Allow' },
		{ value: 'block', label: 'Block' },
		{ value: 'review', label: 'Review' },
		{ value: 'score', label: 'Add Score' }
	];

	onMount(() => {
		loadRules();
	});

	async function loadRules() {
		loading = true;
		try {
			const data = await api.get<{ rules: Rule[] }>('/api/v1/rules');
			rules = data.rules || [];
		} catch (error) {
			console.error('Failed to load rules:', error);
		}
		loading = false;
	}

	function openCreateModal() {
		editingRule = {
			name: '',
			description: '',
			type: 'amount',
			action: 'score',
			priority: 10,
			enabled: true,
			conditions: {},
			score: 0
		};
		isEditing = false;
		showModal = true;
	}

	function openEditModal(rule: Rule) {
		editingRule = { ...rule };
		isEditing = true;
		showModal = true;
	}

	function closeModal() {
		showModal = false;
		editingRule = {};
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
			console.error('Failed to save rule:', error);
		}
	}

	async function deleteRule(id: string) {
		if (!confirm('Are you sure you want to delete this rule?')) return;

		try {
			await api.delete(`/api/v1/rules/${id}`);
			loadRules();
		} catch (error) {
			console.error('Failed to delete rule:', error);
		}
	}

	async function toggleRule(rule: Rule) {
		try {
			const response = await fetch(`${API_BASE}/rules/${rule.id}`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ ...rule, enabled: !rule.enabled })
			});

			if (response.ok) {
				loadRules();
			}
		} catch (error) {
			console.error('Failed to toggle rule:', error);
		}
	}

	function getActionBadgeClass(action: string) {
		switch (action) {
			case 'block':
				return 'badge-danger';
			case 'review':
				return 'badge-warning';
			case 'allow':
				return 'badge-success';
			default:
				return '';
		}
	}
</script>

<div class="container">
	<div class="header-section">
		<div>
			<h2>Rules Management</h2>
			<p>Configure custom rules for fraud detection and AML</p>
		</div>
		<button class="button button-primary" on:click={openCreateModal}>
			+ Create Rule
		</button>
	</div>

	{#if loading}
		<div class="card">
			<p>Loading rules...</p>
		</div>
	{:else if rules.length === 0}
		<div class="card">
			<p>No rules configured. Create your first rule to get started.</p>
		</div>
	{:else}
		<div class="card">
			<table class="table">
				<thead>
					<tr>
						<th>Name</th>
						<th>Type</th>
						<th>Action</th>
						<th>Score</th>
						<th>Priority</th>
						<th>Status</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each rules as rule}
						<tr>
							<td>
								<strong>{rule.name}</strong>
								<br />
								<small style="color: var(--text-secondary)">
									{rule.description}
								</small>
							</td>
							<td>{rule.type}</td>
							<td>
								<span class="badge {getActionBadgeClass(rule.action)}">
									{rule.action}
								</span>
							</td>
							<td>{rule.score}</td>
							<td>{rule.priority}</td>
							<td>
								<button
									class="toggle-btn"
									class:enabled={rule.enabled}
									on:click={() => toggleRule(rule)}
								>
									{rule.enabled ? 'Enabled' : 'Disabled'}
								</button>
							</td>
							<td>
								<button
									class="action-btn"
									on:click={() => openEditModal(rule)}
								>
									Edit
								</button>
								<button
									class="action-btn danger"
									on:click={() => deleteRule(rule.id)}
								>
									Delete
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

{#if showModal}
	<div class="modal-overlay" on:click={closeModal}>
		<div class="modal" on:click|stopPropagation>
			<h3>{isEditing ? 'Edit Rule' : 'Create New Rule'}</h3>

			<div class="form-group">
				<label class="label" for="name">Name</label>
				<input
					id="name"
					class="input"
					type="text"
					bind:value={editingRule.name}
					placeholder="Rule name"
				/>
			</div>

			<div class="form-group">
				<label class="label" for="description">Description</label>
				<input
					id="description"
					class="input"
					type="text"
					bind:value={editingRule.description}
					placeholder="Description"
				/>
			</div>

			<div class="form-group">
				<label class="label" for="type">Type</label>
				<select id="type" class="select" bind:value={editingRule.type}>
					{#each ruleTypes as type}
						<option value={type.value}>{type.label}</option>
					{/each}
				</select>
			</div>

			<div class="form-group">
				<label class="label" for="action">Action</label>
				<select id="action" class="select" bind:value={editingRule.action}>
					{#each ruleActions as action}
						<option value={action.value}>{action.label}</option>
					{/each}
				</select>
			</div>

			<div class="form-group">
				<label class="label" for="score">Risk Score</label>
				<input
					id="score"
					class="input"
					type="number"
					bind:value={editingRule.score}
					placeholder="0-100"
				/>
			</div>

			<div class="form-group">
				<label class="label" for="priority">Priority</label>
				<input
					id="priority"
					class="input"
					type="number"
					bind:value={editingRule.priority}
					placeholder="Lower number = higher priority"
				/>
			</div>

			<div class="modal-actions">
				<button class="button" on:click={closeModal}>Cancel</button>
				<button class="button button-primary" on:click={saveRule}>
					{isEditing ? 'Update' : 'Create'}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.header-section {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 2rem;
	}

	h2 {
		font-size: 2rem;
		margin-bottom: 0.5rem;
	}

	.toggle-btn {
		padding: 0.25rem 0.75rem;
		border: 1px solid var(--border);
		border-radius: 4px;
		background: white;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.toggle-btn.enabled {
		background: var(--success);
		color: white;
		border-color: var(--success);
	}

	.action-btn {
		padding: 0.5rem 1rem;
		margin-right: 0.5rem;
		border: none;
		border-radius: 4px;
		background: var(--primary);
		color: white;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.action-btn.danger {
		background: var(--danger);
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
		padding: 2rem;
		max-width: 600px;
		width: 90%;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal h3 {
		margin-bottom: 1.5rem;
		font-size: 1.5rem;
	}

	.modal-actions {
		display: flex;
		justify-content: flex-end;
		gap: 1rem;
		margin-top: 2rem;
	}
</style>

