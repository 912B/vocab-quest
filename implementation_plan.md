# Vocab Quest - Technical Design Document (TDD)

## 1. 項目概述 (Project Overview)
**Vocab Quest** 是一款極簡主義的 "深空賽博朋克" 風格單詞記憶遊戲。核心目標是通過高強度的視覺反饋和沈浸式交互，讓枯燥的背單詞變得像駕駛星際飛船一樣刺激。

**文檔目標 (Documentation Goal)**:
本技術設計文檔 (TDD) 旨在完全獨立。**任何開發者在完全不了解項目前提下，應能僅憑本文檔直接開發出該遊戲。**

**設計哲學**:
*   **視覺**: 深色模式 (Dark Mode)，霓虹光效 (Neon Glow)，全息投影風格。
*   **交互**: 鍵盤優先 (Keyboard First)，無鼠標操作，極速響應。
*   **核心**: 難度自適應 (Adaptive)，永遠處於 "舒適區邊緣" (Zone of Proximal Development)。
*   **移動端適配 (Mobile)**:
    -   允許 "PERFECT" 印章遮擋鍵盤 (節省空間)。
    -   **[CRITICAL]** 輸入框 (Input Slots) 和 題目 (Definition) 必須始終**浮在印章之上** (Z-Index > Stamp)，確保在印章停留期間仍可清晰看到輸入內容。

---

## 2. 遊戲核心機制 (Core Game Mechanics)

### 2.1 會話循環 (Session Loop)
遊戲不分關卡，而是以 "會話 (Session)" 為單位。
1.  **Load**: 系統從後端加載 **10個單詞** 組成一個 Session。
    -   **選詞策略**: `30%` 新詞 (無記錄) + `70%` 複習詞 (根據遺忘曲線排期)。
2.  **Play**: 玩家逐個攻克單詞 (見 2.2)。
3.  **Summary**: 10詞全部完成後，顯示統計 (準確率、得分)，並提供 "再次出擊 (Play Again)" 按鈕刷新 Session。

### 2.2 單詞交互流程 (Word Interaction Flow)
每個單詞的生命週期包含以下狀態：

1.  **準備 (Ready)**:
    -   顯示中文定義 (`definition`).
    -   顯示英文單詞挖空模板 (Masked Word).
    -   自動播放發音 (TTS).
    -   *視覺*: 隕石 (Asteroid) 開始緩緩下落 (視覺壓力)。

2.  **輸入 (Input)**:
    -   玩家通過物理鍵盤或屏幕虛擬鍵盤輸入字母。
    -   **輸入槽 (Slots)**: 每個字母一個方框，輸入後自動跳轉下一格。
    -   **退格 (Backspace)**: 支持刪除回退。

3.  **判定 (Evaluation)**: (當輸入長度等於單詞長度時觸發)
    -   **CASE A: 正確 (Success)**
        -   **行為**: 播放金色 "PERFECT" 印章動畫，鎖定屏幕。
        -   **數據**: `Proficiency` +1 (Max 5)。
        -   **得分**: `Base(100) + Combo * 10`。
        -   **連擊**: `Combo` +1。
        -   **過場**: 停留 2.5秒 (Tension Delay)，然後自動進入下一個詞。

    -   **CASE B: 錯誤 (Failure)**
        -   **行為**: 屏幕劇烈震動 (CSS Shake)，背景變紅。
        -   **數據**: `Proficiency` 降級 (<=3 歸零, >3 降2級)。
        -   **懲罰**: `Combo` 歸零，印章消失 (如果有的話)。
        -   **修正 (Correction)**: 
            1.  顯示正確拼寫 (綠色高亮) 2秒。
            2.  進入 **強制重試 (Forced Retry)** 狀態。
            3.  出現 "RETRY CONNECTION" 按鈕。
            4.  玩家必須手動輸入一遍正確單詞 (此過程**不加分**，**不計入進度**)。
            5.  輸入正確後，才算通過該詞 (Pass)，進入下一個。

### 2.3 難度自適應系統 (Adaptive Difficulty)
挖空 (Masking) 算法根據單詞的 `proficiency` (熟練度, 0-5) 動態生成。

| 熟練度 (Level) | 狀態 | 遮擋比例 (Mask Ratio) | 描述 |
| :--- | :--- | :--- | :--- |
| **0 - 1** | 新手 (Novice) | **40%** | 保留大量提示 (首尾字母常用)。 |
| **2 - 3** | 熟練 (Skilled) | **60%** | 增加難度，考驗拼寫。 |
| **4 - 5** | 大師 (Master) | **80%** | 接近盲打，僅保留極少提示。 |

