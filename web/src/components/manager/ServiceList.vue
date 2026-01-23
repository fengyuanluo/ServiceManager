<script setup lang="ts">
import type { Lang } from "../../i18n";
import type { ServiceSummary } from "../../types/service";
import StatusPill from "./StatusPill.vue";

const props = defineProps<{
  services: ServiceSummary[];
  selectedId: string | null;
  selectedIds: string[];
  lang: Lang;
}>();

const emit = defineEmits<{
  (event: "select", id: string): void;
  (event: "toggle", id: string, checked: boolean): void;
}>();
</script>

<template>
  <div class="space-y-2">
    <div
      v-for="service in props.services"
      :key="service.name"
      class="flex w-full items-start gap-3 rounded-xl border border-border/60 px-4 py-3 text-left transition"
      :class="
        props.selectedId === service.name
          ? 'bg-muted/80 text-fg shadow-sm shadow-black/5'
          : 'bg-card/70 text-fg/70 hover:bg-muted/60 hover:text-fg'
      "
    >
      <input
        :data-testid="`service-check-${service.name}`"
        type="checkbox"
        class="mt-1 h-3.5 w-3.5 rounded border-border/60 bg-card/70 text-primary"
        :checked="props.selectedIds.includes(service.name)"
        @change="emit('toggle', service.name, ($event.target as HTMLInputElement).checked)"
      />
      <button
        :data-testid="`service-item-${service.name}`"
        type="button"
        class="flex-1 text-left"
        @click="emit('select', service.name)"
      >
        <div class="flex items-center justify-between gap-3">
          <span class="text-sm font-semibold tracking-wide">{{ service.name }}</span>
          <StatusPill :status="service.status" :lang="props.lang" />
        </div>
        <div class="mt-1 text-xs text-fg/40">{{ service.detail ?? "-" }}</div>
      </button>
    </div>
  </div>
</template>
