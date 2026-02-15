// Package barcode_pao provides Go wrappers for the barcode C++ native FFI library.
// It uses the C++ barcode engine directly via FFI for high-speed barcode generation.
//
// Architecture:
//
//	Go code → syscall → barcode_pao.dll/so/dylib → C++ engine
package barcode_pao

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

// Output format constants.
const (
	FormatPNG  = "png"
	FormatJPEG = "jpg"
	FormatSVG  = "svg"
)

// ─── Native library loading ────────────────────────────────────────────────

var (
	libOnce sync.Once
	libErr  error

	procCreate  *syscall.LazyProc
	procDestroy *syscall.LazyProc

	// Common settings
	procSetOutputFormat    *syscall.LazyProc
	procSetForegroundColor *syscall.LazyProc
	procSetBackgroundColor *syscall.LazyProc
	procSetPxAdjustBlack   *syscall.LazyProc
	procSetPxAdjustWhite   *syscall.LazyProc
	procSetFitWidth        *syscall.LazyProc

	// 1D settings
	procSetShowText         *syscall.LazyProc
	procSetTextFontScale    *syscall.LazyProc
	procSetTextGap          *syscall.LazyProc
	procSetTextEvenSpacing  *syscall.LazyProc

	// 2D settings
	procSetStringEncoding *syscall.LazyProc

	// Type-specific settings
	procSetShowStartStop          *syscall.LazyProc
	procSetCodeMode               *syscall.LazyProc
	procSetExtendedGuard          *syscall.LazyProc
	procSetErrorCorrectionLevel   *syscall.LazyProc
	procSetVersion                *syscall.LazyProc
	procSetEncodeMode             *syscall.LazyProc
	procSetCodeSize               *syscall.LazyProc
	procSetEncodeScheme           *syscall.LazyProc
	procSetErrorLevel             *syscall.LazyProc
	procSetColumns                *syscall.LazyProc
	procSetRows                   *syscall.LazyProc
	procSetAspectRatio            *syscall.LazyProc
	procSetYHeight                *syscall.LazyProc
	procSetSymbolType14           *syscall.LazyProc
	procSetSymbolTypeExp          *syscall.LazyProc
	procSetNoOfColumns            *syscall.LazyProc

	// Draw functions
	procDraw1D            *syscall.LazyProc
	procDraw2D            *syscall.LazyProc
	procDraw2DRect        *syscall.LazyProc
	procDrawYubin         *syscall.LazyProc
	procDrawYubinWithWidth *syscall.LazyProc

	// Get results
	procGetBase64    *syscall.LazyProc
	procGetSvg       *syscall.LazyProc
	procIsSvgOutput  *syscall.LazyProc
)

func getNativeDir() string {
	// 1. Try relative to this source file (development time)
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		dir := filepath.Join(filepath.Dir(thisFile), "native")
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	// 2. Try relative to executable
	exePath, err := os.Executable()
	if err == nil {
		dir := filepath.Join(filepath.Dir(exePath), "native")
		if info, err2 := os.Stat(dir); err2 == nil && info.IsDir() {
			return dir
		}
		// Try same directory as executable
		dir = filepath.Dir(exePath)
		if _, err2 := os.Stat(filepath.Join(dir, "barcode_pao.dll")); err2 == nil {
			return dir
		}
	}
	return "native"
}

