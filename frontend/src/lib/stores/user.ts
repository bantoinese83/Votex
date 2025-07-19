import { writable, derived, type Writable } from 'svelte/store';
import { auth } from './auth.js';

export interface User {
	id: string;
	username: string;
	email?: string;
	avatar?: string;
	preferences?: UserPreferences;
}

export interface UserPreferences {
	theme: 'light' | 'dark' | 'system';
	language: string;
	notifications: boolean;
}

export interface UserState {
	user: User | null;
	loading: boolean;
	error: string | null;
}

// Create the main user store
const createUserStore = () => {
	const { subscribe, set, update }: Writable<UserState> = writable({
		user: null,
		loading: false,
		error: null
	});

	return {
		subscribe,
		set,
		update,
		
		// Actions
		setUser: (user: User) => {
			update(state => ({ ...state, user, error: null }));
		},
		
		setLoading: (loading: boolean) => {
			update(state => ({ ...state, loading }));
		},
		
		setError: (error: string) => {
			update(state => ({ ...state, error, loading: false }));
		},
		
		clearError: () => {
			update(state => ({ ...state, error: null }));
		},
		
		updateUser: (updates: Partial<User>) => {
			update(state => ({
				...state,
				user: state.user ? { ...state.user, ...updates } : null
			}));
		},
		
		updatePreferences: (preferences: Partial<UserPreferences>) => {
			update(state => ({
				...state,
				user: state.user ? {
					...state.user,
					preferences: {
						...state.user.preferences,
						...preferences
					}
				} : null
			}));
		},
		
		logout: () => {
			set({ user: null, loading: false, error: null });
		},
		
		// Async actions
		fetchProfile: async () => {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				const response = await fetch('/api/auth/profile', {
					headers: {
						'Authorization': `Bearer ${localStorage.getItem('token')}`
					}
				});
				
				if (!response.ok) {
					throw new Error('Failed to fetch profile');
				}
				
				const data = await response.json();
				if (data.success) {
					update(state => ({ 
						...state, 
						user: data.data, 
						loading: false 
					}));
				} else {
					throw new Error(data.error || 'Failed to fetch profile');
				}
			} catch (error) {
				update(state => ({ 
					...state, 
					error: error instanceof Error ? error.message : 'Unknown error',
					loading: false 
				}));
			}
		},
		
		updateProfile: async (updates: Partial<User>) => {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				const response = await fetch('/api/auth/profile', {
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json',
						'Authorization': `Bearer ${localStorage.getItem('token')}`
					},
					body: JSON.stringify(updates)
				});
				
				if (!response.ok) {
					throw new Error('Failed to update profile');
				}
				
				const data = await response.json();
				if (data.success) {
					update(state => ({
						...state,
						user: { ...state.user, ...data.data },
						loading: false
					}));
				} else {
					throw new Error(data.error || 'Failed to update profile');
				}
			} catch (error) {
				update(state => ({
					...state,
					error: error instanceof Error ? error.message : 'Unknown error',
					loading: false
				}));
			}
		}
	};
};

export const user = createUserStore();

// Derived stores for specific user data
export const currentUser = derived(user, $user => $user.user);
export const isLoggedIn = derived(currentUser, $user => $user !== null);
export const isLoading = derived(user, $user => $user.loading);
export const userError = derived(user, $user => $user.error);

// Derived store for user preferences with defaults
export const userPreferences = derived(currentUser, $user => ({
	theme: $user?.preferences?.theme || 'system',
	language: $user?.preferences?.language || 'en',
	notifications: $user?.preferences?.notifications ?? true
}));

// Derived store for user display name
export const displayName = derived(currentUser, $user => {
	if (!$user) return '';
	return $user.username || $user.email || 'Unknown User';
});

// Derived store for user avatar
export const userAvatar = derived(currentUser, $user => {
	if (!$user?.avatar) {
		// Generate default avatar based on username
		const username = $user?.username || 'U';
		const initials = username.substring(0, 2).toUpperCase();
		return `https://ui-avatars.com/api/?name=${encodeURIComponent(initials)}&background=random&size=40`;
	}
	return $user.avatar;
});

// Auto-sync with auth store
auth.subscribe($auth => {
	if (!$auth.isAuthenticated) {
		user.logout();
	} else if ($auth.isAuthenticated && !currentUser) {
		// Fetch user profile when authenticated
		user.fetchProfile();
	}
});

// Persist user preferences to localStorage
userPreferences.subscribe($preferences => {
	if (typeof window !== 'undefined') {
		localStorage.setItem('userPreferences', JSON.stringify($preferences));
	}
});

// Load user preferences from localStorage on init
if (typeof window !== 'undefined') {
	const savedPreferences = localStorage.getItem('userPreferences');
	if (savedPreferences) {
		try {
			const preferences = JSON.parse(savedPreferences);
			user.updatePreferences(preferences);
		} catch (error) {
			console.warn('Failed to load user preferences:', error);
		}
	}
} 