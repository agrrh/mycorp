<script>
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { listScenarios } from '$lib/api.js';

	/** @type {Array<{namespace: string, items: Array<{full: string, namespace: string, name: string}>}>} */
	let scenarios = $state([]);
	let loading = $state(true);
	let error = $state('');
	/** @type {Array<{full: string, namespace: string, name: string}>} */
	let parsedScenarios = $state([]);

	onMount(async () => {
		await load();
	});

	async function load() {
		loading = true;
		error = '';
		try {
			const list = await listScenarios();
			parsedScenarios = list.map((full) => {
				const [namespace, name] = full.split('/');
				return { full, namespace, name };
			});
			/** @type {Record<string, Array<{full: string, namespace: string, name: string}>>} */
			const groups = {};
			for (const sc of parsedScenarios) {
				(groups[sc.namespace] ||= []).push(sc);
			}
			scenarios = Object.entries(groups)
				.map(([namespace, items]) => ({
					namespace,
					items: items.sort((a, b) => a.name.localeCompare(b.name))
				}))
				.sort((a, b) => a.namespace.localeCompare(b.namespace));
		} catch (e) {
			error = /** @type {Error} */ (e).message;
		} finally {
			loading = false;
		}
	}

	function retry() {
		load();
	}
</script>

<div class="columns is-centered">
	<div class="column is-10">
		<div class="level">
			<div class="level-left">
				<div>
					<h1 class="title">Scenarios</h1>
					<p class="subtitle">Select a scenario to run it</p>
				</div>
			</div>
			<div class="level-right">
				<button class="button" onclick={retry} aria-label="Refresh">
					<span class="icon"><i class="fas fa-rotate"></i></span>
				</button>
			</div>
		</div>

		{#if loading}
			<div class="has-text-centered">
				<span class="icon is-large has-text-grey">
					<i class="fas fa-spinner fa-pulse fa-2x"></i>
				</span>
			</div>
		{:else if error}
			<div class="notification is-danger">
				<button class="delete" aria-label="dismiss" onclick={() => (error = '')}></button>
				<strong>Failed to load scenarios:</strong>
				{error}
				<div class="mt-3">
					<button class="button is-small is-link" onclick={retry}>Try again</button>
				</div>
			</div>
		{:else if scenarios.length === 0}
			<div class="notification is-info">
				No scenarios available. Check your API URL and token in settings.
			</div>
		{:else}
			{#each scenarios as group (group.namespace)}
				<h2 class="title is-4 mt-5">
					<span class="icon has-text-info"><i class="fas fa-folder"></i></span>
					{group.namespace}
				</h2>
				<div class="columns is-multiline">
					{#each group.items as sc (sc.full)}
						<div class="column is-one-third">
							<a class="box card" href={resolve(`/${sc.namespace}/${sc.name}`)}>
								<div class="card-content">
									<p class="title is-5">{sc.name}</p>
									<p class="is-size-7 has-text-grey">
										{sc.namespace}/{sc.name}
									</p>
								</div>
							</a>
						</div>
					{/each}
				</div>
			{/each}
		{/if}
	</div>
</div>
