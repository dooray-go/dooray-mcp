> 해당 기능은 비공식이며 , 커뮤니티 기여로 작성하는 도구입니다. NHN Dooray 서비스에서 제공하는 기능이 아님을 밝혀 둡니다.

# dooray_mcp

Dooray! 를 Claude 등 MCP 호환 AI 클라이언트에서 사용할 수 있도록 해 주는 **MCP (Model Context Protocol) 서버**입니다.

자연어로 메신저를 보내고, 캘린더 일정을 조회·등록·수정·삭제하고, 업무(프로젝트 포스트)를 검색·단건 조회·등록하고, 위키 페이지·댓글·첨부파일을 관리하고, 멤버 정보를 찾을 수 있습니다.

## 주요 기능

| 분류 | 기능 | 관련 도구 |
|------|------|-----------|
| 메신저 | 다른 멤버에게 DM 전송 | `dooray_messenger` |
| 메신저 | 참여 중인 채널 조회·채널 생성 | `dooray_messenger_channels`, `dooray_messenger_channel` |
| 메신저 | 채널 멤버 참여·퇴장 | `dooray_messenger_channel_members` |
| 메신저 | 채널 메시지 전송·수정·삭제·답글 | `dooray_messenger_channel_message`, `dooray_messenger_channel_log` |
| 메신저 | 새 메시지 또는 기존 메시지에서 스레드 생성 | `dooray_messenger_channel_thread` |
| 캘린더 | 내 캘린더 목록 조회 | `dooray_calendar_calendars` |
| 캘린더 | 기간별 일정 조회 | `dooray_calendar_events` |
| 캘린더 | 일정 등록 (종일/반복일정 지원) | `dooray_calendar_post_event` |
| 캘린더 | 일정 수정 | `dooray_calendar_update_event` |
| 캘린더 | 일정 삭제 | `dooray_calendar_delete_event` |
| 계정 | 이름/userCode 로 멤버 검색 | `dooray_account_members` |
| 계정 | 멤버 상세정보 조회 | `dooray_account_member` |
| 프로젝트 | 참여 중인 프로젝트 조회 | `dooray_project` |
| 프로젝트 | 업무(포스트) 검색 (담당자/상태/기한 필터) | `dooray_posts` |
| 프로젝트 | 업무(포스트) 단건 조회 (본문·첨부 포함) | `dooray_post` |
| 프로젝트 | 업무(포스트) 등록 | `dooray_project_post` |
| 위키 | 접근 가능한 위키 목록 | `dooray_wikis` |
| 위키 | 페이지 목록 (한 단계) | `dooray_wiki_pages` |
| 위키 | 페이지 본문 조회 | `dooray_wiki_page` |
| 위키 | 페이지 등록 | `dooray_wiki_page_post` |
| 위키 | 페이지 제목·본문 수정 | `dooray_wiki_page_update` |
| 위키 | 페이지 이동·삭제·참조자 변경 | `dooray_wiki_page_move`, `dooray_wiki_page_delete`, `dooray_wiki_page_referrers_update` |
| 위키 | 댓글 목록·단건 조회 | `dooray_wiki_comments`, `dooray_wiki_comment` |
| 위키 | 댓글 등록·수정·삭제 | `dooray_wiki_comment_post`, `dooray_wiki_comment_update`, `dooray_wiki_comment_delete` |
| 위키 | 공유 링크 조회 | `dooray_wiki_shared_links` |
| 위키 | 위키·페이지 파일 업로드 | `dooray_wiki_file_upload`, `dooray_wiki_page_file_upload` |
| 위키 | 첨부파일 다운로드 | `dooray_wiki_attach_file_download`, `dooray_wiki_page_file_download` |
| 위키 | 페이지 첨부파일 삭제 | `dooray_wiki_page_file_delete` |
| 기타 | 현재 시각 조회 | `os` |

반복 일정은 `daily / weekly / monthly / yearly` 주기, interval, 종료일, 요일/일자 지정까지 지원합니다.

## 이렇게 활용할 수 있습니다

MCP 클라이언트는 한 번의 요청을 처리하면서 여러 Dooray 도구를 순서대로 호출할 수 있습니다. 단순 조회뿐 아니라 조회한 결과를 바탕으로 일정·업무·위키·메신저 작업을 이어서 처리할 수 있습니다.

### 아침 업무 브리핑

현재 시각을 기준으로 오늘 일정과 마감이 가까운 업무를 모아 우선순위를 정리합니다.

