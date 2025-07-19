<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let type: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' = 'text';
	export let placeholder = '';
	export let value = '';
	export let disabled = false;
	export let required = false;
	export let name = '';
	export let id = '';
	export let label = '';
	export let error = '';
	export let fullWidth = false;
	export let autocomplete: string | undefined = undefined;

	const dispatch = createEventDispatcher<{
		input: string;
		change: string;
		focus: FocusEvent;
		blur: FocusEvent;
	}>();

	// Generate unique ID if not provided
	$: uniqueId = id || `${name || type}-${Math.random().toString(36).substr(2, 9)}`;

	// Set appropriate autocomplete based on type and name
	$: autoCompleteValue = autocomplete || (() => {
		switch (type) {
			case 'email':
				return 'email';
			case 'password':
				return name?.includes('password') ? 'current-password' : 'new-password';
			case 'text':
				if (name?.includes('username') || name?.includes('user')) return 'username';
				if (name?.includes('name')) return 'name';
				return 'off';
			default:
				return 'off';
		}
	})();

	function handleInput(event: Event) {
		const target = event.target as HTMLInputElement;
		value = target.value;
		dispatch('input', value);
	}

	function handleChange(event: Event) {
		const target = event.target as HTMLInputElement;
		value = target.value;
		dispatch('change', value);
	}

	function handleFocus(event: FocusEvent) {
		dispatch('focus', event);
	}

	function handleBlur(event: FocusEvent) {
		dispatch('blur', event);
	}

	$: inputClasses = [
		'block w-full rounded-md border-gray-300 shadow-sm transition-colors duration-200 focus:border-blue-500 focus:ring-blue-500 sm:text-sm',
		error ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : '',
		disabled ? 'bg-gray-50 cursor-not-allowed' : 'bg-white',
		fullWidth ? 'w-full' : ''
	].filter(Boolean).join(' ');
</script>

<div class="space-y-1">
	{#if label}
		<label for={uniqueId} class="block text-sm font-medium text-gray-700">
			{label}
			{#if required}
				<span class="text-red-500 ml-1">*</span>
			{/if}
		</label>
	{/if}
	
	<input
		{type}
		{placeholder}
		bind:value
		{disabled}
		{required}
		{name}
		id={uniqueId}
		autocomplete={autoCompleteValue as any}
		class={inputClasses}
		on:input={handleInput}
		on:change={handleChange}
		on:focus={handleFocus}
		on:blur={handleBlur}
		aria-invalid={error ? 'true' : 'false'}
		aria-describedby={error ? `${uniqueId}-error` : undefined}
	/>
	
	{#if error}
		<p id="{uniqueId}-error" class="text-sm text-red-600">
			{error}
		</p>
	{/if}
</div> 