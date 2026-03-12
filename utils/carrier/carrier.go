package carrier

import (
	"fmt"
	"sort"

	"github.com/spf13/cast"

	"github.com/xuri/excelize/v2"
)

type ExcelCarrier struct {
	Titles        []string
	File          *excelize.File
	FilePath      string
	SheetName     string
	NewSheet      bool
	Data          [][]string
	TitleRichText bool
	RowNumber     int
	IsLock        bool
}

var ExcelStyleMap = map[string]*excelize.Style{
	"detail_header": {
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E8E8E8"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "微软雅黑",
			Size:   11,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
		},
	},
	"detail_data": {
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
		},
	},
	"blue_header": {
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00B0F0"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold:   false,
			Family: "黑体-简",
			Size:   12,
			Color:  "#ffffff",
		},

		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "a6a6a6",
				Style: 1,
			},
		},
	},
	"blue_data": {
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   false,
			Family: "黑体-简",
			Size:   12,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "a6a6a6",
				Style: 1,
			},
		},
	},
	"blue_center_data": {
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   false,
			Family: "黑体-简",
			Size:   12,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "a6a6a6",
				Style: 1,
			},
		},
	},
	"blue_num_data": {
		NumFmt: 2,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   false,
			Family: "黑体-简",
			Size:   12,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "a6a6a6",
				Style: 1,
			},
		},
	},
	"blue_right_data": {
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   false,
			Family: "黑体-简",
			Size:   12,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "a6a6a6",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "a6a6a6",
				Style: 1,
			},
		},
	},
	"detail_data_red": {
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FF2D00"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
		},
	},
	"detail_data_orange": {
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFB900"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
		},
	},
	"detail_header_blue": {
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Family: "微软雅黑",
			Size:   11,
			Color:  "#000000",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00B9FF"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
		},
	},
}

func (ec *ExcelCarrier) Write() (err error) {
	defer func() {
		ec.Data = make([][]string, 0)
	}()
	if ec.FilePath != "" {
		return ec.Append()
	}
	if ec.NewSheet {
		ec.File.NewSheet(ec.SheetName)
	} else {
		ec.SheetName = ec.File.GetSheetName(0)
	}

	//err = ec.File.SetSheetFormatPr(ec.SheetName, excelize.DefaultRowHeight(24), excelize.DefaultColWidth(12), excelize.CustomHeight(true))
	//if err != nil {
	//	return err
	//}
	hs, err := ec.File.NewStyle(ExcelStyleMap["blue_header"])
	if err != nil {
		return err
	}
	ds, err := ec.File.NewStyle(ExcelStyleMap["blue_data"])
	if err != nil {
		return err
	}
	streamWriter, err := ec.File.NewStreamWriter(ec.SheetName)
	if err != nil {
		return err
	}

	hl := make([]interface{}, 0)
	if len(ec.Titles) == 0 {
		return nil
	}
	colWidthList := make([][]int, len(ec.Titles))

	for i, v := range ec.Titles {
		hl = append(hl, excelize.Cell{StyleID: hs, Value: v})
		colWidthList[i] = append(colWidthList[i], len(v))
	}
	for _, d := range ec.Data {
		for i, v := range d {
			colWidthList[i] = append(colWidthList[i], len(v))
		}
	}

	if len(colWidthList) > 0 {
		colWidth := getColWidth(colWidthList)
		for i, w := range colWidth {
			err = streamWriter.SetColWidth(i+1, i+1, w)
			if err != nil {
				return err
			}
		}
	}

	err = streamWriter.SetRow("A1", hl, excelize.RowOpts{Height: 24})
	if err != nil {
		return err
	}
	for i, d := range ec.Data {
		row := make([]interface{}, 0)
		for i, v := range d {
			row = append(row, excelize.Cell{StyleID: ds, Value: v})
			colWidthList[i] = append(colWidthList[i], len(v))
		}
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := streamWriter.SetRow(cell, row, excelize.RowOpts{Height: 24}); err != nil {
			continue
		}
	}
	ec.RowNumber += len(ec.Data)

	err = streamWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (ec *ExcelCarrier) Append() (err error) {
	if len(ec.Data) == 0 {
		return nil
	}
	ec.File, err = excelize.OpenFile(ec.FilePath)
	ds, err := ec.File.NewStyle(ExcelStyleMap["blue_data"])
	if err != nil {
		return err
	}
	err = ec.File.SetCellStyle(ec.SheetName, fmt.Sprintf("A%d", ec.RowNumber+2), fmt.Sprintf("%s%d", IntToExcelColumn(len(ec.Data[0])+1), ec.RowNumber+len(ec.Data)+1), ds)
	if err != nil {
		return err
	}
	for i, d := range ec.Data {
		for j, v := range d {
			if err = ec.File.SetCellValue(ec.SheetName, fmt.Sprintf("%s%d", IntToExcelColumn(j+1), ec.RowNumber+i+2), v); err != nil {
				continue
			}
		}
	}
	ec.RowNumber += len(ec.Data)
	return ec.File.SaveAs(ec.FilePath)
}

func IntToExcelColumn(n int) string {
	result := ""
	for n > 0 {
		n--
		result = string('A'+n%26) + result
		n /= 26
	}
	return result
}

func getColWidth(colWidthList [][]int) (colWidth []float64) {
	for _, v := range colWidthList {
		if len(v) > 0 {
			sort.Ints(v)
			width := cast.ToFloat64(v[len(v)-1])
			if width < 10 {
				width = 10
			}
			if width > 255 {
				width = 255
			}
			colWidth = append(colWidth, width)
		} else {
			colWidth = append(colWidth, 10)
		}

	}
	return
}
