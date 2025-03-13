<script lang="ts">
    import {
        validateFileSize,
        getMaxFileSizeDisplay,
        formatFileSize,
    } from "../utils/fileValidation";

    let file: File | null = null;
    let loading = false;
    let error = "";
    let originalFileSize: number | null = null;
    let optimizedFileSize: number | null = null;

    // For use in the template
    const maxSizeDisplay = getMaxFileSizeDisplay();

    function handleFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files.length > 0) {
            const selectedFile = input.files[0];

            // Check file size
            const sizeError = validateFileSize(selectedFile);
            if (sizeError) {
                error = sizeError;
                file = null;
                return;
            }

            file = selectedFile;
            error = ""; // Clear previous errors
            originalFileSize = selectedFile.size; // Store original file size
            optimizedFileSize = null; // Reset optimized file size
        } else {
            file = null;
            originalFileSize = null;
            optimizedFileSize = null;
        }
    }

    async function handleOptimize() {
        if (!file) {
            error = "Please select a file";
            return;
        }

        loading = true;
        error = "";
        optimizedFileSize = null;

        try {
            const formData = new FormData();
            formData.append("file", file);

            const response = await fetch("/v1/optimize", {
                method: "POST",
                body: formData,
            });

            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.message || "Failed to optimize PDF");
            }

            // Get filename from Content-Disposition header or use a default name
            const contentDisposition = response.headers.get(
                "Content-Disposition",
            );
            const fileName = contentDisposition
                ? contentDisposition.split("filename=")[1].replace(/["']/g, "")
                : "optimized.pdf";

            // Create blob from response and trigger download
            const blob = await response.blob();
            optimizedFileSize = blob.size; // Store optimized file size
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
            originalFileSize = null;
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

<div
    class="max-w-md mx-auto p-6 bg-white dark:bg-gray-800 rounded-lg shadow-md"
>
    <h2 class="text-2xl font-semibold mb-4 text-gray-800 dark:text-white">
        Optimize PDF
    </h2>

    <div>
        <label
            for="optimize-file"
            class="block text-sm font-medium text-gray-700 dark:text-gray-300"
            >PDF File (Max: {maxSizeDisplay})</label
        >
        <div class="mt-1">
            <input
                type="file"
                id="optimize-file"
                accept=".pdf"
                on:change={handleFileChange}
                class="block w-full text-sm text-gray-500 dark:text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 dark:file:bg-blue-900 file:text-blue-700 dark:file:text-blue-300 hover:file:bg-blue-100 dark:hover:file:bg-blue-800"
                required
            />
        </div>
    </div>

    {#if error}
        <div
            class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 mb-4"
            role="alert"
        >
            <div class="flex">
                <div class="ml-3">
                    <h3
                        class="text-sm font-medium text-red-800 dark:text-red-200"
                    >
                        {error}
                    </h3>
                </div>
            </div>
        </div>
    {/if}

    {#if originalFileSize && optimizedFileSize}
        <div class="mb-4">
            <p class="text-sm text-gray-700 dark:text-gray-300">
                Original Size: {formatFileSize(originalFileSize)}
            </p>
            <p class="text-sm text-gray-700 dark:text-gray-300">
                Optimized Size: {formatFileSize(optimizedFileSize)}
            </p>
        </div>
    {/if}

    <button
        type="button"
        disabled={loading || !file}
        on:click={handleOptimize}
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 dark:bg-blue-700 hover:bg-blue-700 dark:hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:focus:ring-offset-gray-900 disabled:opacity-50 disabled:cursor-not-allowed mt-4 cursor-pointer"
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
            Optimizing...
        {:else}
            Optimize PDF
        {/if}
    </button>
</div>
