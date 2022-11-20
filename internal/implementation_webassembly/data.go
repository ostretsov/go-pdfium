package implementation_webassembly

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/klippa-app/go-pdfium/structs"
)

type CString struct {
	Pointer uint64
	Free    func()
}

func (p *PdfiumImplementation) CString(input string) (*CString, error) {
	inputLength := uint64(len(input)) + 1

	pointer, err := p.Malloc(inputLength)
	if err != nil {
		return nil, err
	}

	// Write string + null terminator.
	if !p.Module.Memory().Write(p.Context, uint32(pointer), append([]byte(input), byte(0))) {
		return nil, errors.New("could not write CString data")
	}

	return &CString{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
	}, nil
}

type IntPointer struct {
	Pointer uint64
	Free    func()
	Value   func() (int, error)
}

func (p *PdfiumImplementation) IntPointer() (*IntPointer, error) {
	pointer, err := p.Malloc(p.CSizeInt())
	if err != nil {
		return nil, err
	}

	return &IntPointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() (int, error) {
			b, success := p.Module.Memory().Read(p.Context, uint32(pointer), uint32(p.CSizeInt()))
			if !success {
				return 0, errors.New("could not read int data from memory")
			}

			var myInt int32
			err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &myInt)
			if err != nil {
				return 0, err
			}

			return int(myInt), nil
		},
	}, nil
}

type UIntPointer struct {
	Pointer uint64
	Free    func()
	Value   func() (uint, error)
}

func (p *PdfiumImplementation) UIntPointer() (*UIntPointer, error) {
	pointer, err := p.Malloc(p.CSizeUInt())
	if err != nil {
		return nil, err
	}

	return &UIntPointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() (uint, error) {
			b, success := p.Module.Memory().Read(p.Context, uint32(pointer), uint32(p.CSizeUInt()))
			if !success {
				return 0, errors.New("could not read uint data from memory")
			}

			var myInt uint32
			err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &myInt)
			if err != nil {
				return 0, err
			}

			return uint(myInt), nil
		},
	}, nil
}

type UIntArrayPointer struct {
	Pointer uint64
	Size    uint64
	Free    func()
	Value   func() ([]uint, error)
}

func (p *PdfiumImplementation) UIntArrayPointer(size uint64) (*UIntArrayPointer, error) {
	pointer, err := p.Malloc(p.CSizeUInt() * size)
	if err != nil {
		return nil, err
	}

	return &UIntArrayPointer{
		Pointer: pointer,
		Size:    size,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() ([]uint, error) {
			myInts := []uint{}

			for i := 0; i < 5; i++ {
				b, success := p.Module.Memory().Read(p.Context, uint32(pointer+(uint64(i)*p.CSizeUInt())), uint32(p.CSizeUInt()))
				if !success {
					return nil, errors.New("could not read uint array data from memory")
				}

				var myInt uint32
				err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &myInt)
				if err != nil {
					return nil, err
				}
			}

			return myInts, nil
		},
	}, nil
}

type DoublePointer struct {
	Pointer uint64
	Free    func()
	Value   func() (float64, error)
}

func (p *PdfiumImplementation) DoublePointer() (*DoublePointer, error) {
	pointer, err := p.Malloc(p.CSizeDouble())
	if err != nil {
		return nil, err
	}

	return &DoublePointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() (float64, error) {
			val, success := p.Module.Memory().ReadFloat64Le(p.Context, uint32(pointer))
			if !success {
				return 0, errors.New("could not read double data from memory")
			}

			return val, nil
		},
	}, nil
}

type ByteArrayPointer struct {
	Pointer uint64
	Free    func()
	Value   func() ([]byte, error)
}

func (p *PdfiumImplementation) ByteArrayPointer(size uint64) (*ByteArrayPointer, error) {
	pointer, err := p.Malloc(size)
	if err != nil {
		return nil, err
	}

	return &ByteArrayPointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() ([]byte, error) {
			b, success := p.Module.Memory().Read(p.Context, uint32(pointer), uint32(size))
			if !success {
				return nil, errors.New("could not read byte array data from memory")
			}

			return b, nil
		},
	}, nil
}

type LongPointer struct {
	Pointer uint64
	Free    func()
	Value   func() (int64, error)
}

func (p *PdfiumImplementation) LongPointer() (*LongPointer, error) {
	pointer, err := p.Malloc(p.CSizeLong())
	if err != nil {
		return nil, err
	}

	return &LongPointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() (int64, error) {
			b, success := p.Module.Memory().Read(p.Context, uint32(pointer), uint32(p.CSizeLong()))
			if !success {
				return 0, errors.New("could not read long data from memory")
			}

			var myInt int64
			err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &myInt)
			if err != nil {
				return 0, err
			}

			return myInt, nil
		},
	}, nil
}

type ULongPointer struct {
	Pointer uint64
	Free    func()
	Value   func() (uint64, error)
}