```text
오늘 일정과 이번 주까지 마감인 내 업무를 확인해서, 지금 해야 할 일 순서로 정리해 줘.
```

사용 도구: `os` → `dooray_calendar_calendars` → `dooray_calendar_events` → `dooray_project` → `dooray_posts`

### 회의 일정 조율과 참석 안내

안내할 멤버를 찾고 내 일정의 빈 시간을 확인한 뒤 일정을 등록합니다. 등록 결과를 확인하고 해당 멤버에게 메신저로 안내할 수도 있습니다.

```text
김Dooray를 찾아서 내일 오후 일정과 겹치지 않는 30분 회의를 만들고, 회의 시간을 DM으로 알려줘.
```

사용 도구: `dooray_account_members` → `dooray_account_member` → `dooray_calendar_events` → `dooray_calendar_post_event` → `dooray_messenger`

등록한 일정의 시간·제목·장소가 바뀌면 `dooray_calendar_update_event`, 취소할 때는 `dooray_calendar_delete_event`를 이어서 사용할 수 있습니다.

### 프로젝트 협업 채널 개설

프로젝트 참여자를 찾고 전용 채널을 만든 뒤 첫 안내 메시지를 전송합니다. 프로젝트 도중 참여자가 바뀌면 같은 채널에 멤버를 추가하거나 내보낼 수 있습니다.

```text
김Dooray와 이Dooray의 멤버 ID를 찾아 "신규 기능 출시 준비" 비공개 채널을 만들고, 출시 체크리스트를 공유하는 첫 메시지를 보내줘.
```

사용 도구: `dooray_account_members` → `dooray_messenger_channel` → `dooray_messenger_channel_message`; 참여자 변경 시 `dooray_messenger_channel_members`

### 장애·이슈 대응 대화 구성

관련 업무와 위키 문서를 확인하고 참여 중인 메신저 채널을 찾은 뒤, 대응 메시지를 게시합니다. 세부 조사는 스레드로 분리해 대화의 맥락을 유지할 수 있습니다.

```text
참여 중인 프로젝트에서 진행 중인 긴급 업무를 찾고 장애 대응 위키 페이지도 확인해서 요약해 줘. 운영 채널을 찾아 요약을 게시하고, 로그 분석 항목은 별도 스레드로 시작해 줘.
```

사용 도구: `dooray_project` → `dooray_posts` → `dooray_post` → `dooray_wikis` → `dooray_wiki_pages` → `dooray_wiki_page` → `dooray_messenger_channels` → `dooray_messenger_channel_message` → `dooray_messenger_channel_thread`

### 공지 정정과 후속 답변

채널에 보낸 메시지의 `channelId`와 `logId`를 이용해 잘못된 내용을 수정하거나 삭제하고, 특정 메시지에 답글을 남깁니다.

```text
방금 채널에 보낸 공지의 배포 시간을 오후 4시로 고쳐줘. 원문에는 변경 사유를 답글로 남겨줘.
```

사용 도구: `dooray_messenger_channel_message` → `dooray_messenger_channel_log`; 독립적인 후속 논의가 필요하면 `dooray_messenger_channel_thread`

### 업무 선별과 상세 분석

프로젝트의 업무를 담당자·상태·기한으로 필터링하고, 중요한 업무의 본문과 첨부파일 메타데이터를 확인해 진행 상황이나 위험 요소를 정리합니다.

```text
Dooray-잘쓰자 프로젝트에서 이번 주 마감인 내 업무를 찾아줘. 그중 진행 중인 업무의 본문과 첨부파일 목록을 확인해서 막힌 점을 요약해 줘.
```

사용 도구: `dooray_project` → `dooray_posts` → `dooray_post`

### 업무 등록과 담당자 알림

이름으로 멤버 ID를 찾고, 담당자·참조자를 지정한 업무를 만든 다음 생성 결과를 메신저로 알립니다.

```text
김Dooray를 담당자로 지정해서 "위키 운영 가이드 검토" 업무를 만들고, 생성된 업무 내용을 김Dooray에게 DM으로 보내줘.
```

사용 도구: `dooray_account_members` → `dooray_project` → `dooray_project_post` → `dooray_messenger`

### 사내 지식 검색과 요약

접근 가능한 위키를 찾고 페이지 계층을 탐색한 뒤 본문, 참조자, 첨부파일을 함께 읽어 필요한 내용을 요약합니다.

```text
운영 위키에서 장애 대응 절차 페이지를 찾아 본문과 첨부파일을 확인하고 체크리스트로 정리해 줘.
```

