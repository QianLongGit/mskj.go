/***
*                           _ooOoo_
*                          o8888888o
*                          88" . "88
*                          (| -_- |)
*                          O\  =  /O
*                       ____/`---'\____
*                     .'  \\|     |//  `.
*                    /  \\|||  :  |||//  \
*                   /  _||||| -:- |||||-  \
*                   |   | \\\  -  /// |   |
*                   | \_|  ''\---/''  |   |
*                   \  .-\__  `-`  ___/-. /
*                 ___`. .'  /--.--\  `. . __
*              ."" '<  `.___\_<|>_/___.'  >'"".
*             | | :  `- \`.;`\ _ /`;.`/ - ` : | |
*             \  \ `-.   \_ __\ /__ _/   .-` /  /
*        ======`-.____`-.___\_____/___.-`____.-'======
*                           `=---='
*        ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
*                    佛祖保佑       永无BUG
*           Created by wei.qin on 2019-12-11 11:56.
***/

package tools

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

type XlsxRow struct {
	Row  *xlsx.Row
	Data []string
}

func NewRow(row *xlsx.Row, data []string) *XlsxRow {
	return &XlsxRow{
		Row:  row,
		Data: data,
	}
}

func (row *XlsxRow) SetRowTitle() error {
	return generateRow(row.Row, row.Data)
}

func (row *XlsxRow) GenerateRow() error {
	return generateRow(row.Row, row.Data)
}

func generateRow(row *xlsx.Row, rowStr []string) error {
	if rowStr == nil {
		return fmt.Errorf("no data to ganerate xlsx")
	}
	for _, v := range rowStr {
		cell := row.AddCell()
		cell.SetString(v)
	}
	return nil
}
