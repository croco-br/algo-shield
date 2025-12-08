<script lang="ts">
	import '../app.css';
	import { authStore } from '$lib/stores/auth';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let user: any = null;
	let showUserMenu = false;

	authStore.subscribe(value => {
		user = value;
	});

	// Check authentication for protected routes
	$: {
		if ($page.url.pathname !== '/login' && !user) {
			goto('/login');
		}
	}

	function logout() {
		authStore.logout();
		goto('/login');
	}

	function toggleUserMenu() {
		showUserMenu = !showUserMenu;
	}

	// Close menu when clicking outside
	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.user-menu-container')) {
			showUserMenu = false;
		}
	}

	onMount(() => {
		document.addEventListener('click', handleClickOutside);
		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});
</script>

<div class="app">
	{#if user}
		<header>
			<div class="container">
				<nav>
					<div class="nav-brand">
						<h1>🛡️ AlgoShield</h1>
						<p>Fraud Detection & Anti-Money Laundering</p>
					</div>

					<div class="nav-links">
						<a 
							href="/" 
							class="nav-link"
							class:active={$page.url.pathname === '/'}
						>
							Rules
						</a>
						<a 
							href="/synthetic-test" 
							class="nav-link"
							class:active={$page.url.pathname === '/synthetic-test'}
						>
							Synthetic Test
						</a>
					</div>

					<div class="user-menu-container">
						<button class="user-button" on:click={toggleUserMenu}>
							<div class="user-avatar">
								{user.name.charAt(0).toUpperCase()}
							</div>
							<div class="user-info">
								<div class="user-name">{user.name}</div>
								<div class="user-email">{user.email}</div>
							</div>
							<svg class="chevron" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
								<polyline points="6 9 12 15 18 9"></polyline>
							</svg>
						</button>

						{#if showUserMenu}
							<div class="user-menu">
								<div class="user-menu-header">
									<div class="user-name">{user.name}</div>
									<div class="user-email">{user.email}</div>
									<div class="auth-type-badge">
										{user.authType === 'sso' ? 'SSO' : 'Local'}
									</div>
								</div>
								<button class="menu-item logout" on:click={logout}>
									<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
										<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
										<polyline points="16 17 21 12 16 7"></polyline>
										<line x1="21" y1="12" x2="9" y2="12"></line>
									</svg>
									Logout
								</button>
							</div>
						{/if}
					</div>
				</nav>
			</div>
		</header>
	{/if}

	<main class:no-header={!user}>
		<slot />
	</main>
</div>

<style>
	.app {
		min-height: 100vh;
	}

	header {
		background: var(--surface);
		border-bottom: 1px solid var(--border);
		box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
	}

	nav {
		padding: 1.5rem 0;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 2rem;
	}

	.nav-brand h1 {
		font-size: 1.75rem;
		color: var(--text);
		margin-bottom: 0.25rem;
	}

	.nav-brand p {
		color: var(--text-secondary);
		font-size: 0.875rem;
	}

	.nav-links {
		display: flex;
		gap: 0.5rem;
		flex: 1;
		justify-content: center;
	}

	.nav-link {
		padding: 0.75rem 1.5rem;
		border-radius: 6px;
		text-decoration: none;
		color: var(--text-secondary);
		font-weight: 500;
		transition: all 0.2s;
	}

	.nav-link:hover {
		color: var(--text);
		background: rgba(0, 0, 0, 0.05);
	}

	.nav-link.active {
		color: var(--primary);
		background: rgba(102, 126, 234, 0.1);
	}

	.user-menu-container {
		position: relative;
	}

	.user-button {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.5rem;
		border: 1px solid var(--border);
		border-radius: 8px;
		background: white;
		cursor: pointer;
		transition: all 0.2s;
	}

	.user-button:hover {
		background: var(--surface);
	}

	.user-avatar {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		font-size: 1.125rem;
	}

	.user-info {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
	}

	.user-name {
		font-weight: 500;
		color: var(--text);
		font-size: 0.875rem;
	}

	.user-email {
		font-size: 0.75rem;
		color: var(--text-secondary);
	}

	.chevron {
		color: var(--text-secondary);
	}

	.user-menu {
		position: absolute;
		top: calc(100% + 0.5rem);
		right: 0;
		background: white;
		border: 1px solid var(--border);
		border-radius: 8px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		min-width: 240px;
		z-index: 100;
	}

	.user-menu-header {
		padding: 1rem;
		border-bottom: 1px solid var(--border);
	}

	.user-menu-header .user-name {
		font-weight: 600;
		font-size: 1rem;
		margin-bottom: 0.25rem;
	}

	.user-menu-header .user-email {
		font-size: 0.875rem;
		color: var(--text-secondary);
		margin-bottom: 0.5rem;
	}

	.auth-type-badge {
		display: inline-block;
		padding: 0.25rem 0.75rem;
		background: var(--surface);
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
		color: var(--text-secondary);
	}

	.menu-item {
		width: 100%;
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		border: none;
		background: none;
		cursor: pointer;
		font-size: 0.875rem;
		color: var(--text);
		text-align: left;
		transition: background 0.2s;
	}

	.menu-item:hover {
		background: var(--surface);
	}

	.menu-item.logout {
		color: var(--danger);
		border-top: 1px solid var(--border);
	}

	main {
		padding: 2rem 0;
	}

	main.no-header {
		padding: 0;
	}
</style>

