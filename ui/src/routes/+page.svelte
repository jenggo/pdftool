<script lang="ts">
    import TabMenu from "$lib/components/TabMenu.svelte";
    import EncryptDecrypt from "$lib/components/EncryptDecrypt.svelte";
    import RepairPdf from "$lib/components/RepairPdf.svelte";
    import OcrPdf from "$lib/components/OcrPdf.svelte";

    let activeTab: "encrypt" | "ocr" | "repair" = "encrypt";

    function setActiveTab(tab: "encrypt" | "ocr" | "repair") {
        activeTab = tab;
    }
</script>

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
        </div>

        <!-- Tab Menu -->
        <TabMenu {activeTab} {setActiveTab} />

        <div class="space-y-6">
            {#if activeTab === "encrypt"}
                <EncryptDecrypt />
            {:else if activeTab === "repair"}
                <RepairPdf />
            {:else if activeTab === "ocr"}
                <OcrPdf />
            {/if}
        </div>
    </div>
</div>
