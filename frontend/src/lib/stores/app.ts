import { writable, derived } from 'svelte/store';
import type { Writable } from 'svelte/store';

export interface AppState {
	theme: 'light' | 'dark';
	sidebarOpen: boolean;
	notifications: Notification[];
	isLoading: boolean;
}

export interface Notification {
	id: string;
	type: 'success' | 'error' | 'warning' | 'info';
	message: string;
	duration?: number;
}

const createAppStore = () => {
	const { subscribe, set, update }: Writable<AppState> = writable({
		theme: 'light',
		sidebarOpen: false,
		notifications: [],
		isLoading: false
	});

	// Initialize from localStorage
	const init = () => {
		if (typeof window !== 'undefined') {
			const savedTheme = localStorage.getItem('theme') as 'light' | 'dark';
			if (savedTheme) {
				update(state => ({ ...state, theme: savedTheme }));
			}
		}
	};

	const toggleTheme = () => {
		update(state => {
			const newTheme = state.theme === 'light' ? 'dark' : 'light';
			
			// Save to localStorage
			if (typeof window !== 'undefined') {
				localStorage.setItem('theme', newTheme);
			}
			
			return { ...state, theme: newTheme };
		});
	};

	const setTheme = (theme: 'light' | 'dark') => {
		update(state => {
			// Save to localStorage
			if (typeof window !== 'undefined') {
				localStorage.setItem('theme', theme);
			}
			
			return { ...state, theme };
		});
	};

	const toggleSidebar = () => {
		update(state => ({ ...state, sidebarOpen: !state.sidebarOpen }));
	};

	const setSidebarOpen = (open: boolean) => {
		update(state => ({ ...state, sidebarOpen: open }));
	};

	const addNotification = (notification: Omit<Notification, 'id'>) => {
		const id = Date.now().toString();
		const newNotification = { ...notification, id };
		
		update(state => ({
			...state,
			notifications: [...state.notifications, newNotification]
		}));

		// Auto-remove notification after duration
		if (notification.duration !== 0) {
			setTimeout(() => {
				removeNotification(id);
			}, notification.duration || 5000);
		}

		return id;
	};

	const removeNotification = (id: string) => {
		update(state => ({
			...state,
			notifications: state.notifications.filter(n => n.id !== id)
		}));
	};

	const clearNotifications = () => {
		update(state => ({ ...state, notifications: [] }));
	};

	const setLoading = (loading: boolean) => {
		update(state => ({ ...state, isLoading: loading }));
	};

	// Convenience methods for common notification types
	const success = (message: string, duration?: number) => {
		return addNotification({ type: 'success', message, duration });
	};

	const error = (message: string, duration?: number) => {
		return addNotification({ type: 'error', message, duration });
	};

	const warning = (message: string, duration?: number) => {
		return addNotification({ type: 'warning', message, duration });
	};

	const info = (message: string, duration?: number) => {
		return addNotification({ type: 'info', message, duration });
	};

	return {
		subscribe,
		init,
		toggleTheme,
		setTheme,
		toggleSidebar,
		setSidebarOpen,
		addNotification,
		removeNotification,
		clearNotifications,
		setLoading,
		success,
		error,
		warning,
		info
	};
};

export const appStore = createAppStore();

// Derived stores
export const theme = derived(appStore, ($app) => $app.theme);
export const sidebarOpen = derived(appStore, ($app) => $app.sidebarOpen);
export const notifications = derived(appStore, ($app) => $app.notifications);
export const isLoading = derived(appStore, ($app) => $app.isLoading);

// Initialize the store when the module is imported
if (typeof window !== 'undefined') {
	appStore.init();
} 