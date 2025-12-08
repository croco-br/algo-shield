<script lang="ts">
	import '../app.css';
</script>

<div class="app">
	<header>
		<div class="container">
			<nav>
				<h1>🛡️ AlgoShield</h1>
				<p>Fraud Detection & Anti-Money Laundering</p>
			</nav>
		</div>
	</header>

	<main>
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
	}

	h1 {
		font-size: 1.75rem;
		color: var(--text);
		margin-bottom: 0.25rem;
	}

	p {
		color: var(--text-secondary);
		font-size: 0.875rem;
	}

	main {
		padding: 2rem 0;
	}
</style>

