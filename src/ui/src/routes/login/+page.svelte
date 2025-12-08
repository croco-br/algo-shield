<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';

	let email = '';
	let password = '';
	let activeTab: 'local' | 'sso' = 'local';
	let loading = false;

	async function handleLocalLogin() {
		if (!email || !password) {
			alert('Please enter email and password');
			return;
		}
		
		loading = true;
		// Simulate API call
		await new Promise(resolve => setTimeout(resolve, 500));
		
		authStore.login(email, password);
		loading = false;
		goto('/');
	}

	async function handleSSOLogin(provider: string) {
		loading = true;
		// Simulate OAuth flow
		await new Promise(resolve => setTimeout(resolve, 1000));
		
		authStore.loginSSO(provider);
		loading = false;
		goto('/');
	}
</script>

<div class="login-container">
	<div class="login-card">
		<div class="login-header">
			<h1>🛡️ AlgoShield</h1>
			<p>Fraud Detection & Anti-Money Laundering</p>
		</div>

		<div class="tabs">
			<button 
				class="tab" 
				class:active={activeTab === 'local'}
				on:click={() => activeTab = 'local'}
			>
				Local Account
			</button>
			<button 
				class="tab" 
				class:active={activeTab === 'sso'}
				on:click={() => activeTab = 'sso'}
			>
				SSO
			</button>
		</div>

		{#if activeTab === 'local'}
			<form on:submit|preventDefault={handleLocalLogin} class="login-form">
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
			<div class="sso-options">
				<p class="sso-description">Choose your SSO provider to continue</p>
				
				<button 
					class="sso-button google"
					on:click={() => handleSSOLogin('google')}
					disabled={loading}
				>
					<svg width="18" height="18" viewBox="0 0 24 24">
						<path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
						<path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
						<path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
						<path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
					</svg>
					Continue with Google
				</button>

				<button 
					class="sso-button microsoft"
					on:click={() => handleSSOLogin('microsoft')}
					disabled={loading}
				>
					<svg width="18" height="18" viewBox="0 0 24 24">
						<path fill="#f25022" d="M1 1h10v10H1z"/>
						<path fill="#00a4ef" d="M13 1h10v10H13z"/>
						<path fill="#7fba00" d="M1 13h10v10H1z"/>
						<path fill="#ffb900" d="M13 13h10v10H13z"/>
					</svg>
					Continue with Microsoft
				</button>

				<button 
					class="sso-button okta"
					on:click={() => handleSSOLogin('okta')}
					disabled={loading}
				>
					<svg width="18" height="18" viewBox="0 0 24 24">
						<circle fill="#007DC1" cx="12" cy="12" r="12"/>
						<circle fill="#FFFFFF" cx="12" cy="12" r="6"/>
					</svg>
					Continue with Okta
				</button>

				{#if loading}
					<p class="loading-text">Authenticating...</p>
				{/if}
			</div>
		{/if}

		<div class="login-footer">
			<p>Demo mode - No real authentication required</p>
		</div>
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
	}

	.form-group input:focus {
		outline: none;
		border-color: var(--primary);
	}

	.form-group input:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.full-width {
		width: 100%;
	}

	.sso-options {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.sso-description {
		text-align: center;
		color: var(--text-secondary);
		margin-bottom: 0.5rem;
	}

	.sso-button {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.75rem;
		padding: 0.875rem 1.5rem;
		border: 2px solid var(--border);
		border-radius: 6px;
		background: white;
		cursor: pointer;
		font-size: 1rem;
		font-weight: 500;
		transition: all 0.2s;
		color: var(--text);
	}

	.sso-button:hover:not(:disabled) {
		border-color: var(--primary);
		background: var(--surface);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	.sso-button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.loading-text {
		text-align: center;
		color: var(--primary);
		font-weight: 500;
	}

	.login-footer {
		margin-top: 2rem;
		padding-top: 1.5rem;
		border-top: 1px solid var(--border);
		text-align: center;
	}

	.login-footer p {
		color: var(--text-secondary);
		font-size: 0.875rem;
	}
</style>
