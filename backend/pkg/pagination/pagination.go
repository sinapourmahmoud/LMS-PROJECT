package pagination

type Params struct {
	Page  int
	Limit int
}

func (p Params) Offset() int { return (p.Page - 1) * p.Limit }
