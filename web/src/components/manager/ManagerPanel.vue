<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { t, type Lang } from "../../i18n";
import type { ManagerKey } from "../../data/navigation";
import { getServiceDetail, getServiceSummaries } from "../../data/services";
import ServiceDetail from "./ServiceDetail.vue";
import ServiceList from "./ServiceList.vue";

const props = defineProps<{
  managerKey: ManagerKey;
  lang: Lang;
}>();

const selectedId = ref<string | null>(null);

const services = computed(() => getServiceSummaries(props.managerKey));
const detail = computed(() =>
  selectedId.value ? getServiceDetail(props.managerKey, selectedId.value) : null
);

watch(
  () => props.managerKey,
  () => {
    selectedId.value = null;
  }
);

const selectService = (id: string) => {
  selectedId.value = id;
};
</script>

<template>
  <div class="grid gap-6 lg:grid-cols-[360px_1fr]">
    <section class="rounded-2xl border border-border/60 bg-card/70 p-4 shadow-sm">
      <header class="flex items-center justify-between pb-3">
        <span class="text-[11px] uppercase tracking-[0.35em] text-fg/40">
          {{ t("layout.list", props.lang) }}
        </span>
        <span class="text-[10px] uppercase tracking-[0.25em] text-fg/30">
          {{ services.length }}
        </span>
      </header>
      <ServiceList
        :services="services"
        :selected-id="selectedId"
        :lang="props.lang"
        @select="selectService"
      />
    </section>

    <section class="rounded-2xl border border-border/60 bg-card/70 p-4 shadow-sm">
      <ServiceDetail :service="detail" :lang="props.lang" />
    </section>
  </div>
</template>
