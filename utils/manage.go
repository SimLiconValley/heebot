package utils

import (
	"fmt"
	"heebot/config"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// 슈퍼유저 정보를 저장할 구조체 및 맵 정의
type SuperUserSession struct {
	ExpiresAt time.Time
}

var (
	sessionMutex sync.RWMutex
	superUsers   = make(map[string]SuperUserSession)
)

// 1. 명령어 정의 데이터
var SuperUserCommand = &discordgo.ApplicationCommand{
	Name:        "superuser",
	Description: "비밀번호를 입력하여 1시간 동안 슈퍼유저 권한을 획득합니다.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "password",
			Description: "슈퍼유저 인증 비밀번호",
			Required:    true,
		},
	},
}

// 💡 새롭게 추가된 /say 명령어 정의
var SayCommand = &discordgo.ApplicationCommand{
	Name:        "say",
	Description: "봇이 지정한 메시지를 채널에 말하게 합니다. (슈퍼유저 권한 필요)",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "message",
			Description: "봇이 전달할 대사 (줄바꿈은 \\n 사용)",
			Required:    true,
		},
	},
}

var ManageCommand = &discordgo.ApplicationCommand{
	Name:        "manage",
	Description: "관리자 전용 종합 제어 명령어입니다. (슈퍼유저 권한 필요)",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "새로운 채널을 생성합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "type",
					Description: "유형을 선택하세요.",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "텍스트 채널 (Text)", Value: "text"},
						{Name: "음성 채널 (Voice)", Value: "voice"},
						{Name: "카테고리 (Category)", Value: "category"},
						{Name: "뉴스 채널 (News)", Value: "news"},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "새 항목의 이름",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "category",
					Description: "소속시킬 카테고리를 선택하세요. (선택 사항)",
					Required:    false,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildCategory,
					},
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "delete",
			Description: "채널을 삭제합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel, // 💡 자동완성 팝업
					Name:        "channel",
					Description: "삭제할 채널 혹은 카테고리를 선택하세요.",
					Required:    true,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
						discordgo.ChannelTypeGuildVoice,
						discordgo.ChannelTypeGuildCategory,
						discordgo.ChannelTypeGuildNews,
					},
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "rename",
			Description: "채널/카테고리 이름을 변경합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel, // 💡 String에서 Channel로 변경
					Name:        "channel",
					Description: "이름을 바꿀 채널 혹은 카테고리를 선택하세요.",
					Required:    true,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
						discordgo.ChannelTypeGuildVoice,
						discordgo.ChannelTypeGuildCategory,
						discordgo.ChannelTypeGuildNews,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "new_name",
					Description: "새로운 이름",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "move",
			Description: "채널 순서를 변경합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "이동할 채널 혹은 카테고리를 선택하세요.",
					Required:    true,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
						discordgo.ChannelTypeGuildVoice,
						discordgo.ChannelTypeGuildCategory,
						discordgo.ChannelTypeGuildNews,
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "position",
					Description: "원하는 위치 번호 (0이 맨 위)",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "announce",
			Description: "예쁜 임베드 박스 형태로 공지를 발송합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "공지를 올릴 텍스트 또는 뉴스 채널을 선택하세요.",
					Required:    true,
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText, // 일반 텍스트 채널
						discordgo.ChannelTypeGuildNews, // 뉴스/공지 채널
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "공지 제목",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "공지 본문 내용 (줄바꿈은 \\n 사용)",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "delete_message",
			Description: "특정 채널의 메시지/공지를 삭제합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel, // 💡 String에서 Channel로 변경
					Name:        "channel",
					Description: "메시지가 있는 채널을 선택하세요.",
					Required:    true,
					// 🔒 메시지가 생성될 수 있는 채널 유형만 노출
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,  // 일반 텍스트 채널
						discordgo.ChannelTypeGuildNews,  // 뉴스/공지 채널
						discordgo.ChannelTypeGuildVoice, // 음성 채널 (음성 채널 내 채팅 기능 대응)
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message_id",
					Description: "삭제할 메시지 고유 ID (메시지 우클릭 후 ID 복사)",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "timeout",
			Description: "유저를 지정한 시간 동안 뮤트 시킵니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "처벌할 유저를 선택하거나 검색하세요.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "minutes",
					Description: "금지할 시간(분 단위, 0 입력 시 즉시 해제)",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "ban",
			Description: "유저를 서버에서 차단(밴)합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "user_id",
					Description: "차단할 유저의 고유 ID",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "차단 사유 (선택 사항)",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "unban",
			Description: "차단된 유저의 밴을 해제합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "user_id",
					Description: "밴을 해제할 유저의 고유 ID",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "restart",
			Description: "봇 서버를 재시작합니다 (인자 없음).",
		},
	},
}

