# sssh

`~/.ssh/config`의 **Host 별명**으로 SSH에 접속하고, 호스트 블록을 조회·추가·수정·삭제하는 CLI입니다.

인자를 주지 않으면 터미널에서 호스트를 고르는 화면(TUI)이 열립니다. 실제 접속은 시스템 `ssh`를 그대로 실행합니다(`syscall.Exec`).

## 요구 사항

- Go 1.24 이상 (로컬 빌드 시)
- PATH에 `ssh` 명령이 있을 것

## 설치 / 실행

Homebrew (저장소가 tap `dmnkx/sssh`):

```bash
brew tap dmnkx/sssh
brew install sssh
```

버전 태그를 푸시하면 GitHub Actions가 바이너리 릴리스를 만들고 `Formula/sssh.rb`를 같은 저장소에 갱신합니다.

```bash
git tag v0.1.0
git push origin v0.1.0
```

소스에서 직접 빌드:

```bash
go build -o sssh .
./sssh --help
```

테스트:

```bash
go test ./...
```

## 기본 동작

| 명령 | 하는 일 |
|------|---------|
| `sssh` | 선택 가능한 Host가 있으면 TUI로 고른 뒤 접속. Host가 없으면 에러 |
| `sssh <alias>` | 해당 별명으로 바로 접속 |
| `sssh connect <alias>` | 위와 같음 (서브커맨드 형태) |

와일드카드 Host(`Host *`, `Host *.internal` 등)는 목록·접속 대상에서 빠집니다. `Host prod staging`처럼 한 블록에 별명이 여러 개면 각각 따로 고를 수 있습니다.

접속 시 실행되는 명령은 대략 다음과 같습니다.

```text
ssh <alias>
```

별명만 넘기므로 User, Port, IdentityFile, ProxyJump 등은 **OpenSSH가 `~/.ssh/config`를 읽어** 적용합니다. `sssh`가 옵션을 다시 조립하지는 않습니다.

## 공통 플래그

모든 서브커맨드에서 쓸 수 있습니다.

| 플래그 | 설명 |
|--------|------|
| `--config <경로>` | SSH config 파일. 기본값은 `~/.ssh/config` |
| `--print-cmd` | `ssh`를 실행하지 않고, 실행할 명령 문자열만 출력 |

예:

```bash
sssh --config /tmp/myconfig list
sssh --print-cmd prod
sssh connect --print-cmd prod
```

파일이 없으면 빈 config로 취급합니다(에러로 끝내지 않음). 저장할 때는 기존 파일이 있으면 `<경로>.bak` 백업을 남긴 뒤 덮어씁니다. 권한은 디렉터리 `0700`, 파일 `0600`입니다.

## 서브커맨드

### `sssh list`

선택 가능한 Host를 한 줄씩 출력합니다.

```text
별명                   user@hostname:port IdentityFile
```

- User 또는 HostName이 비어 있으면 `-`
- Port가 없으면 `22`

### `sssh add <alias>`

Host 블록을 추가합니다. 같은 별명이 이미 있으면 실패합니다(`sssh edit`를 쓰라는 메시지).

| 플래그 | 필수 | SSH 키워드 |
|--------|------|------------|
| `--host` | 예 | HostName (IP 또는 DNS) |
| `--user` | 아니오 | User |
| `--port` | 아니오 | Port |
| `--key` | 아니오 | IdentityFile |
| `--jump` | 아니오 | ProxyJump |

```bash
sssh add prod --host 10.0.0.1 --user ubuntu --port 22 --key ~/.ssh/id_ed25519
sssh add jumpbox --host 8.8.8.8 --jump bastion
```

별명에 `*`, `?`, `!`가 들어가면 거부합니다.

### `sssh edit <alias>`

있는 Host만 고칩니다. **이번에 넘긴 플래그만** 바꿉니다. 안 준 필드는 그대로입니다.

`--user ""`처럼 빈 값을 주면, 플래그를 준 것으로 보아 해당 필드를 지울 수 있습니다.

```bash
sssh edit prod --user ubuntu
sssh edit prod --host 10.0.0.2 --port 2222
```

### `sssh rm <alias>`

해당 별명을 제거합니다. 블록에 별명이 하나뿐이면 블록 전체를 지웁니다. `Host prod staging`에서 `staging`만 지우면 `Host prod`로 남습니다.

없는 별명이면 에러입니다.

### `sssh show <alias>`

그 Host를 SSH config 형식 텍스트로 출력합니다.

### `sssh connect <alias>`

`sssh <alias>`와 같이 접속합니다.

## TUI (인자 없이 실행)

`sssh`만 실행하면 목록이 나옵니다.

| 키 | 동작 |
|----|------|
| ↑ / ↓ | 커서 이동 |
| 글자 입력 | 별명·HostName·User에 대한 필터 (대소문자 무시) |
| Backspace / Delete | 필터 한 글자 삭제 |
| Enter | 현재 항목으로 접속 |
| Esc, Ctrl+C | 종료 (접속 안 함) |
| `q` | 필터가 비어 있을 때만 종료. 필터 입력 중이면 글자 `q`로 취급 |

필터와 맞는 항목이 없을 때 Enter를 눌러도 아무 일도 없습니다.

## 패키지 구조

```text
main.go                 진입점
internal/cli/           cobra 명령 (list, add, edit, rm, show, connect)
internal/sshcfg/        ~/.ssh/config 파싱·저장
internal/connect/       ssh 경로 찾기, 명령 출력, Exec
internal/tui/           호스트 선택 UI
```

테스트는 각 패키지의 `*_test.go`에 있습니다. Go는 같은 폴더의 테스트만 그 패키지로 묶습니다.

## CI / 배포

워크플로: [`.github/workflows/ci.yml`](.github/workflows/ci.yml)

1. `main`/`master` 푸시와 Pull Request → `go test`, `go vet`, `go build`, GoReleaser 설정 검사
2. `v*` 태그 푸시 → 테스트 성공 후 GoReleaser가 GitHub Release(바이너리·checksum)를 만들고 `Formula/sssh.rb`를 이 저장소에 커밋합니다. URL과 SHA256은 릴리스 에셋에서 자동으로 채웁니다.

## 대화 기록

이 저장소를 만들며 나눈 요청·결정 요약은 [docs/conversation.html](docs/conversation.html)을 브라우저로 열면 됩니다.
