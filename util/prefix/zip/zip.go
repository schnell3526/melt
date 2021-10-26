package zip

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/dsnet/compress/bzip2"
	"github.com/klauspost/compress/zip"
	"github.com/klauspost/compress/zstd"
	"github.com/schnell3526/melt/util"
	"github.com/ulikunitz/xz"
)

// ZipCompressionMethod Compression type
type ZipCompressionMethod uint16

// Compression methods.
// see https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT.
// Note LZMA: Disabled - because 7z isn't able to unpack ZIP+LZMA ZIP+LZMA2 archives made this way - and vice versa.
const (
	Store   ZipCompressionMethod = 0
	Deflate ZipCompressionMethod = 8
	BZIP2   ZipCompressionMethod = 12
	LZMA    ZipCompressionMethod = 14
	ZSTD    ZipCompressionMethod = 93
	XZ      ZipCompressionMethod = 95
)

// Zip provides facilities for operating ZIP archives.
// See https://pkware.cachefly.net/webdocs/casestudies/APPNOTE.TXT.
type Zip struct {
	// The compression level to use, as described
	// in the compress/flate package.
	CompressionLevel int

	// Whether to overwrite existing files; if false,
	// an error is returned if the file exists.
	OverwriteExisting bool

	// Whether to make all the directories necessary
	// to create a zip archive in the desired path.
	MkdirAll bool

	// If enabled, selective compression will only
	// compress files which are not already in a
	// compressed format; this is decided based
	// simply on file extension.
	SelectiveCompression bool

	// A single top-level folder can be implicitly
	// created by the Archive or Unarchive methods
	// if the files to be added to the archive
	// or the files to be extracted from the archive
	// do not all have a common root. This roughly
	// mimics the behavior of archival tools integrated
	// into OS file browsers which create a subfolder
	// to avoid unexpectedly littering the destination
	// folder with potentially many files, causing a
	// problematic cleanup/organization situation.
	// This feature is available for both creation
	// and extraction of archives, but may be slightly
	// inefficient with lots and lots of files,
	// especially on extraction.
	ImplicitTopLevelFolder bool

	// Strip number of leading paths. This feature is available
	// only during unpacking of the entire archive.
	StripComponents int

	// If true, errors encountered during reading
	// or writing a single file will be logged and
	// the operation will continue on remaining files.
	ContinueOnError bool

	// Compression algorithm
	// FileMethod ZipCompressionMethod
	zw   *zip.Writer
	zr   *zip.Reader
	ridx int
	//decinitialized bool
}

// Registering a global decompressor is not reentrant and may panic
func registerDecompressor(zr *zip.Reader) {
	// register zstd decompressor
	zr.RegisterDecompressor(uint16(ZSTD), func(r io.Reader) io.ReadCloser {
		zr, err := zstd.NewReader(r)
		if err != nil {
			return nil
		}
		return zr.IOReadCloser()
	})
	zr.RegisterDecompressor(uint16(BZIP2), func(r io.Reader) io.ReadCloser {
		bz2r, err := bzip2.NewReader(r, nil)
		if err != nil {
			return nil
		}
		return bz2r
	})
	zr.RegisterDecompressor(uint16(XZ), func(r io.Reader) io.ReadCloser {
		xr, err := xz.NewReader(r)
		if err != nil {
			return nil
		}
		return ioutil.NopCloser(xr)
	})
}

// unpack .zip file to destination.
func (z *Zip) Unarchive(source, destination string) error {
	if !util.Exists(destination) && z.MkdirAll {
		err := util.Mkdir(destination, 0755)
		if err != nil {
			return fmt.Errorf("preparing destination: %v", err)
		}
	}

	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("opening source file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("statting source file: %v", err)
	}

	err = z.Open(file, fileInfo.Size())
	if err != nil {
		return fmt.Errorf("opening zip archive for reading: %v", err)
	}

	// TODO
	return nil
}

func (z *Zip) Open(in io.Reader, size int64) error {
	inRdrAt, ok := in.(io.ReaderAt)
	if !ok {
		return fmt.Errorf("reader must be io.ReaderAt")
	}
	if z.zr != nil {
		return fmt.Errorf("zip archive is already open for reading")
	}
	var err error
	z.zr, err = zip.NewReader(inRdrAt, size)
	if err != nil {
		return fmt.Errorf("creating reader: %v", err)
	}
	registerDecompressor(z.zr)
	z.ridx = 0
	return nil
}
