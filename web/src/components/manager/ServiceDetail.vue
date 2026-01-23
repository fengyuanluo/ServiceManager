<script setup lang="ts">
import { t, type Lang } from "../../i18n";
import type { ServiceDetail } from "../../types/service";
import StatusPill from "./StatusPill.vue";

const props = defineProps<{
  service: ServiceDetail | null;
  lang: Lang;
}>();
</script>

<template>
  <div v-if="!props.service" data-testid="detail-empty" class="h-full">
    <div
      class="flex h-full items-center justify-center rounded-2xl border border-dashed border-border/60 bg-card/50 px-6 py-16 text-xs uppercase tracking-[0.35em] text-fg/40"
    >
      {{ t("layout.select", props.lang) }}
    </div>
  </div>
  <div v-else class="space-y-4">
    <div class="flex items-center justify-between gap-4">
      <div>
        <div class="text-xs uppercase tracking-[0.35em] text-fg/40">
          {{ t("layout.detail", props.lang) }}
        </div>
        <div
          data-testid="detail-name"
          class="mt-2 font-display text-xl tracking-wide text-fg"
        >
          {{ props.service.name }}
        </div>
      </div>
      <StatusPill :status="props.service.status" :lang="props.lang" />
    </div>

    <div class="text-sm text-fg/60">
      {{ props.service.description || props.service.summary }}
    </div>

    <dl class="grid grid-cols-2 gap-4 text-xs text-fg/60">
      <div class="space-y-1">
        <dt class="text-[10px] uppercase tracking-[0.25em] text-fg/40">
          {{ t("detail.pid", props.lang) }}
        </dt>
        <dd>{{ props.service.pid ?? "-" }}</dd>
      </div>
      <div class="space-y-1">
        <dt class="text-[10px] uppercase tracking-[0.25em] text-fg/40">
          {{ t("detail.since", props.lang) }}
        </dt>
        <dd>{{ props.service.since ?? "-" }}</dd>
      </div>
      <div class="space-y-1 col-span-2">
        <dt class="text-[10px] uppercase tracking-[0.25em] text-fg/40">
          {{ t("detail.unit", props.lang) }}
        </dt>
        <dd class="break-all">{{ props.service.unitPath ?? "-" }}</dd>
      </div>
    </dl>
  </div>
</template>
