package metrics

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

// Prometheus Prometheus 메트릭 래퍼 / Prometheus metrics wrapper
type Prometheus struct {
	fiberPrometheus *fiberprometheus.FiberPrometheus
}

// NewPrometheus 새 Prometheus 메트릭 인스턴스 생성 / Create new Prometheus metrics instance
func NewPrometheus() *Prometheus {
	// Prometheus 설정 / Prometheus configuration
	prometheus := fiberprometheus.New("spindle")
	
	return &Prometheus{
		fiberPrometheus: prometheus,
	}
}

// Middleware Prometheus 메트릭 미들웨어 반환 / Return Prometheus metrics middleware
func (p *Prometheus) Middleware() fiber.Handler {
	return p.fiberPrometheus.Middleware
}

// RegisterAt 특정 경로에 메트릭 핸들러 등록 / Register metrics handler at specific path
func (p *Prometheus) RegisterAt(app fiber.Router, url string, handlers ...fiber.Handler) {
	p.fiberPrometheus.RegisterAt(app, url, handlers...)
}

// RegisterCustomMetrics 사용자 정의 메트릭 등록 / Register custom metrics
// 향후 비즈니스 메트릭 추가 시 사용 / Use when adding business metrics in the future
func (p *Prometheus) RegisterCustomMetrics() {
	// TODO: 사용자 정의 메트릭 등록 / Register custom metrics
	// 예시: / Examples:
	// - 사용자 생성 카운터 / User creation counter
	// - 데이터베이스 연결 풀 메트릭 / Database connection pool metrics
	// - 캐시 히트/미스 비율 / Cache hit/miss ratio
	// - 비즈니스 이벤트 메트릭 / Business event metrics
	
	// userCreationCounter := prometheus.NewCounterVec(
	//     prometheus.CounterOpts{
	//         Name: "users_created_total",
	//         Help: "Total number of users created",
	//     },
	//     []string{"status"},
	// )
	// prometheus.MustRegister(userCreationCounter)
}

// GetSubsystem 서브시스템별 메트릭 그룹 / Get metrics group by subsystem
// 큰 애플리케이션에서 메트릭을 체계적으로 관리하기 위함 / For systematic metrics management in large applications
func GetSubsystem(name string) string {
	return "spindle_" + name
}

// 메트릭 베스트 프랙티스 가이드 / Metrics best practices guide
// 1. 라벨 카디널리티 제한: 동적 값(user_id 등)을 라벨로 사용하지 않기
//    Limit label cardinality: Don't use dynamic values (like user_id) as labels
// 2. 메트릭 이름 일관성: {namespace}_{subsystem}_{metric_name}_{unit}
//    Metric naming consistency: {namespace}_{subsystem}_{metric_name}_{unit}
// 3. 적절한 메트릭 타입 선택: Counter, Gauge, Histogram, Summary
//    Choose appropriate metric types: Counter, Gauge, Histogram, Summary
// 4. 비즈니스 메트릭과 인프라 메트릭 분리
//    Separate business metrics from infrastructure metrics