<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { t, type Lang } from "../../i18n";
import type { ManagerKey } from "../../data/navigation";
import { getServiceDetail, getServiceLogs, getServiceSummaries } from "../../data/services";
import BatchToolbar from "./BatchToolbar.vue";
import ServiceDetail from "./ServiceDetail.vue";
import ServiceEditor from "./ServiceEditor.vue";
import ServiceList from "./ServiceList.vue";
import LogViewer from "./LogViewer.vue";

const props = defineProps<{
  managerKey: ManagerKey;
  lang: Lang;
}>();

const selectedId = ref<string | null>(null);
const selectedIds = ref<string[]>([]);

const services = computed(() => getServiceSummaries(props.managerKey));
const detail = computed(() =>
  selectedId.value ? getServiceDetail(props.managerKey, selectedId.value) : null
);
const logs = computed(() =>
  selectedId.value ? getServiceLogs(props.managerKey, selectedId.value) : []
);
const allSelected = computed(
  () => services.value.length > 0 && selectedIds.value.length === services.value.length
);

watch(
  () => props.managerKey,
  () => {
    selectedId.value = null;
    selectedIds.value = [];
  }
);

const selectService = (id: string) => {
  selectedId.value = id;
};

const toggleSelection = (id: string, checked: boolean) => {
  if (checked) {
    if (!selectedIds.value.includes(id)) {
      selectedIds.value = [...selectedIds.value, id];
    }
    return;
  }
  selectedIds.value = selectedIds.value.filter((value) => value !== id);
};

const toggleAll = (checked: boolean) => {
  selectedIds.value = checked ? services.value.map((service) => service.id) : [];
};

const handleBatchAction = (_action: string) => {
  // placeholder for future API integration
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
      <BatchToolbar
        :selected-count="selectedIds.length"
        :total-count="services.length"
        :all-selected="allSelected"
        :lang="props.lang"
        @toggle-all="toggleAll"
        @action="handleBatchAction"
      />
      <ServiceList
        :services="services"
        :selected-id="selectedId"
        :selected-ids="selectedIds"
        :lang="props.lang"
        @select="selectService"
        @toggle="toggleSelection"
      />
    </section>

    <section class="space-y-6">
      <div class="rounded-2xl border border-border/60 bg-card/70 p-4 shadow-sm">
        <ServiceDetail :service="detail" :lang="props.lang" />
      </div>
      <div class="rounded-2xl border border-border/60 bg-card/70 p-4 shadow-sm">
        <LogViewer :lines="logs" :lang="props.lang" :active="!!detail" />
      </div>
      <div class="rounded-2xl border border-border/60 bg-card/70 p-4 shadow-sm">
        <ServiceEditor
          :service="detail"
          :manager-key="props.managerKey"
          :lang="props.lang"
        />
      </div>
    </section>
  </div>
</template>
