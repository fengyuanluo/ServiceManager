import type { ManagerKey } from "../data/navigation";

export type ServiceStatus = "running" | "stopped" | "failed" | "unknown";

export interface ServiceSummary {
  name: string;
  status: ServiceStatus;
  detail?: string;
  manager: ManagerKey;
}

export interface ServiceDetail extends ServiceSummary {
  description?: string;
  since?: string;
  pid?: string;
  unitPath?: string;
  user?: string;
  extra?: Record<string, string>;
}

export interface CreateOptions {
  start?: boolean;
  autostart?: boolean;
}

export interface CreateRequest {
  name: string;
  description?: string;
  exec?: string;
  workdir?: string;
  user?: string;
  env?: string;
  template_id?: string;
  raw_text?: string;
  extra?: Record<string, string>;
  options?: CreateOptions;
}
