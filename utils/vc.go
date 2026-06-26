package utils

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// 1. 명령어 정의 데이터
var VcCommand = &discordgo.ApplicationCommand{
	Name:        "vc",
	Description: "음성 채널에 참여하거나, 나갑니다.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "채널",
			Description: "봇을 참여시킬 음성 채널 (입력하지 않으면 음성 채널에서 나갑니다)",
			Required:    false,
			ChannelTypes: []discordgo.ChannelType{
				discordgo.ChannelTypeGuildVoice,
			},
		},
	},
}

// 2. 명령어 실행 로직
func HandleVc(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	guildID := i.GuildID
	options := data.Options

	// 1. 인자(채널)가 넘어오지 않은 경우 -> 퇴장 처리
	if len(options) == 0 {
		// 봇이 현재 이 서버(Guild)의 음성 채널에 연결되어 있는지 확인
		if vc, exists := s.VoiceConnections[guildID]; exists {
			// 응답 먼저 전송
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "👋 음성 채널에서 퇴장합니다.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			// 음성 채널 연결 해제
			vc.Disconnect()
		} else {
			// 연결되어 있지 않은 경우
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ 봇이 현재 어떤 음성 채널에도 참여하고 있지 않습니다.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		return
	}

	// 2. 인자(채널)가 넘어온 경우 -> 기존대로 입장 처리
	selectedChannel := options[0].ChannelValue(s)
	targetChannelID := selectedChannel.ID

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("음성 채널<#%s> 채널에 참여합니다!", targetChannelID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	_, err := s.ChannelVoiceJoin(guildID, targetChannelID, false, false)
	if err != nil {
		log.Printf("음성 채널 참여 실패: %v", err)
	}
}
