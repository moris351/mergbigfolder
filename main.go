package mergbigfolder

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	humanize "github.com/dustin/go-humanize"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type wfi struct {
	path   string
	fi     os.FileInfo
	digest []byte
}

type Walker struct {
	wfiList []wfi
	root    string
}
type ByDigest []wfi

func (w wfi) String() string {
	return fmt.Sprintf("path:%s, digest:%x", w.path, w.digest)
}

func (b ByDigest) Len() int           { return len(b) }
func (b ByDigest) Less(i, j int) bool { return bytes.Compare(b[i].digest, b[j].digest) == -1 }
func (b ByDigest) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func NewWalker(r string) *Walker {
	return &Walker{root: r}
}

func (w *Walker) Walk() error {
	return filepath.Walk(w.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("filepath.Walk failed with err,", err)
			return err
		}

		if info.IsDir() == false && info.Mode().IsRegular() {
			if digest := w.GetDigest(path); digest != nil {
				w.wfiList = append(w.wfiList, wfi{path, info, digest})
			}
			return nil
		}
		return nil
	})
}

// 8KB
const filechunk = 8192

func (w *Walker) GetDigest(path string) []byte {

	start := time.Now()
	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Cannot find file:", path)
		return nil
	}

	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		fmt.Println("Cannot access file:", path)
		return nil
	}

	// Get the filesize
	filesize := info.Size()

	// Calculate the number of blocks
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	fmt.Println("blocks=",blocks)

	// Start hash
	//hash := md5.New()
	hash := sha1.New()

	// Check each block
	for i := uint64(0); i < blocks; i++ {
		// Calculate block size
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		if i%2 == 1 {
			file.Seek(int64(blocksize), 1)
			continue
		}

		// Make a buffer
		buf := make([]byte, blocksize)

		// Make a buffer
		file.Read(buf)

		// Write to the buffer
		hash.Write(buf)
	}

	// Output the results
	digest := hash.Sum(nil)
	end := time.Now()

	fmt.Printf("cost %s %x %s %s\n", end.Sub(start), digest, path, humanize.Bytes(uint64(info.Size())))
	return digest
}

func FindDiff(src string, dst string) error {
	fmt.Println("FindDiff call:")
	srcw := NewWalker(src)
	dstw := NewWalker(dst)

	if err := srcw.Walk(); err != nil {
		return err
	}
	if err := dstw.Walk(); err != nil {
		return err
	}

	sort.Sort(ByDigest(srcw.wfiList))
	fmt.Printf("srcw: %v\n", srcw.wfiList)

	sort.Sort(ByDigest(dstw.wfiList))
	fmt.Printf("dstw: %v\n", dstw.wfiList)

	srcl:=len(srcw.wfiList)
	dstl := len(dstw.wfiList)
	fmt.Println("dstl=", dstl)

	//length:=math.Max(float64(srcl),float64(dstl))
	printResult:=func(isSrc bool,w wfi){
		if !isSrc{
			fmt.Printf("......................")
		}
			
		fmt.Printf("%s %s \n",w.path, humanize.Bytes(uint64(w.fi.Size())))
	}
	for i,j:=0,0;i<srcl || j<dstl;{
		if i<srcl && j< dstl {
			srcwfi:=srcw.wfiList[i]
			dstwfi:=dstw.wfiList[j]
			n:=bytes.Compare(srcwfi.digest,dstwfi.digest)
			if n==0{
				i++
				j++
			}else if n==-1{
				printResult(true,srcwfi)
				i++
			}else if n==1{
				printResult(false,dstwfi)
				j++	
			}
			continue
		}
		if i<srcl{
			printResult(true,srcw.wfiList[i])
			i++
		}
		if j<dstl{
			printResult(false,dstw.wfiList[j])
			j++
		}

	}
	return nil
}
