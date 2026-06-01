<script lang="ts">
  import Button from "$lib/components/ui/button/button.svelte";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import { useSidebar } from "$lib/components/ui/sidebar/index.js";
  import { Trigger } from "$lib/components/ui/tooltip";

  import type { Page } from "./type";

  type Props = {
    pages: Page[];
    currentPage: Page;
    onPageChange: (page: Page) => void;
  };

  let { pages, currentPage, onPageChange }: Props = $props();

  const sidebar = useSidebar();
</script>

<Sidebar.Root collapsible="icon">
  <Sidebar.Header>
    <!-- <Button variant="ghost" onclick={() => sidebar.toggle()}>☰</Button> -->
    <Sidebar.Trigger />
  </Sidebar.Header>

  <Sidebar.Content>
    <Sidebar.Group>
      <Sidebar.GroupContent>
        <Sidebar.Menu>
          {#each pages as page (page.id)}
            {#if page.section === "content"}
              <Sidebar.MenuItem>
                <Sidebar.MenuButton isActive={currentPage.id === page.id}>
                  {#snippet child({ props })}
                    <button onclick={() => onPageChange(page)} {...props}>
                      <page.icon />
                      <span>{page.title}</span>
                    </button>
                  {/snippet}
                </Sidebar.MenuButton>
              </Sidebar.MenuItem>
            {/if}
          {/each}
        </Sidebar.Menu>
      </Sidebar.GroupContent>
    </Sidebar.Group>
  </Sidebar.Content>

  <Sidebar.Footer>
    {#each pages as page (page.id)}
      {#if page.section === "footer"}
        <Sidebar.MenuItem>
          <Sidebar.MenuButton isActive={currentPage.id === page.id}>
            {#snippet child({ props })}
              <button onclick={() => onPageChange(page)} {...props}>
                <page.icon />
                <span>{page.title}</span>
              </button>
            {/snippet}
          </Sidebar.MenuButton>
        </Sidebar.MenuItem>
      {/if}
    {/each}
  </Sidebar.Footer>
</Sidebar.Root>
