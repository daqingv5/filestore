package handler

import (
	"fileStore/meta"
	"fileStore/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//向用户返回数据的ResponseWriter对象和接收用户请求的Request的对象指针
func UploadHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		//返回上传html页面
		data,err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	}else if r.Method == "POST" {
		//接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		FileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}


		newFile, err := os.Create(FileMeta.Location)
		if err!=nil{
			fmt.Printf("Failed to create file, err:%s\n, err.Error()")
			return
		}
		defer newFile.Close()

		FileMeta.FileSize,err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
			return
		}

		newFile.Seek(0,0)
		FileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(FileMeta)

		http.Redirect(w,r,"/file/upload/success",http.StatusFound)

		
	}


}

//上传已完成
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "Upload finished!")
}