package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type link struct {
	port int
	name string
	uri  string
}

func getEnv() {
	envs := os.Environ()
	for _, env := range envs {
		// fmt.Println(env)

		pair := strings.Split(env, "=")
		l := redex(pair[0])
		if l != nil {
			l.uri = pair[1]
			fmt.Println(l)
		}
		// fmt.Println(pair[0])
	}

}

func getContainers() {
	user := "1755739A98556879E22E"
	password := "opCzCxQeb7FFmWsz2LKJ4TuhhmNULcuDUf9AgW8y"
	resp, err := http.Get("http://" + user + ":" + password + "@178.62.207.180:8080/v1/containers/")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
	var jResponse interface{}
	err = json.Unmarshal(body, &jResponse)
	data := jResponse.(map[string]interface{})["data"].([]interface{})
	mapa := make(map[string][]string)
	for _, d := range data {

		dataContainer := d.(map[string]interface{})["data"]
		fields := dataContainer.(map[string]interface{})["fields"]
		ipAddress := fields.(map[string]interface{})["primaryIpAddress"]
		fmt.Println(ipAddress)

		dockerInspect := dataContainer.(map[string]interface{})["dockerInspect"]
		config := dockerInspect.(map[string]interface{})["Config"]
		env := config.(map[string]interface{})["Env"]
		envVar := env.([]interface{})

		for _, e := range envVar {
			re := regexp.MustCompile("APP_NAME=(?P<name>[a-zA-Z0-9_]+)")
			if re.MatchString(e.(string)) {
				tokenizer := re.FindAllStringSubmatch(e.(string), 1)
				mapa[tokenizer[0][1]] = append(mapa[tokenizer[0][1]], ipAddress.(string))
				fmt.Println(mapa)

			}
		}
		// re := regexp.MustCompile("APP_NAME=(?P<name>[a-zA-Z0-9_]+)")
		// str, _ := env.(string)
		// tokenizer := re.FindAllStringSubmatch(str, 1)

		// fmt.Println(tokenizer)
		// fmt.Println(str)
		// fmt.Println(env.([]interface{}).([]string))

	}
}

func redex(env string) *link {
	re := regexp.MustCompile("^(?P<service_name>[a-zA-Z_]+)(_[\\d]+)?_PORT_(?P<service_port>[\\d]+)_TCP_ADDR$")
	tokenizer := re.FindAllStringSubmatch(env, 1)
	if re.MatchString(env) {
		fmt.Println(tokenizer[0][1])
		fmt.Println(tokenizer[0][2])
		fmt.Println(tokenizer[0][3])
		portInt, _ := strconv.Atoi(tokenizer[0][3])
		return &link{
			port: portInt,
			name: tokenizer[0][1],
		}
	}
	return nil
}

func main() {
	getContainers()
}
