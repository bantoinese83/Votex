<script lang="ts">
	import type { PageData } from './$types';

	export let data: PageData;
</script>

<h1>Hi, {data.user.username}!</h1>
<p>Your user ID is {data.user.id}.</p>
<form method="post" action="?/logout">
	<button>Sign out</button>
</form>
