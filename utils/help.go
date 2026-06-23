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
	sb.WriteString("유저분들을 위해 지원하는 슬래시(/) 명령어 목록과 사용법입니다.\n\n")

	// 📌 기본 명령어 섹션
	sb.WriteString("⚙️ **일반 도구**\n")
	sb.WriteString("> `/help` : 현재 보고 계시는 명령어 도움말 창을 화면에 띄워줍니다.\n\n")

	// 📌 yt 명령어 섹션
	sb.WriteString("🎥 **유튜브 다운로더 (`/yt`)**\n")
	sb.WriteString("> `/yt url:<유튜브주소> [quality:<화질/포맷>] [shared:<공개여부>]`\n")
	sb.WriteString("• 단일 유튜브 영상을 고속으로 분석하고 안전하게 다운로드할 수 있는 웹 링크를 생성합니다.\n")
	sb.WriteString("• **quality (선택)**: `1080p`, `720p` (기본값), `480p` 화질 지정 혹은 `MP3 (최고음질)` 추출 포맷을 선택할 수 있습니다.\n")
	sb.WriteString("• **shared (선택)**: `공유` 설정 시 서버원 모두가 볼 수 있으며, 선택하지 않거나 `비공개`로 두면 나에게만 결과가 보입니다.\n")
	sb.WriteString("• ⚠️ 유튜브 재생목록(Playlist) 주소는 이 명령어로 직접 다운로드할 수 없습니다.\n\n")

	// 📌 mty 명령어 섹션
	sb.WriteString("🔍 **유튜브 URL 전개기 (`/mty`)**\n")
	sb.WriteString("> `/mty content:<긴 텍스트> [print_at:<출력 위치>] [count:<출력 개수>]`\n")
	sb.WriteString("• 카톡 공지사항이나 장문의 글 속에서 다른 웹 링크는 거르고 **오직 유튜브 주소만 정확히 골라냅니다.**\n")
	sb.WriteString("• 만약 주소 중에 **유튜브 재생목록(Playlist)** 주소가 포함되어 있다면, 내부 시스템(`yt-dlp`)이 백그라운드에서 작동하여 **안에 숨겨진 모든 개별 영상의 제목과 URL을 전부 파싱하여 일괄 전개**해 줍니다.\n")
	sb.WriteString("• **10개 단위 자동 분할 출력:** 디스코드 글자 수 제한을 초과하지 않도록 10개씩 깔끔하게 페이지를 나누어 가독성 높은 리스트로 출력합니다.\n")
	sb.WriteString("• **유연한 출력 위치 (`print_at`):** 결과를 채널에 '전체 공개', 나만 볼 수 있는 '나만 보기(Ephemeral)', 혹은 깔끔하게 갠톡으로 받는 'DM으로' 중 선택하여 전송받을 수 있습니다. (기본값: 나만 보기)\n")
	sb.WriteString("• **최대 개수 제한 (`count`):** 출력할 URL의 최대 개수를 직접 지정할 수 있으며, `0`을 입력하면 제한 없이 모든 링크를 전개합니다. (기본값: 0)\n")
	sb.WriteString("• 전개된 주소는 `/yt` 명령어로 개별 다운로드할 때 매우 유용하게 쓰일 수 있습니다.\n\n")

	// 📌 파일 서버 및 파기 규칙 안내
	sb.WriteString("⏳ **주의 사항 및 보관 주기**\n")
	sb.WriteString("• 생성된 모든 미디어 다운로드 링크는 서버 보안 및 저장소 관리를 위해 시스템 내부 타이머에 의해 **자동으로 영구 파기**됩니다.\n")
	sb.WriteString("• 만료되기 전이라도 다운로드 창 하단에 있는 **`🗑️ 삭제`** 버튼을 누르면 그 즉시 수동으로 디스크에서 파일을 안전하게 소멸시킬 수 있습니다.\n")

	// 3. 인터랙션 응답 (유저의 shared 선택 여부와 무관하게 도움말은 기본적으로 친 본인에게만 보이도록 비공개 플래그를 설정합니다.)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
			Flags:   discordgo.MessageFlagsEphemeral, // 도움말은 깔끔하게 본인에게만 표시
		},
	})
}
