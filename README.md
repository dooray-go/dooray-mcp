> 해당 기능은 비공식이며 , 커뮤니티 기여로 작성하는 도구입니다. NHN Dooray 서비스에서 제공하는 기능이 아님을 밝혀 둡니다.

# dooray_mcp

Dooray! 를 Claude 등 MCP 호환 AI 클라이언트에서 사용할 수 있도록 해 주는 **MCP (Model Context Protocol) 서버**입니다.

자연어로 메신저를 보내고, 캘린더 일정을 조회·등록·수정·삭제하고, 업무(프로젝트 포스트)를 검색·단건 조회·등록하고, 위키를 조회·작성하고, 멤버 정보를 찾을 수 있습니다.

## 주요 기능

| 분류 | 기능 | 관련 도구 |
|------|------|-----------|
| 메신저 | 다른 멤버에게 DM 전송 | `dooray_messenger` |
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
| 기타 | 현재 시각 조회 | `os` |

반복 일정은 `daily / weekly / monthly / yearly` 주기, interval, 종료일, 요일/일자 지정까지 지원합니다.

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

## 도구 레퍼런스

### `dooray_messenger`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `send` |
| to | O | 수신자의 organizationMemberId |
| message | O | 보낼 메시지 본문 |

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

댓글·파일 업로드/다운로드·페이지 이동·삭제는 이 버전에 포함하지 않습니다.

### `os`

| 파라미터 | 필수 | 설명 |
|----------|------|------|
| operation | O | `get_date_time` — 현재 로컬 시각 반환 |

## 개발

### 의존성

* [github.com/dooray-go/dooray-sdk](https://github.com/dooray-go/dooray-sdk) — Dooray OpenAPI Go 클라이언트 (v0.8.0)
* [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) — MCP 서버 SDK

### 디렉터리 구조

```
.
├── cmd/dooray-mcp/          # MCP 서버 진입점
├── internal/account/        # 멤버 조회 도구
├── internal/calendar/       # 캘린더 조회·등록·수정·삭제
├── internal/messenger/      # 메신저 DM 도구
├── internal/ostool/         # 현재 시각 도구 (`os`)
├── internal/project/        # 프로젝트·업무 도구
├── internal/wiki/           # 위키 목록·조회·등록·수정
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
