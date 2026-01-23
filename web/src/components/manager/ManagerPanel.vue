<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { t, type Lang } from "../../i18n";
import type { ManagerKey } from "../../data/navigation";
import type { ServiceDetail as ServiceDetailType, ServiceSummary } from "../../types/service";
import { ApiError } from "../../api/client";
import { getServiceDetail, getServiceLogs, installManager, listServices, runServiceBatch } from "../../api/services";
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

const services = ref<ServiceSummary[]>([]);
const detail = ref<ServiceDetailType | null>(null);
const logs = ref<string[]>([]);
const pm2User = ref("root");
const missing = ref(false);
const installState = ref<"idle" | "installing" | "failed">("idle");

const currentUser = computed(() =>
  props.managerKey === "pm2" ? pm2User.value.trim() || "root" : ""
);
const installing = computed(() => installState.value === "installing");
const allSelected = computed(
  () => services.value.length > 0 && selectedIds.value.length === services.value.length
);

const isManagerUnavailable = (error: unknown) =>
  error instanceof ApiError && error.message.toLowerCase().includes("manager unavailable");

const loadList = async () => {
  try {
    services.value = await listServices(props.managerKey, currentUser.value);
    missing.value = false;
    installState.value = "idle";
  } catch (error) {
    console.error(error);
    services.value = [];
    detail.value = null;
    logs.value = [];
    missing.value = isManagerUnavailable(error);
    if (!missing.value) {
      installState.value = "idle";
    }
  }
  const nameSet = new Set(services.value.map((service) => service.name));
  selectedIds.value = selectedIds.value.filter((name) => nameSet.has(name));
  if (selectedId.value && !nameSet.has(selectedId.value)) {
    selectedId.value = null;
    detail.value = null;
    logs.value = [];
  }
};

const loadDetail = async (name: string) => {
  const summary = services.value.find((service) => service.name === name);
  try {
    const payload = await getServiceDetail(props.managerKey, name, currentUser.value);
    if (summary) {
      detail.value = {
        ...summary,
        ...payload,
        detail: payload.detail ?? summary.detail
      };
      return;
    }
    detail.value = payload;
  } catch (error) {
    console.error(error);
    detail.value = summary ?? null;
  }
};

const loadLogs = async (name: string) => {
  try {
    logs.value = await getServiceLogs(props.managerKey, name, {
      lines: 500,
      user: currentUser.value
    });
  } catch (error) {
    console.error(error);
    logs.value = [];
  }
};

watch(
  () => props.managerKey,
  () => {
    selectedId.value = null;
    selectedIds.value = [];
    detail.value = null;
    logs.value = [];
    void loadList();
  },
  { immediate: true }
);

watch(pm2User, () => {
  if (props.managerKey !== "pm2") return;
  selectedId.value = null;
  selectedIds.value = [];
  detail.value = null;
  logs.value = [];
  void loadList();
});

const selectService = (id: string) => {
  selectedId.value = id;
  detail.value = null;
  logs.value = [];
  void loadDetail(id);
  void loadLogs(id);
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
  selectedIds.value = checked ? services.value.map((service) => service.name) : [];
};

const handleBatchAction = async (action: string) => {
  if (selectedIds.value.length === 0) return;
  const items = selectedIds.value.map((name) => ({ name, action }));
  try {
    await runServiceBatch(props.managerKey, items, currentUser.value);
  } catch (error) {
    console.error(error);
  } finally {
    void loadList();
  }
};

const handleInstall = async () => {
  if (installing.value) return;
  if (typeof window !== "undefined") {
    const confirmed = window.confirm(t("install.confirm", props.lang));
    if (!confirmed) return;
  }
  installState.value = "installing";
  try {
    await installManager(props.managerKey);
    await loadList();
    if (missing.value) {
      installState.value = "failed";
    }
  } catch (error) {
    console.error(error);
    installState.value = "failed";
  }
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
      <div
        v-if="missing"
        data-testid="install-card"
        class="mb-4 rounded-2xl border border-amber-300/60 bg-amber-200/10 p-3 text-xs text-fg/70 shadow-sm shadow-amber-200/30"
      >
        <div class="flex items-center justify-between gap-3">
          <div class="text-[10px] uppercase tracking-[0.35em] text-amber-700">
            {{ t("install.missing", props.lang) }}
          </div>
          <button
            data-testid="install-action"
            type="button"
            class="rounded-full border border-amber-400/60 bg-amber-200/30 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-amber-700 transition hover:bg-amber-200/50"
            :disabled="installing"
            @click="handleInstall"
          >
            {{ installing ? t("install.installing", props.lang) : t("install.action", props.lang) }}
          </button>
        </div>
        <div v-if="installState === 'failed'" class="mt-2 text-[10px] uppercase tracking-[0.3em] text-rose-500">
          {{ t("install.failed", props.lang) }}
        </div>
      </div>
      <div v-if="props.managerKey === 'pm2'" class="pb-3">
        <label class="flex items-center gap-2 text-[10px] uppercase tracking-[0.3em] text-fg/40">
          {{ t("editor.user", props.lang) }}
          <input
            v-model="pm2User"
            type="text"
            class="h-7 w-full rounded-full border border-border/70 bg-card/70 px-3 text-[11px] normal-case tracking-normal text-fg/70"
          />
        </label>
      </div>
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
