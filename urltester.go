package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type testDef struct {
	URL            string `json:"URL"`
	TestType       string `json:"testType"`
	ExpectedResult string `json:"expectedResult"`
}

type testsURL struct {
	Comment            string    `json:"_comment"`
	AppName            string    `json:"appName"`
	NonRegressionTests []testDef `json:"nonRegressionTests"`
	CheckTests         []testDef `json:"checkTests"`
	ResponseCodeTests  []string  `json:"responseCodeTests"`
}

func main() {
	var testsConf testsURL
	var nr, c, r bool
	var baseURL string
	var failCount int

	if len(os.Args) < 4 {
		fmt.Println("At least 3 arguments required: <path_to_json> <base_url> <test_type>")
		fmt.Println("test_type could be either all, nonregression|nr, check|c responsecode|r")
		os.Exit(1)
	}

	baseURL = os.Args[2]

	fmt.Println(fmt.Sprintf("Parsing configuration json %q", os.Args[1]))
	jsonFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()
	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(b, &testsConf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(fmt.Sprintf("Setting up test's type %v", os.Args[3:]))
	for _, v := range os.Args[3:] {
		switch v {
		case "all":
			nr, c, r = true, true, true
		case "nonregression", "nr":
			nr = true
		case "check", "c":
			c = true
		case "responsecode", "r":
			r = true
		default:
			log.Fatalln("Test type(s) requested paramater is not correct (test_type could be either all, nonregression|nr, check|c responsecode|r)")
		}
	}
	if nr {
		fmt.Println("Non regression tests:")
		for _, v := range testsConf.NonRegressionTests {
			printTest(baseURL+v.URL, v.TestType, v.ExpectedResult)
			res := doTest(baseURL+v.URL, v.TestType, v.ExpectedResult)
			if !res {
				failCount++
				fmt.Println("Failed")
			} else {
				fmt.Println("OK")
			}
		}
	}
	if c {
		fmt.Println("Check tests:")
		for _, v := range testsConf.CheckTests {
			printTest(baseURL+v.URL, v.TestType, v.ExpectedResult)
			res := doTest(baseURL+v.URL, v.TestType, v.ExpectedResult)
			if !res {
				failCount++
				fmt.Println("Failed")
			} else {
				fmt.Println("OK")
			}
		}
	}
	if r {
		fmt.Println("ResponseCode tests:")
		for _, v := range testsConf.ResponseCodeTests {
			printTest(baseURL+v, "HTTP code", "200 OK")

			res := doTest(baseURL+v, "code", "200 OK")
			if !res {
				failCount++
				fmt.Println("Failed")
			} else {
				fmt.Println("OK")
			}
		}
	}
	os.Exit(failCount)
}

// doTest takes URL and performs test against. TestType could be either
// exact: exact match with the string
// inc: response incudes the string
// regex: response match with regex
// code: response returns http code
func doTest(url string, testType string, expRes string) bool {
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	result := string(body)
	switch testType {
	case "", "code":
		if resp.Status != expRes {
			return false
		}
	case "exact":
		if result != expRes {
			return false
		}
	case "inc", "includes":
		if strings.Index(result, expRes) == -1 {
			return false
		}
	case "regex":
		match, err := regexp.MatchString(expRes, result)
		if err != nil {
			log.Fatalf("wrong regular expression %q. Error is %q", expRes, err.Error())
		}
		if !match {
			return false
		}
	}

	return true
}

func printTest(url string, testType string, expRes string) {
	fmt.Printf("%s %s %s: ", url, testType, expRes)
}
