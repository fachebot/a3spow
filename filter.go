package a3spow

type Filter interface {
	Filter(address string) bool
}

type LongRepeatedFilter struct {
	MinLength     int
	Char          *byte
	MaxStartIndex *int
	Reverse       bool
}

func (f *LongRepeatedFilter) Filter(address string) bool {
	length := 1
	maxLength := 1

	if f.Reverse {
		runes := []rune(address)
		for i, j := 0, len(address)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		address = string(runes)
	}

	for i := 0; i < len(address)-1; i++ {
		if address[i] == address[i+1] {
			if f.Char == nil || address[i] == *f.Char {
				length++
				continue
			}
		}

		if f.MaxStartIndex == nil {
			if length >= f.MinLength {
				return true
			}
		} else {
			if (i+1)-length > *f.MaxStartIndex {
				return false
			} else if length >= f.MinLength {
				return true
			}
		}

		if length > maxLength {
			maxLength = length
		}
		length = 1
	}

	if length > maxLength {
		maxLength = length
	}

	return maxLength >= f.MinLength
}
