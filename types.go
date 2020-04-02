package paypal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	// APIBaseSandBox points to the sandbox (for testing) version of the API
	APIBaseSandBox = "https://api.sandbox.paypal.com"

	// APIBaseLive points to the live version of the API
	APIBaseLive = "https://api.paypal.com"

	// RequestNewTokenBeforeExpiresIn is used by SendWithAuth and try to get new Token when it's about to expire
	RequestNewTokenBeforeExpiresIn = time.Duration(60) * time.Second
)

// Possible values for `no_shipping` in InputFields
//
// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
const (
	NoShippingDisplay      uint = 0
	NoShippingHide         uint = 1
	NoShippingBuyerAccount uint = 2
)

// Possible values for `address_override` in InputFields
//
// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
const (
	AddrOverrideFromFile uint = 0
	AddrOverrideFromCall uint = 1
)

// Possible values for `landing_page_type` in FlowConfig
//
// https://developer.paypal.com/docs/api/payment-experience/#definition-flow_config
const (
	LandingPageTypeBilling string = "Billing"
	LandingPageTypeLogin   string = "Login"
)

// Possible value for `allowed_payment_method` in PaymentOptions
//
// https://developer.paypal.com/docs/api/payments/#definition-payment_options
const (
	AllowedPaymentUnrestricted         string = "UNRESTRICTED"
	AllowedPaymentInstantFundingSource string = "INSTANT_FUNDING_SOURCE"
	AllowedPaymentImmediatePay         string = "IMMEDIATE_PAY"
)

// Possible value for `intent` in CreateOrder
//
// https://developer.paypal.com/docs/api/orders/v2/#orders_create
const (
	OrderIntentCapture   string = "CAPTURE"
	OrderIntentAuthorize string = "AUTHORIZE"
)

// Possible values for `category` in Item
//
// https://developer.paypal.com/docs/api/orders/v2/#definition-item
const (
	ItemCategoryDigitalGood  string = "DIGITAL_GOODS"
	ItemCategoryPhysicalGood string = "PHYSICAL_GOODS"
)

// Possible values for `shipping_preference` in ApplicationContext
//
// https://developer.paypal.com/docs/api/orders/v2/#definition-application_context
const (
	ShippingPreferenceGetFromFile        string = "GET_FROM_FILE"
	ShippingPreferenceNoShipping         string = "NO_SHIPPING"
	ShippingPreferenceSetProvidedAddress string = "SET_PROVIDED_ADDRESS"
)

const (
	EventPaymentCaptureCompleted       string = "PAYMENT.CAPTURE.COMPLETED"
	EventPaymentCaptureDenied          string = "PAYMENT.CAPTURE.DENIED"
	EventPaymentCaptureRefunded        string = "PAYMENT.CAPTURE.REFUNDED"
	EventMerchantOnboardingCompleted   string = "MERCHANT.ONBOARDING.COMPLETED"
	EventMerchantPartnerConsentRevoked string = "MERCHANT.PARTNER-CONSENT.REVOKED"
)

const (
	OperationAPIIntegration   string = "API_INTEGRATION"
	ProductExpressCheckout    string = "EXPRESS_CHECKOUT"
	IntegrationMethodPayPal   string = "PAYPAL"
	IntegrationTypeThirdParty string = "THIRD_PARTY"
	ConsentShareData          string = "SHARE_DATA_CONSENT"
)

const (
	FeaturePayment               string = "PAYMENT"
	FeatureRefund                string = "REFUND"
	FeatureFuturePayment         string = "FUTURE_PAYMENT"
	FeatureDirectPayment         string = "DIRECT_PAYMENT"
	FeaturePartnerFee            string = "PARTNER_FEE"
	FeatureDelayFunds            string = "DELAY_FUNDS_DISBURSEMENT"
	FeatureReadSellerDispute     string = "READ_SELLER_DISPUTE"
	FeatureUpdateSellerDispute   string = "UPDATE_SELLER_DISPUTE"
	FeatureDisputeReadBuyer      string = "DISPUTE_READ_BUYER"
	FeatureUpdateCustomerDispute string = "UPDATE_CUSTOMER_DISPUTES"
)

const (
	LinkRelSelf      string = "self"
	LinkRelActionURL string = "action_url"
)

// Possible values for PatchObject.Operation
const (
	Add     string = "add"
	Replace string = "replace"
	Remove  string = "remove"
)

// Possible values for PatchObject.Path when you update product
const (
	Description string = "/description"
	Category    string = "/category"
	ImageUrl    string = "/image_url"
	HomeUrl     string = "/home_url"
)

// Possible values for CreateProduct.Type
const (
	PHYSICAL string = "PHYSICAL"
	DIGITAL  string = "DIGITAL"
	SERVICE  string = "SERVICE"
)