func loadLibrary() error {
	libOnce.Do(func() {
		nativeDir := getNativeDir()

		// Preload dependent DLLs
		for _, dep := range []string{"SDL2.dll", "SDL2_image.dll", "SDL2_ttf.dll"} {
			depPath := filepath.Join(nativeDir, dep)
			if _, err := os.Stat(depPath); err == nil {
				syscall.LoadDLL(depPath)
			}
		}

		dllPath := filepath.Join(nativeDir, "barcode_pao.dll")
		dll := syscall.NewLazyDLL(dllPath)

		// Bind all functions
		procCreate = dll.NewProc("barcode_create")
		procDestroy = dll.NewProc("barcode_destroy")

		procSetOutputFormat = dll.NewProc("barcode_set_output_format")
		procSetForegroundColor = dll.NewProc("barcode_set_foreground_color")
		procSetBackgroundColor = dll.NewProc("barcode_set_background_color")
		procSetPxAdjustBlack = dll.NewProc("barcode_set_px_adjust_black")
		procSetPxAdjustWhite = dll.NewProc("barcode_set_px_adjust_white")
		procSetFitWidth = dll.NewProc("barcode_set_fit_width")

		procSetShowText = dll.NewProc("barcode_set_show_text")
		procSetTextFontScale = dll.NewProc("barcode_set_text_font_scale")
		procSetTextGap = dll.NewProc("barcode_set_text_gap")
		procSetTextEvenSpacing = dll.NewProc("barcode_set_text_even_spacing")

		procSetStringEncoding = dll.NewProc("barcode_set_string_encoding")

		procSetShowStartStop = dll.NewProc("barcode_set_show_start_stop")
		procSetCodeMode = dll.NewProc("barcode_set_code_mode")
		procSetExtendedGuard = dll.NewProc("barcode_set_extended_guard")
		procSetErrorCorrectionLevel = dll.NewProc("barcode_set_error_correction_level")
		procSetVersion = dll.NewProc("barcode_set_version")
		procSetEncodeMode = dll.NewProc("barcode_set_encode_mode")
		procSetCodeSize = dll.NewProc("barcode_set_code_size")
		procSetEncodeScheme = dll.NewProc("barcode_set_encode_scheme")
		procSetErrorLevel = dll.NewProc("barcode_set_error_level")
		procSetColumns = dll.NewProc("barcode_set_columns")
		procSetRows = dll.NewProc("barcode_set_rows")
		procSetAspectRatio = dll.NewProc("barcode_set_aspect_ratio")
		procSetYHeight = dll.NewProc("barcode_set_y_height")
		procSetSymbolType14 = dll.NewProc("barcode_set_symbol_type_14")
		procSetSymbolTypeExp = dll.NewProc("barcode_set_symbol_type_exp")
		procSetNoOfColumns = dll.NewProc("barcode_set_no_of_columns")

		procDraw1D = dll.NewProc("barcode_draw_1d")
		procDraw2D = dll.NewProc("barcode_draw_2d")
		procDraw2DRect = dll.NewProc("barcode_draw_2d_rect")
		procDrawYubin = dll.NewProc("barcode_draw_yubin")
		procDrawYubinWithWidth = dll.NewProc("barcode_draw_yubin_with_width")

		procGetBase64 = dll.NewProc("barcode_get_base64")
		procGetSvg = dll.NewProc("barcode_get_svg")
		procIsSvgOutput = dll.NewProc("barcode_is_svg_output")

		// Verify the DLL can be loaded
		if err := procCreate.Find(); err != nil {
			libErr = fmt.Errorf("failed to load barcode native library from %s: %w", dllPath, err)
		}
	})
	return libErr
}

// ─── Helper functions ──────────────────────────────────────────────────────

func toPtr(s string) uintptr {
	b := append([]byte(s), 0)
	return uintptr(unsafe.Pointer(&b[0]))
}

func fromPtr(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	// Read null-terminated C string
	var buf []byte
	for i := 0; ; i++ {
		b := *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
		if b == 0 {
			break
		}
		buf = append(buf, b)
	}
	return string(buf)
}

