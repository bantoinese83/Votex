import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

export type NotificationType = 'success' | 'error' | 'warning' | 'info';

export interface Notification {
	id: string;
	type: NotificationType;
	title: string;
	message: string;
	duration?: number; // Auto-dismiss after duration (ms), 0 = no auto-dismiss
	persistent?: boolean; // Don't auto-dismiss
	actions?: NotificationAction[];
	createdAt: Date;
}

export interface NotificationAction {
	label: string;
	action: () => void;
	variant?: 'primary' | 'secondary' | 'danger';
}

interface NotificationState {
	notifications: Notification[];
	position: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left' | 'top-center' | 'bottom-center';
	maxNotifications: number;
}

// Create the notification store
const createNotificationStore = () => {
	const { subscribe, set, update } = writable<NotificationState>({
		notifications: [],
		position: 'top-right',
		maxNotifications: 5
	});

	// Generate unique ID for notifications
	const generateId = (): string => {
		return `notification-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
	};

	// Auto-dismiss notifications
	const autoDismiss = (id: string, duration: number) => {
		if (duration > 0) {
			setTimeout(() => {
				dismiss(id);
			}, duration);
		}
	};

	return {
		subscribe,
		
		// Add a new notification
		add: (notification: Omit<Notification, 'id' | 'createdAt'>) => {
			const id = generateId();
			const newNotification: Notification = {
				...notification,
				id,
				createdAt: new Date(),
				duration: notification.duration ?? 5000 // Default 5 seconds
			};

			update(state => {
				const notifications = [...state.notifications, newNotification];
				
				// Limit number of notifications
				if (notifications.length > state.maxNotifications) {
					notifications.splice(0, notifications.length - state.maxNotifications);
				}
				
				return { ...state, notifications };
			});

			// Auto-dismiss if duration is set
			if (newNotification.duration && newNotification.duration > 0 && !newNotification.persistent) {
				autoDismiss(id, newNotification.duration);
			}

			return id;
		},

		// Success notification
		success: (title: string, message?: string, options?: Partial<Notification>) => {
			return this.add({
				type: 'success',
				title,
				message: message || '',
				...options
			});
		},

		// Error notification
		error: (title: string, message?: string, options?: Partial<Notification>) => {
			return this.add({
				type: 'error',
				title,
				message: message || '',
				persistent: true, // Errors are persistent by default
				...options
			});
		},

		// Warning notification
		warning: (title: string, message?: string, options?: Partial<Notification>) => {
			return this.add({
				type: 'warning',
				title,
				message: message || '',
				...options
			});
		},

		// Info notification
		info: (title: string, message?: string, options?: Partial<Notification>) => {
			return this.add({
				type: 'info',
				title,
				message: message || '',
				...options
			});
		},

		// Dismiss a specific notification
		dismiss: (id: string) => {
			update(state => ({
				...state,
				notifications: state.notifications.filter(n => n.id !== id)
			}));
		},

		// Dismiss all notifications
		dismissAll: () => {
			update(state => ({
				...state,
				notifications: []
			}));
		},

		// Update notification position
		setPosition: (position: NotificationState['position']) => {
			update(state => ({ ...state, position }));
		},

		// Set max notifications
		setMaxNotifications: (max: number) => {
			update(state => ({ ...state, maxNotifications: max }));
		},

		// Update a notification
		update: (id: string, updates: Partial<Notification>) => {
			update(state => ({
				...state,
				notifications: state.notifications.map(n => 
					n.id === id ? { ...n, ...updates } : n
				)
			}));
		}
	};
};

export const notifications = createNotificationStore();

// Derived stores
export const currentNotifications = derived(notifications, $notifications => $notifications.notifications);
export const notificationCount = derived(currentNotifications, $notifications => $notifications.length);
export const hasNotifications = derived(notificationCount, $count => $count > 0);

// Notification position
export const notificationPosition = derived(notifications, $notifications => $notifications.position);

// Filtered notifications by type
export const successNotifications = derived(currentNotifications, $notifications => 
	$notifications.filter(n => n.type === 'success')
);

export const errorNotifications = derived(currentNotifications, $notifications => 
	$notifications.filter(n => n.type === 'error')
);

export const warningNotifications = derived(currentNotifications, $notifications => 
	$notifications.filter(n => n.type === 'warning')
);

export const infoNotifications = derived(currentNotifications, $notifications => 
	$notifications.filter(n => n.type === 'info')
);

// Notification utilities
export const getNotificationIcon = (type: NotificationType): string => {
	switch (type) {
		case 'success':
			return '✓';
		case 'error':
			return '✕';
		case 'warning':
			return '⚠';
		case 'info':
			return 'ℹ';
		default:
			return '•';
	}
};

export const getNotificationColor = (type: NotificationType): string => {
	switch (type) {
		case 'success':
			return 'text-green-600 bg-green-50 border-green-200';
		case 'error':
			return 'text-red-600 bg-red-50 border-red-200';
		case 'warning':
			return 'text-yellow-600 bg-yellow-50 border-yellow-200';
		case 'info':
			return 'text-blue-600 bg-blue-50 border-blue-200';
		default:
			return 'text-gray-600 bg-gray-50 border-gray-200';
	}
};

// Auto-cleanup old notifications (older than 1 hour)
if (browser) {
	setInterval(() => {
		const oneHourAgo = new Date(Date.now() - 60 * 60 * 1000);
		notifications.update(state => ({
			...state,
			notifications: state.notifications.filter(n => n.createdAt > oneHourAgo)
		}));
	}, 5 * 60 * 1000); // Check every 5 minutes
}

// Export convenience functions
export const { add, success, error, warning, info, dismiss, dismissAll, setPosition, setMaxNotifications, update } = notifications; 