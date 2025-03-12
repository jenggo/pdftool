<script lang="ts">
    let file: File | null = null;
    let loading = false;
    let error = "";

    function handleFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files.length > 0) {
            file = input.files[0];
        } else {
            file = null;
        }
    }

    async function handleRepair() {
        if (!file) {
            error = "Please select a file";
            return;
        }

        loading = true;
        error = "";

        try {
            const formData = new FormData();
            formData.append("file", file);

            const response = await fetch("/v1/repair", {
                method: "POST",
                body: formData,
            });

            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.message || "Failed to repair PDF");
            }

            // Get filename from Content-Disposition header or use a default name
            const contentDisposition = response.headers.get(
                "Content-Disposition",
            );
            const fileName = contentDisposition
                ? contentDisposition.split("filename=")[1].replace(/["']/g, "")
                : "repaired.pdf";

            // Create blob from response and trigger download
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = url;
            a.download = fileName;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);

            // Reset form
            file = null;
        } catch (err) {
            error =
                err instanceof Error
                    ? err.message
                    : "An unexpected error occurred";
        } finally {
            loading = false;
        }
    }
</script>

<div class="rounded-md bg-yellow-50 dark:bg-yellow-900/30 p-4 mb-6">
    <div class="flex">
        <div class="flex-shrink-0">
            <svg
                class="h-5 w-5 text-yellow-400"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
            >
                <path
                    fill-rule="evenodd"
                    d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                    clip-rule="evenodd"
                />
            </svg>
        </div>
        <div class="ml-3">
            <h3
                class="text-sm font-medium text-yellow-800 dark:text-yellow-200"
            >
                Information
            </h3>
            <div class="mt-2 text-sm text-yellow-700 dark:text-yellow-300">
                <p>
                    This tool attempts to repair corrupted PDF files. Results
                    may vary depending on the level of corruption.
                </p>
            </div>
        </div>
    </div>
</div>

<div>
    <label
        for="repair-file"
        class="block text-sm font-medium text-gray-700 dark:text-gray-300"
        >PDF File</label
    >
    <div class="mt-1">
        <input
            type="file"
            id="repair-file"
            accept=".pdf"
            on:change={handleFileChange}
            class="block w-full text-sm text-gray-500 dark:text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 dark:file:bg-blue-900 file:text-blue-700 dark:file:text-blue-300 hover:file:bg-blue-100 dark:hover:file:bg-blue-800"
            required
        />
    </div>
</div>

{#if error}
    <div class="rounded-md bg-red-50 dark:bg-red-900/50 p-4">
        <div class="flex">
            <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800 dark:text-red-200">
                    {error}
                </h3>
            </div>
        </div>
    </div>
{/if}

<button
    type="button"
    disabled={loading}
    on:click={handleRepair}
    class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 dark:bg-blue-700 hover:bg-blue-700 dark:hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:focus:ring-offset-gray-900 disabled:opacity-50 disabled:cursor-not-allowed"
>
    {#if loading}
        <svg
            class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
        >
            <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
            ></circle>
            <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
        </svg>
        Repairing...
    {:else}
        Repair PDF
    {/if}
</button>