사용 도구: `dooray_wikis` → `dooray_wiki_pages` → `dooray_wiki_page` → `dooray_wiki_page_file_download` 또는 `dooray_wiki_attach_file_download`

### 회의록·가이드 작성과 자료 첨부

문서를 새 페이지로 작성하고 파일을 연결합니다. 위키에 먼저 올린 파일을 새 페이지 생성 시 연결하거나, 만들어진 페이지에 파일을 바로 추가할 수 있습니다.

```text
오늘 회의 내용을 개발 위키에 회의록 페이지로 만들고, 이 Base64 파일을 참고자료로 첨부해 줘. 마지막에 페이지를 다시 조회해서 본문과 첨부 여부를 확인해 줘.
```

사용 도구: `dooray_wiki_file_upload` → `dooray_wiki_page_post` → `dooray_wiki_page`, 또는 `dooray_wiki_page_post` → `dooray_wiki_page_file_upload` → `dooray_wiki_page`

작성한 페이지의 제목이나 본문은 `dooray_wiki_page_update`로 보완할 수 있습니다.

### 위키 리뷰와 피드백 반영

페이지 댓글을 모아 검토 의견을 정리하고, 댓글을 등록·수정·삭제하면서 리뷰 과정을 진행합니다. 페이지에 발급된 공유 링크의 상태도 함께 확인할 수 있습니다.

```text
이 위키 페이지의 댓글과 공유 링크를 확인해 줘. 미해결 의견을 요약하고 "수정 사항을 반영했습니다"라는 댓글을 남겨줘.
```

사용 도구: `dooray_wiki_comments` → `dooray_wiki_comment` → `dooray_wiki_comment_post` / `dooray_wiki_comment_update` / `dooray_wiki_comment_delete` → `dooray_wiki_shared_links`

### 위키 구조와 접근 대상 정리

문서 구조가 바뀌었을 때 페이지를 다른 부모 아래로 이동하고 참조자 목록을 교체합니다. 더 이상 필요 없는 첨부파일이나 페이지도 정리할 수 있습니다.

```text
이 페이지를 "완료된 프로젝트" 아래로 옮기고 참조자를 지정한 멤버들로 교체해 줘. 삭제하기 전에는 현재 페이지와 첨부파일을 먼저 보여줘.
```

사용 도구: `dooray_wiki_page` → `dooray_wiki_page_move` → `dooray_wiki_page_referrers_update`; 정리가 필요하면 `dooray_wiki_page_file_delete` 또는 `dooray_wiki_page_delete`

삭제·메시지 전송·일정 및 업무 등록처럼 Dooray 데이터를 변경하는 요청에는 대상과 내용을 구체적으로 적는 것이 좋습니다. 예를 들어 “방금 조회한 일정 중 ID가 …인 일정만 삭제해 줘”처럼 요청하면 다른 항목을 잘못 변경할 가능성을 줄일 수 있습니다.

## 설치하기

### Homebrew (macOS / Linux)

```bash
brew tap dooray-go/tap
brew install dooray-mcp
```

설치 후 `dooray-mcp` 명령어를 사용할 수 있습니다.

### 직접 다운로드

릴리즈 페이지에서 PC 아키텍처에 해당하는 바이너리를 다운로드 합니다.

* https://github.com/dooray-go/dooray_mcp/releases

### 소스에서 빌드

Go 1.26 이상이 필요합니다.

```bash
git clone https://github.com/dooray-go/dooray_mcp.git
cd dooray_mcp

# 현재 플랫폼용 빌드
go build -o dist/dooray-mcp ./cmd/dooray-mcp

# 또는 모든 플랫폼용 크로스 컴파일
make build-all
```

`make build-all` 실행 시 `dist/` 디렉터리에 다음 바이너리가 생성됩니다.

* `dooray.darwin.amd64`, `dooray.darwin.arm64`
* `dooray.linux.amd64`
* `dooray.windows.amd64.exe`

## Dooray 개인 인증 토큰 발급

1. Dooray! 웹에서 **개인설정 > API > 개인 인증 토큰** 메뉴로 이동합니다.
2. 새 토큰을 생성하고 복사해 둡니다. (토큰은 한 번만 노출됩니다.)
3. 아래 설정에서 `{개인토큰}` 자리에 이 값을 입력합니다.

## 설정하기

### Claude Desktop 에서 사용하기

1. MCP 를 지원하는 Claude.ai 데스크탑 애플리케이션을 설치합니다. → https://claude.ai/download
2. Claude.ai 데스크탑 애플리케이션에서 **설정 > 개발자 > 설정 편집** 을 선택합니다.

   ![img.png](img.png)

