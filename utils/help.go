package utils

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// 1. 명령어 정의 데이터
var HelpCommand = &discordgo.ApplicationCommand{
	Name:        "help",
	Description: "히봇의 사용법 및 명령어 도움말을 표시합니다.",
}

// 2. 명령어 실행 로직
func HandleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var sb strings.Builder

	sb.WriteString("🤖 **히봇 명령어 도움말 가이드** 🤖\n")

	// 📌 기본 명령어 섹션
	sb.WriteString("⚙️ **명령어**\n")
	sb.WriteString("> `/help` : 현재 보고 계시는 명령어 도움말 창을 화면에 띄워줍니다.\n\n")

	//sb.WriteString("> `/yt url:<유튜브주소> [quality:<화질/포맷>] [shared:<공개여부>]` : 유튜브를 다운로드합니다.\n")
	//sb.WriteString("• **quality (선택)**: `1080p`, `720p` (기본값), `480p`, `MP3 (최고음질)` 중 택 1\n")
	//sb.WriteString("• **shared (선택)**: `공유` 설정 시 서버원 모두가 볼 수 있으며, 선택하지 않거나 `비공개`로 두면 나에게만 결과가 보입니다.\n")
	//sb.WriteString("• ⚠️ 유튜브 재생목록(Playlist) 주소는 이 명령어로 직접 다운로드할 수 없습니다.\n\n")

	//sb.WriteString("> `/mty content:<긴 텍스트> [print_at:<출력 위치>] [count:<출력 개수>]`\n")
	//sb.WriteString("• 장문의 글 속에서 다른 웹 링크는 거르고 오직 유튜브 주소를 리스트화합니다.**\n")
	//sb.WriteString("• **유튜브 재생목록(Playlist)** 주소가 포함되어 있다면, 모든 개별 영상의 제목과 URL을 추출합니다.\n")
	//sb.WriteString("• (`print_at`): 결과를 채널에 '전체 공개', 나만 볼 수 있는 '나만 보기', 혹은 깔끔하게 갠톡으로 받는 'DM으로' 중 선택하여 전송받을 수 있습니다. (기본값: 나만 보기)\n")
	//sb.WriteString("• (`count`): 출력할 URL의 최대 개수를 직접 지정할 수 있으며, `0`을 입력하면 제한 없이 모든 링크를 전개합니다. (기본값: 0)\n")

	// 3. 인터랙션 응답 (유저의 shared 선택 여부와 무관하게 도움말은 기본적으로 친 본인에게만 보이도록 비공개 플래그를 설정합니다.)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
			Flags:   discordgo.MessageFlagsEphemeral, // 도움말은 깔끔하게 본인에게만 표시
		},
	})
}