*   **邊界規則**: `1 <= 遮擋數量 < 單詞長度` (至少遮1個，至少留1個)。

---

## 3. 系統架構 (Technical Architecture)

## Coding Standards & Layout Principles

### 1. Engineering Philosophy (Core Rule)
*   **Root Cause Resolution**: Never use "force" (e.g., `position: fixed` hacks, `!important`, `overflow: hidden` to hide bugs) to patch a symptom.
*   **Design-Led Solutions**: If a layout breaks (e.g., content too wide), **redesign the layout logic** (e.g., responsive collapse) rather than creating rigid patches.
*   **Anti-Patch**: Do not "lock" the interface to hide overflow issues. Fix the overflow.

### 2. Hub & Spoke Navigation (SPA)
*   **Central Hub (Dashboard)**: The entry point for all non-game activities (Stats, Profile, Admin).
*   **Game Loop (Isolated)**: The game runs in a dedicated, distraction-free environment.
*   **Unified Store**: A single `Store.js` manages global state (User Session, Current Route).

### 3. Layout & Styling Standards (CRITICAL)
*   **Strict Separation**: JavaScript files **MUST NOT** contain layout definition code (e.g., `element.style.marginTop = '20px'`).
*   **CSS Authority**: All positioning, sizing, margins, and padding must be defined in CSS files (`layout.css`, `theme.css`).
*   **Class-Based State**: JS should only manipulate layout by toggling CSS classes (e.g., `element.classList.add('hidden')`), never by setting inline styles.
*   **Exception**: Complex dynamic scaling (e.g., fitting a game canvas to a specific aspect ratio based on window size) is the only permitted use of inline `transform` or `width/height` in JS.

### 3.1 技術棧 (Tech Stack)
*   **Frontend**: Vanilla JS (ES6 Modules), HTML5, CSS3 (CSS Variables for Theming).
*   **Backend**: Go (Golang) - Standard Library or Lightweight Router.
*   **Database**: SQLite (嵌入式關係數據庫).
*   **Protocol**: REST API (JSON).

### 3.2 數據模型 (Data Models)

#### Word (單詞)
```json
{
  "id": 1,
  "text": "exercise",
  "definition": "鍛煉",
  "difficulty": 1
}
```

#### UserProgress (學習記錄)
**核心原則**: 後端只做數據採集 (Data Collection)，不做複雜評級。
```sql
CREATE TABLE user_progress (
    user_id INT,
    word_id INT,
    attempts INT DEFAULT 0,     -- 總答題次數
    successes INT DEFAULT 0,    -- 正確次數
    last_played_at DATETIME     -- 最後一次練習時間
);
```

### 3.3 數據驅動邏輯 (Data-Driven Logic)

*   **難度計算 (Frontend Derived)**:
    -   熟練度 = `successes / attempts` (正確率) + `attempts` (經驗值).
    -   例如: 練習少於 3 次 -> 新手; 正確率 > 80% -> 大師.

*   **會話生成策略 (Session Strategy - The "3-4-3" Mix)**:
    為了保持 "心流 (Flow)" 體驗，建議每次會話 (10詞) 採用以下黃金配比：
    1.  **30% 新詞 (New)**: `attempts = 0`. 
        *   *目的*: 擴充詞彙量，提供新鮮感。
    2.  **40% 弱點 (Weakness)**: `success_rate < 0.6` (或者歷史錯誤次數高).
        *   *目的*: "刻意練習" (Deliberate Practice)，攻克難點。
    3.  **30% 複習 (Review)**: `last_played > 3 days` (按時間倒序).
        *   *目的*: 增強信心，防止遺忘，同時作為 "熱身" 或 "冷卻" 使用。
    
### 3.4 後端邏輯實現 (Backend Implementation: Review Fixed, Remedial First)
**文件位置**: `server/services/session_strategy.go` (獨立邏輯文件)

根據用戶指示：**"複習詞固定3個，弱點詞多就犧牲新詞，保持 3-4-3 骨架"**。

**優先級隊列 (Priority Queue)**:
1.  **複習詞 (Review)**: **固定 3 個** (只要有足夠的複習庫存)。保證每次都有熱身。
    *   *Limit = 3*
2.  **弱點詞 (Weak)**: **填補剩餘空間**。
    *   *Limit = 10 - count(Review)*
    *   如果錯題很多 (e.g. 7個)，就會佔滿剩下的位置，導致**沒有新詞** (0 New)。
    *   如果錯題適中 (e.g. 4個)，就會剩下 3 個位置給新詞 (即經典 3-4-3)。
3.  **新詞 (New)**: **最後填充**。
    *   *Limit = 10 - count(Review) - count(Weak)*

