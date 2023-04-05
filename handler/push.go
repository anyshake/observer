package handler

import "log"

func HandlePush(options *HandlerOptions) {
	log.Println("推送",
		len(options.Message.Acceleration)*len(options.Message.Acceleration[0].Vertical),
		"条数据")
}
