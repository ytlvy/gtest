package count

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type fsizeInfo struct {
	nfiles int64
	nbytes int64
}

// CountdownManager ss
type CountdownManager struct {
	fileMap        map[string]*fsizeInfo
	nfiles, nbytes int64
	fileSizes      chan fileItem
	newpathchan    chan string
	root           string
}

//RunCount 计数
func (f *CountdownManager) RunCount() {
	ticker := time.NewTicker(1 * time.Second)

	//abort
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	fmt.Println("Commencing countdown. Press return to abort")

count:
	for countdown := 10; countdown > 0; countdown-- {
		select {
		case <-ticker.C:
			fmt.Printf("counting ... %02d\n", countdown)
		case <-abort:
			fmt.Println("Launch aborted")
			break count
		}

	}

	ticker.Stop()
}

var verbose = flag.Bool("v", false, "show verbose progress messages")
var layer = flag.Int("n", 0, "show how many layers")

type fileItem struct {
	path   string
	nbytes int64
}

func (f *CountdownManager) walkDir1(curdir string) {
	for _, entry := range f.dirents(curdir) {
		if entry.IsDir() {
			subdir := filepath.Join(curdir, entry.Name())

			//增加子目录
			f.newpathchan <- subdir
			f.walkDir1(subdir)
		} else {
			f.fileSizes <- fileItem{curdir, entry.Size()}
		}
	}
}

func (f *CountdownManager) addPathInfo(path string) {

	relpath, _ := filepath.Rel(f.root, path)
	paths := strings.Split(relpath, string(os.PathSeparator))
	if len(paths) > 1 {
		return
	}

	sizeinfo, ok := f.fileMap[path]
	if !ok {
		sizeinfo = &fsizeInfo{}
		f.fileMap[path] = sizeinfo
	}
}

//RunDu1 打印目录使用情况
// go run main.go -v "/Users/yanjieguo/Documents"
func (f *CountdownManager) RunDu1() {

	flag.Parse()
	roots := flag.Args()
	if len(roots) < 1 {
		roots = []string{"."}
	}

	f.fileMap = make(map[string]*fsizeInfo)
	f.fileSizes = make(chan fileItem)
	f.newpathchan = make(chan string)

	//abort
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	root := roots[0]
	f.root = root
	go func() {
		f.walkDir1(root)
		close(f.fileSizes)
		close(f.newpathchan)
		fmt.Println("closer run")
	}()

	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	isabort := false
loop:
	for {
		select {
		case item, ok := <-f.fileSizes:
			if !ok {
				break loop
			}

			for path, info := range f.fileMap {
				if strings.HasPrefix(item.path, path) {
					info.nfiles++
					info.nbytes += item.nbytes
				}
			}

			f.nfiles++
			f.nbytes += item.nbytes

		case <-tick:
			fmt.Printf("counting press enter to cancel \n")
		case newpath, ok := <-f.newpathchan:
			if !ok {
				break loop
			}

			f.addPathInfo(newpath)
		case <-abort:
			isabort = true
			fmt.Println("Launch aborted")
			close(abort)
			break loop
		}
	}

	if !isabort {
		f.printDiskUsage()
	}
}

func (f *CountdownManager) printDiskUsage() {
	fmt.Printf("disk usage total %d files  %.1f GB\n", f.nfiles, float64(f.nbytes)/1e9)
	for path, info := range f.fileMap {
		relp, _ := filepath.Rel(f.root, path)
		fmt.Printf("  | %s  %dfiles  %.1fMB\n", relp, info.nfiles, float64(info.nbytes)/1e6)
	}
}

// RunDu 打印目录使用情况
func (f *CountdownManager) RunDu() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) < 1 {
		roots = []string{"."}
	}

	fileSizes := make(chan int64, 0)
	go func() {
		for _, root := range roots {
			f.walkDir(root, fileSizes)
		}
		close(fileSizes)
	}()

	var nfiles, nbytes int64
	for size := range fileSizes {
		nfiles++
		nbytes += size
	}

	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func (f *CountdownManager) walkDir(dir string, fileSizes chan<- int64) {
	for _, entry := range f.dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			f.walkDir(subdir, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

func (f *CountdownManager) currentFiles(dir string) []string {
	mfils := make([]string, 100)
	for _, entry := range f.dirents(dir) {
		filepath := filepath.Join(dir, entry.Name())
		mfils = append(mfils, filepath)
	}
	return mfils
}

func (f *CountdownManager) dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}

	return entries
}
