<template>
  <div class="h-screen w-screen overflow-hidden bg-[#282A36] text-[#F8F8F2] flex flex-col">
    <TabsBar />
    <main class="relative w-full flex-1 flex flex-col gap-2 lg:flex-row p-2 overflow-hidden">
      <EditorPanel :leftWidth="leftWidth" />
      <div
        class="divider-drag hidden md:flex items-center justify-center cursor-col-resize w-1.5 bg-transparent hover:bg-[#44475A]/30 z-10 shrink-0 transition-colors rounded"
        @mousedown="startDrag">
        <div class="w-px h-[30%] bg-[#44475A] pointer-events-none rounded-full"></div>
      </div>
      <PreviewPane />
    </main>

    <div v-if="isInstallingBun"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm">
      <div class="flex flex-col items-center gap-4 text-[#F8F8F2]">
        <svg class="animate-spin h-10 w-10 text-[#BD93F9]" xmlns="http://www.w3.org/2000/svg" fill="none"
          viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
          </path>
        </svg>
        <span class="text-lg font-medium">Installing Bun... Please wait.</span>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { EventsOn } from '../wailsjs/runtime/runtime'
import { InstallBunFlow } from '../wailsjs/go/main/App'
import TabsBar from './features/TabsBar.vue'
import EditorPanel from './features/EditorPanel.vue'
import PreviewPane from './features/PreviewPane.vue'

const leftWidth = ref(300)
let dragging = false
let startX = 0
let startWidth = 0
let rafId: number | null = null
let latestX = 0

const onMouseMove = (e: MouseEvent) => {
  if (!dragging) return
  latestX = e.clientX
  if (rafId !== null) return
  rafId = window.requestAnimationFrame(() => {
    const deltaX = latestX - startX
    leftWidth.value = Math.max(150, Math.min(startWidth + deltaX, window.innerWidth - 150))
    rafId = null
  })
}

const stopDrag = () => {
  dragging = false
  if (rafId !== null) {
    window.cancelAnimationFrame(rafId)
    rafId = null
  }
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', stopDrag)
  document.body.style.userSelect = ''
  document.body.style.cursor = ''
}

const startDrag = (e: MouseEvent) => {
  dragging = true
  startX = e.clientX
  startWidth = leftWidth.value
  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', stopDrag)
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'col-resize'
}


const handleResize = () => {
  if (window.innerWidth <= 768) {
    leftWidth.value = window.innerWidth
  } else if (!dragging) {
    leftWidth.value = window.innerWidth / 2
  }
};

let offBunMissing: () => void;
let offBunInstalling: () => void;
let offBunReady: () => void;
let offBunError: () => void;

const isInstallingBun = ref(false)

onMounted(() => {
  if (window.innerWidth > 768) {
    leftWidth.value = window.innerWidth / 2;
  }
  window.addEventListener("resize", handleResize);

  offBunMissing = EventsOn("runner:bun-missing", async () => {
    try {
      await InstallBunFlow();
    } catch (e) {
      console.error("Failed to start bun install flow:", e);
    }
  });

  offBunInstalling = EventsOn("runner:installing-bun", () => {
    isInstallingBun.value = true;
  });

  offBunReady = EventsOn("runner:ready", () => {
    isInstallingBun.value = false;
  });

  offBunError = EventsOn("runner:error", () => {
    isInstallingBun.value = false;
  });
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", handleResize);
  window.removeEventListener('mousemove', onMouseMove)
  window.removeEventListener('mouseup', stopDrag)
  if (rafId !== null) {
    window.cancelAnimationFrame(rafId)
  }

  if (offBunMissing) offBunMissing();
  if (offBunInstalling) offBunInstalling();
  if (offBunReady) offBunReady();
  if (offBunError) offBunError();
})
</script>

<style scoped>
/* Metadata panel slide transition */
.meta-slide-enter-active,
.meta-slide-leave-active {
  transition: max-height 0.22s ease, opacity 0.18s ease, padding 0.22s ease;
  overflow: hidden;
  max-height: 120px;
}

.meta-slide-enter-from,
.meta-slide-leave-to {
  max-height: 0;
  opacity: 0;
  padding-top: 0;
  padding-bottom: 0;
}

.tabs-scroll {
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.tabs-scroll::-webkit-scrollbar {
  width: 0;
  height: 0;
}
</style>