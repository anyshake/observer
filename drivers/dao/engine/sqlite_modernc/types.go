package sqlite_modernc

// This implementation uses modernc.org/sqlite driver
// Deos not support the following architectures:
// - linux/mips
// - linux/mipsle
// - linux/mips64
// - linux/mips64le
// - windows/386
// - windows/arm
type SQLite struct{}
