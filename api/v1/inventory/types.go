package inventory

type Inventory struct{}

type inventoryBinding struct {
	Format string `form:"format" json:"format" xml:"format"`
}
