package model

import "github.com/jinzhu/gorm"

//在售房屋数据model，数据库字段格式必须为首字母大写的驼峰格式
//在结构体内可以继承gorm.Model的模型 ，gorm.Model 会自带ID, CreatedAt, UpdatedAt, DeletedAt这4个字段
type Selling struct {
	//gorm.Model
	Id         string `gorm:"varchar(64) ;primary_key ;comment: '房子id'"`
	Name       string `gorm:"varchar(64);comment:'小区名称'"`
	TotalPrice int    `gorm:"comment: '房子总价'"`
	UnitPrice  int    `gorm:"comment: '房子单价'"`
	District   string `gorm:"varchar(64); comment:'所属行政区'"`
	Region     string `gorm:"varchar(64); comment:'详细区域'"`
	Area       int    `gorm:"comment:'面积'"`
}

func SaveSellingInfo(Id string, Name string, TotalPrice int, UnitPrice int, District string, Region string, Area int, db *gorm.DB) {
	//db := common.GetDB()
	sellingInfo := Selling{}
	sellingInfo.Id = Id
	sellingInfo.Name = Name
	sellingInfo.TotalPrice = TotalPrice
	sellingInfo.UnitPrice = UnitPrice
	sellingInfo.District = District
	sellingInfo.Region = Region
	sellingInfo.Area = Area
	db.Create(&sellingInfo)
}
