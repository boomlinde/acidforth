package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	pcmFmtTag   = 1
	floatFmtTag = 3
)

type wave struct {
	RiffId     uint32
	RiffSize   uint32
	WaveId     uint32
	Fmt        uint32
	FmtSize    uint32
	FmtTag     uint16
	Channels   uint16
	Rate       uint32
	ByteRate   uint32
	BlockAlign uint16
	Bits       uint16
	DataId     uint32
	DataSize   uint32
}

type Sound struct {
	Rate uint32
	Data [][]float32
}

type FormatError struct {
	Field string
	Value interface{}
}

func (e FormatError) Error() string {
	return fmt.Sprintf("%s=%v", e.Field, e.Value)
}

func NewFormatError(field string, value interface{}) FormatError {
	return FormatError{Field: field, Value: value}
}

func ReadWav(f io.Reader) (*Sound, error) {
	wav := wave{}
	err := binary.Read(f, binary.LittleEndian, &wav)
	if err != nil {
		return nil, err
	}
	if wav.FmtSize != 16 {
		return nil, NewFormatError("FmtSize", wav.FmtSize)
	}
	if (wav.Bits != 8 && wav.Bits != 16 && wav.Bits != 32) || wav.Bits&0x7 != 0 {
		return nil, NewFormatError("Bits", wav.Bits)
	}
	if wav.Bits < 32 && wav.FmtTag == floatFmtTag {
		return nil, NewFormatError("Bits,FmtTag", fmt.Sprintf("%v,%v", wav.Bits, wav.FmtTag))
	}

	bytes := uint32(wav.Bits / 8)
	samples := wav.DataSize / bytes / uint32(wav.Channels)

	sound := Sound{Rate: wav.Rate}
	sound.Data = make([][]float32, wav.Channels)
	for i := range sound.Data {
		sound.Data[i] = make([]float32, samples)
	}

	if wav.FmtTag == pcmFmtTag {
		switch {
		case bytes == 1:
			for i := range sound.Data[0] {
				for j := uint16(0); j < wav.Channels; j++ {
					var read uint8
					err = binary.Read(f, binary.LittleEndian, &read)
					if err != nil {
						return nil, err
					}
					sound.Data[j][i] = float32(read)/128 - 1
				}
			}
		case bytes == 2:
			for i := range sound.Data[0] {
				for j := uint16(0); j < wav.Channels; j++ {
					var read int16
					err = binary.Read(f, binary.LittleEndian, &read)
					if err != nil {
						return nil, err
					}
					sound.Data[j][i] = float32(read) / 32768
				}
			}
		case bytes == 4:
			return nil, NewFormatError("Bits,FmtTag", fmt.Sprintf("%v,%v", wav.Bits, wav.FmtTag))
		}
	} else if wav.FmtTag == floatFmtTag {
		if bytes != 4 {
			return nil, errors.New("Unsupported format")
		}
		for i := range sound.Data[0] {
			for j := uint16(0); j < wav.Channels; j++ {
				if binary.Read(f, binary.LittleEndian, &sound.Data[j][i]) != nil {
					return nil, errors.New("Unsupported format")
				}
			}
		}
	} else {
		return nil, NewFormatError("FmtTag", wav.FmtTag)
	}

	return &sound, nil
}

func ReadWavFile(fName string) (*Sound, error) {
	f, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadWav(f)
}
