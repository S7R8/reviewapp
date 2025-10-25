# ReviewApp

**「使えば使うほど、あなたらしくなるAIアシスタント」**

AIコードレビューツールですが、あなたの過去の判断基準やコーディングスタイルを学習し、一貫性のある「あなたらしい」レビューを提供します。

## 🎯 コアコンセプト

- **ナレッジによる一貫性**: 過去のあなたの判断基準を常に参照
- **成長するAI**: 使うほど賢くなる、あなた専用にカスタマイズ
- **あなたのクローン**: 「あなただったらこう言う」というレビュー

## 🚀 クイックスタート（Windows環境）

### 前提条件

1. **Docker Desktop for Windows**
   - [Docker Desktop](https://www.docker.com/products/docker-desktop/)をインストール
   - WSL2バックエンドを有効化

2. **VS Code**
   - [VS Code](https://code.visualstudio.com/)をインストール
   - 拡張機能「Dev Containers」をインストール

3. **Git for Windows**（オプション）
   - [Git](https://git-scm.com/download/win)をインストール

### セットアップ手順

#### 1. プロジェクトを配置
```powershell
# すでに C:\reviewApp にあることを確認
cd C:\reviewApp
dir