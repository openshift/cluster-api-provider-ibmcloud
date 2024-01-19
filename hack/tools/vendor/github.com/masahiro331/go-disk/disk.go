package disk

import (
	"io"

	"github.com/masahiro331/go-disk/gpt"
	"github.com/masahiro331/go-disk/mbr"
	"github.com/masahiro331/go-disk/types"
	"golang.org/x/xerrors"
)

type Driver interface {
	Next() (types.Partition, error)
}

func NewDriver(sr *io.SectionReader) (Driver, error) {
	m, err := mbr.NewMasterBootRecord(sr)
	if err != nil {
		return nil, xerrors.Errorf("failed to new MBR: %w", err)
	}

	g, err := gpt.NewGUIDPartitionTable(sr)
	if err != nil {
		if m.UniqueMBRDiskSignature != [4]byte{0x00, 0x00, 0x00, 0x00} {
			return m, nil
		}
		return nil, xerrors.Errorf("failed to parse GUID partition table: %w", err)
	}

	return g, nil
}