// Possible values for CreateProduct.Category
const (
	AC_REFRIGERATION_REPAIR                                       string = "AC_REFRIGERATION_REPAIR"
	ACADEMIC_SOFTWARE                                             string = "ACADEMIC_SOFTWARE"
	ACCESSORIES                                                   string = "ACCESSORIES"
	ACCOUNTING                                                    string = "ACCOUNTING"
	ADULT                                                         string = "ADULT"
	ADVERTISING                                                   string = "ADVERTISING"
	AFFILIATED_AUTO_RENTAL                                        string = "AFFILIATED_AUTO_RENTAL"
	AGENCIES                                                      string = "AGENCIES"
	AGGREGATORS                                                   string = "AGGREGATORS"
	AGRICULTURAL_COOPERATIVE_FOR_MAIL_ORDER                       string = "AGRICULTURAL_COOPERATIVE_FOR_MAIL_ORDER"
	AIR_CARRIERS_AIRLINES                                         string = "AIR_CARRIERS_AIRLINES"
	AIRLINES                                                      string = "AIRLINES"
	AIRPORTS_FLYING_FIELDS                                        string = "AIRPORTS_FLYING_FIELDS"
	ALCOHOLIC_BEVERAGES                                           string = "ALCOHOLIC_BEVERAGES"
	AMUSEMENT_PARKS_CARNIVALS                                     string = "AMUSEMENT_PARKS_CARNIVALS"
	ANIMATION                                                     string = "ANIMATION"
	ANTIQUES                                                      string = "ANTIQUES"
	APPLIANCES                                                    string = "APPLIANCES"
	AQUARIAMS_SEAQUARIUMS_DOLPHINARIUMS                           string = "AQUARIAMS_SEAQUARIUMS_DOLPHINARIUMS"
	ARCHITECTURAL_ENGINEERING_AND_SURVEYING_SERVICES              string = "ARCHITECTURAL_ENGINEERING_AND_SURVEYING_SERVICES"
	ART_AND_CRAFT_SUPPLIES                                        string = "ART_AND_CRAFT_SUPPLIES"
	ART_DEALERS_AND_GALLERIES                                     string = "ART_DEALERS_AND_GALLERIES"
	ARTIFACTS_GRAVE_RELATED_AND_NATIVE_AMERICAN_CRAFTS            string = "ARTIFACTS_GRAVE_RELATED_AND_NATIVE_AMERICAN_CRAFTS"
	ARTS_AND_CRAFTS                                               string = "ARTS_AND_CRAFTS"
	ARTS_CRAFTS_AND_COLLECTIBLES                                  string = "ARTS_CRAFTS_AND_COLLECTIBLES"
	AUDIO_BOOKS                                                   string = "AUDIO_BOOKS"
	AUTO_ASSOCIATIONS_CLUBS                                       string = "AUTO_ASSOCIATIONS_CLUBS"
	AUTO_DEALER_USED_ONLY                                         string = "AUTO_DEALER_USED_ONLY"
	AUTO_RENTALS                                                  string = "AUTO_RENTALS"
	AUTO_SERVICE                                                  string = "AUTO_SERVICE"
	AUTOMATED_FUEL_DISPENSERS                                     string = "AUTOMATED_FUEL_DISPENSERS"
	AUTOMOBILE_ASSOCIATIONS                                       string = "AUTOMOBILE_ASSOCIATIONS"
	AUTOMOTIVE                                                    string = "AUTOMOTIVE"
	AUTOMOTIVE_REPAIR_SHOPS_NON_DEALER                            string = "AUTOMOTIVE_REPAIR_SHOPS_NON_DEALER"
	AUTOMOTIVE_TOP_AND_BODY_SHOPS                                 string = "AUTOMOTIVE_TOP_AND_BODY_SHOPS"
	AVIATION                                                      string = "AVIATION"
	BABIES_CLOTHING_AND_SUPPLIES                                  string = "BABIES_CLOTHING_AND_SUPPLIES"
	BABY                                                          string = "BABY"
	BANDS_ORCHESTRAS_ENTERTAINERS                                 string = "BANDS_ORCHESTRAS_ENTERTAINERS"
	BARBIES                                                       string = "BARBIES"
	BATH_AND_BODY                                                 string = "BATH_AND_BODY"
	BATTERIES                                                     string = "BATTERIES"
	BEAN_BABIES                                                   string = "BEAN_BABIES"
	BEAUTY                                                        string = "BEAUTY"
	BEAUTY_AND_FRAGRANCES                                         string = "BEAUTY_AND_FRAGRANCES"
	BED_AND_BATH                                                  string = "BED_AND_BATH"
	BICYCLE_SHOPS_SALES_AND_SERVICE                               string = "BICYCLE_SHOPS_SALES_AND_SERVICE"
	BICYCLES_AND_ACCESSORIES                                      string = "BICYCLES_AND_ACCESSORIES"
	BILLIARD_POOL_ESTABLISHMENTS                                  string = "BILLIARD_POOL_ESTABLISHMENTS"
	BOAT_DEALERS                                                  string = "BOAT_DEALERS"
	BOAT_RENTALS_AND_LEASING                                      string = "BOAT_RENTALS_AND_LEASING"
	BOATING_SAILING_AND_ACCESSORIES                               string = "BOATING_SAILING_AND_ACCESSORIES"
	BOOKS                                                         string = "BOOKS"
	BOOKS_AND_MAGAZINES                                           string = "BOOKS_AND_MAGAZINES"
	BOOKS_MANUSCRIPTS                                             string = "BOOKS_MANUSCRIPTS"
	BOOKS_PERIODICALS_AND_NEWSPAPERS                              string = "BOOKS_PERIODICALS_AND_NEWSPAPERS"
	BOWLING_ALLEYS                                                string = "BOWLING_ALLEYS"
	BULLETIN_BOARD                                                string = "BULLETIN_BOARD"
	BUS_LINE                                                      string = "BUS_LINE"
	BUS_LINES_CHARTERS_TOUR_BUSES                                 string = "BUS_LINES_CHARTERS_TOUR_BUSES"
	BUSINESS                                                      string = "BUSINESS"
	BUSINESS_AND_SECRETARIAL_SCHOOLS                              string = "BUSINESS_AND_SECRETARIAL_SCHOOLS"
	BUYING_AND_SHOPPING_SERVICES_AND_CLUBS                        string = "BUYING_AND_SHOPPING_SERVICES_AND_CLUBS"
	CABLE_SATELLITE_AND_OTHER_PAY_TELEVISION_AND_RADIO_SERVICES   string = "CABLE_SATELLITE_AND_OTHER_PAY_TELEVISION_AND_RADIO_SERVICES"
	CABLE_SATELLITE_AND_OTHER_PAY_TV_AND_RADIO                    string = "CABLE_SATELLITE_AND_OTHER_PAY_TV_AND_RADIO"
	CAMERA_AND_PHOTOGRAPHIC_SUPPLIES                              string = "CAMERA_AND_PHOTOGRAPHIC_SUPPLIES"
	CAMERAS                                                       string = "CAMERAS"
	CAMERAS_AND_PHOTOGRAPHY                                       string = "CAMERAS_AND_PHOTOGRAPHY"
	CAMPER_RECREATIONAL_AND_UTILITY_TRAILER_DEALERS               string = "CAMPER_RECREATIONAL_AND_UTILITY_TRAILER_DEALERS"
	CAMPING_AND_OUTDOORS                                          string = "CAMPING_AND_OUTDOORS"
	CAMPING_AND_SURVIVAL                                          string = "CAMPING_AND_SURVIVAL"
	CAR_AND_TRUCK_DEALERS                                         string = "CAR_AND_TRUCK_DEALERS"
	CAR_AND_TRUCK_DEALERS_USED_ONLY                               string = "CAR_AND_TRUCK_DEALERS_USED_ONLY"
	CAR_AUDIO_AND_ELECTRONICS                                     string = "CAR_AUDIO_AND_ELECTRONICS"
	CAR_RENTAL_AGENCY                                             string = "CAR_RENTAL_AGENCY"
	CATALOG_MERCHANT                                              string = "CATALOG_MERCHANT"
	CATALOG_RETAIL_MERCHANT                                       string = "CATALOG_RETAIL_MERCHANT"
	CATERING_SERVICES                                             string = "CATERING_SERVICES"
	CHARITY                                                       string = "CHARITY"
	CHECK_CASHIER                                                 string = "CHECK_CASHIER"
	CHILD_CARE_SERVICES                                           string = "CHILD_CARE_SERVICES"
	CHILDREN_BOOKS                                                string = "CHILDREN_BOOKS"
	CHIROPODISTS_PODIATRISTS                                      string = "CHIROPODISTS_PODIATRISTS"
	CHIROPRACTORS                                                 string = "CHIROPRACTORS"
	CIGAR_STORES_AND_STANDS                                       string = "CIGAR_STORES_AND_STANDS"
	CIVIC_SOCIAL_FRATERNAL_ASSOCIATIONS                           string = "CIVIC_SOCIAL_FRATERNAL_ASSOCIATIONS"
	CIVIL_SOCIAL_FRAT_ASSOCIATIONS                                string = "CIVIL_SOCIAL_FRAT_ASSOCIATIONS"
	CLOTHING                                                      string = "CLOTHING"
	CLOTHING_ACCESSORIES_AND_SHOES                                string = "CLOTHING_ACCESSORIES_AND_SHOES"
	CLOTHING_RENTAL                                               string = "CLOTHING_RENTAL"
	COFFEE_AND_TEA                                                string = "COFFEE_AND_TEA"
	COIN_OPERATED_BANKS_AND_CASINOS                               string = "COIN_OPERATED_BANKS_AND_CASINOS"
	COLLECTIBLES                                                  string = "COLLECTIBLES"
	COLLECTION_AGENCY                                             string = "COLLECTION_AGENCY"
	COLLEGES_AND_UNIVERSITIES                                     string = "COLLEGES_AND_UNIVERSITIES"
	COMMERCIAL_EQUIPMENT                                          string = "COMMERCIAL_EQUIPMENT"
	COMMERCIAL_FOOTWEAR                                           string = "COMMERCIAL_FOOTWEAR"
	COMMERCIAL_PHOTOGRAPHY                                        string = "COMMERCIAL_PHOTOGRAPHY"
	COMMERCIAL_PHOTOGRAPHY_ART_AND_GRAPHICS                       string = "COMMERCIAL_PHOTOGRAPHY_ART_AND_GRAPHICS"
	COMMERCIAL_SPORTS_PROFESSIONA                                 string = "COMMERCIAL_SPORTS_PROFESSIONA"
	COMMODITIES_AND_FUTURES_EXCHANGE                              string = "COMMODITIES_AND_FUTURES_EXCHANGE"
	COMPUTER_AND_DATA_PROCESSING_SERVICES                         string = "COMPUTER_AND_DATA_PROCESSING_SERVICES"
	COMPUTER_HARDWARE_AND_SOFTWARE                                string = "COMPUTER_HARDWARE_AND_SOFTWARE"
	COMPUTER_MAINTENANCE_REPAIR_AND_SERVICES_NOT_ELSEWHERE_CLAS   string = "COMPUTER_MAINTENANCE_REPAIR_AND_SERVICES_NOT_ELSEWHERE_CLAS"
	CONSTRUCTION                                                  string = "CONSTRUCTION"
	CONSTRUCTION_MATERIALS_NOT_ELSEWHERE_CLASSIFIED               string = "CONSTRUCTION_MATERIALS_NOT_ELSEWHERE_CLASSIFIED"
	CONSULTING_SERVICES                                           string = "CONSULTING_SERVICES"
	CONSUMER_CREDIT_REPORTING_AGENCIES                            string = "CONSUMER_CREDIT_REPORTING_AGENCIES"
	CONVALESCENT_HOMES                                            string = "CONVALESCENT_HOMES"
	COSMETIC_STORES                                               string = "COSMETIC_STORES"
	COUNSELING_SERVICES_DEBT_MARRIAGE_PERSONAL                    string = "COUNSELING_SERVICES_DEBT_MARRIAGE_PERSONAL"
	COUNTERFEIT_CURRENCY_AND_STAMPS                               string = "COUNTERFEIT_CURRENCY_AND_STAMPS"
	COUNTERFEIT_ITEMS                                             string = "COUNTERFEIT_ITEMS"
	COUNTRY_CLUBS                                                 string = "COUNTRY_CLUBS"
	COURIER_SERVICES                                              string = "COURIER_SERVICES"
	COURIER_SERVICES_AIR_AND_GROUND_AND_FREIGHT_FORWARDERS        string = "COURIER_SERVICES_AIR_AND_GROUND_AND_FREIGHT_FORWARDERS"
	COURT_COSTS_ALIMNY_CHILD_SUPT                                 string = "COURT_COSTS_ALIMNY_CHILD_SUPT"
	COURT_COSTS_INCLUDING_ALIMONY_AND_CHILD_SUPPORT_COURTS_OF_LAW string = "COURT_COSTS_INCLUDING_ALIMONY_AND_CHILD_SUPPORT_COURTS_OF_LAW"
	CREDIT_CARDS                                                  string = "CREDIT_CARDS"
	CREDIT_UNION                                                  string = "CREDIT_UNION"
	CULTURE_AND_RELIGION                                          string = "CULTURE_AND_RELIGION"
	DAIRY_PRODUCTS_STORES                                         string = "DAIRY_PRODUCTS_STORES"
	DANCE_HALLS_STUDIOS_AND_SCHOOLS                               string = "DANCE_HALLS_STUDIOS_AND_SCHOOLS"
	DECORATIVE                                                    string = "DECORATIVE"
	DENTAL                                                        string = "DENTAL"
	DENTISTS_AND_ORTHODONTISTS                                    string = "DENTISTS_AND_ORTHODONTISTS"
	DEPARTMENT_STORES                                             string = "DEPARTMENT_STORES"
	DESKTOP_PCS                                                   string = "DESKTOP_PCS"
	DEVICES                                                       string = "DEVICES"
	DIECAST_TOYS_VEHICLES                                         string = "DIECAST_TOYS_VEHICLES"
	DIGITAL_GAMES                                                 string = "DIGITAL_GAMES"
	DIGITAL_MEDIA_BOOKS_MOVIES_MUSIC                              string = "DIGITAL_MEDIA_BOOKS_MOVIES_MUSIC"
	DIRECT_MARKETING                                              string = "DIRECT_MARKETING"
	DIRECT_MARKETING_CATALOG_MERCHANT                             string = "DIRECT_MARKETING_CATALOG_MERCHANT"
	DIRECT_MARKETING_INBOUND_TELE                                 string = "DIRECT_MARKETING_INBOUND_TELE"
	DIRECT_MARKETING_OUTBOUND_TELE                                string = "DIRECT_MARKETING_OUTBOUND_TELE"
	DIRECT_MARKETING_SUBSCRIPTION                                 string = "DIRECT_MARKETING_SUBSCRIPTION"
	DISCOUNT_STORES                                               string = "DISCOUNT_STORES"
	DOOR_TO_DOOR_SALES                                            string = "DOOR_TO_DOOR_SALES"
	DRAPERY_WINDOW_COVERING_AND_UPHOLSTERY                        string = "DRAPERY_WINDOW_COVERING_AND_UPHOLSTERY"
	DRINKING_PLACES                                               string = "DRINKING_PLACES"
	DRUGSTORE                                                     string = "DRUGSTORE"
	DURABLE_GOODS                                                 string = "DURABLE_GOODS"
	ECOMMERCE_DEVELOPMENT                                         string = "ECOMMERCE_DEVELOPMENT"
	ECOMMERCE_SERVICES                                            string = "ECOMMERCE_SERVICES"
	EDUCATIONAL_AND_TEXTBOOKS                                     string = "EDUCATIONAL_AND_TEXTBOOKS"
	ELECTRIC_RAZOR_STORES                                         string = "ELECTRIC_RAZOR_STORES"
	ELECTRICAL_AND_SMALL_APPLIANCE_REPAIR                         string = "ELECTRICAL_AND_SMALL_APPLIANCE_REPAIR"
	ELECTRICAL_CONTRACTORS                                        string = "ELECTRICAL_CONTRACTORS"
	ELECTRICAL_PARTS_AND_EQUIPMENT                                string = "ELECTRICAL_PARTS_AND_EQUIPMENT"
	ELECTRONIC_CASH                                               string = "ELECTRONIC_CASH"
	ELEMENTARY_AND_SECONDARY_SCHOOLS                              string = "ELEMENTARY_AND_SECONDARY_SCHOOLS"
	EMPLOYMENT                                                    string = "EMPLOYMENT"
	ENTERTAINERS                                                  string = "ENTERTAINERS"
	ENTERTAINMENT_AND_MEDIA                                       string = "ENTERTAINMENT_AND_MEDIA"
	EQUIP_TOOL_FURNITURE_AND_APPLIANCE_RENTAL_AND_LEASING         string = "EQUIP_TOOL_FURNITURE_AND_APPLIANCE_RENTAL_AND_LEASING"
	ESCROW                                                        string = "ESCROW"
	EVENT_AND_WEDDING_PLANNING                                    string = "EVENT_AND_WEDDING_PLANNING"
	EXERCISE_AND_FITNESS                                          string = "EXERCISE_AND_FITNESS"
	EXERCISE_EQUIPMENT                                            string = "EXERCISE_EQUIPMENT"
	EXTERMINATING_AND_DISINFECTING_SERVICES                       string = "EXTERMINATING_AND_DISINFECTING_SERVICES"
	FABRICS_AND_SEWING                                            string = "FABRICS_AND_SEWING"
	FAMILY_CLOTHING_STORES                                        string = "FAMILY_CLOTHING_STORES"
	FASHION_JEWELRY                                               string = "FASHION_JEWELRY"
	FAST_FOOD_RESTAURANTS                                         string = "FAST_FOOD_RESTAURANTS"
	FICTION_AND_NONFICTION                                        string = "FICTION_AND_NONFICTION"
	FINANCE_COMPANY                                               string = "FINANCE_COMPANY"
	FINANCIAL_AND_INVESTMENT_ADVICE                               string = "FINANCIAL_AND_INVESTMENT_ADVICE"
	FINANCIAL_INSTITUTIONS_MERCHANDISE_AND_SERVICES               string = "FINANCIAL_INSTITUTIONS_MERCHANDISE_AND_SERVICES"
	FIREARM_ACCESSORIES                                           string = "FIREARM_ACCESSORIES"
	FIREARMS_WEAPONS_AND_KNIVES                                   string = "FIREARMS_WEAPONS_AND_KNIVES"
	FIREPLACE_AND_FIREPLACE_SCREENS                               string = "FIREPLACE_AND_FIREPLACE_SCREENS"
	FIREWORKS                                                     string = "FIREWORKS"
	FISHING                                                       string = "FISHING"
	FLORISTS                                                      string = "FLORISTS"
	FLOWERS                                                       string = "FLOWERS"
	FOOD_DRINK_AND_NUTRITION                                      string = "FOOD_DRINK_AND_NUTRITION"
	FOOD_PRODUCTS                                                 string = "FOOD_PRODUCTS"
	FOOD_RETAIL_AND_SERVICE                                       string = "FOOD_RETAIL_AND_SERVICE"
	FRAGRANCES_AND_PERFUMES                                       string = "FRAGRANCES_AND_PERFUMES"
	FREEZER_AND_LOCKER_MEAT_PROVISIONERS                          string = "FREEZER_AND_LOCKER_MEAT_PROVISIONERS"
	FUEL_DEALERS_FUEL_OIL_WOOD_AND_COAL                           string = "FUEL_DEALERS_FUEL_OIL_WOOD_AND_COAL"
	FUEL_DEALERS_NON_AUTOMOTIVE                                   string = "FUEL_DEALERS_NON_AUTOMOTIVE"
	FUNERAL_SERVICES_AND_CREMATORIES                              string = "FUNERAL_SERVICES_AND_CREMATORIES"
	FURNISHING_AND_DECORATING                                     string = "FURNISHING_AND_DECORATING"
	FURNITURE                                                     string = "FURNITURE"
	FURRIERS_AND_FUR_SHOPS                                        string = "FURRIERS_AND_FUR_SHOPS"
	GADGETS_AND_OTHER_ELECTRONICS                                 string = "GADGETS_AND_OTHER_ELECTRONICS"
	GAMBLING                                                      string = "GAMBLING"
	GAME_SOFTWARE                                                 string = "GAME_SOFTWARE"
	GAMES                                                         string = "GAMES"
	GARDEN_SUPPLIES                                               string = "GARDEN_SUPPLIES"
	GENERAL                                                       string = "GENERAL"
	GENERAL_CONTRACTORS                                           string = "GENERAL_CONTRACTORS"
	GENERAL_GOVERNMENT                                            string = "GENERAL_GOVERNMENT"
	GENERAL_SOFTWARE                                              string = "GENERAL_SOFTWARE"
	GENERAL_TELECOM                                               string = "GENERAL_TELECOM"
	GIFTS_AND_FLOWERS                                             string = "GIFTS_AND_FLOWERS"
	GLASS_PAINT_AND_WALLPAPER_STORES                              string = "GLASS_PAINT_AND_WALLPAPER_STORES"
	GLASSWARE_CRYSTAL_STORES                                      string = "GLASSWARE_CRYSTAL_STORES"
	GOVERNMENT                                                    string = "GOVERNMENT"
	GOVERNMENT_IDS_AND_LICENSES                                   string = "GOVERNMENT_IDS_AND_LICENSES"
	GOVERNMENT_LICENSED_ON_LINE_CASINOS_ON_LINE_GAMBLING          string = "GOVERNMENT_LICENSED_ON_LINE_CASINOS_ON_LINE_GAMBLING"
	GOVERNMENT_OWNED_LOTTERIES                                    string = "GOVERNMENT_OWNED_LOTTERIES"
	GOVERNMENT_SERVICES                                           string = "GOVERNMENT_SERVICES"
	GRAPHIC_AND_COMMERCIAL_DESIGN                                 string = "GRAPHIC_AND_COMMERCIAL_DESIGN"
	GREETING_CARDS                                                string = "GREETING_CARDS"
	GROCERY_STORES_AND_SUPERMARKETS                               string = "GROCERY_STORES_AND_SUPERMARKETS"
	HARDWARE_AND_TOOLS                                            string = "HARDWARE_AND_TOOLS"
	HARDWARE_EQUIPMENT_AND_SUPPLIES                               string = "HARDWARE_EQUIPMENT_AND_SUPPLIES"
	HAZARDOUS_RESTRICTED_AND_PERISHABLE_ITEMS                     string = "HAZARDOUS_RESTRICTED_AND_PERISHABLE_ITEMS"
	HEALTH_AND_BEAUTY_SPAS                                        string = "HEALTH_AND_BEAUTY_SPAS"
	HEALTH_AND_NUTRITION                                          string = "HEALTH_AND_NUTRITION"
	HEALTH_AND_PERSONAL_CARE                                      string = "HEALTH_AND_PERSONAL_CARE"
	HEARING_AIDS_SALES_AND_SUPPLIES                               string = "HEARING_AIDS_SALES_AND_SUPPLIES"
	HEATING_PLUMBING_AC                                           string = "HEATING_PLUMBING_AC"
	HIGH_RISK_MERCHANT                                            string = "HIGH_RISK_MERCHANT"
	HIRING_SERVICES                                               string = "HIRING_SERVICES"
	HOBBIES_TOYS_AND_GAMES                                        string = "HOBBIES_TOYS_AND_GAMES"
	HOME_AND_GARDEN                                               string = "HOME_AND_GARDEN"
	HOME_AUDIO                                                    string = "HOME_AUDIO"
	HOME_DECOR                                                    string = "HOME_DECOR"
	HOME_ELECTRONICS                                              string = "HOME_ELECTRONICS"
	HOSPITALS                                                     string = "HOSPITALS"
	HOTELS_MOTELS_INNS_RESORTS                                    string = "HOTELS_MOTELS_INNS_RESORTS"
	HOUSEWARES                                                    string = "HOUSEWARES"
	HUMAN_PARTS_AND_REMAINS                                       string = "HUMAN_PARTS_AND_REMAINS"
	HUMOROUS_GIFTS_AND_NOVELTIES                                  string = "HUMOROUS_GIFTS_AND_NOVELTIES"
	HUNTING                                                       string = "HUNTING"
	IDS_LICENSES_AND_PASSPORTS                                    string = "IDS_LICENSES_AND_PASSPORTS"
	ILLEGAL_DRUGS_AND_PARAPHERNALIA                               string = "ILLEGAL_DRUGS_AND_PARAPHERNALIA"
	INDUSTRIAL                                                    string = "INDUSTRIAL"
	INDUSTRIAL_AND_MANUFACTURING_SUPPLIES                         string = "INDUSTRIAL_AND_MANUFACTURING_SUPPLIES"
	INSURANCE_AUTO_AND_HOME                                       string = "INSURANCE_AUTO_AND_HOME"
	INSURANCE_DIRECT                                              string = "INSURANCE_DIRECT"
	INSURANCE_LIFE_AND_ANNUITY                                    string = "INSURANCE_LIFE_AND_ANNUITY"
	INSURANCE_SALES_UNDERWRITING                                  string = "INSURANCE_SALES_UNDERWRITING"
	INSURANCE_UNDERWRITING_PREMIUMS                               string = "INSURANCE_UNDERWRITING_PREMIUMS"
	INTERNET_AND_NETWORK_SERVICES                                 string = "INTERNET_AND_NETWORK_SERVICES"
	INTRA_COMPANY_PURCHASES                                       string = "INTRA_COMPANY_PURCHASES"
	LABORATORIES_DENTAL_MEDICAL                                   string = "LABORATORIES_DENTAL_MEDICAL"
	LANDSCAPING                                                   string = "LANDSCAPING"
	LANDSCAPING_AND_HORTICULTURAL_SERVICES                        string = "LANDSCAPING_AND_HORTICULTURAL_SERVICES"
	LAUNDRY_CLEANING_SERVICES                                     string = "LAUNDRY_CLEANING_SERVICES"
	LEGAL                                                         string = "LEGAL"
	LEGAL_SERVICES_AND_ATTORNEYS                                  string = "LEGAL_SERVICES_AND_ATTORNEYS"
	LOCAL_DELIVERY_SERVICE                                        string = "LOCAL_DELIVERY_SERVICE"
	LOCKSMITH                                                     string = "LOCKSMITH"
	LODGING_AND_ACCOMMODATIONS                                    string = "LODGING_AND_ACCOMMODATIONS"
	LOTTERY_AND_CONTESTS                                          string = "LOTTERY_AND_CONTESTS"
	LUGGAGE_AND_LEATHER_GOODS                                     string = "LUGGAGE_AND_LEATHER_GOODS"
	LUMBER_AND_BUILDING_MATERIALS                                 string = "LUMBER_AND_BUILDING_MATERIALS"
	MAGAZINES                                                     string = "MAGAZINES"
	MAINTENANCE_AND_REPAIR_SERVICES                               string = "MAINTENANCE_AND_REPAIR_SERVICES"
	MAKEUP_AND_COSMETICS                                          string = "MAKEUP_AND_COSMETICS"
	MANUAL_CASH_DISBURSEMENTS                                     string = "MANUAL_CASH_DISBURSEMENTS"
	MASSAGE_PARLORS                                               string = "MASSAGE_PARLORS"
	MEDICAL                                                       string = "MEDICAL"
	MEDICAL_AND_PHARMACEUTICAL                                    string = "MEDICAL_AND_PHARMACEUTICAL"
	MEDICAL_CARE                                                  string = "MEDICAL_CARE"
	MEDICAL_EQUIPMENT_AND_SUPPLIES                                string = "MEDICAL_EQUIPMENT_AND_SUPPLIES"
	MEDICAL_SERVICES                                              string = "MEDICAL_SERVICES"
	MEETING_PLANNERS                                              string = "MEETING_PLANNERS"
	MEMBERSHIP_CLUBS_AND_ORGANIZATIONS                            string = "MEMBERSHIP_CLUBS_AND_ORGANIZATIONS"
	MEMBERSHIP_COUNTRY_CLUBS_GOLF                                 string = "MEMBERSHIP_COUNTRY_CLUBS_GOLF"
	MEMORABILIA                                                   string = "MEMORABILIA"
	MEN_AND_BOY_CLOTHING_AND_ACCESSORY_STORES                     string = "MEN_AND_BOY_CLOTHING_AND_ACCESSORY_STORES"
	MEN_CLOTHING                                                  string = "MEN_CLOTHING"
	MERCHANDISE                                                   string = "MERCHANDISE"
	METAPHYSICAL                                                  string = "METAPHYSICAL"
	MILITARIA                                                     string = "MILITARIA"
	MILITARY_AND_CIVIL_SERVICE_UNIFORMS                           string = "MILITARY_AND_CIVIL_SERVICE_UNIFORMS"
	MISC                                                          string = "MISC"
	MISCELLANEOUS_GENERAL_SERVICES                                string = "MISCELLANEOUS_GENERAL_SERVICES"
	MISCELLANEOUS_REPAIR_SHOPS_AND_RELATED_SERVICES               string = "MISCELLANEOUS_REPAIR_SHOPS_AND_RELATED_SERVICES"
	MODEL_KITS                                                    string = "MODEL_KITS"
	MONEY_TRANSFER_MEMBER_FINANCIAL_INSTITUTION                   string = "MONEY_TRANSFER_MEMBER_FINANCIAL_INSTITUTION"
	MONEY_TRANSFER_MERCHANT                                       string = "MONEY_TRANSFER_MERCHANT"
	MOTION_PICTURE_THEATERS                                       string = "MOTION_PICTURE_THEATERS"
	MOTOR_FREIGHT_CARRIERS_AND_TRUCKING                           string = "MOTOR_FREIGHT_CARRIERS_AND_TRUCKING"
	MOTOR_HOME_AND_RECREATIONAL_VEHICLE_RENTAL                    string = "MOTOR_HOME_AND_RECREATIONAL_VEHICLE_RENTAL"
	MOTOR_HOMES_DEALERS                                           string = "MOTOR_HOMES_DEALERS"
	MOTOR_VEHICLE_SUPPLIES_AND_NEW_PARTS                          string = "MOTOR_VEHICLE_SUPPLIES_AND_NEW_PARTS"
	MOTORCYCLE_DEALERS                                            string = "MOTORCYCLE_DEALERS"
	MOTORCYCLES                                                   string = "MOTORCYCLES"
	MOVIE                                                         string = "MOVIE"
	MOVIE_TICKETS                                                 string = "MOVIE_TICKETS"
	MOVING_AND_STORAGE                                            string = "MOVING_AND_STORAGE"
	MULTI_LEVEL_MARKETING                                         string = "MULTI_LEVEL_MARKETING"
	MUSIC_CDS_CASSETTES_AND_ALBUMS                                string = "MUSIC_CDS_CASSETTES_AND_ALBUMS"
	MUSIC_STORE_INSTRUMENTS_AND_SHEET_MUSIC                       string = "MUSIC_STORE_INSTRUMENTS_AND_SHEET_MUSIC"
	NETWORKING                                                    string = "NETWORKING"
	NEW_AGE                                                       string = "NEW_AGE"
	NEW_PARTS_AND_SUPPLIES_MOTOR_VEHICLE                          string = "NEW_PARTS_AND_SUPPLIES_MOTOR_VEHICLE"
	NEWS_DEALERS_AND_NEWSTANDS                                    string = "NEWS_DEALERS_AND_NEWSTANDS"
	NON_DURABLE_GOODS                                             string = "NON_DURABLE_GOODS"
	NON_FICTION                                                   string = "NON_FICTION"
	NON_PROFIT_POLITICAL_AND_RELIGION                             string = "NON_PROFIT_POLITICAL_AND_RELIGION"
	NONPROFIT                                                     string = "NONPROFIT"
	NOVELTIES                                                     string = "NOVELTIES"
	OEM_SOFTWARE                                                  string = "OEM_SOFTWARE"
	OFFICE_SUPPLIES_AND_EQUIPMENT                                 string = "OFFICE_SUPPLIES_AND_EQUIPMENT"
	ONLINE_DATING                                                 string = "ONLINE_DATING"
	ONLINE_GAMING                                                 string = "ONLINE_GAMING"
	ONLINE_GAMING_CURRENCY                                        string = "ONLINE_GAMING_CURRENCY"
	ONLINE_SERVICES                                               string = "ONLINE_SERVICES"
	OOUTBOUND_TELEMARKETING_MERCH                                 string = "OOUTBOUND_TELEMARKETING_MERCH"
	OPHTHALMOLOGISTS_OPTOMETRIST                                  string = "OPHTHALMOLOGISTS_OPTOMETRIST"
	OPTICIANS_AND_DISPENSING                                      string = "OPTICIANS_AND_DISPENSING"
	ORTHOPEDIC_GOODS_PROSTHETICS                                  string = "ORTHOPEDIC_GOODS_PROSTHETICS"
	OSTEOPATHS                                                    string = "OSTEOPATHS"
	OTHER                                                         string = "OTHER"
	PACKAGE_TOUR_OPERATORS                                        string = "PACKAGE_TOUR_OPERATORS"
	PAINTBALL                                                     string = "PAINTBALL"
	PAINTS_VARNISHES_AND_SUPPLIES                                 string = "PAINTS_VARNISHES_AND_SUPPLIES"
	PARKING_LOTS_AND_GARAGES                                      string = "PARKING_LOTS_AND_GARAGES"
	PARTS_AND_ACCESSORIES                                         string = "PARTS_AND_ACCESSORIES"
	PAWN_SHOPS                                                    string = "PAWN_SHOPS"
	PAYCHECK_LENDER_OR_CASH_ADVANCE                               string = "PAYCHECK_LENDER_OR_CASH_ADVANCE"
	PERIPHERALS                                                   string = "PERIPHERALS"
	PERSONALIZED_GIFTS                                            string = "PERSONALIZED_GIFTS"
	PET_SHOPS_PET_FOOD_AND_SUPPLIES                               string = "PET_SHOPS_PET_FOOD_AND_SUPPLIES"
	PETROLEUM_AND_PETROLEUM_PRODUCTS                              string = "PETROLEUM_AND_PETROLEUM_PRODUCTS"
	PETS_AND_ANIMALS                                              string = "PETS_AND_ANIMALS"
	PHOTOFINISHING_LABORATORIES_PHOTO_DEVELOPING                  string = "PHOTOFINISHING_LABORATORIES_PHOTO_DEVELOPING"
	PHOTOGRAPHIC_STUDIOS_PORTRAITS                                string = "PHOTOGRAPHIC_STUDIOS_PORTRAITS"
	PHOTOGRAPHY                                                   string = "PHOTOGRAPHY"
	PHYSICAL_GOOD                                                 string = "PHYSICAL_GOOD"
	PICTURE_VIDEO_PRODUCTION                                      string = "PICTURE_VIDEO_PRODUCTION"
	PIECE_GOODS_NOTIONS_AND_OTHER_DRY_GOODS                       string = "PIECE_GOODS_NOTIONS_AND_OTHER_DRY_GOODS"
	PLANTS_AND_SEEDS                                              string = "PLANTS_AND_SEEDS"
	PLUMBING_AND_HEATING_EQUIPMENTS_AND_SUPPLIES                  string = "PLUMBING_AND_HEATING_EQUIPMENTS_AND_SUPPLIES"
	POLICE_RELATED_ITEMS                                          string = "POLICE_RELATED_ITEMS"
	POLITICAL_ORGANIZATIONS                                       string = "POLITICAL_ORGANIZATIONS"
	POSTAL_SERVICES_GOVERNMENT_ONLY                               string = "POSTAL_SERVICES_GOVERNMENT_ONLY"
	POSTERS                                                       string = "POSTERS"
	PREPAID_AND_STORED_VALUE_CARDS                                string = "PREPAID_AND_STORED_VALUE_CARDS"
	PRESCRIPTION_DRUGS                                            string = "PRESCRIPTION_DRUGS"
	PROMOTIONAL_ITEMS                                             string = "PROMOTIONAL_ITEMS"
	PUBLIC_WAREHOUSING_AND_STORAGE                                string = "PUBLIC_WAREHOUSING_AND_STORAGE"
	PUBLISHING_AND_PRINTING                                       string = "PUBLISHING_AND_PRINTING"
	PUBLISHING_SERVICES                                           string = "PUBLISHING_SERVICES"
	RADAR_DECTORS                                                 string = "RADAR_DECTORS"
	RADIO_TELEVISION_AND_STEREO_REPAIR                            string = "RADIO_TELEVISION_AND_STEREO_REPAIR"
	REAL_ESTATE                                                   string = "REAL_ESTATE"
	REAL_ESTATE_AGENT                                             string = "REAL_ESTATE_AGENT"
	REAL_ESTATE_AGENTS_AND_MANAGERS_RENTALS                       string = "REAL_ESTATE_AGENTS_AND_MANAGERS_RENTALS"
	RELIGION_AND_SPIRITUALITY_FOR_PROFIT                          string = "RELIGION_AND_SPIRITUALITY_FOR_PROFIT"
	RELIGIOUS                                                     string = "RELIGIOUS"
	RELIGIOUS_ORGANIZATIONS                                       string = "RELIGIOUS_ORGANIZATIONS"
	REMITTANCE                                                    string = "REMITTANCE"
	RENTAL_PROPERTY_MANAGEMENT                                    string = "RENTAL_PROPERTY_MANAGEMENT"
	RESIDENTIAL                                                   string = "RESIDENTIAL"
	RETAIL                                                        string = "RETAIL"
	RETAIL_FINE_JEWELRY_AND_WATCHES                               string = "RETAIL_FINE_JEWELRY_AND_WATCHES"
	REUPHOLSTERY_AND_FURNITURE_REPAIR                             string = "REUPHOLSTERY_AND_FURNITURE_REPAIR"
	RINGS                                                         string = "RINGS"
	ROOFING_SIDING_SHEET_METAL                                    string = "ROOFING_SIDING_SHEET_METAL"
	RUGS_AND_CARPETS                                              string = "RUGS_AND_CARPETS"
	SCHOOLS_AND_COLLEGES                                          string = "SCHOOLS_AND_COLLEGES"
	SCIENCE_FICTION                                               string = "SCIENCE_FICTION"
	SCRAPBOOKING                                                  string = "SCRAPBOOKING"
	SCULPTURES                                                    string = "SCULPTURES"
	SECURITIES_BROKERS_AND_DEALERS                                string = "SECURITIES_BROKERS_AND_DEALERS"
	SECURITY_AND_SURVEILLANCE                                     string = "SECURITY_AND_SURVEILLANCE"
	SECURITY_AND_SURVEILLANCE_EQUIPMENT                           string = "SECURITY_AND_SURVEILLANCE_EQUIPMENT"
	SECURITY_BROKERS_AND_DEALERS                                  string = "SECURITY_BROKERS_AND_DEALERS"
	SEMINARS                                                      string = "SEMINARS"
	SERVICE_STATIONS                                              string = "SERVICE_STATIONS"
	SERVICES                                                      string = "SERVICES"
	SEWING_NEEDLEWORK_FABRIC_AND_PIECE_GOODS_STORES               string = "SEWING_NEEDLEWORK_FABRIC_AND_PIECE_GOODS_STORES"
	SHIPPING_AND_PACKING                                          string = "SHIPPING_AND_PACKING"
	SHOE_REPAIR_HAT_CLEANING                                      string = "SHOE_REPAIR_HAT_CLEANING"
	SHOE_STORES                                                   string = "SHOE_STORES"
	SHOES                                                         string = "SHOES"
	SNOWMOBILE_DEALERS                                            string = "SNOWMOBILE_DEALERS"
	SOFTWARE                                                      string = "SOFTWARE"
	SPECIALTY_AND_MISC                                            string = "SPECIALTY_AND_MISC"
	SPECIALTY_CLEANING_POLISHING_AND_SANITATION_PREPARATIONS      string = "SPECIALTY_CLEANING_POLISHING_AND_SANITATION_PREPARATIONS"
	SPECIALTY_OR_RARE_PETS                                        string = "SPECIALTY_OR_RARE_PETS"
	SPORT_GAMES_AND_TOYS                                          string = "SPORT_GAMES_AND_TOYS"
	SPORTING_AND_RECREATIONAL_CAMPS                               string = "SPORTING_AND_RECREATIONAL_CAMPS"
	SPORTING_GOODS                                                string = "SPORTING_GOODS"
	SPORTS_AND_OUTDOORS                                           string = "SPORTS_AND_OUTDOORS"
	SPORTS_AND_RECREATION                                         string = "SPORTS_AND_RECREATION"
	STAMP_AND_COIN                                                string = "STAMP_AND_COIN"
	STATIONARY_PRINTING_AND_WRITING_PAPER                         string = "STATIONARY_PRINTING_AND_WRITING_PAPER"
	STENOGRAPHIC_AND_SECRETARIAL_SUPPORT_SERVICES                 string = "STENOGRAPHIC_AND_SECRETARIAL_SUPPORT_SERVICES"
	STOCKS_BONDS_SECURITIES_AND_RELATED_CERTIFICATES              string = "STOCKS_BONDS_SECURITIES_AND_RELATED_CERTIFICATES"
	STORED_VALUE_CARDS                                            string = "STORED_VALUE_CARDS"
	SUPPLIES                                                      string = "SUPPLIES"
	SUPPLIES_AND_TOYS                                             string = "SUPPLIES_AND_TOYS"
	SURVEILLANCE_EQUIPMENT                                        string = "SURVEILLANCE_EQUIPMENT"
	SWIMMING_POOLS_AND_SPAS                                       string = "SWIMMING_POOLS_AND_SPAS"
	SWIMMING_POOLS_SALES_SUPPLIES_SERVICES                        string = "SWIMMING_POOLS_SALES_SUPPLIES_SERVICES"
	TAILORS_AND_ALTERATIONS                                       string = "TAILORS_AND_ALTERATIONS"
	TAX_PAYMENTS                                                  string = "TAX_PAYMENTS"
	TAX_PAYMENTS_GOVERNMENT_AGENCIES                              string = "TAX_PAYMENTS_GOVERNMENT_AGENCIES"
	TAXICABS_AND_LIMOUSINES                                       string = "TAXICABS_AND_LIMOUSINES"
	TELECOMMUNICATION_SERVICES                                    string = "TELECOMMUNICATION_SERVICES"
	TELEPHONE_CARDS                                               string = "TELEPHONE_CARDS"
	TELEPHONE_EQUIPMENT                                           string = "TELEPHONE_EQUIPMENT"
	TELEPHONE_SERVICES                                            string = "TELEPHONE_SERVICES"
	THEATER                                                       string = "THEATER"
	TIRE_RETREADING_AND_REPAIR                                    string = "TIRE_RETREADING_AND_REPAIR"
	TOLL_OR_BRIDGE_FEES                                           string = "TOLL_OR_BRIDGE_FEES"
	TOOLS_AND_EQUIPMENT                                           string = "TOOLS_AND_EQUIPMENT"
	TOURIST_ATTRACTIONS_AND_EXHIBITS                              string = "TOURIST_ATTRACTIONS_AND_EXHIBITS"
	TOWING_SERVICE                                                string = "TOWING_SERVICE"
	TOYS_AND_GAMES                                                string = "TOYS_AND_GAMES"
	TRADE_AND_VOCATIONAL_SCHOOLS                                  string = "TRADE_AND_VOCATIONAL_SCHOOLS"
	TRADEMARK_INFRINGEMENT                                        string = "TRADEMARK_INFRINGEMENT"
	TRAILER_PARKS_AND_CAMPGROUNDS                                 string = "TRAILER_PARKS_AND_CAMPGROUNDS"
	TRAINING_SERVICES                                             string = "TRAINING_SERVICES"
	TRANSPORTATION_SERVICES                                       string = "TRANSPORTATION_SERVICES"
	TRAVEL                                                        string = "TRAVEL"
	TRUCK_AND_UTILITY_TRAILER_RENTALS                             string = "TRUCK_AND_UTILITY_TRAILER_RENTALS"
	TRUCK_STOP                                                    string = "TRUCK_STOP"
	TYPESETTING_PLATE_MAKING_AND_RELATED_SERVICES                 string = "TYPESETTING_PLATE_MAKING_AND_RELATED_SERVICES"
	USED_MERCHANDISE_AND_SECONDHAND_STORES                        string = "USED_MERCHANDISE_AND_SECONDHAND_STORES"
	USED_PARTS_MOTOR_VEHICLE                                      string = "USED_PARTS_MOTOR_VEHICLE"
	UTILITIES                                                     string = "UTILITIES"
	UTILITIES_ELECTRIC_GAS_WATER_SANITARY                         string = "UTILITIES_ELECTRIC_GAS_WATER_SANITARY"
	VARIETY_STORES                                                string = "VARIETY_STORES"
	VEHICLE_SALES                                                 string = "VEHICLE_SALES"
	VEHICLE_SERVICE_AND_ACCESSORIES                               string = "VEHICLE_SERVICE_AND_ACCESSORIES"
	VIDEO_EQUIPMENT                                               string = "VIDEO_EQUIPMENT"
	VIDEO_GAME_ARCADES_ESTABLISH                                  string = "VIDEO_GAME_ARCADES_ESTABLISH"
	VIDEO_GAMES_AND_SYSTEMS                                       string = "VIDEO_GAMES_AND_SYSTEMS"
	VIDEO_TAPE_RENTAL_STORES                                      string = "VIDEO_TAPE_RENTAL_STORES"
	VINTAGE_AND_COLLECTIBLE_VEHICLES                              string = "VINTAGE_AND_COLLECTIBLE_VEHICLES"
	VINTAGE_AND_COLLECTIBLES                                      string = "VINTAGE_AND_COLLECTIBLES"
	VITAMINS_AND_SUPPLEMENTS                                      string = "VITAMINS_AND_SUPPLEMENTS"
	VOCATIONAL_AND_TRADE_SCHOOLS                                  string = "VOCATIONAL_AND_TRADE_SCHOOLS"
	WATCH_CLOCK_AND_JEWELRY_REPAIR                                string = "WATCH_CLOCK_AND_JEWELRY_REPAIR"
	WEB_HOSTING_AND_DESIGN                                        string = "WEB_HOSTING_AND_DESIGN"
	WELDING_REPAIR                                                string = "WELDING_REPAIR"
	WHOLESALE_CLUBS                                               string = "WHOLESALE_CLUBS"
	WHOLESALE_FLORIST_SUPPLIERS                                   string = "WHOLESALE_FLORIST_SUPPLIERS"
	WHOLESALE_PRESCRIPTION_DRUGS                                  string = "WHOLESALE_PRESCRIPTION_DRUGS"
	WILDLIFE_PRODUCTS                                             string = "WILDLIFE_PRODUCTS"
	WIRE_TRANSFER                                                 string = "WIRE_TRANSFER"
	WIRE_TRANSFER_AND_MONEY_ORDER                                 string = "WIRE_TRANSFER_AND_MONEY_ORDER"
	WOMEN_ACCESSORY_SPECIALITY                                    string = "WOMEN_ACCESSORY_SPECIALITY"
	WOMEN_CLOTHING                                                string = "WOMEN_CLOTHING"
)

