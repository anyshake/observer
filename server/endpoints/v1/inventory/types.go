package inventory

type Inventory struct{}

type request struct {
	Format string `form:"format" json:"format" xml:"format"`
}
