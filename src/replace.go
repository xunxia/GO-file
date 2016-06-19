package main 

import (
	"io/ioutil"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sort"
	"strconv"
)

type replaceNum int

type TextFile struct{
	name string
	replNum replaceNum
	}

type DecriOrder func(f1, f2 *TextFile) bool

func (decriOrder DecriOrder) Sort(textFiles []TextFile){
	tfs :=&fileSorter{
		textFiles: textFiles,
		decriOrder:decriOrder,
		} 
	sort.Sort(tfs)
	}

type fileSorter struct{
	textFiles []TextFile
	decriOrder func(f1,f2 *TextFile) bool
	}

func(s *fileSorter) Len() int{
	return len(s.textFiles)
	}

func(s *fileSorter) Swap(i,j int){
	s.textFiles[i],s.textFiles[j]=s.textFiles[j],s.textFiles[i]
	}

func(s *fileSorter) Less(i,j int) bool{
	return s.decriOrder(&s.textFiles[i], &s.textFiles[j])
	}

func main() {
	findWord:= os.Args[1]	
	fmt.Println("The word to replace is :",findWord)
	r:=regexp.MustCompile(`(?i)`+findWord)
	trackInfor := make([]TextFile,100)
	index:=0
	walk:= func(path string, info os.FileInfo, err error) error{				
		if !info.IsDir(){
			if strings.HasSuffix(path, ".txt"){				
				var numR replaceNum
				num :=0
				data, err:=ioutil.ReadFile(path)
				if err!=nil{
					panic(err)
				}
				dataStr:=string(data)
				dataRep:= r.ReplaceAllString(dataStr, strings.ToUpper(findWord))								
				num = strings.Count(strings.ToLower(dataStr),strings.ToLower(findWord))								
				dataBac:= []byte(dataRep)
				err = ioutil.WriteFile(path, dataBac,0644)				
				if err!= nil{
					panic(err)
					}
				if num>0{
					var file TextFile
					file.name = path
					numR=replaceNum(num)
					file.replNum = numR
					trackInfor[index] = file
					index++
					}				
				}						
			}
		return err
		}	
	er:=filepath.Walk("./", walk)
	if er!= nil{
		panic(er)
		}
	textfile := make([]TextFile, index)
	copy(textfile,trackInfor)
	replNum:=func(f1,f2 *TextFile)bool{
		return f1.replNum<f2.replNum
		}
	decreasingReplNum:=func(f1,f2 *TextFile)bool{
		return !replNum(f1,f2)
		}
	
	DecriOrder(replNum).Sort(textfile)
	DecriOrder(decreasingReplNum).Sort(textfile)
	repData:="Changes"+"         "+"File Name \r\n"
	fmt.Println("Changes"+"         "+"File Name")
	for _,f:=range textfile{
		number:=strconv.Itoa(int(f.replNum))
		repData += number+"               "+f.name+" \r\n"
		fmt.Printf("%d               %s\n", f.replNum, f.name)
		}
	file, err:=os.Create("report.txt")
	if err!= nil{
		panic(err)
		}
	file.WriteString(repData)	
	file.Close()
}