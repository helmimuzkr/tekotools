<script lang="ts">
  import { ModeWatcher } from "mode-watcher";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import Braces from "@lucide/svelte/icons/braces";
  import Monitor from "@lucide/svelte/icons/monitor";
  import SettingsIcon from "@lucide/svelte/icons/settings";

  import type { Page } from "./type";
  import AppSidebar from "./component/AppSidebar.svelte";
  import Tekojar from "./tekojar";
  import JsonataQuery from "./jsonata_query/JsonataQuery.svelte";
  import Settings from "./settings";

  const pages: Page[] = [
    {
      id: "tekojar",
      title: "Tekojar",
      icon: Monitor,
      section: "content",
      component: Tekojar,
    },
    {
      id: "jsonata",
      title: "JSONata",
      icon: Braces,
      section: "content",
      component: JsonataQuery,
    },
    {
      id: "settings",
      title: "Settings",
      icon: SettingsIcon,
      section: "footer",
      component: Settings,
    },
  ];

  let currentPage = $state<Page>(pages[0]);

  function handlePage(page: Page) {
    currentPage = page;
  }
</script>

<!-- Dark Mode Whatcher -->
<ModeWatcher />

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
