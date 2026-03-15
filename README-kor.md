# nlm - NotebookLM CLI

[Google NotebookLM](https://notebooklm.google.com)의 비공식 CLI 클라이언트 (Go).

터미널에서 노트북 관리, 소스 추가, AI 채팅, 콘텐츠 생성을 수행합니다.

> [English](README.md)

## 왜 만들었나?

기존 Python 클라이언트([notebooklm-py](https://github.com/teng-lin/notebooklm-py))가 있지만, Python을 Windows, macOS, Linux에 배포하려면 인터프리터, 가상 환경, 의존성 관리가 필요합니다.

**nlm**은 단일 정적 바이너리입니다. 다운로드해서 바로 실행 — 런타임 의존성 없이 CI/CD, 자동화 스크립트, 크로스 플랫폼 배포에 적합합니다.

## 설치

### 빌드된 바이너리

[Releases](https://github.com/jmk/notebooklm-cli/releases)에서 플랫폼별 바이너리를 다운로드합니다.

### 소스에서 빌드

```bash
git clone https://github.com/jmk/notebooklm-cli.git
cd notebooklm-cli
make build

go install github.com/jmk/notebooklm-cli@latest
```

### 크로스 컴파일

```bash
GOOS=windows GOARCH=amd64 go build -o nlm.exe .
GOOS=linux   GOARCH=amd64 go build -o nlm-linux .
GOOS=darwin  GOARCH=arm64 go build -o nlm-darwin .
```

## 빠른 시작

```bash
nlm auth login --reuse          # 1. 인증 (Chrome 로그인 재사용)
nlm notebook list               # 2. 노트북 목록
nlm use <notebook-id>           # 3. 노트북 선택
nlm chat ask "핵심 내용을 요약해줘"  # 4. AI에게 질문
```

## 인증

### 방법 1: Chrome 쿠키 재사용 (추천)

Chrome의 로컬 쿠키 DB에서 직접 추출합니다. 브라우저를 열지 않고, 기존 로그인에 영향 없습니다.

```bash
nlm auth login --reuse
```

> macOS 키체인 접근 허용 팝업이 뜨면 "허용"을 선택하세요.

### 방법 2: 새 브라우저 로그인

```bash
nlm auth login
```

### 인증 확인 / 삭제

```bash
nlm auth status
nlm auth clear
```

## 명령어

### 노트북

```bash
nlm notebook list                        # 목록
nlm notebook create "연구 노트"           # 생성
nlm notebook rename <id> "새 이름"        # 이름 변경
nlm notebook delete <id>                 # 삭제
nlm nb ls                                # 별칭
```

### 활성 노트북

```bash
nlm use <notebook-id>                    # 이후 명령의 기본 노트북
```

### 소스

```bash
nlm source list                          # 목록
nlm source add https://example.com       # URL 추가
nlm source add ./paper.pdf               # 파일 업로드
nlm source add <url> --wait              # 처리 완료 대기
nlm source get <id>                      # 상세 정보
nlm source refresh <id>                  # 새로고침
nlm source delete <id>                   # 삭제
```

### AI 채팅

```bash
nlm chat ask "이 문서의 핵심 논점은?"
nlm chat ask "3장만 참조해줘" -s <source-id>
nlm chat history
```

### 콘텐츠 생성

```bash
nlm generate audio                       # 오디오
nlm generate audio -i "한국어로"          # 지침 포함
nlm generate report                      # 보고서
nlm generate quiz                        # 퀴즈
nlm generate video                       # 비디오
nlm generate mind-map                    # 마인드맵
nlm generate infographic                 # 인포그래픽
nlm generate slide-deck                  # 슬라이드
nlm gen audio --wait                     # 완료 대기
```

### 아티팩트 / 노트

```bash
nlm artifact list                        # 아티팩트 목록
nlm artifact export <id> output.md       # 파일 저장
nlm note list                            # 노트 목록
nlm note create "제목" "내용"             # 노트 생성
```

### 딥 리서치 / 공유

```bash
nlm research start "X에 대해 조사해줘"
nlm research poll <id>
nlm share status
nlm share set viewer|editor|none
```

## 전역 옵션

| 플래그 | 설명 |
|--------|------|
| `--json` | JSON 출력 |
| `-n, --notebook <id>` | 노트북 ID 지정 |
| `-v, --verbose` | 상세 로그 |

## 설정

`~/.notebooklm/`에 저장됩니다. 환경 변수 `NOTEBOOKLM_AUTH_JSON`으로 인증 정보를 직접 전달할 수도 있습니다.

## 요구 사항

- **macOS**: 전체 지원 (Chrome 쿠키 추출은 키체인 사용)
- **Windows/Linux**: 브라우저 로그인 또는 수동 쿠키 임포트
- Google Chrome + Google 계정 로그인
- Go 1.21+ (빌드 시만)

## 레퍼런스

- **[notebooklm-py](https://github.com/teng-lin/notebooklm-py)** — Python 클라이언트. RPC method ID, 인코딩, 인증 흐름의 기반.
- **[Google batchexecute](https://kovatch.medium.com/deciphering-google-batchexecute-74991e4e446c)** — Google 내부 RPC 프로토콜 문서.

## 기여

1. Fork → 기능 브랜치 생성
2. 테스트 작성 (`make e2e`)
3. `go vet` / `go build` 통과 확인
4. Pull Request 제출

```bash
make build       # 빌드
make e2e-basic   # CLI 테스트 (인증 불필요)
make e2e         # 전체 E2E (인증 필요)
```

## 라이선스

MIT License — [LICENSE](LICENSE) 참조
