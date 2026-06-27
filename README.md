# 🤖 히봇 (heebot)

히봇은 디스코드 서버를 풍요롭게 해줍니다!

## 설치 방법

`go build -o heebot`

`chmod +x heebot`

이후 `heebot` 을 실행해주세요.

## ✨ 주요 기능

### 0. help
- `/help` : 현재 보고 계시는 명령어 도움말 창을 화면에 띄워줍니다.

### 1. vc (현재 기능 꺼둠)
- `/vc <채널>` : 보이스 채널에 직접 참가합니다.(기능추가예정)

### 2. 권한 획득 명령어
* `/superuser <password>`
  * 마스터 비밀번호를 입력해 1시간 동안 관리자 권한을 획득합니다.
  * config.json의 admin_whitelist_ids 에 등록된 user_id만 사용 가능

#### 2.1 소통 명령어
* `/say <message>`
  * 봇이 입력된 대사를 현재 채널에 일반 텍스트로 대신 말합니다. (줄바꿈 사용 시 `\n` 입력)

#### 2.2 서버 종합 제어 명령어 (`/manage`)
`/manage` 명령어는 디스코드의 하위 서브 명령어(Sub-Command) 형태로 구조화되어 있습니다.

| 서브 명령어 | 인자 (Arguments) | 설명 |
| :--- | :--- | :--- |
| `create` | `<type>` `<name>` `[category_id]` | 새로운 채널(text, voice, category, news)을 생성합니다. |
| `delete` | `<channel_id>` | 특정 채널 혹은 카테고리를 완전히 삭제합니다. |
| `rename` | `<channel_id>` `<new_name>` | 지정한 채널의 이름을 변경합니다. |
| `move` | `<channel_id>` `<position>` | 채널의 정렬 순서를 변경합니다. (0이 맨 위) |
| `announce`| `<channel_id>` `<title>` `<description>`| 지정한 채널에 깔끔한 Embed 상자 형태로 공지를 발송합니다. (`\n` 개행 지원) |
| `delete_message` | `<channel_id>` `<message_id>` | 봇이 작성한 공지나 특정 메시지를 원격으로 삭제합니다. |
| `timeout` | `<user_id>` `<minutes>` | 유저를 분 단위로 뮤트합니다. **`0`을 입력하면 즉시 해제**됩니다. |
| `ban` | `<user_id>` `[reason]` | 유저를 서버에서 차단하고 최근 7일간의 메시지를 청소합니다. |
| `unban` | `<user_id>` | 차단 목록(Ban list)에서 특정 유저를 사면(해제)합니다. |
| `restart` | 없음 | 봇 프로세스를 안전하게 종료하고 재시작(Exit 0)합니다. |


### 3. 💬 기본 상호작용 (텍스트 반응)

`안녕` 이라고 쳐보세요. `그래 잘 살고 있니?`라고 즉시 대답합니다.




## 🛠️ 기술 스택 및 아키텍처

### Backend
- `Go (Golang)` 의 `github.com/bwmarrin/discordgo`

