# Changelog

### 2026-09-18 — `feature/messenger-sdk-tools`

- `dooray-sdk`를 v0.9.0으로 올리고 메신저 채널 목록·생성, 멤버 참여·퇴장, 채널 메시지 전송·수정·삭제·답글, 스레드 생성 기능을 MCP 도구로 추가했습니다.
- 기존 DM 도구의 입력 검증과 API 오류 처리를 강화했습니다.
- README에 프로젝트 채널 개설, 장애 대응 스레드, 공지 정정 등 메신저 기능을 다른 Dooray 도구와 조합하는 활용 시나리오를 추가했습니다.
- 검증: 전체 단위 테스트·vet·race 검사·크로스 빌드와 MCP stdio `tools/call`을 통한 실제 채널 목록 조회 및 메시지 생성·수정·답글·스레드 생성·삭제를 통과했으며, 테스트 로그가 남지 않았음을 확인했습니다.

### 2026-09-16 — `feature/document-tool-use-cases`

- README에 캘린더·업무·메신저·위키 도구를 조합한 8가지 활용 시나리오를 추가했습니다.
- 각 시나리오에 바로 사용할 수 있는 자연어 요청 예시와 도구 호출 흐름을 함께 안내합니다.

### 2026-09-16 — `release/v1.5.0-beta.2`

- MCP initialize 서버 버전을 `1.5.0-beta.2`로 올렸습니다.
- Wiki 댓글·첨부파일·공유 링크·페이지 관리 도구 14개를 추가해 Wiki 도구 총 19개를 제공합니다.

### 2026-09-16 — `feature/wiki-sdk-tools`

- `dooray-sdk` v0.8.0의 기존 Wiki API를 연결해 MCP 도구 14개를 추가했습니다.
- 댓글 목록·단건 조회·등록·수정·삭제, 공유 링크 조회를 지원합니다.
- 위키·페이지 파일 업로드, 첨부파일 다운로드, 페이지 첨부파일 삭제를 지원합니다. 파일 데이터는 Base64로 전달합니다.
- 페이지 이동·삭제·참조자 목록 교체를 지원합니다.
- 페이지 생성에 `attachFileIds`, `referrerMemberIds` 옵션을 추가했습니다.
- 검증: MCP stdio `tools/call`을 통한 실제 Dooray API 호출로 신규 도구 14개를 확인했습니다. 댓글 CRUD, 공유 링크 필터, 참조자 초기화, 파일 바이트 일치, 페이지 생성 시 첨부 연결, 이동·삭제를 검증하고 테스트 데이터를 정리했습니다.

### 2026-09-16 — `release/v1.5.0-beta.1`

- MCP initialize 서버 버전을 `1.5.0-beta.1`로 올렸습니다. 위키 도구와 패키지 분리 범위가 커서 프리릴리스입니다.

### 2026-09-16 — `feature/wiki-tools`

- `dooray-sdk` v0.8.0 위키 API로 위키 목록·페이지 목록·단건 조회·등록·수정을 추가했습니다.
- 도구: `dooray_wikis`, `dooray_wiki_pages`, `dooray_wiki_page`, `dooray_wiki_page_post`, `dooray_wiki_page_update`.
- 제목만 또는 본문만 바꾸면 SDK의 title/content 전용 수정 API를 씁니다.
- 댓글·파일·이동·삭제는 포함하지 않습니다.
- 패키지를 `cmd/dooray-mcp`와 `internal/{account,calendar,messenger,ostool,project,wiki}`로 나눴습니다.
- 루트에 커밋돼 있던 `dooray_mcp` 바이너리, `.idea`, `*.iml`, `.DS_Store`를 추적에서 빼고, 빌드 산출물은 `dist/`만 쓰도록 했습니다.

### 2026-09-15 — `release/v1.4.0`

- MCP initialize 서버 버전을 `1.4.0`으로 올렸습니다.

### 2026-09-15 — `feature/get-post-tool`

- `dooray_post`로 업무 한 건을 조회합니다. 본문과 첨부 파일이 포함됩니다.
- `dooray-sdk`를 v0.7.0으로 올려 `GetPostContext`를 사용합니다.
- 프로젝트 ID·업무 ID는 각각 하나만 받으며, 검색은 기존 `dooray_posts`를 그대로 씁니다.

### 2026-09-15 — `feature/bump-mcp-server-version`

- MCP initialize에 광고하는 서버 버전을 `1.0.0`에서 `1.3.0`으로 올렸습니다. 일정 수정·삭제 도구가 포함된 다음 릴리스와 맞춥니다.

### 2026-09-15 — `feature/calendar-update-delete`

- `dooray_calendar_update_event`로 제목·본문·시간·종일 여부·장소를 지정한 필드만 수정합니다. 생략한 필드는 유지됩니다.
- `dooray_calendar_delete_event`로 일정을 삭제합니다. `deleteType`은 `this` / `wholeFromThis` / `whole`입니다. 단건도 `this`가 필요합니다.
- `dooray-sdk`를 v0.4.1에서 v0.6.0으로 올렸습니다. 일정 수정·삭제는 v0.6.0 API를 쓰고, v0.5.0의 `GetPosts` `parent.number`(int) 수정도 함께 들어갑니다.
- 빈 문자열로 제목·본문·장소를 지울 수 없고, 반복 규칙·참석자 변경은 지원하지 않습니다. 반복 회차 ID는 조회된 값을 그대로 씁니다.
- 검증: 단위 테스트, 실제 Dooray 일정 등록→수정→삭제, 하위 업무가 포함된 `dooray_posts` 조회.

### 2026-09-09 — `feature/create-project-task`

- `dooray_project_post`로 업무 제목·본문·담당자·참조자를 지정해 등록할 수 있습니다.
- SDK의 컨텍스트 지원 생성 API를 사용하고, 입력과 API 응답의 성공 여부를 검사합니다.
- 외부 호출 없는 등록 테스트, Codex 로컬 등록 안내와 토큰 파일 Git 제외 규칙을 추가했습니다.
- 검증: 전체 테스트·빌드·`go vet`, MCP 초기화·도구 목록 조회, 실제 프로젝트 조회 및 업무 등록 성공.
