package fs

import (
	"io"

	"github.com/masahiro331/go-disk/types"
	"golang.org/x/xerrors"
)

var (
	_ types.Driver    = &DirectFileSystem{}
	_ types.Partition = &DirectFileSystemPartition{}
)

type DirectFileSystem struct {
	Partition *DirectFileSystemPartition
}

func (d *DirectFileSystem) Next() (types.Partition, error) {
	if d.Partition == nil {
		return nil, io.EOF
	}
	partition := d.Partition
	d.Partition = nil

	return partition, nil
}

type DirectFileSystemPartition struct {
	sectionReader *io.SectionReader
}

func (p DirectFileSystemPartition) Name() string {
	return "0"
}

func (p DirectFileSystemPartition) GetType() []byte {
	return []byte{}
}

func (p DirectFileSystemPartition) GetStartSector() uint64 {
	return uint64(0)
}

func (p DirectFileSystemPartition) Bootable() bool {
	return false
}

func (p DirectFileSystemPartition) GetSize() uint64 {
	return uint64(p.sectionReader.Size())
}

func (p DirectFileSystemPartition) GetSectionReader() io.SectionReader {
	return *p.sectionReader
}

func (p DirectFileSystemPartition) IsSupported() bool {
	return true
}

func NewDirectFileSystem(sr *io.SectionReader) (*DirectFileSystem, error) {
	_, err := sr.Seek(0, io.SeekStart)
	if err != nil {
		return nil, xerrors.Errorf("failed to DirectFileSystem seek offset error: %w", err)
	}
	return &DirectFileSystem{
		Partition: &DirectFileSystemPartition{
			sectionReader: sr,
		},
	}, nil
}

type CheckFsFunc func(r io.Reader) bool

func CheckFileSystems(r *io.SectionReader, checkFsFuncs []CheckFsFunc) (bool, error) {
	for _, checkFsFunc := range checkFsFuncs {
		_, err := r.Seek(0, io.SeekStart)
		if err != nil {
			return false, xerrors.Errorf("failed to seek offset error: %w", err)
		}
		if checkFsFunc(r) {
			return true, nil
		}
	}
	return false, nil
}
