package middleware

const (
	USER_ID      = "user_id"
	BEARER_TOKEN = "bearer_token"
)

const (

	// USER //

	SCOPE_USER_ALL           = "users.*"
	SCOPE_USER_CREATE        = "users.create"
	SCOPE_USER_LIST          = "users.list"
	SCOPE_USER_DETAIL        = "users.detail"
	SCOPE_USER_UPDATE        = "users.update"
	SCOPE_USER_DELETE        = "users.delete"
	SCOPE_USER_DOWNLOAD_DATA = "users.download"
	SCOPE_USER_DASHBOARD     = "users.dashboard"

	//CUSTOMER
	SCOPE_CUSTOMER_DOWNLOAD_DATA = "customers.download"

	//PROFILE
	SCOPE_PROFILE_VIEW           = "profile.view"
	SCOPE_PROFILE_UPDATE         = "profile.update"
	SCOPE_PROFILE_UPLOAD_PICTURE = "profile.upload_picture"

	//FCM TOKEN
	SCOPE_FCM_STORE_TOKEN  = "fcm.store_token"
	SCOPE_FCM_REMOVE_TOKEN = "fcm.remove_token"
	SCOPE_FCM_TEST_NOTIFY  = "fcm.test_notify"

	//ROLE
	SCOPE_ROLE_CREATE   = "roles.create"
	SCOPE_ROLE_UPDATE   = "roles.update"
	SCOPE_ROLE_DELETE   = "roles.delete"
	SCOPE_ROLE_VIEW_ALL = "roles.view_all"
	SCOPE_ROLE_VIEW     = "roles.view"

	//PIN
	SCOPE_PIN_CHECK_PUBLIC  = "pin.check_public"
	SCOPE_PIN_CHANGE_PUBLIC = "pin.change_public"
	SCOPE_PIN_FORGOT_PUBLIC = "pin.forgot_public"
	SCOPE_PIN_RESET_PUBLIC  = "pin.reset_public"
	SCOPE_PIN_CHECK_AGEN    = "pin.check_agen"
	SCOPE_PIN_CHANGE_AGEN   = "pin.change_agen"
	SCOPE_PIN_FORGOT_AGEN   = "pin.forgot_agen"
	SCOPE_PIN_RESET_AGEN    = "pin.reset_agen"

	//ADDRESS
	SCOPE_ADDRESS_CREATE   = "address.create"
	SCOPE_ADDRESS_VIEW_ALL = "address.view_all"
	SCOPE_ADDRESS_VIEW     = "address.view"
	SCOPE_ADDRESS_UPDATE   = "address.update"
	SCOPE_ADDRESS_DELETE   = "address.delete"
	//CUSTOMER MANAGEMENT

	SCOPE_CUSTOMER_DELETE = "customer.delete"
)

const (
	// AGEN
	SCOPE_AGEN_ALL            = "agens.*"
	SCOPE_AGEN_REGISTER       = "agen.register"
	SCOPE_AGEN_GET_PROFILE    = "agen.get_profile"
	SCOPE_AGEN_DOWNLOAD_DATA  = "agen.download"
	SCOPE_AGEN_UPDATE_SELLING = "agen.update_selling_transport"
	SCOPE_AGEN_UPDATE_KTP     = "agen.update_pengajuan_ktp"

	// DASHBOARD
	SCOPE_AGEN_VERIFY_BULK  = "dashboard.agen.verify_bulk"
	SCOPE_AGEN_SUSPEND_BULK = "dashboard.agen.suspend_bulk"
	SCOPE_AGEN_VERIFY       = "dashboard.agen.verify"
	SCOPE_AGEN_SUSPEND      = "dashboard.agen.suspend"
	SCOPE_AGEN_UPDATE       = "dashboard.agen.update"
	SCOPE_AGEN_GET          = "dashboard.agen.get"
	SCOPE_AGEN_GET_ALL      = "dashboard.agen.get_all"
)