// Possible values for Plan.Status
const (
	P_CREATED  string = "CREATED"
	P_INACTIVE string = "INACTIVE"
	P_ACTIVE   string = "ACTIVE"
)

// Possible values for Plan.Status
const (
	S_APPROVAL_PENDING string = "APPROVAL_PENDING"
	S_APPROVED         string = "APPROVED"
	S_ACTIVE           string = "ACTIVE"
	S_SUSPENDED        string = "SUSPENDED"
	S_CANCELLED        string = "CANCELLED"
	S_EXPIRED          string = "EXPIRED"
)

// Possible values for BillingCycle.TenureType and CycleExecution.TenureType
const (
	REGULAR string = "REGULAR"
	TRIAL   string = "TRIAL"
)

// Possible values for Frequency.IntervalUnit
const (
	DAY   string = "DAY"
	WEEK  string = "WEEK"
	MONTH string = "MONTH"
	YEAR  string = "YEAR"
)

// Possible values for PaymentPreferences.SetupFeeFailureAction
const (
	CONTINUE string = "CONTINUE"
	CANCEL   string = "CANCEL"
)

// Possible values for PaymentMethod.PayerSelected
const (
	PAYPAL string = "PAYPAL"
)

// Possible values for PaymentMethod.PayeePreferred
const (
	UNRESTRICTED               string = "UNRESTRICTED"
	IMMEDIATE_PAYMENT_REQUIRED string = "IMMEDIATE_PAYMENT_REQUIRED"
)

