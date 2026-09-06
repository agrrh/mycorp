<script>
	import { onMount } from 'svelte';
	import { resolve } from '$app/paths';
	import { getScenario, runScenario } from '$lib/api.js';

	let { params } = $props();

	let namespace = $derived(params.namespace);
	let name = $derived(params.name);

	/** @type {any} */
	let scenario = $state(null);
	let loading = $state(true);
	let error = $state('');
	let running = $state(false);
	let runError = $state('');
	/** @type {any} */
	let result = $state(null);

	/** @type {Record<string, any>} */
	let values = $state({});
	/** @type {Record<string, string>} */
	let validation = $state({});

	onMount(async () => {
		await load();
	});

	async function load() {
		loading = true;
		error = '';
		try {
			const cli = await getScenario(namespace, name);
			scenario = cli;

			/** @type {Record<string, any>} */
			const init = {};
			for (const input of cli.spec?.inputs || []) {
				if (input.default !== undefined && input.default !== null) {
					init[input.name] = input.default;
				} else if (input.type === 'enum') {
					init[input.name] = input.values?.[0];
				} else {
					init[input.name] = '';
				}
			}
			values = init;
		} catch (e) {
			error = /** @type {Error} */ (e).message;
		} finally {
			loading = false;
		}
	}

	/** @param {any} input @param {any} v */
	function setValue(input, v) {
		values = { ...values, [input.name]: v };
		validation = { ...validation, [input.name]: '' };
	}

	/** @param {any} input @param {any} raw */
	function validateInput(input, raw) {
		const value = raw;

		if (value === '' || value === null || value === undefined) {
			if (input.default !== undefined || input.optional) return '';
			return 'This field is required';
		}

		if (input.regexp) {
			try {
				const re = new RegExp(input.regexp);
				if (!re.test(String(value))) {
					return `Must match ${input.regexp}`;
				}
			} catch {
				// ignore invalid regexp
			}
		}

		if (input.type === 'int' && !Number.isInteger(Number(value))) {
			return 'Must be an integer';
		}
		if (input.type === 'float' && Number.isNaN(Number(value))) {
			return 'Must be a number';
		}

		return '';
	}

	async function submit() {
		/** @type {Record<string, string>} */
		const hasError = {};
		/** @type {Record<string, any>} */
		const payload = {};
		let valid = true;

		for (const input of scenario.spec?.inputs || []) {
			const raw = values[input.name];
			const msg = validateInput(input, raw);
			if (msg) {
				hasError[input.name] = msg;
				valid = false;
				continue;
			}

			if (raw === '' || raw === null || raw === undefined) {
				if (input.default !== undefined) {
					payload[input.name] = input.default;
				}
				continue;
			}

			if (input.type === 'int') {
				payload[input.name] = parseInt(raw, 10);
			} else if (input.type === 'float') {
				payload[input.name] = parseFloat(raw);
			} else if (input.type === 'bool') {
				payload[input.name] = raw === true || raw === 'true' || raw === 1 || raw === '1';
			} else {
				payload[input.name] = raw;
			}
		}

		validation = hasError;
		if (!valid) return;

		running = true;
		runError = '';
		result = null;
		try {
			const res = await runScenario(namespace, name, payload);
			result = res;
		} catch (e) {
			const err = /** @type {Error & { body?: any }} */ (e);
			runError = err.message;
			if (err.body && typeof err.body === 'object') {
				runError = err.body.output || runError;
			}
		} finally {
			running = false;
		}
	}

	/** @param {any} input */
	function typeOf(input) {
		return input.type || (typeof input.default === 'number' ? 'number' : 'string');
	}

	/** @param {any} input */
	function isRequired(input) {
		return !(input.default !== undefined || input.optional);
	}
</script>

<svelte:head>
	<title>MyCorp — {namespace}/{name}</title>
</svelte:head>

