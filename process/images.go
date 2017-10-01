package process

import (
	"net/http"
	"protoactor/data"
	"protoactor/images"
)

//type ImageServer struct {
//	Addr    string //地址
//	handler *ImageHandler //处理
//}
//
//type ImageHandler struct {
//	pattern string //路径
//}

func imageServeHTTP1(w http.ResponseWriter, r *http.Request) {
	//r.ParseMultipartForm(32 << 20) //maxMemory=32<<20,把上传的文件存储在内存和临时文件中
	if r.Method == "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var userid string = r.FormValue("userid")
	if userid == "" {
		http.Error(w, "userid empty", http.StatusBadRequest)
		return
	}
	var password string = r.FormValue("password")
	if password == "" {
		http.Error(w, "password empty", http.StatusBadRequest)
		return
	}
	//user := &data.User{Userid: userid}
	//if !user.PWDIsOK(password) { //角色验证
	//	http.Error(w, "password or userid error", http.StatusBadRequest)
	//	return
	//}
	file, _, err := r.FormFile("uploadfile") //获取文件句柄,上传参数为uploadfile
	if err != nil {
		http.Error(w, "iamge data error", http.StatusBadRequest)
		return
	}
	defer file.Close()
	imageId, err := images.SaveImage(data.Conf.ImageDir, file)
	if err != nil {
		http.Error(w, "save iamge file error", http.StatusNotFound)
		return
	}
	//user.Photo = imageId
	//user.UpdatePhoto()
	w.Write([]byte(imageId))
}

func imageServeHTTP2(w http.ResponseWriter, r *http.Request) {
	//r.ParseMultipartForm(32 << 20)
	if r.Method == "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var userid string = r.FormValue("userid")
	if userid == "" {
		http.Error(w, "userid empty", http.StatusBadRequest)
		return
	}
	var password string = r.FormValue("password")
	if password == "" {
		http.Error(w, "password empty", http.StatusBadRequest)
		return
	}
	var feedbackContent string = r.FormValue("feedback")
	if feedbackContent == "" {
		http.Error(w, "feedback content empty", http.StatusBadRequest)
		return
	}
	//user := &data.User{Userid: userid}
	//if !user.PWDIsOK(password) {
	//	http.Error(w, "password or userid error", http.StatusBadRequest)
	//	return
	//}
	file, _, err := r.FormFile("uploadfile") //上传参数为uploadfile
	if err != nil {
		http.Error(w, "iamge data error", http.StatusBadRequest)
		return
	}
	defer file.Close()
	//imageId, err := images.SaveImage(data.Conf.FeedbackDir, file)
	_, err = images.SaveImage(data.Conf.FeedbackDir, file)
	if err != nil {
		http.Error(w, "save iamge file error", http.StatusNotFound)
		return
	}
	//fb := &data.DataFeedback{
	//	Userid:    userid,
	//	Content:   feedbackContent,
	//	ImagePath: imageId,
	//}
	//if err := fb.Save(); err != nil {
	//	http.Error(w, "save content error", http.StatusNotFound)
	//	return
	//}
	//http.Error(w, "ok", http.StatusOK)
	w.Write([]byte("ok"))
}

//启动服务
func RunImages() {
	imserver1 := images.NewServer(data.Conf.ImageDir, data.Conf.ImagePort)
	imserver1.HandleFunc("/", imageServeHTTP1).Methods("POST")
	go imserver1.Run()
	imserver2 := images.NewServer(data.Conf.FeedbackDir, data.Conf.FeedbackPort)
	imserver2.HandleFunc("/", imageServeHTTP2).Methods("POST")
	go imserver2.Run()
}
