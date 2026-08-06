package io

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/caerbannogwhite/enchanter"
)

// https://cran.r-project.org/web/packages/sas7bdat/vignettes/sas7bdat.pdf
// Example files: https://www2.census.gov/programs-surveys/ahs/2013/AHS%202013%20National%20PUF%20v2.0%20SAS.zip

const _MIN_READ_SIZE = 512

var _MAGIC_NUMBER = []byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xc2, 0xea, 0x81, 0x60,
	0xb3, 0x14, 0x11, 0xcf, 0xbd, 0x92, 0x08, 0x00,
	0x09, 0xc7, 0x31, 0x8c, 0x18, 0x1f, 0x10, 0x11,
}

type Sas7bdatReader struct {
	path   string
	reader io.ReadSeeker
	ctx    *enchanter.Context
	iod    *IoData

	is64Bit    bool
	offset     int
	byteOrder  binary.ByteOrder
	header     *_SAS7BDATFileHeader
	pageSize   int64
	pageCount  int64
	headerSize int64
	pages      []*_SAS7BDATPage
	columns    []*_SAS7BDATColumn
}

func NewSas7bdatReader(ctx *enchanter.Context) *Sas7bdatReader {
	return &Sas7bdatReader{
		path:   "",
		reader: nil,
		ctx:    ctx,
	}
}

func (r *Sas7bdatReader) SetPath(path string) *Sas7bdatReader {
	r.path = path
	return r
}

func (r *Sas7bdatReader) SetReader(reader io.ReadSeeker) *Sas7bdatReader {
	r.reader = reader
	return r
}

func (r *Sas7bdatReader) Read() *IoData {
	if r.reader == nil {
		file, err := os.Open(r.path)
		if err != nil {
			return &IoData{Error: fmt.Errorf("failed to open file: %w", err)}
		}
		r.reader = file
	}

	return r.readSas7bdata()
}

func (r *Sas7bdatReader) readBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)

	bufferSize := _MIN_READ_SIZE
	if bufferSize > n {
		bufferSize = n
	}

	bytesRead := 0
	buffer := make([]byte, bufferSize)
	for {
		b, err := r.reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				copy(bytes[bytesRead:], buffer)
				bytesRead += b
				break
			}
			return nil, err
		}

		copy(bytes[bytesRead:], buffer)
		bytesRead += b
		if bytesRead >= n {
			break
		}
	}

	// Update the offset
	r.offset += bytesRead

	return bytes, nil
}

func (r *Sas7bdatReader) getPointerSize() int {
	if r.is64Bit {
		return 8
	}
	return 4
}

func (r *Sas7bdatReader) readInteger(data []byte, size int) (int64, error) {
	switch size {
	case 1:
		return int64(data[0]), nil
	case 2:
		return int64(r.byteOrder.Uint16(data)), nil
	case 4:
		return int64(r.byteOrder.Uint32(data)), nil
	case 8:
		return int64(r.byteOrder.Uint64(data)), nil
	default:
		return 0, fmt.Errorf("unsupported integer size: %d", size)
	}
}

// PageType defines the different types of pages
type _SAS7BDATPageType int

const (
	PageTypeMeta _SAS7BDATPageType = iota
	PageTypeData
	PageTypeMix
	PageTypeAMD
)

// SubheaderType defines the different types of subheaders
type _SAS7BDATSubheaderType int

const (
	SubheaderRowSize _SAS7BDATSubheaderType = iota
	SubheaderColSize
	SubheaderSubhdrCounts
	SubheaderColText
	SubheaderColName
	SubheaderColAttrs
	SubheaderColFormatAndLabel
	SubheaderColList
	SubheaderUnknown
)

// ColumnType defines the supported data types
type _SAS7BDATColumnType int

const (
	ColumnTypeNumeric _SAS7BDATColumnType = iota
	ColumnTypeCharacter
	ColumnTypeDate
	ColumnTypeTime
	ColumnTypeDateTime
	ColumnTypeCurrency
)

