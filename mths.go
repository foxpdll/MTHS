package main

import (
    "fmt"
    _ "net"
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
    "math/rand"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func handler_files(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile(r.URL.Path[1:])
    fmt.Fprintf(w, "%s", body)
}

func handler_cladm3(w http.ResponseWriter, r *http.Request) {
    sqlst := ""
    idst := r.FormValue("id")
    if idst=="0" {
        sqlst = fmt.Sprintf(`insert into clients values((select max(id)+1 from clients),'%s','%s','%s','%s',%s,%s,%s,'%s','%s');`,r.FormValue("name"),r.FormValue("hs"),r.FormValue("banner"),r.FormValue("url"),r.FormValue("hsver"),r.FormValue("bannerver"),r.FormValue("urlver"),r.FormValue("dtstart"),r.FormValue("dtend"))
    } else {
        sqlst = fmt.Sprintf(`update clients set name='%s',hs='%s',banner='%s',url='%s',hsver=%s,bannerver=%s,urlver=%s,dtstart='%s',dtend='%s' where id=%s;`,r.FormValue("name"),r.FormValue("hs"),r.FormValue("banner"),r.FormValue("url"),r.FormValue("hsver"),r.FormValue("bannerver"), r.FormValue("urlver"),r.FormValue("dtstart"),r.FormValue("dtend"),r.FormValue("id"))
    }
    fmt.Fprintf(w,"<html><head><meta http-equiv='refresh' content='3;/cladm/'></head><body>Zaebis!</body></html>")
    _, err := db.Exec(sqlst)
    checkErr(err)
}

func handler_cladm2(w http.ResponseWriter, r *http.Request) {
    var id int
    var name string
    var hs string
    var banner string
    var url string
    var hsver int
    var bannerver int
    var urlver int
    var dtstart string
    var dtend string
    idu:=strings.Split(r.URL.Path,"/")[2]
    if idu == "0" {
        rows, err := db.Query("select max(id)+1 from clients")
        checkErr(err)
        rows.Next()
        rows.Scan(&id)
        rows.Close()
        name="New"
        hs="files/"+strconv.Itoa(id)+"/index.html"
        banner="files/"+strconv.Itoa(id)+"/banner.jpg"
        url="http://google.ru/"
        hsver=0
        bannerver=0
        urlver=0
        dtstart="1970-01-01 00:00:00"
        dtend="2020-01-01 00:00:00"
    } else {
        rows, err := db.Query("select * from clients where id="+idu)
        checkErr(err)
        rows.Next()
        rows.Scan(&id,&name,&hs,&banner,&url,&hsver,&bannerver,&urlver,&dtstart,&dtend)
        rows.Close()
    }
    fmt.Fprintf(w, "<html><head></head><body>\n")
    fmt.Fprintf(w, "<form method='post' action='/cladm3/'>\n")
    fmt.Fprintf(w, "<table border=1>\n")
    fmt.Fprintf(w, "<tr><td>ID:</td><td><input name='id' type='text' value='"+idu+"' readonly></td></tr>\n")
    fmt.Fprintf(w, "<tr><td>Name:</td><td><input name='name' type='text' value='"+name+"'></td></tr>\n")
    fmt.Fprintf(w, "<tr><td>HotSpotStartPage:</td><td><input name='hs' type='text' value='"+hs+"'></td></tr>\n")
    fmt.Fprintf(w, "<tr><td>Banner:</td><td><input name='banner' type='text' value='"+banner+"'></td></tr>\n")
    fmt.Fprintf(w, "<tr><td>RedirectURL</td><td><input name='url' type='text' value='"+url+"'></td></tr>\n")
    fmt.Fprintf(w, "<tr><td>Ver HS</td><td><input name='hsver' type='text' value='"+strconv.Itoa(hsver)+"'></td></tr>\n")
    fmt.Fprintf(w, "<tr><td>Ver Banner</td><td><input name='bannerver' type='text' value='"+strconv.Itoa(bannerver)+"'></td></tr>\n")
    fmt.Fprintf(w, "<tr><td>Ver URL</td><td><input name='urlver' type='text' value='"+strconv.Itoa(urlver)+"'></td></tr>\n")
    fmt.Fprintf(w, "<tr><td>Start DateTime</td><td><input name='dtstart' type='text' value='"+dtstart+"'></td></tr>\n")
    fmt.Fprintf(w, "<tr><td>End DateTime</td><td><input name='dtend' type='text' value='"+dtend+"'></td></tr>\n")
    fmt.Fprintf(w, "</table>\n")
    fmt.Fprintf(w, "<input type='submit' value='Save'></form>\n")
    fmt.Fprintf(w, "</body></html>\n")
}

func handler_cladm(w http.ResponseWriter, r *http.Request) {
    var id int
    var name string
    var hs string
    var banner string
    var url string
    var hsver int
    var bannerver int
    var urlver int
    var dtstart string
    var dtend string
    fmt.Fprintf(w, "<html><head></head><body>")
    fmt.Fprintf(w, "<a href=/clstat/>Statistics</a>")
    fmt.Fprintf(w, "<table border=1>")
    rows, err := db.Query("select * from clients order by id")
    checkErr(err)
    for rows.Next() {
        rows.Scan(&id,&name,&hs,&banner,&url,&hsver,&bannerver,&urlver,&dtstart,&dtend)
        fmt.Fprintf(w, "<tr><td><a href=/cladm2/"+strconv.Itoa(id)+">"+strconv.Itoa(id)+"</a></td>")
        fmt.Fprintf(w, "<td>%s</td>", name)
        fmt.Fprintf(w, "<td>%s</td>", hs)
        fmt.Fprintf(w, "<td>%s</td>", banner)
        fmt.Fprintf(w, "<td>%s</td>", url)
        fmt.Fprintf(w, "<td>%d</td>", hsver)
        fmt.Fprintf(w, "<td>%d</td>", bannerver)
        fmt.Fprintf(w, "<td>%d</td>", urlver)
        fmt.Fprintf(w, "<td>%s</td>", dtstart)
        fmt.Fprintf(w, "<td>%s</td></tr>", dtend)
    }
    rows.Close()
    fmt.Fprintf(w, "</table>")
    fmt.Fprintf(w, "<a href=/cladm2/0>Add new</a>")
    fmt.Fprintf(w, "</body></html>")
}

func handler_clstat(w http.ResponseWriter, r *http.Request) {
    var name string
    var types string
    var sum int
    fmt.Fprintf(w, "<html><head></head><body>")
    rows, err := db.Query(`select (select name from clients where id=clientid) as name, replace(replace(replace(type,1,"HS"),2,"Banner"),3,"URL") as type, sum(1) as sum from logMedia as lm group by 1,2 order by 1,2;`)
    checkErr(err)
    fmt.Fprintf(w, "<table border=1>")
    for rows.Next() {
        rows.Scan(&name,&types,&sum)
        fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%d</td></tr>",name,types,sum)
    }
    rows.Close()
    fmt.Fprintf(w, "</table>")
    fmt.Fprintf(w, "</body></html>")
}

func handler_hscheckcode(w http.ResponseWriter, r *http.Request) {
//    fmt.Fprintf(w,r.FormValue("identity"))
//    fmt.Fprintf(w,r.FormValue("ip"))
//    fmt.Fprintf(w,r.FormValue("mac"))
//    fmt.Fprintf(w,r.FormValue("phone"))
//    fmt.Fprintf(w,r.FormValue("code"))
    count:=0
    sql := fmt.Sprintf("select count(1) from users where mac='%s' and code='%s' and checked=0;",r.FormValue("mac"),r.FormValue("code"))
//    fmt.Fprintf(w,sql)
    rows, err := db.Query(sql)
    checkErr(err)
    rows.Next()
    rows.Scan(&count)
    rows.Close()
    if count > 0 {
        sql := fmt.Sprintf("update users set checked=1 where mac='%s' and code='%s';",r.FormValue("mac"),r.FormValue("code"))
        _, err := db.Exec(sql)
        checkErr(err)
    }
    fmt.Fprintf(w,`<html><head><meta http-equiv="refresh" content="0;http://google.ru/"></head></html>`)
}

func handler_hssendcode(w http.ResponseWriter, r *http.Request) {
    phone:=r.FormValue("phone")
    phone=strings.Replace(phone," ","",100)
    phone=strings.Replace(phone,"(","",100)
    phone=strings.Replace(phone,")","",100)
    if phone != "" {
        code := strconv.Itoa(rand.Intn(9999))
        sql := fmt.Sprintf("insert into users values('%s','%s','%s','%s','%s',datetime('now'),0);",r.FormValue("identity"),r.FormValue("mac"),phone,code,r.FormValue("ip"))
        fmt.Printf("Send sms to ")
        fmt.Printf(phone)
        fmt.Printf(" ")
        fmt.Printf(code)
        fmt.Printf("\n")
        http.Get("http://somesite/foxsms.php?tel="+phone+"&txt="+code)
//      fmt.Fprintf(w,sql)
        _, err := db.Exec(sql)
        checkErr(err)
        body := getclhs()
        stzam:=fmt.Sprintf(`
                            <form action="/hscheckcode/" id="start" name ="redirect" accept-charset="UTF-8" method="post" autocomplete="off">
                                <input type="hidden" name="identity" value="%s">
                                <input type="hidden" name="ip" value="%s">
                                <input type="hidden" name="mac" value="%s">
                                <input type="hidden" name="phone" class="phone" id="phone" placeholder="пример +7 (xxx) xxx-xx-xx" value="%s"/>
                                <input name="code" id="code"/>
                                <input class="bt" class="send" type="submit"  value="Войти в Internet">
                            </form>`,r.FormValue("identity"),r.FormValue("ip"),r.FormValue("mac"),phone)
        fmt.Fprintf(w, strings.Replace(string(body),"###METAFORM###",stzam,1))
    }
}
//identity text,
//mac text,
//tel text,
//code text,
//ip text,
//regdate text


func handler_hs(w http.ResponseWriter, r *http.Request) {
    var count int
    rows, err := db.Query("select count(1) from users where mac='"+r.FormValue("mac")+"' and mac<>'' and checked=1")
    checkErr(err)
    rows.Next()
    rows.Scan(&count)
    rows.Close()
    body := getclhs()
    if count>0 {
        stzam:=fmt.Sprintf(`
                            <form action="%s" id="start" name ="redirect" accept-charset="UTF-8" method="post" autocomplete="off">
                                <input type="hidden" name="dst" value="%s">
                                <input type="hidden" name="username" value="T-%s">
                                <input type="hidden" name="phone" class="phone" id="phone" placeholder="пример +7 (xxx) xxx-xx-xx" value="+7(999)1111111"/>
                                <input class="bt" class="send" type="submit"  value="Войти в интернет">
                            </form>`,r.FormValue("link-login-only"),getclurl(),r.FormValue("mac"))
        fmt.Fprintf(w, strings.Replace(string(body),"###METAFORM###",stzam,1))
    } else {
        stzam:=fmt.Sprintf(`
                            <form action="/hssendcode/" id="start" name ="redirect" accept-charset="UTF-8" method="post" autocomplete="off">
                                <input type="hidden" name="identity" value="%s">
                                <input type="hidden" name="ip" value="%s">
                                <input type="hidden" name="mac" value="%s">
                                <input name="phone" class="phone" id="phone" placeholder="пример +7 (xxx) xxx-xx-xx"/>
                                <input class="bt" class="send" type="submit"  value="Прислать SMS">
                            </form>`,r.FormValue("identity"),r.FormValue("ip"),r.FormValue("mac"))
        fmt.Fprintf(w, strings.Replace(string(body),"###METAFORM###",stzam,1))
    }
    _, err = db.Exec("insert into log values(datetime('now'),'"+r.FormValue("identity")+"','"+r.FormValue("mac")+"','"+r.FormValue("ip")+"','"+r.RemoteAddr+"','"+r.Header.Get("User-Agent")+"');")
    checkErr(err)
}

func handler_banner(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("select id,banner from shb")
    checkErr(err)
    var st string
    var id int
    rows.Next()
    rows.Scan(&id,&st)
    rows.Close()
    body, _ := ioutil.ReadFile(st)
    fmt.Fprintf(w, "%s", body)
    addLogMedia(id,2)
}

func handler_main(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<html>Pusto!</html>\n")
}

func handler_favicon(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "\n")
}

