import { notifications } from '$lib/stores/notifications';

export interface APIError {
	message: string;
	code?: string;
	details?: Record<string, any>;
	status?: number;
}

export interface ErrorHandlerOptions {
	showNotification?: boolean;
	logToConsole?: boolean;
	retryCount?: number;
	retryDelay?: number;
}

export class ErrorHandler {
	private static instance: ErrorHandler;
	private notificationCallback?: (message: string, type: 'error' | 'warning' | 'info') => void;

	private constructor() {}

	static getInstance(): ErrorHandler {
		if (!ErrorHandler.instance) {
			ErrorHandler.instance = new ErrorHandler();
		}
		return ErrorHandler.instance;
	}

	setNotificationCallback(callback: (message: string, type: 'error' | 'warning' | 'info') => void) {
		this.notificationCallback = callback;
	}

	/**
	 * Parse API error response
	 */
	parseAPIError(response: Response, data?: any): APIError {
		const error: APIError = {
			message: 'An unexpected error occurred',
			status: response.status
		};

		if (data) {
			if (typeof data === 'string') {
				error.message = data;
			} else if (typeof data === 'object') {
				error.message = data.message || data.error || error.message;
				error.code = data.code;
				error.details = data.details;
			}
		}

		// Map common HTTP status codes to user-friendly messages
		switch (response.status) {
			case 400:
				error.message = error.message || 'Invalid request';
				break;
			case 401:
				error.message = error.message || 'Authentication required';
				break;
			case 403:
				error.message = error.message || 'Access denied';
				break;
			case 404:
				error.message = error.message || 'Resource not found';
				break;
			case 409:
				error.message = error.message || 'Conflict occurred';
				break;
			case 422:
				error.message = error.message || 'Validation failed';
				break;
			case 429:
				error.message = error.message || 'Too many requests. Please try again later.';
				break;
			case 500:
				error.message = error.message || 'Internal server error';
				break;
			case 502:
			case 503:
			case 504:
				error.message = error.message || 'Service temporarily unavailable';
				break;
		}

		return error;
	}

	/**
	 * Handle API errors
	 */
	async handleAPIError(response: Response, options: ErrorHandlerOptions = {}): Promise<APIError> {
		const {
			showNotification = true,
			logToConsole = true,
			retryCount = 0,
			retryDelay = 1000
		} = options;

		let data;
		try {
			data = await response.json();
		} catch {
			// If response is not JSON, try to get text
			try {
				data = await response.text();
			} catch {
				data = null;
			}
		}

		const error = this.parseAPIError(response, data);

		if (logToConsole) {
			console.error('API Error:', {
				status: response.status,
				statusText: response.statusText,
				url: response.url,
				error
			});
		}

		if (showNotification && this.notificationCallback) {
			this.notificationCallback(error.message, 'error');
		}

		// Retry logic for certain errors
		if (retryCount > 0 && this.shouldRetry(response.status)) {
			await this.delay(retryDelay);
			return this.handleAPIError(response, { ...options, retryCount: retryCount - 1 });
		}

		return error;
	}

	/**
	 * Handle general errors
	 */
	handleError(error: Error | string, options: ErrorHandlerOptions = {}): void {
		const {
			showNotification = true,
			logToConsole = true
		} = options;

		const message = typeof error === 'string' ? error : error.message;

		if (logToConsole) {
			console.error('Error:', error);
		}

		if (showNotification && this.notificationCallback) {
			this.notificationCallback(message, 'error');
		}
	}

	/**
	 * Handle validation errors
	 */
	handleValidationError(errors: Record<string, string[]>, options: ErrorHandlerOptions = {}): void {
		const {
			showNotification = true,
			logToConsole = true
		} = options;

		const errorMessages = Object.values(errors).flat();
		const message = errorMessages.length > 0 ? errorMessages[0] : 'Validation failed';

		if (logToConsole) {
			console.error('Validation Error:', errors);
		}

		if (showNotification && this.notificationCallback) {
			this.notificationCallback(message, 'error');
		}
	}

	/**
	 * Handle network errors
	 */
	handleNetworkError(error: Error, options: ErrorHandlerOptions = {}): void {
		const {
			showNotification = true,
			logToConsole = true
		} = options;

		const message = 'Network error. Please check your connection and try again.';

		if (logToConsole) {
			console.error('Network Error:', error);
		}

		if (showNotification && this.notificationCallback) {
			this.notificationCallback(message, 'error');
		}
	}

	/**
	 * Retry function with exponential backoff
	 */
	async retry<T>(
		fn: () => Promise<T>,
		maxRetries: number = 3,
		baseDelay: number = 1000
	): Promise<T> {
		let lastError: Error;

		for (let attempt = 0; attempt <= maxRetries; attempt++) {
			try {
				return await fn();
			} catch (error) {
				lastError = error as Error;

				if (attempt === maxRetries) {
					break;
				}

				const delay = baseDelay * Math.pow(2, attempt);
				await this.delay(delay);
			}
		}

		throw lastError!;
	}

	/**
	 * Check if an error should be retried
	 */
	private shouldRetry(status: number): boolean {
		// Retry on server errors and rate limiting
		return status >= 500 || status === 429;
	}

	/**
	 * Delay utility
	 */
	private delay(ms: number): Promise<void> {
		return new Promise(resolve => setTimeout(resolve, ms));
	}

	/**
	 * Create a user-friendly error message
	 */
	createUserFriendlyMessage(error: APIError | Error | string): string {
		if (typeof error === 'string') {
			return error;
		}

		if ('status' in error) {
			// API Error
			return error.message;
		}

		// General Error
		const message = error.message;
		
		// Map common error messages to user-friendly versions
		const friendlyMessages: Record<string, string> = {
			'NetworkError when attempting to fetch resource.': 'Network error. Please check your connection.',
			'Failed to fetch': 'Unable to connect to the server. Please try again.',
			'Request timeout': 'Request timed out. Please try again.',
			'User denied the request': 'Permission denied. Please check your settings.'
		};

		return friendlyMessages[message] || message;
	}
}

// Export singleton instance
export const errorHandler = ErrorHandler.getInstance(); 