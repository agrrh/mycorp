<script>
	import favicon from '$lib/assets/favicon.svg';
	import { resolve } from '$app/paths';
	import { apiUrl, token } from '$lib/api.js';

	let { children } = $props();

	let settingsOpen = $state(false);
	let draftUrl = $state('http://127.0.0.1:8080');
	let draftToken = $state('');

	apiUrl.subscribe((v) => {
		draftUrl = v;
	});
	token.subscribe((v) => {
		draftToken = v;
	});

	function openSettings() {
		settingsOpen = true;
	}

	function saveSettings() {
		apiUrl.set(draftUrl.trim());
		token.set(draftToken.trim());
		settingsOpen = false;
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>MyCorp</title>
</svelte:head>

<nav class="navbar is-dark" aria-label="main navigation">
	<div class="navbar-brand">
		<a class="navbar-item has-text-weight-bold" href={resolve('/')}>MyCorp</a>
	</div>
	<div class="navbar-menu">
		<div class="navbar-end">
			<span class="navbar-item is-size-7 has-text-grey-light">{draftUrl}</span>
			<div class="navbar-item">
				<button class="button is-small is-light" onclick={openSettings}>Settings</button>
			</div>
		</div>
	</div>
</nav>

<main class="section">
	<div class="container is-max-desktop">{@render children()}</div>
</main>

{#if settingsOpen}
	<div class="modal is-active">
		<div
			class="modal-background"
			role="button"
			tabindex="0"
			onclick={() => (settingsOpen = false)}
			onkeydown={(e) => e.key === 'Enter' && (settingsOpen = false)}
		></div>
		<div class="modal-card">
			<header class="modal-card-head">
				<p class="modal-card-title">Settings</p>
				<button class="delete" aria-label="close" onclick={() => (settingsOpen = false)}></button>
			</header>
			<section class="modal-card-body">
				<div class="field">
					<label class="label" for="api-url">API URL</label>
					<div class="control">
						<input
							id="api-url"
							class="input"
							type="text"
							bind:value={draftUrl}
							placeholder="http://127.0.0.1:8080"
						/>
					</div>
				</div>
				<div class="field">
					<label class="label" for="token">Token</label>
					<div class="control">
						<input
							id="token"
							class="input"
							type="password"
							bind:value={draftToken}
							placeholder="X-Token value"
						/>
					</div>
				</div>
			</section>
			<footer class="modal-card-foot">
				<button class="button is-primary" onclick={saveSettings}>Save</button>
				<button class="button" onclick={() => (settingsOpen = false)}>Cancel</button>
			</footer>
		</div>
	</div>
{/if}
