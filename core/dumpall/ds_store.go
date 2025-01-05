// https://github.com/gehaxelt/ds_store/blob/master/ds_store.go
package dumpall

import (
	"bytes"
	"encoding/binary"
	"errors"
	"slack-wails/lib/clients"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

type Block struct {
	Allocator *Allocator
	Offset    uint32
	Size      uint32
	Data      []byte
	Pos       uint32
}

type Allocator struct {
	Data     []byte
	Pos      uint32
	Root     *Block
	Offsets  []uint32
	Toc      map[string]uint32
	FreeList map[uint32][]uint32
}

func NewBlock(a *Allocator, pos uint32, size uint32) (block *Block, err error) {
	if len(a.Data) < int(pos+0x4+size) {
		return nil, errors.New("not enought data")
	}
	block = &Block{Size: size, Allocator: a, Data: a.Data[pos+0x4 : pos+0x4+size]}
	return block, nil
}

func (block *Block) readUint32() (value uint32, err error) {
	if block.Size-block.Pos < 4 {
		return 0, errors.New("not enough bytes to read")
	}
	data := bytes.NewBuffer(block.Data)
	data.Next(int(block.Pos))
	binary.Read(data, binary.BigEndian, &value)
	block.Pos += 4
	return value, nil
}

func (block *Block) readByte() (value byte, err error) {
	if block.Size-block.Pos < 1 {
		return 0, errors.New("not enough bytes to read")
	}
	data := bytes.NewBuffer(block.Data)
	data.Next(int(block.Pos))
	binary.Read(data, binary.BigEndian, &value)
	block.Pos += 1
	return value, nil
}

func (block *Block) readBuf(length int) (buf []byte, err error) {
	if int(block.Size)-int(block.Pos) < length {
		return nil, errors.New("not enough bytes to read")
	}
	data := bytes.NewBuffer(block.Data)
	data.Next(int(block.Pos))
	buf = make([]byte, length)
	binary.Read(data, binary.BigEndian, &buf)
	block.Pos += uint32(length)
	return buf, nil
}

func (block *Block) readFileName() (name string, err error) {
	length, err := block.readUint32()
	if err != nil {
		return "", err
	}
	buf, err := block.readBuf(int(2 * length))
	if err != nil {
		return "", err
	}

	/*
		sid, err := block.readUint32()
		if err != nil {
			return "", err
		}
	*/
	block.skip(4)

	stype, err := block.readBuf(4)
	if err != nil {
		return "", err
	}

	t := string(stype)
	bytesToSkip := -1

	switch {
	case t == "bool":
		bytesToSkip = 1
	case t == "type" || t == "long" || t == "shor":
		bytesToSkip = 4
	case t == "comp" || t == "dutc":
		bytesToSkip = 8
	case t == "blob":
		blen, err := block.readUint32()
		if err != nil {
			break
		}
		bytesToSkip = int(blen)
	case t == "ustr":
		blen, err := block.readUint32()
		if err != nil {
			break
		}
		bytesToSkip = int(2 * blen)
	default:
		break
	}

	if bytesToSkip <= 0 {
		return "", errors.New("unknown file format")
	}
	block.skip(uint32(bytesToSkip))

	name = utf16be2utf8(buf)
	return name, nil
}

func (block *Block) skip(i uint32) {
	block.Pos += i
}

func NewAllocator(data []byte) (a *Allocator, err error) {
	a = &Allocator{Data: data} //bytes.NewBuffer(data)}
	a.Toc = make(map[string]uint32)
	a.FreeList = make(map[uint32][]uint32)

	offset, size, err := a.readHeader()
	if err != nil {
		return nil, err
	}

	a.Root, err = NewBlock(a, offset, size)
	if err != nil {
		return nil, err
	}

	err = a.readOffsets()
	if err != nil {
		return nil, err
	}

	err = a.readToc()
	if err != nil {
		return nil, err
	}

	err = a.readFreeList()
	if err != nil {
		return nil, err
	}

	return a, err
}

func (a *Allocator) GetBlock(bid uint32) (block *Block, err error) {
	if len(a.Offsets) <= int(bid) {
		return nil, errors.New("cannot find key in Offset-Table")
	}
	addr := a.Offsets[bid]

	offset := int(addr) & ^0x1f
	size := 1 << (uint(addr) & 0x1f)

	block, err = NewBlock(a, uint32(offset), uint32(size)) ///+4??
	if err != nil {
		return nil, errors.New("cannot create/read block")
	}
	return block, nil
}

func (a *Allocator) TraverseFromRootNode() (filenames []string, err error) {
	rootBlk, err := a.GetBlock(a.Toc["DSDB"])
	if err != nil {
		return nil, err
	}
	rootNode, err := rootBlk.readUint32()
	if err != nil {
		return nil, err
	}
	/*height, err := rootBlk.readUint32()
	if err != nil {
		return nil, err
	}
	recordsCount, err := rootBlk.readUint32()
	if err != nil {
		return nil, err
	}
	nodesCount, err := rootBlk.readUint32()
	if err != nil {
		return nil, err
	}
	blksize, err := rootBlk.readUint32()
	if err != nil {
		return nil, err
	}*/
	rootBlk.skip(4 * 4)

	return a.Traverse(rootNode)
}

func (a *Allocator) Traverse(bid uint32) (filenames []string, err error) {
	node, err := a.GetBlock(bid)
	if err != nil {
		return nil, err
	}
	nextPtr, err := node.readUint32()
	if err != nil {
		return nil, err
	}
	count, err := node.readUint32()
	if err != nil {
		return nil, err
	}
	if nextPtr > 0 {
		//This may be broken
		for i := 0; i < int(count); i++ {
			next, err := node.readUint32()
			if err != nil {
				return nil, err
			}
			files, err := a.Traverse(next)
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, files...)
			f, err := node.readFileName()
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, f)
		}
		files, err := a.Traverse(nextPtr)
		if err != nil {
			return nil, err
		}
		filenames = append(filenames, files...)
	} else {
		for i := 0; i < int(count); i++ {
			f, err := node.readFileName()
			if err != nil {
				return nil, err
			}
			filenames = append(filenames, f)
		}
	}

	return filenames, nil
}

