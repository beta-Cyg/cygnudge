package judge;

import (
	"fmt"
	"log"
	"path/filepath"
)

func pack(code_file string,uid int,pid string){
	code_file_path,err:=filepath.Abs(code_file)
	if err!=nil{
		//debug
		log.Fatalln(err)
	}
	fmt.Println("%v",code_file_path)
}

