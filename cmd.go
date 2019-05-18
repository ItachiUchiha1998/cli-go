package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/mgutz/ansi"
)

var ExpressData = []byte(`{
	"name": "project-name",
	"version": "1.0.0",
	"description": "",
	"main": "server.js",
	"scripts": {
	  "test": "echo \"Error: no test specified\" && exit 1"
	},
	"license": "ISC",
	"dependencies": {
	  "body-parser": "^1.18.3",
	  "cors": "^2.8.5",
	  "express": "^4.16.4",
	  "morgan": "^1.9.1",
	  "path": "^0.12.7"
	}
}`)

var ExpressRoute = []byte(`const express = require('express');
const router = express.Router();
const bodyParser = require('body-parser');

router.use(bodyParser.urlencoded({ extended: true }))

router.use(function(req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    res.header('Access-Control-Allow-Methods', 'PUT, GET, POST, DELETE, OPTIONS');
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
    res.header('Access-Control-Allow-Credentials', true);
    next();
});

router.get('/', async (req,res) => {
    res.send({success: true})
});


router.post('*',(req,res) => {
    res.status(404).send({success: false,message: 'Page does not exist!'});
});

router.get('*',(req,res) => {
    res.status(404).send({success: false,message: 'Page does not exist!'});
});

module.exports = router;`)

var ExpressServer = []byte(`const express = require('express');
	const bodyParser = require('body-parser');
	const logger = require('morgan');
	const path = require('path');
	const index = require('./routes/index');
	const app = express();
	const port = 3000;

	app.use(logger('dev'));
	app.use(bodyParser.json());
	app.use(bodyParser.urlencoded({ extended: false }));
	app.use('/', index);

	app.use(function(req, res, next) {
		res.header("Access-Control-Allow-Origin", "*");
		res.header('Access-Control-Allow-Methods', 'PUT, GET, POST, DELETE, OPTIONS');
		res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
		next();
	});

	app.listen(port,function(){
		console.log("Listening to port " + port);
	});
`)

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func CreateModule(mod string, data []byte) {
	file, err := os.OpenFile(mod, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func sum(numbers ...int64) int64 {
	res := int64(0)
	for _, num := range numbers {
		res += num
	}

	return res
}

func mul(numbers ...int64) int64 {
	res := int64(1)
	for _, num := range numbers {
		res *= num
	}

	return res
}

func sub(numbers ...int64) int64 {
	res := int64(numbers[0])
	for i := 1; i < len(numbers); i++ {
		res -= int64(numbers[i])
	}

	return res
}

func div(numbers ...int64) int64 {
	res := int64(numbers[0])
	for i := 1; i < len(numbers); i++ {
		res /= int64(numbers[i])
	}

	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(os.Stderr, err)
		}
		err = runCommand(cmdString)
		if err != nil {
			fmt.Println(os.Stderr, err)
		}
	}
}

func runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	switch arrCommandStr[0] {
	case "exit":
		os.Exit(0)
	case "plus":
		if len(arrCommandStr) < 3 {
			return errors.New("2 args required!")
		}
		arrNum := []int64{}
		for i, arg := range arrCommandStr {
			if i == 0 {
				continue
			}
			n, _ := strconv.ParseInt(arg, 10, 64)
			arrNum = append(arrNum, n)
		}
		fmt.Fprintln(os.Stdout, sum(arrNum...))
		return nil
	case "mul":
		if len(arrCommandStr) < 3 {
			return errors.New("2 args required!")
		}
		arrNum := []int64{}
		for i, arg := range arrCommandStr {
			if i == 0 {
				continue
			}
			n, _ := strconv.ParseInt(arg, 10, 64)
			arrNum = append(arrNum, n)
		}
		fmt.Fprintln(os.Stdout, mul(arrNum...))
		return nil
	case "sub":
		if len(arrCommandStr) < 3 {
			return errors.New("2 args required!")
		}
		arrNum := []int64{}
		for i, arg := range arrCommandStr {
			if i == 0 {
				continue
			}
			n, _ := strconv.ParseInt(arg, 10, 64)
			arrNum = append(arrNum, n)
		}
		fmt.Fprintln(os.Stdout, sub(arrNum...))
		return nil
	case "div":
		if len(arrCommandStr) < 3 {
			return errors.New("2 args required!")
		}
		arrNum := []int64{}
		for i, arg := range arrCommandStr {
			if i == 0 {
				continue
			}
			n, _ := strconv.ParseInt(arg, 10, 64)
			arrNum = append(arrNum, n)
		}
		fmt.Fprintln(os.Stdout, div(arrNum...))
		return nil
	case "lookIP":
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			os.Stderr.WriteString("Oops: " + err.Error() + "\n")
			os.Exit(1)
		}

		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					os.Stdout.WriteString(ipnet.IP.String() + "\n")
				}
			}
		}
		return nil
	case "lookOS":
		fmt.Println("OS :", runtime.GOOS)
		return nil
	case "project":
		var f string
		fmt.Print("Enter Framework name(express): ")
		fmt.Scanf("%s", &f)
		if f == "express" {
			var project string
			fmt.Print("Enter Project Name: ")
			fmt.Scanf("%s", &project)
			CreateDirIfNotExist(project)
			CreateDirIfNotExist(project + "/routes")
			CreateModule(project+"/server.js", ExpressServer)
			CreateModule(project+"/package.json", ExpressData)
			CreateModule(project+"/README.MD", []byte("Welcome to "+project))
			CreateModule(project+"/routes/index.js", ExpressRoute)
		}
		return nil
	case "ls":
		dirname := "."
		f, err := os.Open(dirname)
		if err != nil {
			log.Fatal(err)
		}
		files, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				cyan := ansi.ColorFunc("cyan+")
				fmt.Printf(cyan(file.Name()) + "\t")
			} else {
				fmt.Printf(file.Name() + "\t")
			}
		}
		fmt.Println()
		return nil
	}
	cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()

}
