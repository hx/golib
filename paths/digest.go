package paths

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
)

type Digest []byte

func (p *Path) Digest(hash hash.Hash) (d Digest, err error) {
	file, err := p.Open()
	if err != nil {
		return
	}
	_, err = io.Copy(hash, file)
	if err == nil {
		err = file.Close()
	} else {
		_ = file.Close()
	}
	d = make(Digest, hash.Size())
	copy(d, hash.Sum(nil))
	return
}

func (p *Path) MustDigest(hash hash.Hash) Digest { return must1(p.Digest(hash)).(Digest) }

func (p *Path) Sha1Digest() (Digest, error)   { return p.Digest(sha1.New()) }
func (p *Path) Sha224Digest() (Digest, error) { return p.Digest(sha256.New224()) }
func (p *Path) Sha256Digest() (Digest, error) { return p.Digest(sha256.New()) }
func (p *Path) Sha384Digest() (Digest, error) { return p.Digest(sha512.New384()) }
func (p *Path) Sha512Digest() (Digest, error) { return p.Digest(sha512.New()) }

func (p *Path) MustSha1Digest() Digest   { return p.MustDigest(sha1.New()) }
func (p *Path) MustSha224Digest() Digest { return p.MustDigest(sha256.New224()) }
func (p *Path) MustSha256Digest() Digest { return p.MustDigest(sha256.New()) }
func (p *Path) MustSha384Digest() Digest { return p.MustDigest(sha512.New384()) }
func (p *Path) MustSha512Digest() Digest { return p.MustDigest(sha512.New()) }

func (p *Path) DigestIfExists(hash hash.Hash) (d Digest, err error) {
	if !p.Exists() {
		return
	}
	return p.Digest(hash)
}

func (p *Path) MustDigestIfExists(hash hash.Hash) Digest {
	return must1(p.DigestIfExists(hash)).(Digest)
}

func (p *Path) Sha1DigestIfExists() (Digest, error)   { return p.DigestIfExists(sha1.New()) }
func (p *Path) Sha224DigestIfExists() (Digest, error) { return p.DigestIfExists(sha256.New224()) }
func (p *Path) Sha256DigestIfExists() (Digest, error) { return p.DigestIfExists(sha256.New()) }
func (p *Path) Sha384DigestIfExists() (Digest, error) { return p.DigestIfExists(sha512.New384()) }
func (p *Path) Sha512DigestIfExists() (Digest, error) { return p.DigestIfExists(sha512.New()) }

func (p *Path) MustSha1DigestIfExists() Digest   { return p.MustDigestIfExists(sha1.New()) }
func (p *Path) MustSha224DigestIfExists() Digest { return p.MustDigestIfExists(sha256.New224()) }
func (p *Path) MustSha256DigestIfExists() Digest { return p.MustDigestIfExists(sha256.New()) }
func (p *Path) MustSha384DigestIfExists() Digest { return p.MustDigestIfExists(sha512.New384()) }
func (p *Path) MustSha512DigestIfExists() Digest { return p.MustDigestIfExists(sha512.New()) }

func (d Digest) Hex() string { return hex.EncodeToString(d) }

func (d Digest) Equals(other Digest) bool { return bytes.Equal(d, other) }
