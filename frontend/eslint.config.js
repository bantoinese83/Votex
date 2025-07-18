import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import svelte from 'eslint-plugin-svelte';
import prettier from 'eslint-config-prettier';
import storybook from 'eslint-plugin-storybook';

export default tseslint.config(
	eslint.configs.recommended,
	...tseslint.configs.recommended,
	...svelte.configs['flat/recommended'],
	prettier,
	...svelte.configs['flat/prettier'],
	...storybook.configs['flat/recommended'],
	{
		files: ['*.svelte'],
		parser: 'svelte-eslint-parser',
		parserOptions: {
			parser: tseslint.parser
		}
	},
	{
		ignores: ['.svelte-kit/', 'src/lib/paraglide/']
	}
);
