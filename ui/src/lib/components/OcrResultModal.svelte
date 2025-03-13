<script lang="ts">
    import { onMount } from "svelte";

    // biome-ignore lint: false positive
    export let pages: { html: string; index: number }[] = [];
    // biome-ignore lint: false positive
    export let showModal: boolean = false;
    export let closeModal: () => void;

    let currentPage = 0;

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

    onMount(() => {
        window.addEventListener("keydown", handleKeydown);
        return () => {
            window.removeEventListener("keydown", handleKeydown);
        };
    });

    function downloadAllPages() {
        const allContent = pages
            .map((page, idx) => `## Page ${idx + 1}\n\n${page.html}\n\n---\n`)
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
    }

    function copyToClipboard() {
        const text = pages[currentPage]?.html || "";
        navigator.clipboard.writeText(text);
    }
</script>

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

                <!-- Download and Copy buttons -->
                <div class="mt-3 flex justify-between">
                    <button
                        type="button"
                        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                        on:click={downloadAllPages}
                    >
                        Download All Pages
                    </button>
                    <button
                        type="button"
                        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                        on:click={copyToClipboard}
                    >
                        Copy to Clipboard
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}
