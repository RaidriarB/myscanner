package output

import (
	"fmt"
	"myscanner/core/types"
	"myscanner/settings"
	"os"
)

//几种处理扫描结果的方式
// 1. 放入文件中
// 2. 放入数据库
// TODO: 还有什么需求
// 还需要把一些统计信息打印出来。
func ProcessResult(result types.TargetPortBanners) {
	if settings.SAVE_RESULT_TO_FILE {
		saveResultToFile(result)
	}
	if settings.SAVE_RESULT_TO_DB {
		saveResultToDB(result)
	}

}

func saveResultToDB(result types.TargetPortBanners) {
	// TODO: implement
}

func saveResultToFile(result types.TargetPortBanners) {

	filename := settings.RESULT_FILE
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, r := range result {
		line := fmt.Sprintf("[%s],[%v],[%s]\n",
			r.Target.URI(), r.TcpFinger, r.Status)
		//fmt.Println(line)
		file.WriteString(line)
	}

	fmt.Printf("成功保存到文件[%s]中\n", settings.RESULT_FILE)

}
