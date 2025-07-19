import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { userPreferences } from './user.js';

export type Theme = 'light' | 'dark' | 'system';

interface ThemeState {
	theme: Theme;
	actualTheme: 'light' | 'dark';
}

// Create the theme store
const createThemeStore = () => {
	// Get initial theme from localStorage or default to 'system'
	const getInitialTheme = (): Theme => {
		if (!browser) return 'system';
		
		const saved = localStorage.getItem('theme');
		if (saved && ['light', 'dark', 'system'].includes(saved)) {
			return saved as Theme;
		}
		
		return 'system';
	};

	const { subscribe, set, update } = writable<ThemeState>({
		theme: getInitialTheme(),
		actualTheme: 'light'
	});

	// Function to get the actual theme (light/dark) based on system preference
	const getActualTheme = (theme: Theme): 'light' | 'dark' => {
		if (theme === 'system' && browser) {
			return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
		}
		return theme === 'dark' ? 'dark' : 'light';
	};

	// Function to apply theme to document
	const applyTheme = (actualTheme: 'light' | 'dark') => {
		if (!browser) return;
		
		const root = document.documentElement;
		root.classList.remove('light', 'dark');
		root.classList.add(actualTheme);
		root.setAttribute('data-theme', actualTheme);
		
		// Update meta theme-color
		const metaThemeColor = document.querySelector('meta[name="theme-color"]');
		if (metaThemeColor) {
			metaThemeColor.setAttribute('content', actualTheme === 'dark' ? '#1a1a1a' : '#ffffff');
		}
	};

	// Listen for system theme changes
	let mediaQuery: MediaQueryList | null = null;
	if (browser) {
		mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
		mediaQuery.addEventListener('change', () => {
			update(state => {
				if (state.theme === 'system') {
					const actualTheme = getActualTheme('system');
					applyTheme(actualTheme);
					return { ...state, actualTheme };
				}
				return state;
			});
		});
	}

	return {
		subscribe,
		
		setTheme: (theme: Theme) => {
			if (browser) {
				localStorage.setItem('theme', theme);
			}
			
			const actualTheme = getActualTheme(theme);
			applyTheme(actualTheme);
			
			set({ theme, actualTheme });
		},
		
		toggle: () => {
			update(state => {
				const newTheme = state.theme === 'light' ? 'dark' : 'light';
				const actualTheme = getActualTheme(newTheme);
				
				if (browser) {
					localStorage.setItem('theme', newTheme);
				}
				
				applyTheme(actualTheme);
				return { theme: newTheme, actualTheme };
			});
		},
		
		// Initialize theme on first load
		init: () => {
			if (browser) {
				update(state => {
					const actualTheme = getActualTheme(state.theme);
					applyTheme(actualTheme);
					return { ...state, actualTheme };
				});
			}
		}
	};
};

export const theme = createThemeStore();

// Derived store for the actual theme (light/dark)
export const actualTheme = derived(theme, $theme => $theme.actualTheme);

// Derived store for CSS classes
export const themeClasses = derived(actualTheme, $actualTheme => ({
	'data-theme': $actualTheme,
	'class': $actualTheme
}));

// Auto-sync with user preferences
userPreferences.subscribe($preferences => {
	if ($preferences.theme) {
		theme.setTheme($preferences.theme);
	}
});

// Initialize theme on app start
if (browser) {
	theme.init();
}

// Export theme utilities
export const isDark = derived(actualTheme, $actualTheme => $actualTheme === 'dark');
export const isLight = derived(actualTheme, $actualTheme => $actualTheme === 'light');

// Theme color utilities
export const themeColors = derived(actualTheme, $actualTheme => ({
	primary: $actualTheme === 'dark' ? '#3b82f6' : '#2563eb',
	background: $actualTheme === 'dark' ? '#1a1a1a' : '#ffffff',
	surface: $actualTheme === 'dark' ? '#2a2a2a' : '#f8fafc',
	text: $actualTheme === 'dark' ? '#ffffff' : '#1a1a1a',
	textSecondary: $actualTheme === 'dark' ? '#a1a1aa' : '#6b7280',
	border: $actualTheme === 'dark' ? '#3a3a3a' : '#e5e7eb'
})); 