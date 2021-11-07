package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/fatelei/wxacode-go/pkg"
)

func main() {
	var appID string
	var appSecret string
	var scene string

	flag.StringVar(&appID, "app_id", "", "wx app id")
	flag.StringVar(&appSecret, "app_secret", "", "wx app secret")
	flag.StringVar(&scene, "scene", "", "wx app scene")
	flag.Parse()

	if len(appID) == 0 || len(appSecret) == 0 || len(scene) == 0 {
		fmt.Println("app_id or app_secret or scene is required")
		return
	}

	ctx := context.Background()
	client := wxacode.NewWxCodeClient(appID, appSecret)
	rst, err := client.GetAccessToken(ctx)
	if err != nil {
		panic(err)
	}
	rst1, err := client.GenerateQrCode(ctx, rst.AccessToken, "a")
	if err != nil {
		panic(err)
	}
	fmt.Printf("base64 image: %s\n", rst1.B64Image)
}
