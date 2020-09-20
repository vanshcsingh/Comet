package sessiondb

// SessionDB holds the weights for selection policies for contextuuids
type SessionDB interface {
	// TODO: complete interface definition
	Get(string) interface{}
}