// Page represents a single page in the SAS7BDAT file
type _SAS7BDATPage struct {
	Type        int16
	BlockCount  int16
	SubhdrCount int16
	PageType    _SAS7BDATPageType
	Subheaders  []*_SAS7BDATSubheader
	Data        []byte
	Deleted     bool
}

// Subheader represents metadata within a page
type _SAS7BDATSubheader struct {
	// Signature   _SAS7BDATSubheaderSignature
	// Compression _SAS7BDATCompressionType
	Type   _SAS7BDATSubheaderType
	Length int64
	Data   []byte
}

// Column represents a dataset column
type _SAS7BDATColumn struct {
	Index  int
	Name   string
	Label  string
	Format string
	Type   _SAS7BDATColumnType
	Length int
	Offset int64
}

// FileHeader contains the file header information
type _SAS7BDATFileHeader struct {
	a1 int64
	a2 int64

	sasRelease string
	sasServer  string
	osVersion  string
	osMaker    string
	osName     string
}

// PlatformInfo contains platform-specific parsing parameters
type _SAS7BDATPlatformInfo struct {
	Name        string
	Endianness  binary.ByteOrder
	Is64Bit     bool
	IsU64       bool
	PointerSize int
	IntSize     int
}

func (r *Sas7bdatReader) readSas7bdata() *IoData {
	r.iod = &IoData{}

	// Parse file header
	if err := r.parseHeader(); err != nil {
		return &IoData{Error: fmt.Errorf("header parsing failed: %w", err)}
	}

	// Parse pages
	if err := r.parsePages(); err != nil {
		return &IoData{Error: fmt.Errorf("page parsing failed: %w", err)}
	}

	return &IoData{Error: fmt.Errorf("not implemented")}
}

// parseHeader reads and parses the SAS7BDAT file header
func (r *Sas7bdatReader) parseHeader() error {
	const initialHeaderSize = 1024

	// Read the first 1024 bytes for initial header analysis
	headerBytes, err := r.readBytes(initialHeaderSize)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	r.header = &_SAS7BDATFileHeader{}

	// Check magic number (offset 0, length 32)
	if !r.validateMagicNumber(headerBytes[0:32]) {
		return fmt.Errorf("invalid magic number")
	}

	// Get alignment 2
	if headerBytes[32] == 0x33 {
		r.header.a2 = 4
		r.is64Bit = true
	} else {
		r.header.a2 = 0
		r.is64Bit = false
	}

	// Get alignment 1
	if headerBytes[35] == 0x33 {
		r.header.a1 = 4
	} else {
		r.header.a1 = 0
	}

	// Get endianness
	if headerBytes[37] == 0x00 {
		r.byteOrder = binary.BigEndian
	} else if headerBytes[37] == 0x01 {
		r.byteOrder = binary.LittleEndian
	} else {
		return fmt.Errorf("unable to determine byte order")
	}

	// Get OS type
	if headerBytes[39] == 0x01 {
		r.iod.FileMeta.SasOs = "UNIX"
	} else if headerBytes[39] == 0x02 {
		r.iod.FileMeta.SasOs = "WINDOWS"
	}

	// Get encoding (validated but not stored yet)
	if _, err := r.readInteger(headerBytes[70:72], 2); err != nil {
		return fmt.Errorf("failed to read encoding: %w", err)
	}

	// Get SAS file
	if strings.TrimSpace(string(headerBytes[84:92])) != "SAS FILE" {
		return fmt.Errorf("invalid SAS file")
	}

	// Get dataset name
	datasetName := string(headerBytes[92:156])
	r.iod.FileMeta.SasDsName = strings.TrimSpace(datasetName)

	// Get Date Created
	bits := r.byteOrder.Uint64(headerBytes[164+r.header.a1 : 172+r.header.a1])
	r.iod.FileMeta.Created = sasNumericToDateTime(
		math.Float64frombits(bits))

	// Get Date Modified
	bits = r.byteOrder.Uint64(headerBytes[172+r.header.a1 : 180+r.header.a1])
	r.iod.FileMeta.LastModified = sasNumericToDateTime(
		math.Float64frombits(bits))

	// Get Header Length
	r.headerSize, err = r.readInteger(headerBytes[196+r.header.a1:200+r.header.a1], 4)
	if err != nil {
		return fmt.Errorf("failed to read header length: %w", err)
	}

	// Get Page Size
	r.pageSize, err = r.readInteger(headerBytes[200+r.header.a1:204+r.header.a1], 4)
	if err != nil {
		return fmt.Errorf("failed to read page size: %w", err)
	}

	// Get Page Count
	if r.is64Bit {
		r.pageCount, err = r.readInteger(headerBytes[204+r.header.a1:212+r.header.a1], 8)
	} else {
		r.pageCount, err = r.readInteger(headerBytes[204+r.header.a1:208+r.header.a1], 4)
	}
	if err != nil {
		return fmt.Errorf("failed to read page count: %w", err)
	}

	a1a2 := r.header.a1 + r.header.a2
	r.header.sasRelease = string(headerBytes[216+a1a2 : 224+a1a2])
	r.header.sasServer = string(headerBytes[224+a1a2 : 240+a1a2])
	r.header.osVersion = string(headerBytes[240+a1a2 : 256+a1a2])
	r.header.osMaker = string(headerBytes[256+a1a2 : 272+a1a2])
	r.header.osName = string(headerBytes[272+a1a2 : 288+a1a2])

	// Read the rest of the header
	_, err = r.readBytes(int(r.headerSize) - initialHeaderSize)
	if err != nil {
		return fmt.Errorf("failed to read rest of header: %w", err)
	}

	return nil
}

