<template>
    <header style="widows: 1; -webkit-app-region: drag;"
        class="flex h-11 w-full shrink-0 items-center justify-between gap-3 border-b border-[#14161d] bg-gradient-to-b from-[#20222d] to-[#1b1d27] pr-2 select-none">
        <nav class="tabs-scroll flex h-full min-w-0 flex-1 items-stretch overflow-x-auto overflow-y-hidden whitespace-nowrap">
            <button type="button" v-for="tab in tabs" :key="tab.tab"
                style="widows: 2; -webkit-app-region: no-drag;" :class="[
                    'group relative flex h-full min-w-[132px] max-w-[220px] cursor-pointer items-center gap-2 border-r border-[#191A21] px-4 text-[13px] transition-all duration-150',
                    activeTab === tab.tab
                        ? 'bg-[#282A36] text-[#F8F8F2]'
                        : 'bg-transparent text-[#7c86b2] hover:bg-[#282A36]/60 hover:text-[#F8F8F2]'
                ]" @click="playgroundStore.setTabActive(tab.tab)">
                <div v-if="activeTab === tab.tab" class="absolute inset-x-0 top-0 h-0.5 bg-[#BD93F9]" />

                <svg class="h-4 w-4 shrink-0 text-[#8BE9FD]" viewBox="0 0 24 24" fill="currentColor">
                    <path
                        d="M4 3C2.89543 3 2 3.89543 2 5V19C2 20.1046 2.89543 21 4 21H20C21.1046 21 22 20.1046 22 19V5C22 3.89543 21.1046 3 20 3H4ZM12.0001 16.5C10.7416 16.5 9.7719 15.9329 9.30907 15.0067L10.7684 14.1685C11.0264 14.7334 11.4554 15.011 12.0298 15.011C12.6343 15.011 12.9818 14.6734 12.9818 14.2862C12.9818 13.8096 12.4259 13.621 11.5323 13.3827C10.2316 13.0451 8.98064 12.4495 8.98064 10.9701C8.98064 9.61002 10.0962 8.5 11.8512 8.5C13.0625 8.5 14.0157 9.00645 14.4725 9.87023L13.1023 10.7243C12.8342 10.1584 12.457 9.94002 11.8512 9.94002C11.3151 9.94002 10.8782 10.2378 10.8782 10.7544C10.8782 11.3402 11.5932 11.4593 12.457 11.7274C13.8471 12.1345 14.9392 12.75 14.9392 14.2465C14.9392 15.7058 13.7875 16.5 12.0001 16.5ZM16.0392 8.65087H20.5V10.1599H19.0605V16.3312H17.3132V10.1599H16.0392V8.65087Z" />
                </svg>

                <span class="min-w-0 flex-1 truncate text-left font-medium">App {{ tab.tab }}</span>

                <div @click.stop="playgroundStore.removeTab(tab.tab)"
                    style="widows: 2; -webkit-app-region: no-drag;" :class="[
                        'shrink-0 rounded-md p-0.5 transition-all duration-150',
                        tab.tab !== 1
                            ? 'text-[#6272A4] opacity-0 group-hover:opacity-100 hover:bg-[#44475A] hover:text-[#FF5555]'
                            : 'pointer-events-none cursor-default opacity-0'
                    ]">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                        stroke="currentColor" class="h-3.5 w-3.5">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </div>
            </button>

            <button type="button" @click="playgroundStore.createTab"
                style="widows: 2; -webkit-app-region: no-drag;"
                class="ml-2 my-1.5 inline-flex h-8 w-8 items-center justify-center rounded-md text-[#6272A4] transition-colors hover:bg-[#44475A] hover:text-[#F8F8F2]"
                title="New Tab">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                    stroke="currentColor" class="h-4 w-4">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                </svg>
            </button>
        </nav>

        <div style="widows: 2; -webkit-app-region: no-drag;" class="flex items-center gap-2 pl-2">
            <button type="button" @click="isSettingsOpen = true"
                class="inline-flex h-8 w-8 items-center justify-center rounded-md text-[#6272A4] transition-colors hover:bg-[#44475A] hover:text-[#F8F8F2]"
                title="Settings">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                    stroke="currentColor" class="h-5 w-5">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M10.343 3.94c.09-.542.56-.94 1.11-.94h1.093c.55 0 1.02.398 1.11.94l.149.894c.07.424.384.764.78.93.398.164.855.142 1.205-.108l.737-.527a1.125 1.125 0 011.45.12l.773.774c.39.389.44 1.002.12 1.45l-.527.737c-.25.35-.272.806-.107 1.204.165.397.505.71.93.78l.893.15c.543.09.94.56.94 1.109v1.094c0 .55-.397 1.02-.94 1.11l-.893.149c-.425.07-.765.383-.93.78-.165.398-.143.854.107 1.204l.527.738c.32.447.269 1.06-.12 1.45l-.774.773a1.125 1.125 0 01-1.449.12l-.738-.527c-.35-.25-.806-.272-1.203-.107-.397.165-.71.505-.781.929l-.149.894c-.09.542-.56.94-1.11.94h-1.094c-.55 0-1.019-.398-1.11-.94l-.148-.894c-.071-.424-.384-.764-.781-.93-.398-.164-.854-.142-1.204.108l-.738.527c-.447.32-1.06.269-1.45-.12l-.773-.774a1.125 1.125 0 01-.12-1.45l.527-.737c.25-.35.273-.806.108-1.204-.165-.397-.505-.71-.93-.78l-.894-.15c-.542-.09-.94-.56-.94-1.109v-1.094c0-.55.398-1.02.94-1.11l.894-.149c.424-.07.765-.383.93-.78.165-.398.143-.854-.107-1.204l-.527-.738a1.125 1.125 0 01.12-1.45l.773-.773a1.125 1.125 0 011.45-.12l.737.527c.35.25.807.272 1.204.107.397-.165.71-.505.78-.929l.15-.894z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
            </button>

            <div class="flex items-center self-stretch border-l border-[#2a2d39] pl-2">
                <button type="button" @click="WindowMinimise()"
                    class="inline-flex h-full w-10 items-center justify-center text-sm text-[#7c86b2] transition-colors hover:bg-[#2c3040] hover:text-[#F8F8F2]">
                    ─
                </button>
                <button type="button" @click="WindowToggleMaximise()"
                    class="inline-flex h-full w-10 items-center justify-center text-[12px] text-[#7c86b2] transition-colors hover:bg-[#2c3040] hover:text-[#F8F8F2]">
                    □
                </button>
                <button type="button" @click="Quit()"
                    class="inline-flex h-full w-10 items-center justify-center text-sm text-[#c7c9d3] transition-colors hover:bg-[#FF5555] hover:text-white">
                    ✕
                </button>
            </div>
        </div>

        <SettingsModal v-model:isOpen="isSettingsOpen" />
    </header>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import usePlaygroundStore from '@/store/Playground'
import SettingsModal from './SettingsModal.vue'
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime'

const playgroundStore = usePlaygroundStore()
const { tabs, activeTab } = storeToRefs(playgroundStore)

const isSettingsOpen = ref(false)

</script>

<style scoped>
.tabs-scroll::-webkit-scrollbar {
    height: 0;
}

.tabs-scroll {
    scrollbar-width: none;
}
</style>