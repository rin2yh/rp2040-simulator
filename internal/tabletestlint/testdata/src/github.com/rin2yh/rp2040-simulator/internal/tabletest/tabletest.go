package tabletest

type Case[I, W any] struct {
	Name string
	In   I
	Want W
}
