import type { ManagerKey } from "../data/navigation";
import type { ServiceDetail, ServiceStatus, ServiceSummary } from "../types/service";
import { request } from "./client";

type ApiServiceSummary = {
  name: string;
  status: string;
  detail?: string;
};

type ApiServiceDetail = {
  name: string;
  status: string;
  pid?: string;
  since?: string;
  unit_path?: string;
  description?: string;
  user?: string;
  extra?: Record<string, string>;
};

type ActionRequest = {
  name: string;
  action: string;
  user?: string;
};

type BatchItem = {
  name: string;
  action: string;
};

type BatchResult = {
  name: string;
  action: string;
  success: boolean;
  message?: string;
};

type InstallResult = {
  manager: string;
  success: boolean;
  message?: string;
};

type LogQuery = {
  lines?: number;
  keyword?: string;
  user?: string;
};

const buildQuery = (params: Record<string, string | number | undefined>) => {
  const query = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === "") return;
    query.set(key, String(value));
  });
  const queryString = query.toString();
  return queryString ? `?${queryString}` : "";
};

const normalizeStatus = (status?: string): ServiceStatus => {
  const value = status?.toLowerCase();
  if (value === "running" || value === "stopped" || value === "failed" || value === "unknown") {
    return value;
  }
  return "unknown";
};

const toSummary = (manager: ManagerKey, item: ApiServiceSummary): ServiceSummary => ({
  name: item.name,
  status: normalizeStatus(item.status),
  detail: item.detail,
  manager
});

const toDetail = (manager: ManagerKey, item: ApiServiceDetail): ServiceDetail => ({
  name: item.name,
  status: normalizeStatus(item.status),
  detail: item.description,
  description: item.description,
  since: item.since,
  pid: item.pid,
  unitPath: item.unit_path,
  user: item.user,
  extra: item.extra,
  manager
});

export async function listServices(manager: ManagerKey, user = "") {
  const query = buildQuery({ user });
  const payload = await request<ApiServiceSummary[]>(`/${manager}/list${query}`);
  return payload.map((item) => toSummary(manager, item));
}

export async function getServiceDetail(manager: ManagerKey, name: string, user = "") {
  const query = buildQuery({ name, user });
  const payload = await request<ApiServiceDetail>(`/${manager}/detail${query}`);
  return toDetail(manager, payload);
}

export async function getServiceLogs(manager: ManagerKey, name: string, query: LogQuery = {}) {
  const queryString = buildQuery({
    name,
    lines: query.lines,
    keyword: query.keyword,
    user: query.user
  });
  return request<string[]>(`/${manager}/logs${queryString}`);
}

export async function runServiceAction(manager: ManagerKey, payload: ActionRequest) {
  return request<string>(`/${manager}/action`, {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function runServiceBatch(manager: ManagerKey, items: BatchItem[], user = "") {
  return request<BatchResult[]>(`/${manager}/batch`, {
    method: "POST",
    body: JSON.stringify({ items, user })
  });
}

export async function installManager(manager: ManagerKey) {
  return request<InstallResult>(`/${manager}/install`, {
    method: "POST"
  });
}
