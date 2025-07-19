<script lang="ts">
	import { onMount } from 'svelte';
	import { auth, login, logout } from '$lib/stores/auth.js';
	import { user, currentUser, displayName, userAvatar, userPreferences, updatePreferences } from '$lib/stores/user.js';
	import { theme, actualTheme, isDark, isLight, setTheme } from '$lib/stores/theme.js';
	import { notifications, success, error, warning, info } from '$lib/stores/notifications.js';

	// Local state
	let username = '';
	let password = '';
	let email = '';
	let avatar = '';
	let showProfileForm = false;

	// Theme options
	const themeOptions = [
		{ value: 'light', label: 'Light' },
		{ value: 'dark', label: 'Dark' },
		{ value: 'system', label: 'System' }
	];

	// Language options
	const languageOptions = [
		{ value: 'en', label: 'English' },
		{ value: 'es', label: 'Español' },
		{ value: 'fr', label: 'Français' },
		{ value: 'de', label: 'Deutsch' }
	];

	// Handle login
	async function handleLogin() {
		if (!username || !password) {
			error('Login Failed', 'Please enter both username and password');
			return;
		}

		try {
			await login(username, password);
			success('Login Successful', 'Welcome back!');
			username = '';
			password = '';
		} catch (err) {
			error('Login Failed', err instanceof Error ? err.message : 'Unknown error');
		}
	}

	// Handle logout
	function handleLogout() {
		logout();
		info('Logged Out', 'You have been successfully logged out');
	}

	// Handle profile update
	async function handleProfileUpdate() {
		if (!email) {
			warning('Profile Update', 'Please enter an email address');
			return;
		}

		try {
			await user.updateProfile({ email, avatar });
			success('Profile Updated', 'Your profile has been updated successfully');
			showProfileForm = false;
		} catch (err) {
			error('Profile Update Failed', err instanceof Error ? err.message : 'Unknown error');
		}
	}

	// Handle theme change
	function handleThemeChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		setTheme(target.value as 'light' | 'dark' | 'system');
		info('Theme Changed', `Theme changed to ${target.value}`);
	}

	// Handle language change
	function handleLanguageChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		updatePreferences({ language: target.value });
		info('Language Changed', `Language changed to ${target.value}`);
	}

	// Handle notification test
	function testNotification(type: 'success' | 'error' | 'warning' | 'info') {
		switch (type) {
			case 'success':
				success('Success!', 'This is a success notification');
				break;
			case 'error':
				error('Error!', 'This is an error notification');
				break;
			case 'warning':
				warning('Warning!', 'This is a warning notification');
				break;
			case 'info':
				info('Info!', 'This is an info notification');
				break;
		}
	}

	onMount(() => {
		// Show welcome notification
		info('Welcome!', 'This example demonstrates shared state management across multiple stores');
	});
</script>

