<script lang="ts">
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import { errorHandler } from '$lib/utils/error-handler';
	import { healthApi } from '$lib/api/client';

	let showModal = false;
	let username = '';
	let email = '';
	let password = '';
	let loading = false;

	function handleSubmit() {
		loading = true;
		
		// Simulate API call
		setTimeout(() => {
			loading = false;
			showModal = true;
		}, 1000);
	}

	function handleError() {
		errorHandler.handleError('This is a test error message');
	}

	function handleAPIError() {
		// Simulate API error
		const mockResponse = new Response('{"error": "Test API error"}', {
			status: 400,
			statusText: 'Bad Request'
		});
		
		errorHandler.handleAPIError(mockResponse);
	}
</script>

<svelte:head>
	<title>Vortex Template - Production-Ready Full-Stack App</title>
	<meta name="description" content="A comprehensive, production-ready full-stack application featuring Go backend with SvelteKit frontend" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<!-- Header -->
	<header class="bg-white shadow-sm">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between items-center py-6">
				<div class="flex items-center">
					<h1 class="text-2xl font-bold text-gray-900">Vortex Template</h1>
					<span class="ml-2 px-2 py-1 text-xs font-medium bg-green-100 text-green-800 rounded-full">
						Production Ready
					</span>
				</div>
				<nav class="flex space-x-8">
					<a href="#features" class="text-gray-500 hover:text-gray-900">Features</a>
					<a href="#demo" class="text-gray-500 hover:text-gray-900">Demo</a>
					<a href="http://localhost:8080/api/docs" target="_blank" class="text-gray-500 hover:text-gray-900">API Docs</a>
				</nav>
			</div>
		</div>
	</header>

	<!-- Hero Section -->
	<section class="py-20 bg-gradient-to-r from-blue-600 to-purple-600">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
			<h2 class="text-4xl font-bold text-white mb-6">
				Production-Ready Full-Stack Template
			</h2>
			<p class="text-xl text-blue-100 mb-8 max-w-3xl mx-auto">
				Complete with authentication, API documentation, CI/CD pipeline, security features, 
				and reusable UI components. Built with Go backend and SvelteKit frontend.
			</p>
			<div class="flex justify-center space-x-4">
				<Button variant="primary" size="lg" on:click={() => document.getElementById('demo')?.scrollIntoView({ behavior: 'smooth' })}>
					Try Demo
				</Button>
				<Button variant="outline" size="lg" on:click={() => window.open('http://localhost:8080/api/docs', '_blank')}>
					View API Docs
				</Button>
			</div>
		</div>
	</section>

	<!-- Features Section -->
	<section id="features" class="py-20">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<h3 class="text-3xl font-bold text-center text-gray-900 mb-12">Key Features</h3>
			
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
				<Card>
					<div class="text-center">
						<div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mx-auto mb-4">
							<svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
							</svg>
						</div>
						<h4 class="text-lg font-semibold text-gray-900 mb-2">Authentication</h4>
						<p class="text-gray-600">JWT-based auth with password reset, email verification, and secure token handling.</p>
					</div>
				</Card>

				<Card>
					<div class="text-center">
						<div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mx-auto mb-4">
							<svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
							</svg>
						</div>
						<h4 class="text-lg font-semibold text-gray-900 mb-2">API Documentation</h4>
						<p class="text-gray-600">Complete OpenAPI 3.1 specification with interactive Swagger UI.</p>
					</div>
				</Card>

				<Card>
					<div class="text-center">
						<div class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mx-auto mb-4">
							<svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
							</svg>
						</div>
						<h4 class="text-lg font-semibold text-gray-900 mb-2">CI/CD Pipeline</h4>
						<p class="text-gray-600">Automated testing, security scanning, Docker builds, and deployment.</p>
					</div>
				</Card>

				<Card>
					<div class="text-center">
						<div class="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center mx-auto mb-4">
							<svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
							</svg>
						</div>
						<h4 class="text-lg font-semibold text-gray-900 mb-2">Security</h4>
						<p class="text-gray-600">Rate limiting, security headers, input validation, and XSS protection.</p>
					</div>
				</Card>

				<Card>
					<div class="text-center">
						<div class="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center mx-auto mb-4">
							<svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"></path>
							</svg>
						</div>
						<h4 class="text-lg font-semibold text-gray-900 mb-2">UI Components</h4>
						<p class="text-gray-600">Reusable components with TypeScript, Tailwind CSS, and accessibility.</p>
					</div>
				</Card>

				<Card>
					<div class="text-center">
						<div class="w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center mx-auto mb-4">
							<svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
							</svg>
						</div>
						<h4 class="text-lg font-semibold text-gray-900 mb-2">Monitoring</h4>
						<p class="text-gray-600">Health checks, structured logging, error tracking, and performance monitoring.</p>
					</div>
				</Card>
			</div>
		</div>
	</section>

	<!-- Demo Section -->
	<section id="demo" class="py-20 bg-gray-100">
		<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
			<h3 class="text-3xl font-bold text-center text-gray-900 mb-12">Component Demo</h3>
			
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
				<!-- Form Demo -->
				<Card>
					<h4 class="text-xl font-semibold text-gray-900 mb-6">Form Components</h4>
					<form on:submit|preventDefault={handleSubmit} class="space-y-4">
						<Input
							label="Username"
							bind:value={username}
							placeholder="Enter username"
							required
						/>
						<Input
							type="email"
							label="Email"
							bind:value={email}
							placeholder="Enter email"
							required
						/>
						<Input
							type="password"
							label="Password"
							bind:value={password}
							placeholder="Enter password"
							required
						/>
						<div class="flex space-x-3">
							<Button type="submit" variant="primary" {loading} fullWidth>
								{loading ? 'Submitting...' : 'Submit'}
							</Button>
							<Button type="button" variant="outline" on:click={handleError}>
								Test Error
							</Button>
						</div>
					</form>
				</Card>

				<!-- Button Demo -->
				<Card>
					<h4 class="text-xl font-semibold text-gray-900 mb-6">Button Variants</h4>
					<div class="space-y-4">
						<div class="flex flex-wrap gap-3">
							<Button variant="primary">Primary</Button>
							<Button variant="secondary">Secondary</Button>
							<Button variant="outline">Outline</Button>
							<Button variant="ghost">Ghost</Button>
							<Button variant="danger">Danger</Button>
						</div>
						<div class="flex flex-wrap gap-3">
							<Button size="sm">Small</Button>
							<Button size="md">Medium</Button>
							<Button size="lg">Large</Button>
						</div>
						<div class="flex flex-wrap gap-3">
							<Button disabled>Disabled</Button>
							<Button loading>Loading</Button>
						</div>
					</div>
				</Card>
			</div>

			<!-- Error Handling Demo -->
			<Card className="mt-8">
				<h4 class="text-xl font-semibold text-gray-900 mb-6">Error Handling Demo</h4>
				<div class="flex space-x-3">
					<Button variant="outline" on:click={handleAPIError}>
						Simulate API Error
					</Button>
					<Button variant="outline" on:click={() => errorHandler.handleNetworkError(new Error('Network error'))}>
						Simulate Network Error
					</Button>
					<Button variant="outline" on:click={() => errorHandler.handleValidationError({ username: ['Username is required'] })}>
						Simulate Validation Error
					</Button>
				</div>
			</Card>
		</div>
	</section>

	<!-- Footer -->
	<footer class="bg-gray-900 text-white py-12">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="text-center">
				<h3 class="text-2xl font-bold mb-4">Vortex Template</h3>
				<p class="text-gray-400 mb-6">
					A production-ready foundation for your next full-stack application
				</p>
				<div class="flex justify-center space-x-6">
					<a href="/api/docs" class="text-gray-400 hover:text-white">API Documentation</a>
					<a href="https://github.com" class="text-gray-400 hover:text-white">GitHub</a>
					<a href="/health" class="text-gray-400 hover:text-white">Health Check</a>
				</div>
			</div>
		</div>
	</footer>
</div>

<!-- Modal -->
<Modal bind:open={showModal} title="Success!" size="sm">
	<div class="text-center">
		<div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
			<svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
			</svg>
		</div>
		<h3 class="text-lg font-semibold text-gray-900 mb-2">Form Submitted!</h3>
		<p class="text-gray-600 mb-6">
			This demonstrates the modal component and form handling capabilities.
		</p>
	</div>
	
	<svelte:fragment slot="footer">
		<Button variant="primary" on:click={() => showModal = false}>
			Close
		</Button>
	</svelte:fragment>
</Modal>
