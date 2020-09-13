package database

type DbClient interface {
	InsertCoupon(coupon *Coupon) error
	GetCoupon(id string) (coupon *Coupon, err error)
	GetCoupons(filter *GetCouponsQueryFilter) (coupons []Coupon, err error)
}
