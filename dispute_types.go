package paypal

// https://developer.paypal.com/docs/api/customer-disputes/v1/#definition-dispute_info
type DisputeState string

const (
	DisputeStateOpenInquiries            DisputeState = "OPEN_INQUIRIES"              //The dispute is open.
	DisputeStateRequiredAction           DisputeState = "REQUIRED_ACTION"             //The dispute is waiting for a response.
	DisputeStateRequiredOtherPartyAction DisputeState = "REQUIRED_OTHER_PARTY_ACTION" //The dispute is waiting for a response from other party.
	DisputeStateUnderPaypalReview        DisputeState = "UNDER_PAYPAL_REVIEW"         //The dispute is under review with PayPal.
	DisputeStateAppealable               DisputeState = "APPEALABLE"                  //The dispute can be appealed.
	DisputeStateResolved                 DisputeState = "RESOLVED"                    //The dispute is resolved.
)

type AcceptClaimReason string

const (
	AcceptClaimReasonDidNotShipItem   AcceptClaimReason = "DID_NOT_SHIP_ITEM"  //Merchant is accepting customer's claim as they could not ship the item back to the customer
	AcceptClaimReasonTooTimeConsuming AcceptClaimReason = "TOO_TIME_CONSUMING" //Merchant is accepting customer's claim as it is taking too long for merchant to fulfil the order
	AcceptClaimReasonLostInMail       AcceptClaimReason = "LOST_IN_MAIL"       //Merchant is accepting customer's claim as the item is lost in mail or transit
	AcceptClaimReasonNotAbleToWin     AcceptClaimReason = "NOT_ABLE_TO_WIN"    //Merchant is accepting customer's claim as the merchant is not able to find sufficient evidence to win this dispute
	AcceptClaimReasonCompanyPolicy    AcceptClaimReason = "COMPANY_POLICY"     //Merchant is accepting customer’s claims to follow their internal company policy
	AcceptClaimReasonReasonNotSet     AcceptClaimReason = "REASON_NOT_SET"     //This is the default value merchant can use if none of the above reasons apply
)

type AcceptClaimType string

const (
	AcceptClaimTypeRefund                        AcceptClaimType = "REFUND"                            //	The merchant must refund the customer without any item replacement or return. This type is applicable when a merchant is willing to refund the entire dispute amount without any further action from customer. Omit the refund_amount and return_shipping_address parameters from the accept claim call.
	AcceptClaimTypeRefundWithReturn              AcceptClaimType = "REFUND_WITH_RETURN"                //	The customer must return the item to the merchant and then merchant will refund the money. This type is applicable when a merchant is willing to refund the dispute amount and requires the customer to return the item. Include the return_shipping_address parameter in but omit the refund_amount parameter from the accept claim call.
	AcceptClaimTypePartialRefund                 AcceptClaimType = "PARTIAL_REFUND"                    //	The merchant proposes a partial refund for the dispute.This type is applicable when a merchant is willing to refund an amount lesser than dispute amount. Include the refund_amount parameter.
	AcceptClaimTypeRefundWithReturnShipmentLabel AcceptClaimType = "REFUND_WITH_RETURN_SHIPMENT_LABEL" //	The customer must return the item to the merchant and then merchant will refund the money. This type is applicable when a merchant is willing to refund the dispute amount and requires the customer to return the item using the shipment label provided by the merchant. Include the return_shipment_info and return_shipping_address parameter in but omit the refund_amount parameter from the accept claim call.
)

type AdjudicationOutcome string

const (
	AdjudicationOutcomeBuyerFavor  AdjudicationOutcome = "BUYER_FAVOR"  //	Resolves the case in the customer's favor. Outcome is set to RESOLVED_BUYER_FAVOR.
	AdjudicationOutcomeSellerFavor AdjudicationOutcome = "SELLER_FAVOR" //	Resolves the case in the merchant's favor. Outcome is set to RESOLVED_SELLER_FAVOR.
)

type DisputeStatusAction string

const (
	DisputeStatusBuyerEvidence  DisputeStatusAction = "BUYER_EVIDENCE"  //	Changes the status of the dispute to WAITING_FOR_BUYER_RESPONSE.
	DisputeStatusSellerEvidence DisputeStatusAction = "SELLER_EVIDENCE" //	Changes the status of the dispute to WAITING_FOR_SELLER_RESPONSE.
)

type DisputeStatus string