<div class="state-management-example">
	<header class="header">
		<h1>State Management Example</h1>
		<p>Demonstrating shared state across multiple stores and components</p>
	</header>

	<!-- Authentication Section -->
	<section class="section">
		<h2>Authentication State</h2>
		<div class="auth-status">
			{#if $auth.isAuthenticated}
				<div class="user-info">
					<img src={$userAvatar} alt="Avatar" class="avatar" />
					<div class="user-details">
						<p><strong>Welcome, {$displayName}!</strong></p>
						<p>User ID: {$currentUser?.id || 'N/A'}</p>
						<p>Email: {$currentUser?.email || 'Not set'}</p>
					</div>
					<button on:click={handleLogout} class="btn btn-danger">
						Logout
					</button>
				</div>
			{:else}
				<div class="login-form">
					<h3>Login</h3>
					<div class="form-group">
						<label for="username">Username:</label>
						<input
							type="text"
							id="username"
							bind:value={username}
							placeholder="Enter username"
							class="input"
						/>
					</div>
					<div class="form-group">
						<label for="password">Password:</label>
						<input
							type="password"
							id="password"
							bind:value={password}
							placeholder="Enter password"
							class="input"
						/>
					</div>
					<button on:click={handleLogin} class="btn btn-primary">
						Login
					</button>
				</div>
			{/if}
		</div>
	</section>

	<!-- Theme Management Section -->
	<section class="section">
		<h2>Theme Management</h2>
		<div class="theme-controls">
			<div class="form-group">
				<label for="theme">Theme:</label>
				<select id="theme" bind:value={$theme.theme} on:change={handleThemeChange} class="select">
					{#each themeOptions as option}
						<option value={option.value}>{option.label}</option>
					{/each}
				</select>
			</div>
			<div class="theme-info">
				<p><strong>Current Theme:</strong> {$actualTheme}</p>
				<p><strong>Is Dark:</strong> {$isDark ? 'Yes' : 'No'}</p>
				<p><strong>Is Light:</strong> {$isLight ? 'Yes' : 'No'}</p>
			</div>
		</div>
	</section>

	<!-- User Preferences Section -->
	{#if $auth.isAuthenticated}
		<section class="section">
			<h2>User Preferences</h2>
			<div class="preferences">
				<div class="form-group">
					<label for="language">Language:</label>
					<select id="language" bind:value={$userPreferences.language} on:change={handleLanguageChange} class="select">
						{#each languageOptions as option}
							<option value={option.value}>{option.label}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>
						<input
							type="checkbox"
							bind:checked={$userPreferences.notifications}
							class="checkbox"
						/>
						Enable Notifications
					</label>
				</div>
			</div>

			<!-- Profile Update Form -->
			<div class="profile-section">
				<button on:click={() => showProfileForm = !showProfileForm} class="btn btn-secondary">
					{showProfileForm ? 'Cancel' : 'Update Profile'}
				</button>
				
				{#if showProfileForm}
					<form on:submit|preventDefault={handleProfileUpdate} class="profile-form">
						<div class="form-group">
							<label for="email">Email:</label>
							<input
								type="email"
								id="email"
								bind:value={email}
								placeholder="Enter email"
								class="input"
							/>
						</div>
						<div class="form-group">
							<label for="avatar">Avatar URL:</label>
							<input
								type="url"
								id="avatar"
								bind:value={avatar}
								placeholder="Enter avatar URL"
								class="input"
							/>
						</div>
						<button type="submit" class="btn btn-primary">Update Profile</button>
					</form>
				{/if}
			</div>
		</section>
	{/if}

	<!-- Notification Testing Section -->
	<section class="section">
		<h2>Notification Testing</h2>
		<div class="notification-buttons">
			<button on:click={() => testNotification('success')} class="btn btn-success">
				Test Success
			</button>
			<button on:click={() => testNotification('error')} class="btn btn-danger">
				Test Error
			</button>
			<button on:click={() => testNotification('warning')} class="btn btn-warning">
				Test Warning
			</button>
			<button on:click={() => testNotification('info')} class="btn btn-info">
				Test Info
			</button>
		</div>
	</section>

	<!-- State Debug Section -->
	<section class="section">
		<h2>State Debug</h2>
		<div class="debug-info">
			<div class="debug-item">
				<h4>Auth State:</h4>
				<pre>{JSON.stringify($auth, null, 2)}</pre>
			</div>
			<div class="debug-item">
				<h4>User State:</h4>
				<pre>{JSON.stringify($user, null, 2)}</pre>
			</div>
			<div class="debug-item">
				<h4>Theme State:</h4>
				<pre>{JSON.stringify($theme, null, 2)}</pre>
			</div>
			<div class="debug-item">
				<h4>User Preferences:</h4>
				<pre>{JSON.stringify($userPreferences, null, 2)}</pre>
			</div>
		</div>
	</section>
</div>

<style>
	.state-management-example {
		max-width: 800px;
		margin: 0 auto;
		padding: 2rem;
		font-family: system-ui, -apple-system, sans-serif;
	}

	.header {
		text-align: center;
		margin-bottom: 3rem;
	}

	.header h1 {
		color: var(--text-color, #1a1a1a);
		margin-bottom: 0.5rem;
	}

	.header p {
		color: var(--text-secondary, #6b7280);
	}

	.section {
		background: var(--surface-color, #f8fafc);
		border: 1px solid var(--border-color, #e5e7eb);
		border-radius: 8px;
		padding: 1.5rem;
		margin-bottom: 2rem;
	}

	.section h2 {
		color: var(--text-color, #1a1a1a);
		margin-bottom: 1rem;
		font-size: 1.25rem;
	}

	.auth-status {
		display: flex;
		justify-content: center;
	}

	.user-info {
		display: flex;
		align-items: center;
		gap: 1rem;
		background: white;
		padding: 1rem;
		border-radius: 8px;
		border: 1px solid var(--border-color, #e5e7eb);
	}

	.avatar {
		width: 48px;
		height: 48px;
		border-radius: 50%;
		object-fit: cover;
	}

	.user-details p {
		margin: 0.25rem 0;
		color: var(--text-color, #1a1a1a);
	}

	.login-form {
		background: white;
		padding: 1.5rem;
		border-radius: 8px;
		border: 1px solid var(--border-color, #e5e7eb);
		max-width: 400px;
	}

	.login-form h3 {
		margin-bottom: 1rem;
		color: var(--text-color, #1a1a1a);
	}

	.form-group {
		margin-bottom: 1rem;
	}

	.form-group label {
		display: block;
		margin-bottom: 0.5rem;
		font-weight: 500;
		color: var(--text-color, #1a1a1a);
	}

	.input, .select {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid var(--border-color, #e5e7eb);
		border-radius: 4px;
		font-size: 1rem;
		background: white;
		color: var(--text-color, #1a1a1a);
	}

	.input:focus, .select:focus {
		outline: none;
		border-color: var(--primary-color, #2563eb);
		box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
	}

	.checkbox {
		margin-right: 0.5rem;
	}

	.btn {
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 4px;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-primary {
		background: var(--primary-color, #2563eb);
		color: white;
	}

	.btn-primary:hover {
		background: var(--primary-hover, #1d4ed8);
	}

	.btn-secondary {
		background: var(--secondary-color, #6b7280);
		color: white;
	}

	.btn-secondary:hover {
		background: var(--secondary-hover, #4b5563);
	}

	.btn-danger {
		background: #dc2626;
		color: white;
	}

	.btn-danger:hover {
		background: #b91c1c;
	}

	.btn-success {
		background: #16a34a;
		color: white;
	}

	.btn-success:hover {
		background: #15803d;
	}

	.btn-warning {
		background: #d97706;
		color: white;
	}

	.btn-warning:hover {
		background: #b45309;
	}

	.btn-info {
		background: #0891b2;
		color: white;
	}

	.btn-info:hover {
		background: #0e7490;
	}

	.theme-controls, .preferences {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 2rem;
		align-items: start;
	}

	.theme-info p {
		margin: 0.5rem 0;
		color: var(--text-color, #1a1a1a);
	}

	.profile-section {
		margin-top: 1.5rem;
		padding-top: 1.5rem;
		border-top: 1px solid var(--border-color, #e5e7eb);
	}

	.profile-form {
		margin-top: 1rem;
		background: white;
		padding: 1rem;
		border-radius: 8px;
		border: 1px solid var(--border-color, #e5e7eb);
	}

	.notification-buttons {
		display: flex;
		gap: 1rem;
		flex-wrap: wrap;
	}

	.debug-info {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}

	.debug-item {
		background: white;
		padding: 1rem;
		border-radius: 4px;
		border: 1px solid var(--border-color, #e5e7eb);
	}

	.debug-item h4 {
		margin-bottom: 0.5rem;
		color: var(--text-color, #1a1a1a);
	}

	.debug-item pre {
		background: #f8f9fa;
		padding: 0.5rem;
		border-radius: 4px;
		font-size: 0.875rem;
		overflow-x: auto;
		margin: 0;
	}

	/* Dark theme adjustments */
	:global([data-theme="dark"]) .state-management-example {
		--text-color: #ffffff;
		--text-secondary: #a1a1aa;
		--surface-color: #2a2a2a;
		--border-color: #3a3a3a;
		--primary-color: #3b82f6;
		--primary-hover: #2563eb;
		--secondary-color: #6b7280;
		--secondary-hover: #4b5563;
	}

	:global([data-theme="dark"]) .input,
	:global([data-theme="dark"]) .select,
	:global([data-theme="dark"]) .user-info,
	:global([data-theme="dark"]) .login-form,
	:global([data-theme="dark"]) .profile-form,
	:global([data-theme="dark"]) .debug-item {
		background: #1a1a1a;
		border-color: #3a3a3a;
		color: #ffffff;
	}

	:global([data-theme="dark"]) .debug-item pre {
		background: #2a2a2a;
		color: #ffffff;
	}

	@media (max-width: 768px) {
		.state-management-example {
			padding: 1rem;
		}

		.theme-controls, .preferences, .debug-info {
			grid-template-columns: 1fr;
		}

		.notification-buttons {
			flex-direction: column;
		}

		.user-info {
			flex-direction: column;
			text-align: center;
		}
	}
</style> 