func boolToInt(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

// ═════════════════════════════════════════════════════════════════════════════
// Base types
// ═════════════════════════════════════════════════════════════════════════════

// BarcodeBase holds the native handle for all barcode types.
type BarcodeBase struct {
	handle       uintptr
	outputFormat string
}

func newBarcodeBase(typeID int, outputFormat string) (*BarcodeBase, error) {
	if err := loadLibrary(); err != nil {
		return nil, err
	}
	handle, _, _ := procCreate.Call(uintptr(typeID))
	if handle == 0 {
		return nil, fmt.Errorf("failed to create barcode handle for type %d", typeID)
	}
	b := &BarcodeBase{handle: handle, outputFormat: outputFormat}
	b.SetOutputFormat(outputFormat)
	runtime.SetFinalizer(b, func(b *BarcodeBase) {
		if b.handle != 0 {
			procDestroy.Call(b.handle)
			b.handle = 0
		}
	})
	return b, nil
}

// SetOutputFormat sets the output format (png, jpg, svg).
func (b *BarcodeBase) SetOutputFormat(format string) {
	b.outputFormat = format
	procSetOutputFormat.Call(b.handle, toPtr(format))
}

// SetForegroundColor sets the foreground color (RGBA).
func (b *BarcodeBase) SetForegroundColor(r, g, bl, a int) {
	procSetForegroundColor.Call(b.handle, uintptr(r), uintptr(g), uintptr(bl), uintptr(a))
}

// SetBackgroundColor sets the background color (RGBA).
func (b *BarcodeBase) SetBackgroundColor(r, g, bl, a int) {
	procSetBackgroundColor.Call(b.handle, uintptr(r), uintptr(g), uintptr(bl), uintptr(a))
}

func (b *BarcodeBase) getResult() (string, error) {
	isSvg, _, _ := procIsSvgOutput.Call(b.handle)
	if isSvg == 1 {
		ptr, _, _ := procGetSvg.Call(b.handle)
		return fromPtr(ptr), nil
	}
	ptr, _, _ := procGetBase64.Call(b.handle)
	return fromPtr(ptr), nil
}

// Barcode1DBase provides common 1D barcode settings.
type Barcode1DBase struct {
	BarcodeBase
}

// SetShowText sets whether to show text below the barcode.
func (b *Barcode1DBase) SetShowText(show bool) {
	procSetShowText.Call(b.handle, boolToInt(show))
}

// SetTextGap sets the gap between barcode and text.
func (b *Barcode1DBase) SetTextGap(gap float64) {
	procSetTextGap.Call(b.handle, uintptr(*(*uint64)(unsafe.Pointer(&gap))))
}

// SetTextFontScale sets the text font scale.
func (b *Barcode1DBase) SetTextFontScale(scale float64) {
	procSetTextFontScale.Call(b.handle, uintptr(*(*uint64)(unsafe.Pointer(&scale))))
}

// SetTextEvenSpacing sets text even spacing mode.
func (b *Barcode1DBase) SetTextEvenSpacing(even bool) {
	procSetTextEvenSpacing.Call(b.handle, boolToInt(even))
}

// SetFitWidth sets whether to fit the barcode to width.
func (b *Barcode1DBase) SetFitWidth(fit bool) {
	procSetFitWidth.Call(b.handle, boolToInt(fit))
}

// SetPxAdjustBlack sets pixel adjustment for black bars.
func (b *Barcode1DBase) SetPxAdjustBlack(adj int) {
	procSetPxAdjustBlack.Call(b.handle, uintptr(adj))
}

// SetPxAdjustWhite sets pixel adjustment for white bars.
func (b *Barcode1DBase) SetPxAdjustWhite(adj int) {
	procSetPxAdjustWhite.Call(b.handle, uintptr(adj))
}

// Draw generates a 1D barcode and returns Base64 or SVG string.
func (b *Barcode1DBase) Draw(code string, width, height int) (string, error) {
	ret, _, _ := procDraw1D.Call(b.handle, toPtr(code), uintptr(width), uintptr(height))
	if ret != 1 {
		return "", fmt.Errorf("draw failed")
	}
	return b.getResult()
}

// Barcode2DBase provides common 2D barcode settings.
type Barcode2DBase struct {
	BarcodeBase
}

// SetStringEncoding sets the string encoding (utf-8, shift-jis).
func (b *Barcode2DBase) SetStringEncoding(enc string) {
	procSetStringEncoding.Call(b.handle, toPtr(enc))
}

// SetFitWidth sets whether to fit the barcode to width.
func (b *Barcode2DBase) SetFitWidth(fit bool) {
	procSetFitWidth.Call(b.handle, boolToInt(fit))
}

// Draw generates a 2D barcode and returns Base64 or SVG string.
func (b *Barcode2DBase) Draw(code string, size int) (string, error) {
	ret, _, _ := procDraw2D.Call(b.handle, toPtr(code), uintptr(size))
	if ret != 1 {
		return "", fmt.Errorf("draw failed")
	}
	return b.getResult()
}

// ═════════════════════════════════════════════════════════════════════════════
// 1D Barcodes
// ═════════════════════════════════════════════════════════════════════════════

// Code39 generates Code39 barcodes.
type Code39 struct{ Barcode1DBase }

// NewCode39 creates a Code39 barcode generator.
func NewCode39(outputFormat string) *Code39 {
	base, err := newBarcodeBase(0, outputFormat)
	if err != nil {
		panic(err)
	}
	return &Code39{Barcode1DBase{*base}}
}

// SetShowStartStop sets whether to show start/stop characters.
func (b *Code39) SetShowStartStop(show bool) {
	procSetShowStartStop.Call(b.handle, boolToInt(show))
}

// Code93 generates Code93 barcodes.
type Code93 struct{ Barcode1DBase }

// NewCode93 creates a Code93 barcode generator.
func NewCode93(outputFormat string) *Code93 {
	base, err := newBarcodeBase(1, outputFormat)
	if err != nil {
		panic(err)
	}
	return &Code93{Barcode1DBase{*base}}
}

// Code128 generates Code128 barcodes.
type Code128 struct{ Barcode1DBase }

// NewCode128 creates a Code128 barcode generator.
func NewCode128(outputFormat string) *Code128 {
	base, err := newBarcodeBase(2, outputFormat)
	if err != nil {
		panic(err)
	}
	return &Code128{Barcode1DBase{*base}}
}

// SetCodeMode sets the code mode (AUTO, A, B, C).
func (b *Code128) SetCodeMode(mode string) {
	procSetCodeMode.Call(b.handle, toPtr(mode))
}

// GS1128 generates GS1-128 barcodes.
type GS1128 struct{ Barcode1DBase }

// NewGS1128 creates a GS1-128 barcode generator.
func NewGS1128(outputFormat string) *GS1128 {
	base, err := newBarcodeBase(3, outputFormat)
	if err != nil {
		panic(err)
	}
	return &GS1128{Barcode1DBase{*base}}
}

// NW7 generates NW-7 (Codabar) barcodes.
type NW7 struct{ Barcode1DBase }

// NewNW7 creates a NW-7 barcode generator.
func NewNW7(outputFormat string) *NW7 {
	base, err := newBarcodeBase(4, outputFormat)
	if err != nil {
		panic(err)
	}
	return &NW7{Barcode1DBase{*base}}
}

// SetShowStartStop sets whether to show start/stop characters.
func (b *NW7) SetShowStartStop(show bool) {
	procSetShowStartStop.Call(b.handle, boolToInt(show))
}

// ITF generates ITF (Interleaved 2 of 5) barcodes.
type ITF struct{ Barcode1DBase }

// NewITF creates an ITF barcode generator.
func NewITF(outputFormat string) *ITF {
	base, err := newBarcodeBase(11, outputFormat)
	if err != nil {
		panic(err)
	}
	return &ITF{Barcode1DBase{*base}}
}

// Matrix2of5 generates Matrix 2 of 5 barcodes.
type Matrix2of5 struct{ Barcode1DBase }

// NewMatrix2of5 creates a Matrix 2 of 5 barcode generator.
func NewMatrix2of5(outputFormat string) *Matrix2of5 {
	base, err := newBarcodeBase(5, outputFormat)
	if err != nil {
		panic(err)
	}
	return &Matrix2of5{Barcode1DBase{*base}}
}

// NEC2of5 generates NEC 2 of 5 barcodes.
type NEC2of5 struct{ Barcode1DBase }

// NewNEC2of5 creates a NEC 2 of 5 barcode generator.
func NewNEC2of5(outputFormat string) *NEC2of5 {
	base, err := newBarcodeBase(6, outputFormat)
	if err != nil {
		panic(err)
	}
	return &NEC2of5{Barcode1DBase{*base}}
}

// Jan8 generates JAN-8 (EAN-8) barcodes.
type Jan8 struct{ Barcode1DBase }

// NewJAN8 creates a JAN-8 barcode generator.
func NewJAN8(outputFormat string) *Jan8 {
	base, err := newBarcodeBase(7, outputFormat)
	if err != nil {
		panic(err)
	}
	return &Jan8{Barcode1DBase{*base}}
}

// SetExtendedGuard sets whether to use extended guard bars.
func (b *Jan8) SetExtendedGuard(ext bool) {
	procSetExtendedGuard.Call(b.handle, boolToInt(ext))
}

// Jan13 generates JAN-13 (EAN-13) barcodes.
type Jan13 struct{ Barcode1DBase }

// NewJAN13 creates a JAN-13 barcode generator.
func NewJAN13(outputFormat string) *Jan13 {
	base, err := newBarcodeBase(8, outputFormat)
	if err != nil {
		panic(err)
	}
	return &Jan13{Barcode1DBase{*base}}
}

// SetExtendedGuard sets whether to use extended guard bars.
func (b *Jan13) SetExtendedGuard(ext bool) {
	procSetExtendedGuard.Call(b.handle, boolToInt(ext))
}

// UPCA generates UPC-A barcodes.
type UPCA struct{ Barcode1DBase }

// NewUPCA creates a UPC-A barcode generator.
func NewUPCA(outputFormat string) *UPCA {
	base, err := newBarcodeBase(9, outputFormat)
	if err != nil {
		panic(err)
	}
	return &UPCA{Barcode1DBase{*base}}
}

// SetExtendedGuard sets whether to use extended guard bars.
func (b *UPCA) SetExtendedGuard(ext bool) {
	procSetExtendedGuard.Call(b.handle, boolToInt(ext))
}

// UPCE generates UPC-E barcodes.
type UPCE struct{ Barcode1DBase }

// NewUPCE creates a UPC-E barcode generator.
func NewUPCE(outputFormat string) *UPCE {
	base, err := newBarcodeBase(10, outputFormat)
	if err != nil {
		panic(err)
	}
	return &UPCE{Barcode1DBase{*base}}
}

// SetExtendedGuard sets whether to use extended guard bars.
func (b *UPCE) SetExtendedGuard(ext bool) {
	procSetExtendedGuard.Call(b.handle, boolToInt(ext))
}

// ═════════════════════════════════════════════════════════════════════════════
// GS1 DataBar
// ═════════════════════════════════════════════════════════════════════════════

// GS1DataBar14 generates GS1 DataBar 14 barcodes.
type GS1DataBar14 struct{ Barcode1DBase }

// NewGS1DataBar14 creates a GS1 DataBar 14 barcode generator.
func NewGS1DataBar14(outputFormat string) *GS1DataBar14 {
	base, err := newBarcodeBase(12, outputFormat)
	if err != nil {
		panic(err)
	}
	return &GS1DataBar14{Barcode1DBase{*base}}
}

// SetSymbolType sets the symbol type (OMNIDIRECTIONAL, STACKED, STACKED_OMNIDIRECTIONAL).
func (b *GS1DataBar14) SetSymbolType(symbolType string) {
	procSetSymbolType14.Call(b.handle, toPtr(symbolType))
}

// GS1DataBarLimited generates GS1 DataBar Limited barcodes.
type GS1DataBarLimited struct{ Barcode1DBase }

// NewGS1DataBarLimited creates a GS1 DataBar Limited barcode generator.
func NewGS1DataBarLimited(outputFormat string) *GS1DataBarLimited {
	base, err := newBarcodeBase(13, outputFormat)
	if err != nil {
		panic(err)
	}
	return &GS1DataBarLimited{Barcode1DBase{*base}}
}

// GS1DataBarExpanded generates GS1 DataBar Expanded barcodes.
type GS1DataBarExpanded struct{ Barcode1DBase }

// NewGS1DataBarExpanded creates a GS1 DataBar Expanded barcode generator.
func NewGS1DataBarExpanded(outputFormat string) *GS1DataBarExpanded {
	base, err := newBarcodeBase(14, outputFormat)
	if err != nil {
		panic(err)
	}
	return &GS1DataBarExpanded{Barcode1DBase{*base}}
}

// SetSymbolType sets the symbol type (UNSTACKED, STACKED).
func (b *GS1DataBarExpanded) SetSymbolType(symbolType string) {
	procSetSymbolTypeExp.Call(b.handle, toPtr(symbolType))
}

// SetNoOfColumns sets the number of columns for stacked version.
func (b *GS1DataBarExpanded) SetNoOfColumns(cols int) {
	procSetNoOfColumns.Call(b.handle, uintptr(cols))
}

// ═════════════════════════════════════════════════════════════════════════════
// Special Barcodes
// ═════════════════════════════════════════════════════════════════════════════

// YubinCustomer generates Japanese postal customer barcodes.
type YubinCustomer struct {
	BarcodeBase
}

// NewYubinCustomer creates a YubinCustomer barcode generator.
func NewYubinCustomer(outputFormat string) *YubinCustomer {
	base, err := newBarcodeBase(15, outputFormat)
	if err != nil {
		panic(err)
	}
	return &YubinCustomer{*base}
}

// SetPxAdjustBlack sets pixel adjustment for black bars.
func (b *YubinCustomer) SetPxAdjustBlack(adj int) {
	procSetPxAdjustBlack.Call(b.handle, uintptr(adj))
}

// SetPxAdjustWhite sets pixel adjustment for white bars.
func (b *YubinCustomer) SetPxAdjustWhite(adj int) {
	procSetPxAdjustWhite.Call(b.handle, uintptr(adj))
}

// Draw generates a postal barcode. Width is auto-calculated.
func (b *YubinCustomer) Draw(code string, height int) (string, error) {
	ret, _, _ := procDrawYubin.Call(b.handle, toPtr(code), uintptr(height))
	if ret != 1 {
		return "", fmt.Errorf("draw failed")
	}
	return b.getResult()
}

// DrawWithWidth generates a postal barcode with explicit width.
func (b *YubinCustomer) DrawWithWidth(code string, width, height int) (string, error) {
	ret, _, _ := procDrawYubinWithWidth.Call(b.handle, toPtr(code), uintptr(width), uintptr(height))
	if ret != 1 {
		return "", fmt.Errorf("draw failed")
	}
	return b.getResult()
}

// ═════════════════════════════════════════════════════════════════════════════
// 2D Barcodes
// ═════════════════════════════════════════════════════════════════════════════

// QR generates QR codes.
type QR struct{ Barcode2DBase }

// NewQRCode creates a QR code generator.
func NewQRCode(outputFormat string) *QR {
	base, err := newBarcodeBase(16, outputFormat)
	if err != nil {
		panic(err)
	}
	return &QR{Barcode2DBase{*base}}
}

// SetErrorCorrectionLevel sets the error correction level (L, M, Q, H).
func (b *QR) SetErrorCorrectionLevel(level string) {
	procSetErrorCorrectionLevel.Call(b.handle, toPtr(level))
}

// SetVersion sets QR version (0=auto, 1-40).
func (b *QR) SetVersion(version int) {
	procSetVersion.Call(b.handle, uintptr(version))
}

// SetEncodeMode sets the encode mode (NUMERIC, ALPHANUMERIC, BYTE, KANJI).
func (b *QR) SetEncodeMode(mode string) {
	procSetEncodeMode.Call(b.handle, toPtr(mode))
}

// DataMatrix generates DataMatrix barcodes.
type DataMatrix struct{ Barcode2DBase }

// NewDataMatrix creates a DataMatrix barcode generator.
func NewDataMatrix(outputFormat string) *DataMatrix {
	base, err := newBarcodeBase(17, outputFormat)
	if err != nil {
		panic(err)
	}
	return &DataMatrix{Barcode2DBase{*base}}
}

// SetCodeSize sets the code size (AUTO, 10x10, 12x12, etc.).
func (b *DataMatrix) SetCodeSize(size string) {
	procSetCodeSize.Call(b.handle, toPtr(size))
}

// SetEncodeScheme sets the encode scheme (AUTO, ASCII, C40, TEXT, X12, EDIFACT, BASE256).
func (b *DataMatrix) SetEncodeScheme(scheme string) {
	procSetEncodeScheme.Call(b.handle, toPtr(scheme))
}

// PDF417 generates PDF417 barcodes.
type PDF417 struct{ Barcode2DBase }

// NewPDF417 creates a PDF417 barcode generator.
func NewPDF417(outputFormat string) *PDF417 {
	base, err := newBarcodeBase(18, outputFormat)
	if err != nil {
		panic(err)
	}
	return &PDF417{Barcode2DBase{*base}}
}

// SetErrorLevel sets the error correction level (-1=auto, 0-8).
func (b *PDF417) SetErrorLevel(level int) {
	procSetErrorLevel.Call(b.handle, uintptr(level))
}

// SetColumns sets the number of columns.
func (b *PDF417) SetColumns(cols int) {
	procSetColumns.Call(b.handle, uintptr(cols))
}

// SetRows sets the number of rows.
func (b *PDF417) SetRows(rows int) {
	procSetRows.Call(b.handle, uintptr(rows))
}

// SetAspectRatio sets the aspect ratio.
func (b *PDF417) SetAspectRatio(ratio float64) {
	procSetAspectRatio.Call(b.handle, uintptr(*(*uint64)(unsafe.Pointer(&ratio))))
}

// SetYHeight sets the Y height.
func (b *PDF417) SetYHeight(yHeight int) {
	procSetYHeight.Call(b.handle, uintptr(yHeight))
}

// Draw generates a PDF417 barcode (width × height).
func (b *PDF417) Draw(code string, width, height int) (string, error) {
	ret, _, _ := procDraw2DRect.Call(b.handle, toPtr(code), uintptr(width), uintptr(height))
	if ret != 1 {
		return "", fmt.Errorf("draw failed")
	}
	return b.getResult()
}

// ═════════════════════════════════════════════════════════════════════════════
// Product Info
// ═════════════════════════════════════════════════════════════════════════════

// GetProductName returns the product name.
func GetProductName() string { return "barcode-pao (Go)" }

// GetVersion returns the version.
func GetVersion() string { return "0.0.1" }

// GetManufacturer returns the manufacturer.
func GetManufacturer() string { return "Pao" }
