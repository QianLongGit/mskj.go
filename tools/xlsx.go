package tools

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/tealeg/xlsx"
	"path"
)

// xlsx(excel)表格导入/导出工具

// 读取文档, 指定开始行开始列
func Read(excelFileName string, startRow int) ([]interface{}, error) {
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		logs.Error(fmt.Sprintf("读取文档失败, 文件路径: %s", excelFileName))
		return nil, fmt.Errorf("读取文档失败, 文件路径: %s", excelFileName)
	}
	var sheets []interface{}
	for _, sheet := range xlFile.Sheets {
		var res [2]interface{}
		var data [][]string
		for rowIndex, row := range sheet.Rows {
			if rowIndex < startRow {
				continue
			}
			var rowData []string
			for _, cell := range row.Cells {
				text := cell.String()
				rowData = append(rowData, text)
			}
			data = append(data, rowData)
		}
		res[0] = sheet.Name
		res[1] = data
		sheets = append(sheets, res)
	}
	return sheets, nil
}

// 重写sheet
// excelSrc源文件路径
// excelDest目标文件路径
// sheetData[0] {string} sheet名称
// sheetData[1] {[][]string} 需要添加的数据
func ReWriteSheet(excelSrc string, excelDest string, sheetData [2]interface{}) error {
	xlFile, err := xlsx.OpenFile(excelSrc)
	if err != nil {
		logs.Error(fmt.Sprintf("读取文档失败, 文件路径: %s", excelSrc))
		return fmt.Errorf("读取文档失败, 文件路径: %s", excelSrc)
	}
	l := len(sheetData)
	if l != 2 {
		return fmt.Errorf("数据长度不符合要求, 重写sheet失败")
	}
	sheetName, ok := sheetData[0].(string)
	if !ok {
		return fmt.Errorf("sheet名称类型必须为string")
	}
	data, ok := sheetData[1].([][]string)
	if !ok {
		return fmt.Errorf("sheet数据类型必须为[][]string")
	}
	for _, sheet := range xlFile.Sheets {
		if sheet.Name == sheetName {
			// 添加数据到sheet
			for _, rowData := range data {
				row := sheet.AddRow()
				for _, cellData := range rowData {
					cell := row.AddCell()
					cell.Value = cellData
				}
			}
			break
		}
	}
	// 创建目标目录
	_, err = CreatePathIfNotExists(path.Dir(excelDest))
	if err != nil {
		return err
	}
	// 文件保存
	xlFile.Save(excelDest)
	return nil
}
