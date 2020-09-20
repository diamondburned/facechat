export async function fetchx(url, opts) {
	let resp = await fetch(url, opts)
	if (!resp.ok) {
		let msg = `Unexpected status code ${resp.status}`

		let body = await resp.json()
		if (body.error) {
			msg += `: ${body.error}`
		}

		throw msg
	}

	return resp
}
