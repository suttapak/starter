package service

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/suttapak/starter/errs"
	"github.com/xuri/excelize/v2"
)

type ExcelStyle uint

const (
	CurrencyStyle ExcelStyle = iota + 1
)

type (
	ExcelDataFrame struct {
		Header []string
		Rows   [][]any
	}

	RowStyle struct {
		Sheet, Col string
		Style      ExcelStyle
	}

	Excel interface {
		Build(df *ExcelDataFrame, st ...RowStyle) (*bytes.Buffer, error)
	}
	excel struct{}
)

// Build สร้างไฟล์ Excel จากข้อมูลที่ได้รับในรูปแบบ ExcelDataFrame โดยจะกำหนดชื่อคอลัมน์ในแถวแรก
// และเติมข้อมูลแต่ละแถวลงในไฟล์ Excel จากนั้นปรับขนาดคอลัมน์ให้เหมาะสม และคืนค่าเป็น bytes.Buffer
// หากเกิดข้อผิดพลาดระหว่างการสร้างหรือเขียนข้อมูล จะคืนค่า error กลับมา
func (e excel) Build(df *ExcelDataFrame, st ...RowStyle) (*bytes.Buffer, error) {

	if df == nil {
		return nil, errs.ErrBadRequest
	}
	f := excelize.NewFile()

	currencyStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 4, // built-in Excel number format for currency with 2 decimals
	})
	if err != nil {
		return nil, err
	}

	const sheet = "Sheet1"

	for i, header := range df.Header {
		cell := fmt.Sprintf("%s1", string(rune(65+i))) // Column A, B, C, etc.
		if err := f.SetCellValue(sheet, cell, header); err != nil {
			return nil, err
		}
	}

	for i, row := range df.Rows {
		for j, c := range row {
			cell := fmt.Sprintf("%s", string(rune(65+j))) // Column A, B, C, etc.
			if err := f.SetCellValue(sheet, fmt.Sprintf("%s%d", cell, i+2), c); err != nil {
				return nil, err
			}
		}
	}
	if err := e.setColumnFit(sheet, f); err != nil {
		return nil, err
	}

	for _, s := range st {
		if s.Style == CurrencyStyle {
			if err := f.SetColStyle(s.Sheet, s.Col, currencyStyle); err != nil {
				return nil, err
			}
		}
	}

	return f.WriteToBuffer()

}

func (e excel) setColumnFit(sheet string, x *excelize.File) error {
	cols, err := x.GetCols(sheet)
	if err != nil {
		return err
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 2 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return err
		}
		if err := x.SetColWidth(sheet, name, name, float64(largestWidth)); err != nil {
			return err
		}
	}

	rows, err := x.GetRows(sheet)
	if err != nil {
		return err
	}
	for idx := range rows {
		if err := x.SetRowHeight(sheet, idx+1, 18); err != nil {
			return err
		}
	}
	return nil
}

func NewExcelService() Excel {
	return excel{}
}
