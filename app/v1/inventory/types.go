package inventory

type Inventory struct{}

type Binding struct {
	Format string `form:"format" json:"format" xml:"format"`
}
