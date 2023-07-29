package main

import (
	"batch48/connection"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Project struct {
	Id        int
	Name      string
	StarDate  string
	EndDate   string
	Duration  string
	Detail    string
	Playstore bool
	Android   bool
	Java      bool
	React     bool
}
type Users struct {
	Id             int
	Name           string
	Email          string
	HashedPassword string
}
type UserLoginSession struct {
	IsLogin bool
	Name    string
}

var userLoginSession = UserLoginSession{}

var dataProject = []Project{}

func main() {
	e := echo.New()
	connection.DataBaseConnect()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("PDI IS NUMBER ONE"))))

	// e = echo package nyaa...
	// get = method yg akan djalankan...
	// npoin nya / routing...

	// untuk mengirim folder routing statis
	e.Static("/public", "public")
	// rout home
	e.GET("/hello", helloWolrd)
	e.GET("/index", home)

	// rout project
	e.GET("/addproject-data", project)
	e.GET("/formaddproject", formproject)
	e.GET("/project-detail/:id", projectDetail)
	e.GET("/edit-addproject/:id", editProject)
	e.POST("/edit-addproject/:id", submitEditedProject)
	e.POST("/delete-addproject/:id", deleteProject)
	e.POST("/addproject", submitProject)

	// rout contact
	e.GET("/contact", contact)

	// rout testimonial
	e.GET("/testimonial", testimonial)

	// rout login
	e.GET("/form-login", formLogin)
	e.POST("/form-login", Login)

	// rout logout
	e.POST("/logout", Logout)

	// rout register
	e.GET("/form-register", formRegister)
	e.POST("/form-register", Register)

	// rout port my
	e.Logger.Fatal(e.Start(":1234"))

}
func helloWolrd(c echo.Context) error {
	return c.String(http.StatusOK, "hello, ibab")
}

var userData = Users{}

func home(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/index.html")

	sess, _ := session.Get("session", c)

	datas := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
		"DataSession":  userData,
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"massage": err.Error()})
	}

	return tmpl.Execute(c.Response(), datas)
}
func project(c echo.Context) error {

	var tmpl, err = template.ParseFiles("views/addproject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, duration, detail, playstore, android, java, react FROM tb_project")

	var result []Project
	for data.Next() {
		var each = Project{}
		err := data.Scan(&each.Id, &each.Name, &each.StarDate, &each.EndDate, &each.Duration, &each.Detail, &each.Playstore, &each.Android, &each.Java, &each.React)

		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		result = append(result, each)

	}
	

	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userLoginSession.IsLogin = false
	} else {
		userLoginSession.IsLogin = true
		userLoginSession.Name = sess.Values["name"].(string)
	}



	Projects := map[string]interface{}{
		"Projects": result,
		"UserLoginSession": userLoginSession,
	}

	return tmpl.Execute(c.Response(), Projects)
}

func formproject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/formaddproject.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

// func untuk rout post nyaaa
func submitProject(c echo.Context) error {
	name := c.FormValue("input-project-title")
	detail := c.FormValue("input-description")
	starDate := c.FormValue("input-startdate")
	endDate := c.FormValue("input-enddate")
	diffUse := countDuration(starDate, endDate)

	// checkbox
	var nodejs bool
	if c.FormValue("nodejs") == "checked" {
		nodejs = true
	}

	var reactjs bool
	if c.FormValue("reactjs") == "checked" {
		reactjs = true
	}

	var nextjs bool
	if c.FormValue("nextjs") == "checked" {
		nextjs = true
	}

	var typescript bool
	if c.FormValue("typescript") == "checked" {
		typescript = true
	}

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (name, detail, start_date, end_date, duration, playstore, android, java, react) VALUES ($1 , $2, $3, $4, $5, $6, $7, $8, $9)", name, detail, starDate, endDate, diffUse, nodejs, reactjs, nextjs, typescript)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	fmt.Println(dataProject)

	return c.Redirect(http.StatusMovedPermanently, "/addproject-data")
}

