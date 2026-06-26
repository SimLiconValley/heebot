package utils

import (
	"heebot/config"

	"github.com/bwmarrin/discordgo"
)

// 1. 명령어 정의 데이터
var SayCommand = &discordgo.ApplicationCommand{
	Name:        "say",
	Description: "암호가 일치하면 봇이 대신 메시지를 전송합니다.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "password",
			Description: "명령어를 실행하기 위한 비밀번호",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "content",
			Description: "봇이 대신 말할 내용",
			Required:    true,
		},
	},
}

// 2. 명령어 실행 로직
func HandleSay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	options := data.Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	inputPassword := optionMap["password"].StringValue()
	inputMessage := optionMap["content"].StringValue()

	// 암호 검증
	if inputPassword == config.AppConfig.AdminPassword {
		// 성공 메시지는 명령어를 친 사람에게만 비밀스럽게 표시 (Flags: 64)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "성공",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		// 봇이 해당 채널에 일반 채팅으로 메시지 전송
		s.ChannelMessageSend(i.ChannelID, inputMessage)
	} else {
		// 암호가 틀렸을 때도 본인에게만 에러 표시
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ 암호가 올바르지 않습니다. 권한이 없습니다.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

}
