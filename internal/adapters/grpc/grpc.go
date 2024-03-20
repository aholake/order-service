package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/aholake/order-proto/golang/order"
	"github.com/aholake/order-service/internal/application/core/domain"
	"github.com/aholake/order-service/internal/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	api  ports.APIPort
	port int
	order.UnimplementedOrderServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{
		api:  api,
		port: port,
	}
}

func (a Adapter) Run() {
	var err error

	shutdown := initTracer()
	defer shutdown()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("Unable to listen on port %d, error: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	order.RegisterOrderServer(grpcServer, a)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc on port %d, error: %v", a.port, err)
	}
}

func (a Adapter) Create(ctx context.Context, orderRequest *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	tracer := otel.Tracer("demo1TracerName")

	ctx, parentSpan := tracer.Start(
		ctx,
		"parentSpanName",
		trace.WithAttributes(attribute.String("parentAttributeKey1", "parentAttributeValue1")))

	parentSpan.AddEvent("ParentSpan-Event")
	defer parentSpan.End()
	orderItems := []domain.OrderItem{}
	ctx, childSpan := tracer.Start(
		ctx,
		"childSpanName",
		trace.WithAttributes(attribute.String("childAttributeKey1", "childAttributeValue1")))
	for _, oi := range orderRequest.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: oi.ProductCode,
			UnitPrice:   oi.UnitPrice,
			Quantity:    oi.Quantity,
		})
	}
	time.Sleep(time.Second)
	childSpan.AddEvent("ChildSpan-Event")
	childSpan.End()
	newOrder := domain.NewOrder(orderRequest.UserId, orderItems)
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	res, err := a.api.PlaceOrder(timeoutCtx, newOrder)
	cancel()
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		OrderId: res.ID,
	}, nil
}

func newExporter(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("otel-otlp-go-service"),
			attribute.String("application", "otel-otlp-go-app"),
		),
	)
}

func newTraceProvider(res *resource.Resource, bsp sdktrace.SpanProcessor) *sdktrace.TracerProvider {
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	return tracerProvider
}

func initTracer() func() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	res, err := newResource(ctx)
	reportErr(err, "failed to create res")

	conn, err := grpc.DialContext(ctx, "localhost:4317", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	reportErr(err, "failed to create gRPC connection to collector")

	// Set up a trace exporter
	traceExporter, err := newExporter(ctx, conn)
	reportErr(err, "failed to create trace exporter")

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := newTraceProvider(res, batchSpanProcessor)
	otel.SetTracerProvider(tracerProvider)

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		reportErr(tracerProvider.Shutdown(ctx), "failed to shutdown TracerProvider")
		cancel()
	}
}

func reportErr(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}