func (r *Sas7bdatReader) validateMagicNumber(magic []byte) bool {
	for i, b := range magic {
		if _MAGIC_NUMBER[i] != b {
			return false
		}
	}
	return true
}

func (r *Sas7bdatReader) detectPlatform() (*_SAS7BDATPlatformInfo, error) {
	platform := &_SAS7BDATPlatformInfo{
		Endianness: r.byteOrder,
	}

	// Analyze header fields to determine platform characteristics
	if len(r.header.osName) > 0 {
		platform.Name = r.header.osName

		// Platform-specific configurations
		switch {
		case strings.Contains(platform.Name, "WIN"):
			platform.Is64Bit = r.is64Bit
			platform.PointerSize = r.getPointerSize()
		case strings.Contains(platform.Name, "Linux"):
			platform.Is64Bit = r.is64Bit
			platform.PointerSize = r.getPointerSize()
		case strings.Contains(platform.Name, "SunOS"):
			platform.IsU64 = true
			platform.Is64Bit = true
			platform.PointerSize = 8
			platform.Endianness = binary.BigEndian
		default:
			// Default assumptions for unknown platforms
			platform.PointerSize = 4
		}
	}

	return platform, nil
}

// parsePages reads and processes all pages in the file
func (r *Sas7bdatReader) parsePages() error {
	r.pages = make([]*_SAS7BDATPage, r.pageCount)

	for i := int64(0); i < r.pageCount; i++ {
		page, err := r.parsePage()
		if err != nil {
			return fmt.Errorf("failed to parse page %d: %w", i, err)
		}
		r.pages[i] = page
	}

	return nil
}

