# Vocab Quest - 技术设计文档 (TDD)

## 1. 项目概述 (Project Overview)
**Vocab Quest** 是一款极简主义的 "深空赛博朋克" 风格单词记忆游戏。
核心目标是通过高强度的视觉反馈和沉浸式交互，让枯燥的背单词变得像驾驶星际飞船一样刺激。

**设计哲学**:
*   **视觉**: 深色模式 (Dark Mode)，霓虹光效 (Neon Glow)，全息投影风格。
*   **交互**: 键盘优先 (Keyboard First)，无鼠标操作，极速响应。
*   **核心**: 难度自适应 (Adaptive)，永远处于 "舒适区边缘" (Zone of Proximal Development)。
*   **移动端适配**: 针对移动设备优化，输入框始终位于视觉中心。

---

## 2. 游戏核心机制 (Core Mechanics)

### 2.1 会话循环 (Session Loop)
游戏以 **Session (会话)** 为单位，不分关卡。
1.  **生成 (Load)**: 每次加载 **10个单词**。
    -   **30% 新词 (New)**: 从未接触过的单词。
    -   **40% 弱点 (Weak)**: 历史错误率高的单词。
    -   **30% 复习 (Review)**: 根据遗忘曲线需要复习的单词。
2.  **游玩 (Play)**: 玩家逐个攻克单词。
3.  **结算 (Summary)**: 显示准确率、得分，更新数据库进度。

### 2.2 交互流程 (Interaction Flow)
*   **准备**: 显示中文定义，单词挖空 (如 `e_er__se`)，播放发音。
*   **输入**: 玩家输入字母填充空缺。
*   **判定**:
    -   **正确 (Perfect)**: 金色印章动画，得分 (+Combo)，进入下一词。
    -   **错误 (Fail)**: 屏幕震动，背景变红，Combo 清零。必须**强制重打一遍**正确拼写才能继续 (不计分)。
*   **双重过关机制 (Double-Pass Drill)**:
    -   所有未掌握的单词必须在一个 Session 内**连续两次**正确回答才能算作 "临时掌握" 并移除队列，否则会不断在该 Session 中循环出现。

---

## 3. 系统架构 (System Architecture)

### 3.1 技术栈 (Tech Stack)
*   **Frontend**: 原生 JavaScript (ES6 Modules), HTML5, CSS3 Vars.
*   **Backend**: Go (Golang) 标准库 + SQLite。
*   **安全**: `bcrypt` 密码哈希。
*   **部署**: 单二进制文件 (Server) + 静态资源。

### 3.2 目录结构
```text
client/
├── src/
│   ├── core/           # 核心 (App, Store, Audio)
│   ├── items/          # 组件 (Topbar, Button)
│   ├── modules/        # 业务 (Game, Dashboard, Auth)
│   └── utils/          # 工具
admin/                  # 后台管理前端
server/
├── handlers/           # API 路由处理
├── services/           # 业务逻辑 (学习引擎, 统计)
├── repositories/       # 数据库操作
└── main.go             # 入口
```

### 3.3 数据库模型
*   **Users**: `id`, `username`, `password` (Hash), `role` (admin/user).
*   **Dictionaries**: `id`, `name`, `description`, `is_active`.
*   **Words**: `id`, `dictionary_id`, `text`, `definition`.
*   **UserProgress**: `user_id`, `word_id`, `attempts`, `successes`, `next_review_at`.

---

## 4. 功能模块详解 (Features)

### 4.1 词库管理 (Dictionary System)
*   **多词库支持**: 允许创建、激活多个词典。
*   **Excel 导入**: 支持从 `.xlsx` 文件批量导入单词 (Admin)。
    -   格式: A列=单词, B列=定义, C列=难度。
*   **词库切换**: 
    -   顶栏 (Topbar) 下拉菜单实时切换词库。
    -   选择自动保存至 `localStorage`。
    -   切换后游戏引擎和统计页面自动过滤数据，仅显示当前词库内容。

### 4.2 安全与配置 (Security & Config)
*   **密码安全**: 所有用户密码均使用 `bcrypt` 算法加密存储。
*   **登录状态持有**: 
    -   使用 `localStorage` 存储 Session 标识。
    -   **自动登录**: 刷新页面或重新打开时，自动跳过登录页进入游戏。
*   **服务器配置 (`server_config.json`)**:
    -   **Port**: 可自定义服务端口 (默认 8081)。
    -   **Admin Reset**: 通过配置文件重置管理员密码 (重启后自动失效，保证安全)。

### 4.3 后台管理 (Admin Dashboard)
*   **仪表盘**: 查看系统总用户数、总单词数。
*   **用户管理**: 查看所有用户，**删除用户**。
*   **词库管理**: 上传 Excel，管理词库列表。

---

## 5. 部署与版本控制 (Deployment & Git)
*   **版本控制**: 已初始化 Git 仓库，配置了 `.gitignore` 排除敏感文件 (配置、数据库)。
*   **运行**:
    ```bash
    go run server/main.go
    ```
*   **访问**: 浏览器打开 `http://localhost:8081`。

---

## 6. 后续规划 (Roadmap)
*   **[P2] 移动端 PWA**: 支持离线缓存，添加到主屏幕。
*   **[P3] 多人对战**: 实时单词拼写竞速。
