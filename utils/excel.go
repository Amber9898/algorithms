package utils

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
)

func WriteSortExcel(strList []string){
	f:= excelize.NewFile()
	//表头第一行
	categories := map[string]string{"A1": "条目"}
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	line := 2
	for _, str := range strList{
		values := map[string]interface{}{
			fmt.Sprintf("A%d", line): str,
		}
		for k, v := range values {
			f.SetCellValue("Sheet1", k, v)
		}
		line++
	}

	// 根据指定路径保存文件
	err := f.SaveAs("./t_LeaguePKNumber.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}