<div class="columns">
	<div class="column is-7">
		<nav class="breadcrumb" aria-label="breadcrumbs">
			<ul>
				<li><a href={resolve('/')}>Scenarios</a></li>
				<li><span class="px-1">{namespace}</span></li>
				<li class="is-active"><span aria-current="page" class="px-1">{name}</span></li>
			</ul>
		</nav>

		{#if loading}
			<div class="has-text-centered my-6">
				<span class="icon is-large has-text-grey">
					<i class="fas fa-spinner fa-pulse fa-2x"></i>
				</span>
			</div>
		{:else if error}
			<div class="notification is-danger">
				<button class="delete" aria-label="dismiss" onclick={() => (error = '')}></button>
				<strong>Failed to load scenario:</strong>
				{error}
			</div>
		{:else if scenario}
			<h1 class="title">{name}</h1>
			<p class="subtitle">
				<span class="tag is-medium">{scenario.metadata.namespace}/{scenario.metadata.name}</span>
			</p>

			{#if runError}
				<div class="notification is-danger">
					<button class="delete" aria-label="dismiss" onclick={() => (runError = '')}></button>
					{runError}
				</div>
			{/if}

			{#if scenario.spec?.inputs?.length}
				<form
					onsubmit={(e) => {
						e.preventDefault();
						submit();
					}}
				>
					{#each scenario.spec.inputs as input (input.name)}
						<div class="field">
							<label class="label" for={input.name}>
								{input.name}
								{#if !isRequired(input)}
									<span class="has-text-grey-light is-size-7">(optional)</span>
								{/if}
							</label>
							{#if input.type === 'enum'}
								<div class="control">
									<div class="select">
										<select
											id={input.name}
											value={values[input.name]}
											onchange={(e) => setValue(input, e.currentTarget.value)}
										>
											{#each input.values || [] as v (v)}
												<option value={v}>{v}</option>
											{/each}
										</select>
									</div>
								</div>
							{:else if typeOf(input) === 'bool' || input.default === true || input.default === false}
								<div class="control">
									<label class="checkbox">
										<input
											type="checkbox"
											checked={!!values[input.name]}
											onchange={(e) => setValue(input, e.currentTarget.checked)}
										/>
										<span class="ml-1">true</span>
									</label>
								</div>
							{:else}
								<div class="control">
									<input
										class="input"
										id={input.name}
										type={typeOf(input) === 'number' ? 'number' : 'text'}
										value={values[input.name]}
										placeholder={input.description || ''}
										oninput={(e) => setValue(input, e.currentTarget.value)}
									/>
								</div>
							{/if}
							{#if validation[input.name]}
								<p class="help is-danger">{validation[input.name]}</p>
							{/if}
							{#if input.description}
								<p class="help">{input.description}</p>
							{/if}
						</div>
					{/each}
					<div class="field">
						<div class="control">
							<button class="button is-primary" type="submit" disabled={running}>
								{#if running}
									<span class="icon"><i class="fas fa-spinner fa-pulse"></i></span>
									<span>Running…</span>
								{:else}
									<span class="icon"><i class="fas fa-play"></i></span>
									<span>Run scenario</span>
								{/if}
							</button>
						</div>
					</div>
				</form>
			{:else}
				<div class="field">
					<div class="control">
						<button class="button is-primary" onclick={submit} disabled={running}>
							{#if running}
								<span class="icon"><i class="fas fa-spinner fa-pulse"></i></span>
								<span>Running…</span>
							{:else}
								<span class="icon"><i class="fas fa-play"></i></span>
								<span>Run scenario</span>
							{/if}
						</button>
					</div>
				</div>
			{/if}
		{/if}
	</div>

	<div class="column is-5">
		{#if result}
			<div class="box">
				<p class="title is-5">Result</p>
				{#if result.success}
					<div class="notification is-success">
						<strong>Success</strong>
					</div>
				{:else}
					<div class="notification is-danger">
						<strong>{result.status === 'error' ? 'Error' : 'Failed'}</strong>
					</div>
				{/if}
				{#if result.output}
					<pre class="mt-3">{result.output}</pre>
				{/if}
				{#if result.results}
					<details class="mt-3">
						<summary class="is-size-7 has-text-grey">Raw results</summary>
						<pre class="mt-2"><code>{result.results}</code></pre>
					</details>
				{/if}
			</div>
		{/if}
	</div>
</div>
