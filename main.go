package main

import (
	"fmt"
	"log"
	"math"
	"syscall"
	"unsafe"
)

const (
	MEM_COMMIT  = 0x1000;
	MEM_RESERVE = 0x2000;
	MEM_RELEASE = 0x8000;
	PAGE_EXECUTE_READWRITE = 0x40;
	MEM = 200;
	PAGE_SIZE = 40;
)

var (
	kernel32     = syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc = kernel32.MustFindProc("VirtualAlloc");
	VirtualFree  = kernel32.MustFindProc("VirtualFree");
)

type Allocator struct {
	size uintptr;
	adress uintptr;
}

func (allocator Allocator) memAlloc() (uintptr, error) {
	addr, _, msg := VirtualAlloc.Call(0, allocator.size, MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	if msg != nil {
		fmt.Println(msg)
		return addr, nil
	} else {
		return 0, nil
	}
}

func (allocator Allocator) requestMemory (value int, pages []Page, freeMemory int, startPage *byte) ([]Page, int){
	if len(pages) == 0 {
		if value > PAGE_SIZE {
			fmt.Println("We need more than one page");
		} else {
			usedSpace := value;
			freeSpace := 40 - usedSpace;
			page := Page{
				startPage: startPage,
				sizePage: PAGE_SIZE,
				currentBlockSize: value,
				usedSpace: value,
				freeSpace: freeSpace,
			}

			fmt.Printf("%+v\n", page);
			fmt.Println("Current page test: ", page);
			pages = append(pages, page);
			fmt.Println("All current pages 2: ", pages);
			freeMemory -= value;
			return pages, freeMemory
		}
	} else {
		fmt.Println("pages containt pages")
		for i:=0; i < len(pages); i ++ {
			fmt.Println(" =>", pages[i]);

			if pages[i].currentBlockSize == value {
				fmt.Println("WE HAVE PAGE WITH CURRENT BLOCK SIZE", pages[i])
				fmt.Println("SEEE", pages[i].currentBlockSize);
				pages[i].usedSpace += value;
				pages[i].freeSpace -= value;
				freeMemory -= value;
				return pages, freeMemory;
			}
		}

		for i:= 0; i < len(pages); i++ {
			fmt.Println("pages", pages[i].currentBlockSize);
			fmt.Println("VALUE", value);
			fmt.Println("LOOP PAGE", pages[i]);
			if pages[i].currentBlockSize == value {
				fmt.Println("We have page with current block size", pages[i]);
				fmt.Println("SEEE", pages[i].currentBlockSize);
				pages[i].usedSpace += value;
				pages[i].freeSpace -= value;
				freeMemory -= value;
				return pages, freeMemory;
			} else {
				fmt.Println("We need to divided new page")
				usedSpace := value;
				freeSpace := 40 - usedSpace;
				page := Page{
					startPage: startPage,
					sizePage: PAGE_SIZE,
					currentBlockSize: value,
					usedSpace: value,
					freeSpace: freeSpace,
				}
				pages = append(pages, page);
				fmt.Println("All current pages 2: ", pages);
				freeMemory -= value;
				return pages, freeMemory
			}
		}
		return nil, 0
	}
	return nil, 0
}

type Page struct {
	startPage *byte;
	//nextPage *byte;
	sizePage int;
	currentBlockSize int;
	usedSpace int;
	freeSpace int;
	//status string;
}

func (page *Page) Init(startPage *byte, sizePage int, currectBlockSize int,
	usedSpace int, freeSpace int) {
	page.startPage = startPage;
	//page.nextPage = nextPage;
	page.sizePage = sizePage;
	page.currentBlockSize = currectBlockSize;
	page.usedSpace = usedSpace;
	page.freeSpace = freeSpace;
	//page.status = "Free";
}

func checkedSize (value float64) (float64){
	var res float64;
	res = math.Pow(2, math.Ceil(math.Log(value)/math.Log(2)));
	return res;
}

func main () {
	var allocator Allocator;
	allocator.size = MEM;
	addr, err:= allocator.memAlloc();
	if err != nil {
		log.Fatal(err)
	}
	allocator.adress = addr;
	_addr:= (*byte)(unsafe.Pointer(allocator.adress))
	memory := (*[MEM]byte)(unsafe.Pointer(addr));
	var freeMemory int;
	freeMemory = MEM;
	fmt.Println(allocator);
	fmt.Println(memory);
	var pages []Page;
	var addValue bool;
	addValue = true;
	for addValue {
		var value float64;
		fmt.Print("Enter size of new block: ");
		fmt.Scan(&value);
		next := checkedSize(value);
		var converValue int = int(next);
		//var pages []Page;
		_pages, freeMemory := allocator.requestMemory(converValue, pages,freeMemory, _addr);
		pages = _pages;

		for i:= 0; i < len(pages); i ++ {
			fmt.Println("Page: ",i, pages[i]);
		}
		fmt.Println("Free space:= ", freeMemory)

		if freeMemory < 0 {
			addValue = false;
		}
	}

}