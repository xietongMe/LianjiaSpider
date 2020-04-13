package model

//已售房屋数据model，数据库字段格式必须为首字母大写的驼峰格式
//在结构体内可以继承gorm.Model的模型 ，gorm.Model 会自带ID, CreatedAt, UpdatedAt, DeletedAt这4个字段
type Sold struct {
	//gorm.Model
	Id         string `gorm:"varchar(64) ;primary_key ;comment: '房子id'"`
	Name       string `gorm:"varchar(64);comment:'小区名称'"`
	TotalPrice int    `gorm:"comment: '房子总价'"`
	UnitPrice  int    `gorm:"comment: '房子单价'"`
	District   string `gorm:"varchar(64); comment:'所属行政区'"`
	//Region     string `gorm:"varchar(64); comment:'具体地区'"`
	SoldYear  string `gorm:"varchar(32)  ;comment: '交易年份'"`
	SoldMonth string `gorm:"varchar(32)  ;comment: '交易月份'"`
	Area      int    `gorm:"comment:'面积'"`
}
