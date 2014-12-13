package core

type Output interface {
	Process(results []*Result)
}
