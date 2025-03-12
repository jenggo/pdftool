<script lang="ts">
    import { parse } from "marked";
    import OcrResultModal from "./OcrResultModal.svelte";

    let file: File | null = null;
    let loading = false;
    let error = "";
    let showModal = false;
    let pages: { html: string; index: number }[] = [];

    function handleFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files.length > 0) {
            file = input.files[0];
        } else {
            file = null;
        }
    }

    async function handleOCR() {
        if (!file) {
            error = "Please select a file";
            return;
        }

        loading = true;
        error = "";
        pages = [];

        try {
            const formData = new FormData();
            formData.append("file", file);

            const response = await fetch("/v1/ocr", {
                method: "POST",
                body: formData,
            });

            const result = await response.json();

            if (!response.ok) {
                throw new Error(result.message || "Failed to process OCR");
            }

            pages = await Promise.all(
                result.data.pages.map(
                    async (page: { markdown: string; index: number }) => ({
                        html: await parse(page.markdown),
                        index: page.index,
                    }),
                ),
            );

            showModal = true;
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

    function closeModal() {
        showModal = false;
    }
</script>

<div>
    <label
        for="ocr-file"
        class="block text-sm font-medium text-gray-700 dark:text-gray-300"
        >PDF File</label
    >
    <div class="mt-1">
        <input
            type="file"
            id="ocr-file"
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
    on:click={handleOCR}
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
        Processing OCR...
    {:else}
        Extract Text
    {/if}
</button>

<OcrResultModal {pages} {showModal} {closeModal} />
