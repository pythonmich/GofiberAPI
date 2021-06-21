package database

type QueryInterface interface {

}

var _ QueryInterface = (*Queries)(nil)