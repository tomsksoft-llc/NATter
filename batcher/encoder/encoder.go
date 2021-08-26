package encoder

type Encoder interface {
	Marshal(msgs [][]byte) ([]byte, error)
}
