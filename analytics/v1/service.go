package v1

import (
	"context"
	"ecommerce/analytics/store"
	pb "ecommerce/codegen/analytics/v1/proto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type AnalyticsService struct {
	pb.UnimplementedAnalyticsServer
	analyticsStore *store.AnalyticsStore
}

func NewAnalyticsService(conn *pgx.Conn) *AnalyticsService {
	return &AnalyticsService{analyticsStore: store.NewAnalyticsStore(conn)}
}

func (s *AnalyticsService) GetTotalSales(ctx context.Context, empty *emptypb.Empty) (*pb.GetTotalSalesResponse, error) {
	log.Printf("Received GetTotalSales request")
	amount, err := s.analyticsStore.GetTotalSales(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetTotalSalesResponse{
		Amount: amount,
	}, nil
}

func (s *AnalyticsService) GetSalesByProduct(ctx context.Context, request *pb.GetProductRequest) (*pb.GetSalesByProductResponse, error) {
	log.Println("Received GetSalesByProduct Request")
	id, _ := uuid.Parse(request.ProductId)
	product, err := s.analyticsStore.GetSalesByProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.GetSalesByProductResponse{
		Quantity:  product.Quantity,
		Amount:    product.Amount,
		ProductId: product.ProductId,
	}, nil
}

func (s *AnalyticsService) GetTopCustomers(ctx context.Context, empty *emptypb.Empty) (*pb.GetTopCustomersResponse, error) {
	log.Println("Received GetTopCustomers Request")
	customers, err := s.analyticsStore.GetTopCustomers(ctx)
	if err != nil {
		return nil, err
	}

	var customerStatsList []*pb.CustomerStats
	for _, customer := range customers {
		customerStats := store.CustomerStructToProto(customer)
		customerStatsList = append(customerStatsList, customerStats)
	}
	return &pb.GetTopCustomersResponse{
		Stats: customerStatsList,
	}, nil
}