3. `claude_desktop_config.json` 을 다음과 같이 편집합니다.

```json
{
  "mcpServers": {
    "dooray": {
      "command": "dooray-mcp",
      "args": [
        "--token",
        "{개인토큰}"
      ]
    }
  },
  "globalShortcut": ""
}
```

> Homebrew 로 설치한 경우 `dooray-mcp` 를 그대로 사용합니다. 직접 다운로드한 경우 바이너리의 전체 경로를 입력하세요.
> 예: `"/Users/me/bin/dooray.darwin.arm64"`

4. Claude Desktop 을 재시작하면 Dooray 도구가 인식됩니다.

### Claude Code (CLI) 에서 사용하기

Claude Code 에 MCP 서버로 등록하면 터미널에서 바로 Dooray 기능을 사용할 수 있습니다.

1. 프로젝트 단위 등록

```bash
claude mcp add dooray -- dooray-mcp --token {개인토큰}
```

2. 전역으로 등록

```bash
claude mcp add --scope user dooray -- dooray-mcp --token {개인토큰}
```

3. 등록 확인

```bash
claude mcp list
```

4. 사용 예

```bash
claude "오늘 내 캘린더 일정을 알려줘"
```

### Codex에서 사용하기 (로컬 빌드)

저장소 루트에서 `go build -o dist/dooray-mcp ./cmd/dooray-mcp`로 빌드합니다. 같은 디렉터리의 `.dooray-token` 파일에 개인 토큰만 저장하세요. 이 파일은 Git에서 제외됩니다.

아래 `/absolute/path/dooray_mcp`를 실제 저장소 경로로 바꿔 등록합니다.

```bash
codex mcp add dooray-local -- /bin/sh -c 'token=$(cat "/absolute/path/dooray_mcp/.dooray-token") || exit 1; exec "/absolute/path/dooray_mcp/dist/dooray-mcp" --token "$token"'
codex mcp get dooray-local
```

Codex의 새 세션에서 도구를 사용합니다. 코드 변경 후 같은 경로에 다시 빌드하고 MCP 서버를 재시작하면 변경 사항이 적용됩니다.

## 사용 예시

### 메신저

```
오늘 내 일정을 중요한 순으로 정렬해서 김XX 에게 메신저로 보내 줘.
```

```
정만티에게 "회의 시작합니다" 라고 DM 보내줘.
```

### 캘린더 조회

```
내일 일정 중에 중요한 일정은 뭐야?
```

```
이번 주 금요일 오후에 비어있는 시간 알려줘.
```

### 일정 등록

```
내일 아침 9시에 자유수영 할 예정이야. 두시간 일정 등록해줘.
```

```
다음주부터 매주 월,수,금 오전 10시 스크럼 30분 일정 등록해줘. 4주간 반복.
```

```
매달 1일 오전 10시에 월간 리뷰 일정 등록, 올해 12월까지.
```

### 일정 수정 / 삭제

```
내일 아침 수영 일정 제목을 "자유수영(변경)"으로 바꿔줘.
```

```
방금 조회한 스크럼 일정만 삭제해줘.
```

![img_1.png](img_1.png)

### 업무 / 프로젝트

```
나는 정만티야 Dooray-잘쓰자 프로젝트에서 가장 급한일을 알려줘
```

```
Dooray-잘쓰자 프로젝트에서 내게 할당된 업무 중 이번 주 마감인 것만 알려줘.
```

```
지난 30일간 생성된 내 업무를 상태별로 정리해 줘.
```

```
Dooray-잘쓰자 프로젝트에 "MCP 업무 등록 테스트" 업무를 만들어줘. 본문은 "등록 기능 확인"으로 해줘.
```

```
이 업무 본문 보여줘. https://nhnent.dooray.com/task/3787724725029315943/4413565643388656467
```

### 위키

```
내가 볼 수 있는 위키 목록 보여줘.
```

```
공지사항 위키 홈 페이지 본문 읽어줘.
```

```
그 위키에 "MCP 위키 테스트" 페이지를 만들고 본문은 "등록 확인"으로 해줘.
```

```
그 페이지의 댓글을 보여줘. 확인 후 "검토 완료했습니다" 댓글을 달아줘.
```

```
그 페이지의 첨부파일 목록과 공유 링크를 보여줘.
```

## 도구 레퍼런스

### `dooray_messenger`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `send` |
| to | O | 수신자의 organizationMemberId |
| message | O | 보낼 메시지 본문 |

