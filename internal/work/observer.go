package work

// Observer receives deterministic lifecycle events so a UI can render progress
// without scraping output. Every callback is optional.
type Observer struct {
	UnitStarted func(Unit)
	Reviewing   func(Unit)
	Retrying    func(Unit, string)
	UnitDone    func(UnitStats)
}
