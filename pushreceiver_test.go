package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetBodyAndSignatureGET(t *testing.T) {

	reqGet := httptest.NewRequest(http.MethodGet, "http://example.jp/?signature=testsign&body=testbody", nil)

	sig, body, err := getBodyAndSignature(reqGet)

	if sig != "testsign" || body != "testbody" || err != nil {
		t.Errorf("GET Error sig:[%v] body:[%v] err:[%v]", sig, body, err)
	}
}

func TestGetBodyAndSignaturePOST(t *testing.T) {

	testBody := "{'request':'body'}"
	//reqBody := bytes.NewBufferString(testBody)
	reqBody := strings.NewReader(testBody)
	reqPost := httptest.NewRequest(http.MethodPost, "http://example.jp/unused", reqBody)

	reqPost.Header.Set("Content-Type", "application/json")
	reqPost.Header.Set("X-Signature", "testsign")

	sig, body, err := getBodyAndSignature(reqPost)

	if sig != "testsign" || body != testBody || err != nil {
		t.Errorf("POST Error sig:[%v] body:[%v] err:[%v]", sig, body, err)
	}
}

func TestGetBodyAndSignatureFORM(t *testing.T) {

	testBody := "2145f42602c68433a544f08f9e28efd0aef45c15bc3"
	data := url.Values{}
	data.Set("signature", "testsign")
	data.Set("body", testBody)
	reqForm := httptest.NewRequest(http.MethodPost, "http://example.jp/unused", strings.NewReader(data.Encode()))
	reqForm.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	sig, body, err := getBodyAndSignature(reqForm)

	if sig != "testsign" || body != testBody || err != nil {
		t.Errorf("Form Error sig:[%v] body:[%v] err:[%v]", sig, body, err)
	}
}

func TestGenerateSignatureMD5(t *testing.T) {
	expect := "cf7afe9ca5a522cdf061d65022e7f297"
	result, err := generateSignature('1', "shareKey", "testBody")

	if err != nil {
		t.Fatalf("generate Error:%v", err)
	}

	if result != expect {
		t.Errorf("MD5[%s] != [%s]\n", expect, result)
	}
}

func TestGenerateSignatureSHA256(t *testing.T) {
	expect := "93afe018bbc231e6ad8e64ab96e16df7780d977e04b681362ab9bbf213801bf5"
	result, err := generateSignature('5', "shareKey", "testBody")

	if err != nil {
		t.Fatalf("generate Error:%v", err)
	}

	if result != expect {
		t.Errorf("SHA256[%s] != [%s]\n", expect, result)
	}
}

func TestGenerateSignatureSHA512(t *testing.T) {
	expect := "4f7a30f207e42145f42602c68433a544f08f9e28efd0aef45c15bc35bb49e9ec777377dc4b26f7580c731d0f757436db934df28d79e8bb613f4a2b9988752532"
	result, err := generateSignature('6', "shareKey", "testBody")

	if err != nil {
		t.Fatalf("generate Error:%v", err)
	}

	if result != expect {
		t.Errorf("SHA512[%s] != [%s]\n", expect, result)
	}
}

func TestGenerateSignatureFail(t *testing.T) {
	_, err := generateSignature('Z', "shareKey", "testBody")

	if err == nil {
		t.Fatal("Unexpected successful!")
	}
}
