/**
 * @file barcode_ffi.h
 * @brief C FFI API for Barcode library (dart:ffi compatible)
 *
 * Provides a C-compatible API wrapping the C++ barcode library.
 * All functions use extern "C" linkage and opaque handles.
 * Designed for use with Flutter's dart:ffi.
 *
 * Supported barcode types (19):
 *   1D: Code39, Code93, Code128, GS1_128, NW7, Matrix2of5, NEC2of5,
 *       Jan8, Jan13, UPC_A, UPC_E, ITF
 *   GS1 DataBar: GS1DataBar14, GS1DataBarLimited, GS1DataBarExpanded
 *   Special: YubinCustomer
 *   2D: QR, DataMatrix, PDF417
 */

#ifndef BARCODE_FFI_H
#define BARCODE_FFI_H

#include <stdint.h>

#ifdef _WIN32
#define FFI_EXPORT __declspec(dllexport)
#else
#define FFI_EXPORT __attribute__((visibility("default")))
#endif

#ifdef __cplusplus
extern "C" {
#endif

/** Opaque handle to a barcode instance */
typedef void* BarcodeHandle;

/* =========================================================================
 * Barcode type IDs (used for factory creation)
 * ========================================================================= */
#define BC_CODE39                  0
#define BC_CODE93                  1
#define BC_CODE128                 2
#define BC_GS1_128                 3
#define BC_NW7                     4
#define BC_MATRIX2OF5              5
#define BC_NEC2OF5                 6
#define BC_JAN8                    7
#define BC_JAN13                   8
#define BC_UPC_A                   9
#define BC_UPC_E                  10
#define BC_ITF                    11
#define BC_GS1_DATABAR_14         12
#define BC_GS1_DATABAR_LIMITED    13
#define BC_GS1_DATABAR_EXPANDED   14
#define BC_YUBIN_CUSTOMER         15
#define BC_QR                     16
#define BC_DATAMATRIX             17
#define BC_PDF417                 18

/* =========================================================================
 * Create / Destroy
 * ========================================================================= */

/** Create a barcode instance by type ID. Returns NULL on invalid type. */
FFI_EXPORT BarcodeHandle barcode_create(int type_id);

/** Destroy a barcode instance and free resources. */
FFI_EXPORT void barcode_destroy(BarcodeHandle handle);

/* =========================================================================
 * Common settings (BarcodeBase)
 * ========================================================================= */

FFI_EXPORT void barcode_set_output_format(BarcodeHandle h, const char* format);
FFI_EXPORT void barcode_set_foreground_color(BarcodeHandle h, int r, int g, int b, int a);
FFI_EXPORT void barcode_set_background_color(BarcodeHandle h, int r, int g, int b, int a);
FFI_EXPORT void barcode_set_px_adjust_black(BarcodeHandle h, int adjust);
FFI_EXPORT void barcode_set_px_adjust_white(BarcodeHandle h, int adjust);
FFI_EXPORT void barcode_set_fit_width(BarcodeHandle h, int fit);

/* =========================================================================
 * 1D barcode settings (BarcodeBase1D)
 * ========================================================================= */

FFI_EXPORT void barcode_set_show_text(BarcodeHandle h, int show);
FFI_EXPORT void barcode_set_text_font_scale(BarcodeHandle h, double scale);
FFI_EXPORT void barcode_set_text_gap(BarcodeHandle h, double scale);
FFI_EXPORT void barcode_set_text_even_spacing(BarcodeHandle h, int even);

/* =========================================================================
 * 2D barcode settings (BarcodeBase2D)
 * ========================================================================= */

FFI_EXPORT void barcode_set_string_encoding(BarcodeHandle h, const char* encoding);

/* =========================================================================
 * Type-specific settings
 * ========================================================================= */

/* Code39, NW7 */
FFI_EXPORT void barcode_set_show_start_stop(BarcodeHandle h, int show);

/* Code128 */
FFI_EXPORT void barcode_set_code_mode(BarcodeHandle h, const char* mode);

/* Jan8, Jan13, UPC_A, UPC_E */
FFI_EXPORT void barcode_set_extended_guard(BarcodeHandle h, int extended);

/* QR */
FFI_EXPORT void barcode_set_error_correction_level(BarcodeHandle h, const char* level);
FFI_EXPORT void barcode_set_version(BarcodeHandle h, int version);
FFI_EXPORT void barcode_set_encode_mode(BarcodeHandle h, const char* mode);

/* DataMatrix */
FFI_EXPORT void barcode_set_code_size(BarcodeHandle h, const char* size);
FFI_EXPORT void barcode_set_encode_scheme(BarcodeHandle h, const char* scheme);

/* PDF417 */
FFI_EXPORT void barcode_set_error_level(BarcodeHandle h, int level);
FFI_EXPORT void barcode_set_columns(BarcodeHandle h, int columns);
FFI_EXPORT void barcode_set_rows(BarcodeHandle h, int rows);
FFI_EXPORT void barcode_set_aspect_ratio(BarcodeHandle h, double ratio);
FFI_EXPORT void barcode_set_y_height(BarcodeHandle h, int y_height);

/* GS1 DataBar 14 */
FFI_EXPORT void barcode_set_symbol_type_14(BarcodeHandle h, const char* type);
FFI_EXPORT const char* barcode_get_symbol_type_14(BarcodeHandle h);
FFI_EXPORT int barcode_encode_14(BarcodeHandle h, const char* content);
FFI_EXPORT const char* barcode_calculate_check_digit_14(const char* src);

/* GS1 DataBar Expanded */
FFI_EXPORT void barcode_set_symbol_type_exp(BarcodeHandle h, const char* type);
FFI_EXPORT void barcode_set_no_of_columns(BarcodeHandle h, int columns);

/* =========================================================================
 * Draw functions
 * ========================================================================= */

/** Draw 1D barcode (code, width, height). Returns 1 on success. */
FFI_EXPORT int barcode_draw_1d(BarcodeHandle h, const char* code, int width, int height);

/** Draw 2D barcode (code, size). Returns 1 on success. */
FFI_EXPORT int barcode_draw_2d(BarcodeHandle h, const char* code, int size);

/** Draw 2D barcode (code, width, height). Returns 1 on success. */
FFI_EXPORT int barcode_draw_2d_rect(BarcodeHandle h, const char* code, int width, int height);

/** Draw YubinCustomer (code, height only). Returns 1 on success. */
FFI_EXPORT int barcode_draw_yubin(BarcodeHandle h, const char* code, int height);

/** Draw YubinCustomer with width (code, width, height). Returns 1 on success. */
FFI_EXPORT int barcode_draw_yubin_with_width(BarcodeHandle h, const char* code, int width, int height);

/** Draw GS1-128 convenience barcode. Returns 1 on success. */
FFI_EXPORT int barcode_draw_convenience(BarcodeHandle h, const char* code, int width, int height);

/** Draw GS1 DataBar Expanded stacked. Returns 1 on success. */
FFI_EXPORT int barcode_draw_stacked(BarcodeHandle h, const char* code, int width, int height);

/* =========================================================================
 * Get results (call after a successful draw)
 * ========================================================================= */

/** Get Base64-encoded image string. Returns empty string on failure. */
FFI_EXPORT const char* barcode_get_base64(BarcodeHandle h);

/** Get SVG string. Returns empty string on failure. */
FFI_EXPORT const char* barcode_get_svg(BarcodeHandle h);

/** Get raw image data. Sets *out_size. Returns NULL on failure. */
FFI_EXPORT const uint8_t* barcode_get_image_data(BarcodeHandle h, int* out_size);

/** Check if current output mode is SVG. Returns 1 if SVG. */
FFI_EXPORT int barcode_is_svg_output(BarcodeHandle h);

#ifdef __cplusplus
}
#endif

#endif /* BARCODE_FFI_H */
