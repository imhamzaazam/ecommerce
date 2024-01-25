package store

import pb "ecommerce/codegen/analytics/v1/proto"

func CustomerStructToProto(customer Customer) *pb.CustomerStats {
	return &pb.CustomerStats{
		CustomerId:    customer.Customers.String(),
		TotalQuantity: customer.Quantity,
		TotalAmount:   customer.Amount,
	}
}