### 메신저 채널 조회·생성

| 도구 | operation | 설명 |
|------|-----------|------|
| `dooray_messenger_channels` | `find_channels` | 내가 참여 중인 채널 목록 조회 |
| `dooray_messenger_channel` | `create_channel` | 채널 생성 |

채널 생성 파라미터:

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| idType | O | `email` 또는 `member-id`. `memberIds`의 식별자 종류 |
| type | O | `direct` 또는 `private` |
| memberIds | O | 채널에 포함할 이메일 또는 organizationMemberId 문자열 배열 |
| capacity | X | 채널 정원. Dooray API에 문자열로 전달 |
| title | X | 채널 제목 |

`create_channel` 응답의 `result.id`가 이후 작업에 사용하는 `channelId`입니다. `find_channels`는 별도 검색·페이지 파라미터 없이 참여 중인 전체 채널을 반환합니다.

### `dooray_messenger_channel_members`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `join_members` 또는 `leave_members` |
| channelId | O | 대상 채널 ID |
| memberIds | O | 참여 또는 퇴장시킬 organizationMemberId 문자열 배열 |

### `dooray_messenger_channel_message`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `send_message` |
| channelId | O | 메시지를 보낼 채널 ID |
| text | O | 메시지 본문 |

응답의 `result.id`는 메시지의 `logId`, `result.channelId`는 채널 ID입니다. 수정·삭제·답글·기존 메시지에서 스레드 생성 시 이 두 값을 사용합니다.

### `dooray_messenger_channel_log`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `update_log`, `delete_log`, `reply_log` 중 하나 |
| channelId | O | 대상 채널 ID |
| logId | O | 대상 메시지 ID |
| text | 조건부 | 수정하거나 답글로 보낼 본문. `delete_log`에는 사용하지 않음 |

### `dooray_messenger_channel_thread`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `create_thread` 또는 `create_thread_from_log` |
| channelId | O | 대상 채널 ID |
| text | O | 새 루트 메시지 또는 기존 메시지에서 시작할 스레드 메시지 |
| threadText | X | `create_thread`에서 루트 메시지와 함께 보낼 첫 스레드 메시지 |
| logId | 조건부 | `create_thread_from_log`의 대상 메시지 ID |

스레드 생성 응답의 `result.channelId`는 요청에 사용한 채널 ID와 다를 수 있습니다. 응답의 `result.channelId`와 `result.id`를 각각 후속 작업의 `channelId`와 `logId`로 사용하세요.

```json
{
  "operation": "create_thread_from_log",
  "channelId": "3986497069711082184",
  "logId": "3986497071236383013",
  "text": "이 항목의 원인을 스레드에서 분석하겠습니다."
}
```

### `dooray_calendar_calendars`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `find_calendars` |

### `dooray_calendar_events`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `find_events` |
| calendars | X | 조회할 캘린더 ID 목록 (쉼표 구분) |
| timeMin | O | 시작 시각 (ISO 8601, 예: `2025-04-11T00:00:00+09:00`) |
| timeMax | O | 종료 시각 (ISO 8601) |

### `dooray_calendar_post_event`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `create_event` |
| calendarId | X | 등록 대상 캘린더 ID |
| subject | O | 일정 제목 |
| content | O | 일정 본문 (text/html) |
| startedAt | O | 시작 시각 (ISO 8601) |
| endedAt | O | 종료 시각 (ISO 8601) |
| wholeDayFlag | X | 종일 일정 여부. `true` 인 경우 날짜만 지정 (예: `2025-04-11+09:00`) |
| recurrenceFrequency | X | `daily` / `weekly` / `monthly` / `yearly` |
| recurrenceInterval | X | 반복 간격 (기본 1) |
| recurrenceUntil | X | 반복 종료일 (ISO 8601) |
| recurrenceByday | X | 반복 요일. 예: `MO,WE,FR`, 월간은 `1MO`, `-1FR` 등 |
| recurrenceBymonth | X | 반복 월(1-12), 연간 반복용 |
| recurrenceBymonthday | X | 반복 일(1-31), 월간/연간 반복용 |
| recurrenceTimezoneName | X | 타임존 (기본 `Asia/Seoul`) |

### `dooray_calendar_update_event`

