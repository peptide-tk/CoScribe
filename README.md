# CoScribe - リアルタイム共同執筆ツール

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.22-blue)
![React Version](https://img.shields.io/badge/React-18-blue)
![TypeScript](https://img.shields.io/badge/TypeScript-5.0-blue)

Go WebSocket + React TypeScript で実装したリアルタイム共同文書編集アプリケーション。

## 機能

- **リアルタイム共同編集**: 複数ユーザーの同時編集とリアルタイム同期
- **WebSocket 通信**: gorilla/websocket による低遅延双方向通信
- **楽観ロック**: バージョン管理による競合検出・解決
- **自動保存**: 文書の自動永続化
- **マルチルーム**: 文書別の独立した編集セッション

---

## 技術スタック

### バックエンド

- **Go 1.22** - gorilla/websocket, Gin
- **並行処理** - Goroutine/Channel による多クライアント処理
- **Hub/Room/Client** パターンによる接続管理

### フロントエンド

- **React 18** + TypeScript
- **Vite** - 高速ビルド
- **Custom Hooks** - WebSocket 統合

---

## 構成

```
CoScribe/
├── cmd/server/           # エントリーポイント
├── internal/
│   ├── ws/              # WebSocket処理
│   ├── document/        # 文書管理
│   └── database/        # データ
└── web/                 # React フロントエンド
    ├── src/components/
    ├── src/hooks/
    └── src/types/
```

---

## 起動

```bash
make dev-full

```

- フロントエンド: http://localhost:3000
- バックエンド: http://localhost:8080

---
