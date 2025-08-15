# Spindle API (fiber-gorm-starter)

Go Fiber v2와 GORM으로 구축된 프로덕션 준비 완료 REST API 스켈레톤입니다. MySQL과 PostgreSQL을 모두 지원합니다.

## 주요 기능

- 🚀 **Go Fiber v2** - 빠른 HTTP 프레임워크
- 🗄️ **GORM** - MySQL/PostgreSQL 지원하는 강력한 ORM
- 🔄 **데이터베이스 마이그레이션** - golang-migrate 통합
- 📊 **관측성** - Prometheus 메트릭, 구조화된 로깅 (zap)
- 🛡️ **보안** - CORS, 보안 헤더, API 키 인증
- 🐳 **Docker** - distroless 이미지를 사용한 다단계 빌드
- 🧪 **테스팅** - 유닛 테스트 및 통합 테스트
- 📖 **API 문서** - Swagger/OpenAPI 통합
- ⚡ **성능** - k6 로드 테스팅 스크립트
- 🔧 **CI/CD** - 린팅, 테스팅, 보안 스캔이 포함된 GitHub Actions

## 빠른 시작

### 사전 요구사항

- Go 1.22+
- Docker & Docker Compose
- Make (선택사항, 편의 명령어용)

### 로컬 개발

1. **저장소 클론**
   ```bash
   git clone https://github.com/kyungseok-lee/fiber-gorm-starter.git
   cd fiber-gorm-starter
   ```

2. **환경 파일 복사**
   ```bash
   cp .env.example .env
   ```

3. **데이터베이스 시작 (MySQL)**
   ```bash
   docker-compose --profile mysql up -d
   ```

4. **데이터베이스 마이그레이션 실행**
   ```bash
   ./scripts/migrate.sh up
   ```

5. **데이터베이스 시드 (선택사항)**
   ```bash
   go run scripts/seed.go
   ```

6. **애플리케이션 시작**
   ```bash
   go run cmd/server/main.go
   ```

7. **API 접근**
   - API: http://localhost:8080
   - 헬스체크: http://localhost:8080/health
   - Swagger 문서: http://localhost:8080/docs/index.html
   - 메트릭: http://localhost:8080/metrics

### PostgreSQL 사용하기

PostgreSQL로 전환하려면:

1. **환경변수 업데이트**
   ```bash
   # .env 파일에서
   DB_DRIVER=postgres
   DB_PORT=5432
   ```

2. **PostgreSQL 시작**
   ```bash
   docker-compose --profile postgres up -d
   ```

3. **마이그레이션 실행**
   ```bash
   ./scripts/migrate.sh up
   ```

## API 엔드포인트

### 사용자
- `GET /v1/users` - 페이지네이션을 포함한 사용자 목록
- `GET /v1/users/:id` - ID로 사용자 조회
- `POST /v1/users` - 새 사용자 생성
- `PUT /v1/users/:id` - 사용자 업데이트
- `DELETE /v1/users/:id` - 사용자 삭제

### 시스템
- `GET /health` - 헬스 체크
- `GET /ready` - 준비 상태 체크 (의존성 포함)
- `GET /metrics` - Prometheus 메트릭
- `GET /docs/*` - Swagger 문서 (개발환경만)

## 데이터베이스 관리

### 마이그레이션

```bash
# 모든 마이그레이션 적용
./scripts/migrate.sh up

# 특정 개수의 마이그레이션 적용
./scripts/migrate.sh up 1

# 마이그레이션 롤백
./scripts/migrate.sh down 1

# 마이그레이션 상태 확인
./scripts/migrate.sh status

# 새 마이그레이션 생성
./scripts/migrate.sh create add_user_profile
```

### 데이터베이스 전환

애플리케이션은 MySQL과 PostgreSQL을 모두 지원합니다. 환경변수 설정으로 전환할 수 있습니다:

**MySQL:**
```env
DB_DRIVER=mysql
DB_PORT=3306
# MySQL 연결 문자열 형식은 자동으로 처리됩니다
```

