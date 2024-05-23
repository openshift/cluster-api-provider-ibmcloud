package mbr

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"

	"github.com/masahiro331/go-disk/types"
	"golang.org/x/xerrors"
)

const (
	SIGNATURE = 0xAA55
	Sector    = 512
)

/*
# Master Boot Record Spec
https://uefi.org/sites/default/files/resources/UEFI%20Spec%202.8B%20May%202020.pdf
p. 112
Master Boot Record always 512 bytes.
+-------------------------------+
|         Name           | Byte |
+------------------------+------+
| Bootstrap Code Area    | 440  |
| UniqueMBRDiskSignature | 4    |
| Unknown                | 2    |
| Partion 1              | 16   |
| Partion 2              | 16   |
| Partion 3              | 16   |
| Partion 4              | 16   |
| Boot Recore Sigunature | 2    |
+-------------------------------+

# Partion Spec
+-------------------+------+----------------------------------------------------------+
|        Name       | Byte |                        Description                       |
+-------------------+------+----------------------------------------------------------+
| Boot Indicator    | 1    | Boot Partion                                             |
| Staring CHS value | 3    | Starting sector of the partition in Cylinder Head Sector |
| Partition type    | 1    | FileSystem used by the partition	                      |
| Ending CHS values | 3    | Ending sector of the partition in Cylinder Head Sector   |
| Starting Sector   | 4    | Starting sector of the active partition                  |
| Partition Size    | 4    | Represents partition size in sectors                     |
+-------------------+------+----------------------------------------------------------+


ref: https://www.ijais.org/research/volume10/number8/sadi-2016-ijais-451541.pdf
*/

var InvalidSignature = xerrors.New("Invalid master boot record signature")

type MasterBootRecord struct {
	BootCodeArea           [440]byte
	UniqueMBRDiskSignature [4]byte
	Unknown                [2]byte
	Partitions             [4]Partition
	Signature              uint16

	currentPartition *Partition
	sectionReader    *io.SectionReader
}

type CHS [3]byte

type Partition struct {
	Boot     bool
	StartCHS CHS
	Type     byte
	EndCHS   CHS

	StartSector uint32
	Size        uint32
	index       int

	off           int64
	sectionReader *io.SectionReader
}

func (m *MasterBootRecord) Next() (types.Partition, error) {
	index := 0
	if m.currentPartition != nil {
		m.currentPartition.sectionReader = nil
		index = m.currentPartition.index + 1
	}
	if len(m.Partitions) <= index {
		return nil, io.EOF
	}

	m.currentPartition = &m.Partitions[index]
	offset := int64(m.currentPartition.GetStartSector()) * 512
	_, err := m.sectionReader.Seek(offset, 0)
	if err != nil {
		return nil, xerrors.Errorf("failed to seek partition(%d): %w", m.currentPartition.Index(), err)
	}
	m.currentPartition.sectionReader = io.NewSectionReader(m.sectionReader, offset, int64(m.currentPartition.GetSize()*512))

	return m.currentPartition, nil
}

func (p Partition) Index() int {
	return p.index
}

func (p Partition) Name() string {
	// TODO: add extension with type

	return strconv.Itoa(int(p.index))
}

func (p Partition) GetType() []byte {
	return []byte{p.Type}
}

func (p Partition) GetStartSector() uint64 {
	return uint64(p.StartSector)
}

func (p Partition) Bootable() bool {
	return p.Boot
}

func (p Partition) GetSize() uint64 {
	return uint64(p.Size)
}

func (p Partition) GetSectionReader() io.SectionReader {
	return *p.sectionReader
}

func NewMasterBootRecord(sr *io.SectionReader) (*MasterBootRecord, error) {
	buf := make([]byte, Sector)
	size, err := sr.Read(buf)
	if err != nil {
		return nil, xerrors.Errorf("failed to read mbr error: %w", err)
	}
	if size != Sector {
		return nil, xerrors.Errorf("binary size error: actual(%d), expected(%d)", Sector, size)
	}

	r := bytes.NewReader(buf)
	mbr := MasterBootRecord{sectionReader: sr}

	if err := binary.Read(r, binary.LittleEndian, &mbr.UniqueMBRDiskSignature); err != nil {
		return nil, xerrors.Errorf("failed to parse unique MBR disk signature: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &mbr.Unknown); err != nil {
		return nil, xerrors.Errorf("failed to parse unknown: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &mbr.BootCodeArea); err != nil {
		return nil, xerrors.Errorf("failed to parse boot code: %w", err)
	}

	for i := 0; i < len(mbr.Partitions); i++ {
		if err := binary.Read(r, binary.LittleEndian, &mbr.Partitions[i].Boot); err != nil {
			return nil, xerrors.Errorf("failed to parse partition[%d] Boot: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &mbr.Partitions[i].StartCHS); err != nil {
			return nil, xerrors.Errorf("failed to parse partition[%d] StartCHS: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &mbr.Partitions[i].Type); err != nil {
			return nil, xerrors.Errorf("failed to parse partition[%d] Type: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &mbr.Partitions[i].EndCHS); err != nil {
			return nil, xerrors.Errorf("failed to parse partition[%d] EndCHS: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &mbr.Partitions[i].StartSector); err != nil {
			return nil, xerrors.Errorf("failed to parse partition[%d] StartSector: %w", i, err)
		}
		if err := binary.Read(r, binary.LittleEndian, &mbr.Partitions[i].Size); err != nil {
			return nil, xerrors.Errorf("failed to parse partition[%d] Size: %w", i, err)
		}
		mbr.Partitions[i].index = i
	}

	if err := binary.Read(r, binary.LittleEndian, &mbr.Signature); err != nil {
		return nil, xerrors.Errorf("failed to parse signature: %w", err)
	}
	if mbr.Signature != SIGNATURE {
		return nil, InvalidSignature
	}

	for i := 0; i < len(mbr.Partitions); i++ {
		if mbr.Partitions[i].Type != 0x05 && mbr.Partitions[i].Type != 0x0f {
			continue
		}
		_, err := sr.Seek(int64(mbr.Partitions[i].StartSector)<<9, 0)
		if err != nil {
			return nil, xerrors.Errorf("failed to seek to extended boot record: %w", err)
		}
		_, err = NewMasterBootRecord(sr)
		if xerrors.Is(InvalidSignature, err) {
			mbr.Partitions[i].StartSector = mbr.Partitions[i].StartSector + 2
			mbr.Partitions[i].Size = mbr.Partitions[i].Size - 2
		} else {
			// TODO: Support Extended Master Boot Record
			return nil, xerrors.New("unsupported extended master boot record")
		}
	}
	return &mbr, nil
}

func (p Partition) IsSupported() bool {
	return true
}
