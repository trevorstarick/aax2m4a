package main

import "fmt"
import "sync"
import "runtime"
import "io/ioutil"
import "path/filepath"
import "os/exec"

var tasks chan string
var wg sync.WaitGroup

func worker() {
	for f := range tasks {
		f = f[0 : len(f)-4]
		cmd := exec.Command("ffmpeg", "-y", "-activation_bytes", "deadbeef", "-i", f+".aax", "-c:a", "copy", "-vn", "tmp/"+f+".m4a")
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	wg.Done()
}

func main() {
	tasks = make(chan string)

	for i := 0; i <= runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker()
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".aax" {
			fmt.Println(f.Name())
			tasks <- f.Name()
		}
	}

	close(tasks)

	wg.Wait()
}
