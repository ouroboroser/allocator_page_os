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
	MEM = 10000;
	PAGE_SIZE = 1000;
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
			freeSpace := PAGE_SIZE - usedSpace;
			page := Page{
				startPage: startPage,
				sizePage: PAGE_SIZE,
				currentBlockSize: value,
				usedSpace: value,
				freeSpace: freeSpace,
			}
			pages = append(pages, page);
			freeMemory -= value;
			return pages, freeMemory
		}
	} else {
		for i:=0; i < len(pages); i ++ {
			fmt.Println("Work here 1")
			if pages[i].currentBlockSize == value && pages[i].freeSpace > value {
				_freeSpace := pages[i].freeSpace;
				fmt.Println("FREE SPACE: ", _freeSpace);
				if pages[i].freeSpace < value {
					fmt.Println(pages);
					fmt.Println("Free space: ", pages[i].freeSpace, "value: ", value);
					fmt.Println("not enough memory")
					usedSpace := value;
					freeSpace := PAGE_SIZE - usedSpace;
					fmt.Println("Word here 2")
					page := Page{
						startPage: startPage,
						sizePage: PAGE_SIZE,
						currentBlockSize: value,
						usedSpace: value,
						freeSpace: freeSpace,
					}
					pages = append(pages, page);
					freeMemory -= value;
					return pages, freeMemory
				}

			}

			if pages[i].currentBlockSize == value {
				if pages[i].freeSpace > value {
					pages[i].usedSpace += value;
					pages[i].freeSpace -= value;
					freeMemory -= value;
					return pages, freeMemory;
				}
			}

		}

		for i:= 0; i < len(pages); i++ {
			usedSpace := value;
			freeSpace := PAGE_SIZE - usedSpace;
			page := Page{
				startPage: startPage,
				sizePage: PAGE_SIZE,
				currentBlockSize: value,
				usedSpace: value,
				freeSpace: freeSpace,
			}
			pages = append(pages, page);
			freeMemory -= value;
			return pages, freeMemory
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
			fmt.Printf("%+v\n", pages[i]);
		}
		fmt.Println("Free space:= ", freeMemory)

		if freeMemory < 0 {
			addValue = false;
		}
	}
}