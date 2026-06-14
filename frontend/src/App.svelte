<script lang="ts">
  import { onMount } from "svelte";
  import { ModeWatcher } from "mode-watcher";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import Braces from "@lucide/svelte/icons/braces";
  import Monitor from "@lucide/svelte/icons/monitor";
  import SettingsIcon from "@lucide/svelte/icons/settings";
  import type { Page } from "./type";
  import AppSidebar from "@/shared/components/AppSidebar.svelte";
  import { Tekojar, initTekojar } from "./features/tekojar";
  import Setting from "./features/settings";
  import JsonataQuery from "./features/jsonata_query/JsonataQuery.svelte";
  import { Toaster } from "$lib/components/ui/sonner";

  const pages: Page[] = [
    {
      id: "tekojar",
      title: "Tekojar",
      icon: Monitor,
      section: "content",
      component: Tekojar,
      onInit: initTekojar,
    },
    {
      id: "jsonata",
      title: "JSONata",
      icon: Braces,
      section: "content",
      component: JsonataQuery,
      onInit: null,
    },
    {
      id: "settings",
      title: "Settings",
      icon: SettingsIcon,
      section: "footer",
      component: Setting,
      onInit: null,
    },
  ];

  onMount(() => {
    pages.forEach((p) => {
      p.onInit?.();
    });
  });

  let currentPage = $state<Page>(pages[0]);

  function handlePage(page: Page) {
    currentPage = page;
  }
</script>

<!-- Dark Mode Whatcher -->
<ModeWatcher />

<!-- Toast -->
<Toaster position="top-center" richColors={true} />

<div class="flex h-screen">
  <Sidebar.Provider open={false}>
    <AppSidebar {pages} {currentPage} onPageChange={handlePage} />
    <Sidebar.Inset class="flex flex-col h-full min-h-0 min-w-0 ">
      {#if currentPage.component}
        <currentPage.component />
      {/if}
    </Sidebar.Inset>
  </Sidebar.Provider>
</div>
