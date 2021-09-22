package paths

import "bytes"

func (p *Path) BytesAreEqual(other *Path) (isEqual bool, err error) {
	sizeA, err := p.Size()
	if err != nil {
		return
	}
	sizeB, err := other.Size()
	if err != nil {
		return
	}
	if sizeA != sizeB {
		return false, nil
	}
	bytesA, err := p.ReadBytes()
	if err != nil {
		return
	}
	bytesB, err := other.ReadBytes()
	if err != nil {
		return
	}
	return bytes.Equal(bytesA, bytesB), nil
}

func (p *Path) MustBytesAreEqual(other *Path) bool { return must1(p.BytesAreEqual(other)).(bool) }

func (p *Path) BytesIfExistsAreEqual(other *Path) (isEqual bool, err error) {
	sizeA, err := p.SizeIfExists()
	if err != nil {
		return
	}
	sizeB, err := other.SizeIfExists()
	if err != nil || sizeA != sizeB {
		return
	}
	bytesA, err := p.ReadBytesIfExists()
	if err != nil {
		return
	}
	bytesB, err := other.ReadBytesIfExists()
	if err != nil {
		return
	}
	return bytes.Equal(bytesA, bytesB), nil
}

func (p *Path) MustBytesIfExistsAreEqual(other *Path) bool {
	return must1(p.BytesIfExistsAreEqual(other)).(bool)
}

func (p *Path) BytesAreEqualToBytes(other []byte) (isEqual bool, err error) {
	size, err := p.Size()
	if err != nil || int(size) != len(other) {
		return
	}
	b, err := p.ReadBytes()
	if err != nil {
		return
	}
	return bytes.Equal(b, other), nil
}

func (p *Path) MustBytesAreEqualToBytes(other []byte) bool {
	return must1(p.BytesAreEqualToBytes(other)).(bool)
}

func (p *Path) BytesIfExistsAreEqualToBytes(other []byte) (isEqual bool, err error) {
	b, err := p.ReadBytesIfExists()
	if err != nil {
		return
	}
	return bytes.Equal(b, other), nil
}

func (p *Path) MustBytesIfExistsAreEqualToBytes(other []byte) bool {
	return must1(p.BytesIfExistsAreEqualToBytes(other)).(bool)
}
