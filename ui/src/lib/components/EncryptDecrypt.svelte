<script lang="ts">
    import {
        validateFileSize,
        getMaxFileSizeDisplay,
    } from "../utils/fileValidation";

    // biome-ignore lint: false positive
    let isEncrypt = true;
    let file: File | null = null;
    let password = "";
    let loading = false;
    let error = "";
    // biome-ignore lint: false positive
    let showPassword = false;

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
        } else {
            file = null;
        }
    }

    async function handleEncryptDecrypt() {
        if (!file || !password) {
            error = "Please select a file and enter a password";
            return;
        }

        loading = true;
        error = "";

        try {
            const formData = new FormData();
            formData.append("file", file);
            formData.append("pdf_password", password);

            const endpoint = isEncrypt ? "/v1/encrypt" : "/v1/decrypt";
            const response = await fetch(endpoint, {
                method: "POST",
                body: formData,
            });

            if (!response.ok) {
                const data = await response.json();
                throw new Error(
                    data.message ||
                        `Failed to ${isEncrypt ? "encrypt" : "decrypt"} PDF`,
                );
            }

            // Get filename from Content-Disposition header or use a default name
            const contentDisposition = response.headers.get(
                "Content-Disposition",
            );
            const fileName = contentDisposition
                ? contentDisposition.split("filename=")[1].replace(/["']/g, "")
                : isEncrypt
                  ? "encrypted.pdf"
                  : "decrypted.pdf";

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
            password = "";
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
        {isEncrypt ? "Encrypt PDF" : "Decrypt PDF"}
    </h2>

    <!-- Encrypt/Decrypt Switch -->
    <div class="flex items-center justify-center mb-6">
        <label class="inline-flex relative items-center cursor-pointer">
            <input
                type="checkbox"
                class="sr-only peer"
                bind:checked={isEncrypt}
                id="encrypt-decrypt-switch"
            />
            <div
                class="w-14 h-7 bg-gray-200 rounded-full peer dark:bg-gray-700
                       peer-checked:after:translate-x-full peer-checked:after:border-white
                       after:content-[''] after:absolute after:top-0.5 after:left-[2px]
                       after:bg-white after:border-gray-300 after:border after:rounded-full
                       after:h-6 after:w-6 after:transition-all dark:border-gray-600
                       peer-checked:bg-blue-600"
            ></div>
            <span
                class="ml-3 text-sm font-medium text-gray-900 dark:text-gray-300"
            >
                {isEncrypt ? "Encrypt" : "Decrypt"}
            </span>
        </label>
    </div>

    <div>
        <label
            for="file"
            class="block text-sm font-medium text-gray-700 dark:text-gray-300"
            >PDF File (Max: {maxSizeDisplay})</label
        >
        <div class="mt-1">
            <input
                type="file"
                id="file"
                accept=".pdf"
                on:change={handleFileChange}
                class="block w-full text-sm text-gray-500 dark:text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 dark:file:bg-blue-900 file:text-blue-700 dark:file:text-blue-300 hover:file:bg-blue-100 dark:hover:file:bg-blue-800"
                required
            />
        </div>
    </div>

    <div class="mb-4">
        <label
            for="pdf_password"
            class="block text-sm font-medium text-gray-700 dark:text-gray-300"
            >{isEncrypt ? "Set Password" : "Enter Password"}</label
        >
        <div class="mt-1 relative">
            <input
                type={showPassword ? "text" : "password"}
                id="pdf_password"
                autocomplete="off"
                bind:value={password}
                class="shadow-sm focus:ring-blue-500 focus:border-blue-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white pr-10"
                required
            />
            <div
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-sm leading-5"
            >
                <button
                    type="button"
                    class="text-gray-500 hover:text-gray-700 focus:outline-none"
                    on:click={() => (showPassword = !showPassword)}
                >
                    {#if showPassword}
                        <svg
                            class="h-5 w-5"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                            ></path>
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7 1.274 4.057-1.512 9-5.458 9-3.947 0-6.733-4.943-8.007-9z"
                            ></path>
                        </svg>
                    {:else}
                        <svg
                            class="h-5 w-5"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M13.875 18.825A10.05 10.05 0 0012 19.5c-.67.005-1.335-.11-2-.315m-1.375-1.575A10.07 10.07 0 0012 13.5c.67-.005 1.335.11 2 .315m4.125 1.575A10.07 10.07 0 0012 16.5c-.67.005-1.335-.11-2-.315M21 12c0 4.556-4.03 8.25-9 8.25a9.764 9.764 0 01-2.555-.337m-2.415.337c-4.556 0-8.25-4.03-8.25-9 0-4.556 4.03-8.25 9-8.25a9.764 9.764 0 012.555.337m2.415-.337C16.97 3 21 6.694 21 12z"
                            ></path>
                        </svg>
                    {/if}
                </button>
            </div>
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

    <button
        type="button"
        disabled={loading || !file}
        on:click={handleEncryptDecrypt}
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
            {isEncrypt ? "Encrypting..." : "Decrypting..."}
        {:else}
            {isEncrypt ? "Encrypt PDF" : "Decrypt PDF"}
        {/if}
    </button>
</div>
