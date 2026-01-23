export const messages = {
  zh: {
    app: {
      title: "服务管理器"
    },
    nav: {
      systemd: "Systemd",
      sysv: "SysV",
      openrc: "OpenRC",
      supervisor: "Supervisor",
      pm2: "PM2"
    },
    layout: {
      list: "列表",
      detail: "详情",
      select: "选择服务"
    },
    action: {
      menu: "打开菜单",
      language: "切换语言"
    },
    theme: {
      white: "纯白",
      dark: "深色",
      pink: "浅粉",
      cyan: "淡青",
      paper: "书卷"
    },
    status: {
      running: "运行中",
      stopped: "已停止",
      failed: "异常",
      unknown: "未知"
    },
    detail: {
      pid: "PID",
      since: "启动时间",
      unit: "路径"
    }
  },
  en: {
    app: {
      title: "Service Manager"
    },
    nav: {
      systemd: "Systemd",
      sysv: "SysV",
      openrc: "OpenRC",
      supervisor: "Supervisor",
      pm2: "PM2"
    },
    layout: {
      list: "List",
      detail: "Detail",
      select: "Select"
    },
    action: {
      menu: "Open menu",
      language: "Switch language"
    },
    theme: {
      white: "White",
      dark: "Dark",
      pink: "Pink",
      cyan: "Cyan",
      paper: "Paper"
    },
    status: {
      running: "Running",
      stopped: "Stopped",
      failed: "Failed",
      unknown: "Unknown"
    },
    detail: {
      pid: "PID",
      since: "Since",
      unit: "Unit"
    }
  }
} as const;
