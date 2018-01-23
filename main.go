package mergbigfolder

import(
	"fmt"
	"crypto/sha1"
	"os"
	"io"
	"math"
	"path/filepath"
)

type wfi struct{
	path string
	fi os.FileInfo
	digest []byte
}

type Walker struct{
	wfiList []wfi
	root string
}

func NewWalker(r string) *Walker{
	return &Walker{root:r}
}

func(w *Walker)Walk()error{
	return filepath.Walk(w.root,func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() == false && info.Mode().IsRegular(){
			if digest:=w.GetDigest(path); digest !=nil{
				w.wfiList=append(w.wfiList,wfi{path,info,digest})
			}
			return nil			
		}
		return nil
	})
}

// 8KB
const filechunk = 8192

func(w *Walker)GetDigest(path string)[]byte{

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

	// Start hash
	//hash := md5.New()
	hash := sha1.New()

	// Check each block
	for i := uint64(0); i < blocks; i++ {
		// Calculate block size
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		if i%2 == 0 {
			file.Seek(int64(blocksize),1)
			continue
		}

		// Make a buffer
		buf := make([]byte, blocksize)

		// Make a buffer
		file.Read(buf)

		
		// Write to the buffer
		io.WriteString(hash, string(buf))
	}

	// Output the results
	//fmt.Printf("%x\n", hash.Sum(nil))
	return hash.Sum(nil)
}


func FindDiff(src string,dst string){
	srcw:=NewWalker(src)
	dstw:=NewWalker(dst)

	srcw.Walk()
	dstw.Walk()
}
