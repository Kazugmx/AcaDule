# AcaDule Academic X Schedule 学生用スケジュール・タスク管理ツール
## クライアント
- CLI tool
- Vue
## サーバー
- Go or Kotlin
  - Kotlinで開発 -> Go Fiber等へシフト
## 機能
- タスク管理
  - タスク作成
  - ステータス管理
  - タグ機能
  - サブタスク
  - 通知
- スケジュール管理
  - カレンダー
  - 時間割登録
  - 複合型でタスク表示
- 学習支援ツール
  - タイマーツール
  - 学習ログ
- 同期機能
- アカウント管理機能
  - ユーザ作成
  - ユーザ削除
  - ログイン・ログアウト
## 機能詳細
### タスク管理
- タスク/データ構造 
```yaml
task:
    id: UUID
    targetName: String(varchar len100)
    description: String 
    progress: String(enum)
    deadline: datetime
    hasDone: boolean
```
