import { writable } from 'svelte/store';

/** @type {import('svelte/store').Writable<string>} */
export const apiUrl = writable(
	typeof sessionStorage !== 'undefined'
		? (sessionStorage.getItem('api_url') ?? 'http://127.0.0.1:8080')
		: 'http://127.0.0.1:8080'
);

/** @type {import('svelte/store').Writable<string>} */
export const token = writable(
	typeof sessionStorage !== 'undefined' ? (sessionStorage.getItem('token') ?? '') : ''
);

apiUrl.subscribe((v) => {
	if (typeof sessionStorage !== 'undefined') sessionStorage.setItem('api_url', v);
});
token.subscribe((v) => {
	if (typeof sessionStorage !== 'undefined') sessionStorage.setItem('token', v);
});

let tokenValue = '';
let urlValue = 'http://127.0.0.1:8080';

token.subscribe((v) => {
	tokenValue = v;
});
apiUrl.subscribe((v) => {
	urlValue = v;
});

/** @returns {Record<string, string>} */
function baseHeaders() {
	/** @type {Record<string, string>} */
	const headers = { 'Content-Type': 'application/json' };
	if (tokenValue) {
		headers['X-Token'] = tokenValue;
	}
	return headers;
}

/**
 * @param {string} path
 * @param {{ method?: string; body?: string; headers?: Record<string, string> }} options
 * @returns {Promise<any>}
 */
async function request(path, options = {}) {
	const res = await fetch(`${urlValue}${path}`, {
		...options,
		headers: { ...baseHeaders(), ...(options.headers || {}) }
	});

	let body = null;
	const text = await res.text();
	if (text) {
		try {
			body = JSON.parse(text);
		} catch {
			body = text;
		}
	}

	if (!res.ok) {
		const err = /** @type {Error & { status?: number; body?: any }} */ (
			new Error(body && typeof body === 'string' ? body : `Request failed (${res.status})`)
		);
		err.status = res.status;
		err.body = body;
		throw err;
	}

	return body;
}

/** @returns {Promise<string[]>} */
export function listScenarios() {
	return request('/scenarios/');
}

/**
 * @param {string} namespace
 * @param {string} name
 * @returns {Promise<any>}
 */
export function getScenario(namespace, name) {
	return request(`/scenarios/${namespace}/${name}/_cli`);
}

/**
 * @param {string} namespace
 * @param {string} name
 * @param {Record<string, any>} inputs
 * @returns {Promise<any>}
 */
export function runScenario(namespace, name, inputs) {
	return request(`/scenarios/${namespace}/${name}`, {
		method: 'POST',
		body: JSON.stringify(inputs)
	});
}