func (p *PdfiumImplementation) ULongPointer() (*ULongPointer, error) {
	pointer, err := p.Malloc(p.CSizeULong())
	if err != nil {
		return nil, err
	}

	return &ULongPointer{
		Pointer: pointer,
		Free: func() {
			p.Free(pointer)
		},
		Value: func() (uint64, error) {
			val, success := p.Module.Memory().ReadUint64Le(p.Context, uint32(pointer))
			if !success {
				return 0, errors.New("could not read long data from memory")
			}

			return val, nil
		},
	}, nil
}

func (p *PdfiumImplementation) Malloc(size uint64) (uint64, error) {
	results, err := p.Functions["malloc"].Call(p.Context, size)
	if err != nil {
		return 0, err
	}

	pointer := results[0]

	ok := p.Module.Memory().Write(p.Context, uint32(results[0]), make([]byte, size))
	if !ok {
		return 0, errors.New("could not write nulls to memory")
	}

	return pointer, nil
}

func (p *PdfiumImplementation) Free(pointer uint64) error {
	_, err := p.Functions["free"].Call(p.Context, pointer)
	if err != nil {
		return err
	}
	return nil
}

func (p *PdfiumImplementation) CSizeUInt() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 4
}

func (p *PdfiumImplementation) CSizeInt() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 4
}

func (p *PdfiumImplementation) CSizeFloat() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 4
}

func (p *PdfiumImplementation) CSizeDouble() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 8
}

func (p *PdfiumImplementation) CSizeLong() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 8
}

func (p *PdfiumImplementation) CSizeULong() uint64 {
	// @todo: implement on pdfium/emscripten side?
	return 8
}

func (p *PdfiumImplementation) CSizeStructFS_MATRIX() uint64 {
	// FS_MATRIX is 6 * float (a, b, c, d, e, f).
	return p.CSizeFloat() * 6
}

func (p *PdfiumImplementation) CStructFS_MATRIX(in *structs.FPDF_FS_MATRIX) (uint64, func() (*structs.FPDF_FS_MATRIX, error), error) {
	pointer, err := p.Malloc(p.CSizeStructFS_MATRIX())
	if err != nil {
		return 0, nil, err
	}

	if in != nil {
		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.A) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}

		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.B) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}

		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.C) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}

		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.D) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}

		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.E) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}

		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.F) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}
	}

	return pointer, func() (*structs.FPDF_FS_MATRIX, error) {
		a, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		b, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer+(p.CSizeFloat()*1)))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		c, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer+(p.CSizeFloat()*2)))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		d, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer+(p.CSizeFloat()*3)))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		e, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer+(p.CSizeFloat()*4)))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		f, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer+(p.CSizeFloat()*5)))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		return &structs.FPDF_FS_MATRIX{
			A: a,
			B: b,
			C: c,
			D: d,
			E: e,
			F: f,
		}, nil
	}, nil
}

func (p *PdfiumImplementation) CSizeStructFS_RECTF() uint64 {
	// FS_RECTF is 4 * float (left, top, right, bottom).
	return p.CSizeFloat() * 4
}

func (p *PdfiumImplementation) CStructFS_RECTF(in *structs.FPDF_FS_RECTF) (uint64, func() (*structs.FPDF_FS_RECTF, error), error) {
	pointer, err := p.Malloc(p.CSizeStructFS_RECTF())
	if err != nil {
		return 0, nil, err
	}

	if in != nil {
		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Left) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}

		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Top) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}

		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Right) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}

		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Bottom) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}
	}

	return pointer, func() (*structs.FPDF_FS_RECTF, error) {
		left, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		top, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer+(p.CSizeFloat()*1)))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		right, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer+(p.CSizeFloat()*2)))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		bottom, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer+(p.CSizeFloat()*3)))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		return &structs.FPDF_FS_RECTF{
			Left:   left,
			Top:    top,
			Right:  right,
			Bottom: bottom,
		}, nil
	}, nil
}

func (p *PdfiumImplementation) CSizeStructFS_SIZEF() uint64 {
	// FS_SIZEF is 2 * float (width, height).
	return p.CSizeFloat() * 2
}

func (p *PdfiumImplementation) CStructFS_SIZEF(in *structs.FPDF_FS_SIZEF) (uint64, func() (*structs.FPDF_FS_SIZEF, error), error) {
	pointer, err := p.Malloc(p.CSizeStructFS_SIZEF())
	if err != nil {
		return 0, nil, err
	}

	if in != nil {
		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Width) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}

		if !p.Module.Memory().WriteFloat32Le(p.Context, uint32(pointer), in.Height) {
			p.Free(pointer)
			return 0, nil, errors.New("could not write float data to memory")
		}
	}

	return pointer, func() (*structs.FPDF_FS_SIZEF, error) {
		width, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		height, ok := p.Module.Memory().ReadFloat32Le(p.Context, uint32(pointer+(p.CSizeFloat()*1)))
		if ok {
			return nil, errors.New("could not read float data from memory")
		}

		return &structs.FPDF_FS_SIZEF{
			Width:  width,
			Height: height,
		}, nil
	}, nil
}