**PostgreSQL:**
```env
DB_DRIVER=postgres
DB_PORT=5432
DB_SSL_MODE=disable
# PostgreSQL 연결 문자열 형식은 자동으로 처리됩니다
```

## Docker 배포

### Docker Compose 사용

**MySQL과 함께:**
```bash
docker-compose --profile mysql --profile app up -d
```

**PostgreSQL과 함께:**
```bash
docker-compose --profile postgres --profile app up -d
```

### Docker 이미지 빌드

```bash
# 이미지 빌드
docker build -t fiber-gorm-starter .

# 컨테이너 실행
docker run -p 8080:8080 --env-file .env fiber-gorm-starter
```

## 개발 도구

### Make 명령어

```bash
# 애플리케이션 실행
make run

# 테스트 실행
make test

# 린터 실행
make lint

# 바이너리 빌드
make build

# Swagger 문서 생성
make swag

# 데이터베이스 작업
make migrate-up
make migrate-down
make seed

# Docker 작업
make docker-up     # MySQL
make docker-up-pg  # PostgreSQL
make docker-down

# 코드 포맷팅
make fmt

# 모든 검사 실행 (lint + test + build)
make check
```

### 코드 생성

Swagger 문서 생성:
```bash
# swag 설치
go install github.com/swaggo/swag/cmd/swag@latest

# 문서 생성
swag init -g cmd/server/main.go -o ./docs
```

## 테스팅

### 유닛 테스트
```bash
go test -v ./...
```

### 통합 테스트
```bash
# 테스트 데이터베이스 시작
docker-compose --profile mysql up -d

# 데이터베이스와 함께 테스트 실행
ENV=test DB_NAME=fiber_gorm_starter_test go test -v ./...
```

### 로드 테스팅
```bash
# k6 설치
# https://k6.io/docs/getting-started/installation/

# 성능 테스트 실행
k6 run scripts/k6/users-smoke.js
```

## 설정

### 환경변수

| 변수 | 설명 | 기본값 |
|------|------|---------|
| `ENV` | 환경 (local/dev/prod) | `local` |
| `PORT` | 서버 포트 | `8080` |
| `DB_DRIVER` | 데이터베이스 드라이버 (mysql/postgres) | `mysql` |
| `DB_HOST` | 데이터베이스 호스트 | `localhost` |
| `DB_PORT` | 데이터베이스 포트 | `3306` |
| `DB_USER` | 데이터베이스 사용자 | `user` |
| `DB_PASS` | 데이터베이스 비밀번호 | `password` |
| `DB_NAME` | 데이터베이스 이름 | `fiber_gorm_starter` |
| `DB_SSL_MODE` | SSL 모드 (postgres 전용) | `disable` |
| `DB_MAX_OPEN` | 최대 열린 연결 수 | `25` |
| `DB_MAX_IDLE` | 최대 유휴 연결 수 | `10` |
| `DB_MAX_LIFETIME` | 연결 최대 생존 시간 | `300s` |
| `API_KEY` | 인증용 API 키 | `` |
| `LOG_LEVEL` | 로깅 레벨 | `info` |
| `METRICS_ENABLED` | Prometheus 메트릭 활성화 | `true` |
| `PPROF_ENABLED` | pprof 엔드포인트 활성화 | `false` |

### 데이터베이스 연결 풀

워크로드에 따라 연결 풀 설정을 최적화하세요:

```env
# 고트래픽 애플리케이션용
DB_MAX_OPEN=100
DB_MAX_IDLE=20
DB_MAX_LIFETIME=600s

# 저트래픽 애플리케이션용
DB_MAX_OPEN=10
DB_MAX_IDLE=5
DB_MAX_LIFETIME=300s
```

## 모니터링 & 관측성

### 메트릭

Prometheus 메트릭은 `/metrics`에서 확인할 수 있습니다:
- HTTP 요청 지속시간 및 횟수
- HTTP 요청 크기
- 데이터베이스 연결 풀 통계
- Go 런타임 메트릭

### 헬스 체크