const (
	DisputeStatusOpen                     DisputeStatus = "OPEN"                        //	The dispute is open.
	DisputeStatusWaitingForBuyerResponse  DisputeStatus = "WAITING_FOR_BUYER_RESPONSE"  //	The dispute is waiting for a response from the customer.
	DisputeStatusWaitingForSellerResponse DisputeStatus = "WAITING_FOR_SELLER_RESPONSE" //	The dispute is waiting for a response from the merchant.
	DisputeStatusUnderReview              DisputeStatus = "UNDER_REVIEW"                //	The dispute is under review with PayPal.
	DisputeStatusResolved                 DisputeStatus = "RESOLVED"                    //	The dispute is resolved.
	DisputeStatusOther                    DisputeStatus = "OTHER"                       //	The default status if the dispute does not have one of the other statuses.
)

type MakeOfferType string

const (
	MakeOfferTypeRefund                   MakeOfferType = "REFUND"                     //	The merchant must refund the customer without any item replacement or return. This offer type is valid in the inquiry phase and occurs when a merchant is willing to refund a specific amount. Buyer acceptance is needed for partial refund offers and dispute is auto closed for full refunds. Include the offer_amount but omit the return_shipping_address parameters from the make offer request.
	MakeOfferTypeRefundWithReturn         MakeOfferType = "REFUND_WITH_RETURN"         //	The customer must return the item to the merchant and then merchant will refund the money. This offer type is valid in the inquiry phase and occurs when a merchant is willing to refund a specific amount and requires the customer to return the item. Include the return_shipping_address parameter and the offer_amount parameter in the make offer request.
	MakeOfferTypeRefundWithReplacement    MakeOfferType = "REFUND_WITH_REPLACEMENT"    //	The merchant must do a refund and then send a replacement item to the customer. This offer type is valid in the inquiry phase when a merchant is willing to refund a specific amount and send the replacement item. Include the offer_amount parameter in the make offer request.
	MakeOfferTypeReplacementWithoutRefund MakeOfferType = "REPLACEMENT_WITHOUT_REFUND" //	The merchant must send a replacement item to the customer with no additional refunds. This offer type is valid in the inquiry phase when a merchant is willing to replace the item without any refund. Omit the offer_amount parameter from the make offer request.
)

type AcknowledgementType string

const (
	AcknowledgementTypeItemReceived            AcknowledgementType = "ITEM_RECEIVED"              //	The merchant must refund the customer without any item replacement or return. This offer type is valid in the inquiry phase and occurs when a merchant is willing to refund a specific amount. Buyer acceptance is needed for partial refund offers and dispute is auto closed for full refunds. Include the offer_amount but omit the return_shipping_address parameters from the make offer request.
	AcknowledgementTypeItemNotReceived         AcknowledgementType = "ITEM_NOT_RECEIVED"          //	The merchant has not received the item.
	AcknowledgementTypeDamaged                 AcknowledgementType = "DAMAGED"                    //	The items returned by the customer were damaged.
	AcknowledgementTypeEmptyPackageOrDifferent AcknowledgementType = "EMPTY_PACKAGE_OR_DIFFERENT" //	The package was empty or the goods were different from what was expected.
	AcknowledgementTypeMissingItems            AcknowledgementType = "MISSING_ITEMS"              //	The package did not have all the items that were expected.
)

type EvidenceType string

