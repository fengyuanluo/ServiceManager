import { messages } from "./messages";

export type Lang = keyof typeof messages;

export function t(key: string, lang: Lang) {
  const parts = key.split(".");
  let current: unknown = messages[lang];
  for (const part of parts) {
    if (!current || typeof current !== "object" || !(part in current)) {
      return key;
    }
    current = (current as Record<string, unknown>)[part];
  }
  return typeof current === "string" ? current : key;
}
