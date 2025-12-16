<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { api } from '$lib/api';

	let email = '';
	let password = '';
	let name = '';
	let activeTab: 'login' | 'register' = 'login';
	let loading = false;
	let error = '';

	async function handleLogin() {
		if (!email || !password) {
			error = 'Please enter email and password';
			return;
		}

		loading = true;
		error = '';

		try {
			const response = await api.post<{ token: string; user: any }>('/api/v1/auth/login', {
				email,
				password
			});

			await authStore.setToken(response.token);
			goto('/');
		} catch (e: any) {
			error = e.message || 'Login failed. Please try again.';
		} finally {
			loading = false;
		}
	}

	async function handleRegister() {
		if (!email || !password || !name) {
			error = 'Please fill in all fields';
			return;
		}

		if (password.length < 6) {
			error = 'Password must be at least 6 characters';
			return;
		}

		loading = true;
		error = '';

		try {
			const response = await api.post<{ token: string; user: any }>('/api/v1/auth/register', {
				email,
				password,
				name
			});

			await authStore.setToken(response.token);
			goto('/');
		} catch (e: any) {
			error = e.message || 'Registration failed. Please try again.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="login-container">
	<div class="login-card">
		<div class="login-header">
			<div class="brand-header">
				<img src="/gopher.png" alt="AlgoShield" class="brand-icon" />
				<h1>AlgoShield</h1>
			</div>
			<p>Fraud Detection & Anti-Money Laundering</p>
		</div>

		<div class="tabs">
			<button 
				class="tab" 
				class:active={activeTab === 'login'}
				on:click={() => { activeTab = 'login'; error = ''; }}
			>
				Login
			</button>
			<button 
				class="tab" 
				class:active={activeTab === 'register'}
				on:click={() => { activeTab = 'register'; error = ''; }}
			>
				Register
			</button>
		</div>

		{#if error}
			<div class="error-message">
				{error}
			</div>
		{/if}

		{#if activeTab === 'login'}
			<form on:submit|preventDefault={handleLogin} class="login-form">
				<div class="form-group">
					<label for="email">Email</label>
					<input
						id="email"
						type="email"
						bind:value={email}
						placeholder="user@example.com"
						required
						disabled={loading}
					/>
				</div>

				<div class="form-group">
					<label for="password">Password</label>
					<input
						id="password"
						type="password"
						bind:value={password}
						placeholder="••••••••"
						required
						disabled={loading}
					/>
				</div>

				<button type="submit" class="button button-primary full-width" disabled={loading}>
					{loading ? 'Signing in...' : 'Sign In'}
				</button>
			</form>
		{:else}
			<form on:submit|preventDefault={handleRegister} class="login-form">
				<div class="form-group">
					<label for="reg-name">Name</label>
					<input
						id="reg-name"
						type="text"
						bind:value={name}
						placeholder="Your Name"
						required
						disabled={loading}
					/>
				</div>

				<div class="form-group">
					<label for="reg-email">Email</label>
					<input
						id="reg-email"
						type="email"
						bind:value={email}
						placeholder="user@example.com"
						required
						disabled={loading}
					/>
				</div>

				<div class="form-group">
					<label for="reg-password">Password</label>
					<input
						id="reg-password"
						type="password"
						bind:value={password}
						placeholder="••••••••"
						required
						minlength="6"
						disabled={loading}
					/>
					<small>Minimum 6 characters</small>
				</div>

				<button type="submit" class="button button-primary full-width" disabled={loading}>
					{loading ? 'Creating account...' : 'Create Account'}
				</button>
			</form>
		{/if}
	</div>
</div>

<style>
	.login-container {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		padding: 2rem;
	}

	.login-card {
		background: white;
		border-radius: 12px;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
		width: 100%;
		max-width: 450px;
		padding: 2.5rem;
	}

	.login-header {
		text-align: center;
		margin-bottom: 2rem;
	}

	.brand-header {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
		margin-bottom: 0.5rem;
	}

	.brand-icon {
		width: 48px;
		height: 48px;
		object-fit: contain;
	}

	.login-header h1 {
		font-size: 2rem;
		margin-bottom: 0.5rem;
		color: var(--text);
	}

	.login-header p {
		color: var(--text-secondary);
		font-size: 0.875rem;
	}

	.tabs {
		display: flex;
		gap: 0.5rem;
		margin-bottom: 2rem;
		border-bottom: 2px solid var(--border);
	}

	.tab {
		flex: 1;
		padding: 0.75rem 1rem;
		border: none;
		background: none;
		cursor: pointer;
		font-size: 1rem;
		color: var(--text-secondary);
		border-bottom: 2px solid transparent;
		margin-bottom: -2px;
		transition: all 0.2s;
	}

	.tab:hover {
		color: var(--text);
	}

	.tab.active {
		color: var(--primary);
		border-bottom-color: var(--primary);
		font-weight: 600;
	}

	.error-message {
		background: #fee;
		border: 1px solid #fcc;
		border-radius: 6px;
		padding: 0.75rem 1rem;
		margin-bottom: 1.5rem;
		color: #c33;
		font-size: 0.875rem;
	}

	.login-form {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.form-group label {
		font-weight: 500;
		color: var(--text);
	}

	.form-group input {
		padding: 0.75rem 1rem;
		border: 1px solid var(--border);
		border-radius: 6px;
		font-size: 1rem;
		transition: border-color 0.2s;
		color: var(--text);
	}

	.form-group input:focus {
		outline: none;
		border-color: var(--primary);
	}

	.form-group input::selection {
		background-color: var(--primary);
		color: white;
	}

	.form-group input::-moz-selection {
		background-color: var(--primary);
		color: white;
	}

	.form-group input:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.form-group small {
		font-size: 0.75rem;
		color: var(--text-secondary);
	}

	.button {
		padding: 0.875rem 1.5rem;
		border: none;
		border-radius: 6px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.button-primary {
		background: var(--primary);
		color: white;
	}

	.button-primary:hover:not(:disabled) {
		background: #5568d3;
	}

	.button-primary:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.full-width {
		width: 100%;
	}
</style>
