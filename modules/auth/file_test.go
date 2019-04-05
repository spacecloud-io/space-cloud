package auth
import (
	"testing"
	"github.com/spaceuptech/space-cloud/config"
)

//{testName : "Success" , path : "/" , module:&Module{ fileRules :map[string]*config.FileRule{"create":&config.FileRule{ Prefix : "/" ,Rule : map[string]*config.Rule{"rule": &config.Rule{Rule:"allow"}} ,} ,"delete":&config.FileRule{ Prefix : "/" ,Rule : map[string]*config.Rule{"rule": &config.Rule{Rule:"allow"}} ,} ,"read":&config.FileRule{ Prefix : "/" ,Rule : map[string]*config.Rule{"rule": &config.Rule{Rule:"allow"}} ,}} } } ,

func TestGetFileRule(t *testing.T){
	
	fileRule := &config.FileRule{
		Prefix : "/" ,
		Rule : map[string]*config.Rule{"rule": &config.Rule{Rule:"allow"}} ,
	}

	fileRule1 := &config.FileRule{
		Prefix : "/folder" ,
		Rule : map[string]*config.Rule{"rule": &config.Rule{Rule:"allow"}} ,
	}

	var mod = []struct{
		module *Module
		testName string
		path string
	}{
		{testName : "Success" , path : "/" , module:&Module{ fileRules :map[string]*config.FileRule{"create":fileRule ,"delete":fileRule ,"read":fileRule} } } ,
		{testName : "Success" , path : "/folder" , module:&Module{ fileRules :map[string]*config.FileRule{"create":fileRule ,"delete":fileRule ,"read":fileRule} } } ,
		{testName : "Success" , path : "/folder" , module:&Module{ fileRules :map[string]*config.FileRule{"create":fileRule1 ,"delete":fileRule1 ,"read":fileRule1} } } ,
		{testName : "Success" , path : "/folder/file" , module:&Module{ fileRules :map[string]*config.FileRule{"create":fileRule1 ,"delete":fileRule1 ,"read":fileRule1} } } ,
	
		{testName : "Fail" , path : "/NewFolder/file" , module:&Module{ fileRules :map[string]*config.FileRule{"create":fileRule1 ,"delete":fileRule1 ,"read":fileRule1} } } ,
	}

	for _,test := range mod {
		t.Run(test.testName, func(t *testing.T) {
		
		 data, rules , err1 := (test.module).getFileRule(test.path)
		 if test.testName == "Success" {
			if err1 != nil {
				t.Error(data,rules,err1) 
			}
		 } else {
			if err1 != nil {
				t.Error(data,rules,err1) 
			}
		 }
		})
	}
}