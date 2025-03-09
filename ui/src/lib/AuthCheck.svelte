<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";

    onMount(async () => {
        if ($page.url.pathname === "/login") {
            return;
        }

        try {
            const response = await fetch("/check-auth");
            if (!response.ok) {
                goto("/login");
            }
        } catch (error) {
            goto("/login");
        }
    });
</script>
