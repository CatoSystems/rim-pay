# CLICK (BNM) Provider

CLICK is BNM's e-commerce payment solution, built on the TagPay/"Payfair"
`online.php` hosted-page system. The customer is redirected to a TagPay-hosted
payment page where they authenticate (PIN + SMS OTP, or PIN on IVR), and TagPay
notifies your server of the result.

The flow has three steps:

1. **Session** — the library requests a session ID from TagPay.
2. **Order form** — the library builds the hidden form fields your web page
   POSTs (from the customer's browser) to the TagPay payment page.
3. **Notification** — TagPay calls your Notification URL (server-to-server) with
   the transaction result, which you parse with `HandleClickNotification`.

## Configuration

```go
client, _ := rimpay.NewClient(rimpay.DefaultConfig())

err := client.AddClickProvider(rimpay.ProviderConfig{
    BaseURL: "https://tagpay.example",      // TagPay online base URL (from BNM)
    Timeout: 30 * time.Second,
    Credentials: map[string]string{
        "merchant_id": "0896353536734538",  // 16-digit merchant ID (required)
    },
    Options: map[string]interface{}{
        "brand_name": "My Shop",            // optional, shown on the TagPay page
    },
})
```

| Config | Required | Description |
|--------|----------|-------------|
| `BaseURL` | ✅ | TagPay online payment base URL |
| `Credentials["merchant_id"]` | ✅ | 16-digit merchant identifier from BNM |
| `Timeout` | ✅ | HTTP timeout (must be positive) |
| `Options["brand_name"]` | ❌ | Brand shown on the hosted payment page |

The merchant's Notify/Accept/Decline/Cancel URLs and IP whitelist are configured
administratively on the TagPay side, not by this library.

## Creating a payment

```go
resp, err := client.ProcessClickPayment(ctx, &rimpay.ClickPaymentRequest{
    PhoneNumber: phoneNum, // optional; if set, the customer can't change it
    Amount:      money.New(decimal.NewFromFloat(15.00), money.MRU),
    Reference:   "PURCHASE0987",
    Description: "Online purchase",
    SuccessURL:  "https://shop.example/ok",
    FailureURL:  "https://shop.example/fail",
    CancelURL:   "https://shop.example/cancel",
})
```

`resp.PaymentURL` is the TagPay page to POST to, and `resp.Metadata["form_data"]`
(a `url.Values`) holds the hidden form fields. Amounts are sent in cents and the
currency as the ISO 4217 numeric code (MRU = `929`). All order-form parameter
names and values are lower case, per the TagPay spec.

## Handling the notification

TagPay calls your Notification URL with a server-to-server GET. In your handler,
copy the query parameters into `ClickNotificationData`:

```go
q := r.URL.Query()
status, err := client.HandleClickNotification(&rimpay.ClickNotificationData{
    Status:      q.Get("status"),      // OK / NOK
    PurchaseRef: q.Get("purchaseref"),
    Amount:      q.Get("amount"),
    Currency:    q.Get("currency"),
    ClientID:    q.Get("clientid"),
    ClientName:  q.Get("cname"),
    Mobile:      q.Get("mobile"),
    PaymentRef:  q.Get("paymentref"),
    PayID:       q.Get("payid"),
    Timestamp:   q.Get("timestamp"),
    IPAddress:   q.Get("ipaddr"),
    Error:       q.Get("error"),
    Reason:      q.Get("reason"),
})
```

Status mapping:

| `status` | `error` | Result |
|----------|---------|--------|
| `OK` | — | `PaymentStatusSuccess` |
| `NOK` | `CANCEL` | `PaymentStatusCancelled` |
| `NOK` | `AUTHENTICATION` | `PaymentStatusFailed` |
| `NOK` | `PAYMENT_FAILED` | `PaymentStatusFailed` |
| `NOK` | `EXPIRED_SESSION` | `PaymentStatusExpired` |
| `NOK` | other/empty | `PaymentStatusFailed` |

> **Security:** verify the notified `amount`/`purchaseref` against your own order
> record, and restrict the notification endpoint to TagPay's IP address. The
> library does not enforce these.