// detail project
func projectDetail(c echo.Context) error {
	// strconcov/string converter = untuk conver tipe data lain jadi string
	id := c.Param("id")
	var tmpl, errTmp = template.ParseFiles("views/project-detail.html")

	idToInt, _ := strconv.Atoi(id)

	var ProjectDetail = Project{}

	// query row untuk nge get satu datanya
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, duration, detail, playstore, android, java, react FROM tb_project WHERE id=$1", idToInt).Scan(&ProjectDetail.Id, &ProjectDetail.Name, &ProjectDetail.StarDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Detail, &ProjectDetail.Playstore, &ProjectDetail.Android, &ProjectDetail.Java, &ProjectDetail.React)

	fmt.Println("ini data blog detail", err)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	data := map[string]interface{}{
		"Id":      id,
		"Project": ProjectDetail,
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": errTmp.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}
func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_project WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Name, &ProjectDetail.StarDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Detail, &ProjectDetail.Playstore, &ProjectDetail.Android, &ProjectDetail.Java, &ProjectDetail.React)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message query": err.Error()})
	}
	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, errTmp = template.ParseFiles("views/editproject.html")
	if errTmp != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message rout": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)

}
func submitEditedProject(c echo.Context) error {

	// Menangkap Id dari Query Params
	id, _ := strconv.Atoi(c.Param("id"))

	name := c.FormValue("input-project-title")
	starDate := c.FormValue("input-startdate")
	endDate := c.FormValue("input-enddate")
	duration := countDuration(starDate, endDate)

	detail := c.FormValue("input-description")
	// checkbox
	var playstore bool
	if c.FormValue("nodejs") == "checked" {
		playstore = true
	}

	var android bool
	if c.FormValue("reactjs") == "checked" {
		android = true
	}

	var java bool
	if c.FormValue("nextjs") == "checked" {
		java = true
	}

	var react bool
	if c.FormValue("typescript") == "checked" {
		react = true
	}

	// dataProject[id] = editedProject
	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_project SET name=$1, start_date=$2, end_date=$3, duration=$4, detail=$5, playstore=$6, android=$7, java=$8, react=$9 WHERE id=$10", name, starDate, endDate, duration, detail, playstore, android, java, react, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/addproject-data")
}
func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	// dataProject = append(dataProject[:id], dataProject[id+1:]...)

	// query delete to data base
	connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
	return c.Redirect(http.StatusMovedPermanently, "/addproject-data")
}
func countDuration(d1 string, d2 string) string {
	date1, _ := time.Parse("2006-01-02", d1)
	date2, _ := time.Parse("2006-01-02", d2)

	diff := date2.Sub(date1)
	days := int(diff.Hours() / 24)
	weeks := days / 7
	months := days / 30
	years := months / 12

	if days < 7 {
		return strconv.Itoa(days) + " Hari"
	}
	if weeks < 4 {
		return strconv.Itoa(weeks) + " Minggu"
	}
	if months < 12 {
		return strconv.Itoa(months) + " Bulan"
	}
	return strconv.Itoa(years) + " Tahun"

}

func contact(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}
func testimonial(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/testimonial.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), nil)
}
func formLogin(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/formlogin.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	sesi, errsesi := session.Get("session", c)

	if errsesi != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	fmt.Println("message:", sesi.Values["message"])
	fmt.Print("status:", sesi.Values["status"])

	flash := map[string]interface{}{
		"FlashMessage": sesi.Values["message"],
		"FlashStatus":  sesi.Values["status"],
	}

	delete(sesi.Values, "message")
	delete(sesi.Values, "status")
	sesi.Save(c.Request(), c.Response())

	return tmpl.Execute(c.Response(), flash)
}
func Login(c echo.Context) error {
	inputEmail := c.FormValue("inputEmail")
	inputPassword := c.FormValue("inputPassword")

	fmt.Printf(inputPassword)

	// struct
	user := Users{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, email, password FROM tb_user WHERE email=$1", inputEmail).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword)

	if err != nil {
		return redirectWithMessage(c, "Login Gagal", false, "/form-login")
	}

	errPassword := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(inputPassword))

	if errPassword != nil {
		return redirectWithMessage(c, "Login Gagal", false, "/form-login")
	}

	sesi, _ := session.Get("session", c)
	sesi.Options.MaxAge = 10800   //3jam
	sesi.Values["message"] = "Login Succes!!!"
	sesi.Values["status"] = true
	sesi.Values["name"] =user.Name
	sesi.Values["email"] = user.Email
	sesi.Values["id"] = user.Id
	sesi.Values["isLogin"] = true
	sesi.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/index")
	

	// return c.Redirect(http.StatusMovedPermanently, "/index")
}
func formRegister(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/formregister.html")

	sesi, errsesi := session.Get("session", c)

	if errsesi != nil {
		return c.JSON(http.StatusInternalServerError, errsesi.Error())
	}

	flash := map[string]interface{}{
		"FlashMessage": sesi.Values["message"],
		"FlashStatus":  sesi.Values["status"],
	}
	delete(sesi.Values, "message")
	delete(sesi.Values, "status")
	sesi.Save(c.Request(), c.Response())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return tmpl.Execute(c.Response(), flash)
}
func Register(c echo.Context) error {
	inputName := c.FormValue("inputName")
	inputEmail := c.FormValue("inputEmail")
	inputPassword := c.FormValue("inputPassword")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), 10)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	fmt.Println(inputName, inputEmail, inputPassword)

	regis, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (name, email, password) VALUES($1, $2, $3)", inputName, inputEmail, hashedPassword)
	fmt.Println("masuk kesini mas", regis.RowsAffected())

	if err != nil {
		return redirectWithMessage(c, "Registrasi Gagal", false, "/form-register")
	}

	return redirectWithMessage(c, "Registtrasi Berhasil", true, "/form-login")
}
func Logout(c echo.Context) error {
	sesi, _ := session.Get("session", c)
	sesi.Options.MaxAge = -1
	sesi.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/index")
}
func redirectWithMessage(c echo.Context, message string, status bool, patch string) error {
	sesi, errsesi := session.Get("session", c)

	if errsesi != nil {
		return c.JSON(http.StatusInternalServerError, errsesi.Error())
	}
	sesi.Values["message"] = message
	sesi.Values["status"] = status
	sesi.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, patch)
}
