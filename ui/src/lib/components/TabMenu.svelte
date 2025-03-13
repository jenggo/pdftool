<script lang="ts">
    export let activeTab: "encrypt" | "ocr" | "repair" | "optimize";
    export let setActiveTab: (
        tab: "encrypt" | "ocr" | "repair" | "optimize",
    ) => void;
</script>

<nav
    class="bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-600"
    aria-label="Tabs"
>
    <ul class="flex flex-wrap justify-center sm:justify-start -mb-px">
        <li>
            <button
                class={`inline-block py-4 px-4 text-sm font-medium text-center rounded-t-lg border-b-2 ${
                    activeTab === "encrypt"
                        ? "border-blue-500 text-blue-600 dark:text-blue-400"
                        : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300"
                }`}
                on:click={() => setActiveTab("encrypt")}
            >
                Encrypt/Decrypt
            </button>
        </li>
        <li>
            <button
                class={`inline-block py-4 px-4 text-sm font-medium text-center rounded-t-lg border-b-2 ${
                    activeTab === "repair"
                        ? "border-blue-500 text-blue-600 dark:text-blue-400"
                        : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300"
                }`}
                on:click={() => setActiveTab("repair")}
            >
                Repair
            </button>
        </li>
        <li>
            <button
                class={`inline-block py-4 px-4 text-sm font-medium text-center rounded-t-lg border-b-2 ${
                    activeTab === "ocr"
                        ? "border-blue-500 text-blue-600 dark:text-blue-400"
                        : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300"
                }`}
                on:click={() => setActiveTab("ocr")}
            >
                Extract Text (OCR)
            </button>
        </li>
        <li>
            <button
                class={`inline-block py-4 px-4 text-sm font-medium text-center rounded-t-lg border-b-2 ${
                    activeTab === "optimize"
                        ? "border-blue-500 text-blue-600 dark:text-blue-400"
                        : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300"
                }`}
                on:click={() => setActiveTab("optimize")}
            >
                Optimize
            </button>
        </li>
    </ul>
</nav>