// 2. /superuser 명령어 핸들러
func HandleSuperUser(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	userID := i.Member.User.ID
	username := i.Member.User.Username
	data := i.ApplicationCommandData()
	inputPassword := data.Options[0].StringValue()

	// 화이트리스트 검사
	if !config.IsAdminWhiteList(userID) {
		fmt.Printf("[%s] 🚨 경고: 화이트리스트 미등록 유저 인증 시도 거부 - 유저: %s(%s)\n", time.Now().Format("2006-01-02 15:04:05"), username, userID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ 실패",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	// 비밀번호 검사
	if inputPassword != config.AppConfig.AdminPassword {
		fmt.Printf("[%s] ⚠️ 슈퍼유저 권한 획득 실패 - 유저: %s(%s) | 사유: 비밀번호 불일치\n", time.Now().Format("2006-01-02 15:04:05"), username, userID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ 실패",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	expiration := time.Now().Add(1 * time.Hour)

	sessionMutex.Lock()
	superUsers[userID] = SuperUserSession{ExpiresAt: expiration}
	sessionMutex.Unlock()

	fmt.Printf("[%s] 🔑 슈퍼유저 권한 획득 성공 - 유저: %s(%s) | 만료시간: %s\n", time.Now().Format("2006-01-02 15:04:05"), username, userID, expiration.Format("15:04:05"))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("인증 성공! 지금부터 **1시간 동안(%s까지)** 유효합니다.", expiration.Format("15:04:05")),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// 💡 3. 새롭게 추가된 /say 명령어 핸들러
func HandleSay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	userID := i.Member.User.ID
	username := i.Member.User.Username

	// 🔒 [보안 검증] 슈퍼유저 유효성 체크
	sessionMutex.RLock()
	session, exists := superUsers[userID]
	sessionMutex.RUnlock()

	if !exists || time.Now().After(session.ExpiresAt) {
		if exists {
			sessionMutex.Lock()
			delete(superUsers, userID)
			sessionMutex.Unlock()
		}

		fmt.Printf("[%s] 🛑 권한 없는 거부 - 유저: %s(%s)가 /say 명령을 시도함.\n", time.Now().Format("2006-01-02 15:04:05"), username, userID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ 권한이 없습니다.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	data := i.ApplicationCommandData()
	messageContent := data.Options[0].StringValue()

	// 줄바꿈 매핑 처리 (\n 지원)
	messageContent = strings.ReplaceAll(messageContent, "\\n", "\n")

	// 1단계: 명령어를 입력한 유저에게는 내부 응답 성공 처리(Ephemeral 플래그로 관리자 본인에게만 메시지 전송 성공 알림 표시)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ 메시지를 성공적으로 전달했습니다.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	// 2단계: 실제로 봇이 해당 채널에 메시지 발송
	_, err := s.ChannelMessageSend(i.ChannelID, messageContent)
	if err != nil {
		fmt.Printf("[%s] ❌ /say 발송 에러: %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		return
	}

	// 3단계: 로깅
	fmt.Printf("[%s] 💬 [슈퍼유저 활성] 실행 유저: %s(%s) | /say 메시지: %s\n", time.Now().Format("2006-01-02 15:04:05"), username, userID, messageContent)
}

// 4. /manage 종합 명령어 핸들러
func HandleManage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	userID := i.Member.User.ID
	username := i.Member.User.Username

	sessionMutex.RLock()
	session, exists := superUsers[userID]
	sessionMutex.RUnlock()

	if !exists || time.Now().After(session.ExpiresAt) {
		if exists {
			sessionMutex.Lock()
			delete(superUsers, userID)
			sessionMutex.Unlock()
		}

		fmt.Printf("[%s] 🛑 권한 없는 거부 - 유저: %s(%s)가 /manage 명령을 시도함.\n", time.Now().Format("2006-01-02 15:04:05"), username, userID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ 권한이 없습니다.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	data := i.ApplicationCommandData()
	options := data.Options

	subCommand := options[0]

	subOptionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, opt := range subCommand.Options {
		subOptionMap[opt.Name] = opt
	}

	var responseMessage string
	isRestartPending := false

	switch subCommand.Name {
	case "create":
		channelTypeStr := strings.ToLower(subOptionMap["type"].StringValue())
		channelName := subOptionMap["name"].StringValue()

		var categoryID string
		// 💡 category_id -> category 로 변경 및 ChannelValue(s).ID 추출
		if opt, ok := subOptionMap["category"]; ok {
			categoryID = opt.ChannelValue(s).ID
		}

		var discordChannelType discordgo.ChannelType
		switch channelTypeStr {
		case "text":
			discordChannelType = discordgo.ChannelTypeGuildText
		case "voice":
			discordChannelType = discordgo.ChannelTypeGuildVoice
		case "category":
			discordChannelType = discordgo.ChannelTypeGuildCategory
		case "news":
			discordChannelType = discordgo.ChannelTypeGuildNews
		default:
			responseMessage = "❌ 잘못된 유형입니다. (text, voice, category, news)"
		}

		if responseMessage == "" {
			channel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
				Name:     channelName,
				Type:     discordChannelType,
				ParentID: categoryID,
			})
			if err != nil {
				responseMessage = "❌ 생성 실패: " + err.Error()
			} else {
				responseMessage = "✅ 채널이 생성되었습니다: <#" + channel.ID + ">"
			}
		}

	case "delete":
		// 💡 channel_id -> channel 로 변경 및 ChannelValue(s).ID 추출
		targetID := subOptionMap["channel"].ChannelValue(s).ID
		if _, err := s.ChannelDelete(targetID); err != nil {
			responseMessage = "❌ 삭제 실패: " + err.Error()
		} else {
			responseMessage = "✅ 채널이 완전히 삭제되었습니다."
		}

	case "rename":
		// 💡 channel_id -> channel 로 변경 및 ChannelValue(s).ID 추출
		targetID := subOptionMap["channel"].ChannelValue(s).ID
		newName := subOptionMap["new_name"].StringValue()
		_, err := s.ChannelEdit(targetID, &discordgo.ChannelEdit{Name: newName})
		if err != nil {
			responseMessage = "❌ 변경 실패: " + err.Error()
		} else {
			responseMessage = "✅ 채널명이 **" + newName + "**(으)로 변경되었습니다."
		}

	case "move":
		// 💡 channel_id -> channel 로 변경 및 ChannelValue(s).ID 추출
		targetID := subOptionMap["channel"].ChannelValue(s).ID
		pos := int(subOptionMap["position"].IntValue())
		_, err := s.ChannelEdit(targetID, &discordgo.ChannelEdit{Position: &pos})
		if err != nil {
			responseMessage = "❌ 순서 변경 실패: " + err.Error()
		} else {
			responseMessage = "✅ 채널이 " + strconv.Itoa(pos) + "번 위치로 이동했습니다."
		}

	case "announce":
		// 💡 channel_id -> channel 로 변경 및 ChannelValue(s).ID 추출
		targetChannelID := subOptionMap["channel"].ChannelValue(s).ID
		title := subOptionMap["title"].StringValue()
		description := strings.ReplaceAll(subOptionMap["description"].StringValue(), "\\n", "\n")

		embed := &discordgo.MessageEmbed{
			Title:       title,
			Description: description,
			Color:       0x00FF00,
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "히봇 공지",
			},
		}

		_, err := s.ChannelMessageSendEmbed(targetChannelID, embed)
		if err != nil {
			responseMessage = "❌ 공지 발송 실패: " + err.Error()
		} else {
			responseMessage = "✅ 지정된 채널에 임베드 공지를 발송했습니다."
		}

	case "delete_message":
		// 💡 channel_id -> channel 로 변경 및 ChannelValue(s).ID 추출 (message_id는 유지)
		targetChannelID := subOptionMap["channel"].ChannelValue(s).ID
		targetMessageID := subOptionMap["message_id"].StringValue()

		err := s.ChannelMessageDelete(targetChannelID, targetMessageID)
		if err != nil {
			responseMessage = "❌ 메시지 삭제 실패: " + err.Error()
		} else {
			responseMessage = "✅ 해당 메시지가 성공적으로 삭제되었습니다."
		}

	case "timeout":
		// 💡 user_id -> user 로 변경 및 UserValue(s).ID 추출
		targetUserID := subOptionMap["user"].UserValue(s).ID
		minutes := subOptionMap["minutes"].IntValue()

		var until *time.Time
		if minutes > 0 {
			timeoutDuration := time.Now().Add(time.Duration(minutes) * time.Minute)
			until = &timeoutDuration
		} else {
			until = nil
		}

		err := s.GuildMemberTimeout(i.GuildID, targetUserID, until)
		if err != nil {
			responseMessage = "❌ 타임아웃 처리 실패: " + err.Error()
		} else {
			if minutes > 0 {
				responseMessage = "✅ 해당 유저(ID: " + targetUserID + ")를 " + strconv.FormatInt(minutes, 10) + "분 동안 타임아웃 처리했습니다."
			} else {
				responseMessage = "✅ 해당 유저(ID: " + targetUserID + ")의 타임아웃을 즉시 해제했습니다."
			}
		}

	case "ban":
		// 💡 user_id -> user 로 변경 및 UserValue(s).ID 추출
		targetUserID := subOptionMap["user_id"].StringValue()
		reason := "관리자 수동 차단"
		if opt, ok := subOptionMap["reason"]; ok {
			reason = opt.StringValue()
		}

		err := s.GuildBanCreateWithReason(i.GuildID, targetUserID, reason, 7)
		if err != nil {
			responseMessage = "❌ 차단 실패: " + err.Error()
		} else {
			responseMessage = "✅ 유저(ID: " + targetUserID + ")를 완전히 차단했습니다. (사유: " + reason + ")"
		}

	case "unban":
		// ⚠️ unban은 기존 기획대로 String 형태의 user_id를 유지합니다.
		targetUserID := subOptionMap["user_id"].StringValue()

		err := s.GuildBanDelete(i.GuildID, targetUserID)
		if err != nil {
			responseMessage = "❌ 차단 해제 실패: " + err.Error()
		} else {
			responseMessage = "✅ 유저(ID: " + targetUserID + ")의 서버 차단(밴)을 해제했습니다."
		}

	case "restart":
		responseMessage = "♻️ 봇 서버를 재시작합니다..."
		isRestartPending = true
	}

	fmt.Printf("[%s] 🛠️ [슈퍼유저 활성] 실행 유저: %s(%s) | 명령: /manage %s\n", time.Now().Format("2006-01-02 15:04:05"), username, userID, subCommand.Name)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: responseMessage,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if isRestartPending {
		go func() {
			time.Sleep(1500 * time.Millisecond)
			os.Exit(0)
		}()
	}
}
