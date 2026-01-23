<script setup lang="ts">
import type { Lang } from "../../i18n";
import type { ServiceSummary } from "../../types/service";
import StatusPill from "./StatusPill.vue";

const props = defineProps<{
  services: ServiceSummary[];
  selectedId: string | null;
  lang: Lang;
}>();

const emit = defineEmits<{
  (event: "select", id: string): void;
}>();
</script>

<template>
  <div class="space-y-2">
    <button
      v-for="service in props.services"
      :key="service.id"
      :data-testid="`service-item-${service.id}`"
      type="button"
      class="w-full rounded-xl border border-border/60 px-4 py-3 text-left transition"
      :class="
        props.selectedId === service.id
          ? 'bg-muted/80 text-fg shadow-sm shadow-black/5'
          : 'bg-card/70 text-fg/70 hover:bg-muted/60 hover:text-fg'
      "
      @click="emit('select', service.id)"
    >
      <div class="flex items-center justify-between gap-3">
        <span class="text-sm font-semibold tracking-wide">{{ service.name }}</span>
        <StatusPill :status="service.status" :lang="props.lang" />
      </div>
      <div class="mt-1 text-xs text-fg/40">{{ service.summary }}</div>
    </button>
  </div>
</template>
