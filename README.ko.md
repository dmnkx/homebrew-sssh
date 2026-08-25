[English](README.md) | [한국어](README.ko.md)

# 🔑 sssh

> `~/.ssh/config`의 Host 별명을 터미널에서 고른 뒤 SSH로 접속하고, 호스트 블록을 조회·추가·수정·삭제합니다.

`sssh`는 OpenSSH를 감싼 작은 CLI입니다. 인자를 주지 않으면 호스트를 고르는 TUI가 열립니다. 실제 접속은 시스템 `ssh`를 그대로 실행합니다(`syscall.Exec`). SSH를 다시 구현하지 않습니다.

```sh
sssh
sssh prod
sssh add prod --host 10.0.0.1 --user ubuntu
```

## sssh란?

이미 `~/.ssh/config`에 있는 **Host 별명**으로 접속하고, 그 블록을 커맨드라인에서 관리합니다.

와일드카드 Host(`Host *`, `Host *.internal` 등)는 목록·접속 대상에서 빠집니다. `Host prod staging`처럼 한 블록에 별명이 여러 개면 각각 따로 고를 수 있습니다.

접속 시 실행되는 명령은 대략 다음과 같습니다.

```text
ssh <alias>
```

별명만 넘기므로 User, Port, IdentityFile, ProxyJump 등은 **OpenSSH가 config를 읽어** 적용합니다. `sssh`가 옵션을 다시 조립하지는 않습니다.

## 주요 기능

- **TUI 선택** — 인자 없이 `sssh`를 실행하고, 필터한 뒤 Enter로 접속
- **바로 접속** — `sssh <alias>` 또는 `sssh connect <alias>`
- **Config CRUD** — `list`, `add`, `edit`, `rm`, `show`로 `~/.ssh/config` 관리
- **OpenSSH 그대로** — 실제 `ssh`를 Exec하므로 기존 config가 그대로 동작
- **안전한 저장** — 덮어쓰기 전 `<경로>.bak` 백업, 디렉터리 `0700`, 파일 `0600`
- **Homebrew + GitHub Releases** — 로컬에 Go가 없어도 설치 가능

## 설치

**Homebrew**

```sh
brew tap dmnkx/sssh
brew install sssh
```

탭은 이 저장소(`github.com/dmnkx/homebrew-sssh`) `main`의 [`Formula/sssh.rb`](Formula/sssh.rb)를 사용합니다. 예전 클론에 Formula가 없으면:

```sh
brew untap dmnkx/sssh
brew tap dmnkx/sssh
brew install sssh
```

이미 탭해 두었다면:

```sh
brew update
brew install sssh
```

**소스에서 빌드** (Go 1.24+)

```sh
git clone https://github.com/dmnkx/homebrew-sssh.git
cd homebrew-sssh
go build -o sssh .
./sssh --help
```

PATH에 `ssh`가 있어야 합니다.

## 빠른 시작

```sh
# 대화형 호스트 선택
sssh

# 별명으로 접속
sssh prod

# 접속하지 않고 ssh 명령만 보기
sssh --print-cmd prod

# 선택 가능한 호스트 목록
sssh list

# Host 블록 추가 / 수정 / 조회 / 삭제
sssh add prod --host 10.0.0.1 --user ubuntu --port 22 --key ~/.ssh/id_ed25519
sssh edit prod --user ubuntu
sssh show prod
sssh rm prod
```

다른 config 파일:

```sh
sssh --config /tmp/myconfig list
```

파일이 없으면 빈 config로 취급합니다(에러로 끝내지 않음). 저장할 때는 기존 파일이 있으면 먼저 백업합니다.

## 명령

| 명령 | 하는 일 |
|------|---------|
| `sssh` | 선택 가능한 Host가 있으면 TUI로 고른 뒤 접속. Host가 없으면 에러 |
| `sssh <alias>` | 해당 별명으로 바로 접속 |
| `sssh connect <alias>` | 위와 같음 |
| `sssh list` | 선택 가능한 Host를 한 줄씩 출력 |
| `sssh add <alias>` | Host 블록 추가 (같은 별명이 있으면 실패 — `edit` 사용) |
| `sssh edit <alias>` | 이번에 넘긴 플래그만 변경. `--user ""`로 필드 삭제 가능 |
| `sssh rm <alias>` | 해당 별명 제거. 블록에 별명이 하나뿐이면 블록 전체 삭제 |
| `sssh show <alias>` | 그 Host를 SSH config 형식 텍스트로 출력 |

### `sssh add` 플래그

| 플래그 | 필수 | SSH 키워드 |
|--------|------|------------|
| `--host` | 예 | HostName (IP 또는 DNS) |
| `--user` | 아니오 | User |
| `--port` | 아니오 | Port |
| `--key` | 아니오 | IdentityFile |
| `--jump` | 아니오 | ProxyJump |

별명에 `*`, `?`, `!`가 들어가면 거부합니다.

### 공통 플래그

| 플래그 | 설명 |
|--------|------|
| `--config <경로>` | SSH config 파일. 기본값은 `~/.ssh/config` |
| `--print-cmd` | `ssh`를 실행하지 않고, 실행할 명령 문자열만 출력 |

## TUI

| 키 | 동작 |
|----|------|
| ↑ / ↓ | 커서 이동 |
| 글자 입력 | 별명·HostName·User 필터 (대소문자 무시) |
| Backspace / Delete | 필터 한 글자 삭제 |
| Enter | 현재 항목으로 접속 (맞는 항목이 없으면 아무 일도 없음) |
| Esc, Ctrl+C | 종료 (접속 안 함) |
| `q` | 필터가 비어 있을 때만 종료. 필터 입력 중이면 글자 `q` |

## 패키지 구조

```text
main.go                 진입점
internal/cli/           cobra 명령 (list, add, edit, rm, show, connect)
internal/sshcfg/        ~/.ssh/config 파싱·저장
internal/connect/       ssh 경로 찾기, 명령 출력, Exec
internal/tui/           호스트 선택 UI
```

테스트는 각 패키지의 `*_test.go`에 있습니다.

```sh
go test ./...
```

## CI / 배포

워크플로: [`.github/workflows/ci.yml`](.github/workflows/ci.yml)

1. `main`/`master` 푸시와 Pull Request → `go test`, `go vet`, `go build`, GoReleaser 설정 검사
2. `v*` 태그 푸시 → 테스트 성공 후 GoReleaser가 GitHub Release 바이너리를 올리고, `Formula/sssh.rb`가 그 에셋 URL·SHA256을 가리키도록 `main`에 커밋합니다.

## 라이선스

[MIT](LICENSE)
