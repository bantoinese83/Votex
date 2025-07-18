<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let variant: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger' = 'primary';
	export let size: 'sm' | 'md' | 'lg' = 'md';
	export let disabled = false;
	export let loading = false;
	export let type: 'button' | 'submit' | 'reset' = 'button';
	export let fullWidth = false;

	const dispatch = createEventDispatcher<{
		click: MouseEvent;
	}>();

	function handleClick(event: MouseEvent) {
		if (!disabled && !loading) {
			dispatch('click', event);
		}
	}

	$: variantClasses = {
		primary: 'bg-blue-600 hover:bg-blue-700 text-white shadow-sm',
		secondary: 'bg-gray-600 hover:bg-gray-700 text-white shadow-sm',
		outline: 'border border-gray-300 hover:bg-gray-50 text-gray-700',
		ghost: 'hover:bg-gray-100 text-gray-700',
		danger: 'bg-red-600 hover:bg-red-700 text-white shadow-sm'
	};

	$: sizeClasses = {
		sm: 'px-3 py-1.5 text-sm',
		md: 'px-4 py-2 text-sm',
		lg: 'px-6 py-3 text-base'
	};

	$: buttonClasses = [
		'inline-flex items-center justify-center font-medium rounded-md transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2',
		variantClasses[variant],
		sizeClasses[size],
		fullWidth ? 'w-full' : '',
		disabled || loading ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'
	].filter(Boolean).join(' ');
</script>

<button
	{type}
	class={buttonClasses}
	{disabled}
	on:click={handleClick}
	on:keydown={(e) => e.key === 'Enter' && handleClick(e as any)}
>
	{#if loading}
		<svg class="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
			<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
			<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
		</svg>
	{/if}
	<slot />
</button> 