package contextpacket

// StaticReadTargets returns the frozen mode-specific read targets. It is exposed
// for the S3 CDL equivalence test; it is retired when context packets consume the
// compiled EIR mode declarations instead.
func StaticReadTargets(mode string) []string {
	return append([]string(nil), staticModeReadTargets(mode)...)
}

// StrictMode reports whether the mode is a strict single-scope mode. Exposed for
// the same equivalence test and retired the same way.
func StrictMode(mode string) bool { return strictMode(mode) }
