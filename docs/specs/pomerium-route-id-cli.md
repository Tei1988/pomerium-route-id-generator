# 概要
Pomeriumの構成ファイル（config.yaml）を読み込み、定義されている各ルートに対応するユニークな route_id を計算して出力するCLIツール。

## 計算ロジック
route_id の計算は、Pomerium内部の pkg/config/policy.go に実装されている generateRouteID メソッドと同一のアルゴリズム（xxh3ハッシュ）を使用しなければならない。

ハッシュ計算の入力順序:
1. from (文字列)
1. prefix (文字列)
1. path (文字列)
1. regex (文字列)
1. アクション種別コード (1バイト):
  - To (Weighted URLs) が設定されている場合: 1
  - Redirect が設定されている場合: 2
  - Response が設定されている場合: 3
1. アクションの詳細データ:
   - To の場合: URLの数、および各URLの構成要素（Scheme, Host, Path, LbWeightなど）を特定の順序で書き込む。
   - Redirect の場合: 各フィールド（HTTPSRedirect, SchemeRedirectなど）をポインタの有無を含めて書き込む。
   - Response の場合: ステータスコードとボディ内容を書き込む。

## CLI インターフェース仕様
- コマンド名: pomerium-route-id-gen
- 必須フラグ:
  - --config <path>: Pomeriumの構成ファイル（YAML/JSON）を指定する。
- オプションフラグ:
  - --output, -o <format>: 出力のフォーマットをtextかJSONかを指定する。デフォルトはtext
- 動作:
   1. 指定されたファイルを Options 構造体としてパースする。
   1. routes および policy セクションに含まれるすべてのポリシーを走査する。
   1. 各ポリシーに対して上記のハッシュ計算を実行し、16進数文字列（Hex）に変換する。
- 出力フォーマット: 標準出力に指定された出力形式で表示する。
   ```text
   RouteID: <route-id>
   From: <from>
   ---
   RouteID: <route-id>
   From: <from>
   ```
   ```json
   [
     {
       "route_id": "<route-id>",
       "from": "<from>"
     },
     {
       "route_id": "<route-id>",
       "from": "<from>"
     }
   ]
   ```

## 実装上の注意
- Pomeriumの既存パッケージ（github.com/pomerium/pomerium/pkg/config など）を再利用して、計算ロジックの乖離を防ぐこと。
- パース時にYAMLの構文エラーがある場合は、適切にエラーを表示すること。
