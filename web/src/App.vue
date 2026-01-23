<script setup lang="ts">
import { computed, ref } from "vue";
import { login } from "./api/auth";
import { useAuthState } from "./api/client";
import { t, type Lang } from "./i18n";
import {
  managerNav,
  themes,
  type ManagerKey,
  type ThemeId
} from "./data/navigation";
import LoginPanel from "./components/auth/LoginPanel.vue";
import MobileDrawer from "./components/layout/MobileDrawer.vue";
import SideNav from "./components/layout/SideNav.vue";
import TopBar from "./components/layout/TopBar.vue";
import ManagerPanel from "./components/manager/ManagerPanel.vue";

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
const { isAuthenticated } = useAuthState();
const loginError = ref("");
const loginLoading = ref(false);

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

const handleLogin = async (payload: {
  username: string;
  password: string;
  expire: number;
}) => {
  loginLoading.value = true;
  loginError.value = "";
  try {
    await login(payload.username, payload.password, payload.expire);
  } catch (error) {
    loginError.value =
      error instanceof Error && error.message
        ? error.message
        : t("auth.error", lang.value);
  } finally {
    loginLoading.value = false;
  }
};
</script>

<template>
  <div
    data-testid="app-shell"
    class="shell-bg min-h-screen bg-bg text-fg"
    :class="themeClass"
  >
    <LoginPanel
      v-if="!isAuthenticated"
      :lang="lang"
      :themes="themes"
      :active-theme="theme"
      :loading="loginLoading"
      :error="loginError"
      @toggle-lang="toggleLang"
      @set-theme="setTheme"
      @submit="handleLogin"
    />

    <div v-else class="flex min-h-screen">
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
          <ManagerPanel :manager-key="activeManager" :lang="lang" />
        </main>
      </div>
    </div>

    <MobileDrawer
      v-if="isAuthenticated"
      :open="drawerOpen"
      :items="managerNav"
      :active="activeManager"
      :lang="lang"
      @close="drawerOpen = false"
      @select="selectManager"
    />
  </div>
</template>
