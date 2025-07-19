<script lang="ts">
	export let padding: 'none' | 'sm' | 'md' | 'lg' = 'md';
	export let shadow: 'none' | 'sm' | 'md' | 'lg' = 'md';
	export let hover = false;
	export let fullWidth = false;
	export let className = '';

	$: paddingClasses = {
		none: '',
		sm: 'p-3',
		md: 'p-6',
		lg: 'p-8'
	};

	$: shadowClasses = {
		none: '',
		sm: 'shadow-sm',
		md: 'shadow-md',
		lg: 'shadow-lg'
	};

	$: cardClasses = [
		'bg-white rounded-lg border border-gray-200',
		paddingClasses[padding],
		shadowClasses[shadow],
		hover ? 'transition-shadow duration-200 hover:shadow-lg' : '',
		fullWidth ? 'w-full' : '',
		className
	].filter(Boolean).join(' ');
</script>

<div class={cardClasses}>
	<slot />
</div> 