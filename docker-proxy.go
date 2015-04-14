package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

var templateNginx = `
{{range $index, $element := .}}
upstream {{$index}} {
    ip_hash;
    {{range $element}}
    server {{.}}:5000;
    {{end}}
}
{{end}}


server {
    listen 80;
    server_name pcf.es;

    root /usr/share/nginx/html;

    {{range $index, $element := .}}
    location /{{$index}} {
        proxy_pass http://{{$index}};
    }
    {{end}}
}`

func getContainers() map[string][]string {
	// user := "1755739A98556879E22E"
	user := os.Getenv("USER")
	// password := "opCzCxQeb7FFmWsz2LKJ4TuhhmNULcuDUf9AgW8y"
	password := os.Getenv("PASSWORD")
	// ip := "178.62.207.180"
	ip := os.Getenv("IP")
	resp, err := http.Get("http://" + user + ":" + password + "@" + ip + ":8080/v1/containers/")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var jResponse interface{}
	err = json.Unmarshal(body, &jResponse)
	data := jResponse.(map[string]interface{})["data"].([]interface{})
	mapa := make(map[string][]string)
	for _, d := range data {
		dataContainer := d.(map[string]interface{})["data"]
		fields := dataContainer.(map[string]interface{})["fields"]
		ipAddress := fields.(map[string]interface{})["primaryIpAddress"]
		dockerInspect := dataContainer.(map[string]interface{})["dockerInspect"]
		config := dockerInspect.(map[string]interface{})["Config"]
		env := config.(map[string]interface{})["Env"]
		envVar := env.([]interface{})
		for _, e := range envVar {
			re := regexp.MustCompile("APP_NAME=(?P<name>[a-zA-Z0-9_]+)")
			if re.MatchString(e.(string)) {
				tokenizer := re.FindAllStringSubmatch(e.(string), 1)
				mapa[tokenizer[0][1]] = append(mapa[tokenizer[0][1]], ipAddress.(string))
			}
		}
	}

	return mapa
}

func templating(mapa map[string][]string) {
	tmpl, err := template.New("test").Parse(templateNginx)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("/etc/nginx/sites-enabled/proxy.conf")
	err = tmpl.Execute(f, mapa)
}

func main() {
	fmt.Println(getContainers())
	templating(getContainers())
}
