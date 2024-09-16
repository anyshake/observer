package tables

import (
	"github.com/anyshake/observer/drivers/dao"
	"github.com/anyshake/observer/drivers/dao/array"
)

type AdcCount struct {
	dao.BaseTable
	Timestamp  int64            `gorm:"column:timestamp;index;not null;unique"`
	Z_Axis     array.Int32Array `gorm:"column:z_axis;type:text"`
	E_Axis     array.Int32Array `gorm:"column:e_axis;type:text"`
	N_Axis     array.Int32Array `gorm:"column:n_axis;type:text"`
	SampleRate int              `gorm:"column:sample_rate;not null"`
}

func (t AdcCount) GetModel() any {
	return &AdcCount{}
}

func (t AdcCount) GetName() string {
	return "adc_count"
}
