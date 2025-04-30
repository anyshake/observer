package seekbuf

import "bytes"

type Buffer struct {
	b bytes.Buffer
	p int
	n int
}
