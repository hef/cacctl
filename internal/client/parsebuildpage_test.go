package client

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestParseBuildPage(t *testing.T) {
	f, err := os.Open("testdata/build_page.html")
	if err != nil {
		panic(err)
	}

	form, err := parseBuildPage(f)
	if err != nil {
		t.Errorf("error parsing build page: %s", err)
	}

	if form.infra != "146a053564029f" {
		t.Errorf("got unexpected infra, expected %s, got %s", "146a053564029f", form.infra)
	}
	if form.token != "1c074d9ed538ed" {
		t.Errorf("got unexpected token, expected %s, got %s", "1c074d9ed538ed", form.token)
	}
}

func TestParseBuildPageWithMacButton(t *testing.T) {
	f, err := os.Open("testdata/build_page_with_mac_pro_button.html")
	if err != nil {
		panic(err)
	}

	form, err := parseBuildPage(f)
	if err != nil {
		t.Errorf("error parsing build page: %s", err)
	}

	if form.infra != "710704960be81" {
		t.Errorf("got unexpected infra, expected %s, got %s", "710704960be81", form.infra)
	}
	if form.token != "1de9a246e7fa01" {
		t.Errorf("got unexpected token, expected %s, got %s", "1de9a246e7fa01", form.token)
	}
}

func TestParseBuildPageV4(t *testing.T) {
	f, err := os.Open("testdata/build_page_v4.html")
	if err != nil {
		panic(err)
	}
	expectedForm := parsedBuildFormV4{
		cpu:       []int32{1, 2},
		ramMb:     []int32{512, 1024, 1536, 2048, 2560, 3072, 4096},
		storageGb: []int32{10, 20, 30, 40},
		os: map[string]int32{
			"CentOS 7.9 64bit":        151,
			"CentOS 8.3 64bit":        144,
			"Debian 9.13 64bit":       145,
			"FreeBSD 12.2 64bit":      135,
			"Ubuntu 18.04 LTS 64bit":  146,
			"Windows 10 2004 64bit":   124,
			"Windows Server 2012 Std": 125,
			"Windows Server 2016 Std": 126,
			"Windows Server 2019 Std": 127,
		},
		totalBuilds:       50,
		totalCpu:          2,
		totalRam:          4096,
		totalSsd:          40,
		percentUsedBuilds: 4,
		percentUsedCpu:    0,
		percentUsedRam:    0,
		percentUsedSSd:    0,
		infra:             "f4cbc572ed888",
		token:             "1b7f7f4dc33757",
		buildSubmit:       "161bcae4c9d113",
	}

	form, err := parseBuildPageV4(f)
	if err != nil {
		t.Errorf("error parsing build page: %s", err)
	}

	if !reflect.DeepEqual(expectedForm.cpu, form.cpu) {
		t.Errorf("expected cpu %v, got %v", expectedForm.cpu, form.cpu)
	}

	if !reflect.DeepEqual(expectedForm.ramMb, form.ramMb) {
		t.Errorf("expected ramMb %v, got %v", expectedForm.ramMb, form.ramMb)
	}

	if !reflect.DeepEqual(expectedForm.storageGb, form.storageGb) {
		t.Errorf("expected storageGb %v, got %v", expectedForm.storageGb, form.storageGb)
	}

	if !reflect.DeepEqual(expectedForm.os, form.os) {
		t.Errorf("expected os %v, got %v", expectedForm.os, form.os)
	}

	if expectedForm.totalBuilds != form.totalBuilds {
		t.Errorf("expected totalBuilds %v, got %v", expectedForm.totalBuilds, form.totalBuilds)
	}

	if expectedForm.totalCpu != form.totalCpu {
		t.Errorf("expected totalCpu %v, got %v", expectedForm.totalCpu, form.totalCpu)
	}

	if expectedForm.totalRam != form.totalRam {
		t.Errorf("expected totalRam %v, got %v", expectedForm.totalRam, form.totalRam)
	}

	if expectedForm.totalSsd != form.totalSsd {
		t.Errorf("expected totalSsd %v, got %v", expectedForm.totalSsd, form.totalSsd)
	}

	if expectedForm.percentUsedBuilds != form.percentUsedBuilds {
		t.Errorf("expected percentUsedBuilds %v, got %v", expectedForm.percentUsedBuilds, form.percentUsedBuilds)
	}

	if expectedForm.percentUsedCpu != form.percentUsedCpu {
		t.Errorf("expected percentUsedCpu %v, got %v", expectedForm.percentUsedCpu, form.percentUsedCpu)
	}

	if expectedForm.percentUsedRam != form.percentUsedRam {
		t.Errorf("expected percentUsedRam %v, got %v", expectedForm.percentUsedRam, form.percentUsedRam)
	}

	if expectedForm.percentUsedSSd != form.percentUsedSSd {
		t.Errorf("expected percentUsedSSd %v, got %v", expectedForm.percentUsedSSd, form.percentUsedSSd)
	}

	if expectedForm.buildSubmit != form.buildSubmit {
		t.Errorf("expected build-submit: %s, got %s", expectedForm.buildSubmit, form.buildSubmit)
	}

	if expectedForm.infra != form.infra {
		t.Errorf("expected infra: %s, got %s", expectedForm.infra, form.infra)
	}

	if expectedForm.token != form.token {
		t.Errorf("expected token: %s, got %s", expectedForm.token, form.token)
	}
}

func FuzzParseBuildPage(f *testing.F) {

	sampleFiles := []string{
		"testdata/build_page_with_mac_pro_button.html",
		"testdata/build_page_v4.html",
		"testdata/build_page.html",
	}

	for _, sampleFile := range sampleFiles {
		sampleReader, _ := os.Open(sampleFile)
		sample, _ := io.ReadAll(sampleReader)
		f.Add(sample)
	}

	f.Fuzz(func(t *testing.T, s []byte) {
		r := bytes.NewReader(s)
		out, err := parseBuildPageV4(r)
		if err != nil && out != nil {
			t.Errorf("%q, %v", out, err)
		}
	})
}

func TestXParseBuildPageV4(t *testing.T) {

	testTable := []struct {
		name  string
		input []byte
	}{
		{"nullbyte", []byte{0}},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)
			out, err := parseBuildPageV4(r)
			if err != nil && out != nil {
				t.Errorf("%q, %v", out, err)
			}
		})
	}
}
