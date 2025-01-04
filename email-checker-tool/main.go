package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main(){
    scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan(){
		checkDomain(scanner.Text())
	}

	err := scanner.Err()
	if err != nil{
		log.Fatal("could not read from input %v",err)
	}
  
}

func checkDomain(domain string){
	var hasMX,hasSPF,hasDMARC bool
	var spfRecord,dmarcRecord string

	mxRecord,err := net.LookupMX(domain)

	if err != nil{
		log.Printf("error in looking up for mx %v",err)
	}

	if len(mxRecord) > 0{
		hasMX = true
	}

	txtRecords,err := net.LookupTXT(domain)

	if err != nil{
		log.Printf("error in looking up for txt %v",err)
	}

	for _,record := range txtRecords{
		if strings.HasPrefix(record,"v=spf1"){
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords,err := net.LookupTXT("_dmarc." + domain)

	if err != nil{
		log.Printf("error in looking up for dmarc %v",err)
	}

	for _,record := range dmarcRecords{
		if strings.HasPrefix(record,"v=DMARC1"){
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%v ,%v ,%v ,%v ,%v ,%v",domain,hasMX,hasSPF,spfRecord,hasDMARC,dmarcRecord)



}