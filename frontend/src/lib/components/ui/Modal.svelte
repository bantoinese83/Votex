<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	export let open = false;
	export let title = '';
	export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
	export let closeOnOverlayClick = true;
	export let closeOnEscape = true;

	const dispatch = createEventDispatcher<{
		close: void;
		open: void;
	}>();

	let modalElement: HTMLDivElement;

	function handleOverlayClick(event: MouseEvent) {
		if (closeOnOverlayClick && event.target === event.currentTarget) {
			close();
		}
	}

	function handleOverlayKeydown(event: KeyboardEvent) {
		if (closeOnOverlayClick && event.key === 'Enter') {
			close();
		}
	}

	function handleKeydown(event: KeyboardEvent) {
		if (closeOnEscape && event.key === 'Escape') {
			close();
		}
	}

	function close() {
		open = false;
		dispatch('close');
	}

	function openModal() {
		open = true;
		dispatch('open');
	}

	$: sizeClasses = {
		sm: 'max-w-md',
		md: 'max-w-lg',
		lg: 'max-w-2xl',
		xl: 'max-w-4xl'
	};

	onMount(() => {
		if (closeOnEscape) {
			document.addEventListener('keydown', handleKeydown);
		}

		return () => {
			if (closeOnEscape) {
				document.removeEventListener('keydown', handleKeydown);
			}
		};
	});

	// Focus management
	$: if (open && modalElement) {
		// Focus the modal when it opens
		setTimeout(() => {
			modalElement?.focus();
		}, 100);
	}
</script>

{#if open}
	<div
		bind:this={modalElement}
		class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 transition-opacity duration-200"
		on:click={handleOverlayClick}
		on:keydown={handleOverlayKeydown}
		role="dialog"
		aria-modal="true"
		aria-labelledby="modal-title"
		tabindex="-1"
	>
		<div class="relative w-full {sizeClasses[size]} mx-4">
			<div class="bg-white rounded-lg shadow-xl transform transition-all duration-200">
				<!-- Header -->
				{#if title}
					<div class="flex items-center justify-between p-6 border-b border-gray-200">
						<h2 id="modal-title" class="text-lg font-semibold text-gray-900">
							{title}
						</h2>
						<button
							type="button"
							class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
							on:click={close}
							aria-label="Close modal"
						>
							<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
							</svg>
						</button>
					</div>
				{/if}

				<!-- Content -->
				<div class="p-6">
					<slot />
				</div>

				<!-- Footer -->
				{#if $$slots.footer}
					<div class="flex items-center justify-end gap-3 p-6 border-t border-gray-200">
						<slot name="footer" />
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if} 