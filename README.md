# pomerium-route-id-gen

Pomeriumの構成ファイルから内部的な `route_id` を事前に算出するためのCLIツールです。

## インストール

```bash
go build -o pomerium-route-id-gen ./cmd/pomerium-route-id-gen
```

## 使い方

```bash
./pomerium-route-id-gen --config config.yaml
```

### オプション

- `--config <path>`: Pomeriumの構成ファイル（YAML/JSON）へのパス（必須）
- `--output, -o <format>`: 出力形式を指定します。 `text` または `json` （デフォルト: `text`）

### 使用例

#### テキスト出力（デフォルト）

```bash
$ ./pomerium-route-id-gen --config config.yaml
RouteID: fd1bc82bb90632f6
From: https://a.example.com
---
RouteID: 50e0bc30f656b7e4
From: https://b.example.com
```

#### JSON出力

```bash
$ ./pomerium-route-id-gen --config config.yaml --output json
[
  {
    "route_id": "fd1bc82bb90632f6",
    "from": "https://a.example.com"
  },
  {
    "route_id": "50e0bc30f656b7e4",
    "from": "https://b.example.com"
  }
]
```

## 開発

### テストの実行

```bash
make test
```

### ビルド

```bash
make build
```