func (a *Allocator) readFreeList() error {
	for i := 0; i < 32; i++ {
		blkcount, err := a.Root.readUint32()
		if err != nil {
			return err
		}
		if blkcount == 0 {
			continue
		}
		a.FreeList[uint32(i)] = make([]uint32, 0)
		for k := 0; k < int(blkcount); k++ {
			val, err := a.Root.readUint32()
			if err != nil {
				return err
			}
			if val == 0 {
				continue
			}
			a.FreeList[uint32(i)] = append(a.FreeList[uint32(i)], val)
		}
	}
	return nil
}

func (a *Allocator) readToc() error {
	toccount, err := a.Root.readUint32()
	if err != nil {
		return err
	}
	for i := toccount; i > 0; i-- {
		tlen, err := a.Root.readByte()
		if err != nil {
			return err
		}
		name, err := a.Root.readBuf(int(tlen))
		if err != nil {
			return err
		}
		value, err := a.Root.readUint32()
		if err != nil {
			return err
		}
		a.Toc[string(name)] = value
	}
	return nil
}

func (a *Allocator) readOffsets() error {
	count, err := a.Root.readUint32()
	if err != nil {
		return err
	}
	a.Root.skip(4)

	for offcount := int(count); offcount > 0; offcount -= 256 {
		for i := 0; i < 256; i++ {
			val, err := a.Root.readUint32()
			if err != nil {
				return err
			}
			if val == 0 {
				continue
			}
			a.Offsets = append(a.Offsets, val)
		}
	}
	return nil
}

func (a *Allocator) readHeader() (offset uint32, size uint32, err error) {
	data := bytes.NewBuffer(a.Data)
	if data.Len() < 32 {
		return offset, size, errors.New("header not long enough")
	}
	var magic1, magic, offset2 uint32

	binary.Read(data, binary.BigEndian, &magic1)
	if magic1 != 1 {
		return offset, size, errors.New("wrong magic bytes")
	}
	a.Pos += 4

	binary.Read(data, binary.BigEndian, &magic)
	if magic != 0x42756431 {
		return offset, size, errors.New("wrong magic bytes")
	}
	a.Pos += 4

	binary.Read(data, binary.BigEndian, &offset)
	a.Pos += 4
	binary.Read(data, binary.BigEndian, &size)
	a.Pos += 4
	binary.Read(data, binary.BigEndian, &offset2)
	if offset != offset2 {
		return offset, size, errors.New("offets do not match")
	}
	a.Pos += 4

	return offset, size, nil
}

func utf16be2utf8(utf16be []byte) string {
	n := len(utf16be)
	// 使用 unsafe.Slice 替代 reflect.SliceHeader
	shorts := unsafe.Slice((*uint16)(unsafe.Pointer(&utf16be[0])), n/2)
	// shorts may need byte-swapping
	for i := 0; i < n; i += 2 {
		shorts[i/2] = (uint16(utf16be[i]) << 8) | uint16(utf16be[i+1])
	}

	// Convert to []byte
	count := 0
	for i := 0; i < len(shorts); i++ {
		r := rune(shorts[i])
		if utf16.IsSurrogate(r) {
			i++
			r = utf16.DecodeRune(r, rune(shorts[i]))
		}
		count += utf8.RuneLen(r)
	}
	buf := make([]byte, count)
	bi := 0
	for i := 0; i < len(shorts); i++ {
		r := rune(shorts[i])
		if utf16.IsSurrogate(r) {
			i++
			r = utf16.DecodeRune(r, rune(shorts[i]))
		}
		bi += utf8.EncodeRune(buf[bi:], r)
	}
	return string(buf)
}

func ExtractDSStore(url string) ([]string, error) {
	_, body, err := clients.NewSimpleGetRequest(url, clients.NewHttpClient(nil, true))
	if err != nil {
		return nil, err
	}
	a, err := NewAllocator(body)
	if err != nil {
		return nil, err
	}
	filenames, err := a.TraverseFromRootNode()
	if err != nil {
		return nil, err
	}

	var urlRoot = strings.TrimSuffix(url, ".DS_Store")
	var result = []string{}
	for _, f := range filenames {
		result = append(result, urlRoot+f)
	}
	return result, nil
}