// parsePage reads and parses a single page from the file
func (r *Sas7bdatReader) parsePage() (*_SAS7BDATPage, error) {
	pageData, err := r.readBytes(int(r.pageSize))
	if err != nil {
		return nil, fmt.Errorf("failed to read page: %w", err)
	}

	page := &_SAS7BDATPage{
		Data: pageData,
	}

	// Parse page header (first 24 bytes typically)
	page.Type = int16(r.byteOrder.Uint16(pageData[16:18]))
	page.BlockCount = int16(r.byteOrder.Uint16(pageData[18:20]))
	page.SubhdrCount = int16(r.byteOrder.Uint16(pageData[20:22]))

	// Determine page type based on content
	// page.PageType = f.determinePageType(page)

	// Parse subheaders if present
	if page.SubhdrCount > 0 {
		if err := r.parseSubheaders(page); err != nil {
			return nil, err
		}
	}

	return page, nil
}

func (r *Sas7bdatReader) parseSubheaders(page *_SAS7BDATPage) error {
	offset := 24 // Standard page header size

	for i := 0; i < int(page.SubhdrCount); i++ {
		if offset+12 > len(page.Data) {
			break
		}

		subheader := &_SAS7BDATSubheader{}

		// Parse subheader pointer (offset, length, compression, type)
		subheaderOffset := r.byteOrder.Uint32(page.Data[offset : offset+4])
		subheaderLength := r.byteOrder.Uint32(page.Data[offset+4 : offset+8])
		// compressionInfo := page.Data[offset+8]
		// subheaderType := page.Data[offset+9]

		subheader.Length = int64(subheaderLength)
		// subheader.Compression = CompressionType(compressionInfo)
		// subheader.Type = f.identifySubheaderType(subheaderType)

		// Extract subheader data
		if subheaderOffset < uint32(len(page.Data)) &&
			subheaderOffset+subheaderLength <= uint32(len(page.Data)) {
			subheader.Data = page.Data[subheaderOffset : subheaderOffset+subheaderLength]
		}

		page.Subheaders = append(page.Subheaders, subheader)
		offset += 12 // Move to next subheader pointer
	}

	return nil
}

// NumericFormat represents different numeric storage formats
type NumericFormat struct {
	Size     int
	SignBits int
	ExpBits  int
	MantBits int
}

var numericFormats = map[int]NumericFormat{
	3: {24, 1, 11, 12}, // 24-bit format
	4: {32, 1, 11, 20}, // 32-bit format
	5: {40, 1, 11, 28}, // 40-bit format
	6: {48, 1, 11, 36}, // 48-bit format
	7: {56, 1, 11, 44}, // 56-bit format
	8: {64, 1, 11, 52}, // 64-bit format
}

func (r *Sas7bdatReader) parseNumericValue(data []byte, format NumericFormat) (float64, error) {
	if len(data) < format.Size {
		return 0, fmt.Errorf("insufficient data for numeric format")
	}

	// Handle different numeric formats based on size
	switch format.Size {
	case 8:
		// Standard IEEE 754 double precision
		bits := r.byteOrder.Uint64(data[:8])
		return math.Float64frombits(bits), nil
	case 4:
		// IEEE 754 single precision
		bits := r.byteOrder.Uint32(data[:4])
		return float64(math.Float32frombits(bits)), nil
	default:
		// Custom SAS numeric formats
		return r.parseCustomNumeric(data[:format.Size], format)
	}
}

func (r *Sas7bdatReader) parseCustomNumeric(data []byte, format NumericFormat) (float64, error) {
	// Extract sign, exponent, and mantissa based on format specification
	var value uint64
	for i, b := range data {
		value |= uint64(b) << (8 * (len(data) - 1 - i))
	}

	// Apply SAS-specific numeric conversion logic
	sign := (value >> (format.ExpBits + format.MantBits)) & 1
	exponent := (value >> format.MantBits) & ((1 << format.ExpBits) - 1)
	mantissa := value & ((1 << format.MantBits) - 1)

	// Convert to IEEE 754 equivalent
	if exponent == 0 && mantissa == 0 {
		return 0.0, nil
	}

	// Adjust for SAS-specific bias and scaling
	adjustedExp := int64(exponent) - (1 << (format.ExpBits - 1)) + 1023
	adjustedMant := mantissa << (52 - format.MantBits)

	result := (sign << 63) | (uint64(adjustedExp) << 52) | adjustedMant
	return math.Float64frombits(result), nil
}

