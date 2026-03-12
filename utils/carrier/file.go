package carrier

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

func Test() {
	options := make(map[string]interface{})
	options["user_id"] = 4
	options["user_name"] = "admin"
	err := WriteOptionToExcel("./upload/bomber/export/ExportData_20231114133022.xlsx", options)
	if err != nil {
		fmt.Println("err", err.Error())
	}
}

func Test2() {
	err, str := ParseExcelOption("/Users/ck/Documents/go/src/bbkdevadmin/bbm/upload/files/users-20230926131119-failed.xlsx")
	if err != nil {
		fmt.Println("err", err.Error())
		return
	}
	fmt.Println("strstrstr", str)
}

func WriteOptionToExcel(filePath string, option interface{}) error {
	gwd, _ := os.Getwd()
	filePath = strings.ReplaceAll(filePath, "./", gwd+"/")
	// 获取文件目录和文件名
	dir, filename := filepath.Split(filePath)

	// 去掉文件名的后缀
	baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))

	// 创建新的文件夹
	newDir := filepath.Join(dir, baseFilename)

	err := os.MkdirAll(newDir, 0755)
	if err != nil {
		return err
	}
	newFilePath := filepath.Join(dir, baseFilename+".zip")

	cmd := exec.Command("mv", filePath, newFilePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建文件失败 command %v %v out %v", cmd.String(), err.Error(), string(out))
	}

	cmd = exec.Command("unzip", newFilePath, "-d", newDir)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("unzip 失败 command %v %v", cmd.String(), err.Error())
	}

	// 删除压缩文件
	err = os.Remove(newFilePath)
	if err != nil {
		return fmt.Errorf("删除压缩文件失败 file %v %v", newFilePath, err.Error())
	}

	content, _ := json.Marshal(option)
	err = ioutil.WriteFile(filepath.Join(newDir, "options.txt"), content, 0644)
	if err != nil {
		return fmt.Errorf("写入文件内容失败 file %v %v", string(content[:]), err.Error())
	}

	zipCmd := exec.Command("zip", "-r", fmt.Sprintf("%s.zip", baseFilename), "./")
	zipCmd.Dir = newDir
	zipOut, zipEr := zipCmd.CombinedOutput()
	if zipEr != nil {
		return fmt.Errorf("重新打包失败 command %v %v %v", zipCmd.String(), zipEr.Error(), string(zipOut))
	}

	mvCmd := exec.Command("mv", fmt.Sprintf("%s.zip", baseFilename), fmt.Sprintf("../%s.xlsx", baseFilename))
	mvCmd.Dir = newDir
	mvOut, mvEr := mvCmd.CombinedOutput()
	if mvEr != nil {
		return fmt.Errorf("移动打包文件失败 command %v %v %v", mvCmd.String(), mvEr.Error(), string(mvOut))
	}

	if err = os.RemoveAll(newDir); err != nil {
		return fmt.Errorf("删除文件夹失败 %v %v", newDir, err.Error())
	}
	return nil
}

func ParseExcelOption(filePath string) (error, string) {
	// 获取文件目录和文件名
	dir, filename := filepath.Split(filePath)
	// 去掉文件名的后缀
	baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))

	// 创建新的文件夹
	newDir := filepath.Join(dir, baseFilename)
	defer os.RemoveAll(newDir)

	err := os.Mkdir(newDir, 0755)
	if err != nil {
		return err, ""
	}

	// 移动并重命名文件
	newFilePath := filepath.Join(dir, baseFilename+".zip")
	err = os.Rename(filePath, newFilePath)
	if err != nil {
		return err, ""
	}

	// 解压文件
	cmd := exec.Command("unzip", newFilePath, "-d", newDir)
	err = cmd.Run()
	if err != nil {
		return err, ""
	}

	content, err := os.ReadFile(fmt.Sprintf("%s/options.txt", newDir))
	if err != nil {
		return err, ""
	}

	// 移动并重命名文件
	err = os.Rename(newFilePath, filePath)
	if err != nil {
		return err, ""
	}

	return nil, string(content[:])
}

func GenerateImageWithText(text string, path string) error {
	// 创建一个大小为 400x400 像素的透明背景图像
	dc := gg.NewContext(400, 400)
	dc.SetRGB(1, 1, 1) // 设置背景色为透明

	// 设置字体颜色为灰色
	dc.SetColor(color.RGBA{169, 169, 169, 255})
	//dc.SetColor(color.Black)

	// 将文字绘制在图像上，朝向为斜向右上方
	dc.Rotate(-45)                                // 旋转文本为斜向右上方
	dc.DrawStringAnchored(text, 100, 300, 1, 0.1) // 在中心绘制文字

	// 保存生成的图像
	if err := dc.SavePNG(path); err != nil {
		return nil
	}
	return nil
}
