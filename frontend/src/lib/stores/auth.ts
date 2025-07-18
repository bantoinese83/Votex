import { writable, derived } from 'svelte/store';
import type { Writable } from 'svelte/store';

export interface User {
	id: string;
	username: string;
}

export interface AuthState {
	user: User | null;
	token: string | null;
	isLoading: boolean;
	error: string | null;
}

const createAuthStore = () => {
	const { subscribe, set, update }: Writable<AuthState> = writable({
		user: null,
		token: null,
		isLoading: false,
		error: null
	});

	// Initialize from localStorage
	const init = () => {
		if (typeof window !== 'undefined') {
			const token = localStorage.getItem('auth_token');
			const userStr = localStorage.getItem('auth_user');
			
			if (token && userStr) {
				try {
					const user = JSON.parse(userStr);
					set({ user, token, isLoading: false, error: null });
				} catch (error) {
					console.error('Failed to parse stored user data:', error);
					logout();
				}
			}
		}
	};

	const login = async (username: string, password: string) => {
		update(state => ({ ...state, isLoading: true, error: null }));

		try {
			const response = await fetch('/api/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ username, password })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Login failed');
			}

			const { token, user } = data.data;

			// Store in localStorage
			if (typeof window !== 'undefined') {
				localStorage.setItem('auth_token', token);
				localStorage.setItem('auth_user', JSON.stringify(user));
			}

			set({ user, token, isLoading: false, error: null });
			return { success: true };
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : 'Login failed';
			set({ user: null, token: null, isLoading: false, error: errorMessage });
			return { success: false, error: errorMessage };
		}
	};

	const register = async (username: string, password: string) => {
		update(state => ({ ...state, isLoading: true, error: null }));

		try {
			const response = await fetch('/api/auth/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ username, password })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.error || 'Registration failed');
			}

			const { token, user } = data.data;

			// Store in localStorage
			if (typeof window !== 'undefined') {
				localStorage.setItem('auth_token', token);
				localStorage.setItem('auth_user', JSON.stringify(user));
			}

			set({ user, token, isLoading: false, error: null });
			return { success: true };
		} catch (error) {
			const errorMessage = error instanceof Error ? error.message : 'Registration failed';
			set({ user: null, token: null, isLoading: false, error: errorMessage });
			return { success: false, error: errorMessage };
		}
	};

	const logout = () => {
		// Clear localStorage
		if (typeof window !== 'undefined') {
			localStorage.removeItem('auth_token');
			localStorage.removeItem('auth_user');
		}

		set({ user: null, token: null, isLoading: false, error: null });
	};

	const refreshProfile = async () => {
		const state = get({ subscribe });
		if (!state.token) return;

		try {
			const response = await fetch('/api/auth/profile', {
				headers: {
					'Authorization': `Bearer ${state.token}`
				}
			});

			if (response.ok) {
				const data = await response.json();
				const user = data.data;
				
				// Update stored user data
				if (typeof window !== 'undefined') {
					localStorage.setItem('auth_user', JSON.stringify(user));
				}

				update(state => ({ ...state, user }));
			} else {
				// Token might be invalid, logout
				logout();
			}
		} catch (error) {
			console.error('Failed to refresh profile:', error);
		}
	};

	const clearError = () => {
		update(state => ({ ...state, error: null }));
	};

	return {
		subscribe,
		init,
		login,
		register,
		logout,
		refreshProfile,
		clearError
	};
};

export const authStore = createAuthStore();

// Derived stores for easier access to specific state
export const isAuthenticated = derived(authStore, ($auth) => !!$auth.user && !!$auth.token);
export const currentUser = derived(authStore, ($auth) => $auth.user);
export const authToken = derived(authStore, ($auth) => $auth.token);
export const authLoading = derived(authStore, ($auth) => $auth.isLoading);
export const authError = derived(authStore, ($auth) => $auth.error);

// Initialize the store when the module is imported
if (typeof window !== 'undefined') {
	authStore.init();
} 