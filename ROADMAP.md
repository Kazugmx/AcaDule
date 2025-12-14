# 🧭 AcaDule Project Roadmap

> Academic × Schedule — 学生の学習管理とスケジュールを統合するプラットフォーム  
> Backend (Kotlin / Ktor) + CLI (Go) + PWA (Next.js)

---

## 🎯 Phase 1: MVP (Minimum Viable Product)

✅ **目標**: 最小構成で動作する学習タスク管理システム  
CLI・Webいずれからも操作可能なタスクAPIを提供する。

### ✅ 実装範囲
- [x] Backend 基盤構築 (Ktor / Exposed / HikariCP)
- [x] JWT 認証 (AuthService)
- [x] TaskService  
  - [x] `/task` GET/POST  
  - [x] `/task/{id}` GET/POST/DELETE  
- [x] PostgreSQL 接続設定
- [x] BCrypt パスワードハッシュ化
- [ ] CLI (Go) 基本コマンド  
  - [ ] `todo list`
  - [ ] `todo add`
  - [ ] `todo view`
  - [ ] `todo update`
  - [ ] `todo delete`

---

## 🧩 Phase 2: PWA Integration (Web Interface)

🎯 **目標**: Webブラウザから操作できるNext.js PWAを構築。  
ローカル保存と同期をサポート。

### 実装予定
- [ ] Next.js + TypeScript + TailwindCSS
- [ ] PWA 対応（オフラインキャッシュ＋通知）
- [ ] JWT 認証連携
- [ ] タスク一覧ページ
- [ ] 進捗ステータス表示 (`TaskProgress`)
- [ ] API通信 (`/task`, `/login`, `/register`)

---

## 📚 Phase 3: Advanced Modules

🎯 **目標**: 学生の日常スケジュールと学習を統合。

### 実装予定
- [ ] Timetable（授業スケジュール管理）
- [ ] StudySession（集中時間ログ）
- [ ] Subtask 機能
- [ ] Tags（科目分類）

---

## ⚖️ License

- **Project License:** Apache License 2.0  
  - Compatible with Ktor / Exposed stack  
  - Allows commercial and open collaboration

---

## 🛠️ Tech Stack Overview

| Layer    | Technology              | Description          |
|----------|-------------------------|----------------------|
| Backend  | Kotlin / Ktor / Exposed | REST API             |
| Database | PostgreSQL              | Persistent storage   |
| Auth     | JWT + BCrypt            | Secure login         |
| CLI      | Go                      | Terminal task client |
| Web      | Next.js (PWA)           | Web interface        |
| Infra    | Docker / k8s-ready      | Deployable stack     |

---

## 🗓️ Planned Milestones

| Version  | Target      | Description            |
|----------|-------------|:-----------------------|
| `v0.1.0` | MVP完成       | CLI + Task API 基本動作確認  |
| `v0.2.0` | PWAプロトタイプ   | Webクライアントとの統合          |
| `v0.3.0` | Timetable統合 | 学習スケジュール機能実装           |
| `v1.0.0` | 公開版         | 全機能安定化・デプロイ            |

---

## 📦 OSS Components Used

| Library          | Repository                                                              | License            | Purpose                      |
|------------------|-------------------------------------------------------------------------|--------------------|------------------------------|
| **Ktor**         | [ktorio/ktor](https://github.com/ktorio/ktor)                           | Apache License 2.0 | Kotlin Web Framework         |
| **Exposed**      | [JetBrains/Exposed](https://github.com/JetBrains/Exposed)               | Apache License 2.0 | Kotlin ORM / SQL DSL         |
| **BCrypt**       | [patrickfav/bcrypt](https://github.com/patrickfav/bcrypt)               | Apache License 2.0 | Password Hashing             |
| **HikariCP**     | [brettwooldridge/HikariCP](https://github.com/brettwooldridge/HikariCP) | Apache License 2.0 | Database Connection Pool     |
| **Next.js**      | [vercel/next.js](https://github.com/vercel/next.js)                     | MIT                | React Framework              |
| **TailwindCSS**  | [tailwindlabs/tailwindcss](https://github.com/tailwindlabs/tailwindcss) | MIT                | Utility-first CSS            |

---

*Maintained by [kazugmx](https://github.com/kazugmx)*