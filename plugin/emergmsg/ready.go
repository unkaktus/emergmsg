package emergmsg

import "context"

// Ready implements the ready.Readiness interface, once this flips to true CoreDNS
// assumes this plugin is ready for queries; it is not checked again.
func (e *Emergmsg) Ready() bool {
	err := e.rdb.Ping(context.Background()).Err()
	return err == nil
}
