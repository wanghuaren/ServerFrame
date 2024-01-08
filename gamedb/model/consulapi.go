package model

import "gameutils/gameuts"

func initConsul() {
	regConfig := gameuts.DiscoveryConfig{}
	regConfig.ID = Conf.String("consul_number")
	regConfig.Name = Conf.String("micro_name")
	regConfig.Tags = []string{Conf.String("micro_name") + "_" + Conf.String("consul_number")}
	regConfig.Address = LocalIP
	regConfig.Port = Conf.Int("micro_port")

	gameuts.RegisterService(regConfig, Conf.String("consul_address"))
}
