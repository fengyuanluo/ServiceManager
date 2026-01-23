import type { ManagerKey } from "../data/navigation";

export type ServiceStatus = "running" | "stopped" | "failed" | "unknown";

export interface ServiceSummary {
  id: string;
  name: string;
  status: ServiceStatus;
  summary: string;
  manager: ManagerKey;
}

export interface ServiceDetail extends ServiceSummary {
  description?: string;
  since?: string;
  pid?: number;
  unitPath?: string;
}
