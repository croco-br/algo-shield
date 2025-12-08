<script lang="ts">
	import { onMount } from 'svelte';

	interface Field {
		name: string;
		type: 'string' | 'number' | 'boolean' | 'date' | 'uuid' | 'email' | 'ip';
		required: boolean;
		format?: string;
	}

	interface TestResult {
		timestamp: string;
		data: any;
		score: number;
		action: string;
		status: 'success' | 'error';
	}

	let fields: Field[] = [
		{ name: 'transaction_id', type: 'uuid', required: true },
		{ name: 'amount', type: 'number', required: true },
		{ name: 'user_email', type: 'email', required: true }
	];

	let testResults: TestResult[] = [];
	let loading = false;
	let shouldStop = false;
	let currentProgress = 0;
	let numberOfEvents = 1;
	let showAddField = false;
	let newField: Partial<Field> = { name: '', type: 'string', required: true };

	const fieldTypes = [
		{ value: 'string', label: 'String' },
		{ value: 'number', label: 'Number' },
		{ value: 'boolean', label: 'Boolean' },
		{ value: 'date', label: 'Date/Time' },
		{ value: 'uuid', label: 'UUID' },
		{ value: 'email', label: 'Email' },
		{ value: 'ip', label: 'IP Address' }
	];

	function addField() {
		if (newField.name && newField.type) {
			fields = [...fields, newField as Field];
			newField = { name: '', type: 'string', required: true };
			showAddField = false;
		}
	}

	function removeField(index: number) {
		fields = fields.filter((_, i) => i !== index);
	}

	function generateValue(field: Field): any {
		switch (field.type) {
			case 'uuid':
				return crypto.randomUUID();
			case 'email':
				return `user${Math.floor(Math.random() * 10000)}@example.com`;
			case 'number':
				return Math.floor(Math.random() * 10000) + 100;
			case 'boolean':
				return Math.random() > 0.5;
			case 'date':
				return new Date().toISOString();
			case 'ip':
				return `${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}`;
			case 'string':
			default:
				return `value_${Math.random().toString(36).substring(7)}`;
		}
	}

	function generateSyntheticData(): any {
		const data: any = {};
		fields.forEach(field => {
			data[field.name] = generateValue(field);
		});
		return data;
	}

	async function runTest() {
		loading = true;
		shouldStop = false;
		currentProgress = 0;
		testResults = [];

		try {
			for (let i = 0; i < numberOfEvents; i++) {
				// Check if user requested stop
				if (shouldStop) {
					console.log('Test stopped by user');
					break;
				}

				const syntheticData = generateSyntheticData();
				
				// Simulate API call
				await new Promise(resolve => setTimeout(resolve, 300));
				
				const mockScore = Math.floor(Math.random() * 100);
				const mockAction = mockScore > 80 ? 'block' : mockScore > 50 ? 'review' : 'allow';
				
				testResults = [...testResults, {
					timestamp: new Date().toISOString(),
					data: syntheticData,
					score: mockScore,
					action: mockAction,
					status: 'success'
				}];

				currentProgress = i + 1;
			}
		} catch (error) {
			console.error('Test failed:', error);
		}

		loading = false;
		currentProgress = 0;
	}

	function stopTest() {
		shouldStop = true;
	}

	function exportResults() {
		const dataStr = JSON.stringify(testResults, null, 2);
		const blob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `synthetic-test-${Date.now()}.json`;
		link.click();
		URL.revokeObjectURL(url);
	}

	function clearResults() {
		testResults = [];
	}
</script>

