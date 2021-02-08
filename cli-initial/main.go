package main

import (
	"fmt"

	"github.com/gocarina/gocsv"	

	"time"

	"sync"
	
	"strconv"

	"os"

	"log"



)


var wg = sync.WaitGroup{}

var ch = make(chan string,1000)



// "Title,Message 1,Message 2,Stream Delay,Run Times\nCLI Invoker Name,First Message,Second Msg,2,10"
type CliStreamerRecord struct {
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

type CliRunnerRecord struct {
	// How many streamer will run.
	Run         string `csv:"Run"`
	Title       string `csv:"Title"`
	Message1    string `csv:"Message 1"`
	Message2    string `csv:"Message 2"`
	StreamDelay int    `csv:"Stream Delay"`
	RunTimes    int    `csv:"Run Times"`
}

func (cliRunnerRecord CliRunnerRecord) CliStreamerRecord() CliStreamerRecord {
	return CliStreamerRecord{
		Title:       cliRunnerRecord.Title,
		Message1:    cliRunnerRecord.Message1,
		Message2:    cliRunnerRecord.Message2,
		StreamDelay: cliRunnerRecord.StreamDelay,
		RunTimes:    cliRunnerRecord.RunTimes,
	}
}
// function i wrote 
func (cliStreamerRecord CliStreamerRecord) runCliStreamerRecord(invokerNo int) {
      
      for i := 0 ; i < cliStreamerRecord.RunTimes ; i++ {

         // send msg1 to channel 
         ch <- fmt.Sprintf("CLI Invoker %v -> %v \n", invokerNo , cliStreamerRecord.Message1 ) 

         time.Sleep( time.Duration(cliStreamerRecord.StreamDelay) * time.Second)

         // send msg2 to channel 
         ch <- fmt.Sprintf("CLI Invoker %v -> %v \n", invokerNo , cliStreamerRecord.Message2 ) 
   	
      }

      wg.Done()

}

func (cliRunnerRecord CliRunnerRecord) CliStreamerRecordCsv() string {
	cliStreamerRecords := []CliStreamerRecord{cliRunnerRecord.CliStreamerRecord()}

	out, err := gocsv.MarshalString(cliStreamerRecords)

	if err != nil{
		panic(err)
	}

	return out
}


func Csv(cliRunners *[]CliRunnerRecord) string {
	out, err := gocsv.MarshalString(cliRunners)

	if err != nil {
		panic(err)
	}

	return out
}



func main() {

    // var args string 
	args := "Run,Title,Message 1,Message 2,Stream Delay,Run Times\n1,First Streamer,First Message,Second Msg,2,10\n2,Second Streamer,First Message,Second Msg,2,10"
	// rawArgs := os.Args
	// args := rawArgs[1]
	// fmt.Println(args)
	// fmt.Println(reflect.TypeOf(args))
	// wg.Add(1)

     // go logger()

	var cliRunners []CliRunnerRecord
	
	gocsv.UnmarshalString(
		args,
		&cliRunners)

	fmt.Print(Csv(&cliRunners))
	
	fmt.Println("---------------------------------")
	
	for i , runner := range cliRunners {

		// fmt.Println(i, ":")
		//fmt.Print(runner.CliStreamerRecordCsv())
		run_times , err := strconv.ParseUint(runner.Run, 10, 32)
		
		if err != nil {
			
			fmt.Println("non integer value")
			
			return 
		}
		
		for j := 0 ; j < int(run_times) ; j++ {

			wg.Add(1)
			
			go func(invokerNo int) {

				runner.CliStreamerRecord().runCliStreamerRecord(invokerNo)
			
			}(i+1)
		
		}

	
	}

	go func(){

		wg.Wait()

		close(ch)

	}()

	// open file 
	file , err := os.Create("out.text")
	// check for open file 
	if err != nil {
		log.Fatalf("could not open file")
		return 
	}
	
	defer file.Close()


	for r:= range ch {
		fmt.Print(r)
		if _ ,err := file.WriteString(r) ; err != nil {
			log.Fatalf("could not write to file")
		}
	}

}