const (
	// General Scopes
	SCOPE_PRODUCT_ALL                 = "products.*"
	SCOPE_PRODUCT_CREATE              = "products.create"
	SCOPE_PRODUCT_UPDATE              = "products.update"
	SCOPE_PRODUCT_DELETE              = "products.delete"
	SCOPE_PRODUCT_GET                 = "products.get"
	SCOPE_PRODUCT_GET_ALL             = "products.get_all"
	SCOPE_PRODUCT_DOWNLOAD            = "products.download"
	SCOPE_PRODUCT_NEED_REVIEW         = "products.need_review"
	SCOPE_PRODUCT_ACCEPT_REVIEW       = "products.accept_review"
	SCOPE_PRODUCT_DECLINE_REVIEW      = "products.decline_review"
	SCOPE_PRODUCT_ACCEPT_REVIEW_BULK  = "products.accept_review_bulk"
	SCOPE_PRODUCT_DECLINE_REVIEW_BULK = "products.decline_review_bulk"

	// Agen Scopes
	SCOPE_AGEN_PRODUCT_REQUEST     = "agen.products.request"
	SCOPE_AGEN_PRODUCT_REQUEST_GET = "agen.products.request_get"
	SCOPE_AGEN_PRODUCT_REQUEST_V2  = "agen.products.request_v2"

	// Store Product Scopes
	SCOPE_STORE_PRODUCT_ALL    = "store_products.*"
	SCOPE_STORE_PRODUCT_CREATE = "store_products.create"
	SCOPE_STORE_PRODUCT_GET    = "store_products.get"
	SCOPE_STORE_PRODUCT_UPDATE = "store_products.update"
	SCOPE_STORE_PRODUCT_DELETE = "store_products.delete"

	// Store Management Scopes
	SCOPE_STORE_ALL             = "store.*"
	SCOPE_STORE_CREATE          = "store.create"
	SCOPE_STORE_UPDATE          = "store.update"
	SCOPE_STORE_DELETE          = "store.delete"
	SCOPE_STORE_GET             = "store.get"
	SCOPE_STORE_GET_ALL         = "store.get_all"
	SCOPE_STORE_OPERATION_HOURS = "store.operation_hours"
)

const (
	SCOPE_VOUCHER_ALL           = "vouchers.*"
	SCOPE_VOUCHER_CREATE        = "vouchers.create"
	SCOPE_VOUCHER_UPDATE        = "vouchers.update"
	SCOPE_VOUCHER_UPDATE_STATUS = "vouchers.update-status"
	SCOPE_VOUCHER_GET           = "vouchers.get"
	SCOPE_VOUCHER_GET_ALL       = "vouchers.get-all"
	SCOPE_VOUCHER_GET_ACTIVE    = "vouchers.get-active"
	SCOPE_VOUCHER_DELETE        = "vouchers.delete"
	SCOPE_VOUCHER_GET_AVAILABLE = "vouchers.get-available"
	SCOPE_VOUCHER_GET_CLAIMED   = "vouchers.get-claimed"
	SCOPE_VOUCHER_CLAIM         = "vouchers.claim"
)

const (
	// General Orders Scopes
	SCOPE_ORDERS_ALL = "orders.*"

	// Customer Scopes
	SCOPE_CUSTOMER_CARI_AGEN             = "customer.cari-agen"
	SCOPE_CUSTOMER_CHECKOUT              = "customer.checkout"
	SCOPE_CUSTOMER_CHANGE_PAYMENT        = "customer.change-payment"
	SCOPE_CUSTOMER_ORDERS_GET            = "customer.orders.get"
	SCOPE_CUSTOMER_ORDER_GET             = "customer.order.get"
	SCOPE_CUSTOMER_ORDER_CANCEL          = "customer.order.cancel"
	SCOPE_CUSTOMER_ORDER_ACCEPT_DELIVERY = "customer.order.accept-delivery"
	SCOPE_CUSTOMER_ORDER_GIVE_RATING     = "customer.order.give-rating"

	// Agen Scopes
	SCOPE_AGEN_ORDERS_GET            = "agen.orders.get"
	SCOPE_AGEN_ORDER_GET             = "agen.order.get"
	SCOPE_AGEN_ORDER_PRINT_INVOICE   = "agen.order.print-invoice"
	SCOPE_AGEN_ORDER_ACCEPT          = "agen.order.accept"
	SCOPE_AGEN_ORDER_DECLINE         = "agen.order.decline"
	SCOPE_AGEN_ORDER_DELIVER         = "agen.order.deliver"
	SCOPE_AGEN_ORDER_CANCEL          = "agen.order.cancel"
	SCOPE_AGEN_ORDER_FINISH_DELIVERY = "agen.order.finish-delivery"
	SCOPE_AGEN_STORE_RATINGS_GET     = "agen.store-ratings.get"

	// Dashboard Scopes
	SCOPE_DASHBOARD_ORDERS_GET_ALL  = "dashboard.orders.get-all"
	SCOPE_DASHBOARD_ORDERS_DOWNLOAD = "dashboard.orders.download"
	SCOPE_DASHBOARD_ORDER_GET       = "dashboard.order.get"
	SCOPE_DASHBOARD_ORDER_DELETE    = "dashboard.order.delete"
	SCOPE_DASHBOARD_CARTS_GET       = "dashboard.carts.get"

	// Public Scopes
	SCOPE_PUBLIC_ORDER_UPDATE = "public.orders.update"
	SCOPE_PUBLIC_CART_CREATE  = "public.carts.create"
	SCOPE_PUBLIC_CART_GET     = "public.carts.get"
	SCOPE_PUBLIC_CART_UPDATE  = "public.carts.update"

	// General Carts Scopes
	SCOPE_CARTS_ALL = "carts.*"

	// Cart Scopes
	SCOPE_CUSTOMER_CART_GET      = "customer.carts.get"
	SCOPE_CUSTOMER_CART_ADD      = "customer.carts.add"
	SCOPE_CUSTOMER_CART_DECREASE = "customer.carts.decrease"
	SCOPE_CUSTOMER_CART_UPDATE   = "customer.carts.update"
	SCOPE_CUSTOMER_CART_REMOVE   = "customer.carts.remove"
)