// Possible values for PaymentMethod.Category
const (
	CUSTOMER_PRESENT_SINGLE_PURCHASE string = "CUSTOMER_PRESENT_SINGLE_PURCHASE"
	CUSTOMER_NOT_PRESENT_RECURRING   string = "CUSTOMER_NOT_PRESENT_RECURRING"
	CUSTOMER_PRESENT_RECURRING_FIRST string = "CUSTOMER_PRESENT_RECURRING_FIRST"
	CUSTOMER_PRESENT_UNSCHEDULED     string = "CUSTOMER_PRESENT_UNSCHEDULED"
	CUSTOMER_NOT_PRESENT_UNSCHEDULED string = "CUSTOMER_NOT_PRESENT_UNSCHEDULED"
	MAIL_ORDER_TELEPHONE_ORDER       string = "MAIL_ORDER_TELEPHONE_ORDER"
)

// Possible values for CardResponseWithBillingAddress.Brand
const (
	VISA            string = "VISA"
	MASTERCARD      string = "MASTERCARD"
	DISCOVER        string = "DISCOVER"
	AMEX            string = "AMEX"
	SOLO            string = "SOLO"
	JCB             string = "JCB"
	STAR            string = "STAR"
	DELTA           string = "DELTA"
	SWITCH          string = "SWITCH"
	MAESTRO         string = "MAESTRO"
	CB_NATIONALE    string = "CB_NATIONALE"
	CONFIGOGA       string = "CONFIGOGA"
	CONFIDIS        string = "CONFIDIS"
	ELECTRON        string = "ELECTRON"
	CETELEM         string = "CETELEM"
	CHINA_UNION_PAY string = "CHINA_UNION_PAY"
)

// Possible values for CardResponseWithBillingAddress.Type
const (
	CREDIT  string = "CREDIT"
	DEBIT   string = "DEBIT"
	PREPAID string = "PREPAID"
	UNKNOWN string = "UNKNOWN"
)

// Possible values for FailedPaymentDetails.ReasonCode
const (
	PAYMENT_DENIED                       string = "PAYMENT_DENIED"
	INTERNAL_SERVER_ERROR                string = "INTERNAL_SERVER_ERROR"
	PAYEE_ACCOUNT_RESTRICTED             string = "PAYEE_ACCOUNT_RESTRICTED"
	PAYER_ACCOUNT_RESTRICTED             string = "PAYER_ACCOUNT_RESTRICTED"
	PAYER_CANNOT_PAY                     string = "PAYER_CANNOT_PAY"
	SENDING_LIMIT_EXCEEDED               string = "SENDING_LIMIT_EXCEEDED"
	TRANSACTION_RECEIVING_LIMIT_EXCEEDED string = "TRANSACTION_RECEIVING_LIMIT_EXCEEDED"
	CURRENCY_MISMATCH                    string = "CURRENCY_MISMATCH"
)

// Possible values for CaptureAuthorizedPaymentOnSubscriptionRequest.CaptureType
const (
	OUTSTANDING_BALANCE string = "OUTSTANDING_BALANCE"
)