**偽代碼**:
```go
session = []

// 1. Review (Fixed 3)
review_words = query(Review, limit=3)
session.append(review_words)

// 2. Weak (Fill Remainder)
remaining = 10 - len(session)
weak_words = query(Weak, limit=remaining)
session.append(weak_words)

// 3. New (Fill Remainder)
remaining = 10 - len(session)
if remaining > 0:
    new_words = query(New, limit=remaining)
    session.append(new_words)
```

---

## 4. UI/UX 規範 (Interface Specs)

### 4.1 全局導航 (Global Topbar)
所有頁面頂部常駐黑色半透明 Bar。
*   **左側**: LOGO (Vocab Quest)。
*   **中間 (HUD)**: 
    -   Mission Progress: `3/10` (當前/總數).
    -   Score: `1250 XP`.
    -   Combo: `HYPER-FLUX x4` (僅 >1 時顯示).
*   **右側**: 
    -   Stats (圖表圖標).
    -   Settings/Admin (齒輪, 僅 Admin 可見).
    -   Logout (退出).

### 4.2 視覺風格 (Visual Style)
*   **調色板**:
    -   **Neon Cyan**: `#00f3ff` (UI 邊框, 正常狀態)
    -   **Neon Green**: `#00ff9d` (正確, 成功)
    -   **Neon Red**: `#ff3333` (錯誤, 危險)
    -   **Neon Gold**: `#ffd700` (印章, 高分)
    -   **Background**: `#050510` (深空黑)
*   **字體**: 無襯線字體 (System UI), 大寫為主, 寬字間距 (Letter Spacing).

---

## 5. 部署與運行 (Deployment)
1.  **Backend**: `go run server/main.go` (監聽 8081).
2.  **Frontend**: 靜態文件由 Go 服務器在 `/` 路由下服務。
3.  **Database**: 啟動時自動檢查 `vocab.db`，若不存在則自動種子數據 (Seed)。

---
## 6. 重構與架構升級 (Refactoring & Architecture)
**User Approved**: 2025-12-29
**Goal**: 解決 "有機生長" 帶來的代碼混亂，建立可維護的 SPA (Single Page Application) 架構。

### 6.1 核心原則 (Core Principles)
1.  **單一入口 (Single Entry Point)**: 所有頁面邏輯由 `ClientApp` 統一調度，不再依賴散落在 HTML 中的 `<script>` 標籤。
2.  **組件化 (Component-Based)**: 所有 UI 必須繼承自 `BaseComponent`，擁有標準生命週期 (`mount`, `unmount`, `render`)。
3.  **狀態集中 (Centralized State)**: 全局狀態 (User, Audio, Config) 由 `Store` 管理，禁止全局變量污染。
4.  **資源統一 (Unified Assets)**: 圖片、音頻路徑由 `AssetManager` 統一解析，徹底解決 `../../` 路徑地獄。

### 6.2 目錄結構規範 (Directory Structure)
```text
client/
├── src/
│   ├── core/           # 核心引擎
│   │   ├── App.js      # 主控制器 (Router, Init)
│   │   ├── Component.js # 組件基類
│   │   ├── Store.js    # 全局狀態
│   │   └── Audio.js    # 音頻引擎
│   ├── items/          # 通用組件 (UI Kit)
│   │   ├── Button.js
│   │   ├── Topbar.js
│   │   └── Modal.js
│   ├── modules/        # 業務模塊 (Pages)
│   │   ├── auth/       # 登錄相關
│   │   ├── game/       # 遊戲核心
│   │   └── dashboard/  # 數據看板
│   └── utils/          # 工具函數
### Dictionary Selector
- **Frontend**: `client/src/items/Topbar.js`
    - Add Dropdown `<select id="dict-select">`.
    - Fetch API `GET /dictionaries`.
    - Load saved selection from `localStorage.getItem('vq_dict_id')`.
    - On Change: Save to `localStorage` and `window.location.reload()` (Simplest way to reset game engine).
- **Frontend**: `client/src/modules/game/engine.js`
    - In `Game.init`, read `localStorage.getItem('vq_dict_id')`.
    - Call `/api/session?dictionary_id=...`.
- **Backend**:
    - `GameHandler.GetSession`: Parse `dictionary_id` param.
    - `LearningEngine.GenerateSession`: Accept `dictionaryID`.
    - `ProgressRepository`: Update `GetDue`, `GetNew`, `GetReviewAhead` to filter by `AND w.dictionary_id = ?` if ID > 0.

### Stats Sync (Dictionary Filter)
- **Frontend**: `client/src/modules/dashboard.js`
    - Read `localStorage.getItem('vq_dict_id')`.
    - `API.get('/stats?dictionary_id=...')`.
- **Backend**:
    - `GameHandler.HandleGetStats`: Parse `dictionary_id`.
    - `services.GetUserStats`: Accept `dictionaryID`.
    - **SQL Updates**:
        - Join `user_progress` with `words` to filter by dictionary.
        - Filter `Count(*)` from `words` by dictionary.

### Auto-Login (Persistence)
- **Frontend**: `client/index.html`
    - In `init()`: Check `API.userId`.
    - If valid, skip login form and immediately call `Engine.init()`.
    - Fallback: If `Engine.init` fails (e.g. 401 session invalid), alert and show login.
├── assets/             # 靜態資源 (Audio, Images)
├── index.html          # 唯一入口 HTML
└── styles/             # 全局樣式
```

