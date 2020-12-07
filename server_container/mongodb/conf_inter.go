package mongodb

import (
	"gin-cladder/conf/elite/extend"
)

func InitMongodb() {
	extend.Service("mongodb")
}