type (
	// JSONTime overrides MarshalJson method to format in ISO8601
	JSONTime time.Time

	// Address struct
	Address struct {
		Line1       string `json:"line1"`
		Line2       string `json:"line2, omitempty"`
		City        string `json:"city"`
		CountryCode string `json:"country_code"`
		PostalCode  string `json:"postal_code, omitempty"`
		State       string `json:"state, omitempty"`
		Phone       string `json:"phone, omitempty"`
	}

	// AgreementDetails struct
	AgreementDetails struct {
		OutstandingBalance AmountPayout `json:"outstanding_balance"`
		CyclesRemaining    int          `json:"cycles_remaining,string"`
		CyclesCompleted    int          `json:"cycles_completed,string"`
		NextBillingDate    time.Time    `json:"next_billing_date"`
		LastPaymentDate    time.Time    `json:"last_payment_date"`
		LastPaymentAmount  AmountPayout `json:"last_payment_amount"`
		FinalPaymentDate   time.Time    `json:"final_payment_date"`
		FailedPaymentCount int          `json:"failed_payment_count,string"`
	}

	// Amount struct
	Amount struct {
		Currency string  `json:"currency"`
		Total    string  `json:"total"`
		Details  Details `json:"details, omitempty"`
	}

	// AmountPayout struct
	AmountPayout struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	}

	// ApplicationContext represents the application context, which customizes
	// the payer experience during the subscription approval process with PayPal.
	// ShippingPreference represents the location from which the shipping address is derived.
	// The possible values are:
	// ------------------------------------------------------------------------------------------------------------
	// | GET_FROM_FILE        | Get the customer-provided shipping address on the PayPal site.                    |
	// | NO_SHIPPING          | Redacts the shipping address from the PayPal site. Recommended for digital goods. |
	// | SET_PROVIDED_ADDRESS | Get the merchant-provided address. The customer cannot change this                |
	// |                      | address on the PayPal site. If merchant does not pass an address, customer        |
	// |                      | can choose the address on PayPal pages.                                           |
	// ------------------------------------------------------------------------------------------------------------
	// UserAction configures the label name to Continue or Subscribe Now for subscription consent experience.
	// The possible values are:
	// ------------------------------------------------------------------------------------------------------
	// | CONTINUE      | After you redirect the customer to the PayPal subscription consent page,			|
	// |               | a Continue button appears. Use this option when you want to control the activation |
	// |               | of the subscription and do not want PayPal to activate the subscription.			|
	// | SUBSCRIBE_NOW | After you redirect the customer to the PayPal subscription consent page, 			|
	// |               | a Subscribe Now button appears. Use this option when you want PayPal to activate   |
	// |               | the subscription.																	|
	// ------------------------------------------------------------------------------------------------------
	ApplicationContext struct {
		BrandName          string         `json:"brand_name, omitempty"`
		Locale             string         `json:"locale, omitempty"`
		LandingPage        string         `json:"landing_page, omitempty"`
		ShippingPreference string         `json:"shipping_preference, omitempty"` //default: GET_FROM_FILE
		UserAction         string         `json:"user_action, omitempty"`         //default: SUBSCRIBE_NOW
		PaymentMethod      *PaymentMethod `json:"payment_method, omitempty"`
		ReturnURL          string         `json:"return_url"`
		CancelURL          string         `json:"cancel_url"`
	}

	// Authorization struct
	Authorization struct {
		ID               string                `json:"id, omitempty"`
		CustomID         string                `json:"custom_id, omitempty"`
		InvoiceID        string                `json:"invoice_id, omitempty"`
		Status           string                `json:"status, omitempty"`
		StatusDetails    *CaptureStatusDetails `json:"status_details, omitempty"`
		Amount           *PurchaseUnitAmount   `json:"amount, omitempty"`
		SellerProtection *SellerProtection     `json:"seller_protection, omitempty"`
		CreateTime       *time.Time            `json:"create_time, omitempty"`
		UpdateTime       *time.Time            `json:"update_time, omitempty"`
		ExpirationTime   *time.Time            `json:"expiration_time, omitempty"`
		Links            []Link                `json:"links, omitempty"`
	}

	// AuthorizeOrderResponse .
	AuthorizeOrderResponse struct {
		CreateTime    *time.Time             `json:"create_time, omitempty"`
		UpdateTime    *time.Time             `json:"update_time, omitempty"`
		ID            string                 `json:"id, omitempty"`
		Status        string                 `json:"status, omitempty"`
		Intent        string                 `json:"intent, omitempty"`
		PurchaseUnits []PurchaseUnitRequest  `json:"purchase_units, omitempty"`
		Payer         *PayerWithNameAndPhone `json:"payer, omitempty"`
	}

	// AuthorizeOrderRequest - https://developer.paypal.com/docs/api/orders/v2/#orders_authorize
	AuthorizeOrderRequest struct {
		PaymentSource      *PaymentSource     `json:"payment_source, omitempty"`
		ApplicationContext ApplicationContext `json:"application_context, omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-platform_fee
	PlatformFee struct {
		Amount *Money          `json:"amount, omitempty"`
		Payee  *PayeeForOrders `json:"payee, omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-payment_instruction
	PaymentInstruction struct {
		PlatformFees     []PlatformFee `json:"platform_fees, omitempty"`
		DisbursementMode string        `json:"disbursement_mode, omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#authorizations_capture
	PaymentCaptureRequest struct {
		InvoiceID      string `json:"invoice_id, omitempty"`
		NoteToPayer    string `json:"note_to_payer, omitempty"`
		SoftDescriptor string `json:"soft_descriptor, omitempty"`
		Amount         *Money `json:"amount, omitempty"`
		FinalCapture   bool   `json:"final_capture, omitempty"`
	}

	SellerProtection struct {
		Status            string   `json:"status, omitempty"`
		DisputeCategories []string `json:"dispute_categories, omitempty"`
	}

	// https://developer.paypal.com/docs/api/payments/v2/#definition-capture_status_details
	CaptureStatusDetails struct {
		Reason string `json:"reason, omitempty"`
	}

	PaymentCaptureResponse struct {
		Status           string                `json:"status, omitempty"`
		StatusDetails    *CaptureStatusDetails `json:"status_details, omitempty"`
		ID               string                `json:"id, omitempty"`
		Amount           *Money                `json:"amount, omitempty"`
		InvoiceID        string                `json:"invoice_id, omitempty"`
		FinalCapture     bool                  `json:"final_capture, omitempty"`
		DisbursementMode string                `json:"disbursement_mode, omitempty"`
		Links            []Link                `json:"links, omitempty"`
	}

	// CaptureOrderRequest - https://developer.paypal.com/docs/api/orders/v2/#orders_capture
	CaptureOrderRequest struct {
		PaymentSource *PaymentSource `json:"payment_source"`
	}

	// BatchHeader struct
	BatchHeader struct {
		Amount            *AmountPayout      `json:"amount, omitempty"`
		Fees              *AmountPayout      `json:"fees, omitempty"`
		PayoutBatchID     string             `json:"payout_batch_id, omitempty"`
		BatchStatus       string             `json:"batch_status, omitempty"`
		TimeCreated       *time.Time         `json:"time_created, omitempty"`
		TimeCompleted     *time.Time         `json:"time_completed, omitempty"`
		SenderBatchHeader *SenderBatchHeader `json:"sender_batch_header, omitempty"`
	}

	// BillingAgreement struct
	BillingAgreement struct {
		Name                        string               `json:"name, omitempty"`
		Description                 string               `json:"description, omitempty"`
		StartDate                   JSONTime             `json:"start_date, omitempty"`
		Plan                        BillingPlan          `json:"plan, omitempty"`
		Payer                       Payer                `json:"payer, omitempty"`
		ShippingAddress             *ShippingAddress     `json:"shipping_address, omitempty"`
		OverrideMerchantPreferences *MerchantPreferences `json:"override_merchant_preferences, omitempty"`
	}

	// BillingPlan struct
	BillingPlan struct {
		ID                  string               `json:"id, omitempty"`
		Name                string               `json:"name, omitempty"`
		Description         string               `json:"description, omitempty"`
		Type                string               `json:"type, omitempty"`
		PaymentDefinitions  []PaymentDefinition  `json:"payment_definitions, omitempty"`
		MerchantPreferences *MerchantPreferences `json:"merchant_preferences, omitempty"`
	}

	// Capture struct
	Capture struct {
		Amount         *Amount    `json:"amount, omitempty"`
		IsFinalCapture bool       `json:"is_final_capture"`
		CreateTime     *time.Time `json:"create_time, omitempty"`
		UpdateTime     *time.Time `json:"update_time, omitempty"`
		State          string     `json:"state, omitempty"`
		ParentPayment  string     `json:"parent_payment, omitempty"`
		ID             string     `json:"id, omitempty"`
		Links          []Link     `json:"links, omitempty"`
	}

	// ChargeModel struct
	ChargeModel struct {
		Type   string       `json:"type, omitempty"`
		Amount AmountPayout `json:"amount, omitempty"`
	}

	// Client represents a Paypal REST API Client
	Client struct {
		sync.Mutex
		Client               *http.Client
		ClientID             string
		Secret               string
		APIBase              string
		Log                  io.Writer // If user set log file name all requests will be logged there
		Token                *TokenResponse
		tokenExpiresAt       time.Time
		returnRepresentation bool
	}

	// CreditCard struct
	CreditCard struct {
		ID                 string   `json:"id, omitempty"`
		PayerID            string   `json:"payer_id, omitempty"`
		ExternalCustomerID string   `json:"external_customer_id, omitempty"`
		Number             string   `json:"number"`
		Type               string   `json:"type"`
		ExpireMonth        string   `json:"expire_month"`
		ExpireYear         string   `json:"expire_year"`
		CVV2               string   `json:"cvv2, omitempty"`
		FirstName          string   `json:"first_name, omitempty"`
		LastName           string   `json:"last_name, omitempty"`
		BillingAddress     *Address `json:"billing_address, omitempty"`
		State              string   `json:"state, omitempty"`
		ValidUntil         string   `json:"valid_until, omitempty"`
	}

	// CreditCards GET /v1/vault/credit-cards
	CreditCards struct {
		Items      []CreditCard `json:"items"`
		Links      []Link       `json:"links"`
		TotalItems int          `json:"total_items"`
		TotalPages int          `json:"total_pages"`
	}

	// CreditCardToken struct
	CreditCardToken struct {
		CreditCardID string `json:"credit_card_id"`
		PayerID      string `json:"payer_id, omitempty"`
		Last4        string `json:"last4, omitempty"`
		ExpireYear   string `json:"expire_year, omitempty"`
		ExpireMonth  string `json:"expire_month, omitempty"`
	}

	// CreditCardsFilter struct
	CreditCardsFilter struct {
		PageSize int
		Page     int
	}

	// CreditCardField PATCH /v1/vault/credit-cards/credit_card_id
	CreditCardField struct {
		Operation string `json:"op"`
		Path      string `json:"path"`
		Value     string `json:"value"`
	}

	// Currency struct
	Currency struct {
		Currency string `json:"currency, omitempty"`
		Value    string `json:"value, omitempty"`
	}

	// Details structure used in Amount structures as optional value
	Details struct {
		Subtotal         string `json:"subtotal, omitempty"`
		Shipping         string `json:"shipping, omitempty"`
		Tax              string `json:"tax, omitempty"`
		HandlingFee      string `json:"handling_fee, omitempty"`
		ShippingDiscount string `json:"shipping_discount, omitempty"`
		Insurance        string `json:"insurance, omitempty"`
		GiftWrap         string `json:"gift_wrap, omitempty"`
	}

	// ErrorResponseDetail struct
	ErrorResponseDetail struct {
		Field string `json:"field"`
		Issue string `json:"issue"`
		Links []Link `json:"link"`
	}

	// ErrorResponse https://developer.paypal.com/docs/api/errors/
	ErrorResponse struct {
		Response        *http.Response        `json:"-"`
		Name            string                `json:"name"`
		DebugID         string                `json:"debug_id"`
		Message         string                `json:"message"`
		InformationLink string                `json:"information_link"`
		Details         []ErrorResponseDetail `json:"details"`
	}

	// ExecuteAgreementResponse struct
	ExecuteAgreementResponse struct {
		ID               string           `json:"id"`
		State            string           `json:"state"`
		Description      string           `json:"description, omitempty"`
		Payer            Payer            `json:"payer"`
		Plan             BillingPlan      `json:"plan"`
		StartDate        time.Time        `json:"start_date"`
		ShippingAddress  ShippingAddress  `json:"shipping_address"`
		AgreementDetails AgreementDetails `json:"agreement_details"`
		Links            []Link           `json:"links"`
	}

	// ExecuteResponse struct
	ExecuteResponse struct {
		ID           string        `json:"id"`
		Links        []Link        `json:"links"`
		State        string        `json:"state"`
		Payer        PaymentPayer  `json:"payer"`
		Transactions []Transaction `json:"transactions, omitempty"`
	}

	// FundingInstrument struct
	FundingInstrument struct {
		CreditCard      *CreditCard      `json:"credit_card, omitempty"`
		CreditCardToken *CreditCardToken `json:"credit_card_token, omitempty"`
	}

	// Item struct
	Item struct {
		Name        string `json:"name"`
		UnitAmount  *Money `json:"unit_amount, omitempty"`
		Tax         *Money `json:"tax, omitempty"`
		Quantity    string `json:"quantity"`
		Description string `json:"description, omitempty"`
		SKU         string `json:"sku, omitempty"`
		Category    string `json:"category, omitempty"`
	}

	// ItemList struct
	ItemList struct {
		Items           []Item           `json:"items, omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address, omitempty"`
	}

	// Link struct
	Link struct {
		Href        string `json:"href"`
		Rel         string `json:"rel, omitempty"`
		Method      string `json:"method, omitempty"`
		Description string `json:"description, omitempty"`
		Enctype     string `json:"enctype, omitempty"`
	}

	// PurchaseUnitAmount struct
	PurchaseUnitAmount struct {
		Currency  string                       `json:"currency_code"`
		Value     string                       `json:"value"`
		Breakdown *PurchaseUnitAmountBreakdown `json:"breakdown, omitempty"`
	}

	// PurchaseUnitAmountBreakdown struct
	PurchaseUnitAmountBreakdown struct {
		ItemTotal        *Money `json:"item_total, omitempty"`
		Shipping         *Money `json:"shipping, omitempty"`
		Handling         *Money `json:"handling, omitempty"`
		TaxTotal         *Money `json:"tax_total, omitempty"`
		Insurance        *Money `json:"insurance, omitempty"`
		ShippingDiscount *Money `json:"shipping_discount, omitempty"`
		Discount         *Money `json:"discount, omitempty"`
	}

	// Money represents the amount. For regular pricing, it is limited to a 20% increase
	// from the current amount and the change is applicable for both existing and future subscriptions. For trial period pricing,
	// there is no limit or constraint in changing the amount and the change is applicable only on future subscriptions.
	// The value, which might be:
	// An integer for currencies like JPY that are not typically fractional.
	// A decimal fraction for currencies like TND that are subdivided into thousandths.
	// For the required number of decimal places for a currency code, see https://developer.paypal.com/docs/api/reference/currency-codes/.
	// https://developer.paypal.com/docs/api/orders/v2/#definition-money
	Money struct {
		Currency string `json:"currency_code"`
		Value    string `json:"value"`
	}

	// PurchaseUnit struct
	PurchaseUnit struct {
		ReferenceID string              `json:"reference_id"`
		Amount      *PurchaseUnitAmount `json:"amount, omitempty"`
	}

	// TaxInfo used for orders.
	TaxInfo struct {
		TaxID     string `json:"tax_id, omitempty"`
		TaxIDType string `json:"tax_id_type, omitempty"`
	}

	// PhoneWithTypeNumber struct for PhoneWithType
	PhoneWithTypeNumber struct {
		NationalNumber string `json:"national_number, omitempty"`
	}

	// PhoneWithType struct used for orders
	PhoneWithType struct {
		PhoneType   string               `json:"phone_type, omitempty"`
		PhoneNumber *PhoneWithTypeNumber `json:"phone_number, omitempty"`
	}

	// CreateOrderPayerName create order payer name
	CreateOrderPayerName struct {
		GivenName string `json:"given_name, omitempty"`
		Surname   string `json:"surname, omitempty"`
	}

	// CreateOrderPayer used with create order requests
	CreateOrderPayer struct {
		Name         *CreateOrderPayerName          `json:"name, omitempty"`
		EmailAddress string                         `json:"email_address, omitempty"`
		PayerID      string                         `json:"payer_id, omitempty"`
		Phone        *PhoneWithType                 `json:"phone, omitempty"`
		BirthDate    string                         `json:"birth_date, omitempty"`
		TaxInfo      *TaxInfo                       `json:"tax_info, omitempty"`
		Address      *ShippingDetailAddressPortable `json:"address, omitempty"`
	}

	// PurchaseUnitRequest struct
	PurchaseUnitRequest struct {
		ReferenceID    string              `json:"reference_id, omitempty"`
		Amount         *PurchaseUnitAmount `json:"amount"`
		Payee          *PayeeForOrders     `json:"payee, omitempty"`
		Description    string              `json:"description, omitempty"`
		CustomID       string              `json:"custom_id, omitempty"`
		InvoiceID      string              `json:"invoice_id, omitempty"`
		SoftDescriptor string              `json:"soft_descriptor, omitempty"`
		Items          []Item              `json:"items, omitempty"`
		Shipping       *ShippingDetail     `json:"shipping, omitempty"`
	}

	// MerchantPreferences struct
	MerchantPreferences struct {
		SetupFee                *AmountPayout `json:"setup_fee, omitempty"`
		ReturnURL               string        `json:"return_url, omitempty"`
		CancelURL               string        `json:"cancel_url, omitempty"`
		AutoBillAmount          string        `json:"auto_bill_amount, omitempty"`
		InitialFailAmountAction string        `json:"initial_fail_amount_action, omitempty"`
		MaxFailAttempts         string        `json:"max_fail_attempts, omitempty"`
	}

	// Order struct
	Order struct {
		ID            string         `json:"id, omitempty"`
		Status        string         `json:"status, omitempty"`
		Intent        string         `json:"intent, omitempty"`
		PurchaseUnits []PurchaseUnit `json:"purchase_units, omitempty"`
		Links         []Link         `json:"links, omitempty"`
		CreateTime    *time.Time     `json:"create_time, omitempty"`
		UpdateTime    *time.Time     `json:"update_time, omitempty"`
	}

	// CaptureAmount struct
	CaptureAmount struct {
		ID       string              `json:"id, omitempty"`
		CustomID string              `json:"custom_id, omitempty"`
		Amount   *PurchaseUnitAmount `json:"amount, omitempty"`
	}

	// CapturedPayments has the amounts for a captured order
	CapturedPayments struct {
		Captures []CaptureAmount `json:"captures, omitempty"`
	}

	// CapturedPurchaseItem are items for a captured order
	CapturedPurchaseItem struct {
		Quantity    string `json:"quantity"`
		Name        string `json:"name"`
		SKU         string `json:"sku, omitempty"`
		Description string `json:"description, omitempty"`
	}

	// CapturedPurchaseUnit are purchase units for a captured order
	CapturedPurchaseUnit struct {
		Items    []CapturedPurchaseItem `json:"items, omitempty"`
		Payments *CapturedPayments      `json:"payments, omitempty"`
	}

	// PayerWithNameAndPhone struct
	PayerWithNameAndPhone struct {
		Name         *CreateOrderPayerName `json:"name, omitempty"`
		EmailAddress string                `json:"email_address, omitempty"`
		Phone        *PhoneWithType        `json:"phone, omitempty"`
		PayerID      string                `json:"payer_id, omitempty"`
	}

	// CaptureOrderResponse is the response for capture order
	CaptureOrderResponse struct {
		ID            string                 `json:"id, omitempty"`
		Status        string                 `json:"status, omitempty"`
		Payer         *PayerWithNameAndPhone `json:"payer, omitempty"`
		PurchaseUnits []CapturedPurchaseUnit `json:"purchase_units, omitempty"`
	}

	// Payer struct
	Payer struct {
		PaymentMethod      string              `json:"payment_method"`
		FundingInstruments []FundingInstrument `json:"funding_instruments, omitempty"`
		PayerInfo          *PayerInfo          `json:"payer_info, omitempty"`
		Status             string              `json:"payer_status, omitempty"`
	}

	// PayerInfo struct
	PayerInfo struct {
		Email           string           `json:"email, omitempty"`
		FirstName       string           `json:"first_name, omitempty"`
		LastName        string           `json:"last_name, omitempty"`
		PayerID         string           `json:"payer_id, omitempty"`
		Phone           string           `json:"phone, omitempty"`
		ShippingAddress *ShippingAddress `json:"shipping_address, omitempty"`
		TaxIDType       string           `json:"tax_id_type, omitempty"`
		TaxID           string           `json:"tax_id, omitempty"`
		CountryCode     string           `json:"country_code"`
	}

	// PaymentDefinition struct
	PaymentDefinition struct {
		ID                string        `json:"id, omitempty"`
		Name              string        `json:"name, omitempty"`
		Type              string        `json:"type, omitempty"`
		Frequency         string        `json:"frequency, omitempty"`
		FrequencyInterval string        `json:"frequency_interval, omitempty"`
		Amount            AmountPayout  `json:"amount, omitempty"`
		Cycles            string        `json:"cycles, omitempty"`
		ChargeModels      []ChargeModel `json:"charge_models, omitempty"`
	}

	// PaymentOptions struct
	PaymentOptions struct {
		AllowedPaymentMethod string `json:"allowed_payment_method, omitempty"`
	}

	// PaymentPatch PATCH /v2/payments/payment/{payment_id)
	PaymentPatch struct {
		Operation string      `json:"op"`
		Path      string      `json:"path"`
		Value     interface{} `json:"value"`
	}

	// PaymentPayer struct
	PaymentPayer struct {
		PaymentMethod string     `json:"payment_method"`
		Status        string     `json:"status, omitempty"`
		PayerInfo     *PayerInfo `json:"payer_info, omitempty"`
	}

	// PaymentResponse structure
	PaymentResponse struct {
		ID           string        `json:"id"`
		State        string        `json:"state"`
		Intent       string        `json:"intent"`
		Payer        Payer         `json:"payer"`
		Transactions []Transaction `json:"transactions"`
		Links        []Link        `json:"links"`
	}

	// PaymentSource represents the payment source definitions
	PaymentSource struct {
		Card  *PaymentSourceCard  `json:"card"`
		Token *PaymentSourceToken `json:"token"`
	}

	// PaymentSourceCard represents card details
	// SecurityCode represents the three- or four-digit security code of the card. Also known as the CVV, CVC, CVN, CVE, or CID.
	PaymentSourceCard struct {
		ID             string           `json:"id, omitempty"`
		Name           string           `json:"name, omitempty"`
		Number         string           `json:"number"`
		Expiry         string           `json:"expiry"`
		SecurityCode   string           `json:"security_code, omitempty"`
		LastDigits     string           `json:"last_digits, omitempty"`
		CardType       string           `json:"card_type, omitempty"`
		BillingAddress *AddressPortable `json:"billing_address, omitempty"`
	}

	// AddressPortable represents address details
	// More info -> https://developer.paypal.com/docs/api/subscriptions/v1/#definition-address_portable
	AddressPortable struct {
		AddressLine1 string `json:"address_line_1"`
		AddressLine2 string `json:"address_line_2"`
		AddressLine3 string `json:"address_line_3"`
		AdminArea4   string `json:"admin_area_4"`
		AdminArea3   string `json:"admin_area_3"`
		AdminArea2   string `json:"admin_area_2"`
		AdminArea1   string `json:"admin_area_1"`
		PostalCode   string `json:"postal_code"`
		CountryCode  string `json:"country_code"`
	}

	// PaymentSourceToken structure
	PaymentSourceToken struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}

	// Payout struct
	Payout struct {
		SenderBatchHeader *SenderBatchHeader `json:"sender_batch_header"`
		Items             []PayoutItem       `json:"items"`
	}

	// PayoutItem struct
	PayoutItem struct {
		RecipientType string        `json:"recipient_type"`
		Receiver      string        `json:"receiver"`
		Amount        *AmountPayout `json:"amount"`
		Note          string        `json:"note, omitempty"`
		SenderItemID  string        `json:"sender_item_id, omitempty"`
	}

	// PayoutItemResponse struct
	PayoutItemResponse struct {
		PayoutItemID      string        `json:"payout_item_id"`
		TransactionID     string        `json:"transaction_id"`
		TransactionStatus string        `json:"transaction_status"`
		PayoutBatchID     string        `json:"payout_batch_id, omitempty"`
		PayoutItemFee     *AmountPayout `json:"payout_item_fee, omitempty"`
		PayoutItem        *PayoutItem   `json:"payout_item"`
		TimeProcessed     *time.Time    `json:"time_processed, omitempty"`
		Links             []Link        `json:"links"`
		Error             ErrorResponse `json:"errors, omitempty"`
	}

	// PayoutResponse struct
	PayoutResponse struct {
		BatchHeader *BatchHeader         `json:"batch_header"`
		Items       []PayoutItemResponse `json:"items"`
		Links       []Link               `json:"links"`
	}

	// RedirectURLs struct
	RedirectURLs struct {
		ReturnURL string `json:"return_url, omitempty"`
		CancelURL string `json:"cancel_url, omitempty"`
	}

	// Refund struct
	Refund struct {
		ID            string     `json:"id, omitempty"`
		Amount        *Amount    `json:"amount, omitempty"`
		CreateTime    *time.Time `json:"create_time, omitempty"`
		State         string     `json:"state, omitempty"`
		CaptureID     string     `json:"capture_id, omitempty"`
		ParentPayment string     `json:"parent_payment, omitempty"`
		UpdateTime    *time.Time `json:"update_time, omitempty"`
	}

	// RefundResponse .
	RefundResponse struct {
		ID     string              `json:"id, omitempty"`
		Amount *PurchaseUnitAmount `json:"amount, omitempty"`
		Status string              `json:"status, omitempty"`
	}

	// Related struct
	Related struct {
		Sale          *Sale          `json:"sale, omitempty"`
		Authorization *Authorization `json:"authorization, omitempty"`
		Order         *Order         `json:"order, omitempty"`
		Capture       *Capture       `json:"capture, omitempty"`
		Refund        *Refund        `json:"refund, omitempty"`
	}

	// Sale struct
	Sale struct {
		ID                        string     `json:"id, omitempty"`
		Amount                    *Amount    `json:"amount, omitempty"`
		TransactionFee            *Currency  `json:"transaction_fee, omitempty"`
		Description               string     `json:"description, omitempty"`
		CreateTime                *time.Time `json:"create_time, omitempty"`
		State                     string     `json:"state, omitempty"`
		ParentPayment             string     `json:"parent_payment, omitempty"`
		UpdateTime                *time.Time `json:"update_time, omitempty"`
		PaymentMode               string     `json:"payment_mode, omitempty"`
		PendingReason             string     `json:"pending_reason, omitempty"`
		ReasonCode                string     `json:"reason_code, omitempty"`
		ClearingTime              string     `json:"clearing_time, omitempty"`
		ProtectionEligibility     string     `json:"protection_eligibility, omitempty"`
		ProtectionEligibilityType string     `json:"protection_eligibility_type, omitempty"`
		Links                     []Link     `json:"links, omitempty"`
	}

	// SenderBatchHeader struct
	SenderBatchHeader struct {
		EmailSubject  string `json:"email_subject"`
		SenderBatchID string `json:"sender_batch_id, omitempty"`
	}

	// ShippingAddress struct
	ShippingAddress struct {
		RecipientName string `json:"recipient_name, omitempty"`
		Type          string `json:"type, omitempty"`
		Line1         string `json:"line1"`
		Line2         string `json:"line2, omitempty"`
		City          string `json:"city"`
		CountryCode   string `json:"country_code"`
		PostalCode    string `json:"postal_code, omitempty"`
		State         string `json:"state, omitempty"`
		Phone         string `json:"phone, omitempty"`
	}

	// ShippingDetailAddressPortable represents address details
	// More info -> https://developer.paypal.com/docs/api/subscriptions/v1/#definition-shipping_detail.address_portable
	ShippingDetailAddressPortable struct {
		AddressLine1 string `json:"address_line_1, omitempty"`
		AddressLine2 string `json:"address_line_2, omitempty"`
		AdminArea1   string `json:"admin_area_1, omitempty"`
		AdminArea2   string `json:"admin_area_2, omitempty"`
		PostalCode   string `json:"postal_code, omitempty"`
		CountryCode  string `json:"country_code, omitempty"`
	}

	// ShippingDetailsName represents shipping details full name
	ShippingDetailsName struct {
		FullName string `json:"full_name, omitempty"`
	}

	// ShippingDetail represents the shipping details
	ShippingDetail struct {
		Name    *ShippingDetailsName           `json:"name, omitempty"`
		Address *ShippingDetailAddressPortable `json:"address, omitempty"`
	}

	expirationTime int64

	// TokenResponse is for API response for the /oauth2/token endpoint
	TokenResponse struct {
		RefreshToken string         `json:"refresh_token"`
		Token        string         `json:"access_token"`
		Type         string         `json:"token_type"`
		ExpiresIn    expirationTime `json:"expires_in"`
	}

	// Since it is not used i change it @gligor
	// Transaction represents details about transaction
	Transaction struct {
		ID                  string               `json:"id, omitempty"`                    //Read only
		Status              string               `json:"status, omitempty"`                //Read only
		AmountWithBreakdown *AmountWithBreakdown `json:"amount_with_breakdown, omitempty"` //Read only
		PayerName           *Name                `json:"payer_name, omitempty"`            //Read only
		PayerEmail          string               `json:"payer_email, omitempty"`           //Read only
		Time                string               `json:"time, omitempty"`                  //Read only
	}

	// AmountWithBreakdown represents the breakdown details for the amount. Includes the gross, tax, fee, and shipping amounts.
	AmountWithBreakdown struct {
		GrossAmount    *Money `json:"gross_amount"`               //Read only
		FeeAmount      *Money `json:"fee_amount, omitempty"`      //Read only
		ShippingAmount *Money `json:"shipping_amount, omitempty"` //Read only
		TaxAmount      *Money `json:"tax_amount, omitempty"`      //Read only
		NetAmount      *Money `json:"net_amount"`                 //Read only
	}

	//Payee struct
	Payee struct {
		Email string `json:"email"`
	}

	// PayeeForOrders struct
	PayeeForOrders struct {
		EmailAddress string `json:"email_address, omitempty"`
		MerchantID   string `json:"merchant_id, omitempty"`
	}

	// UserInfo struct
	UserInfo struct {
		ID              string   `json:"user_id"`
		Name            string   `json:"name"`
		GivenName       string   `json:"given_name"`
		FamilyName      string   `json:"family_name"`
		Email           string   `json:"email"`
		Verified        bool     `json:"verified, omitempty,string"`
		Gender          string   `json:"gender, omitempty"`
		BirthDate       string   `json:"birthdate, omitempty"`
		ZoneInfo        string   `json:"zoneinfo, omitempty"`
		Locale          string   `json:"locale, omitempty"`
		Phone           string   `json:"phone_number, omitempty"`
		Address         *Address `json:"address, omitempty"`
		VerifiedAccount bool     `json:"verified_account, omitempty,string"`
		AccountType     string   `json:"account_type, omitempty"`
		AgeRange        string   `json:"age_range, omitempty"`
		PayerID         string   `json:"payer_id, omitempty"`
	}

	// WebProfile represents the configuration of the payment web payment experience
	//
	// https://developer.paypal.com/docs/api/payment-experience/
	WebProfile struct {
		ID           string       `json:"id, omitempty"`
		Name         string       `json:"name"`
		Presentation Presentation `json:"presentation, omitempty"`
		InputFields  InputFields  `json:"input_fields, omitempty"`
		FlowConfig   FlowConfig   `json:"flow_config, omitempty"`
	}

	// Presentation represents the branding and locale that a customer sees on
	// redirect payments
	//
	// https://developer.paypal.com/docs/api/payment-experience/#definition-presentation
	Presentation struct {
		BrandName  string `json:"brand_name, omitempty"`
		LogoImage  string `json:"logo_image, omitempty"`
		LocaleCode string `json:"locale_code, omitempty"`
	}

	// InputFields represents the fields that are displayed to a customer on
	// redirect payments
	//
	// https://developer.paypal.com/docs/api/payment-experience/#definition-input_fields
	InputFields struct {
		AllowNote       bool `json:"allow_note, omitempty"`
		NoShipping      uint `json:"no_shipping, omitempty"`
		AddressOverride uint `json:"address_override, omitempty"`
	}

	// FlowConfig represents the general behaviour of redirect payment pages
	//
	// https://developer.paypal.com/docs/api/payment-experience/#definition-flow_config
	FlowConfig struct {
		LandingPageType   string `json:"landing_page_type, omitempty"`
		BankTXNPendingURL string `json:"bank_txn_pending_url, omitempty"`
		UserAction        string `json:"user_action, omitempty"`
	}

	VerifyWebhookResponse struct {
		VerificationStatus string `json:"verification_status, omitempty"`
	}

	WebhookEvent struct {
		ID              string    `json:"id"`
		CreateTime      time.Time `json:"create_time"`
		ResourceType    string    `json:"resource_type"`
		EventType       string    `json:"event_type"`
		Summary         string    `json:"summary, omitempty"`
		Resource        Resource  `json:"resource"`
		Links           []Link    `json:"links"`
		EventVersion    string    `json:"event_version, omitempty"`
		ResourceVersion string    `json:"resource_version, omitempty"`
	}

	Resource struct {
		// Payment Resource type
		ID                     string                  `json:"id, omitempty"`
		Status                 string                  `json:"status, omitempty"`
		StatusDetails          *CaptureStatusDetails   `json:"status_details, omitempty"`
		Amount                 *PurchaseUnitAmount     `json:"amount, omitempty"`
		UpdateTime             string                  `json:"update_time, omitempty"`
		CreateTime             string                  `json:"create_time, omitempty"`
		ExpirationTime         string                  `json:"expiration_time, omitempty"`
		SellerProtection       *SellerProtection       `json:"seller_protection, omitempty"`
		FinalCapture           bool                    `json:"final_capture, omitempty"`
		SellerPayableBreakdown *CaptureSellerBreakdown `json:"seller_payable_breakdown, omitempty"`
		NoteToPayer            string                  `json:"note_to_payer, omitempty"`
		// merchant-onboarding Resource type
		PartnerClientID string `json:"partner_client_id, omitempty"`
		MerchantID      string `json:"merchant_id, omitempty"`
		// Common
		Links []Link `json:"links, omitempty"`
	}

	CaptureSellerBreakdown struct {
		GrossAmount         PurchaseUnitAmount  `json:"gross_amount"`
		PayPalFee           PurchaseUnitAmount  `json:"paypal_fee"`
		NetAmount           PurchaseUnitAmount  `json:"net_amount"`
		TotalRefundedAmount *PurchaseUnitAmount `json:"total_refunded_amount, omitempty"`
	}

	ReferralRequest struct {
		TrackingID            string                 `json:"tracking_id"`
		PartnerConfigOverride *PartnerConfigOverride `json:"partner_config_override,omitemtpy"`
		Operations            []Operation            `json:"operations, omitempty"`
		Products              []string               `json:"products, omitempty"`
		LegalConsents         []Consent              `json:"legal_consents, omitempty"`
	}

	PartnerConfigOverride struct {
		PartnerLogoURL       string `json:"partner_logo_url, omitempty"`
		ReturnURL            string `json:"return_url, omitempty"`
		ReturnURLDescription string `json:"return_url_description, omitempty"`
		ActionRenewalURL     string `json:"action_renewal_url, omitempty"`
		ShowAddCreditCard    *bool  `json:"show_add_credit_card, omitempty"`
	}

	Operation struct {
		Operation                string              `json:"operation"`
		APIIntegrationPreference *IntegrationDetails `json:"api_integration_preference, omitempty"`
	}

	IntegrationDetails struct {
		RestAPIIntegration *RestAPIIntegration `json:"rest_api_integration, omitempty"`
	}

	RestAPIIntegration struct {
		IntegrationMethod string            `json:"integration_method"`
		IntegrationType   string            `json:"integration_type"`
		ThirdPartyDetails ThirdPartyDetails `json:"third_party_details"`
	}

	ThirdPartyDetails struct {
		Features []string `json:"features"`
	}

	Consent struct {
		Type    string `json:"type"`
		Granted bool   `json:"granted"`
	}

	DefaultResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	// CreateProductRequest represents body parameters needed to create PayPal product
	// Type represents the product type. Indicates whether the product is physical or tangible goods, or a service. The allowed values are:
	// ---------------------------------------------------------
	// | PHYSICAL | Physical goods.							   |
	// | DIGITAL | Digital goods.							   |
	// | SERVICE | A service. For example, technical support. |
	// ---------------------------------------------------------
	// You can see category allowed values in PayPal docs -> https://developer.paypal.com/docs/api/catalog-products/v1/#products-create-request-body
	CreateProductRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Type        string `json:"type"` //default: PHYSICAL
		Category    string `json:"category"`
		ImageUrl    string `json:"image_url"`
		HomeUrl     string `json:"home_url"`
	}

	// Product represents PayPal product
	// Type represents the product type. Indicates whether the product is physical or tangible goods, or a service. The allowed values are:
	// ---------------------------------------------------------
	// | PHYSICAL | Physical goods.							   |
	// | PHYSICAL | Digital goods.							   |
	// | PHYSICAL | A service. For example, technical support. |
	// ---------------------------------------------------------
	// You can see category allowed values in PayPal docs -> https://developer.paypal.com/docs/api/catalog-products/v1/#products-create-request-body
	Product struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Type        string  `json:"type"` //default: PHYSICAL
		Category    string  `json:"category"`
		ImageUrl    string  `json:"image_url"`
		HomeUrl     string  `json:"home_url"`
		CreateTime  string  `json:"creat_time, omitempty"` //Read only
		UpdateTime  string  `json:"updat_time, omitempty"` //Read only
		Links       []*Link `json:"links, omitempty"`      //Read only
	}

	// ListProductsRequest represents query params for list products call
	ListProductsRequest struct {
		PageSize      uint64 `json:"page_size"`      //default: 10 min:1 max:20
		Page          uint64 `json:"page"`           //default: 1 min:1 max:100000
		TotalRequired bool   `json:"total_required"` //default: false
	}

	// ListProductsResponse represents the response od list products
	ListProductsResponse struct {
		TotalItems uint64     `json:"total_items, omitempty"` //min: 0, max: 500000000
		TotalPages uint64     `json:"total_pages, omitempty"` //min: 0, max: 100000000
		Products   []*Product `json:"products"`
		Links      []*Link    `json:"links"` //Read only
	}

	// PatchObject represents the object used for updating PayPal objects
	PatchObject struct {
		Operation string `json:"op"`
		Path      string `json:"path"`
		Value     string `json:"value"`
	}

	//Resource v1 for old hooks
	SuspendBillingAgreementV1Response struct {
		ID       string             `json:"id"`
		Resource BillingAgreementV1 `json:"resource"`
	}

	BillingAgreementV1 struct {
		ID               string `json:"id"`
		PlanID           string `json:"plan_id"`
		Status           string `json:"status"`
		StatusUpdateTime string `json:"status_update_time"`
	}

	// CreatePlan represents body parameters needed to create PayPal plan
	// Status represents the initial state of the plan. Allowed input values are CREATED and ACTIVE. The allowed values are:
	// ----------------------------------------------------------------------------------------------
	// | CREATED  | The plan was created. You cannot create subscriptions for a plan in this state. |
	// | INACTIVE | The plan is inactive.															|
	// | ACTIVE   | The plan is active. You can only create subscriptions for a plan in this state. |
	// ----------------------------------------------------------------------------------------------
	CreatePlan struct {
		ProductID          string              `json:"product_id"`
		Name               string              `json:"name"`
		Status             string              `json:"status"` //default: ACTIVE
		Description        string              `json:"description, omitempty"`
		BillingCycles      []*BillingCycle     `json:"billing_cycles"`
		PaymentPreferences *PaymentPreferences `json:"payment_preferences"`
		Taxes              *Taxes              `json:"taxes, omitempty"`
		QuantitySupported  bool                `json:"quantity_supported, omitempty"`
	}

	// BillingCycle represents the cycles for billing the subscription
	// The tenure type of the billing cycle. In case of a plan having trial period, only 1 trial period is allowed per plan. The possible values are:
	// --------------------------------------
	// | REGULAR | A regular billing cycle. |
	// | TRIAL   | A trial billing cycle.   |
	// --------------------------------------
	BillingCycle struct {
		PricingScheme *PricingScheme `json:"pricing_scheme, omitempty"` //Free Trial Cycle doesn't require scheme
		Frequency     Frequency      `json:"frequency"`
		TenureType    string         `json:"tenure_type"`
		Sequence      uint64         `json:"sequence"`                //min: 0, max: 99
		TotalCycles   uint64         `json:"total_cycles, omitempty"` //default: 1, min: 0, max: 999
	}

	// PricingScheme represents the active pricing scheme for this billing cycle.
	// A free trial billing cycle does not require a pricing scheme.
	PricingScheme struct {
		Version    uint64 `json:"version, omitempty"` //Read only
		FixedPrice *Money `json:"fixed_price"`
		CreateTime string `json:"create_time, omitempty"` //Read only
		UpdateTime string `json:"update_time, omitempty"` //Read only
	}

	// Frequency represents the frequency details for this billing cycle.
	// Interval unit is the interval at which the subscription is charged or billed
	// These are the possible combinations
	// --------------------------------------
	// | Interval unit | Max Interval count |
	// --------------------------------------
	// | DAY           | 365                |
	// | WEEK          | 52                 |
	// | MONTH         | 12                 |
	// | YEAR          | 1                  |
	// --------------------------------------
	Frequency struct {
		IntervalUnit  string `json:"interval_unit"`
		IntervalCount uint64 `json:"interval_count, omitempty"`
	}

	// PaymentPreferences represents the payment preferences for a subscription
	// SetupFeeFailureAction represents the action to take on the subscription if the initial payment for the setup fails. The possible values are:
	// -------------------------------------------------------------------------------------
	// | CONTINUE | Continues the subscription if the initial payment for the setup fails. |
	// | CANCEL   | Cancels the subscription if the initial payment for the setup fails.   |
	// -------------------------------------------------------------------------------------
	PaymentPreferences struct {
		AutoBillOutstanding     bool   `json:"auto_bill_outstanding, omitempty"` //default true
		SetupFee                *Money `json:"setup_fee, omitempty"`
		SetupFeeFailureAction   string `json:"setup_fee_failure_action, omitempty"`  //default: CANCEL
		PaymentFailureThreshold uint64 `json:"payment_failure_threshold, omitempty"` //default: 0, min: 0, max: 999
	}

	// Taxes represents the tax details
	Taxes struct {
		Percentage string `json:"percentage"`
		Inclusive  bool   `json:"inclusive, omitempty"` //default: true
	}

	// Plan represents the details for the subscription plan for paying
	// Status represents the initial state of the plan. Allowed input values are CREATED and ACTIVE. The allowed values are:
	// ----------------------------------------------------------------------------------------------
	// | CREATED  | The plan was created. You cannot create subscriptions for a plan in this state. |
	// | INACTIVE | The plan is inactive.															|
	// | ACTIVE   | The plan is active. You can only create subscriptions for a plan in this state. |
	// ----------------------------------------------------------------------------------------------
	Plan struct {
		ID                 string              `json:"id"`
		ProductID          string              `json:"product_id"`
		Name               string              `json:"name"`
		Status             string              `json:"status"`
		Description        string              `json:"description, omitempty"`
		BillingCycles      []*BillingCycle     `json:"billing_cycles"`
		PaymentPreferences *PaymentPreferences `json:"payment_preferences"`
		Taxes              *Taxes              `json:"taxes, omitempty"`
		QuantitySupported  bool                `json:"quantity_supported, omitempty"`
		CreateTime         string              `json:"create_time"` //Read only
		UpdateTime         string              `json:"update_time"` //Read only
		Links              []*Link             `json:"links"`       //Read only
	}

	// ListPlansParams represents query params for list products call
	ListPlansParams struct {
		ProductID     string `json:"product_id, omitempty"`
		PageSize      uint64 `json:"page_size, omitempty"`      //default: 10, min:1, max:20
		Page          uint64 `json:"page, omitempty"`           //default: 1, min:1, max:100000
		TotalRequired bool   `json:"total_required, omitempty"` //default: false
	}

	// ListPlansResponse represents the list of the plans
	ListPlansResponse struct {
		TotalItems uint64  `json:"total_items, omitempty"` //min: 0, max: 500000000
		TotalPages uint64  `json:"total_pages, omitempty"` //min: 0, max: 100000000
		Products   []*Plan `json:"plans"`
		Links      []*Link `json:"links"` //Read only
	}

	// UpdatePricingSchemasListRequest represents an array of pricing schemes to update plans billing cycles
	UpdatePricingSchemasListRequest struct {
		PricingSchemes []*UpdatePricingSchemaRequest `json:"pricing_schemes"`
	}

	// UpdatePricingSchemaRequest represents a pricing schema to update plans billing cycle
	UpdatePricingSchemaRequest struct {
		BillingCycleSequence uint64         `json:"billing_cycle_sequence"` //min: 1, max: 99
		PricingScheme        *PricingScheme `json:"pricing_scheme"`
	}

	// PaymentMethod represents the customer and merchant payment preferences
	// Currently only PAYPAl payment is supported
	// PayeePreferred represents the merchant-preferred payment sources. The possible values are:
	// -------------------------------------------------------------------------------------------------
	// | UNRESTRICTED 				| Accepts any type of payment from the customer.				   |
	// | IMMEDIATE_PAYMENT_REQUIRED | Accepts only immediate payment from the customer. For example,   |
	// | 							| credit card, PayPal balance, or instant ACH. Ensures that at the |
	// | 							| time of capture, the payment does not have the `pending` status. |
	// -------------------------------------------------------------------------------------------------
	// Category provides context (e.g. frequency of payment (Single, Recurring) along with whether
	// (Customer is Present, Not Present) for the payment being processed. For Card and PayPal Vaulted/Billing
	// Agreement transactions, this helps specify the appropriate indicators to the networks
	// (e.g. Mastercard, Visa) which ensures compliance as well as ensure a better auth-rate.
	// For bank processing, indicates to clearing house whether the transaction is recurring or not depending
	// on the option chosen. The possible values are:
	// --------------------------------------------------------------------------------------------------------------
	// | CUSTOMER_PRESENT_SINGLE_PURCHASE | If the payments is an e-commerce payment initiated by the customer.     |
	// | 								  | Customer typically enters payment information (e.g. card number,  	 	|
	// | 								  | approves payment within the PayPal Checkout flow) and such information  |
	// | 								  | that has not been previously stored on file.							|
	// | CUSTOMER_NOT_PRESENT_RECURRING   | Subsequent recurring payments (e.g. subscriptions with a fixed amount	|
	// |				 				  | on a predefined schedule when customer is not present.					|
	// | CUSTOMER_PRESENT_RECURRING_FIRST | The first payment initiated by the customer which is expected to be		|
	// | 								  | followed by a series of subsequent recurring payments transactions		|
	// | 								  | (e.g. subscriptions with a fixed amount on a predefined schedule when	|
	// | 								  | customer is not present. This is typically used for scenarios where		|
	// | 								  | customer stores credentials and makes a purchase on a given date and	|
	// |                                  | also set’s up a subscription.											|
	// | CUSTOMER_PRESENT_UNSCHEDULED     | Also known as (card-on-file) transactions. Payment details are stored	|
	// | 								  | to enable checkout with one-click, or simply to streamline the checkout |
	// | 								  | process. For card transaction customers typically do not enter a CVC	|
	// | 								  | number as part of the Checkout process.									|
	// | CUSTOMER_NOT_PRESENT_UNSCHEDULED | Unscheduled payments that are not recurring on a predefined schedule	|
	// | 								  | (e.g. balance top-up).													|
	// | MAIL_ORDER_TELEPHONE_ORDER       | Payments that are initiated by the customer via the merchant by mail 	|
	// | 								  | or telephone.															|
	// --------------------------------------------------------------------------------------------------------------
	PaymentMethod struct {
		PayerSelected  string `json:"payer_selected, omitempty"`  //default: PAYPAL
		PayeePreferred string `json:"payee_preferred, omitempty"` //default: UNRESTRICTED
		Category       string `json:"category, omitempty"`        //default: CUSTOMER_PRESENT_SINGLE_PURCHASE
	}

	// CreateSubscriptionRequest represents body parameters needed to create PayPal subscription
	CreateSubscriptionRequest struct {
		PlanID             string              `json:"plan_id"`
		StartTime          string              `json:"start_time, omitempty"` //default: current time
		Quantity           string              `json:"quantity, omitempty"`
		ShippingAmount     *Money              `json:"shipping_amount, omitempty"`
		Subscriber         *SubscriberRequest  `json:"subscriber, omitempty"`
		AutoRenewal        bool                `json:"auto_renewal, omitempty"`
		ApplicationContext *ApplicationContext `json:"application_context, omitempty"`
	}

	// SubscriberRequest represents the subscriber details
	SubscriberRequest struct {
		Name            *PayerName      `json:"name, omitempty, omitempty"`
		EmailAddress    string          `json:"email_address, omitempty"`
		PayerID         string          `json:"payer_id, omitempty"` //Read only
		ShippingAddress *ShippingDetail `json:"shipping_address, omitempty"`
		PaymentSource   *PaymentSource  `json:"payment_source, omitempty"`
	}

	// Subscriber represents the subscriber details
	Subscriber struct {
		Name            *Name                  `json:"name, omitempty, omitempty"`
		EmailAddress    string                 `json:"email_address, omitempty"`
		PayerID         string                 `json:"payer_id, omitempty"` //Read only
		ShippingAddress *ShippingDetail        `json:"shipping_address, omitempty"`
		PaymentSource   *PaymentSourceResponse `json:"payment_source, omitempty"`
	}

	// Name represents payer details
	Name struct {
		Prefix            string `json:"prefix, omitempty"`
		GivenName         string `json:"given_name, omitempty"`
		Surname           string `json:"surname, omitempty"`
		MiddleName        string `json:"middle_name, omitempty"`
		Suffix            string `json:"suffix, omitempty"`
		AlternateFullName string `json:"alternate_full_name, omitempty"`
		FullName          string `json:"full_name, omitempty"`
	}

	// PaymentSourceResponse represents the payment source definitions
	PaymentSourceResponse struct {
		Card *CardResponseWithBillingAddress `json:"card"`
	}

	// CardResponseWithBillingAddress represents card details
	// Brand represents the card brand or network. Typically used in the response. The possible values are:
	// ------------------------------------------------------
	// | VISA			 | Visa card. 						|
	// | MASTERCARD		 | Mastecard card. 					|
	// | DISCOVER		 | Discover card. 					|
	// | AMEX			 | American Express card. 			|
	// | SOLO			 | Solo debit card. 				|
	// | JCB			 | Japan Credit Bureau card. 		|
	// | STAR			 | Military Star card. 				|
	// | DELTA			 | Delta Airlines card. 			|
	// | SWITCH			 | Switch credit card. 				|
	// | MAESTRO 		 | Maestro credit card. 			|
	// | CB_NATIONALE	 | Carte Bancaire (CB) credit card. |
	// | CONFIGOGA		 | Configoga credit card. 			|
	// | CONFIDIS 		 | Confidis credit card. 			|
	// | ELECTRON 		 | Visa Electron credit card.		|
	// | CETELEM 		 | Cetelem credit card. 			|
	// | CHINA_UNION_PAY | China union pay credit card. 	|
	// ------------------------------------------------------
	// Type represents the payment card type. The possible values are:
	// ---------------------------------------------
	// | CREDIT  | A credit card. 				   |
	// | DEBIT   | A debit card. 				   |
	// | PREPAID | A Prepaid card. 				   |
	// | UNKNOWN | Card type cannot be determined. |
	// ---------------------------------------------
	CardResponseWithBillingAddress struct {
		LastDigit      string           `json:"last_digit, omitempty"` //Read only
		Brand          string           `json:"brand, omitempty"`      //Read only
		Type           string           `json:"type, omitempty"`       //Read only
		Name           string           `json:"name, omitempty"`
		BillingAddress *AddressPortable `json:"billing_address, omitempty"`
	}

	// PayerName represents payer name details
	PayerName struct {
		GivenName string `json:"given_name, omitempty"`
		Surname   string `json:"surname, omitempty"`
	}

	// Subscription represents the subscription details
	Subscription struct {
		ID               string                   `json:"id, omitempty"`
		Status           string                   `json:"status, omitempty"`
		StatusChangeNote string                   `json:"status_change_note, omitempty"`
		StatusUpdateTime string                   `json:"status_update_time, omitempty"`
		PlanID           string                   `json:"plan_id, omitempty"`
		StartTime        string                   `json:"start_time, omitempty"`
		Quantity         string                   `json:"quantity, omitempty"`
		ShippingAmount   *Money                   `json:"shipping_amount, omitempty"`
		Subscriber       *Subscriber              `json:"subscriber, omitempty"`
		BillingInfo      *SubscriptionBillingInfo `json:"billing_info, omitempty"` //Read only
		CreateTime       string                   `json:"created_time"`            //Read only
		UpdateTime       string                   `json:"update_time"`             //Read only
		Links            []*Link                  `json:"links"`                   //Read only
	}

	// SubscriptionBillingInfo represents billing details for subscription
	// The number of consecutive payment failures. Resets to 0 after a successful payment. If this reaches the payment_failure_threshold value,
	// the subscription updates to the SUSPENDED state.
	SubscriptionBillingInfo struct {
		OutstandingBalance  *Money               `json:"outstanding_balance"`
		CycleExecutions     []*CycleExecution    `json:"cycle_executions, omitempty"` //Read only
		LastPayment         LastPaymentDetails   `json:"last_payment, omitempty"`     //Read only
		NextBillingTime     string               `json:"next_billing_time"`           //Read only
		FinalPaymentTime    string               `json:"final_payment_time"`          //Read only
		FailedPaymentsCount uint64               `json:"failed_payments_count"`       //min: 0, max: 999
		LastFailedPayment   FailedPaymentDetails `json:"last_failed_payment"`         //Read only
	}

	// CycleExecution represents details about billing cycles executions
	// The number of times this billing cycle runs. Trial billing cycles can only have a value of 1 for total_cycles. Regular billing cycles
	// can either have infinite cycles (value of 0 for total_cycles) or a finite number of cycles (value between 1 and 999 for total_cycles).
	// TenureType represents the type of the billing cycle. The possible values are:
	// --------------------------------------
	// | REGULAR | A regular billing cycle. |
	// | TRIAL   | A trial billing cycle.   |
	// --------------------------------------
	CycleExecution struct {
		TenureType                  string `json:"tenure_type"`                               //Read only
		Sequence                    uint64 `json:"sequence"`                                  //min: 0, max: 99
		CyclesCompleted             uint64 `json:"cycles_completed"`                          //min: 0, max: 9999 Read only
		CyclesRemaining             uint64 `json:"cycles_remaining, omitempty"`               //min: 0, max: 9999 Read only
		CurrentPricingSchemeVersion uint64 `json:"current_pricing_scheme_version, omitempty"` //min: 0, max: 99 Read only
		TotalCycles                 uint64 `json:"total_cycles, omitempty"`                   //min: 0, max: 999 Read only
	}

	// LastPaymentDetails represents details for the last payment
	LastPaymentDetails struct {
		Amount *Money `json:"amount"` //Read only
		Time   string `json:"time"`   //Read only
	}

	// FailedPaymentDetails represents details about failed payment
	// ReasonCode represents the reason code for the payment failure. The possible values are:
	// -----------------------------------------------------------------------------------------------------------------
	// | PAYMENT_DENIED 		 			  | PayPal declined the payment due to one or more customer issues.        |
	// | INTERNAL_SERVER_ERROR   			  | An internal server error has occurred.								   |
	// | PAYEE_ACCOUNT_RESTRICTED			  | The payee account is not in good standing and cannot receive payments. |
	// | PAYER_ACCOUNT_RESTRICTED			  | The payer account is not in good standing and cannot make payments.	   |
	// | PAYER_CANNOT_PAY        			  | Payer cannot pay for this transaction. 								   |
	// | SENDING_LIMIT_EXCEEDED  			  | The transaction exceeds the payer's sending limit.                     |
	// | TRANSACTION_RECEIVING_LIMIT_EXCEEDED | The transaction exceeds the receiver's receiving limit.				   |
	// | CURRENCY_MISMATCH                    | The transaction is declined due to a currency mismatch.				   |
	// -----------------------------------------------------------------------------------------------------------------
	FailedPaymentDetails struct {
		Amount               *Money `json:"amount"`                             //Read only
		Time                 string `json:"time"`                               //Read only
		ReasonCode           string `json:"reason_code, omitempty"`             //Read only
		NextPaymentRetryTime string `json:"next_payment_retry_time, omitempty"` //Read only
	}

	// ShowSubscriptionRequest represents query parameters for show subscription call
	// Fields represents list of fields that are to be returned in the response.
	// Possible value for fields is last_failed _payment.
	ShowSubscriptionRequest struct {
		Fields string `json:"fields, omitempty"`
	}

	// UpdateSubscriptionStatusRequest represents body parameters for activate subscription
	UpdateSubscriptionStatusRequest struct {
		Reason string `json:"reason"`
	}

	// CaptureAuthorizedPaymentOnSubscriptionRequest represents body parameter for capturing authorized payment on subscription
	// CaptureType represents the type of capture. The allowed values are:
	// ---------------------------------------------------------------------------------
	// | OUTSTANDING_BALANCE | The outstanding balance that the subscriber must clear. |
	// ---------------------------------------------------------------------------------
	CaptureAuthorizedPaymentOnSubscriptionRequest struct {
		Note        string `json:"note"`
		CaptureType string `json:"capture_type"`
		Amount      *Money `json:"amount"`
	}

	// ListTransactionsForSubscriptionRequest represents query parameters for list transactions for subscription
	ListTransactionsForSubscriptionRequest struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	//TransactionsList represents list of transactions
	TransactionsList struct {
		Transactions []*Transaction `json:"transactions, omitempty"`
		TotalItems   uint64         `json:"total_items, omitempty"`
		TotalPages   uint64         `json:"total_pages, omitempty"`
		Links        []*Link        `json:"links, omitempty"` //Read only
	}

	// ReviseSubscriptionRequest represents body parameters for revise subscription
	// (update quantity of product or service in subscription)
	ReviseSubscriptionRequest struct {
		PlanID             string              `json:"plan_id, omitempty"`
		Quantity           string              `json:"quantity, omitempty"`
		ShippingAmount     *Money              `json:"shipping_amount, omitempty"`
		ShippingAddress    *ShippingDetail     `json:"shipping_address, omitempty"`
		ApplicationContext *ApplicationContext `json:"application_context, omitempty"`
	}

	// ReviseSubscriptionResponse represents response from revise subscription
	// (update quantity of product or service in subscription)
	ReviseSubscriptionResponse struct {
		PlanID          string          `json:"plan_id, omitempty"`
		Quantity        string          `json:"quantity, omitempty"`
		EffectiveTime   string          `json:"effective_time, omitempty"` //Read only
		ShippingAmount  *Money          `json:"shipping_amount, omitempty"`
		ShippingAddress *ShippingDetail `json:"shipping_address, omitempty"`
		Links           []*Link         `json:"links, omitempty"` //Read only
	}
)

// Error method implementation for ErrorResponse struct
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %s, %+v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message, r.Details)
}

// MarshalJSON for JSONTime
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, time.Time(t).UTC().Format(time.RFC3339))
	return []byte(stamp), nil
}

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	*e = expirationTime(i)
	return nil
}
