<script setup lang="ts">
import { computed, ref } from "vue";
import { t, type Lang } from "./i18n";
import {
  managerNav,
  themes,
  type ManagerKey,
  type ThemeId
} from "./data/navigation";
import MobileDrawer from "./components/layout/MobileDrawer.vue";
import SideNav from "./components/layout/SideNav.vue";
import TopBar from "./components/layout/TopBar.vue";

const storageKey = {
  lang: "sm_lang",
  theme: "sm_theme"
};

const isLang = (value: string | null): value is Lang =>
  value === "zh" || value === "en";

const isTheme = (value: string | null): value is ThemeId =>
  themes.some((theme) => theme.id === value);

const readStorage = (key: string) => {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(key);
};

const writeStorage = (key: string, value: string) => {
  if (typeof window === "undefined") return;
  localStorage.setItem(key, value);
};

const storedLang = readStorage(storageKey.lang);
const storedTheme = readStorage(storageKey.theme);

const lang = ref<Lang>(isLang(storedLang) ? storedLang : "zh");
const theme = ref<ThemeId>(isTheme(storedTheme) ? storedTheme : "white");
const activeManager = ref<ManagerKey>(managerNav[0].key);
const drawerOpen = ref(false);

const themeClass = computed(() => `theme-${theme.value}`);
const title = computed(() => t("app.title", lang.value));

const toggleLang = () => {
  lang.value = lang.value === "zh" ? "en" : "zh";
  writeStorage(storageKey.lang, lang.value);
};

const setTheme = (next: ThemeId) => {
  theme.value = next;
  writeStorage(storageKey.theme, theme.value);
};

const selectManager = (key: ManagerKey) => {
  activeManager.value = key;
  drawerOpen.value = false;
};
</script>

<template>
  <div
    data-testid="app-shell"
    class="shell-bg min-h-screen bg-bg text-fg"
    :class="themeClass"
  >
    <div class="flex min-h-screen">
      <SideNav
        :items="managerNav"
        :active="activeManager"
        :lang="lang"
        @select="selectManager"
      />

      <div class="flex min-h-screen flex-1 flex-col">
        <TopBar
          :title="title"
          :lang="lang"
          :themes="themes"
          :active-theme="theme"
          @toggle-lang="toggleLang"
          @set-theme="setTheme"
          @toggle-drawer="drawerOpen = !drawerOpen"
        />

        <main class="flex-1 px-4 py-6 lg:px-8 lg:py-8">
          <div class="grid gap-6 lg:grid-cols-[360px_1fr]">
            <section class="rounded-2xl border border-border/60 bg-card/70 p-4 shadow-sm">
              <header class="flex items-center justify-between pb-3">
                <span class="text-[11px] uppercase tracking-[0.35em] text-fg/40">
                  {{ t("layout.list", lang) }}
                </span>
              </header>
              <div class="space-y-3">
                <div
                  v-for="index in 6"
                  :key="index"
                  class="h-8 rounded-lg bg-muted/70"
                ></div>
              </div>
            </section>

            <section class="rounded-2xl border border-border/60 bg-card/70 p-4 shadow-sm">
              <header class="flex items-center justify-between pb-3">
                <span class="text-[11px] uppercase tracking-[0.35em] text-fg/40">
                  {{ t("layout.detail", lang) }}
                </span>
              </header>
              <div class="grid gap-3">
                <div class="h-44 rounded-xl bg-muted/70"></div>
                <div class="h-20 rounded-xl bg-muted/60"></div>
              </div>
            </section>
          </div>
        </main>
      </div>
    </div>

    <MobileDrawer
      :open="drawerOpen"
      :items="managerNav"
      :active="activeManager"
      :lang="lang"
      @close="drawerOpen = false"
      @select="selectManager"
    />
  </div>
</template>
