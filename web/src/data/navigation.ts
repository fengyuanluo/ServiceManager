export const managerNav = [
  { key: "systemd", labelKey: "nav.systemd" },
  { key: "sysv", labelKey: "nav.sysv" },
  { key: "openrc", labelKey: "nav.openrc" },
  { key: "supervisor", labelKey: "nav.supervisor" },
  { key: "pm2", labelKey: "nav.pm2" }
] as const;

export type ManagerKey = (typeof managerNav)[number]["key"];

export const themes = [
  { id: "white", labelKey: "theme.white", swatch: "hsl(0 0% 100%)" },
  { id: "dark", labelKey: "theme.dark", swatch: "hsl(222 20% 12%)" },
  { id: "pink", labelKey: "theme.pink", swatch: "hsl(340 70% 55%)" },
  { id: "cyan", labelKey: "theme.cyan", swatch: "hsl(190 70% 45%)" },
  { id: "paper", labelKey: "theme.paper", swatch: "hsl(35 35% 70%)" }
] as const;

export type ThemeId = (typeof themes)[number]["id"];
