# NotebookLM CLI (nlm)

Google NotebookLM 비공식 Go CLI 클라이언트.

## 빌드
```bash
make build    # ./nlm 바이너리 생성
make install  # $GOPATH/bin에 설치
make test     # 테스트 실행
```

## 구조
- `cmd/` - Cobra CLI 명령어 (root, auth, notebook, source, chat, artifact, generate, note, research, share, completion)
- `internal/rpc/` - batchexecute 프로토콜 (encoder, decoder, caller, method IDs)
- `internal/auth/` - 인증 (browser login via rod, storage_state.json, token extraction)
- `internal/api/` - 비즈니스 로직 API (notebooks, sources, chat, artifacts, notes, research, sharing)
- `internal/model/` - 도메인 모델
- `internal/config/` - 설정 관리 (~/.notebooklm/)
- `internal/output/` - 터미널 출력 (lipgloss table, JSON)

## RPC 프로토콜
- 엔드포인트: `https://notebooklm.google.com/_/LabsTailwindUi/data/batchexecute`
- 인증: 쿠키 + SNlM0e(CSRF) + FdrFJe(세션ID)
- 요청: `f.req=[[["method_id", "json_params", null, "generic"]]]&at=csrf&`
- 응답: `)]}'` 접두사 제거 → chunked 파싱 → wrb.fr 매칭
- Method ID는 Python 레퍼런스(teng-lin/notebooklm-py)에서 추출

## 의존성
- spf13/cobra - CLI
- go-rod/rod - 브라우저 자동화
- charmbracelet/lipgloss - 터미널 스타일링
