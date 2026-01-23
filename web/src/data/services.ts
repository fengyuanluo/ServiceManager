import type { ManagerKey } from "./navigation";
import type { ServiceDetail, ServiceSummary } from "../types/service";

const systemdServices: ServiceSummary[] = [
  {
    name: "sshd",
    status: "running",
    detail: "Secure Shell",
    manager: "systemd"
  },
  {
    name: "nginx",
    status: "running",
    detail: "Web Server",
    manager: "systemd"
  }
];

const sysvServices: ServiceSummary[] = [
  {
    name: "cron",
    status: "running",
    detail: "Scheduler",
    manager: "sysv"
  }
];

const openrcServices: ServiceSummary[] = [
  {
    name: "net",
    status: "running",
    detail: "Network",
    manager: "openrc"
  }
];

const supervisorServices: ServiceSummary[] = [
  {
    name: "app",
    status: "stopped",
    detail: "Worker",
    manager: "supervisor"
  }
];

const pm2Services: ServiceSummary[] = [
  {
    name: "node-api",
    status: "running",
    detail: "Node API",
    manager: "pm2"
  }
];

const serviceSummaries: Record<ManagerKey, ServiceSummary[]> = {
  systemd: systemdServices,
  sysv: sysvServices,
  openrc: openrcServices,
  supervisor: supervisorServices,
  pm2: pm2Services
};

const serviceDetails: Record<ManagerKey, Record<string, ServiceDetail>> = {
  systemd: {
    sshd: {
      name: "sshd",
      status: "running",
      detail: "Secure Shell",
      manager: "systemd",
      description: "OpenSSH server daemon",
      since: "2025-12-18 09:12",
      pid: "1223",
      unitPath: "/lib/systemd/system/ssh.service"
    },
    nginx: {
      name: "nginx",
      status: "running",
      detail: "Web Server",
      manager: "systemd",
      description: "Reverse proxy and static web server",
      since: "2025-12-20 14:30",
      pid: "2045",
      unitPath: "/lib/systemd/system/nginx.service"
    }
  },
  sysv: {
    cron: {
      name: "cron",
      status: "running",
      detail: "Scheduler",
      manager: "sysv",
      description: "System task scheduler"
    }
  },
  openrc: {
    net: {
      name: "net",
      status: "running",
      detail: "Network",
      manager: "openrc",
      description: "Network stack"
    }
  },
  supervisor: {
    app: {
      name: "app",
      status: "stopped",
      detail: "Worker",
      manager: "supervisor",
      description: "Supervisor managed worker"
    }
  },
  pm2: {
    "node-api": {
      name: "node-api",
      status: "running",
      detail: "Node API",
      manager: "pm2",
      description: "PM2 managed Node service"
    }
  }
};

const serviceLogs: Record<ManagerKey, Record<string, string[]>> = {
  systemd: {
    sshd: [
      "[info] sshd started",
      "[info] key exchange completed",
      "[warn] invalid user",
      "[info] session closed"
    ],
    nginx: [
      "[info] worker started",
      "[info] new connection",
      "[error] upstream timeout"
    ]
  },
  sysv: {
    cron: ["[info] cron wakeup", "[info] job executed"]
  },
  openrc: {
    net: ["[info] interface up", "[info] dhcp lease ok"]
  },
  supervisor: {
    app: ["[warn] app stopped", "[info] waiting for restart"]
  },
  pm2: {
    "node-api": ["[info] pm2 online", "[error] restart triggered"]
  }
};

export const getServiceSummaries = (manager: ManagerKey) =>
  serviceSummaries[manager] ?? [];

export const getServiceDetail = (manager: ManagerKey, id: string) =>
  serviceDetails[manager]?.[id] ?? null;

export const getServiceLogs = (manager: ManagerKey, id: string) =>
  serviceLogs[manager]?.[id] ?? [];