const (
	EvidenceTypeProofOfDamage                                   EvidenceType = "PROOF_OF_DAMAGE"                                       //	Documentation supporting the claim that the item is damaged.
	EvidenceTypeThirdpartyProofForDamageOrSignificantDifference EvidenceType = "THIRDPARTY_PROOF_FOR_DAMAGE_OR_SIGNIFICANT_DIFFERENCE" //	Proof should be provided by an unbiased third-party, such as a dealer, appraiser or another individual or organisation that's qualified in the area of the item in question (other than yourself), and detail the extent of the damage or clearly explain how the item received significantly differs from the item advertised.
	EvidenceTypeDelcaration                                     EvidenceType = "DECLARATION"                                           //	Signed declaration about the information provided.
	EvidenceTypeProofOfMissingItems                             EvidenceType = "PROOF_OF_MISSING_ITEMS"                                //	Image of open box with returned items and shipping label clearly visible.
	EvidenceTypeProofOfEmptyPackageOrDifferentItem              EvidenceType = "PROOF_OF_EMPTY_PACKAGE_OR_DIFFERENT_ITEM"              //	Image of empty box or returned items that are different from what were expected and shipping label clearly visible.
	EvidenceTypeProofOfItemNotReceived                          EvidenceType = "PROOF_OF_ITEM_NOT_RECEIVED"                            //	Any proof about the non receipt of the item, such as screenshot of tracking info.
	EvidenceTypeProofOfFulfillment                              EvidenceType = "PROOF_OF_FULFILLMENT"                                  //	Proof of fulfillment should be a copy of the actual shipping label on the package that shows the destination address and the shipping company's stamp to verify the shipment date.
	EvidenceTypeProofOfRefund                                   EvidenceType = "PROOF_OF_REFUND"                                       //	Proof of refund issued to the buyer
	EvidenceTypeProofOfDeliverySignature                        EvidenceType = "PROOF_OF_DELIVERY_SIGNATURE"                           //	Proof of delivery signature.
	EvidenceTypeProofOfReceiptCopy                              EvidenceType = "PROOF_OF_RECEIPT_COPY"                                 //	Copy of original receipt or invoice.
	EvidenceTypeReturnPolicy                                    EvidenceType = "RETURN_POLICY"                                         //	Copy of terms and conditions,contract or store return policy
	EvidenceTypeBillingAgreement                                EvidenceType = "BILLING_AGREEMENT"                                     //	Copy of billing agreement.
	EvidenceTypeProofOfReshipment                               EvidenceType = "PROOF_OF_RESHIPMENT"                                   //	Proof of reshipment should be a copy of the actual shipping label on the package that shows the destination address and the shipping company's stamp to verify the reshipment date.
	EvidenceTypeItemDescription                                 EvidenceType = "ITEM_DESCRIPTION"                                      //	A copy of the original description of the item or service
	EvidenceTypePoliceReport                                    EvidenceType = "POLICE_REPORT"                                         //	Copy of the police report filed.
	EvidenceTypeAffidavit                                       EvidenceType = "AFFIDAVIT"                                             //	More information has to be provided about the claim using the affidavit.
	EvidenceTypePaidWithOtherMethod                             EvidenceType = "PAID_WITH_OTHER_METHOD"                                //	Document showing item/service was paid by another payment method.
	EvidenceTypeCopyOfContract                                  EvidenceType = "COPY_OF_CONTRACT"                                      //	Copy of contract if applicable.
	EvidenceTypeTerminalAtmReceipt                              EvidenceType = "TERMINAL_ATM_RECEIPT"                                  //	Copy of terminal/ATM receipt.
	EvidenceTypePriceDifferenceReason                           EvidenceType = "PRICE_DIFFERENCE_REASON"                               //	Explanation of what the price difference is related to (increased tip amount, shipping charges, taxes, etc).
	EvidenceTypeSourceConversionRate                            EvidenceType = "SOURCE_CONVERSION_RATE"                                //	Source of expected conversion rate or fee.
	EvidenceTypeBankStatement                                   EvidenceType = "BANK_STATEMENT"                                        //	Bank/Credit statement showing withdrawal transaction.
	EvidenceTypeCreditDueReason                                 EvidenceType = "CREDIT_DUE_REASON"                                     //	The credit due reason.
	EvidenceTypeRequestCreditReceipt                            EvidenceType = "REQUEST_CREDIT_RECEIPT"                                //	The request credit receipt.
	EvidenceTypeProofOfReturn                                   EvidenceType = "PROOF_OF_RETURN"                                       //	Proof of shipment or postage that shows you returned this item to your seller and should be a copy of the actual shipping label used.
	EvidenceTypeCreate                                          EvidenceType = "CREATE"                                                //	Additional evidence information during case creation.
	EvidenceTypeChangeReason                                    EvidenceType = "CHANGE_REASON"                                         //	The evidence related to the reason change.
	EvidenceTypeProofOfRefundOutsidePaypa                       EvidenceType = "PROOF_OF_REFUND_OUTSIDE_PAYPA"                         //L	Document should show that the seller issued a refund outside Paypal.
	EvidenceTypeReceiptOfMerchandise                            EvidenceType = "RECEIPT_OF_MERCHANDISE"                                //	Check with buyer if item Delivered (seller provided Proof of Shipping)
	EvidenceTypeCustomsDocument                                 EvidenceType = "CUSTOMS_DOCUMENT"                                      //	Document confirming that the item has been confiscated.
	EvidenceTypeCustomsFeeReceipt                               EvidenceType = "CUSTOMS_FEE_RECEIPT"                                   //	Custom fees receipt paid by the buyer
	EvidenceTypeInformationOnResolution                         EvidenceType = "INFORMATION_ON_RESOLUTION"                             //	Any resolution reached with the seller should be communicated to PayPal.
	EvidenceTypeAdditionalInformationOfItem                     EvidenceType = "ADDITIONAL_INFORMATION_OF_ITEM"                        //	Any additional information of the item purchased.
	EvidenceTypeDetailsOfPurchase                               EvidenceType = "DETAILS_OF_PURCHASE"                                   //	Specific details of a purchase made under a particular transaction has to be given.
	EvidenceTypeProofOfSignificantDifference                    EvidenceType = "PROOF_OF_SIGNIFICANT_DIFFERENCE"                       //	More information required on how the item was damaged or was significantly different from the item advertised.
	EvidenceTypeProofOfSoftwareOrServiceNotAsDescribed          EvidenceType = "PROOF_OF_SOFTWARE_OR_SERVICE_NOT_AS_DESCRIBED"         //	Any screenshot or download/usage log showing that the software or service was unavailable or non-functional.
	EvidenceTypeProofOfConfiscation                             EvidenceType = "PROOF_OF_CONFISCATION	Documentation"                   // from a third party or organization that evaluated this item that confirms they confiscated it.
	EvidenceTypeCopyOfLawEnforcementAgencyReport                EvidenceType = "COPY_OF_LAW_ENFORCEMENT_AGENCY_REPORT"                 //	Report filed with a law enforcement agency or government organization. Examples of such agencies are - Internet Crime Complaint Center (www.ic3.gov), state Consumer Protection office, state police or a Federal law enforcement agency such as the FBI or Postal Inspection Service.
	EvidenceTypeAdditionalProofOfShipment                       EvidenceType = "ADDITIONAL_PROOF_OF_SHIPMENT"                          //	Additional proof of shipment such as a packing list, detailed invoice, or shipping manifest to confirm that all items have been shipped. Include carrier name and tracking number if available.
	EvidenceTypeProofOfDenialByCarrier                          EvidenceType = "PROOF_OF_DENIAL_BY_CARRIER"                            //	Documentation from the carrier should confirm the reason why they refuse to ship the item in question and the extent of the original damage.
	EvidenceTypeValidSupportingDocument                         EvidenceType = "VALID_SUPPORTING_DOCUMENT"                             //	The document you have provided doesn't support your claim that the item is Significantly Not as Described. Please provide a document to clearly show how the item received significantly differs from the item advertised.
	EvidenceTypeLegibleSupportingDocument                       EvidenceType = "LEGIBLE_SUPPORTING_DOCUMENT"                           //	The document you have provided is illegible, unclear, or too dark to read. Please provide a document that is legible and clear to read.
	EvidenceTypeReturnTrackingInformation                       EvidenceType = "RETURN_TRACKING_INFORMATION"                           //	Online tracking information for remaining items that have to be shipped to the seller.
	EvidenceTypeDeliveryReceipt                                 EvidenceType = "DELIVERY_RECEIPT"                                      //	Confirmation that the item has been received.
	EvidenceTypeProofOfInstoreReceipt                           EvidenceType = "PROOF_OF_INSTORE_RECEIPT"                              //	In-store receipt or online verification should clearly show that the buyer picked up the item.
	EvidenceTypeAdditionalTrackingInformation                   EvidenceType = "ADDITIONAL_TRACKING_INFORMATION"                       //	Tracking information should include the carrier name, online tracking number and the website where the shipment can be tracked.
	EvidenceTypeProofOfShipmentPostage                          EvidenceType = "PROOF_OF_SHIPMENT_POSTAGE"                             //	Proof of shipment or postage should be a copy of the actual shipping label on the package that shows the destination address and the carrier's stamp to verify the shipment date.
	EvidenceTypeOnlineTrackingInformation                       EvidenceType = "ONLINE_TRACKING_INFORMATION"                           //	Online tracking information to confirm delivery of item.
	EvidenceTypeProofOfInstoreRefund                            EvidenceType = "PROOF_OF_INSTORE_REFUND"                               //	Proof should be an in-store refund receipt or company documentation that clearly shows a completed refund for the transaction.
	EvidenceTypeProofForSoftwareOrServiceDelivered              EvidenceType = "PROOF_FOR_SOFTWARE_OR_SERVICE_DELIVERED"               //	Proof should be compelling evidence to prove that the item or service was as described and was delivered to the buyer. Include information such as transaction ID, invoice ID, name or email to associate the evidence with the buyer along with date and time of the delivery.
	EvidenceTypeReturnAddressForShipping                        EvidenceType = "RETURN_ADDRESS_FOR_SHIPPING"                           //	Return address is required for the buyer to ship the merchandise back to the seller.
	EvidenceTypeCopyOfTheEparcelManifest                        EvidenceType = "COPY_OF_THE_EPARCEL_MANIFEST"                          //	To validate a claim, a copy of the eparcel manifest showing the buyer's address from Australia Post is required.
	EvidenceTypeCopyOfShippingManifest                          EvidenceType = "COPY_OF_SHIPPING_MANIFEST"                             //	The shipping manifest must show the buyer's address and can be obtained from the carrier.
	EvidenceTypeAppealAffidavit                                 EvidenceType = "APPEAL_AFFIDAVIT"                                      //	Appeal affidavit is needed to make an appeal for any case outcome.
	EvidenceTypeReceiptOfReplacement                            EvidenceType = "RECEIPT_OF_REPLACEMENT"                                //	Check with buyer if the replacement of the item sent by the seller was received
	EvidenceTypeCopyOfDriversLicense                            EvidenceType = "COPY_OF_DRIVERS_LICENSE"                               //	Need Copy of Drivers license.
	EvidenceTypeAccountChangeInformation                        EvidenceType = "ACCOUNT_CHANGE_INFORMATION"                            //	Additional Details about how account was accessed/what was changed.
	EvidenceTypeDeliveryAddress                                 EvidenceType = "DELIVERY_ADDRESS"                                      //	Address where item was supposed to be delivered.
	EvidenceTypeConfirmationOfResolution                        EvidenceType = "CONFIRMATION_OF_RESOLUTION"                            //	Confirmation that item was received and issue resolved.
	EvidenceTypeMerchantResponse                                EvidenceType = "MERCHANT_RESPONSE"                                     //	Copy of merchant's response when the resolution was attempted.
	EvidenceTypePermissionDescription                           EvidenceType = "PERMISSION_DESCRIPTION"                                //	A Detailed description about the account or card level permission given to another person.
	EvidenceTypeStatusOfMerchandise                             EvidenceType = "STATUS_OF_MERCHANDISE"                                 //	Details of the merchandise's current location.
	EvidenceTypeLostCardDetails                                 EvidenceType = "LOST_CARD_DETAILS"                                     //	Details of where and when the card was lost/stolen?.
	EvidenceTypeLastValidTransactionDetails                     EvidenceType = "LAST_VALID_TRANSACTION_DETAILS"                        //	Details of the last valid transaction made on the card.
	EvidenceTypeAdditionalProofOfReturn                         EvidenceType = "ADDITIONAL_PROOF_OF_RETURN"                            //	Document to confirm that the item to be returned to the seller has been shipped.
	EvidenceTypeOrderDetails                                    EvidenceType = "ORDER_DETAILS"                                         //	Order details or photos of the item/service/booking/digital download available on the website.
	EvidenceTypeListingUrl                                      EvidenceType = "LISTING_URL"                                           //	The website URLs of the item/service/booking/digital download.
	EvidenceTypeShippingInsurance                               EvidenceType = "SHIPPING_INSURANCE"                                    //	Insurance information of the shipped item.
	EvidenceTypeBuyerResponse                                   EvidenceType = "BUYER_RESPONSE"                                        //	Copy of the buyer response or any other document showing buyer's understanding of the item/service/booking/digital download purchased.
	EvidenceTypePhotosOfShippedItem                             EvidenceType = "PHOTOS_OF_SHIPPED_ITEM"                                //	Photos of the item that were shipped to the buyer.
	EvidenceTypeOther                                           EvidenceType = "OTHER"                                                 //	Other.
	EvidenceTypeCancellationDetails                             EvidenceType = "CANCELLATION_DETAILS"                                  //	Cancellation details information, for example- cancellation date, cancellation number.
	EvidenceTypeMerchantContactDetails                          EvidenceType = "MERCHANT_CONTACT_DETAILS"                              //	Merchant contacted details information, for example- merchant contacted time, mode of contacting.
	EvidenceTypeItemDetails                                     EvidenceType = "ITEM_DETAILS"                                          //	Item related details information, for example- expected delivery date.
	EvidenceTypeCorrectRecipientInformation                     EvidenceType = "CORRECT_RECIPIENT_INFORMATION"                         //	Proof of the correct recipient name, email or phone.
	EvidenceTypeExplanationOfFundsNotDelivered                  EvidenceType = "EXPLANATION_OF_FUNDS_NOT_DELIVERED"                    //	Details on why the funds were not delivered to the correct recipient.
)
