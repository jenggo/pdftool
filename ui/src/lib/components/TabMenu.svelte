<script lang="ts">
    export let activeTab: "encrypt" | "ocr" | "repair";
    export let setActiveTab: (tab: "encrypt" | "ocr" | "repair") => void;
</script>

<div class="mb-6 border-b border-gray-200">
    <nav class="-mb-px flex" aria-label="Tabs">
        <button
            class={`w-1/2 py-4 px-1 text-center border-b-2 font-medium text-sm cursor-pointer ${
                activeTab === "encrypt"
                    ? "border-blue-500 text-blue-600 dark:text-blue-400"
                    : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300"
            }`}
            on:click={() => setActiveTab("encrypt")}
        >
            Encrypt/Decrypt PDF
        </button>
        <button
            class={`w-1/3 py-4 px-1 text-center border-b-2 font-medium text-sm cursor-pointer ${
                activeTab === "repair"
                    ? "border-blue-500 text-blue-600 dark:text-blue-400"
                    : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300"
            }`}
            on:click={() => setActiveTab("repair")}
        >
            Repair PDF
        </button>
        <button
            class={`w-1/2 py-4 px-1 text-center border-b-2 font-medium text-sm cursor-pointer ${
                activeTab === "ocr"
                    ? "border-blue-500 text-blue-600 dark:text-blue-400"
                    : "border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300"
            }`}
            on:click={() => setActiveTab("ocr")}
        >
            Extract Text (OCR)
        </button>
    </nav>
</div>
