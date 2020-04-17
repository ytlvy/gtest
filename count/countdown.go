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
	sizeChan       chan fileItem
	pathchan       chan string
	root           string
	closingChan    chan struct{} // signal channel
	closedChan     chan struct{}
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
		Loop1:
			for {
				select {
				case <-f.closingChan:
					return
				case f.pathchan <- subdir:
					f.walkDir1(subdir)
					break Loop1
				default:
				}
			}
		} else {
		Loop:
			for {
				select {
				case <-f.closingChan:
					return
				case f.sizeChan <- fileItem{curdir, entry.Size()}:
					break Loop
				default:
				}
			}
		}
	}
}

func (f *CountdownManager) addPathInfo(path string) {
	// fmt.Println("addnewpath")
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

func (f *CountdownManager) stopDu() {
	select {
	case f.closingChan <- struct{}{}:
		<-f.closedChan
	case <-f.closedChan:
	}
}

//RunDu1 打印目录使用情况
// go run main.go -v "/Users/yanjieguo/Documents"
func (f *CountdownManager) RunDu1() {

	startTime := time.Now()

	flag.Parse()
	roots := flag.Args()
	if len(roots) < 1 {
		roots = []string{"."}
	}

	f.fileMap = make(map[string]*fsizeInfo)
	f.sizeChan = make(chan fileItem)
	f.pathchan = make(chan string)
	f.closingChan = make(chan struct{})
	f.closedChan = make(chan struct{})

	var ticker *time.Ticker

	//abort
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	defer func() {
		close(abort)
		close(f.closingChan)
		close(f.closedChan)
		close(f.sizeChan)
		close(f.pathchan)
		fmt.Println("closer run")
	}()

	root := roots[0]
	f.root = root
	go func() {
		f.walkDir1(root)
		f.stopDu()
	}()

	if *verbose {
		ticker = time.NewTicker(500 * time.Millisecond)
		defer func() {
			ticker.Stop()
		}()
	}
	isabort := false
loop:
	for {
		if *verbose {
			select {
			case <-ticker.C:
				fmt.Printf("counting press enter to cancel \n")
			default:
			}
		}

		select {
		case <-f.closingChan:
			break loop
		case item, ok := <-f.sizeChan:
			if !ok {
				break loop
			}
			// fmt.Printf("recieve data %v\n", item)

			for path, info := range f.fileMap {
				if strings.HasPrefix(item.path, path) {
					info.nfiles++
					info.nbytes += item.nbytes
				}
			}

			f.nfiles++
			f.nbytes += item.nbytes

		case newpath, ok := <-f.pathchan:
			if !ok {
				break loop
			}

			f.addPathInfo(newpath)

		case <-abort:
			isabort = true
			fmt.Println("Launch aborted")
			f.stopDu()
			break loop
		}
	}

	if !isabort {
		f.printDiskUsage()
	}
	fmt.Printf("usage time %f\n", time.Now().Sub(startTime).Seconds())
}

func (f *CountdownManager) printDiskUsage() {
	fmt.Printf("disk usage total %d files  %.1f GB\n", f.nfiles, float64(f.nbytes)/1e9)
	for path, info := range f.fileMap {
		relp, _ := filepath.Rel(f.root, path)
		fmt.Printf("  | %s  %dfiles  %.1fMB\n", relp, info.nfiles, float64(info.nbytes)/1e6)
	}
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
