package models

import "time"

type OsmNode struct {
	Id        int64     `gorm:"primaryKey;autoIncrement:true;column:id" json:"id"`
	Lat       float64   `gorm:"column:lat" json:"lat"`
	Lon       float64   `gorm:"column:lon" json:"lon"`
	Tags      string    `gorm:"column:tags" json:"tags"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}

func (OsmNode) TableName() string {
	return "osm_nodes"
}
