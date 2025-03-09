<script lang="ts">
    import "../app.css";
    import AuthCheck from "$lib/AuthCheck.svelte";
    // biome-ignore lint: false positive
    let { children } = $props();
</script>

<AuthCheck />
{@render children()}
