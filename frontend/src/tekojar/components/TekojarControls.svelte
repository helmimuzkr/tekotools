<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";

  type Props = {
    selected: string;
  };

  let { selected }: Props = $props();
</script>

<div class="flex justify-between items-center">
  <div class="flex items-center gap-2">
    <span class="font-bold">{selected}</span>
    <Badge>Running</Badge>
  </div>
  <div class="flex gap-2">
    <Button variant="default" size="sm">Start</Button>
    <Button variant="destructive" size="sm">Stop</Button>
    <Button variant="outline" size="sm">Restart</Button>
  </div>
</div>