const (

	// General Permissions
	SCOPE_ACCOUNT_ALL = "account.*"

	// Admin Balance and Revenue
	SCOPE_DASHBOARD_ADMIN_BALANCE_VIEW   = "dashboard.admin.balance.view"
	SCOPE_DASHBOARD_DAILY_REVENUE_VIEW   = "dashboard.daily.revenue.view"
	SCOPE_DASHBOARD_ADMIN_CREATE_REVENUE = "dashboard.admin.create.revenue"

	// Agen and Customer Finance
	SCOPE_AGEN_BALANCE_VIEW     = "agen.balance.view"
	SCOPE_AGEN_FINANCE_VIEW     = "agen.finance.view"
	SCOPE_CUSTOMER_BALANCE_VIEW = "customer.balance.view"

	// Mutation History and Injection
	SCOPE_HISTORY_MUTATION_VIEW    = "history.mutation.view"
	SCOPE_HISTORY_INJECT_USER_VIEW = "history.inject.user.view"
	SCOPE_HISTORY_INJECT_VIEW      = "history.inject.view"
	SCOPE_HISTORY_INJECT_MANAGE    = "history.inject.manage"

	// Disbursement
	SCOPE_DISBURSEMENT_VIEW         = "disbursement.view"
	SCOPE_DISBURSEMENT_DOWNLOAD     = "disbursement.download"
	SCOPE_DISBURSEMENT_HISTORY_VIEW = "disbursement.history.view"
	SCOPE_DISBURSEMENT_REQUEST      = "disbursement.request"
	SCOPE_DISBURSEMENT_APPROVE      = "disbursement.approve"
	SCOPE_DISBURSEMENT_REJECT       = "disbursement.reject"
	SCOPE_DISBURSEMENT_BULK_APPROVE = "disbursement.bulk.approve"
	SCOPE_DISBURSEMENT_BULK_REJECT  = "disbursement.bulk.reject"

	// Payments
	SCOPE_PAYMENT_DOKU_CALLBACK    = "payment.doku.callback"
	SCOPE_PAYMENT_DOKU_STATUS_VIEW = "payment.doku.status.view"

	// Bank Account
	SCOPE_BANK_ACCOUNT_VIEW   = "bank.account.view"
	SCOPE_BANK_ACCOUNT_MANAGE = "bank.account.manage"

	// Dashboard History Mutations
	SCOPE_DASHBOARD_DISBURSEMENT_COUNT_VIEW   = "dashboard.disbursement.count.view"
	SCOPE_DASHBOARD_HISTORY_MUTATION_CUSTOMER = "dashboard.history.mutation.customer"
	SCOPE_DASHBOARD_HISTORY_MUTATION_AGEN     = "dashboard.history.mutation.agen"
)

const (

	// STOCK
	SCOPE_STOCK_ALL    = "stocks.*"
	SCOPE_STOCK_CREATE = "stocks.create"
	SCOPE_STOCK_LIST   = "stocks.list"
	SCOPE_STOCK_DETAIL = "stocks.detail"
	SCOPE_STOCK_UPDATE = "stocks.update"
	SCOPE_STOCK_DELETE = "stocks.delete"

	// PAYMENT
	SCOPE_PAYMENT_ALL    = "payments.*"
	SCOPE_PAYMENT_CREATE = "payments.create"
	SCOPE_PAYMENT_LIST   = "payments.list"
	SCOPE_PAYMENT_DETAIL = "payments.detail"
	SCOPE_PAYMENT_UPDATE = "payments.update"
	SCOPE_PAYMENT_DELETE = "payments.delete"
)
