# barcode-pao-go

クロスプラットフォーム バーコード生成ライブラリ for Go（Native FFI版）

## 概要

`barcode-pao-go` は、C++ バーコードエンジンを FFI（Foreign Function Interface）で直接呼び出す高速なGoパッケージです。ネイティブコードを直接実行するため、高速なバーコード生成が可能です。

## 必要条件

- Go 1.21以上
- Windows（ネイティブDLL同梱）

## 対応バーコード（18種）

### 1次元バーコード（11種）
- **Code39** - 英数字対応の汎用バーコード
- **Code93** - Code39の拡張版
- **Code128** - 全ASCII文字対応の高密度バーコード
- **GS1-128** - 物流・流通向けバーコード（コンビニ収納代行対応）
- **NW-7 (Codabar)** - 血液銀行・宅配便向けバーコード
- **Matrix 2 of 5** - 工業用バーコード
- **NEC 2 of 5** - NECが開発した2 of 5系バーコード
- **JAN-8** - 日本の商品コード（8桁）
- **JAN-13** - 日本の商品コード（13桁）
- **UPC-A** - 北米の商品コード（12桁）
- **UPC-E** - UPC-Aの短縮版（8桁）

### GS1 DataBar（3種）
- **GS1 DataBar 14** - 標準型（オムニ/スタック対応）
- **GS1 DataBar Limited** - 限定型
- **GS1 DataBar Expanded** - 拡張型（スタック対応）

### 2次元バーコード（3種）
- **QRコード** - 日本発の2次元コード
- **DataMatrix** - 工業用途の2次元コード
- **PDF417** - 運転免許証等で使用される2次元コード

### 特殊バーコード（1種）
- **郵便カスタマバーコード** - 日本郵便の住所表示バーコード

## インストール

```bash
go get github.com/pao-xx/barcode-pao-go
```

## 使用例

### QRコード生成

```go
package main

import (
	"fmt"
	"os"

	barcode "github.com/pao-xx/barcode-pao-go"
)

func main() {
	// QRコードインスタンスを作成
	qr := barcode.NewQRCode(barcode.FormatPNG)

	// エラー訂正レベルを設定（L/M/Q/H）
	qr.SetErrorCorrectionLevel("H")

	// Base64エンコードされた画像を取得
	base64Image, err := qr.Draw("https://example.com", 200)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	fmt.Println(base64Image)
}
```

### Code128バーコード生成

```go
// Code128インスタンスを作成
code128 := barcode.NewCode128(barcode.FormatSVG)

// テキスト表示を有効化
code128.SetShowText(true)

// バーコード生成
svgData, err := code128.Draw("ABC-12345", 300, 100)
if err != nil {
	log.Fatal(err)
}
```

### 色のカスタマイズ

```go
code39 := barcode.NewCode39(barcode.FormatPNG)

// 前景色（バーの色）をRGBAで設定
code39.SetForegroundColor(0, 0, 128, 255) // 紺色

// 背景色をRGBAで設定
code39.SetBackgroundColor(255, 255, 200, 255) // 薄黄色

base64Image, err := code39.Draw("12345", 200, 80)
```

### GS1-128 コンビニ収納代行バーコード

```go
gs1 := barcode.NewGS1128(barcode.FormatPNG)
gs1.SetShowText(true)

// 標準料金代理収納用バーコード
convenienceCode := "9101234567890123456789012345678901234567890123"
base64Image, err := gs1.Draw(convenienceCode, 400, 100)
```

### 郵便カスタマバーコード

```go
yubin := barcode.NewYubinCustomer(barcode.FormatPNG)

// 郵便番号 + 住所表示番号
code := "1000001-1-2-3"
base64Image, err := yubin.Draw(code, 50) // 高さのみ指定
```

## API リファレンス

### 共通メソッド（全バーコードクラス）

| メソッド | 説明 |
|---------|------|
| `SetOutputFormat(format)` | 出力フォーマットを設定（"png", "jpg", "svg"）|
| `SetForegroundColor(r, g, b, a)` | 前景色（バーの色）を設定 |
| `SetBackgroundColor(r, g, b, a)` | 背景色を設定 |
| `Draw(code, width, height)` | Base64エンコードされた画像またはSVGを返す |

### 1次元バーコード固有メソッド

| メソッド | 説明 |
|---------|------|
| `SetShowText(show)` | バーコード下のテキスト表示 |
| `SetTextFontScale(scale)` | テキストのフォントサイズスケール |
| `SetTextGap(scale)` | バーとテキストの間隔 |
| `SetFitWidth(fit)` | 幅に合わせてバーを調整 |
| `SetPxAdjustBlack(adjust)` | 黒バーのピクセル調整 |
| `SetPxAdjustWhite(adjust)` | 白バーのピクセル調整 |

### 2次元バーコード固有メソッド

| クラス | メソッド | 説明 |
|--------|---------|------|
| QR | `SetErrorCorrectionLevel(level)` | エラー訂正レベル（L/M/Q/H）|
| QR | `SetVersion(version)` | バージョン（0=自動, 1-40）|
| QR | `SetEncodeMode(mode)` | エンコードモード（NUMERIC/ALPHANUMERIC/BYTE/KANJI）|
| DataMatrix | `SetCodeSize(size)` | シンボルサイズ（"AUTO", "10x10"など）|
| DataMatrix | `SetEncodeScheme(scheme)` | エンコードスキーム（AUTO/ASCII/C40/TEXT/X12/EDIFACT/BASE256）|
| PDF417 | `SetErrorLevel(level)` | エラー訂正レベル（-1=自動, 0-8）|
| PDF417 | `SetColumns(columns)` | 列数 |
| PDF417 | `SetRows(rows)` | 行数 |

## 出力フォーマット

| フォーマット | 説明 |
|-------------|------|
| `png` | PNG画像（デフォルト）|
| `jpg` / `jpeg` | JPEG画像 |
| `svg` | SVGベクター画像 |

## WASM版との違い

| | Native FFI版 | WASM版 |
|---|------------|--------|
| パッケージ | `barcode-pao-go` | `barcode-pao-wasm-go` |
| 実行方式 | C++ DLL/SO を直接呼び出し | Node.js 経由で WASM 実行 |
| 速度 | 高速 | やや遅い |
| 依存 | ネイティブDLL | Node.js |
| API | 同じ | 同じ |

## ライセンス

MIT License

## 関連パッケージ

- [barcode-pao-wasm-go (pkg.go.dev)](https://pkg.go.dev/github.com/pao-xx/barcode-pao-wasm-go) - Go WASM版
- [barcode-pao-wasm (Python)](https://pypi.org/project/barcode-pao-wasm/) - Python WASM版
- [barcode-pao-wasm (Rust)](https://crates.io/crates/barcode-pao-wasm) - Rust WASM版
