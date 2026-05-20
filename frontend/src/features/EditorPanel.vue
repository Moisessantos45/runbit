<template>
    <section
        class="panel-left flex h-full flex-1 flex-col overflow-hidden rounded-md border border-[#44475A] bg-[#21222C]"
        :style="{ '--left-width': props.leftWidth + 'px' }">
        <div class="flex-1 overflow-auto py-2">
            <div class="h-full" :style="{ height: editorHeight ? editorHeight + 'px' : '100%' }">
                <vue-monaco-editor v-model:value="currentSnippet" language="typescript" theme="vs-dark"
                    :options="editorOptions" @beforeMount="handleBeforeMount" @mount="handleEditorMount" class="h-full" />
            </div>
        </div>
    </section>
</template>

<script lang="ts" setup>
import { ref, onBeforeUnmount, onMounted, computed, watch, nextTick } from 'vue'
import { storeToRefs } from "pinia";
import { EventsEmit } from "../../wailsjs/runtime/runtime";
import usePlaygroundStore from "../store/Playground";

const props = defineProps<{
    leftWidth: number
}>()

const playgroundStore = usePlaygroundStore();
const { codeSnippets, activeTab } = storeToRefs(playgroundStore);

const editorHeight = ref(0)
let debounceTimer: ReturnType<typeof setTimeout> | null = null
let isProgrammatic = false;

watch(activeTab, async (newTab) => {
    isProgrammatic = true;
    await nextTick();
    setTimeout(() => {
        isProgrammatic = false;
    }, 50);

    const hasOutput = playgroundStore.stdoutLines.some(line => line.tab === newTab);
    if (!hasOutput) {
        emitRun();
    }
});

const currentSnippet = computed({
    get() {
        return codeSnippets.value.find(snippet => snippet.tab === activeTab.value)?.code || ''
    },
    set(newCode: string) {
        const snippet = codeSnippets.value.find(snippet => snippet.tab === activeTab.value)
        if (snippet) {
            snippet.code = newCode
        }
    }
})

const editorOptions = {
    automaticLayout: true,
    minimap: { enabled: false },
    fontSize: 14,
    scrollBeyondLastLine: false,
    scrollbar: {
        vertical: 'hidden',
        horizontal: 'hidden',
    },
}

const emitRun = () => {
    if (!currentSnippet.value.trim()) return

    EventsEmit("runner:run", {
        code: currentSnippet.value,
        lang: "ts",
        tab: activeTab.value,
    })
}

const handleBeforeMount = (monaco: any) => {
    monaco.languages.typescript.typescriptDefaults.setCompilerOptions({
        target: monaco.languages.typescript.ScriptTarget.ESNext,
        allowNonTsExtensions: true,
        moduleResolution: monaco.languages.typescript.ModuleResolutionKind.NodeJs,
        module: monaco.languages.typescript.ModuleKind.ESNext,
        noEmit: true,
        esModuleInterop: true,
    });
}

const handleEditorMount = (editor: any) => {
    const updateHeight = () => {
        const contentHeight = editor.getContentHeight()
        editorHeight.value = contentHeight
        const layoutInfo = editor.getLayoutInfo()
        editor.layout({ width: layoutInfo.width, height: contentHeight })
    }

    updateHeight()
    editor.onDidContentSizeChange(updateHeight)

    editor.onDidChangeModelContent(() => {
        currentSnippet.value = editor.getValue()

        if (isProgrammatic) return;

        if (debounceTimer) clearTimeout(debounceTimer)

        debounceTimer = setTimeout(() => {
            emitRun()
        }, 500)
    })

    editor.onDidBlurEditorText(() => {
        if (debounceTimer) {
            clearTimeout(debounceTimer)
            debounceTimer = null
            emitRun()
        }
    })
}

onMounted(() => {
    emitRun()
})

onBeforeUnmount(() => {
    if (debounceTimer) clearTimeout(debounceTimer)
})
</script>

<style scoped>
@media (min-width: 768px) {
    .panel-left {
        width: var(--left-width) !important;
        flex: 0 0 var(--left-width);
    }
}

:deep(.monaco-editor),
:deep(.monaco-editor-background),
:deep(.monaco-editor .margin) {
    background-color: #21222C !important;
}


:deep(.monaco-editor .line-numbers) {
    color: #6272A4;
}
</style>