조회한 일정의 `calendar.id`와 `id`를 사용해 지정한 필드만 수정합니다. 생략한 필드는 유지됩니다.

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `update_event` |
| calendarId | O | 대상 캘린더 ID |
| eventId | O | 대상 일정 ID (조회 결과의 `id`) |
| subject | X | 변경할 제목 |
| content | X | 변경할 본문 (text/html) |
| startedAt / endedAt | X | 변경할 시작·종료 시각 (ISO 8601) |
| wholeDayFlag | X | 종일 일정 여부 |
| location | X | 변경할 장소 |

변경할 필드를 하나 이상 지정해야 합니다. 공식 API가 빈 값이 아닌 필드만 수정하므로 빈 문자열로 제목·본문·장소를 지우는 기능은 제공하지 않습니다. 반복 규칙 변경이나 반복 일정 전체 수정 범위는 지원하지 않습니다.

```json
{
  "operation": "update_event",
  "calendarId": "cal-123",
  "eventId": "evt-456",
  "subject": "자유수영(변경)"
}
```

### `dooray_calendar_delete_event`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `delete_event` |
| calendarId | O | 대상 캘린더 ID |
| eventId | O | 대상 일정 ID (반복 일정은 조회된 회차 ID) |
| deleteType | O | `this`: 해당 일정만, `wholeFromThis`: 해당 회차와 이후 반복 일정, `whole`: 반복 일정 전체 |

일반 단건 일정은 `deleteType: "this"`로 삭제합니다. 반복 일정의 회차 ID에는 날짜 접미사가 포함될 수 있으므로 조회된 ID를 그대로 사용합니다.

```json
{
  "operation": "delete_event",
  "calendarId": "cal-123",
  "eventId": "evt-456",
  "deleteType": "this"
}
```

