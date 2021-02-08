package main 

import (
	"testing"
	"strings"
	"io/ioutil"
)

var dummy_file_content string

func init(){

    t , _ := ioutil.ReadFile("test.conf")
    dummy_file_content = string(t)
	blockRecords = []BlockRecord{
		{ "{" , 0 } , 
		{ "}" , 5 } ,
		{ "{" , 7 } ,
		{ "{" , 12} ,
		{ "}" , 14} ,
		{ "}" , 15} ,
	}
}
	

func TestTopLevelBlocks(t *testing.T)  {

	arr  := strings.Split(dummy_file_content,"\n")

	rslt := TopLevelBlocks(blockRecords,0,len(arr))

    expected := []int{ 0 , 5, 7 , 15  }
	
	for i := 0 ; i < len(expected) ; i+=1 {
	
		if expected[i] != rslt[i]{
	
			t.Fatalf("expected %v got %v",expected,rslt)
	
		}
	}

	return  
}

func TestGetNginxBlocks(t *testing.T) {
     
     // test file contensts 
     ngBlocks := GetNginxBlocks(dummy_file_content)

     arr  := strings.Split(dummy_file_content,"\n")

     if ngBlocks.AllContents != dummy_file_content {
     	t.Fatalf("AllContent does not match file content")
     }
     
     // test lilnes 
     for i := 0 ; i < len(arr) ; i+=1 {

     	 if *( (*ngBlocks.AllLines)[i] ) != arr[i] {

     	 	t.Fatalf(" NginxBlocks and File lines do not match ")
     	 } 
     }

    // fmt.Println(ngBlocks.blocks)
    // test nested levels 
   //  for i := 0 ; i < len(*ngBlocks.blocks) ; i+=1 {
   //  	block :=  (*(ngBlocks.blocks))[i] 
	 	// fmt.Println(  block.TotalBlocksInside )
   //  }

    return 
}


func TestGetNginxBlock(t *testing.T) {

	arr  := strings.Split(dummy_file_content,"\n") 
    
    var parr []*string 
	
	for i := 0 ; i < len(arr) ; i++ {
		parr = append(parr,&arr[i])
	}

	rslt := GetNginxBlock( &parr , 0 , 5 , 0)
    
    if rslt.TotalBlocksInside != 0 || rslt.StartLine != arr[0] || rslt.EndLine != arr[5] {

    	t.Fatalf("expected TotalBlocksInside , StartLine , EndLine %v %v %v vs got %v %v %v ",
    		0 , arr[0] , arr[5] , rslt.TotalBlocksInside , rslt.StartLine , rslt.EndLine, 
    	)
    }

	rslt = GetNginxBlock( &parr , 7 , 15 , 0)
    
    if rslt.TotalBlocksInside != 1 || rslt.StartLine != arr[7] || rslt.EndLine != arr[15] {

    	t.Fatalf("expected TotalBlocksInside , StartLine , EndLine %v %v %v vs got %v %v %v ",
    		0 , arr[0] , arr[5] , rslt.TotalBlocksInside , rslt.StartLine , rslt.EndLine, 
    	)
    }

    // fmt.Println(rslt.TotalBlocksInside)


	return 
}