func getclhs ()(string){
    var id int
    var st string
    rows, err := db.Query("select id,hs from shs")
    checkErr(err)
    rows.Next()
    rows.Scan(&id,&st)
    rows.Close()
    addLogMedia(id,1)
    body, _ := ioutil.ReadFile(st)
    return string(body)
}

func getclurl ()(string){
    var id int
    var st string
    rows, err := db.Query("select id,url from shu")
    checkErr(err)
    rows.Next()
    rows.Scan(&id,&st)
    rows.Close()
    addLogMedia(id,3)
    return st
}

func addLogMedia (id int,t int){
    fmt.Printf("Logging Media "+strconv.Itoa(id)+","+strconv.Itoa(t)+"\n")
    _, err := db.Exec("insert into logMedia values("+strconv.Itoa(id)+",datetime('now'),"+strconv.Itoa(t)+");")
    checkErr(err)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

var (
    db *sql.DB
)

func main() {
    var err error
    db, err = sql.Open("sqlite3", "./mths.db")
    checkErr(err)
    _, err = db.Exec(`create table if not exists clients (
id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
name text,
hs text,
banner text,
url text,
hsver int DEFAULT 0,
bannerver int DEFAULT 0,
urlver int DEFAULT 0,
dtstart text DEFAULT "1970-01-01 00:00",
dtend text DEFAULT "2020-01-01 00:00"
);`)
    checkErr(err)
    _, err = db.Exec(`create table if not exists users (
identity text,
mac text,
tel text,
code text,
ip text,
regdate text,
checked int
);`)
    checkErr(err)
    _, err = db.Exec("create table if not exists log (dt text, identity text, mac text, ip1 text, ip2 text, useragent text);")
    checkErr(err)
    _, err = db.Exec("create table if not exists logMedia (clientid int, dt text, type int);")
    checkErr(err)
    _, err = db.Exec(`create view if not exists shs as SELECT id,hs FROM clients where datetime("now")>dtstart and datetime("now")<dtend ORDER BY RANDOM()%hsver desc LIMIT 1;`)
    checkErr(err)
    _, err = db.Exec(`create view if not exists shb as SELECT id,banner FROM clients where datetime("now")>dtstart and datetime("now")<dtend ORDER BY RANDOM()%bannerver desc LIMIT 1;`)
    checkErr(err)
    _, err = db.Exec(`create view if not exists shu as SELECT id,url FROM clients where datetime("now")>dtstart and datetime("now")<dtend ORDER BY RANDOM()%urlver desc LIMIT 1;`)
    checkErr(err)
    http.HandleFunc("/hs/", handler_hs)
    http.HandleFunc("/hssendcode/", handler_hssendcode)
    http.HandleFunc("/hscheckcode/", handler_hscheckcode)
    http.HandleFunc("/files/", handler_files)
    http.HandleFunc("/banner/", handler_banner)
    http.HandleFunc("/cladm/", handler_cladm)
    http.HandleFunc("/cladm2/", handler_cladm2)
    http.HandleFunc("/cladm3/", handler_cladm3)
    http.HandleFunc("/clstat/", handler_clstat)
    http.HandleFunc("/favicon.ico", handler_favicon)
    http.HandleFunc("/", handler_main)
    fmt.Printf("Starting Service\n")
    http.ListenAndServe(":8181", nil)
    fmt.Printf("Stoping Service\n")
}
