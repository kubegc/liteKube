package help

var (
	DefaultChildOffset int    = 4
	DefaultTipOffset   int    = 4 //
	DefaultValueType   string = "unknown"
	DefaultSortPrint   bool   = true
)

type Distance struct {
	childOffset        int // block offset to parent
	tipOffset          int // tip offset to type
	maxKeyLength       int // max value-key  length
	maxValueTypeLength int // max value type length
}

func NewDistance() Distance {
	return Distance{
		tipOffset:          DefaultTipOffset,
		childOffset:        DefaultChildOffset,
		maxKeyLength:       0,
		maxValueTypeLength: 0,
	}
}

func (distance *Distance) UpdateTip(name, kind string) {
	if len(name) > distance.maxKeyLength {
		distance.maxKeyLength = len(name)
	}

	if len(kind) > distance.maxValueTypeLength {
		distance.maxValueTypeLength = len(kind)
	}
}

func (distance *Distance) MaxKeyLength() int {
	return distance.maxKeyLength
}

func (distance *Distance) MaxValueTypeLength() int {
	return distance.maxValueTypeLength
}

func (distance *Distance) ChildOffseth() int {
	return distance.childOffset
}

func (distance *Distance) TipOffset() int {
	return distance.tipOffset
}