- `/health` - 기본 애플리케이션 상태
- `/ready` - 데이터베이스 연결성을 포함한 준비 상태 체크

### 로깅

zap을 사용한 구조화된 로깅:
- 상관관계 ID를 포함한 요청/응답 로깅
- 환경별 다른 로그 레벨
- PII 안전 로깅 방식

### 프로파일링

디버깅을 위한 pprof 활성화 (개발환경에서만):
```env
PPROF_ENABLED=true
```

`/debug/pprof/`에서 프로파일링 엔드포인트에 접근

## 보안

### 보안 헤더
- HSTS with preload
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- X-XSS-Protection: 1; mode=block

### CORS 정책
- 개발환경: 모든 오리진 허용
- 프로덕션: 특정 도메인으로 제한 (미들웨어에서 설정)

### API 인증
간단한 API 키 인증 (JWT로 확장 가능):
```bash
curl -H "Authorization: Bearer your-api-key" http://localhost:8080/v1/users
```

## 성능 최적화

### 연결 풀 튜닝
- `db_connections_open` 메트릭 모니터링
- 동시 로드에 따라 `DB_MAX_OPEN` 조정
- 연결 재활용을 위한 적절한 `DB_MAX_LIFETIME` 설정

### 로그 샘플링
고트래픽 애플리케이션에서는 로그 샘플링 구현:
```go
// 예시: 프로덕션에서 요청의 10%만 샘플링
if env == "prod" && rand.Float64() > 0.1 {
    return // 로깅 건너뛰기
}
```

### 캐싱 전략
캐싱을 위한 Redis 추가 (인프라 준비됨):
```bash
docker-compose --profile redis --profile mysql --profile app up -d
```

## 배포 전략

### 블루-그린 배포
1. 현재 버전과 함께 새 버전 배포
2. 로드 밸런서를 사용하여 트래픽 전환
3. 메트릭 모니터링 후 필요시 롤백

### 카나리 배포
1. 새 버전으로 소량의 트래픽 라우팅
2. 메트릭이 정상이면 점진적으로 트래픽 증가
3. 결과에 따라 전체 배포 또는 롤백

### 롤링 업데이트
1. 인스턴스를 하나씩 업데이트
2. 진행하기 전에 헬스 체크 대기
3. 무중단 배포 유지

## 향후 개선사항

### 인증 및 권한 부여
- [ ] JWT 토큰 기반 인증
- [ ] 역할 기반 접근 제어 (RBAC)
- [ ] OAuth2/OIDC 통합
- [ ] 사용자별 속도 제한

### 캐싱
- [ ] Redis 통합
- [ ] 캐시 사이드 패턴 구현
- [ ] 캐시 무효화 전략

### 데이터베이스
- [ ] 읽기 전용 복제본 지원
- [ ] 데이터베이스 샤딩
- [ ] 쿼리 최적화 및 인덱싱
- [ ] 감사 로그를 포함한 소프트 삭제

### 관측성
- [ ] OpenTelemetry 분산 추적
- [ ] 사용자 정의 비즈니스 메트릭
- [ ] 로그 집계 (ELK 스택)
- [ ] APM 통합

### 테스팅
- [ ] Pact를 사용한 계약 테스팅
- [ ] 카오스 엔지니어링 테스트
- [ ] 성능 회귀 테스트
- [ ] 보안 침투 테스트

### DevOps
- [ ] Kubernetes 매니페스트
- [ ] Helm 차트
- [ ] ArgoCD GitOps
- [ ] Infrastructure as Code (Terraform)

## 기여하기

1. 저장소 포크
2. 기능 브랜치 생성: `git checkout -b feature/new-feature`
3. 변경사항 작성 및 테스트 추가
4. 린팅 및 테스트 실행: `make check`
5. 규칙적인 커밋 형식으로 커밋
6. 푸시 후 Pull Request 생성

## 라이선스

이 프로젝트는 MIT 라이선스 하에 라이선스가 부여됩니다 - 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

---

**Spindle** - 프로덕션 준비 완료 Go API 스타터 템플릿