<script lang="ts">
    import TabMenu from "$lib/components/TabMenu.svelte";
    import EncryptDecrypt from "$lib/components/EncryptDecrypt.svelte";
    import RepairPdf from "$lib/components/RepairPdf.svelte";
    import OcrPdf from "$lib/components/OcrPdf.svelte";
    import OptimizePdf from "$lib/components/OptimizePdf.svelte";

    let activeTab: "encrypt" | "ocr" | "repair" | "optimize" = "encrypt";

    function setActiveTab(tab: "encrypt" | "ocr" | "repair" | "optimize") {
        activeTab = tab;
    }
</script>

<main class="min-h-screen bg-gray-100 dark:bg-gray-900 py-6 sm:py-12">
    <div class="container mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                PDF Tools
            </h1>
        </div>

        <div
            class="bg-white dark:bg-gray-800 rounded-lg shadow-xl overflow-hidden max-w-sm sm:max-w-lg md:max-w-xl lg:max-w-2xl xl:max-w-3xl mx-auto"
        >
            <TabMenu {activeTab} {setActiveTab} />

            <div class="p-4 sm:p-6">
                {#if activeTab === "encrypt"}
                    <EncryptDecrypt />
                {:else if activeTab === "repair"}
                    <RepairPdf />
                {:else if activeTab === "ocr"}
                    <OcrPdf />
                {:else if activeTab === "optimize"}
                    <OptimizePdf />
                {/if}
            </div>
        </div>
    </div>
</main>
