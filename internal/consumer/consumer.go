package consumer

type Consumer interface {
	Start(string) error
}
