<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { t, type Lang } from "../../i18n";
import type { ManagerKey } from "../../data/navigation";
import type { ServiceDetail } from "../../types/service";

const props = defineProps<{
  service: ServiceDetail | null;
  managerKey: ManagerKey;
  lang: Lang;
}>();

const mode = ref<"form" | "raw">("form");
const rawText = ref("");

const form = ref({
  name: "",
  description: "",
  exec: "",
  workdir: "",
  user: "root",
  env: ""
});

const isPm2 = computed(() => props.managerKey === "pm2");

const buildRaw = (service: ServiceDetail | null) => {
  if (!service) return "";
  if (props.managerKey === "systemd") {
    return `[Unit]\nDescription=${service.description ?? service.detail ?? ""}\n\n[Service]\nExecStart=${form.value.exec || "/usr/sbin/" + service.name}\nUser=${form.value.user}\n\n[Install]\nWantedBy=multi-user.target\n`;
  }
  if (props.managerKey === "pm2") {
    return `{\n  \"name\": \"${service.name}\",\n  \"script\": \"${form.value.exec || "app.js"}\",\n  \"user\": \"${form.value.user}\"\n}\n`;
  }
  return `# ${service.name}\nexec=${form.value.exec || service.name}\nuser=${form.value.user}\n`;
};

const resetForm = (service: ServiceDetail | null) => {
  if (!service) {
    form.value = {
      name: "",
      description: "",
      exec: "",
      workdir: "",
      user: "root",
      env: ""
    };
    rawText.value = "";
    return;
  }
  form.value = {
    name: service.name,
    description: service.description ?? service.detail ?? "",
    exec: "",
    workdir: "/",
    user: service.user || "root",
    env: ""
  };
  rawText.value = buildRaw(service);
};

watch(
  () => props.service,
  (service) => {
    resetForm(service);
  },
  { immediate: true }
);

const toggleMode = () => {
  mode.value = mode.value === "form" ? "raw" : "form";
  if (mode.value === "raw" && props.service) {
    rawText.value = buildRaw(props.service);
  }
};
</script>

<template>
  <div class="space-y-4">
    <header class="flex items-center justify-between">
      <div class="text-[11px] uppercase tracking-[0.35em] text-fg/40">
        {{ t("layout.edit", props.lang) }}
      </div>
      <button
        data-testid="toggle-raw"
        type="button"
        class="rounded-full border border-border/70 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-fg/70 transition hover:text-fg"
        @click="toggleMode"
      >
        {{ mode === "form" ? t("action.raw", props.lang) : t("action.form", props.lang) }}
      </button>
    </header>

    <div v-if="!props.service" class="rounded-2xl border border-dashed border-border/60 bg-card/50 px-6 py-16 text-center text-[10px] uppercase tracking-[0.35em] text-fg/40">
      {{ t("layout.select", props.lang) }}
    </div>

    <div v-else>
      <div v-if="mode === 'form'" data-testid="editor-form" class="space-y-4">
        <div class="grid gap-3 md:grid-cols-2">
          <label class="space-y-2 text-xs text-fg/60">
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("editor.name", props.lang) }}
            </span>
            <input
              v-model="form.name"
              type="text"
              class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
              disabled
            />
          </label>
          <label class="space-y-2 text-xs text-fg/60">
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("editor.user", props.lang) }}
            </span>
            <select
              v-model="form.user"
              class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
            >
              <option value="root">root</option>
              <option value="app">app</option>
              <option value="node">node</option>
            </select>
          </label>
        </div>
        <label class="space-y-2 text-xs text-fg/60">
          <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
            {{ t("editor.description", props.lang) }}
          </span>
          <input
            v-model="form.description"
            type="text"
            class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
            data-testid="editor-description"
          />
        </label>
        <div class="grid gap-3 md:grid-cols-2">
          <label class="space-y-2 text-xs text-fg/60">
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("editor.exec", props.lang) }}
            </span>
            <input
              v-model="form.exec"
              type="text"
              class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
            />
          </label>
          <label class="space-y-2 text-xs text-fg/60">
            <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
              {{ t("editor.workdir", props.lang) }}
            </span>
            <input
              v-model="form.workdir"
              type="text"
              class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
            />
          </label>
        </div>
        <label class="space-y-2 text-xs text-fg/60">
          <span class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
            {{ t("editor.env", props.lang) }}
          </span>
          <textarea
            v-model="form.env"
            rows="3"
            class="w-full rounded-xl border border-border/60 bg-card/70 px-3 py-2 text-sm"
          ></textarea>
        </label>

        <div class="rounded-2xl border border-border/60 bg-card/60 p-4 text-xs text-fg/60">
          <div class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
            {{ t("editor.extension", props.lang) }}
          </div>
          <div class="mt-3 grid gap-3 md:grid-cols-2">
            <div class="space-y-1">
              <div class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
                {{ t("editor.unit", props.lang) }}
              </div>
              <div class="break-all text-xs">
                {{ props.service.unitPath ?? "-" }}
              </div>
            </div>
            <div v-if="isPm2" class="space-y-1">
              <div class="text-[10px] uppercase tracking-[0.3em] text-fg/40">
                {{ t("editor.pm2User", props.lang) }}
              </div>
              <div class="text-xs">{{ form.user }}</div>
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-2">
          <button
            data-testid="editor-reset"
            type="button"
            class="rounded-full border border-border/70 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-fg/70 transition hover:text-fg"
            @click="resetForm(props.service)"
          >
            {{ t("action.reset", props.lang) }}
          </button>
          <button
            data-testid="editor-save"
            type="button"
            class="rounded-full border border-primary/60 bg-primary/10 px-3 py-1 text-[11px] uppercase tracking-[0.3em] text-primary transition hover:bg-primary/20 disabled:cursor-not-allowed disabled:opacity-60"
            disabled
          >
            {{ t("action.save", props.lang) }}
          </button>
        </div>
        <div class="text-right text-[10px] uppercase tracking-[0.3em] text-fg/40">
          保存功能暂未开放
        </div>
      </div>

      <div v-else data-testid="editor-raw" class="space-y-3">
        <textarea
          v-model="rawText"
          rows="10"
          class="w-full rounded-2xl border border-border/60 bg-card/70 p-3 font-mono text-xs"
        ></textarea>
      </div>
    </div>
  </div>
</template>
