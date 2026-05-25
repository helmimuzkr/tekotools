<script lang="ts">
  import { Button } from "$lib/components/ui/button";

  type Props = {
    selected: string;
  };

  type Service = {
    name: string;
  };

  let services = $state<Service[]>([
    {
      name: "user-service",
    },
    {
      name: "product-service",
    },
  ]);

  let { selected }: Props = $props();
</script>

<nav class="flex flex-col w-52 h-full border-r p-2 gap-1">
  <p class="text-xs text-muted-foreground px-2 py-1">Services</p>
  {#each services as service (service.name)}
    <Button
      class="w-full justify-start"
      size={"sm"}
      variant="ghost"
      onclick={() => (selected = service.name)}
    >
      {service.name}
    </Button>
  {/each}
</nav>