func (r *Sas7bdatReader) parseTimeValue(data []byte, format string) (time.Time, error) {
	numericValue, err := r.parseNumericValue(data, numericFormats[len(data)])
	if err != nil {
		return time.Time{}, err
	}

	sasEpoch := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)

	switch format {
	case "DATE":
		// Days since January 1, 1960
		return sasEpoch.AddDate(0, 0, int(numericValue)), nil
	case "TIME":
		// Seconds since midnight
		hours := int(numericValue) / 3600
		minutes := (int(numericValue) % 3600) / 60
		seconds := int(numericValue) % 60
		return time.Date(0, 1, 1, hours, minutes, seconds, 0, time.UTC), nil
	case "DATETIME":
		// Seconds since January 1, 1960
		return sasEpoch.Add(time.Duration(numericValue) * time.Second), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported time format: %s", format)
	}
}

func (r *Sas7bdatReader) adaptToU64Format() {
	// if f.header.U64 {
	// 	// Adjust parsing parameters for u64 format
	// 	f.adjustFieldOffsets()
	// 	f.adjustSubheaderParsing()
	// }
}

// GetDataset extracts the complete dataset as a slice of maps
func (r *Sas7bdatReader) getDataset() ([]map[string]interface{}, error) {
	if len(r.columns) == 0 {
		return nil, fmt.Errorf("no columns defined")
	}

	var dataset []map[string]interface{}

	for _, page := range r.pages {
		if page.PageType == PageTypeData || page.PageType == PageTypeMix {
			rows, err := r.extractRowsFromPage(page)
			if err != nil {
				return nil, fmt.Errorf("failed to extract rows: %w", err)
			}
			dataset = append(dataset, rows...)
		}
	}

	return dataset, nil
}

func (r *Sas7bdatReader) extractRowsFromPage(page *_SAS7BDATPage) ([]map[string]interface{}, error) {
	var rows []map[string]interface{}

	// Determine data start offset (after subheaders)
	// dataOffset := f.calculateDataOffset(page)
	// rowSize := f.header.RowLength

	// // Extract individual rows
	// for offset := dataOffset; offset+int(rowSize) <= len(page.Data); offset += int(rowSize) {
	// 	row := make(map[string]interface{})
	// 	rowData := page.Data[offset : offset+int(rowSize)]

	// 	// Parse each column value
	// 	for _, col := range f.columns {
	// 		if col.Offset+int64(col.Length) <= int64(len(rowData)) {
	// 			colData := rowData[col.Offset : col.Offset+int64(col.Length)]
	// 			value, err := f.parseColumnValue(colData, col)
	// 			if err != nil {
	// 				return nil, fmt.Errorf("failed to parse column %s: %w", col.Name, err)
	// 			}
	// 			row[col.Name] = value
	// 		}
	// 	}

	// 	rows = append(rows, row)
	// }

	return rows, nil
}

func (r *Sas7bdatReader) parseColumnValue(data []byte, col *_SAS7BDATColumn) (interface{}, error) {
	switch col.Type {
	case ColumnTypeNumeric:
		format := numericFormats[col.Length]
		return r.parseNumericValue(data, format)
	// case ColumnTypeCharacter:
	// 	return f.parseStringValue(data), nil
	case ColumnTypeDate, ColumnTypeTime, ColumnTypeDateTime:
		return r.parseTimeValue(data, col.Format)
	default:
		return nil, fmt.Errorf("unsupported column type: %v", col.Type)
	}
}

// getColumnInfo returns metadata about dataset columns
func (r *Sas7bdatReader) getColumnInfo() []*_SAS7BDATColumn {
	return r.columns
}

// getFileInfo returns general file metadata
func (r *Sas7bdatReader) getFileInfo() *_SAS7BDATFileHeader {
	return r.header
}