### 6.3 遷移路線圖 (Migration Roadmap)
1.  **Foundation**: 建立 `src/core` 基礎設施 (App, Store, Component)。
2.  **Auth Module**: 重寫登錄頁面為 `AuthModule`。
3.  **Game Module**: 移植遊戲引擎到 `GameModule`，剝離 UI 代碼。
4.  **Cleanup**: 刪除 `client/_legacy` 及散落的 JS 文件。

---

## 7. 組件開發規範 (Component Standards)

### 7.1 BaseComponent 簽名
```javascript
class BaseComponent {
    constructor(container) {
        this.container = container;
        this.element = null;
    }

    // 必須實現
    render() { return '<div>...</div>'; }

    // 生命週期
    mount() {
        this.element = document.createElement('div');
        this.element.innerHTML = this.render();
        this.container.appendChild(this.element);
        this.onMount(); // 綁定事件
    }

    unmount() {
        this.onUnmount(); // 解綁事件
        if (this.element) this.element.remove();
    }
}
```

### 7.2 樣式管理
*   **Scoped CSS**: 盡量使用組件內的 Class 前綴 (e.g. `.comp-button-primary`).
*   **Global Variables**: 顏色、字體必須使用 `theme.css` 變量。
*   **Z-Index**: 統一在 `theme.css` 中定義 Z-Index 層級變量，禁止手寫 `9999`。

---

---

## 8. 後台功能擴展 (Admin Features Expansion)

### 8.1 Excel 文件導入 (Excel Import Feature)
允許管理員上傳 `.xlsx` 文件批量導入單詞，相較於純文本更易於管理和編輯。

#### 8.1.1 後端接口 (Backend Endpoint)
*   **路由**: `POST /api/dictionaries/import` (Multipart/Form-Data)
*   **參數**:
    -   `file`: `.xlsx` 文件。
    -   `dictionary_id`: 目標詞典 ID (Form Value)。
*   **Excel 格式規範**:
    -   **Sheet 1**: 默認讀取第一個工作表。
    -   **第一行**: 表頭 (Header)，程序將忽略第一行。
    -   **列結構**:
        -   **A列 (Column A)**: 單詞 (Word Text)
        -   **B列 (Column B)**: 定義 (Definition)
        -   **C列 (Column C)**: 難度 (Difficulty) [可選, 默認為1]
*   **邏輯**:
    -   使用 `excelize` 庫解析文件。
    -   遍歷行 (從第2行開始)。
    -   跳過空行。
    -   插入數據庫 (忽略重複)。

#### 8.1.2 前端界面 (Library Module)
*   **位置**: `client/admin/modules/library.js`
*   **UI 元素**:
    -   按鈕: `[UPLOAD EXCEL]` (c-button-neon)。
    -   文件選擇器: `<input type="file" accept=".xlsx">` (隱藏或美化)。
    -   提示: "請上傳 Excel 文件，A列單詞，B列定義"。
*   **集成**:
    -   使用 `FormData` 異步上傳。
    -   上傳成功後顯示 "成功導入 N 個單詞" 並刷新列表。

---

### 8.2 用戶刪除功能 (User Deletion)
允許管理員刪除用戶賬號。

#### 8.2.1 後端接口 (Backend Endpoint)
*   **路由**: `DELETE /api/users` (Query Param: `id`)
*   **權限**: Admin Only
*   **邏輯**:
    -   檢查 ID 是否存在。
    -   從數據庫物理刪除 (DELETE FROM users WHERE id = ?).

#### 8.2.2 前端界面 (Users Module)
*   **位置**: `client/admin/modules/users.js`
*   **UI 元素**:
    -   在用戶列表每一行增加 [DELETE] 按鈕 (紅色)。
    -   點擊後彈出確認框 (Confirm)。
    -   確認後調用 API 並刷新列表。

---
*End of Specification*
