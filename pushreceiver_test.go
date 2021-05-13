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
	expect := "51fa79f49ac6d9fcd3af620550f9ddb8"
	result, err := generateSignature('1', "shareKey", "testBody")

	if err != nil {
		t.Fatalf("generate Error:%v", err)
	}

	if result != expect {
		t.Errorf("MD5[%s] != [%s]\n", expect, result)
	}
}

func TestGenerateSignatureSHA256(t *testing.T) {
	expect := "f166187fa46c79233a07a54699087f139b3d209bc3090c8a83cc2f42bca3a61e"
	result, err := generateSignature('5', "shareKey", "testBody")

	if err != nil {
		t.Fatalf("generate Error:%v", err)
	}

	if result != expect {
		t.Errorf("SHA256[%s] != [%s]\n", expect, result)
	}
}

func TestGenerateSignatureSHA512(t *testing.T) {
	expect := "582396bac0f2c5bfbcfe4a95007075281c59d2c526fea5bad385d9a6d102d88258326fc2f55465adac86997f288e5174ff9427215eca1cbe9f41dacd7e1f0f91"
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
