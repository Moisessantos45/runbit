<template>
    <section class="flex h-full flex-1 flex-col overflow-hidden rounded-md border border-[#44475A] bg-[#21222C] relative">
        <div v-if="isLoading"
            class="absolute inset-0 z-10 flex items-center justify-center bg-[#21222C]/80 backdrop-blur-sm">
            <div class="flex items-center gap-3 text-[#F8F8F2] text-sm">
                <span class="spinner" aria-hidden="true"></span>
                Loading preview...
            </div>
        </div>
        <div class="flex-1 overflow-auto px-4 py-3">
            <div class="space-y-2 text-sm font-mono leading-6">
                <div v-for="(line, index) in currentTabLineItems" :key="index"
                    class="whitespace-pre-wrap wrap-break-word rounded-sm px-2" :class="line.event === 'stderr' || line.event === 'error'
                        ? 'bg-red-500/10 text-red-300'
                        : 'text-[#F8F8F2]'
                        ">
                    {{ line.result }}
                </div>
            </div>
        </div>
    </section>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, computed } from "vue";
import { storeToRefs } from "pinia";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import usePlaygroundStore from "../store/Playground";
import { getNumber, getString } from "@/helpers/formatters";

const playgroundStore = usePlaygroundStore();
const { stdoutLines, activeTab, loadingByTab } = storeToRefs(playgroundStore);
const { setTabLoading } = playgroundStore;

type RunnerPayload = {
    tab?: unknown;
    result?: unknown;
};

const pushLine = (event: string, data: RunnerPayload) => {
    const text = getString(data?.result);
    const tab = getNumber(data?.tab);

    setTabLoading(tab, false);

    if (!text.trim()) return;

    const lines = text.split("\n").filter((line) => line.trim().length > 0);

    for (const line of lines) {
        stdoutLines.value.push({
            event,
            result: line,
            tab,
        });
    }
};

const currentTabLineItems = computed(() => {
    return stdoutLines.value.filter((line) => line.tab === activeTab.value);
});

const isLoading = computed(() => {
    return Boolean(loadingByTab.value[activeTab.value]);
});

const offStart = EventsOn("runner:start", (payload) => {
    const tab = getNumber(payload?.tab);
    stdoutLines.value = stdoutLines.value.filter((line) => line.tab !== tab);
    setTabLoading(tab, true);
});

const offStdout = EventsOn("runner:stdout", (payload) => {
    pushLine("stdout", payload);
});

const offStderr = EventsOn("runner:stderr", (payload) => {
    pushLine("stderr", payload);
});

const offError = EventsOn("runner:error", (payload) => {
    pushLine("error", payload);
});

onBeforeUnmount(() => {
    offStart();
    offStdout();
    offStderr();
    offError();
});
</script>

<style scoped>
.spinner {
    width: 16px;
    height: 16px;
    border-radius: 9999px;
    border: 2px solid rgba(248, 248, 242, 0.3);
    border-top-color: #F8F8F2;
    animation: spin 0.8s linear infinite;
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}
</style>