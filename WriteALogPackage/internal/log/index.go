package log

import (
	"github.com/tysonmote/gommap"
	"io"
	"os"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	entWidth        = offWidth + posWidth // entry width
)

type index struct {
	file   *os.File
	memMap gommap.MMap
	size   uint64
}

// newIndex creates an index for the given file.
// We create the index and save the current size of the file,
// so we can track the amount of data in the
// index file as we add index entries. We grow the file to the max index size before
// memory-mapping the file and then return the created index to the caller.
func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}

	idx.size = uint64(fi.Size())
	if err = os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}

	// Create memory-mapping
	if idx.memMap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}

	return idx, nil
}

// Close makes sure the memory-mapped file has synced its data to the persisted
// file and that the persisted file has flushed its contents to stable storage. Then
// it truncates the persisted file to the amount of data that’s actually in it and
// closes the file
func (i *index) Close() error {
	if err := i.memMap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}
	return i.file.Close()
}

// Read takes in an offset and returns the associated record’s position in
// the store.
func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}

	if in == -1 {
		out = uint32((i.size / entWidth) - 1)
	} else {
		out = uint32(in)
	}

	pos = uint64(out) * entWidth
	if i.size < pos+entWidth {
		return 0, 0, io.EOF
	}

	out = enc.Uint32(i.memMap[pos : pos+offWidth])
	pos = enc.Uint64(i.memMap[pos+offWidth : pos+entWidth])

	return out, pos, nil
}

// Write appends the given offset and position to the index
func (i *index) Write(off uint32, pos uint64) error {
	// validate that we have space to write the entry
	if uint64(len(i.memMap)) < i.size+entWidth {
		return io.EOF
	}

	// Encode the offset and position -> write to memory-mapped
	enc.PutUint32(i.memMap[i.size:i.size+offWidth], off)
	enc.PutUint64(i.memMap[i.size+offWidth:i.size+entWidth], pos)
	i.size += entWidth

	return nil
}

func (i *index) Name() string {
	return i.file.Name()
}
