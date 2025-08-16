package initialize

import (
	"fmt"

	"Taskly.com/m/global"
	"Taskly.com/m/package/cloudinary"
)

func NewCloudinary() {
	cloud_name := global.CloudinarySetting.CloudName
	api_key := global.CloudinarySetting.ApiKey
	api_secret := global.CloudinarySetting.ApiSecret
	fmt.Printf("cloudinary ,%s ,%s,%s", cloud_name, api_key, api_secret)
	Cloudinary, err := cloudinary.InitCloudinary(cloud_name, api_key, api_secret)
	if err != nil {
		fmt.Printf("Err connect to cloudinary, %w", err)
	}
	global.Cloudinary = Cloudinary
}
