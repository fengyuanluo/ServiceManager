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
      select: "选择服务",
      logs: "日志",
      edit: "编辑"
    },
    batch: {
      all: "全选",
      selected: "已选"
    },
    action: {
      menu: "打开菜单",
      language: "切换语言",
      create: "新增",
      start: "启动",
      stop: "停止",
      restart: "重启",
      delete: "删除",
      autostart: "自启",
      raw: "原始",
      form: "表单",
      close: "关闭",
      save: "保存",
      reset: "重置"
    },
    create: {
      title: "新增服务",
      placeholder: "填写信息",
      submitting: "创建中",
      failed: "创建失败",
      section: {
        template: "模板",
        selected: "已选",
        form: "表单",
        raw: "原始",
        options: "选项"
      },
      step: {
        template: "模板",
        form: "配置",
        confirm: "确认"
      }
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
    },
    logs: {
      search: "搜索",
      limit: "行数",
      level: "级别",
      levelAll: "全部",
      levelInfo: "信息",
      levelWarn: "警告",
      levelError: "错误",
      export: "导出",
      empty: "暂无日志"
    },
    install: {
      missing: "未检测到",
      action: "安装",
      installing: "安装中",
      failed: "失败",
      confirm: "确认安装该管理器？"
    },
    editor: {
      name: "名称",
      description: "描述",
      exec: "执行",
      workdir: "目录",
      user: "用户",
      env: "环境",
      extension: "扩展",
      unit: "路径",
      pm2User: "PM2用户"
    },
    auth: {
      title: "登录",
      username: "用户名",
      password: "密码",
      submit: "登录",
      error: "登录失败"
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
      select: "Select",
      logs: "Logs",
      edit: "Edit"
    },
    batch: {
      all: "All",
      selected: "Selected"
    },
    action: {
      menu: "Open menu",
      language: "Switch language",
      create: "Create",
      start: "Start",
      stop: "Stop",
      restart: "Restart",
      delete: "Delete",
      autostart: "Autostart",
      raw: "Raw",
      form: "Form",
      close: "Close",
      save: "Save",
      reset: "Reset"
    },
    create: {
      title: "Create Service",
      placeholder: "Fill in",
      submitting: "Creating",
      failed: "Create failed",
      section: {
        template: "Template",
        selected: "Selected",
        form: "Form",
        raw: "Raw",
        options: "Options"
      },
      step: {
        template: "Template",
        form: "Config",
        confirm: "Confirm"
      }
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
    },
    logs: {
      search: "Search",
      limit: "Lines",
      level: "Level",
      levelAll: "All",
      levelInfo: "Info",
      levelWarn: "Warn",
      levelError: "Error",
      export: "Export",
      empty: "Empty"
    },
    install: {
      missing: "Not detected",
      action: "Install",
      installing: "Installing",
      failed: "Failed",
      confirm: "Install this manager?"
    },
    editor: {
      name: "Name",
      description: "Desc",
      exec: "Exec",
      workdir: "Dir",
      user: "User",
      env: "Env",
      extension: "Extra",
      unit: "Unit",
      pm2User: "PM2 User"
    },
    auth: {
      title: "Sign in",
      username: "Username",
      password: "Password",
      submit: "Sign in",
      error: "Sign in failed"
    }
  }
} as const;
