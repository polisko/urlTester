package main

import "testing"

func Test_doTest(t *testing.T) {
	type args struct {
		url      string
		testType string
		expRes   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Positive exact match", args: args{url: "https://gevo-util-dev.eu-de.mybluemix.net/test/isRunning", testType: "exact", expRes: "true"}, want: true},
		{name: "Positive code", args: args{url: "https://gevo-util-dev.eu-de.mybluemix.net/test/login/login.html", testType: "code", expRes: "200 OK"}, want: true},
		{name: "Positive includes", args: args{url: "https://gevo-util-dev.eu-de.mybluemix.net/test/login/login.html", testType: "inc", expRes: `<!-- <div id="login"></div> -->`}, want: true},
		{name: "Positive regex", args: args{url: "https://gevo-util-dev.eu-de.mybluemix.net/test/login/login.html", testType: "regex", expRes: `^<!DOCTYPE html>(?s).*</html>$`}, want: true},
		{name: "Negative exact match", args: args{url: "https://gevo-util-dev.eu-de.mybluemix.net/test/isRunning", testType: "exact", expRes: "false"}, want: false},
		{name: "Negative code", args: args{url: "https://gevo-util-dev.eu-de.mybluemix.net/test/login/", testType: "code", expRes: "200 OK"}, want: false},
		{name: "Negative includes", args: args{url: "https://gevo-util-dev.eu-de.mybluemix.net/test/login/login.html", testType: "inc", expRes: `<!-- <div id=""></div> -->`}, want: false},
		{name: "Negative regex", args: args{url: "https://gevo-util-dev.eu-de.mybluemix.net/test/login/login.html", testType: "regex", expRes: `^<html>(?s).*</html>$`}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doTest(tt.args.url, tt.args.testType, tt.args.expRes); got != tt.want {
				t.Errorf("doTest() = %v, want %v", got, tt.want)
			}
		})
	}
}
