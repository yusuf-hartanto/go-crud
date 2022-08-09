package TaskController

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	Model "project/models"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func sqliteDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to conect database")
	}
	return db
}

func renderTemplateHTML(htmlTmp string, w http.ResponseWriter, data interface{}) {
	files := []string{
		"views/" + htmlTmp + ".html",
		"views/base.html",
	}
	tmpt, err := template.ParseFiles(files...)
	if err != nil {
		panic("Error Template: " + err.Error())
	}

	errExce := tmpt.ExecuteTemplate(w, "base", data)
	if errExce != nil {
		panic("Error Excecute : " + errExce.Error())
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) string {
	file, _, err := r.FormFile("profile_image")
	if file == nil {
		return ""
	}

	if err != nil {
		panic("Error FormFile: " + err.Error())
	}
	defer file.Close()

	tmpFile, errTmp := ioutil.TempFile("assets/images", "cover-*.png")
	if errTmp != nil {
		panic("Error TmpFile: " + err.Error())
	}
	defer tmpFile.Close()

	fileBytes, errFB := ioutil.ReadAll(file)
	if errFB != nil {
		panic("Error ReadAll: " + errFB.Error())
	}

	tmpFile.Write(fileBytes)
	var fileName = strings.Replace(tmpFile.Name(), "assets", "static", -1)
	fileName = strings.Replace(fileName, "\\", "/", -1)
	return fileName
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := sqliteDB()
	var tasks []Model.Tasks
	db.Find(&tasks)
	datas := map[string]interface{}{
		"Tasks": tasks,
	}
	renderTemplateHTML("index", w, datas)
}

func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := sqliteDB()
	if r.Method == "POST" {
		filename := uploadFile(w, r)
		task := Model.Tasks{
			Task:         r.FormValue("task"),
			Assignee:     r.FormValue("assignee"),
			Deadline:     r.FormValue("deadline"),
			Description:  r.FormValue("decription"),
			ProfileImage: filename,
		}
		db.Create(&task)

		http.Redirect(w, r, "/", http.StatusFound)
	}
	renderTemplateHTML("create", w, nil)
}

func Show(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	db := sqliteDB()
	task := Model.Tasks{}
	err := db.First(&task, params.ByName("id")).Error
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	datas := map[string]interface{}{
		"Tasks": task,
	}
	renderTemplateHTML("show", w, datas)
}

func Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	db := sqliteDB()
	task := Model.Tasks{}
	err := db.First(&task, params.ByName("id")).Error
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method == "POST" {
		filename := uploadFile(w, r)
		if filename != "" {
			task.ProfileImage = filename
		}
		task.Task = r.FormValue("task")
		task.Assignee = r.FormValue("assignee")
		task.Deadline = r.FormValue("deadline")
		task.Description = r.FormValue("decription")
		db.Save(&task)

		http.Redirect(w, r, "/", http.StatusFound)
	}
	datas := map[string]interface{}{
		"Tasks": task,
	}
	renderTemplateHTML("update", w, datas)
}

func Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	db := sqliteDB()
	task := Model.Tasks{}
	err := db.First(&task, params.ByName("id")).Error
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	db.Delete(&task, params.ByName("id"))
	http.Redirect(w, r, "/", http.StatusFound)
}