수정·삭제 API의 근거: [Dooray 공식 서비스 API](https://helpdesk.dooray.com/share/pages/9wWo-xwiR66BO5LGshgVTg/2939987647631384419). `dooray-sdk` v0.6.0의 `UpdateEventContext`(PUT) / `DeleteEventContext`(POST `.../events/{eventId}/delete`)를 사용합니다.

### `dooray_account_members` / `dooray_account_member`

| 도구 | operation | 설명 |
|------|-----------|------|
| `dooray_account_members` | `find_member_id` | 이름 또는 userCode 로 멤버 검색 |
| `dooray_account_member` | `find_member_details` | memberId 로 상세정보(닉네임, 이름 등) 조회 |

### `dooray_project`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `find_projects` |
| type | O | `public` / `private` |
| state | O | `active` / `archived` |
| scope | O | `private` / `public` |

### `dooray_posts`

업무(포스트) 검색 도구. `projectId` 만 필수이며, 나머지는 필터로 사용됩니다.

| 파라미터 | 설명 |
|----------|------|
| operation | `find_posts` (필수) |
| projectId | 프로젝트 ID (필수, 쉼표로 여러 개 지정 가능) |
| page / size | 페이지(기본 0), 페이지 크기(기본 20, 최대 100) |
| fromEmailAddress | 보낸 사람 이메일로 필터 |
| fromMemberIds | 작성자 memberId (쉼표 구분) |
| toMemberIds | 담당자 memberId |
| toMemberSize | 담당자 수 (0: 미지정, 1: 단일 담당자) |
| ccMemberIds | 참조자 memberId |
| tagIds | 태그 ID |
| parentPostId | 상위 업무 ID (하위 업무 조회) |
| postNumber | 업무 번호 |
| postWorkflowClasses | `backlog` / `registered` / `working` / `closed` |
| postWorkflowIds | 워크플로 ID |
| milestoneIds | 마일스톤 ID |
| subjects | 제목 키워드 |
| createdAt / updatedAt / dueAt | 날짜 필터. `today`, `thisweek`, `prev-30d`, `next-7d`, 또는 ISO8601 구간 `~` 형식 |
| order | 정렬: `postDueAt`, `postUpdatedAt`, `createdAt` (내림차순은 `-` 접두사) |

### `dooray_post`

업무 한 건의 상세(본문·첨부 포함)를 조회합니다. `dooray_posts` 결과의 `id` 또는 Dooray 업무 URL의 두 번째 ID를 `postId`로 넘깁니다.

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `get_post` |
| projectId | O | 프로젝트 ID 하나 |
| postId | O | 업무 ID 하나 |

```json
{
  "operation": "get_post",
  "projectId": "3787724725029315943",
  "postId": "4413565643388656467"
}
```

성공하면 본문(`body`)과 첨부(`files` / `fileIdList`)가 포함된 Dooray API 응답을 반환합니다.

### `dooray_project_post`

지정한 프로젝트에 새 업무를 등록합니다. 프로젝트 ID는 `dooray_project`로, 담당자·참조자의 멤버 ID는 계정 조회 도구로 확인할 수 있습니다.

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `create_post` |
| projectId | O | 등록할 프로젝트 ID 하나 |
| subject | O | 업무 제목 |
| content | O | 업무 본문 |
| mimeType | X | `text/x-markdown` (기본값) / `text/html` |
| toMemberIds | X | 담당자 organizationMemberId (쉼표 구분) |
| ccMemberIds | X | 참조자 organizationMemberId (쉼표 구분) |

```json
{
  "operation": "create_post",
  "projectId": "1234567890",
  "subject": "MCP 업무 등록 테스트",
  "content": "등록 기능 확인",
  "toMemberIds": "1111111111,2222222222"
}
```

성공하면 생성된 업무 ID를 포함한 Dooray API 응답을 반환합니다.

### `dooray_wikis`

접근 가능한 위키 목록을 조회합니다. 홈 페이지 ID는 `result[].home.pageId`입니다.

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `find_wikis` |
| page | X | 페이지 번호 (기본 0) |
| size | X | 페이지 크기 |

### `dooray_wiki_pages`

위키 페이지를 한 단계만 나열합니다. `parentPageId`를 생략하면 루트 아래 페이지를 반환합니다.

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `find_pages` |
| wikiId | O | 위키 ID (`dooray_wikis`) |
| parentPageId | X | 부모 페이지 ID |

### `dooray_wiki_page`

페이지 본문·참조자·첨부·이미지를 조회합니다.

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `get_page` |
| pageId | O | 페이지 ID |
| wikiId | X | 위키 ID. 있으면 위키 스코프 API를 사용합니다 |

```json
{
  "operation": "get_page",
  "wikiId": "100",
  "pageId": "1001"
}
```

### `dooray_wiki_page_post`

위키 페이지를 등록합니다.

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `create_page` |
| wikiId | O | 위키 ID |
| subject | O | 페이지 제목 |
| content | O | 페이지 본문 |
| parentPageId | X | 부모 페이지 ID |
| mimeType | X | `text/x-markdown` (기본값) / `text/html` |
| attachFileIds | X | 위키 파일 업로드로 얻은 첨부파일 ID의 문자열 배열 |
| referrerMemberIds | X | 참조자의 organizationMemberId 문자열 배열 |

### `dooray_wiki_page_update`

제목·본문 중 보낸 필드만 수정합니다. 둘 다 보내면 전체 수정 API를 씁니다.

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `update_page` |
| wikiId | O | 위키 ID |
| pageId | O | 페이지 ID |
| subject | X | 새 제목 |
| content | X | 새 본문 |
| mimeType | X | 본문 형식. content가 있을 때, 기본 `text/x-markdown` |

### 위키 댓글

모든 댓글 도구는 `operation`, `wikiId`, `pageId`가 필수입니다.

| 도구 | operation | 추가 파라미터 |
|------|-----------|---------------|
| `dooray_wiki_comments` | `find_comments` | `page`, `size` (선택) |
| `dooray_wiki_comment` | `get_comment` | `commentId` (필수) |
| `dooray_wiki_comment_post` | `create_comment` | `content` (필수) |
| `dooray_wiki_comment_update` | `update_comment` | `commentId`, `content` (필수) |
| `dooray_wiki_comment_delete` | `delete_comment` | `commentId` (필수) |

`page`는 0부터 시작하는 정수이고, `size`를 생략하면 Dooray API 기본값을 사용합니다. 댓글 등록·수정은 SDK의 댓글 본문 형식을 사용합니다.

```json
{
  "operation": "create_comment",
  "wikiId": "100",
  "pageId": "1001",
  "content": "검토 완료했습니다."
}
```

### `dooray_wiki_shared_links`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `find_shared_links` |
| wikiId / pageId | O | 위키·페이지 ID |
| page / size | X | 0부터 시작하는 페이지 번호 / 페이지 크기 |
| valid | X | `true`: 유효한 링크, `false`: 유효하지 않은 링크. 생략하면 필터 없이 조회 |

### 페이지 이동·삭제·참조자 변경

세 도구 모두 `operation`, `wikiId`, `pageId`가 필수입니다.

| 도구 | operation | 추가 파라미터 |
|------|-----------|---------------|
| `dooray_wiki_page_move` | `move_page` | `targetParentPageId` (필수), `targetWikiId`, `beforePageId`, `withChildren` (선택) |
| `dooray_wiki_page_delete` | `delete_page` | 없음 |
| `dooray_wiki_page_referrers_update` | `update_referrers` | `referrerMemberIds` (필수 문자열 배열) |

이동할 부모 페이지 ID는 페이지 조회 결과에서 선택합니다. 다른 위키로 옮길 때 `targetWikiId`를 지정하고, `beforePageId`로 배치 순서를 지정할 수 있습니다. `withChildren`은 하위 페이지 포함 여부이며 생략하면 API 기본 동작을 따릅니다.

참조자 변경은 기존 목록 전체를 교체합니다. `referrerMemberIds: []`로 모든 참조자를 제거할 수 있습니다.

### 위키 첨부파일

첨부파일 목록은 `dooray_wiki_page`의 `files`와 `images`에서 조회합니다. 모든 파일 도구는 `operation`, `wikiId`가 필수입니다.

| 도구 | operation | 추가 필수 파라미터 |
|------|-----------|--------------------|
| `dooray_wiki_file_upload` | `upload_wiki_file` | `filename`, `contentBase64` |
| `dooray_wiki_page_file_upload` | `upload_page_file` | `pageId`, `filename`, `contentBase64` |
| `dooray_wiki_attach_file_download` | `download_attach_file` | `attachFileId` |
| `dooray_wiki_page_file_download` | `download_page_file` | `pageId`, `fileId` |
| `dooray_wiki_page_file_delete` | `delete_page_file` | `pageId`, `fileId` |

업로드의 선택 파라미터 `fileType`은 `general` (기본값) 또는 `inline_image`입니다. `contentBase64`에는 파일 바이트를 표준 Base64로 인코딩한 문자열을 전달합니다. 서버의 로컬 파일 경로를 읽거나 쓰지 않습니다.

위키에 먼저 업로드한 파일은 반환된 첨부파일 ID를 `dooray_wiki_page_post.attachFileIds`에 넣어 새 페이지에 연결합니다. 기존 페이지에는 `dooray_wiki_page_file_upload`를 사용합니다. `attachFileId`와 페이지의 `fileId`는 서로 다른 API의 식별자이므로 해당 응답의 ID를 사용하세요.

다운로드 결과는 `contentBase64`, `contentType`, `statusCode`를 담은 JSON 텍스트입니다. 클라이언트에서 Base64를 디코딩하면 원본 파일 바이트를 얻습니다.

```json
{
  "operation": "upload_page_file",
  "wikiId": "100",
  "pageId": "1001",
  "filename": "hello.txt",
  "contentBase64": "SGVsbG8=",
  "fileType": "general"
}
```

### `os`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `get_date_time` — 현재 로컬 시각 반환 |

## 개발

### 의존성

* [github.com/dooray-go/dooray-sdk](https://github.com/dooray-go/dooray-sdk) — Dooray OpenAPI Go 클라이언트 (v0.9.0)
* [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) — MCP 서버 SDK

### 디렉터리 구조

```
.
├── cmd/dooray-mcp/          # MCP 서버 진입점
├── internal/account/        # 멤버 조회 도구
├── internal/calendar/       # 캘린더 조회·등록·수정·삭제
├── internal/messenger/      # 메신저 DM·채널·메시지·스레드 도구
├── internal/ostool/         # 현재 시각 도구 (`os`)
├── internal/project/        # 프로젝트·업무 도구
├── internal/wiki/           # 위키 페이지·댓글·첨부파일·공유 링크 도구
├── internal/mcptest/        # 테스트용 MCP 서버 헬퍼
├── Makefile
└── CHANGELOG.md
```

### 테스트

```bash
go test ./...
```

### 빌드

```bash
make build-all   # 모든 타겟(darwin/linux/windows) 빌드
make clean       # dist/ 제거
```

## 문제 해결

* **`token must be set!!` 로그 후 종료**: `--token` 인자가 지정되지 않았습니다. 설정 파일이나 CLI 인자를 확인하세요.
* **Claude Desktop 에서 도구가 보이지 않음**: 설정 편집 후 Claude Desktop 을 완전히 종료했다가 다시 실행해야 합니다.
* **일정 등록 시 시간 파싱 오류**: `startedAt`, `endedAt` 은 반드시 ISO 8601 형식이어야 합니다 (예: `2025-04-11T09:00:00+09:00`). 종일 일정은 `2025-04-11+09:00` 형태로 지정합니다.

## 변경 이력

[CHANGELOG.md](CHANGELOG.md)

## 라이선스 / 기여

이슈와 PR 은 언제나 환영합니다. → https://github.com/dooray-go/dooray_mcp
