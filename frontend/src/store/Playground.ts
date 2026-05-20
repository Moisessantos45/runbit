import { defineStore } from "pinia";
import { reactive, ref } from "vue";

type TabItem = {
  tab: number;
  active: boolean;
};

type CodeSnippet = {
  tab: number;
  code: string;
};

type OutputLine = {
  event: string;
  result: string;
  tab: number;
};

const usePlaygroundStore = defineStore("playground", () => {
  const tabs = reactive<TabItem[]>([{ tab: 1, active: true }]);

  const activeTab = ref(1);

  const codeSnippets = reactive<CodeSnippet[]>([
    {
      tab: 1,
      code: `function greet(name: string) {\n  return \`Hello, \${name}!\`;\n}\n\nconsole.log(greet("World"));`,
    },
  ]);

  const stdoutLines = reactive<OutputLine[]>([]);
  const loadingByTab = ref<Record<number, boolean>>({ 1: false });

  const createTab = () => {
    const nextTab =
      tabs.length > 0 ? Math.max(...tabs.map((item) => item.tab)) + 1 : 1;

    for (const item of tabs) {
      item.active = false;
    }

    tabs.push({
      tab: nextTab,
      active: true,
    });

    codeSnippets.push({
      tab: nextTab,
      code: `console.log("New tab");`,
    });

    loadingByTab.value[nextTab] = false;

    activeTab.value = nextTab;
  };

  const setTabActive = (tab: number) => {
    activeTab.value = tab;

    for (const item of tabs) {
      item.active = item.tab === tab;
    }
  };

  const setTabLoading = (tab: number, isLoading: boolean) => {
    loadingByTab.value[tab] = isLoading;
  };

  const removeTab = (tab: number) => {
    if (tabs.length === 1) return;

    const tabIndex = tabs.findIndex((item) => item.tab === tab);
    if (tabIndex === -1) return;

    tabs.splice(tabIndex, 1);

    const snippetIndex = codeSnippets.findIndex((item) => item.tab === tab);
    if (snippetIndex !== -1) {
      codeSnippets.splice(snippetIndex, 1);
    }

    for (let i = stdoutLines.length - 1; i >= 0; i--) {
      if (stdoutLines[i].tab === tab) {
        stdoutLines.splice(i, 1);
      }
    }

    delete loadingByTab.value[tab];

    if (activeTab.value === tab && tabs.length > 0) {
      setTabActive(tabs[tabs.length - 1].tab);
    }
  };

  return {
    tabs,
    activeTab,
    codeSnippets,
    stdoutLines,
    loadingByTab,
    createTab,
    setTabActive,
    setTabLoading,
    removeTab,
  };
});

export default usePlaygroundStore;
