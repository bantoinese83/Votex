import { authStore } from '$lib/stores/auth';
import { appStore } from '$lib/stores/app';
import { get } from 'svelte/store';

export interface ApiResponse<T = any> {
	success: boolean;
	data?: T;
	error?: string;
	message?: string;
}

export interface ApiError {
	message: string;
	status: number;
	code?: string;
}

class ApiClient {
	private baseURL: string;

	constructor(baseURL: string = '/api') {
		this.baseURL = baseURL;
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<ApiResponse<T>> {
		const url = `${this.baseURL}${endpoint}`;
		
		// Get auth token
		const authState = get(authStore);
		const headers: Record<string, string> = {
			'Content-Type': 'application/json',
			...options.headers as Record<string, string>
		};

		// Add auth header if token exists
		if (authState.token) {
			headers['Authorization'] = `Bearer ${authState.token}`;
		}

		const config: RequestInit = {
			...options,
			headers
		};

		try {
			appStore.setLoading(true);
			
			const response = await fetch(url, config);
			const data = await response.json();

			if (!response.ok) {
				// Handle authentication errors
				if (response.status === 401) {
					authStore.logout();
					appStore.error('Session expired. Please login again.');
				}

				throw new Error(data.error || `HTTP ${response.status}: ${response.statusText}`);
			}

			return data;
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : 'Network error';
			appStore.error(errorMessage);
			throw error;
		} finally {
			appStore.setLoading(false);
		}
	}

	// GET request
	async get<T>(endpoint: string, params?: Record<string, any>): Promise<ApiResponse<T>> {
		const url = new URL(`${this.baseURL}${endpoint}`, window.location.origin);
		
		if (params) {
			Object.entries(params).forEach(([key, value]) => {
				if (value !== undefined && value !== null) {
					url.searchParams.append(key, String(value));
				}
			});
		}

		return this.request<T>(endpoint + url.search);
	}

	// POST request
	async post<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'POST',
			body: data ? JSON.stringify(data) : undefined
		});
	}

	// PUT request
	async put<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'PUT',
			body: data ? JSON.stringify(data) : undefined
		});
	}

	// PATCH request
	async patch<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'PATCH',
			body: data ? JSON.stringify(data) : undefined
		});
	}

	// DELETE request
	async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
		return this.request<T>(endpoint, {
			method: 'DELETE'
		});
	}

	// Upload file
	async upload<T>(endpoint: string, file: File, onProgress?: (progress: number) => void): Promise<ApiResponse<T>> {
		const formData = new FormData();
		formData.append('file', file);

		const authState = get(authStore);
		const headers: Record<string, string> = {};

		if (authState.token) {
			headers['Authorization'] = `Bearer ${authState.token}`;
		}

		const xhr = new XMLHttpRequest();
		
		return new Promise((resolve, reject) => {
			xhr.upload.addEventListener('progress', (event) => {
				if (event.lengthComputable && onProgress) {
					const progress = (event.loaded / event.total) * 100;
					onProgress(progress);
				}
			});

			xhr.addEventListener('load', () => {
				if (xhr.status >= 200 && xhr.status < 300) {
					try {
						const data = JSON.parse(xhr.responseText);
						resolve(data);
					} catch (error) {
						reject(new Error('Invalid JSON response'));
					}
				} else {
					reject(new Error(`HTTP ${xhr.status}: ${xhr.statusText}`));
				}
			});

			xhr.addEventListener('error', () => {
				reject(new Error('Network error'));
			});

			xhr.open('POST', `${this.baseURL}${endpoint}`);
			
			if (authState.token) {
				xhr.setRequestHeader('Authorization', `Bearer ${authState.token}`);
			}

			xhr.send(formData);
		});
	}
}

// Create and export the API client instance
export const apiClient = new ApiClient();

// Auth-specific API methods
export const authApi = {
	login: (username: string, password: string) =>
		apiClient.post<{ token: string; user: any }>('/auth/login', { username, password }),
	
	register: (username: string, password: string) =>
		apiClient.post<{ token: string; user: any }>('/auth/register', { username, password }),
	
	profile: () => apiClient.get<any>('/auth/profile'),
	
	logout: () => apiClient.post('/auth/logout')
};

// Health check
export const healthApi = {
	check: () => apiClient.get<{ status: string; message: string }>('/health')
}; 