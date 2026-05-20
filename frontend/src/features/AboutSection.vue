<template>
  <div class="mt-8 border-t border-[#44475A] pt-6">
    <div class="flex items-start justify-between mb-4">
      <div>
        <h3 class="text-lg font-semibold text-[#F8F8F2]">RunBit</h3>
        <p class="text-sm text-[#6272A4]">A fast, modern JavaScript/TypeScript playground</p>
        <p class="text-xs text-[#6272A4] mt-1">
          Developed by <span class="text-[#8BE9FD]">Moises Santos Hernandez</span>
        </p>
      </div>
      <div class="text-right flex flex-col items-end gap-1">
        <span class="text-sm font-medium text-[#8BE9FD]">Version {{ currentVersion }}</span>
        <button @click="openRepo" class="text-xs text-[#BD93F9] hover:text-[#D6ACFF] hover:underline transition-colors flex items-center gap-1">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-3 h-3">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
          </svg>
          Repository
        </button>
      </div>
    </div>

    <div class="bg-[#21222C] rounded-md border border-[#44475A] p-4 flex flex-col gap-3">
      <div class="flex items-center justify-between">
        <div>
          <h4 class="text-sm font-medium text-[#F8F8F2]">Application Updates</h4>
          <p class="text-xs text-[#6272A4] mt-1">Check for the latest features and bug fixes.</p>
        </div>
        <button 
          @click="checkUpdate" 
          :disabled="isChecking"
          class="rounded-md bg-[#50FA7B] px-4 py-2 text-sm font-semibold text-[#282A36] hover:bg-[#69FF94] disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center gap-2"
        >
          <svg v-if="isChecking" class="animate-spin h-4 w-4 text-[#282A36]" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
          </svg>
          Check for updates
        </button>
      </div>

      <div v-if="updateResult" class="mt-2 p-3 rounded bg-[#282A36] border" :class="updateResult.hasUpdate ? 'border-[#50FA7B]' : 'border-[#44475A]'">
        <div v-if="updateResult.hasUpdate" class="flex flex-col gap-2">
          <p class="text-sm text-[#F8F8F2]">
            <span class="text-[#50FA7B] font-medium">New update available!</span> Version {{ updateResult.latest }} is ready to download.
          </p>
          <a :href="updateResult.downloadUrl" target="_blank" class="text-xs text-[#8BE9FD] hover:underline self-start">
            Download {{ updateResult.platform }} package
          </a>
        </div>
        <div v-else>
          <p class="text-sm text-[#F8F8F2]">You are up to date! RunBit is running the latest version.</p>
        </div>
      </div>
      
      <div v-if="errorMsg" class="mt-2 flex flex-col gap-2">
        <p class="text-xs text-[#FF5555]">Error: {{ errorMsg }}</p>
        <button 
          @click="openRepo" 
          class="text-xs text-[#8BE9FD] hover:text-[#50FA7B] hover:underline self-start flex items-center gap-1 transition-colors"
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-3 h-3">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
          </svg>
          Open GitHub repository instead
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetAppVersion, CheckUpdate } from '../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'

const currentVersion = ref('dev')
const isChecking = ref(false)
const updateResult = ref<any>(null)
const errorMsg = ref('')

onMounted(async () => {
  try {
    currentVersion.value = await GetAppVersion()
  } catch (e) {
    console.error('Failed to get app version:', e)
  }
})

const checkUpdate = async () => {
  if (isChecking.value) return
  
  isChecking.value = true
  updateResult.value = null
  errorMsg.value = ''
  
  try {
    const result = await CheckUpdate()
    updateResult.value = result
  } catch (e: any) {
    console.error('Update check failed:', e)
    errorMsg.value = String(e)
  } finally {
    isChecking.value = false
  }
}

const openRepo = () => {
  BrowserOpenURL('https://github.com/Moisessantos45/runbit')
}
</script>
