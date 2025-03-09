<script lang="ts">
    import { parse } from "marked";

    // biome-ignore lint: false positive
    let activeTab: "encrypt" | "ocr" = "encrypt";
    let fileInput: HTMLInputElement;
    let file: File | null = null;
    let password = "";
    let loading = false;
    let error = "";
    let ocrResult = "";
    let showModal = false;
    let currentPage = 0;
    let pages: { html: string; index: number }[] = [];

    function handleFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files.length > 0) {
            file = input.files[0];
        } else {
            file = null;
        }
    }

    async function handleEncrypt() {
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

            const response = await fetch("/v1/encrypt", {
                method: "POST",
                body: formData,
            });

            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.message || "Failed to encrypt PDF");
            }

            // Get filename from Content-Disposition header or use a default name
            const contentDisposition = response.headers.get(
                "Content-Disposition",
            );
            const fileName = contentDisposition
                ? contentDisposition.split("filename=")[1].replace(/["']/g, "")
                : "encrypted.pdf";

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

    async function handleOCR() {
        if (!file) {
            error = "Please select a file";
            return;
        }

        loading = true;
        error = "";
        ocrResult = "";
        pages = [];
        currentPage = 0;

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
        currentPage = 0;
    }

    function handleKeydown(event: KeyboardEvent) {
        if (!showModal) return;

        if (event.key === "ArrowLeft" && currentPage > 0) {
            currentPage--;
        } else if (
            event.key === "ArrowRight" &&
            currentPage < pages.length - 1
        ) {
            currentPage++;
        } else if (event.key === "Escape") {
            closeModal();
        }
    }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if showModal}
    <div
        class="fixed inset-0 z-50 overflow-y-auto"
        aria-labelledby="modal-title"
        role="dialog"
        aria-modal="true"
    >
        <div
            class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
        >
            <!-- Background overlay -->
            <button
                aria-label="close"
                class="fixed inset-0 bg-gray-500 dark:bg-gray-900 bg-opacity-75 dark:bg-opacity-75 transition-opacity"
                on:click={closeModal}
            ></button>

            <!-- Modal panel -->
            <div
                class="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-5xl sm:w-full sm:p-6"
            >
                <div class="absolute top-0 right-0 pt-4 pr-4">
                    <button
                        type="button"
                        class="bg-white rounded-md text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                        on:click={closeModal}
                    >
                        <span class="sr-only">Close</span>
                        <svg
                            class="h-6 w-6"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M6 18L18 6M6 6l12 12"
                            />
                        </svg>
                    </button>
                </div>

                <!-- Content -->
                <div class="mt-3 sm:mt-5">
                    <h3
                        class="text-lg leading-6 font-medium text-gray-900 mb-4"
                    >
                        OCR Result - Page {currentPage + 1} of {pages.length}
                    </h3>
                    <div class="mt-2 max-h-[60vh] overflow-y-auto">
                        <div class="prose prose-sm max-w-none">
                            {@html pages[currentPage]?.html || ""}
                        </div>
                    </div>
                </div>

                <!-- Pagination controls -->
                <div class="mt-5 sm:mt-6 flex justify-between items-center">
                    <button
                        type="button"
                        class="inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={currentPage === 0}
                        on:click={() => currentPage--}
                    >
                        Previous
                    </button>

                    <span class="text-sm text-gray-500">
                        Page {currentPage + 1} of {pages.length}
                    </span>

                    <button
                        type="button"
                        class="inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={currentPage === pages.length - 1}
                        on:click={() => currentPage++}
                    >
                        Next
                    </button>
                </div>

                <!-- Download button -->
                <div class="mt-3 flex justify-end">
                    <button
                        type="button"
                        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                        on:click={() => {
                            const allContent = pages
                                .map(
                                    (page, idx) =>
                                        `## Page ${idx + 1}\n\n${page.html}\n\n---\n`,
                                )
                                .join("\n");
                            const blob = new Blob([allContent], {
                                type: "text/plain",
                            });
                            const url = window.URL.createObjectURL(blob);
                            const a = document.createElement("a");
                            a.href = url;
                            a.download = "ocr-result.txt";
                            document.body.appendChild(a);
                            a.click();
                            window.URL.revokeObjectURL(url);
                            document.body.removeChild(a);
                        }}
                    >
                        Download All Pages
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}

<div
    class="min-h-screen bg-gray-100 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8"
>
    <div
        class="max-w-xl mx-auto bg-white dark:bg-gray-800 rounded-lg shadow-lg p-8"
    >
        <div class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                PDF Tools
            </h1>
            <p class="mt-2 text-gray-600 dark:text-gray-300">
                Encrypt or extract text from your PDF files
            </p>
        </div>

        <!-- Tab Menu -->
        <div class="mb-6 border-b border-gray-200">
            <nav class="-mb-px flex" aria-label="Tabs">
                <button
                    class={`w-1/2 py-4 px-1 text-center border-b-2 font-medium text-sm cursor-pointer ${
                        activeTab === "encrypt"
                            ? "border-blue-500 text-blue-600 dark:text-blue-400"
                            : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300"
                    }`}
                    on:click={() => (activeTab = "encrypt")}
                >
                    Encrypt PDF
                </button>
                <button
                    class={`w-1/2 py-4 px-1 text-center border-b-2 font-medium text-sm cursor-pointer ${
                        activeTab === "ocr"
                            ? "border-blue-500 text-blue-600 dark:text-blue-400"
                            : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300"
                    }`}
                    on:click={() => (activeTab = "ocr")}
                >
                    Extract Text (OCR)
                </button>
            </nav>
        </div>

        <form
            on:submit|preventDefault={activeTab === "encrypt"
                ? handleEncrypt
                : handleOCR}
            class="space-y-6"
        >
            <div>
                <label
                    for="file"
                    class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                    >PDF File</label
                >
                <div class="mt-1">
                    <input
                        type="file"
                        id="file"
                        accept=".pdf"
                        bind:this={fileInput}
                        on:change={handleFileChange}
                        class="block w-full text-sm text-gray-500 dark:text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 dark:file:bg-blue-900 file:text-blue-700 dark:file:text-blue-300 hover:file:bg-blue-100 dark:hover:file:bg-blue-800"
                        required
                    />
                </div>
            </div>

            {#if activeTab === "encrypt"}
                <div>
                    <label
                        for="pdf_password"
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                        >Password</label
                    >
                    <div class="mt-1">
                        <input
                            type="password"
                            id="pdf_password"
                            autocomplete="off"
                            bind:value={password}
                            class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                            required
                        />
                    </div>
                </div>
            {/if}

            {#if error}
                <div class="rounded-md bg-red-50 dark:bg-red-900/50 p-4">
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
                type="submit"
                disabled={loading}
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
                    {activeTab === "encrypt"
                        ? "Encrypting..."
                        : "Processing OCR..."}
                {:else}
                    {activeTab === "encrypt" ? "Encrypt PDF" : "Extract Text"}
                {/if}
            </button>
        </form>
    </div>
</div>
