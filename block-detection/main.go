package main 

import (
  "fmt"
  "io/ioutil"
  "log"
  "strings"
  // "reflect"
)



func (ngBlock *NginxBlock) IsBlock(line string) bool {
  // TODO Solve it using regex

  return true 
}

func (ngBlock *NginxBlock) IsLine(line string) bool {
  // TODO Solve it using regex

  return true 
}

func (ngBlock *NginxBlock) HasComment(line string) bool {
  // TODO Solve it using regex
  return true 
}


type BlockRecord struct {
     block string 
     lineNo int 
}


var blockRecords []BlockRecord = []BlockRecord{ };


type NginxBlock struct {
  StartLine   string
  EndLine     string
  AllContents string
  // split lines by \n on AllContents,
  // use make to create *[],
  // first create make([]*Type..)
  // then use &var to make it *
  AllLines          *[]*string
  NestedBlocks      []*NginxBlock
  TotalBlocksInside int
}



type NginxBlocks struct {
  blocks      *[]*NginxBlock
  AllContents string
  // split lines by \n on AllContents
  AllLines *[]*string
}



func GetNginxBlock(
  lines *[]*string,
  startIndex,
  endIndex,
  recursionMax int,
) *NginxBlock {

      nginxBlock := &NginxBlock{}
      
      nginxBlock.StartLine = *((*lines)[startIndex])
      nginxBlock.EndLine   = *((*lines)[endIndex])

      tLines := (*lines)[startIndex:endIndex+1] 
      nginxBlock.AllLines  = &tLines
      
      // fmt.Println(nginxBlock.StartLine)
      // fmt.Println(nginxBlock.EndLine)

      blocks := TopLevelBlocks( blockRecords,startIndex + 1 , endIndex - 1)
      // vartNestedBlocks *[]*NginxBlock
      for i := 0 ; i < len(blocks) / 2; i+=2 {
          nginxBlock.NestedBlocks = append( nginxBlock.NestedBlocks , GetNginxBlock( 
               lines ,
               blocks[i],
               blocks[i+1],
               0 ,
          ))
      }
      nginxBlock.TotalBlocksInside = len(nginxBlock.NestedBlocks)
      // fmt.Println(nginxBlock.TotalBlocksInside)
      return nginxBlock 
}


func GetNginxBlocks(configContent string) *NginxBlocks { 
  // 
  
      nginxBlocks := &NginxBlocks{}
      
      nginxBlocks.AllContents = configContent
      
      arr := strings.Split( configContent , "\n" )
      
      var parr []*string ;
      
      // put the lines in pointer array and collect block positions 
      for i := 0 ; i < len(arr) ; i++ {

          t := arr[i]
        
          parr = append(parr,&t)

          // if string contains { or } store in blockRecord with line no 
          if idx := strings.Index(t,"{") ; idx != -1 {

               blockRecords = append( blockRecords , BlockRecord{ "{" ,i } )
          
          }else if idx := strings.Index(t,"}") ; idx != -1 {

                blockRecords = append( blockRecords , BlockRecord{ "}" ,i } )
          
          }
      }

      nginxBlocks.AllLines = &parr 

      blocks_idxs := TopLevelBlocks(blockRecords,0,len(parr)-1)
     
      var tBlocks []*NginxBlock 

      for i := 0 ; i <= len(blocks_idxs) / 2 ; i+=2{

          t := GetNginxBlock( &parr , blocks_idxs[i] , blocks_idxs[i+1] , 0)

          tBlocks = append( tBlocks , t)

      } 
      
      nginxBlocks.blocks = &tBlocks 

      return nginxBlocks

}



func TopLevelBlocks(blockRecords []BlockRecord,startLine , endLine int) []int {
     
     t := 0
     
     rslt := []int{}

     for i := 0 ; i < len( blockRecords ) ;  i++ {

         if !(blockRecords[i].lineNo >= startLine && blockRecords[i].lineNo <= endLine) {
           continue
         } 

         if blockRecords[i].block == "{" {
    
            if t == 0 {
              rslt = append( rslt , blockRecords[i].lineNo )
            }

            t += 1
 
         }else {
            
            t -= 1 

            if t == 0 {
              rslt = append( rslt , blockRecords[i].lineNo )
            }

         }

     }


     return rslt  
}



func main(){

  fmt.Println("block detection")
  // read the nginx conf file 
  file_content , err := ioutil.ReadFile("nginx.conf")
  
  if err != nil {
    log.Fatalf("could not open file")
  } 

  t :=   GetNginxBlocks(string(file_content))
  
  fmt.Println(t)

 
}
