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

func (allocator Allocator) requestMemory (memory *[MEM]byte, value int, pages []Page, freeMemory int, startPage *byte) ([]Page, int, *byte){
	var nextPage *byte;
	if freeMemory < value {
		fmt.Errorf("Not enough memory");
		return nil, 0, nil;
	}
	if len(pages) == 0 {
		if value > PAGE_SIZE {
			cutValue := float64(value)/float64(PAGE_SIZE) + 1;
			cutValue = math.Round(cutValue);
			g := int((PAGE_SIZE * cutValue));
			next := &memory[g];
			page := Page{
				startPage: startPage,
				nextPage: next,
				sizePage: PAGE_SIZE,
				currentBlockSize: value,
				usedSpace: value,
				freeSpace: 0,
				status: "Diveded",
			}
			pages = append(pages, page);
			freeMemory -= value;
			return pages, freeMemory, nil
		} else {
			next := &memory[PAGE_SIZE];
			usedSpace := value;
			freeSpace := PAGE_SIZE - usedSpace;
			page := Page{
				startPage: startPage,
				nextPage: next,
				sizePage: PAGE_SIZE,
				currentBlockSize: value,
				usedSpace: value,
				freeSpace: freeSpace,
				status: "Diveded",
			}
			pages = append(pages, page);
			freeMemory -= value;
			return pages, freeMemory, next;
		}
	} else {
		if value > PAGE_SIZE {
			page := Page{
				startPage: startPage,
				sizePage: PAGE_SIZE,
				currentBlockSize: value,
				usedSpace: value,
				freeSpace: 0,
				status: "Diveded",
			}
			fmt.Println("Page: ", page);
		}
		for i:=0; i < len(pages); i ++ {
			if value > PAGE_SIZE {
			} else {
				if pages[i].currentBlockSize == value && pages[i].freeSpace > value {
					if pages[i].freeSpace < value {
						page := Page{
							startPage: startPage,
							sizePage: PAGE_SIZE,
							currentBlockSize: value,
							usedSpace: value,
							freeSpace: 1,
							status: "Diveded",
						}
						pages = append(pages, page);
						freeMemory -= value;
						return pages, freeMemory, nextPage;
					}

				}
			}


			if pages[i].currentBlockSize == value {
				if pages[i].freeSpace > value {
					pages[i].usedSpace += value;
					pages[i].freeSpace -= value;
					freeMemory -= value;
					return pages, freeMemory, nextPage;
				}
			}

		}
		for i:= 0; i < len(pages); i++ {
			usedSpace := value;
			freeSpace := PAGE_SIZE - usedSpace;
			tmp := 0;
			for i:= 0; i < MEM; i ++ {
				if &memory[i] == startPage {
					tmp = i;
				}
			}
			nextPage = &memory[tmp+PAGE_SIZE+1];
			var status string;
			var t int;

			if value > PAGE_SIZE {
				t = 0;
				status = "Multiple page block";
			} else {
				t = freeSpace;
				status = "Diveded";
			}
			page := Page{
				startPage: startPage,
				nextPage: nextPage,
				sizePage: PAGE_SIZE,
				currentBlockSize: value,
				usedSpace: value,
				freeSpace: t,
				status: status,
			}
			pages = append(pages, page);
			freeMemory -= value;
			return pages, freeMemory, nextPage
		}
		return nil, 0 , nextPage;
	}
	return nil, 0, nextPage;
}

func (allocator Allocator) freeAlloc() error {
	addr, _, msg := VirtualFree.Call(allocator.adress, 0 , MEM_RELEASE)
	if addr == 0 {
		fmt.Println(msg)
		return msg;
	} else {
		fmt.Println(msg)
		return nil
	}
}

type Page struct {
	startPage *byte;
	nextPage *byte;
	sizePage int;
	currentBlockSize int;
	usedSpace int;
	freeSpace int;
	status string;
}

func (page *Page) Init(startPage *byte, nextPage *byte, sizePage int, currectBlockSize int,
	usedSpace int, freeSpace int, status string) {
	page.startPage = startPage;
	page.nextPage = nextPage;
	page.sizePage = sizePage;
	page.currentBlockSize = currectBlockSize;
	page.usedSpace = usedSpace;
	page.freeSpace = freeSpace;
	page.status = status;
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
	var pages []Page;
	var addValue bool;
	addValue = true;
	for addValue {
		var value float64;
		fmt.Print("Enter size of new block: ");
		fmt.Scan(&value);
		next := checkedSize(value);
		var converValue int = int(next);
		_pages, _freeMemory, _next := allocator.requestMemory(memory, converValue, pages,freeMemory, _addr);
		_addr = _next;
		freeMemory = _freeMemory;
		pages = _pages;

		for i:= 0; i < len(pages); i ++ {
			fmt.Printf("%+v\n", pages[i]);
		}
		fmt.Println("Free space:= ", freeMemory)
		if freeMemory == 0 {
			addValue = false;
		}
	}
	allocator.freeAlloc();
}