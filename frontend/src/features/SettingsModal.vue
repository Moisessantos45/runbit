<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center bg-black/55 backdrop-blur-[3px">
    <div
      class="flex max-h-[80vh] w-full max-w-lg flex-col overflow-hidden rounded-2xl border border-[#44475A] bg-[#282A36] shadow-[0_24px_80px_rgba(0,0,0,0.45)]">
      <div class="flex items-center justify-between border-b border-[#44475A] px-4 py-3 shrink-0">
        <h2 class="text-lg font-semibold text-[#F8F8F2]">Packages Settings</h2>
        <button @click="close"
          class="rounded-md p-1 text-[#6272A4] hover:bg-[#44475A] hover:text-[#FF5555] transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"
            stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-4 flex-1 overflow-y-auto">
        <div class="mb-6">
          <label class="block text-sm font-medium text-[#8BE9FD] mb-2">Install a dependency</label>
          <div class="flex gap-2">
            <input v-model="newPackage" @keyup.enter="installPackage" type="text" placeholder="e.g. lodash"
              class="flex-1 rounded-md border border-[#44475A] bg-[#21222C] px-3 py-2 text-sm text-[#F8F8F2] placeholder-[#6272A4] focus:border-[#BD93F9] focus:outline-none focus:ring-1 focus:ring-[#BD93F9]"
              :disabled="isLoading" />
            <button @click="installPackage" :disabled="!newPackage.trim() || isLoading"
              class="rounded-md bg-[#BD93F9] px-4 py-2 text-sm font-semibold text-[#282A36] hover:bg-[#D6ACFF] disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center gap-2">
              <svg v-if="isInstalling" class="animate-spin h-4 w-4 text-[#282A36]" xmlns="http://www.w3.org/2000/svg"
                fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                </path>
              </svg>
              Install
            </button>
          </div>
          <p v-if="errorMsg" class="mt-2 text-xs text-[#FF5555]">{{ errorMsg }}</p>
        </div>

        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="block text-sm font-medium text-[#8BE9FD]">Installed Packages</label>
            <button @click="fetchPackages" class="text-xs text-[#6272A4] hover:text-[#F8F8F2] transition-colors"
              title="Refresh">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                stroke="currentColor" class="w-4 h-4" :class="{ 'animate-spin': isFetching }">
                <path stroke-linecap="round" stroke-linejoin="round"
                  d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
              </svg>
            </button>
          </div>

          <div v-if="isFetching && packages.length === 0" class="flex justify-center py-4">
            <svg class="animate-spin h-5 w-5 text-[#BD93F9]" xmlns="http://www.w3.org/2000/svg" fill="none"
              viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
              </path>
            </svg>
          </div>
          <ul v-else-if="packages.length > 0"
            class="divide-y divide-[#44475A] border border-[#44475A] rounded-md bg-[#21222C]">
            <li v-for="pkg in packages" :key="pkg.name" class="flex items-center justify-between px-3 py-2">
              <div>
                <span class="text-sm font-medium text-[#F8F8F2]">{{ pkg.name }}</span>
                <span class="ml-2 text-xs text-[#6272A4]">{{ pkg.version }}</span>
              </div>
              <button @click="removePackage(pkg.name)" :disabled="isLoading"
                class="rounded p-1 text-[#6272A4] hover:bg-[#44475A] hover:text-[#FF5555] disabled:opacity-50 transition-colors"
                title="Remove package">
                <svg v-if="removingPkg === pkg.name" class="animate-spin h-4 w-4 text-[#FF5555]"
                  xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                  </path>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2"
                  stroke="currentColor" class="w-4 h-4">
                  <path stroke-linecap="round" stroke-linejoin="round"
                    d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                </svg>
              </button>
            </li>
          </ul>
          <div v-else class="text-sm text-[#6272A4] py-4 text-center border border-[#44475A] border-dashed rounded-md">
            No packages installed yet.
          </div>
        </div>

        <AboutSection />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { AddPackage, RemovePackage, ListInstalledPackages } from '../../wailsjs/go/main/App'
import AboutSection from './AboutSection.vue'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits(['update:isOpen'])

const newPackage = ref('')
const packages = ref<any[]>([])
const isFetching = ref(false)
const isInstalling = ref(false)
const removingPkg = ref('')
const errorMsg = ref('')

const isLoading = computed(() => isInstalling.value || removingPkg.value !== '')

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    fetchPackages()
    errorMsg.value = ''
    newPackage.value = ''
  }
})

const close = () => {
  emit('update:isOpen', false)
}

const fetchPackages = async () => {
  try {
    isFetching.value = true
    errorMsg.value = ''
    const result = await ListInstalledPackages()
    packages.value = result || []
  } catch (e: any) {
    console.error('Error fetching packages:', e)
    errorMsg.value = String(e)
  } finally {
    isFetching.value = false
  }
}

const installPackage = async () => {
  const pkg = newPackage.value.trim()
  if (!pkg || isLoading.value) return

  try {
    isInstalling.value = true
    errorMsg.value = ''
    await AddPackage(pkg)
    newPackage.value = ''
    await fetchPackages()
  } catch (e: any) {
    console.error('Error installing package:', e)
    errorMsg.value = String(e)
  } finally {
    isInstalling.value = false
  }
}

const removePackage = async (name: string) => {
  if (isLoading.value) return

  try {
    removingPkg.value = name
    errorMsg.value = ''
    await RemovePackage(name)
    await fetchPackages()
  } catch (e: any) {
    console.error('Error removing package:', e)
    errorMsg.value = String(e)
  } finally {
    removingPkg.value = ''
  }
}

import { computed } from 'vue'
</script>