<div class="container">
	<div class="header-section">
		<div>
			<h2>Synthetic Data Testing</h2>
			<p>Generate and test synthetic events with custom field schemas</p>
		</div>
	</div>

	<div class="grid">
		<!-- Field Configuration Panel -->
		<div class="card">
			<h3>Event Schema</h3>
			<p class="subtitle">Define the fields for your synthetic events</p>

			<div class="fields-list">
				{#each fields as field, index}
					<div class="field-item">
						<div class="field-info">
							<div class="field-name">
								{field.name}
								{#if field.required}
									<span class="required-badge">Required</span>
								{/if}
							</div>
							<div class="field-type">{field.type}</div>
						</div>
						<button 
							class="remove-btn"
							on:click={() => removeField(index)}
							title="Remove field"
						>
							✕
						</button>
					</div>
				{/each}

				{#if showAddField}
					<div class="add-field-form">
						<input
							class="input"
							type="text"
							bind:value={newField.name}
							placeholder="Field name"
						/>
						<select class="select" bind:value={newField.type}>
							{#each fieldTypes as type}
								<option value={type.value}>{type.label}</option>
							{/each}
						</select>
						<label class="checkbox-label">
							<input type="checkbox" bind:checked={newField.required} />
							Required
						</label>
						<div class="button-group">
							<button class="button button-sm" on:click={addField}>Add</button>
							<button class="button button-sm" on:click={() => showAddField = false}>Cancel</button>
						</div>
					</div>
				{:else}
					<button 
						class="button button-outline"
						on:click={() => showAddField = true}
					>
						+ Add Field
					</button>
				{/if}
			</div>
		</div>

		<!-- Test Configuration Panel -->
		<div class="card">
			<h3>Test Configuration</h3>
			<p class="subtitle">Configure and run your tests</p>

			<div class="form-group">
				<label for="numEvents">Number of Events</label>
				<input
					id="numEvents"
					class="input"
					type="number"
					min="1"
					max="1000"
					bind:value={numberOfEvents}
					disabled={loading}
				/>
			</div>

			{#if loading}
				<div class="progress-section">
					<div class="progress-info">
						<span>Processing events...</span>
						<span class="progress-count">{currentProgress} / {numberOfEvents}</span>
					</div>
					<div class="progress-bar">
						<div 
							class="progress-fill" 
							style="width: {(currentProgress / numberOfEvents) * 100}%"
						></div>
					</div>
				</div>

				<button 
					class="button button-danger full-width"
					on:click={stopTest}
				>
					⏹ Stop Test
				</button>
			{:else}
				<button 
					class="button button-primary full-width"
					on:click={runTest}
					disabled={fields.length === 0}
				>
					▶ Run Test
				</button>
			{/if}

			{#if testResults.length > 0}
				<div class="test-stats">
					<div class="stat">
						<div class="stat-value">{testResults.length}</div>
						<div class="stat-label">Events Tested</div>
					</div>
					<div class="stat">
						<div class="stat-value">
							{testResults.filter(r => r.action === 'block').length}
						</div>
						<div class="stat-label">Blocked</div>
					</div>
					<div class="stat">
						<div class="stat-value">
							{testResults.filter(r => r.action === 'review').length}
						</div>
						<div class="stat-label">Review</div>
					</div>
					<div class="stat">
						<div class="stat-value">
							{testResults.filter(r => r.action === 'allow').length}
						</div>
						<div class="stat-label">Allowed</div>
					</div>
				</div>

				<div class="button-group">
					<button class="button button-outline" on:click={exportResults}>
						Export Results
					</button>
					<button class="button button-outline" on:click={clearResults}>
						Clear
					</button>
				</div>
			{/if}
		</div>
	</div>

	<!-- Results Panel -->
	{#if testResults.length > 0}
		<div class="card results-card">
			<h3>Test Results</h3>
			
			<div class="results-list">
				{#each testResults as result, index}
					<div class="result-item">
						<div class="result-header">
							<span class="result-number">#{index + 1}</span>
							<span class="result-time">{new Date(result.timestamp).toLocaleTimeString()}</span>
							<span class="badge badge-{result.action}">{result.action}</span>
							<span class="result-score">Score: {result.score}</span>
						</div>
						<details class="result-details">
							<summary>View Event Data</summary>
							<pre>{JSON.stringify(result.data, null, 2)}</pre>
						</details>
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>

<style>
	.header-section {
		margin-bottom: 2rem;
	}

	h2 {
		font-size: 2rem;
		margin-bottom: 0.5rem;
	}

	h3 {
		font-size: 1.25rem;
		margin-bottom: 0.5rem;
	}

	.subtitle {
		color: var(--text-secondary);
		font-size: 0.875rem;
		margin-bottom: 1.5rem;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
		gap: 1.5rem;
		margin-bottom: 1.5rem;
	}

	.fields-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.field-item {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0.75rem;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 6px;
	}

	.field-info {
		flex: 1;
	}

	.field-name {
		font-weight: 500;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.field-type {
		color: var(--text-secondary);
		font-size: 0.875rem;
	}

	.required-badge {
		font-size: 0.75rem;
		padding: 0.125rem 0.5rem;
		background: var(--warning);
		color: white;
		border-radius: 3px;
	}

	.remove-btn {
		padding: 0.25rem 0.5rem;
		border: none;
		background: none;
		color: var(--danger);
		cursor: pointer;
		font-size: 1.25rem;
		line-height: 1;
	}

	.remove-btn:hover {
		color: var(--danger);
		opacity: 0.8;
	}

	.add-field-form {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		padding: 1rem;
		background: var(--surface);
		border: 2px dashed var(--border);
		border-radius: 6px;
	}

	.checkbox-label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
	}

	.button-group {
		display: flex;
		gap: 0.5rem;
	}

	.button-sm {
		padding: 0.5rem 1rem;
		font-size: 0.875rem;
	}

	.full-width {
		width: 100%;
	}

	.progress-section {
		margin-bottom: 1rem;
		padding: 1rem;
		background: var(--surface);
		border-radius: 6px;
		border: 1px solid var(--border);
	}

	.progress-info {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.75rem;
		font-size: 0.875rem;
	}

	.progress-count {
		font-weight: 600;
		color: var(--primary);
	}

	.progress-bar {
		width: 100%;
		height: 8px;
		background: var(--border);
		border-radius: 4px;
		overflow: hidden;
	}

	.progress-fill {
		height: 100%;
		background: linear-gradient(90deg, var(--primary), #764ba2);
		transition: width 0.3s ease;
	}

	.test-stats {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 1rem;
		margin-top: 1.5rem;
		padding: 1.5rem;
		background: var(--surface);
		border-radius: 6px;
	}

	.stat {
		text-align: center;
	}

	.stat-value {
		font-size: 2rem;
		font-weight: 700;
		color: var(--primary);
	}

	.stat-label {
		font-size: 0.875rem;
		color: var(--text-secondary);
		margin-top: 0.25rem;
	}

	.results-card {
		margin-top: 1.5rem;
	}

	.results-list {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.result-item {
		border: 1px solid var(--border);
		border-radius: 6px;
		padding: 1rem;
	}

	.result-header {
		display: flex;
		align-items: center;
		gap: 1rem;
		margin-bottom: 0.5rem;
	}

	.result-number {
		font-weight: 700;
		color: var(--text-secondary);
	}

	.result-time {
		color: var(--text-secondary);
		font-size: 0.875rem;
	}

	.result-score {
		margin-left: auto;
		font-weight: 500;
	}

	.badge-allow {
		background: var(--success);
	}

	.badge-review {
		background: var(--warning);
	}

	.badge-block {
		background: var(--danger);
	}

	.result-details {
		margin-top: 0.75rem;
	}

	.result-details summary {
		cursor: pointer;
		color: var(--primary);
		font-size: 0.875rem;
		user-select: none;
	}

	.result-details pre {
		margin-top: 0.75rem;
		padding: 1rem;
		background: var(--surface);
		border-radius: 4px;
		overflow-x: auto;
		font-size: 0.875rem;
	}
</style>
