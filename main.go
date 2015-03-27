// server project main.go
package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"server/models"
	"server/session"
	"server/wol"
	_ "strconv"
	"strings"
	"time"
)

const (
	COOKIE_NAME = "sessionId"
	BUFFER_SIZE = 1024
)

var size string
var b bool
var m map[string]bool
var inMemorySession *session.Session
var adminusername, adminpassword string

func unescape(x string) interface{} {
	b = true
	return template.HTML(x)
}
func isAdmin(r *http.Request) bool {
	b = true
	cookie, _ := r.Cookie(COOKIE_NAME)
	if cookie != nil {
		name := inMemorySession.Get(cookie.Value)
		if name == adminusername {
			return true
		}
	}
	return false
}
func computersHandler(rnd render.Render, r *http.Request) {
	b = true
	if isAdmin(r) {
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		sqlStmt := `
			create table computers(id INTEGER PRIMARY KEY, name TEXT UNIQUE, ip TEXT UNIQUE, mac TEXT UNIQUE);
			`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println(err)
		}
		rows, err := db.Query("select Id, Name, IP, Mac from computers")
		computers := []models.Computer{}
		for rows.Next() {
			var computer = models.Computer{}
			rows.Scan(&computer.Id, &computer.Name, &computer.IP, &computer.Mac)
			computers = append(computers, computer)
		}
		rnd.HTML(200, "allcomputers", computers)
	} else {
		rnd.Redirect("/admin")
	}
}
func computersPostHandler(rnd render.Render, r *http.Request) {
	b = true
	if isAdmin(r) {
		if r.FormValue("poweron") != "" {
			a := strings.Split(r.FormValue("poweron"), ",")
			db, err := sql.Open("sqlite3", "./database.db")
			if err != nil {
				fmt.Println(err.Error())
			}
			defer db.Close()
			sqlStmt := `
			create table computers(id INTEGER PRIMARY KEY, name TEXT UNIQUE, ip TEXT UNIQUE, mac TEXT UNIQUE);
			`
			_, err = db.Exec(sqlStmt)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(a)

			for _, value := range a {
				var mac string
				db.QueryRow("select Mac from computers where Id=" + value).Scan(&mac)
				wol.MagicWake(mac)
			}

		}
	}
	rnd.Redirect("/computers")
}
func settingsHandler(rnd render.Render, r *http.Request) {
	b = true
	if isAdmin(r) {
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		sqlStmt := `
			create table computers(id INTEGER PRIMARY KEY, name TEXT UNIQUE, ip TEXT UNIQUE, mac TEXT UNIQUE);
			`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println(err)
		}
		rows, err := db.Query("select Id, Name, IP, Mac from computers")
		computers := []models.Computer{}
		for rows.Next() {
			var computer = models.Computer{}
			rows.Scan(&computer.Id, &computer.Name, &computer.IP, &computer.Mac)
			computers = append(computers, computer)
		}
		rnd.HTML(200, "settings", computers)

	} else {
		rnd.Redirect("/admin")
	}

}
func settingsPostHandler(rnd render.Render, r *http.Request) {
	b = true
	value1 := r.FormValue("Name-text")
	value2 := r.FormValue("IP-text")
	value3 := r.FormValue("Mac-text")
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	sqlStmt := `
	create table computers(id INTEGER PRIMARY KEY, name TEXT UNIQUE, ip TEXT UNIQUE, mac TEXT UNIQUE);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println(err)
	}
	rw, err := db.Exec("insert into computers(Name, IP, Mac) values('" + value1 + "', '" + value2 + "', '" + value3 + "')")
	fmt.Println(rw)
	fmt.Println(err)
	rnd.Redirect("/settings")
}
func deleteSettingsPostHandler(r *http.Request) {
	b = true
	str := r.FormValue("str")
	a := strings.Split(str, ",")
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	sqlStmt := `
	create table computers(id INTEGER PRIMARY KEY, name TEXT UNIQUE, ip TEXT UNIQUE, mac TEXT UNIQUE);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println(err)
	}
	s := "where id=" + a[0]
	for key, value := range a {
		if key != 0 {
			s = s + " or id=" + value
		}
	}
	fmt.Println(s)
	rw, err := db.Exec("delete from computers " + s)
	fmt.Println(rw)
}
func indexHandler(rnd render.Render, r *http.Request) {
	b = true
	if isAdmin(r) {
		rnd.Redirect("/active")

	} else {
		rnd.Redirect("/admin")
	}
}
func adminAuth(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	b = true
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username)
	fmt.Println(password)
	if username == adminusername && password == adminpassword {
		sessionId := inMemorySession.Init(username)
		cookie := &http.Cookie{
			Name:    COOKIE_NAME,
			Value:   sessionId,
			Expires: time.Now().Add(1000 * time.Hour),
		}
		fmt.Println("Correct")
		http.SetCookie(w, cookie)
		rnd.Redirect("/")
	} else {
		rnd.Redirect("/admin")
		fmt.Println("Incorrect")
	}
}
func adminHandler(rnd render.Render, r *http.Request) {
	b = true
	if isAdmin(r) {
		rnd.Redirect("/")
	} else {
		rnd.HTML(200, "login", nil)
	}
}
func activeHandler(rnd render.Render, r *http.Request) {
	if isAdmin(r) {
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		sqlStmt := `
			create table computers(id INTEGER PRIMARY KEY, name TEXT UNIQUE, ip TEXT UNIQUE, mac TEXT UNIQUE);
			`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println(err)
		}
		s := ""
		for key, _ := range m {
			s = s + "' or IP='" + key
		}
		computers := []models.Computer{}
		if len(s) > 4 {
			s = s[5:]
			s = "select Id, Name, IP, Mac from computers where " + s + "'"
			fmt.Println(s)
			rows, _ := db.Query(s)
			for rows.Next() {
				var computer = models.Computer{}
				rows.Scan(&computer.Id, &computer.Name, &computer.IP, &computer.Mac)
				computers = append(computers, computer)
			}
		}
		rnd.HTML(200, "active", computers)
		/*if err == nil {
			rnd.HTML(200, "active", computers)
		} else {
			rnd.HTML(200, "main", nil)
		}*/
	} else {
		rnd.Redirect("/admin")
	}
}
func activePostHandler(rnd render.Render, r *http.Request) {
	b = true
	if isAdmin(r) {
		if r.FormValue("poweroff") != "" || r.FormValue("reboot") != "" {
			var s string
			if r.FormValue("poweroff") != "" {
				s = r.FormValue("poweroff")
			} else {
				s = r.FormValue("reboot")
			}
			a := strings.Split(s, ",")
			fmt.Println(a)
			db, err := sql.Open("sqlite3", "./database.db")
			if err != nil {
				fmt.Println(err.Error())
			}
			defer db.Close()
			sqlStmt := `
			create table computers(id INTEGER PRIMARY KEY, name TEXT UNIQUE, ip TEXT UNIQUE, mac TEXT UNIQUE);
			`
			_, err = db.Exec(sqlStmt)
			if err != nil {
				fmt.Println(err)
			}
			for _, value := range a {
				var IP string
				db.QueryRow("select IP from computers where Id=" + value).Scan(&IP)

				if r.FormValue("poweroff") != "" {
					connect(IP, "poweroff")
					delete(m, IP)
				}
				if r.FormValue("reboot") != "" {
					connect(IP, "reboot ")
					delete(m, IP)
				}
			}

		}
	}
	rnd.Redirect("/active")
}
func fileUploadHandler(rnd render.Render, r *http.Request) {
	err := r.ParseMultipartForm(100000)
	if err != nil {
		fmt.Println(err.Error())
	}
	mf := r.MultipartForm
	x := r.FormValue("sendfile")
	a := strings.Split(x, ",")
	fmt.Println(a)
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	sqlStmt := `
	create table computers(id INTEGER PRIMARY KEY, name TEXT UNIQUE, ip TEXT UNIQUE, mac TEXT UNIQUE);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println(err)
	}
	files := mf.File["myfiles"]
	for i, _ := range files {
		fmt.Println(files[i].Filename)
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dst, err := os.Create("assets/files/" + files[i].Filename)

		defer dst.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if _, err := io.Copy(dst, file); err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, value := range a {
			var IP string
			db.QueryRow("select IP from computers where Id=" + value).Scan(&IP)
			connect(IP, "down assets/files/"+files[i].Filename)
		}
	}
	rnd.Redirect("/active")

}

func connect(ip string, command string) {
	conn, err := net.Dial("tcp", ip+":3003")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn.Write([]byte(command))
	conn.Close()
}

func remoteViewHandler(rnd render.Render, r *http.Request) {
	//a := strings.Split(r.FormValue("id"), ",")
	if isAdmin(r) {
		b = false
		a := r.FormValue("id")
		fmt.Println(a)
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer db.Close()
		sqlStmt := `
	create table computers(id INTEGER PRIMARY KEY, name TEXT UNIQUE, ip TEXT UNIQUE, mac TEXT UNIQUE);
	`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			fmt.Println(err)
		}

		var IP string
		db.QueryRow("select IP from computers where Id=" + a).Scan(&IP)
		connect(IP, "view")
		fmt.Println(IP)
		rnd.HTML(200, "computer", nil)
	}
}
func request_handler(conn net.Conn) {

	for {
		s := conn.RemoteAddr().String()
		i := strings.Index(s, ":")
		s = s[:i]
		//fmt.Println(s)
		line, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			return
		}
		cmd := strings.TrimRight(string(line), "\n")
		if strings.Contains(cmd, "Disconnected") {
			fmt.Println("Client disconnected: " + s)
			delete(m, s)
		} else if strings.Contains(cmd, "Connected") {
			fmt.Println("Client connected: " + s)
			m[s] = true
		} else {
			fmt.Println("undefined command " + cmd)
		}
	}
}
func launchServer() {
	psock, err := net.Listen("tcp", ":3001")
	if err != nil {
		return
	}
	for {
		conn, err := psock.Accept()
		if err != nil {
			return
		}
		go request_handler(conn)
	}
}
func screen(r *http.Request) {
	s := r.RemoteAddr
	i := strings.Index(s, ":")
	s = s[:i]
	err := r.ParseMultipartForm(100000)
	if err != nil {
		fmt.Println(err.Error())
	}
	mf := r.MultipartForm
	files := mf.File["file"]
	for i, _ := range files {
		fmt.Println(files[i].Filename)
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dst, err := os.Create("assets/files/" + files[i].Filename)

		defer dst.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if _, err := io.Copy(dst, file); err != nil {
			fmt.Println(err.Error())
			return
		}

	}
	if !b {
		connect(s, "view")
	}

}

func main() {
	b = false
	m = make(map[string]bool)
	go launchServer()
	inMemorySession = session.NewSession()
	a := os.Args[1:]
	adminusername = a[0]

	adminpassword = a[1]
	fmt.Println(adminusername)
	fmt.Println(adminpassword)
	m := martini.Classic()
	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	unescapeFuncMap := template.FuncMap{"unescape": unescape}
	m.Use(render.Renderer(render.Options{
		Directory:  "theme",                             // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".html"},                   // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON

	}))
	fmt.Println("Server is running on port 3000!")
	m.Get("/admin", adminHandler)
	m.Get("/", indexHandler)
	m.Get("/settings", settingsHandler)
	m.Get("/computers", computersHandler)
	m.Get("/active", activeHandler)
	m.Get("/rview", remoteViewHandler)
	m.Post("/upload", fileUploadHandler)
	m.Post("/active", activePostHandler)
	m.Post("/computers", computersPostHandler)
	m.Post("/settings", settingsPostHandler)
	m.Post("/delsettings", deleteSettingsPostHandler)
	m.Post("/admin", adminAuth)
	m.Post("/screen", screen)
	m.Run()
}
