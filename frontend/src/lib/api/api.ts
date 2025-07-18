const BASE_URL = 'http://localhost:8080'; // This should be in an environment variable

export async function getHealth() {
	const response = await fetch(`${BASE_URL}/`);
	if (!response.ok) {
		throw new Error('Failed to fetch health');
	}
	return response.text();
}
