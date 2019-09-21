package roc

type ROCPath struct {
	strs []string
	pos  int
}

func NewROCPath(strs []string) *ROCPath {
	res := &ROCPath{
		strs: strs,
	}
	return res
}

func (this *ROCPath) GetPos() int {
	return this.pos
}

func (this *ROCPath) Get(pos int) string {
	if pos < 0 || pos >= len(this.strs) {
		return ""
	}
	return this.strs[pos]
}

func (this *ROCPath) Move() string {
	if this.pos >= len(this.strs) {
		return ""
	}
	this.pos++
	return this.strs[this.pos]
}

func (this *ROCPath) Reset() {
	this.pos = 0
}
